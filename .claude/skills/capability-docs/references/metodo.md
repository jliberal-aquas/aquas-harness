# Método de documentación · procedencia y decisiones

Reglas depuradas el 2026-11-26 del kit editorial «kit-libros-pliego» (método PLIEGO con pedagogía y didáctica tipo CS50). El kit está pensado para libros PDF; aquí solo queda lo que sirve para documentación técnica en Markdown. Esta copia es autocontenida: la skill no depende del kit ni de otras skills.

## Qué se tomó del kit
| Regla | Origen en el kit | Adaptación |
|---|---|---|
| Progresión en capas: problema → piezas → mecanismo → bordes | `METODO.md` §A.3 | Las capas son niveles de página: resumen → proceso → detalle → hoja. |
| Lo concreto antes del nombre: situación, imagen, frase intuitiva, formalismo | `METODO.md` §B | Unidad de explicación de 6 pasos. |
| Ningún término antes de su definición, también en rótulos | `METODO.md` Paso 4 | Primera mención enlazada o en negrita; los diagramas cuentan; los títulos no. |
| Caso de principio a fin único | `METODO.md` §A.1 | Sin personajes; ejecutado de verdad. |
| Toda figura responde a una pregunta; sin diagrama decorativo | `METODO.md` §C.6 | Igual. |
| Diagramas verticales, ~10 nodos, etiquetas ≤3 palabras | `METODO.md` §C.4 y §D | `flowchart TD`; más nodos = página hija. |
| Demostración rotulada real o sintética | `METODO.md` §B | El caso conductor siempre es real. |
| Causa con «queríamos X; ocurrió Y; significa Z» | `METODO.md` §B | Igual. |
| Controles deterministas donde se pueda | `tools/orden_conceptos.py`, `tools/controles.py` | Reescritos como `docscheck` en Go; ahora también leen diagramas. |

## Unidad de explicación (6 pasos)
1. **Pregunta** que resuelve el proceso.
2. **Ejemplo** concreto, tomado del caso conductor.
3. **Diagrama**, si el proceso tiene pasos o decisiones.
4. **Nombre** del concepto, recién ahora.
5. **Causa**: por qué funciona así.
6. **Enlace** a la página o sección siguiente.

El kit traía 11 pasos; se quitaron predicción, analogía, error fértil, contraste y ejercicio porque el lector de esta documentación observa lo construido, no practica.

## Colisiones del kit y cómo se resolvieron
- **Orden del nombre**: el publisher v1.3.3 pone el nombre antes de la demostración; `METODO.md` después. Se adopta *después*.
- **Dirección de diagramas**: `METODO.md` exige vertical; el ensayo usa horizontal. Se adopta vertical (`TD`).
- **Texto alternativo**: el kit lo exige vacío porque el pie describe la imagen. Se adopta descriptivo, por accesibilidad.
- **Error fértil** (cuatro nombres distintos): no aplica; esta documentación no describe errores.
- **Idioma**: el kit fija español de Chile; aquí castellano técnico neutro con trato de tú.

## Qué se descartó y por qué
- Todo lo de impresión: portadas, tipografías, folios, caja de 126 mm, preprensa, huecos de página, ensamblaje de volúmenes.
- Personajes y escenas narrativas: no encajan en documentación de referencia.
- Ejercicios y bloque «Hazlo tú»: el lector no ejecuta tareas.
- Dependencias del motor (MathJax, WeasyPrint, PyMuPDF): sin uso aquí, así que su vigencia no afecta.

## Andragogía
El kit no la trata como tal. Se aplica en dos bloques fijos al abrir cada página: **Para qué sirve** (relevancia) y **Qué necesitas saber antes** (respeta lo que el lector ya sabe y le dice dónde aprender lo que falta).
