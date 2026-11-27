package canon

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// govSettings es la configuración de gobierno; por omisión declara Claude Code como único arnés.
const govSettings = ".aquas-governance/config/settings.json"

// tree crea archivos bajo una carpeta temporal. Claves: rutas con «/».
// Si no trae configuración de gobierno, agrega una con el arnés claude-code.
func tree(t *testing.T, files map[string]string) string {
	t.Helper()
	if _, ok := files[govSettings]; !ok {
		files[govSettings] = `{"harness": ["claude-code"]}`
	}
	dir := t.TempDir()
	for name, content := range files {
		p := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func read(t *testing.T, dir, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(name)))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// run compila sin -check y falla la prueba si hay error. Devuelve la salida.
func run(t *testing.T, dir string) string {
	t.Helper()
	var out bytes.Buffer
	if _, err := Run(dir, false, &out); err != nil {
		t.Fatalf("Run: %v\n%s", err, out.String())
	}
	return out.String()
}

func TestMarkdownGobiernoAntesQueProducto(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/canon/human-communication/rules.md":   "# Gobierno\nregla G\n",
		".product-governance/canon/human-communication/rules.md": "# Producto\nregla P\n",
	})
	run(t, dir)
	got := read(t, dir, ".claude/rules/human-communication/rules.md")
	g, p := strings.Index(got, "regla G"), strings.Index(got, "regla P")
	if g < 0 || p < 0 || g > p {
		t.Fatalf("orden incorrecto o contenido faltante:\n%s", got)
	}
	if !strings.HasPrefix(got, "<!-- Generado por canonc") {
		t.Errorf("falta la cabecera de archivo generado:\n%s", got)
	}
	if !strings.Contains(got, "regla G\n\n# Producto") {
		t.Errorf("las capas deben ir separadas por una línea en blanco:\n%s", got)
	}
}

func TestArchivoEnUnaSolaCapaSeCopia(t *testing.T) {
	dir := tree(t, map[string]string{
		".product-governance/canon/insights/rules.md": "solo producto",
	})
	run(t, dir)
	if got := read(t, dir, ".claude/rules/insights/rules.md"); !strings.HasSuffix(got, "solo producto\n") {
		t.Fatalf("contenido inesperado:\n%s", got)
	}
}

func TestBOMyCRLFSeNormalizan(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/canon/x/rules.md":   "\xEF\xBB\xBFlínea 1\r\nlínea 2\r\n",
		".product-governance/canon/x/rules.md": "\xEF\xBB\xBFlínea 3\r\n",
	})
	run(t, dir)
	got := read(t, dir, ".claude/rules/x/rules.md")
	if strings.Contains(got, "\xEF\xBB\xBF") || strings.Contains(got, "\r") {
		t.Fatalf("quedó BOM o CRLF: %q", got)
	}
}

func TestSkillSeCopiaIntacta(t *testing.T) {
	skill := "---\nname: go\n---\n\n# Go\r\n"
	dir := tree(t, map[string]string{
		".aquas-governance/skills/go/SKILL.md":       skill,
		".aquas-governance/skills/go/ref/extra.md":   "extra",
		".product-governance/skills/propia/SKILL.md": "---\nname: propia\n---\n",
	})
	run(t, dir)
	if got := read(t, dir, ".claude/skills/go/SKILL.md"); got != skill {
		t.Errorf("la skill debe copiarse byte a byte, sin cabecera ni normalizar:\n got  %q\n want %q", got, skill)
	}
	if got := read(t, dir, ".claude/skills/go/ref/extra.md"); got != "extra" {
		t.Errorf("faltó un archivo interno de la skill: %q", got)
	}
	if got := read(t, dir, ".claude/skills/propia/SKILL.md"); !strings.HasPrefix(got, "---") {
		t.Errorf("faltó la skill de producto: %q", got)
	}
}

func TestSkillRepetidaEnDosCapasFalla(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/skills/go/SKILL.md":   "gob",
		".product-governance/skills/go/SKILL.md": "prod",
	})
	_, err := Run(dir, false, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "skills/go/SKILL.md") {
		t.Fatalf("se esperaba error por skill repetida, llegó: %v", err)
	}
}

func TestConfigVaALaCarpetaAquasDelArnes(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/canon/x/rules.md": "G",
	})
	run(t, dir)
	if got := read(t, dir, ".claude/aquas/settings.json"); !strings.Contains(got, `"claude-code"`) {
		t.Fatalf("configuración inesperada:\n%s", got)
	}
}

func TestArnesDesconocidoFalla(t *testing.T) {
	dir := tree(t, map[string]string{
		govSettings: `{"harness": ["claude-code", "inventado"]}`,
	})
	_, err := Run(dir, false, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), `"inventado"`) {
		t.Fatalf("se esperaba error por arnés desconocido, llegó: %v", err)
	}
}

func TestSinArnesFalla(t *testing.T) {
	dir := tree(t, map[string]string{
		govSettings:                          `{"os": "windows"}`,
		".aquas-governance/canon/x/rules.md": "G",
	})
	if _, err := Run(dir, false, &bytes.Buffer{}); err == nil {
		t.Fatal("se esperaba error si la configuración no declara harness")
	}
}

func TestJSONSeFusionaYProductoGana(t *testing.T) {
	dir := tree(t, map[string]string{
		govSettings: `{"harness": ["claude-code"], "a": 1, "anidado": {"x": "gob", "y": "gob"}, "lista": [1, 2]}`,
		".product-governance/config/settings.json": `{"b": 2, "anidado": {"y": "prod"}, "lista": [3]}`,
	})
	run(t, dir)
	got := read(t, dir, ".claude/aquas/settings.json")
	var m map[string]any
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("salida no es JSON válido: %v\n%s", err, got)
	}
	want := map[string]any{
		"harness": []any{"claude-code"},
		"a":       1.0,
		"b":       2.0,
		"anidado": map[string]any{"x": "gob", "y": "prod"},
		"lista":   []any{3.0},
	}
	if !reflect.DeepEqual(m, want) {
		t.Errorf("fusión incorrecta:\n got  %v\n want %v", m, want)
	}
	if !strings.Contains(got, "\n  \"a\": 1,") {
		t.Errorf("se esperaba sangría de dos espacios y claves ordenadas:\n%s", got)
	}
}

func TestJSONVacioCuentaComoObjetoVacio(t *testing.T) {
	dir := tree(t, map[string]string{
		".product-governance/config/settings.json": "",
	})
	run(t, dir)
	if got := read(t, dir, ".claude/aquas/settings.json"); !strings.Contains(got, `"claude-code"`) {
		t.Fatalf("contenido inesperado:\n%s", got)
	}
}

func TestJSONInvalidoFallaConRuta(t *testing.T) {
	dir := tree(t, map[string]string{
		".product-governance/config/settings.json": `{"a": }`,
	})
	_, err := Run(dir, false, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), ".product-governance/config/settings.json") {
		t.Fatalf("se esperaba error con la ruta del archivo, llegó: %v", err)
	}
}

func TestSegundaEjecucionNoCambiaNadaYCheckPasa(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/canon/x/rules.md":   "G",
		".aquas-governance/skills/go/SKILL.md": "S",
	})
	run(t, dir)
	var out bytes.Buffer
	stale, err := Run(dir, true, &out)
	if err != nil || stale {
		t.Fatalf("check debía pasar: stale=%v err=%v\n%s", stale, err, out.String())
	}
}

func TestCheckDetectaDesactualizadoSinEscribir(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/canon/x/rules.md": "G",
		".claude/rules/x/rules.md":           "viejo",
	})
	stale, err := Run(dir, true, &bytes.Buffer{})
	if err != nil || !stale {
		t.Fatalf("check debía detectar diferencia: stale=%v err=%v", stale, err)
	}
	if got := read(t, dir, ".claude/rules/x/rules.md"); got != "viejo" {
		t.Fatalf("check no debe escribir; quedó: %q", got)
	}
}

func TestHuerfanoSoloEnCarpetasPropias(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/canon/x/rules.md": "G",
		".claude/rules/viejo/rules.md":       "huérfano",
		".claude/settings.local.json":        "{}",
	})
	out := run(t, dir)
	if !strings.Contains(out, ".claude/rules/viejo/rules.md no viene de ninguna capa") {
		t.Errorf("faltó el aviso de huérfano:\n%s", out)
	}
	if strings.Contains(out, "settings.local.json") {
		t.Errorf("no debe opinar sobre archivos de .claude que no son suyos:\n%s", out)
	}
	if read(t, dir, ".claude/rules/viejo/rules.md") != "huérfano" {
		t.Error("el huérfano no debía borrarse")
	}
}

func TestRaizSinCapasFalla(t *testing.T) {
	if _, err := Run(t.TempDir(), false, &bytes.Buffer{}); err == nil {
		t.Fatal("se esperaba error con una raíz vacía")
	}
}
