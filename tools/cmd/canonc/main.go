// Command canonc compila las reglas canónicas, la configuración y las skills
// para cada arnés declarado en «harness» de config/settings.json.
//
// Lee cada capa (primero gobierno, después producto). Si dos capas traen el
// mismo archivo, une sus aportes; las skills se copian y deben ser únicas.
// Para claude-code:
//
//	.aquas-governance/canon/<x>/<archivo>   ┐
//	.product-governance/canon/<x>/<archivo> ┴→ .claude/rules/<x>/<archivo>
//	.aquas-governance/config/<archivo>      ┐
//	.product-governance/config/<archivo>    ┴→ .claude/aquas/<archivo>
//	.<capa>/skills/<nombre>/...               → .claude/skills/<nombre>/...
//
// Uso:
//
//	canonc [-root carpeta] [-check]
package main

import (
	"flag"
	"fmt"
	"os"

	"aquas-tools/tools/internal/canon"
)

// main lee las opciones, ejecuta canon.Run y traduce el resultado a código de salida.
func main() {
	root := flag.String("root", ".", "carpeta raíz del arnés (contiene .aquas-governance y .product-governance)")
	check := flag.Bool("check", false, "solo verifica que el arnés esté al día; no escribe nada y sale con código 1 si hay diferencias")
	flag.Parse()

	stale, err := canon.Run(*root, *check, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	if *check && stale {
		fmt.Fprintln(os.Stderr, "el arnés está desactualizado: ejecuta canonc sin -check")
		os.Exit(1)
	}
}
