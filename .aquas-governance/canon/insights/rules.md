# Insights: aprender y evitar reincidencia · v0.3
Misión: comunicar conocimiento nuevo y protección comprobable. Sin elogios ni frases genéricas.

## Qué cuenta como insight
- Cambia cómo detectamos, prevenimos o explicamos algo. «Aprendimos mucho» no cuenta.
- Puede nacer de un fallo, una prueba exitosa, un límite o un patrón observado.
- Sin aprendizaje nuevo: «Sin insight nuevo». No inventarlo.

## Forma completa (en el cuerpo, solo si cambia un control o una decisión)
Orden fijo: HECHO · CAUSA · LECCIÓN · CONTROL · EVIDENCIA · ESTADO.
- HECHO: qué ocurrió, alcance y qué hizo mal el agente, si aplica. Sin excusas.
- CAUSA: con etiqueta de certeza de 01. Sin prueba: NO CONFIRMADO.
- Reincidencia: si ya existía un control, por qué no detectó o no impidió este caso.
- LECCIÓN: regla específica; cuándo aplica y cuándo no.
- CONTROL: en castellano, qué comprueba y qué error detecta o impide.
  - Preferir determinista: condición precisa, resultado esperado y prueba reproducible.
  - Si solo cabe probabilístico: señal observada, umbral, ventana de tiempo, acción ante alerta y riesgo de falsos positivos y falsos negativos.
- EVIDENCIA: prueba ideal = un caso que debe fallar y uno que debe pasar, con resultados reales. Indicar archivo, prueba o registro donde vive el control.
- ESTADO (uno): IMPLEMENTADO Y PROBADO · IMPLEMENTADO SIN PROBAR · PROPUESTO · SIN CONTROL (insight solo explicativo).

## Límites de lo que se afirma
- «Propuse una prueba» no equivale a «no volverá a pasar».
- Solo IMPLEMENTADO Y PROBADO permite decir que el control funciona, y solo para los casos probados.
- Sin registro persistente, no afirmar que otros agentes o sesiones heredaron el aprendizaje.

## Forma breve (campo 💡 del cierre)
- Tras un fallo: «Aprendimos [causa o límite]. Para evitar [fallo], [control]. Evidencia: [prueba/enlace o pendiente]. Estado: [estado].»
- Tras un éxito o un límite: «Aprendimos [hecho]. Desde ahora, [cambio en cómo trabajamos]. Evidencia: [prueba/enlace o pendiente]. Estado: [estado].»