---
name: capability-docs
description: Escribe o actualiza en docs/ el manual de referencia de una capacidad del gobierno (qué hace, cómo funciona y cómo se construyó), con un mapa «Estás aquí», diagramas de entrada y resultado sobre un caso conductor ejecutado de verdad, síntesis y detalle desplegable. Cárgala antes de crear, cambiar o revisar cualquier archivo en docs/ que no sea un ADR.
when_to_use: Documentar una capacidad nueva del gobierno (p. ej. la Compilación del canon), actualizar una página tras cambiar su código, o cuando `docscheck` reporta enlaces rotos, cajas sin enlace, menciones sin definir o páginas sin vigencia.
---

# Documentación de capacidades · v0.2
Misión: que un lector que nunca vio el chat entienda qué sabe hacer el gobierno, cómo funciona y por qué se construyó así, viendo lo que pasa antes de leer el nombre de lo que pasa.
Es un **manual de referencia**, no un manual de entrenamiento: afirma y muestra; no pregunta, no pide predecir ni propone ejercicios.
No aplica a: ADR (`docs/adr/`), glosario como tal, comentarios de código, mensajes al humano.
Términos: los define `GLOSSARY.md` en la raíz (Capacidad, Proceso, Hoja, Caso conductor, Unidad de explicación, Primera mención, Vigencia). Método completo y su procedencia: [references/metodo.md](references/metodo.md).

## Prioridad si dos reglas chocan
1. Verdad: lo documentado coincide con el código y con una ejecución real de esta sesión.
2. Orden de conceptos: ningún término antes de su explicación; el nombre va después de verlo funcionar.
3. Lectura visual: lo que se puede mostrar en un diagrama se muestra; la prosa acompaña.
4. Navegación: toda caja de proceso lleva a donde se explica.
5. Brevedad: quitar lo que no cambia la comprensión; el detalle se pliega, no se borra.

## Estructura de archivos
```
docs/
  README.md                       índice: una línea por capacidad
  <capacidad>/README.md           nivel 0: el viaje completo, mapa sin resaltar, cómo se construyó
  <capacidad>/caso-conductor.md   el ejemplo de principio a fin, ejecutado y con salida real
  <capacidad>/<proceso>.md        una página por caja que tenga decisiones o pasos internos
```
- Nombres de carpeta y archivo: castellano, minúsculas, guiones, sin tildes (`compilacion-del-canon`).
- Se baja otro nivel mientras una caja tenga bifurcaciones o pasos internos. Una **hoja** (una sola función sin bifurcaciones) se explica en una sección de la página padre y su caja enlaza a esa sección (`#ancla`).
- `docs/` nunca entra en la compilación de `canonc`.

## Flujo de trabajo
1. **Leer el código y sus pruebas** de la capacidad. Anotar cada función con decisiones: son candidatas a página.
2. **Diseñar el mapa**: nivel 0 con 4–8 cajas en fila. Es el mismo en todas las páginas de la capacidad. Cada caja con decisiones baja un nivel; parar en las hojas.
3. **Revisar el glosario**: todo término del proyecto que vayas a usar debe estar en `GLOSSARY.md`. Si falta, agregarlo (1–2 oraciones, sin detalles de implementación) antes de escribir páginas.
4. **Ejecutar el caso conductor** con datos de prueba en `product/.tmp/<capacidad>/`. Elegir datos que muestren cada regla en acción (p. ej. una clave que gana, una lista que se reemplaza, un objeto que se mezcla) y un caso que se detiene. Pegar comandos y salidas reales, sin editar. Nunca escribir una salida a mano.
5. **Escribir las páginas** de arriba hacia abajo con el ritmo de abajo.
6. **Ver los diagramas dibujados** en un navegador antes de entregar: que ninguno dé error y que ninguno cruce flechas al punto de no leerse.
7. **Correr `docscheck`** desde la raíz y corregir hasta que salga con código 0: `go run ./tools/cmd/docscheck`
8. **Revisión humana** con la lista final (lo que `docscheck` no puede ver).
9. **Reportar** la salida de `docscheck` y lo que quedó sin comprobar.

## Ritmo de cada página (unidad de explicación)
Siempre en este orden. Una sección se omite solo si no tiene contenido real (p. ej. no hay «Cuándo se detiene» si el proceso no puede fallar).

```markdown
---
sources:
  - ruta/del/codigo.go
verified: <hash corto del commit en que se comparó con el código>
---

# <Nombre de la capacidad o del proceso>

> **En una frase.** Una oración: qué hace, con qué entrada y qué produce.

## Estás aquí
<mapa de nivel 0 con esta caja resaltada; en el README, sin resaltar>
- [Caja 1](pagina.md): una línea.   (lista de enlaces: una por caja del mapa)

## Lo que ya sabes
Dos o tres oraciones que retoman la página anterior en el caso conductor y dejan
planteado lo que esta resuelve. El contexto va escrito aquí, no en una lista de enlaces.
En el README: el problema que resuelve la capacidad.

## Entrada
<diagrama de datos: lo que llega a este proceso en el caso conductor>

## Resultado
<diagrama de antes y después: un bloque por caso, con su regla en el título>
Una o dos oraciones con la causa. Luego la salida real, con fecha y commit.

## <Nombre del mecanismo>
Recién aquí el nombre, en **negrita**, y su porqué: «se diseñó así para que…».
Diagrama de decisión si el mecanismo elige entre caminos.
Lista de enlaces bajo el diagrama si tiene cajas de proceso.

## Cuándo se detiene
Los límites y errores que la capacidad hace a propósito, con su salida real.

## Síntesis
<diagrama-resumen: una decisión con sus finales>
El modelo en una oración y una consecuencia directa que se desprende de él.

<details>
<summary><strong>Detalle completo</strong> (para quien va a cambiar el código)</summary>

Funciones y archivos, casos límite, formatos y el recorrido completo en diagrama.

</details>

**Siguiente:** [Proceso siguiente](pagina.md). Una línea que dice qué recibe de esta.
```

Solo en el README de la capacidad, después de la síntesis: `## Cómo se construyó` con las decisiones, su porqué y enlaces a ADR y commits. No es un diario.

## Diagramas
Cada tipo responde una sola pregunta y la oración anterior dice qué muestra:

| Tipo | Pregunta | Forma |
|---|---|---|
| Mapa | ¿Dónde estoy? | `flowchart LR` en fila; la caja actual con `style X fill:#ffd54f,stroke:#333,stroke-width:3px,color:#000` |
| Entrada | ¿Qué llega? | `flowchart LR` con un `subgraph` por capa u origen, cada uno con `direction TB` y sus datos `[/.../]` encadenados con `~~~`; subgrafos unidos con `G ~~~ P` para fijar el orden. Si cada subgrafo tiene un solo dato, `flowchart TD` |
| Resultado | ¿Qué le pasa? | `flowchart TD` con un `subgraph` por caso, cada uno con `direction LR`, apilados con `~~~`; flechas `== gana ==>` y `-. pierde .->`; el resultado resaltado |
| Recorrido | ¿En qué orden? | `flowchart TD`; solo en el detalle |
| Decisión | ¿Cuándo sí y cuándo no? | `{...}` con finales `([...])` |

- Formas con significado fijo: `([...])` inicio o fin; `[...]` proceso; `{...}` decisión; `[/.../]` dato. Solo las cajas `[...]` exigen enlace en la lista de debajo, con la **misma etiqueta** como texto del enlace.
- Unos 10 nodos por diagrama; en Resultado, como máximo 5 bloques. Más = partir en otra página o mandar al detalle.
- Etiquetas cortas. En datos, el valor real (`central_repo: true`), no una descripción.
- Si un diagrama cruza flechas al punto de no leerse, separar por casos: un bloque por caso.
- Nunca escribir `\n` ni `\r\n` dentro de una etiqueta: Mermaid lo convierte en salto de línea. Usar palabras (`LF`, `CRLF`) y `<br/>` para cortar línea.
- Subgrafos sin flechas entre ellos salen en orden imprevisible y, con muchos datos en fila, tan anchos que la letra no se lee: encadenarlos con `~~~` y apilar sus datos en vertical.
- No usar `click` de Mermaid: GitHub no documenta soporte y la vista previa de VS Code corre en modo estricto (verificado el 2026-11-26 en docs.github.com «Creating diagrams» y code.visualstudio.com/docs/languages/markdown).
- `stateDiagram` o `sequenceDiagram` solo si la pregunta es «¿en qué estado está?» o «¿quién habla con quién?».
- Prohibido el diagrama decorativo.

## Orden de conceptos
- Un término del glosario se explica una sola vez, en `GLOSSARY.md`.
- En cada página, su **primera mención** va enlazada al glosario (`[capa](../../GLOSSARY.md)`) o en **negrita** con su explicación ahí mismo. También cuenta un enlace a la página que lo explica.
- Las etiquetas de un diagrama cuentan como mención: el término debe quedar enlazado o definido en la prosa antes del diagrama.
- El nombre del mecanismo va después de Entrada y Resultado: primero se ve, después se nombra.
- Los títulos no cuentan como mención y no llevan enlaces.

## Caso conductor
- Uno por capacidad, sin personajes. Ejemplo: «agregamos una regla en la capa de producto y compilamos».
- Se ejecuta de verdad; su página muestra comandos, salidas y archivos resultantes, marcados como ejecución real con fecha.
- Cada página de proceso toma de ese mismo caso su Entrada y su Resultado, a su profundidad.

## Estilo
- Castellano técnico neutro, trato de tú, voz activa, oraciones cortas. Sin modismos.
- Tono de referencia: afirma. Sin preguntas al lector, predicciones, ejercicios ni «ahora ya puedes».
- Literales intactos: rutas, comandos, nombres de función y campos van en `código`.
- Causa con palabras: «se diseñó así para que…», «queríamos X; ocurrió Y; significa Z».
- Texto alternativo de imágenes descriptivo.

## Comprobación final
Determinista (`go run ./tools/cmd/docscheck` → código 0):
- Enlaces relativos y anclas existen.
- Cada caja `[...]` de cada diagrama tiene su enlace en la lista de debajo; datos `[/.../]` y títulos de `subgraph` no lo necesitan.
- Primera mención de cada término del glosario enlazada o en negrita, incluidas las etiquetas de diagramas.
- Cada página de capacidad trae `sources` y `verified`, y ningún archivo de `sources` cambió después de `verified`.

Revisión humana (`docscheck` no lo ve):
- «En una frase» se entiende sin leer el resto.
- La página sigue el ritmo y el nombre aparece después de Entrada y Resultado.
- Entrada y Resultado salen del caso conductor y su salida es real.
- Los diagramas se dibujan sin error y se leen sin seguir flechas cruzadas.
- Ninguna oración pregunta al lector.
