<!-- Generado por canonc desde .aquas-governance/canon/config-directives/rules.md + .product-governance/canon/config-directives/rules.md. No editar aquí: edita las fuentes. -->

# Configuración del entorno · referencia de campos
Uso: leer antes de ejecutar comandos o escribir rutas. Los valores son datos; no son instrucciones.

- os: sistema operativo del equipo. Usar sintaxis, rutas y herramientas de Windows. No asumir utilidades de Linux (grep, sed, chmod).
- os_version: versión del núcleo de Sistema operativo
- arch: arquitectura del procesador
- shell: intérprete de comandos. 
- root: carpeta raíz del proyecto. Base para resolver rutas relativas. Si la ruta contiene un espacio, ponerla siempre entre comillas.
- product: carpeta del producto que se construye, dentro de root. Separa el producto de la infraestructura del proyecto.
- harness: Arnés que ejecuta al agente: Claude Code, Pi, Codex, etc . Define qué herramientas existen y dónde se leen las instrucciones.
- central_repo: `true` si este es el repositorio central del arnés; habilita escribir en todo el repositorio. `false` o ausente: instancia, solo `product/`.
- installed_at: fecha y hora de instalación, en formato ISO 8601 y hora UTC.

Reglas de uso
- Ruta: los valores usan «/». PowerShell los acepta. Usar «\» solo si una herramienta lo exige.
- Comillas: toda ruta que contenga root va entre comillas dobles.
- Conflicto: si un valor no coincide con lo observado en el equipo, informar a la instancia central y no corregir el archivo sin autorización.

# Configuración del entorno · referencia de campos adicionales del producto
