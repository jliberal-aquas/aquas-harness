---
name: capability-docs
description: Escribe o actualiza en docs/ la documentación de una capacidad del gobierno (qué es, para qué sirve, cómo funciona y cómo se construyó), de lo abstracto a lo detallado, con diagramas de flujo Mermaid enlazados página a página y un caso conductor ejecutado de verdad. Cárgala antes de crear, cambiar o revisar cualquier archivo en docs/ que no sea un ADR.
when_to_use: Documentar una capacidad nueva del gobierno (p. ej. la Compilación del canon), actualizar una página tras cambiar su código, o cuando `docscheck` reporta enlaces rotos, cajas sin enlace, menciones sin definir o páginas sin vigencia.
---

# Documentación de capacidades · v0.1
Misión: que un lector que nunca vio el chat entienda qué sabe hacer el gobierno, para qué, cómo funciona y cómo se construyó, bajando de lo abstracto a lo detallado sin encontrar un concepto antes de su explicación.
No aplica a: ADR (`docs/adr/`), glosario como tal, comentarios de código, mensajes al humano.
Términos: los define `GLOSSARY.md` en la raíz (Capacidad, Proceso, Hoja, Caso conductor, Unidad de explicación, Primera mención, Vigencia). Método completo y su procedencia: [references/metodo.md](references/metodo.md).

## Prioridad si dos reglas chocan
1. Verdad: lo documentado coincide con el código y con una ejecución real de esta sesión.
2. Orden de conceptos: ningún término antes de su explicación.
3. Navegación: toda caja de un diagrama lleva a donde se explica.
4. Brevedad: quitar lo que no cambia la comprensión.

## Estructura de archivos
```
docs/
  README.md                       índice: una línea por capacidad
  <capacidad>/README.md           nivel 0: resumen, diagrama abstracto, para qué, cómo se construyó
  <capacidad>/caso-conductor.md   el ejemplo de principio a fin, ejecutado y con salida real
  <capacidad>/<proceso>.md        una página por caja que tenga decisiones o pasos internos
```
- Nombres de carpeta y archivo: castellano, minúsculas, guiones, sin tildes (`compilacion-del-canon`).
- Se baja otro nivel mientras una caja tenga bifurcaciones o pasos internos. Una **hoja** (una sola función sin bifurcaciones) se explica en una sección de la página padre y su caja enlaza a esa sección (`#ancla`).
- `docs/` nunca entra en la compilación de `canonc`.

## Flujo de trabajo
1. **Leer el código y sus pruebas** de la capacidad. Anotar cada función con decisiones: son candidatas a página.
2. **Diseñar el árbol de diagramas**: nivel 0 con 4–8 cajas; cada caja con decisiones baja un nivel. Parar en las hojas.
3. **Revisar el glosario**: todo término del proyecto que vayas a usar debe estar en `GLOSSARY.md`. Si falta, agregarlo (definición de 1–2 oraciones, sin detalles de implementación) antes de escribir páginas.
4. **Ejecutar el caso conductor** con datos de prueba en `product/.tmp/<capacidad>/`. Pegar comandos y salidas reales, sin editar. Nunca escribir una salida a mano.
5. **Escribir las páginas** de arriba hacia abajo con la plantilla de abajo.
6. **Correr `docscheck`** desde la raíz y corregir hasta que salga con código 0:
   `go run ./tools/cmd/docscheck`
7. **Revisión humana pedagógica** con la lista final (lo que `docscheck` no puede ver).
8. **Reportar** la salida de `docscheck` y lo que quedó sin comprobar.

## Plantilla de página de capacidad o proceso
```markdown
---
sources:
  - ruta/del/codigo.go
verified: <hash corto del commit en que se comparó con el código>
---

# <Nombre de la capacidad o del proceso>

> **Resumen.** Dos a cuatro oraciones: qué hace, con qué entrada, qué produce y por qué importa.

## Para qué sirve
El problema concreto que resuelve, con el caso conductor. Sin jerga nueva.

## Qué necesitas saber antes
Lista de enlaces a las páginas o términos del glosario que esta página da por sabidos. «Nada» en el nivel 0 si no hay.

## Cómo funciona
<diagrama Mermaid>

- [Etiqueta exacta de la caja 1](pagina-hija.md): una línea.
- [Etiqueta exacta de la caja 2](#seccion-de-la-hoja): una línea.

### <Una sección por caja, en el orden del flujo>
Unidad de explicación (6 pasos): pregunta → ejemplo → diagrama si hace falta → nombre → causa → enlace.

## Cómo se construyó            (solo en el README de la capacidad)
Decisiones con su porqué y enlaces a ADR y commits. No es un diario.
```

## Diagramas
- Siempre `flowchart TD`, unos 10 nodos como máximo, etiquetas de 3 palabras o menos. Más nodos = partir en una página hija.
- Formas con significado fijo: `([...])` inicio o fin; `[...]` proceso; `{...}` decisión. Solo las cajas `[...]` exigen enlace.
- Debajo de cada diagrama, la lista de enlaces: un ítem por caja `[...]`, con la **misma etiqueta** como texto del enlace. Es la vía de navegación que siempre funciona.
- No usar `click` de Mermaid: GitHub no documenta soporte y la vista previa de VS Code corre en modo estricto (verificado el 2026-11-26 en docs.github.com «Creating diagrams» y code.visualstudio.com/docs/languages/markdown).
- `stateDiagram` o `sequenceDiagram` solo si la pregunta es «¿en qué estado está?» o «¿quién habla con quién?».
- Todo diagrama responde a una pregunta escrita en la oración anterior. Prohibido el diagrama decorativo.

## Orden de conceptos
- Un término del glosario se explica una sola vez, en `GLOSSARY.md`.
- En cada página, su **primera mención** va enlazada al glosario (`[capa](../../GLOSSARY.md)`) o en **negrita** con su explicación ahí mismo. También cuenta un enlace a la página que lo explica.
- Las etiquetas de un diagrama cuentan como mención: el término debe quedar enlazado o definido en la prosa antes del diagrama.
- Antes de nombrar un concepto, describirlo con palabras llanas; el nombre va después de la demostración.
- Los títulos no cuentan como mención y no llevan enlaces.

## Caso conductor
- Uno por capacidad, sin personajes. Ejemplo: «agregamos una regla en la capa de producto y compilamos».
- Se ejecuta de verdad; su página muestra comandos, salidas y archivos resultantes, marcados como ejecución real con fecha.
- Cada página de proceso vuelve a ese mismo caso a su profundidad («en el caso conductor, esta caja recibe…»).

## Estilo
- Castellano técnico neutro, trato de tú, voz activa, oraciones cortas. Sin modismos.
- Literales intactos: rutas, comandos, nombres de función y campos van en `código`.
- Causa con palabras: «queríamos X; ocurrió Y; significa Z».
- No hay sección de errores: se documenta lo que la capacidad hace.
- Texto alternativo de imágenes descriptivo.

## Comprobación final
Determinista (`go run ./tools/cmd/docscheck` → código 0):
- Enlaces relativos y anclas existen.
- Cada caja `[...]` de cada diagrama tiene su enlace en la lista de debajo.
- Primera mención de cada término del glosario enlazada o en negrita, incluidas las etiquetas de diagramas.
- Cada página de capacidad trae `sources` y `verified`, y ningún archivo de `sources` cambió después de `verified`.

Revisión humana (`docscheck` no lo ve):
- El resumen se entiende sin leer el resto.
- Cada sección sigue la unidad de explicación y el nombre va después del ejemplo.
- El caso conductor aparece en cada página y su salida es real.
- Cada diagrama responde a una pregunta escrita.
