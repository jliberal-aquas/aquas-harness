// Command canonc compila las reglas canónicas del arnés.
//
// Lee cada capa (primero gobierno, después producto) y escribe en dist
// un archivo por nombre, con el aporte de todas las capas que lo tengan:
//
//	.aquas-governance/canon/<x>/<archivo>   ┐
//	.product-governance/canon/<x>/<archivo> ┴→ dist/<x>/<archivo>
//	.aquas-governance/config/<archivo>      ┐
//	.product-governance/config/<archivo>    ┴→ dist/config/<archivo>
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

func main() {
	root := flag.String("root", ".", "carpeta raíz del arnés (contiene .aquas-governance, .product-governance y dist)")
	check := flag.Bool("check", false, "solo verifica que dist esté al día; no escribe nada y sale con código 1 si hay diferencias")
	flag.Parse()

	stale, err := canon.Run(*root, *check, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	if *check && stale {
		fmt.Fprintln(os.Stderr, "dist está desactualizado: ejecuta canonc sin -check")
		os.Exit(1)
	}
}
