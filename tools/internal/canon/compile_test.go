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

// tree crea archivos bajo una carpeta temporal. Claves: rutas con «/».
func tree(t *testing.T, files map[string]string) string {
	t.Helper()
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

func TestMarkdownGobiernoAntesQueProducto(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/canon/human-communication/rules.md":   "# Gobierno\nregla G\n",
		".product-governance/canon/human-communication/rules.md": "# Producto\nregla P\n",
	})
	if _, err := Run(dir, false, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	got := read(t, dir, "dist/human-communication/rules.md")
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
	if _, err := Run(dir, false, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if got := read(t, dir, "dist/insights/rules.md"); !strings.HasSuffix(got, "solo producto\n") {
		t.Fatalf("contenido inesperado:\n%s", got)
	}
}

func TestBOMyCRLFSeNormalizan(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/canon/x/rules.md":   "\uFEFFlínea 1\r\nlínea 2\r\n",
		".product-governance/canon/x/rules.md": "\uFEFFlínea 3\r\n",
	})
	if _, err := Run(dir, false, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	got := read(t, dir, "dist/x/rules.md")
	if strings.Contains(got, "\uFEFF") || strings.Contains(got, "\r") {
		t.Fatalf("quedó BOM o CRLF: %q", got)
	}
}

func TestJSONSeFusionaYProductoGana(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/config/settings.json":   `{"a": 1, "anidado": {"x": "gob", "y": "gob"}, "lista": [1, 2]}`,
		".product-governance/config/settings.json": `{"b": 2, "anidado": {"y": "prod"}, "lista": [3]}`,
	})
	if _, err := Run(dir, false, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	got := read(t, dir, "dist/config/settings.json")
	var m map[string]any
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("salida no es JSON válido: %v\n%s", err, got)
	}
	want := map[string]any{
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
		".aquas-governance/config/settings.json":   "",
		".product-governance/config/settings.json": `{"b": 2}`,
	})
	if _, err := Run(dir, false, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if got := read(t, dir, "dist/config/settings.json"); !strings.Contains(got, `"b": 2`) {
		t.Fatalf("contenido inesperado:\n%s", got)
	}
}

func TestJSONInvalidoFallaConRuta(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/config/settings.json": `{"a": }`,
	})
	_, err := Run(dir, false, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), ".aquas-governance/config/settings.json") {
		t.Fatalf("se esperaba error con la ruta del archivo, llegó: %v", err)
	}
}

func TestSegundaEjecucionNoCambiaNadaYCheckPasa(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/canon/x/rules.md":       "G",
		".product-governance/config/settings.json": `{"b": 2}`,
	})
	if _, err := Run(dir, false, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	stale, err := Run(dir, true, &out)
	if err != nil || stale {
		t.Fatalf("check debía pasar: stale=%v err=%v\n%s", stale, err, out.String())
	}
}

func TestCheckDetectaDesactualizadoSinEscribir(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/canon/x/rules.md": "G",
		"dist/x/rules.md":                    "viejo",
	})
	stale, err := Run(dir, true, &bytes.Buffer{})
	if err != nil || !stale {
		t.Fatalf("check debía detectar diferencia: stale=%v err=%v", stale, err)
	}
	if got := read(t, dir, "dist/x/rules.md"); got != "viejo" {
		t.Fatalf("check no debe escribir; dist quedó: %q", got)
	}
}

func TestHuerfanoSeAvisaYNoSeBorra(t *testing.T) {
	dir := tree(t, map[string]string{
		".aquas-governance/canon/x/rules.md": "G",
		"dist/viejo/rules.md":                "huérfano",
	})
	var out bytes.Buffer
	if _, err := Run(dir, false, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "dist/viejo/rules.md no viene de ninguna capa") {
		t.Errorf("faltó el aviso de huérfano:\n%s", out.String())
	}
	if read(t, dir, "dist/viejo/rules.md") != "huérfano" {
		t.Error("el huérfano no debía borrarse")
	}
}

func TestRaizSinCapasFalla(t *testing.T) {
	if _, err := Run(t.TempDir(), false, &bytes.Buffer{}); err == nil {
		t.Fatal("se esperaba error con una raíz vacía")
	}
}
