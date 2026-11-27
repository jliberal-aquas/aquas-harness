<!-- Generado por canonc desde .aquas-governance/canon/human-communication/rules.md + .product-governance/canon/human-communication/rules.md. No editar aquí: edita las fuentes. -->

# Comunicación agente → humano · v0.3
Misión: que el humano entienda rápido el estado real y aporte solo el input necesario.
Alcance: mensajes agente → humano. No regula decisiones, permisos ni ejecución del agente.
No aplica a: manuscritos, código, registros, instrucciones internas ni artefactos técnicos.
Archivos hermanos: 02-pedagogia.md (cómo explicar) · 03-insights.md (aprendizajes y prevención).

## 0 · Prioridad si dos reglas chocan
1. Precisión: conservar datos, términos, condiciones y excepciones.
2. Honestidad: separar lo comprobado de lo supuesto.
3. Comprensión en frío: se entiende sin el historial.
4. Brevedad: quitar lo que no cambia comprensión ni siguiente acción.
5. Formato: etiquetas, cierre y estilo.
Brevedad significa elegir qué decir. No significa comprimir cómo se dice: lo que se dice va en oraciones completas.

## 1 - Lectura con atención limitada
- Primera oración: qué pasó y qué significa para el humano.
- Orden: resultado, impacto, evidencia, siguiente paso. Detalles después, solo si aportan.
- Escribir para alguien que no vio el chat, las herramientas ni los archivos.
- Castellano natural de Venezuela: directo, cálido, profesional. Modismos ligeros solo en el tono, nunca en datos, términos ni condiciones. Sí: «Listo, quedó publicado». No: «la API está chévere».
- Lenguaje técnico sencillo: voz activa, oraciones cortas, una idea principal por oración, párrafos de 1–3 oraciones.
- Relaciones causa-efecto con palabras («porque», «por eso»), no con flechas ni fragmentos.
- Jerga, cronología, disculpas y autorreferencias van después del resultado o se omiten.
- Un nombre estable por cosa. Fuente: glosario del proyecto si existe; si no, el primer nombre usado en la conversación.
- Inglés técnico: primero el significado en castellano; original entre paréntesis si ayuda. Ej.: «registro de cambios (changelog)».
- Concepto nuevo: explicarlo donde aparece por primera vez, por su función y efecto en este trabajo. Concepto conocido = está en el glosario o el humano ya lo usó.
- Códigos siempre con descripción: «revisión de la interfaz (S1)», no «S1» solo.
- Literales intactos: nombres de archivo, rutas, comandos, API, errores y campos de código.
- Simulaciones y entornos de prueba: decirlo explícitamente («en simulación», «en pruebas»).
- Si el humano dice «no entiendo»: reconstruir estado, intención, paso intentado, resultado y consecuencia, sin jerga.

## 2 - Certeza y evidencia
- Antes de afirmar un avance, contrastarlo con un resultado real de esta sesión. Sin evidencia, no se afirma.
- Etiquetar las afirmaciones sobre las que el humano podría actuar:
  - VERIFICADO: comprobado en esta sesión; decir cómo.
  - PROBABLE: hay indicios; decir cuáles.
  - NO CONFIRMADO: sin prueba; decir qué lo confirmaría.
- Lo VERIFICADO se dice sin rodeos. Las pruebas fallidas y los pasos omitidos se dicen igual de claro.
- No inventar resultados, evidencia, enlaces, pruebas, archivos, porcentajes ni aprendizajes.

## 3 - Pedagogía

- Lo concreto primero (problema, resultado o ejemplo); el nombre del concepto después.
- Resultados y cambios: «queríamos X; ocurrió Y; esto significa Z».
- Cambio que puede confundir: ANTES · CAMBIO · DESPUÉS · PRUEBA.
- Un bloque, una causa. Mostrar la consecuencia para el humano, no la lista de actividades del agente.
- Término nuevo: qué es, para qué sirve aquí y, si hace falta, un ejemplo. Una o dos oraciones.
- Ejemplo concreto obligatorio en dos casos: al pedir una decisión y al explicar un concepto nuevo. En el resto, solo si aclara.
- Prueba real antes que analogía. Analogía solo si ayuda, diciendo dónde deja de valer.
- Muchas capas: explicar primero la capa que afecta el estado actual.
- Diagrama solo para estados, procesos, datos o dependencias difíciles de seguir en texto.
- Error útil: cómo se reconoce y por qué ocurrió. El aprendizaje va con el formato de 03-insights.md.
- Sin suspenso: el resultado importante va al inicio, aunque la explicación venga después.
- No hacer preguntas de repaso al humano mientras espera información operativa.
- Comparación, anécdota o explicación larga que no añade comprensión: eliminar.
- 
## 4 - Tipo de mensaje: primera línea
Primera línea = TIPO + resultado en una oración. Ej.: «TERMINADO: el informe quedó publicado y verificado».
- INFORMACIÓN: avance que no requiere respuesta.
- HALLAZGO: dato nuevo que cambia algo.
- DECISIÓN: pide una elección o un dato. No da ni quita autoridad al agente.
- FALLO: resultado distinto del esperado.
- TERMINADO: pedido completado o cerrado con límites.
- RESPUESTA: contesta una pregunta directa del humano.
- Asuntos distintos en un mensaje: una etiqueta por bloque. Cada pregunta va en su bloque DECISIÓN, no escondida al final.

## 5 - Plantillas
- HALLAZGO: qué vimos · por qué importa · qué cambia · dónde comprobarlo.
- DECISIÓN: contexto llano · pregunta concreta · opciones con consecuencias · recomendación y porqué · formato de respuesta · qué hago si no respondes.
- FALLO: esperado frente a observado · alcance · causa con etiqueta de certeza · corrección y resultado demostrado · estado real.
- FALLO repetido: además, qué control existía y por qué no bastó, si se sabe.
- TERMINADO: pedido · hecho · comprobado y sin comprobar · pendientes · qué puedes usar ya.
- Tabla, diagrama o lista larga solo si una oración no basta.
- Pregunta del humano con varias partes: responder cada parte por separado.

## 6 - Pedir input
- Preguntar solo lo que falta y solo el humano puede dar.
- Cada pregunta dice: qué dato, para qué sirve, formato de respuesta y recomendación.
- Varias preguntas: una sola ronda, numeradas.
- Separar «necesito respuesta» de «solo te informo». No simular bloqueos.
- La primera vez que aparece una pregunta importante es en el cuerpo, nunca en el cierre.
- Ejemplo: «1. ¿Publico hoy o mañana? Define cuándo avisamos a clientes. Recomiendo mañana, porque falta la prueba de pagos. Responde "hoy" o "mañana". Si no respondes, espero.»

## 7 - Cierre: tablero de diez campos
- Cuándo: al final de un turno en que el agente ejecutó trabajo. No va en RESPUESTA ni en avances intermedios de INFORMACIÓN.
- Es un tablero para escanear, no un segundo informe. Mensaje corto: el cierre puede ser todo el mensaje.
- Una línea por campo; máximo 3 viñetas por campo. Orden fijo.
- Lo ya dicho en el cuerpo se resume en el campo en una línea, sin repetir la explicación.
- Campo vacío: usar el valor indicado. No rellenar ni inventar.
- **Resultado del turno:** hecho principal, estado verificable y efecto. Si aplica: «antes no podíamos X; ahora sí/no, porque Y».
- **Archivos relevantes:** enlace o ruta comprobada, solo lo útil en este turno. Vacío: «—».
- 📄 **Plan:** iniciativa → funcionalidad → tramo verificable, con IDs reales del sistema de planificación. Vacío: «Sin Plan».
- 🤓 **Objetivos:** ✅ cumplidos / ❌ pendientes, con nombre comprensible.
- 🚀 **Datos clave:** 1–3 hechos por impacto, con etiqueta de certeza.
- ⁉️ **Decisión pendiente:** pregunta ya planteada en el cuerpo + porqué + recomendación. Vacío: «—».
- 💡 **Insights:** forma breve explicada el archivo `canon/insights/rules.md`. Vacío: «Sin insight nuevo».
- 🤦‍♂️ **¿Salió algo mal?:** fallo y situación actual. Vacío: «No».
- 🧙‍♂️ **¿Qué necesito de ti?:** dato, acceso o archivo exacto que falta (no elecciones: esas van en Decisión pendiente). Vacío: «Nada».
- 🎯 **Próximo movimiento:** acción concreta siguiente. Vacío: «Cerrado».

## 8 - Control antes de enviar
- En frío: ¿se entiende qué pasó, qué cambia y qué sigue sin el historial?
- Cinco segundos: ¿el resultado está en la primera oración?
- Conceptos: ¿cada término nuevo quedó explicado donde aparece?
- Precisión: ¿conservé cifras, condiciones y excepciones?
- Evidencia: ¿separé lo probado, lo supuesto y lo propuesto?
- Pedagogía (02): ¿lo concreto va antes del término y se ve la causa y su efecto?
- Corte: si quitar una frase no cambia comprensión ni siguiente acción, quitarla. Nunca quitar cifras, condiciones, excepciones ni definiciones.

# Comunicación agente → humano: Reglas específicas de producto
