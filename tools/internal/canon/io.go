package canon

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"os"
	"path"
)

// Toda operación de disco pasa por os.Root: ninguna ruta, enlace simbólico
// ni «..» puede leer o escribir fuera de la carpeta raíz.

func readFile(root *os.Root, name string) ([]byte, error) {
	return root.ReadFile(name)
}

// writeFileAtomic escribe en un archivo temporal y lo renombra,
// para que dist nunca quede con un archivo a medio escribir.
func writeFileAtomic(root *os.Root, name string, data []byte) error {
	if err := root.MkdirAll(path.Dir(name), 0o755); err != nil {
		return err
	}
	tmp := name + ".tmp"
	if err := root.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := root.Rename(tmp, name); err != nil {
		_ = root.Remove(tmp)
		return err
	}
	return nil
}

// decodeObject lee un objeto JSON. Un archivo vacío cuenta como objeto vacío.
// encoding/json/v2 rechaza claves duplicadas, así que un error de tipeo no
// se pierde en silencio.
func decodeObject(data []byte) (map[string]any, error) {
	obj := map[string]any{}
	if len(bytes.TrimSpace(data)) == 0 {
		return obj, nil
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	if obj == nil { // el archivo contenía null
		obj = map[string]any{}
	}
	return obj, nil
}

// decodeHarness lee la lista «harness» de la configuración. Ignora el resto de claves.
func decodeHarness(data []byte) ([]string, error) {
	var s struct {
		Harness []string `json:"harness"`
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return s.Harness, nil
}

// encodeObject escribe JSON con claves ordenadas y sangría de dos espacios,
// para que la misma entrada produzca siempre los mismos bytes.
func encodeObject(v map[string]any) ([]byte, error) {
	out, err := json.Marshal(v,
		json.Deterministic(true),
		jsontext.Multiline(true),
		jsontext.WithIndent("  "),
	)
	if err != nil {
		return nil, err
	}
	return append(out, '\n'), nil
}
