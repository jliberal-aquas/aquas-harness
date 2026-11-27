package docs

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// glossary es un glosario mínimo con dos términos, uno de dos palabras.
const glossary = "# G\n\n**Capa**:\nUna fuente.\n\n**Caso conductor**:\nEl ejemplo.\n"

// meta es una cabecera de vigencia válida para las pruebas en memoria.
const meta = "---\nsources:\n  - a.go\nverified: abc123\n---\n"

// site arma un repositorio en memoria con el glosario y las páginas dadas.
// La vigencia nunca está vencida salvo que stale traiga cambios.
func site(pages map[string]string, stale []string) Site {
	md := map[string]string{"GLOSSARY.md": glossary}
	for k, v := range pages {
		md[k] = v
	}
	return Site{
		Pages:    pages,
		Glossary: "GLOSSARY.md",
		Markdown: md,
		Exists: func(p string) bool {
			_, ok := md[p]
			return ok || p == "a.go"
		},
		Stale: func(string, []string) ([]string, error) { return stale, nil },
	}
}

// messages devuelve los mensajes de los hallazgos, uno por línea, para comparar.
func messages(fs []Finding) string {
	var b strings.Builder
	for _, f := range fs {
		b.WriteString(f.String() + "\n")
	}
	return b.String()
}

// TestCheck recorre cada comprobación con un caso que debe fallar y otro que debe pasar.
func TestCheck(t *testing.T) {
	const p = "docs/cap/README.md"
	tests := []struct {
		name  string
		page  string
		stale []string
		want  []string // fragmentos esperados; vacío = sin hallazgos
	}{
		{
			name: "valid page passes",
			page: meta + "# T\n\nUna [capa](../../GLOSSARY.md) y otra capa.\n\n```mermaid\nflowchart TD\n  I([Inicio]) --> A[Leer capas]\n  A --> D{¿Listo?}\n```\n\n- [Leer capas](#t): lee.\n",
		},
		{
			name: "broken link fails",
			page: meta + "# T\n\nVer [otra](no-existe.md).\n",
			want: []string{`enlace roto: "no-existe.md" no existe`},
		},
		{
			name: "missing anchor fails",
			page: meta + "# T\n\nVer [arriba](#otro).\n",
			want: []string{`no hay un título con el ancla #otro`},
		},
		{
			name: "accented anchor passes",
			page: meta + "# Cómo funciona\n\nVer [arriba](#cómo-funciona).\n",
		},
		{
			name: "box without link fails",
			page: meta + "# T\n\n```mermaid\nflowchart TD\n  A[Uno] --> B[Dos]\n```\n\n- [Uno](#t)\n",
			want: []string{`la caja "Dos" no tiene enlace`},
		},
		{
			name: "start, end and decision need no link",
			page: meta + "# T\n\n```mermaid\nflowchart TD\n  I([Inicio]) --> D{Decide} --> F((Fin))\n```\n",
		},
		{
			name: "data boxes and subgraph titles need no link",
			page: meta + "# T\n\n```mermaid\nflowchart TD\n  subgraph G[\"Gobierno\"]\n    direction LR\n    a1[/\"x: true\"/] --> c1[/x: false/]\n  end\n```\n",
		},
		{
			name: "process box inside a subgraph still needs a link",
			page: meta + "# T\n\n```mermaid\nflowchart TD\n  subgraph G[\"Gobierno\"]\n    a1[/x/] --> B[Leer]\n  end\n```\n",
			want: []string{`la caja "Leer" no tiene enlace`},
		},
		{
			name: "term first mention without link fails",
			page: meta + "# T\n\nCada capa aporta. Ver [capa](../../GLOSSARY.md).\n",
			want: []string{`primera mención de "Capa" sin enlace ni negrita`},
		},
		{
			name: "bold definition and plural pass",
			page: meta + "# T\n\nUna **capa** es una fuente. Hay dos capas.\n",
		},
		{
			name: "multiword term in diagram before prose fails",
			page: meta + "# T\n\n```mermaid\nflowchart TD\n  A[Caso conductor]\n```\n\n- [Caso conductor](#t)\n",
			want: []string{`"Caso conductor" aparece en un diagrama antes`},
		},
		{
			name: "term in heading and inline code is ignored",
			page: meta + "# Capa\n\nEl campo `capa` no cuenta.\n",
		},
		{
			name: "missing freshness header fails",
			page: "# T\n",
			want: []string{"falta la cabecera de vigencia"},
		},
		{
			name:  "stale sources fail",
			page:  meta + "# T\n",
			stale: []string{"commit def456"},
			want:  []string{"desactualizada: las fuentes cambiaron después de abc123 (commit def456)"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := messages(Check(site(map[string]string{p: tt.page}, tt.stale)))
			if len(tt.want) == 0 && got != "" {
				t.Fatalf("got findings, want none:\n%s", got)
			}
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("got:\n%s\nwant a finding containing %q", got, w)
				}
			}
		})
	}
}

// TestIndexAndADRNeedNoFreshness comprueba que docs/README.md y los ADR no exigen cabecera.
func TestIndexAndADRNeedNoFreshness(t *testing.T) {
	pages := map[string]string{
		"docs/README.md":     "# Índice\n",
		"docs/adr/0001-x.md": "# Decisión sobre la capa\n\nLa capa no se enlaza en un ADR.\n",
	}
	if got := messages(Check(site(pages, nil))); got != "" {
		t.Fatalf("got findings, want none:\n%s", got)
	}
}

// TestRunDetectsChangeAfterVerified usa un repositorio git real: la página pasa
// en el commit verificado y queda desactualizada al cambiar su fuente.
func TestRunDetectsChangeAfterVerified(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git no está instalado")
	}
	dir := t.TempDir()
	write := func(name, content string) {
		t.Helper()
		p := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitRun := func(args ...string) string {
		t.Helper()
		out, err := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=t", "-c", "user.email=t@t"}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	gitRun("init", "-q")
	write("GLOSSARY.md", glossary)
	write("a.go", "package a\n")
	gitRun("add", ".")
	gitRun("commit", "-qm", "base")
	head := gitRun("rev-parse", "--short", "HEAD")
	write("docs/cap/README.md", "---\nsources:\n  - a.go\nverified: "+head+"\n---\n# T\n")

	var out bytes.Buffer
	n, err := Run(context.Background(), dir, &out)
	if err != nil || n != 0 {
		t.Fatalf("got %d findings, err %v, want 0:\n%s", n, err, out.String())
	}

	write("a.go", "package a // cambio\n")
	out.Reset()
	n, err = Run(context.Background(), dir, &out)
	if err != nil || n != 1 || !strings.Contains(out.String(), "sin commit:") {
		t.Fatalf("got %d findings, err %v, want 1 uncommitted change:\n%s", n, err, out.String())
	}
}
