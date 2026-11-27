package docs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

// docsDir es la carpeta de documentación y glossaryFile el glosario, ambos desde la raíz.
const (
	docsDir      = "docs"
	glossaryFile = "GLOSSARY.md"
)

// Run lee docs/ y el glosario bajo root, revisa todo y escribe un hallazgo por línea.
// Devuelve cuántos hallazgos hubo; err solo si no se pudo leer o consultar git.
func Run(ctx context.Context, rootPath string, w io.Writer) (int, error) {
	root, err := os.OpenRoot(rootPath)
	if err != nil {
		return 0, fmt.Errorf("no pude abrir la raíz %q: %w", rootPath, err)
	}
	defer root.Close()

	site, err := load(root)
	if err != nil {
		return 0, err
	}
	site.Stale = func(verified string, sources []string) ([]string, error) {
		return gitChanges(ctx, rootPath, verified, sources)
	}

	findings := Check(site)
	for _, f := range findings {
		fmt.Fprintln(w, f)
	}
	fmt.Fprintf(w, "%d páginas revisadas, %d hallazgos\n", len(site.Pages), len(findings))
	return len(findings), nil
}

// load lee todas las páginas .md de docs/ y el glosario.
func load(root *os.Root) (Site, error) {
	s := Site{
		Pages:    map[string]string{},
		Glossary: glossaryFile,
		Markdown: map[string]string{},
		Exists: func(p string) bool {
			_, err := root.Stat(p)
			return err == nil
		},
	}
	err := fs.WalkDir(root.FS(), docsDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		data, err := root.ReadFile(p)
		if err != nil {
			return err
		}
		s.Pages[p] = string(data)
		s.Markdown[p] = string(data)
		return nil
	})
	if err != nil {
		return s, fmt.Errorf("recorriendo %s: %w", docsDir, err)
	}
	data, err := root.ReadFile(glossaryFile)
	if errors.Is(err, fs.ErrNotExist) {
		return s, fmt.Errorf("falta %s en la raíz: ahí se definen los términos", glossaryFile)
	}
	if err != nil {
		return s, fmt.Errorf("leyendo %s: %w", glossaryFile, err)
	}
	s.Markdown[glossaryFile] = string(data)
	return s, nil
}

// gitChanges lista los commits y cambios sin commit que tocaron sources
// después de verified. Error si verified no es un commit del repositorio.
func gitChanges(ctx context.Context, root, verified string, sources []string) ([]string, error) {
	if _, err := git(ctx, root, "rev-parse", "--verify", "--quiet", verified+"^{commit}"); err != nil {
		return nil, fmt.Errorf("%q no es un commit de este repositorio", verified)
	}
	args := append([]string{"log", "--format=commit %h", verified + "..HEAD", "--"}, sources...)
	logged, err := git(ctx, root, args...)
	if err != nil {
		return nil, err
	}
	args = append([]string{"status", "--porcelain", "--"}, sources...)
	pending, err := git(ctx, root, args...)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, l := range strings.Split(logged+pending, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			if !strings.HasPrefix(l, "commit ") {
				l = "sin commit: " + l
			}
			out = append(out, l)
		}
	}
	return out, nil
}

// git ejecuta un comando git en root y devuelve su salida estándar.
func git(ctx context.Context, root string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", root}, args...)...)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %v: %s", args[0], err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}
