// Package canon compila las reglas canónicas del arnés: une cada archivo
// de gobierno con su par de producto y escribe el resultado en dist.
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

// sections asocia cada subcarpeta de una capa con su destino dentro de dist.
var sections = []struct{ src, dst string }{
	{src: "canon", dst: "."},
	{src: "config", dst: "config"},
}

const distDir = "dist"

// bom es la marca de orden de bytes UTF-8 que agregan algunos editores de Windows.
var bom = []byte("\xEF\xBB\xBF")

// part es el aporte de una capa a un archivo de salida.
type part struct {
	layer string // nombre de la capa
	src   string // ruta relativa a la raíz, con «/»
	data  []byte
}

// output es un archivo de dist con sus aportes en orden jerárquico.
type output struct {
	dst   string // ruta relativa a dist, con «/»
	parts []part
}

// Run compila todas las capas en dist. Con check no escribe nada.
// Devuelve stale = true si algún archivo de dist no coincide con lo compilado.
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

	var errs []error
	for _, dst := range slices.Sorted(maps.Keys(outs)) {
		o := outs[dst]
		target := path.Join(distDir, dst)

		want, err := o.build()
		if err != nil {
			errs = append(errs, err)
			continue
		}
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

	for _, orphan := range orphans(root, outs) {
		fmt.Fprintf(w, "aviso: %s no viene de ninguna capa; revísalo o bórralo\n", orphan)
	}
	return stale, errors.Join(errs...)
}

// collect recorre las capas en orden jerárquico y agrupa los archivos por destino.
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
				rel := strings.TrimPrefix(p, base+"/")
				dst := path.Join(s.dst, rel)
				o := outs[dst]
				if o == nil {
					o = &output{dst: dst}
					outs[dst] = o
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

// build combina los aportes según el tipo de archivo.
func (o *output) build() ([]byte, error) {
	switch path.Ext(o.dst) {
	case ".md", ".txt":
		return concatText(o), nil
	case ".json":
		return mergeJSON(o.parts)
	}
	if len(o.parts) == 1 {
		return o.parts[0].data, nil
	}
	return nil, fmt.Errorf("%s: no sé combinar archivos %q de varias capas", o.dst, path.Ext(o.dst))
}

func (o *output) layerNames() string {
	names := make([]string, len(o.parts))
	for i, p := range o.parts {
		names[i] = p.layer
	}
	return strings.Join(names, " + ")
}

// concatText une los textos en orden jerárquico, separados por una línea en blanco.
// En Markdown agrega un comentario que indica las fuentes y avisa que no se edite.
func concatText(o *output) []byte {
	var b bytes.Buffer
	if path.Ext(o.dst) == ".md" {
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

// mergeJSON fusiona objetos JSON en orden jerárquico.
// Objetos anidados se fusionan clave por clave; cualquier otro valor
// (texto, número, lista) de una capa posterior reemplaza al anterior.
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

// orphans lista archivos de dist que ninguna capa produce. No los borra.
func orphans(root *os.Root, outs map[string]*output) []string {
	var found []string
	_ = fs.WalkDir(root.FS(), distDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if _, ok := outs[strings.TrimPrefix(p, distDir+"/")]; !ok {
			found = append(found, p)
		}
		return nil
	})
	return found
}
