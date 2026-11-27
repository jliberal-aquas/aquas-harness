// Package docs revisa la documentación de capacidades de docs/: enlaces,
// cajas de diagramas sin enlace, términos del glosario usados antes de
// explicarse y páginas sin vigencia declarada.
package docs

import (
	"fmt"
	"maps"
	"path"
	"regexp"
	"slices"
	"strings"
	"unicode"
)

// Finding es un problema encontrado en una página.
type Finding struct {
	Page string // ruta relativa a la raíz, con «/»
	Line int    // 1 en adelante; 0 si es de toda la página
	Msg  string
}

// String da el formato «ruta:línea: mensaje» que entienden los editores.
func (f Finding) String() string {
	return fmt.Sprintf("%s:%d: %s", f.Page, f.Line, f.Msg)
}

// Site es lo que el núcleo necesita saber del repositorio, ya leído del disco.
type Site struct {
	Pages    map[string]string // páginas a revisar: ruta → texto
	Glossary string            // ruta del glosario, p. ej. «GLOSSARY.md»
	Markdown map[string]string // todo Markdown que se puede enlazar (incluye Pages y el glosario)
	Exists   func(p string) bool
	// Stale devuelve los cambios de sources posteriores a verified; error si verified no existe.
	Stale func(verified string, sources []string) ([]string, error)
}

// capabilityPage dice si una página pertenece a una capacidad (docs/<x>/...),
// y por tanto debe declarar su vigencia. El índice y los ADR no la declaran.
func capabilityPage(p string) bool {
	parts := strings.Split(p, "/")
	return len(parts) >= 3 && parts[0] == "docs" && parts[1] != "adr"
}

// conceptPage dice si a la página se le aplican diagramas y orden de conceptos.
func conceptPage(p string) bool {
	return !strings.HasPrefix(p, "docs/adr/")
}

// Check revisa todas las páginas y devuelve los hallazgos ordenados por página y línea.
func Check(s Site) []Finding {
	terms := glossaryTerms(s.Markdown[s.Glossary])
	var out []Finding
	for _, p := range slices.Sorted(maps.Keys(s.Pages)) {
		doc := parse(s.Pages[p])
		out = append(out, checkLinks(s, p, doc)...)
		if conceptPage(p) {
			out = append(out, checkDiagrams(p, doc)...)
			out = append(out, checkTerms(p, doc, terms)...)
		}
		if capabilityPage(p) {
			out = append(out, checkFreshness(s, p, doc)...)
		}
	}
	slices.SortStableFunc(out, func(a, b Finding) int {
		if c := strings.Compare(a.Page, b.Page); c != 0 {
			return c
		}
		return a.Line - b.Line
	})
	return out
}

// line es una línea de la página con su clase.
type line struct {
	n       int    // número de línea en el archivo
	text    string // contenido
	kind    lineKind
	blockID int // índice del bloque de código al que pertenece, si es de código
}

type lineKind int

const (
	prose   lineKind = iota // texto normal
	heading                 // título «#»
	fence                   // apertura o cierre de bloque de código
	code                    // dentro de un bloque que no es mermaid
	mermaid                 // dentro de un bloque mermaid
)

// doc es una página ya separada en cabecera y líneas clasificadas.
type doc struct {
	sources  []string
	verified string
	hasMeta  bool
	lines    []line
}

// bom es la marca de orden de bytes UTF-8 que agregan algunos editores de Windows.
const bom = "\xEF\xBB\xBF"

// parse separa la cabecera «---» y clasifica cada línea según el bloque en que está.
func parse(text string) doc {
	text = strings.ReplaceAll(strings.TrimPrefix(text, bom), "\r\n", "\n")
	raw := strings.Split(text, "\n")
	var d doc
	start := 0
	if len(raw) > 0 && raw[0] == "---" {
		for i := 1; i < len(raw); i++ {
			if raw[i] == "---" {
				d.hasMeta = true
				d.sources, d.verified = parseMeta(raw[1:i])
				start = i + 1
				break
			}
		}
	}
	inBlock, isMermaid, block := false, false, 0
	for i := start; i < len(raw); i++ {
		l := line{n: i + 1, text: raw[i]}
		trimmed := strings.TrimSpace(raw[i])
		switch {
		case strings.HasPrefix(trimmed, "```"):
			l.kind = fence
			if !inBlock {
				block++
				isMermaid = strings.TrimSpace(strings.TrimPrefix(trimmed, "```")) == "mermaid"
			}
			inBlock = !inBlock
			l.blockID = block
		case inBlock && isMermaid:
			l.kind, l.blockID = mermaid, block
		case inBlock:
			l.kind, l.blockID = code, block
		case strings.HasPrefix(trimmed, "#"):
			l.kind = heading
		default:
			l.kind = prose
		}
		d.lines = append(d.lines, l)
	}
	return d
}

// parseMeta lee el subconjunto de YAML que usan las páginas: «sources:» con
// una lista de «- ruta» y «verified: hash». Ignora cualquier otra clave.
func parseMeta(lines []string) (sources []string, verified string) {
	inSources := false
	for _, l := range lines {
		t := strings.TrimSpace(l)
		switch {
		case strings.HasPrefix(t, "sources:"):
			inSources = true
		case inSources && strings.HasPrefix(t, "- "):
			sources = append(sources, strings.TrimSpace(strings.TrimPrefix(t, "- ")))
		case strings.HasPrefix(t, "verified:"):
			inSources = false
			verified = strings.TrimSpace(strings.TrimPrefix(t, "verified:"))
		default:
			inSources = false
		}
	}
	return sources, verified
}

// linkRe encuentra enlaces Markdown «[texto](destino)»; no admite corchetes anidados.
var linkRe = regexp.MustCompile(`\[([^\[\]]*)\]\(([^()\s]*)\)`)

// boldRe encuentra negritas «**texto**».
var boldRe = regexp.MustCompile(`\*\*([^*]+)\*\*`)

// inlineCodeRe encuentra código en línea «`x`», que nunca cuenta como mención.
var inlineCodeRe = regexp.MustCompile("`[^`]*`")

// checkLinks comprueba que cada enlace relativo apunte a un archivo existente
// y, si trae ancla, a un título de ese archivo.
func checkLinks(s Site, page string, d doc) []Finding {
	var out []Finding
	for _, l := range d.lines {
		if l.kind != prose && l.kind != heading {
			continue
		}
		text := inlineCodeRe.ReplaceAllString(l.text, "")
		for _, m := range linkRe.FindAllStringSubmatch(text, -1) {
			if msg := resolveLink(s, page, m[2]); msg != "" {
				out = append(out, Finding{page, l.n, msg})
			}
		}
	}
	return out
}

// resolveLink devuelve un mensaje si el destino no existe; vacío si está bien.
// Los enlaces externos (con esquema) no se revisan.
func resolveLink(s Site, page, target string) string {
	if target == "" || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
		return ""
	}
	file, anchor, _ := strings.Cut(target, "#")
	dest := page
	if file != "" {
		dest = path.Clean(path.Join(path.Dir(page), file))
		if dest == ".." || strings.HasPrefix(dest, "../") {
			return fmt.Sprintf("enlace %q sale del repositorio", target)
		}
		if !s.Exists(dest) {
			return fmt.Sprintf("enlace roto: %q no existe", target)
		}
	}
	if anchor == "" {
		return ""
	}
	text, ok := s.Markdown[dest]
	if !ok {
		return fmt.Sprintf("enlace %q: el ancla solo se puede revisar en archivos .md", target)
	}
	if !slices.Contains(anchors(text), anchor) {
		return fmt.Sprintf("enlace %q: no hay un título con el ancla #%s", target, anchor)
	}
	return ""
}

// anchors calcula las anclas de los títulos de una página como lo hace GitHub:
// minúsculas, sin puntuación, espacios a guiones y sufijo -1, -2 si se repiten.
func anchors(text string) []string {
	var out []string
	seen := map[string]int{}
	for _, l := range parse(text).lines {
		if l.kind != heading {
			continue
		}
		slug := slugify(strings.TrimLeft(strings.TrimSpace(l.text), "#"))
		if n := seen[slug]; n > 0 {
			out = append(out, fmt.Sprintf("%s-%d", slug, n))
		} else {
			out = append(out, slug)
		}
		seen[slug]++
	}
	return out
}

// slugify convierte un título en ancla: conserva letras, números, «-» y «_».
func slugify(title string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(title)) {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' || r == '_':
			b.WriteRune(r)
		case r == ' ':
			b.WriteRune('-')
		}
	}
	return b.String()
}

// nodeRe encuentra definiciones de nodos Mermaid: id seguido de su forma.
// Grupos: 1 id, 2 apertura, 3 etiqueta. Las formas largas van primero.
var nodeRe = regexp.MustCompile(`(?:^|[\s;&>-])([A-Za-z0-9_]+)(\(\[|\[\[|\(\(|\{\{|\[/|\[|\(|\{)"?([^\]\)\}"]*)"?`)

// checkDiagrams exige que cada caja de proceso «[...]», «(...)» o «[[...]]» de
// un diagrama tenga un enlace con su misma etiqueta en la lista que lo sigue.
// Inicio y fin «([...])», decisiones «{...}», círculos «((...))», datos «[/.../]»
// y títulos de «subgraph» no lo necesitan.
func checkDiagrams(page string, d doc) []Finding {
	var out []Finding
	for i := 0; i < len(d.lines); i++ {
		l := d.lines[i]
		if l.kind != fence || i+1 >= len(d.lines) || d.lines[i+1].kind != mermaid {
			continue
		}
		boxes := map[string]int{} // etiqueta → línea
		var order []string
		j := i + 1
		for ; j < len(d.lines) && d.lines[j].kind == mermaid; j++ {
			t := strings.TrimSpace(d.lines[j].text)
			if strings.HasPrefix(t, "click ") || strings.HasPrefix(t, "classDef ") || strings.HasPrefix(t, "class ") || strings.HasPrefix(t, "style ") || strings.HasPrefix(t, "subgraph ") {
				continue
			}
			for _, m := range nodeRe.FindAllStringSubmatch(t, -1) {
				if m[2] != "[" && m[2] != "(" && m[2] != "[[" {
					continue
				}
				label := strings.TrimSpace(m[3])
				if _, ok := boxes[label]; !ok && label != "" {
					boxes[label] = d.lines[j].n
					order = append(order, label)
				}
			}
		}
		links := map[string]bool{}
		for k := j + 1; k < len(d.lines) && d.lines[k].kind == prose; k++ {
			for _, m := range linkRe.FindAllStringSubmatch(d.lines[k].text, -1) {
				links[normalize(m[1])] = true
			}
		}
		for _, label := range order {
			if !links[normalize(label)] {
				out = append(out, Finding{page, boxes[label], fmt.Sprintf("la caja %q no tiene enlace en la lista bajo el diagrama", label)})
			}
		}
		i = j
	}
	return out
}

// normalize compara etiquetas sin distinguir mayúsculas ni espacios repetidos.
func normalize(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

// glossaryTerms lee los términos del glosario: líneas «**Término**:».
func glossaryTerms(text string) []string {
	var terms []string
	for _, l := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "**") && strings.HasSuffix(t, "**:") {
			terms = append(terms, strings.TrimSuffix(strings.TrimPrefix(t, "**"), "**:"))
		}
	}
	return terms
}

// termRe arma la expresión que encuentra un término como palabra completa,
// sin distinguir mayúsculas y aceptando plural en «s» o «es».
func termRe(term string) *regexp.Regexp {
	return regexp.MustCompile(`(?i)(?:^|[^\p{L}\p{N}_])(` + regexp.QuoteMeta(term) + `(?:es|s)?)(?:[^\p{L}\p{N}_]|$)`)
}

// checkTerms exige que la primera mención de cada término del glosario en la
// página esté dentro de un enlace o en negrita. Una mención en un diagrama
// antes de esa explicación también es un hallazgo. Los títulos no cuentan.
func checkTerms(page string, d doc, terms []string) []Finding {
	var out []Finding
	res := make([]*regexp.Regexp, len(terms))
	for i, t := range terms {
		res[i] = termRe(t)
	}
	done := make([]bool, len(terms))
	for _, l := range d.lines {
		switch l.kind {
		case mermaid:
			for i, re := range res {
				if !done[i] && re.MatchString(l.text) {
					done[i] = true
					out = append(out, Finding{page, l.n, fmt.Sprintf("%q aparece en un diagrama antes de enlazarse o definirse en la prosa", terms[i])})
				}
			}
		case prose:
			text := blankInlineCode(l.text)
			introduced := spans(text)
			for i, re := range res {
				if done[i] {
					continue
				}
				loc := re.FindStringSubmatchIndex(text)
				if loc == nil {
					continue
				}
				done[i] = true
				if !inside(introduced, loc[2]) {
					out = append(out, Finding{page, l.n, fmt.Sprintf("primera mención de %q sin enlace ni negrita", terms[i])})
				}
			}
		}
	}
	return out
}

// blankInlineCode reemplaza el código en línea por espacios para conservar posiciones.
func blankInlineCode(s string) string {
	return inlineCodeRe.ReplaceAllStringFunc(s, func(m string) string { return strings.Repeat(" ", len(m)) })
}

// spans devuelve los tramos [inicio, fin) del texto de enlaces y de negritas.
func spans(s string) [][2]int {
	var out [][2]int
	for _, m := range linkRe.FindAllStringSubmatchIndex(s, -1) {
		out = append(out, [2]int{m[2], m[3]})
	}
	for _, m := range boldRe.FindAllStringSubmatchIndex(s, -1) {
		out = append(out, [2]int{m[2], m[3]})
	}
	return out
}

// inside dice si la posición cae dentro de alguno de los tramos.
func inside(spans [][2]int, pos int) bool {
	for _, sp := range spans {
		if pos >= sp[0] && pos < sp[1] {
			return true
		}
	}
	return false
}

// checkFreshness exige «sources» y «verified», que cada fuente exista y que no
// haya cambiado después del commit verificado.
func checkFreshness(s Site, page string, d doc) []Finding {
	if !d.hasMeta || len(d.sources) == 0 || d.verified == "" {
		return []Finding{{page, 1, "falta la cabecera de vigencia con «sources» y «verified»"}}
	}
	var out []Finding
	for _, src := range d.sources {
		if !s.Exists(src) {
			out = append(out, Finding{page, 1, fmt.Sprintf("la fuente %q no existe", src)})
		}
	}
	changes, err := s.Stale(d.verified, d.sources)
	if err != nil {
		return append(out, Finding{page, 1, fmt.Sprintf("no pude comprobar la vigencia desde %s: %v", d.verified, err)})
	}
	for _, c := range changes {
		out = append(out, Finding{page, 1, fmt.Sprintf("desactualizada: las fuentes cambiaron después de %s (%s)", d.verified, c)})
	}
	return out
}
