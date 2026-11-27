// Package canon compila las reglas canónicas, la configuración y las skills
// de cada capa (gobierno, producto) y las escribe donde las lee cada arnés.
package canon

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path"
	"slices"
	"strings"
)

// layer es una fuente de reglas. El orden de layers es la jerarquía:
// la primera capa se escribe primero; la última tiene la palabra final.
type layer struct {
	name string // nombre legible para los mensajes
	dir  string // carpeta relativa a la raíz
}

var layers = []layer{
	{name: "governance", dir: ".aquas-governance"},
	{name: "product", dir: ".product-governance"},
}

// kind es el tipo de contenido. Cada arnés decide en qué carpeta va cada tipo.
type kind string

const (
	kindRules  kind = "rules"
	kindConfig kind = "config"
	kindSkills kind = "skills"
)

// sections asocia cada subcarpeta de una capa con el tipo que aporta.
// merge = true une los aportes de todas las capas; merge = false copia el
// archivo tal cual y exige que una sola capa lo traiga (una skill no se
// puede pegar con otra sin romper su cabecera).
var sections = []struct {
	src   string
	kind  kind
	merge bool
}{
	{src: "canon", kind: kindRules, merge: true},
	{src: "config", kind: kindConfig, merge: true},
	{src: "skills", kind: kindSkills, merge: false},
}

// harnesses dice, por cada arnés soportado, en qué carpeta de la raíz va cada
// tipo. Soportar otro arnés es agregar una fila aquí; la lógica no cambia.
var harnesses = map[string]map[kind]string{
	"claude-code": {
		kindRules:  ".claude/rules",
		kindSkills: ".claude/skills",
		kindConfig: ".claude/aquas",
	},
}

// settingsKey es el archivo de configuración, ya unido, que declara los arneses.
const settingsKey = string(kindConfig) + "/settings.json"

// bom es la marca de orden de bytes UTF-8 que agregan algunos editores de Windows.
var bom = []byte("\xEF\xBB\xBF")

// part es el aporte de una capa a un archivo de salida.
type part struct {
	layer string // nombre de la capa
	src   string // ruta relativa a la raíz, con «/»
	data  []byte
}

// output es un archivo compilado con sus aportes en orden jerárquico.
type output struct {
	kind  kind
	rel   string // ruta dentro de la carpeta del tipo, con «/»
	merge bool
	parts []part
}

// key identifica el archivo sin depender del arnés: «<tipo>/<rel>».
func (o *output) key() string { return path.Join(string(o.kind), o.rel) }

// Run compila las capas para cada arnés declarado; con check no escribe nada.
// stale = true si algún archivo difiere; si algo no compila, no escribe ninguno.
func Run(rootPath string, check bool, w io.Writer) (stale bool, err error) {
	root, err := os.OpenRoot(rootPath)
	if err != nil {
		return false, fmt.Errorf("no pude abrir la raíz %q: %w", rootPath, err)
	}
	defer root.Close()

	outs, err := collect(root)
	if err != nil {
		return false, err
	}
	if len(outs) == 0 {
		return false, fmt.Errorf("no encontré archivos en %q: revisa -root", rootPath)
	}
	built, err := buildAll(outs)
	if err != nil {
		return false, err
	}
	names, err := harnessNames(built[settingsKey])
	if err != nil {
		return false, err
	}

	var errs []error
	for _, name := range names {
		s, err := sync(root, harnesses[name], outs, built, check, w)
		stale = stale || s
		errs = append(errs, err)
	}
	return stale, errors.Join(errs...)
}

// sync deja las carpetas de un arnés iguales a lo compilado y avisa de los huérfanos.
func sync(root *os.Root, dirs map[kind]string, outs map[string]*output, built map[string][]byte, check bool, w io.Writer) (stale bool, err error) {
	var errs []error
	owned := map[string]bool{}
	for _, k := range slices.Sorted(maps.Keys(outs)) {
		o := outs[k]
		target := path.Join(dirs[o.kind], o.rel)
		owned[target] = true
		want := built[k]

		have, err := readFile(root, target)
		switch {
		case err == nil && bytes.Equal(have, want):
			fmt.Fprintf(w, "sin cambios     %s\n", target)
			continue
		case err != nil && !errors.Is(err, fs.ErrNotExist):
			errs = append(errs, fmt.Errorf("leyendo %s: %w", target, err))
			continue
		}

		stale = true
		if check {
			fmt.Fprintf(w, "desactualizado  %s\n", target)
			continue
		}
		if err := writeFileAtomic(root, target, want); err != nil {
			errs = append(errs, fmt.Errorf("escribiendo %s: %w", target, err))
			continue
		}
		fmt.Fprintf(w, "escrito         %s  (%s)\n", target, o.layerNames())
	}

	for _, orphan := range orphans(root, dirs, owned) {
		fmt.Fprintf(w, "aviso: %s no viene de ninguna capa; revísalo o bórralo\n", orphan)
	}
	return stale, errors.Join(errs...)
}

// collect recorre las capas en orden jerárquico y agrupa los archivos por clave.
func collect(root *os.Root) (map[string]*output, error) {
	outs := map[string]*output{}
	for _, l := range layers {
		for _, s := range sections {
			base := path.Join(l.dir, s.src)
			if _, err := root.Stat(base); errors.Is(err, fs.ErrNotExist) {
				continue
			} else if err != nil {
				return nil, fmt.Errorf("leyendo %s: %w", base, err)
			}

			err := fs.WalkDir(root.FS(), base, func(p string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() || strings.HasPrefix(d.Name(), ".") {
					return err
				}
				data, err := readFile(root, p)
				if err != nil {
					return err
				}
				o := &output{kind: s.kind, rel: strings.TrimPrefix(p, base+"/"), merge: s.merge}
				if prev, ok := outs[o.key()]; ok {
					o = prev
				} else {
					outs[o.key()] = o
				}
				o.parts = append(o.parts, part{layer: l.name, src: p, data: data})
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("recorriendo %s: %w", base, err)
			}
		}
	}
	return outs, nil
}

// buildAll compila todos los archivos. Reúne todos los errores para mostrarlos juntos.
func buildAll(outs map[string]*output) (map[string][]byte, error) {
	built := make(map[string][]byte, len(outs))
	var errs []error
	for _, k := range slices.Sorted(maps.Keys(outs)) {
		data, err := outs[k].build()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		built[k] = data
	}
	return built, errors.Join(errs...)
}

// build combina los aportes según la sección y el tipo de archivo.
func (o *output) build() ([]byte, error) {
	if !o.merge {
		if len(o.parts) > 1 {
			return nil, fmt.Errorf("%s viene de varias capas (%s): debe existir en una sola", o.key(), o.layerNames())
		}
		return o.parts[0].data, nil
	}
	switch path.Ext(o.rel) {
	case ".md", ".txt":
		return concatText(o), nil
	case ".json":
		return mergeJSON(o.parts)
	}
	if len(o.parts) == 1 {
		return o.parts[0].data, nil
	}
	return nil, fmt.Errorf("%s: no sé combinar archivos %q de varias capas", o.key(), path.Ext(o.rel))
}

// layerNames devuelve las capas que aportan al archivo, en orden: «governance + product».
func (o *output) layerNames() string {
	names := make([]string, len(o.parts))
	for i, p := range o.parts {
		names[i] = p.layer
	}
	return strings.Join(names, " + ")
}

// harnessNames lee la lista «harness» de la configuración unida y comprueba
// que cada arnés tenga destino. Un arnés desconocido es error, no se salta.
func harnessNames(settings []byte) ([]string, error) {
	if settings == nil {
		return nil, fmt.Errorf("falta %s en las capas: ahí se declaran los arneses", settingsKey)
	}
	names, err := decodeHarness(settings)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", settingsKey, err)
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("%s no declara ningún arnés en «harness»", settingsKey)
	}
	for _, n := range names {
		if _, ok := harnesses[n]; !ok {
			return nil, fmt.Errorf("%s: el arnés %q no está soportado (soportados: %s)",
				settingsKey, n, strings.Join(slices.Sorted(maps.Keys(harnesses)), ", "))
		}
	}
	return names, nil
}

// concatText une los textos en orden jerárquico, separados por una línea en blanco.
// En Markdown agrega un comentario que indica las fuentes y avisa que no se edite.
func concatText(o *output) []byte {
	var b bytes.Buffer
	if path.Ext(o.rel) == ".md" {
		srcs := make([]string, len(o.parts))
		for i, p := range o.parts {
			srcs[i] = p.src
		}
		fmt.Fprintf(&b, "<!-- Generado por canonc desde %s. No editar aquí: edita las fuentes. -->\n\n", strings.Join(srcs, " + "))
	}

	wrote := false
	for _, p := range o.parts {
		text := normalizeText(p.data)
		if text == "" {
			continue
		}
		if wrote {
			b.WriteString("\n\n")
		}
		b.WriteString(text)
		wrote = true
	}
	b.WriteString("\n")
	return b.Bytes()
}

// normalizeText quita la marca BOM, usa saltos de línea LF y recorta el final,
// para que la unión no deje un BOM a mitad de archivo ni líneas mezcladas.
func normalizeText(data []byte) string {
	data = bytes.TrimPrefix(data, bom)
	s := strings.ReplaceAll(string(data), "\r\n", "\n")
	return strings.TrimRight(s, " \t\n")
}

// mergeJSON fusiona los objetos JSON de las capas en orden jerárquico
// y devuelve el resultado con formato estable.
func mergeJSON(parts []part) ([]byte, error) {
	merged := map[string]any{}
	for _, p := range parts {
		obj, err := decodeObject(bytes.TrimPrefix(p.data, bom))
		if err != nil {
			return nil, fmt.Errorf("%s: JSON inválido: %w", p.src, err)
		}
		deepMerge(merged, obj)
	}
	return encodeObject(merged)
}

// deepMerge copia src sobre dst: los objetos anidados se unen clave por clave;
// cualquier otro valor (texto, número, lista) de src reemplaza al de dst.
func deepMerge(dst, src map[string]any) {
	for k, sv := range src {
		if sm, ok := sv.(map[string]any); ok {
			if dm, ok := dst[k].(map[string]any); ok {
				deepMerge(dm, sm)
				continue
			}
		}
		dst[k] = sv
	}
}

// orphans lista, sin borrarlos, los archivos de las carpetas del compilador que
// ninguna capa produce. El resto del arnés (p. ej. .claude/settings.json) no se mira.
func orphans(root *os.Root, dirs map[kind]string, owned map[string]bool) []string {
	var found []string
	for _, dir := range slices.Sorted(maps.Values(dirs)) {
		_ = fs.WalkDir(root.FS(), dir, func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() { // carpeta ausente o ilegible: no hay huérfanos que avisar
				return nil
			}
			if !owned[p] {
				found = append(found, p)
			}
			return nil
		})
	}
	return found
}
