# Método de documentación · procedencia y decisiones

Reglas depuradas el 2026-11-26 del kit editorial «kit-libros-pliego» (método PLIEGO con pedagogía y didáctica tipo CS50). El kit está pensado para libros PDF; aquí solo queda lo que sirve para documentación técnica en Markdown. Esta copia es autocontenida: la skill no depende del kit ni de otras skills.

## Qué se tomó del kit
| Regla | Origen en el kit | Adaptación |
|---|---|---|
| Progresión en capas: problema → piezas → mecanismo → bordes | `METODO.md` §A.3 | Las capas son niveles de página: resumen → proceso → detalle → hoja. |
| Lo concreto antes del nombre: situación, imagen, frase intuitiva, formalismo | `METODO.md` §B | Ritmo de página: Entrada y Resultado antes del nombre. |
| Ningún término antes de su definición, también en rótulos | `METODO.md` Paso 4 | Primera mención enlazada o en negrita; los diagramas cuentan; los títulos no. |
| Caso de principio a fin único | `METODO.md` §A.1 | Sin personajes; ejecutado de verdad. |
| Toda figura responde a una pregunta; sin diagrama decorativo | `METODO.md` §C.6 | Igual. |
| Diagramas verticales, ~10 nodos, etiquetas ≤3 palabras | `METODO.md` §C.4 y §D | `flowchart TD`; más nodos = página hija. |
| Demostración rotulada real o sintética | `METODO.md` §B | El caso conductor siempre es real. |
| Causa con «queríamos X; ocurrió Y; significa Z» | `METODO.md` §B | Igual. |
| Controles deterministas donde se pueda | `tools/orden_conceptos.py`, `tools/controles.py` | Reescritos como `docscheck` en Go; ahora también leen diagramas. |

## Unidad de explicación · v0.2 (ritmo de página)
1. **En una frase**: el resultado, para quien solo lee una línea.
2. **Estás aquí**: el mapa de nivel 0 con la caja actual resaltada. Da el todo antes de la parte.
3. **Lo que ya sabes**: el contexto escrito en la página, retomando la anterior.
4. **Entrada** y 5. **Resultado**: lo concreto, en diagramas de datos del caso conductor.
6. **Nombre y causa**: el concepto, recién después de verlo funcionar.
7. **Cuándo se detiene**: el contraste, con errores reales hechos a propósito.
8. **Síntesis**: diagrama-resumen y el modelo en una oración.
9. **Detalle completo** plegado y **Siguiente**: la granularidad la elige el lector sin cambiar de página.

Origen del cambio (2026-11-27): la v0.1 (pregunta → ejemplo → diagrama → nombre → causa → enlace) funcionaba como árbol de consulta. Repetía una plantilla en cada página, mandaba el contexto a otras páginas y sus diagramas mostraban los `if` del código, no lo que le pasa a los datos. Se probó un prototipo con predicción y preguntas al lector; el humano lo descartó porque esto es un manual de referencia, no de entrenamiento. Se conservaron el mapa, el contexto escrito, los diagramas de entrada y resultado, el contraste y la síntesis.

Del kit se recuperan el contraste (como «Cuándo se detiene») y la síntesis. Siguen fuera la predicción, la analogía y el ejercicio, porque el lector consulta lo construido, no practica. Un manual de entrenamiento queda para cuando el producto esté listo.

## Colisiones del kit y cómo se resolvieron
- **Orden del nombre**: el publisher v1.3.3 pone el nombre antes de la demostración; `METODO.md` después. Se adopta *después*.
- **Dirección de diagramas**: `METODO.md` exige vertical; el ensayo usa horizontal. Se adopta vertical (`TD`), salvo el mapa, que va en fila (`LR`), y los bloques de Resultado, que van en fila dentro de un apilado vertical.
- **Texto alternativo**: el kit lo exige vacío porque el pie describe la imagen. Se adopta descriptivo, por accesibilidad.
- **Error fértil** (cuatro nombres distintos): se adopta un solo nombre, «Cuándo se detiene», y solo para los límites que la capacidad hace a propósito, con su salida real.
- **Idioma**: el kit fija español de Chile; aquí castellano técnico neutro con trato de tú.

## Qué se descartó y por qué
- Todo lo de impresión: portadas, tipografías, folios, caja de 126 mm, preprensa, huecos de página, ensamblaje de volúmenes.
- Personajes y escenas narrativas: no encajan en documentación de referencia.
- Ejercicios y bloque «Hazlo tú»: el lector no ejecuta tareas.
- Dependencias del motor (MathJax, WeasyPrint, PyMuPDF): sin uso aquí, así que su vigencia no afecta.

## Andragogía
El kit no la trata como tal. Se aplica en dos bloques fijos al abrir cada página: **Estás aquí** (ubica la parte dentro del todo) y **Lo que ya sabes** (respeta lo que el lector ya sabe y le escribe ahí mismo el contexto, sin mandarlo a otra página). El **Detalle completo** plegado deja que cada lector elija su profundidad.
