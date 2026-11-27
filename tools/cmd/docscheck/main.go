// Command docscheck revisa la documentación de capacidades en docs/:
//   - enlaces relativos y anclas que existen;
//   - cada caja de proceso de un diagrama Mermaid con su enlace debajo;
//   - primera mención de cada término de GLOSSARY.md enlazada o en negrita;
//   - cabecera de vigencia («sources», «verified») sin cambios posteriores en git.
//
// Uso:
//
//	docscheck [-root carpeta]
//
// Sale con 0 si no hay hallazgos, 1 si los hay y 2 si no pudo revisar.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"aquas-tools/tools/internal/docs"
)

// main lee las opciones, ejecuta docs.Run y traduce el resultado a código de salida.
func main() {
	root := flag.String("root", ".", "carpeta raíz del repositorio (contiene docs/ y GLOSSARY.md)")
	flag.Parse()

	n, err := docs.Run(context.Background(), *root, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	if n > 0 {
		os.Exit(1)
	}
}
