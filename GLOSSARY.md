# Gobierno AQUAS

Arnés de gobernanza de Silicon Saint: reglas, configuración y skills que dirigen a los agentes, y la documentación que explica lo que el gobierno sabe hacer.

## Estructura del gobierno

**Repositorio central**:
El repositorio del que nace el arnés y donde se escribe su gobierno; se declara con `central_repo: true` en la configuración. Los demás repositorios son instancias.
_Avoid_: root repo, repo raíz, repo principal

**Instancia**:
Un repositorio que recibe el arnés instalado desde el repositorio central y solo escribe en su carpeta de producto.
_Avoid_: copia, satélite

**Capa**:
Una fuente de reglas, configuración y skills: primero gobierno, después producto; la última tiene la palabra final.
_Avoid_: nivel, fuente

**Canon**:
Las reglas escritas a mano en las capas; es lo único que se edita, nunca su versión compilada.
_Avoid_: reglas fuente, master

**Arnés**:
El programa que ejecuta al agente (por ejemplo Claude Code) y lee el canon ya compilado en su propia carpeta.
_Avoid_: harness, runtime

**Sección**:
Una subcarpeta de una capa con un tipo de contenido: `canon` (reglas), `config` (configuración) o `skills`.
_Avoid_: tipo, categoría

**Aporte**:
Lo que una capa trae para un mismo archivo de salida; un archivo compilado reúne los aportes de todas las capas en orden.
_Avoid_: parte, fragmento

**Huérfano**:
Un archivo dentro de las carpetas del arnés que ninguna capa produce; se avisa, nunca se borra solo.
_Avoid_: sobrante, basura

## Documentación de capacidades

**Capacidad**:
Algo que el gobierno sabe hacer hoy, documentado en tres partes: para qué sirve, cómo funciona y cómo se construyó. Ejemplo: la Compilación del canon.
_Avoid_: feature, funcionalidad (esa palabra nombra una unidad de trabajo del plan, no lo que se documenta)

**Proceso**:
Una caja del diagrama de flujo de una capacidad; tiene su propia página con un diagrama más detallado.
_Avoid_: paso, módulo, etapa

**Hoja**:
Un proceso que corresponde a una sola función sin bifurcaciones; ahí termina el descenso de páginas.
_Avoid_: nodo final, detalle

**Caso conductor**:
El único ejemplo de principio a fin de una capacidad, ejecutado de verdad, que reaparece en cada página a la profundidad que toca.
_Avoid_: ejemplo end-to-end, caso de uso, demo

**Unidad de explicación**:
El orden fijo con que se explica cada proceso: pregunta, ejemplo, diagrama, nombre, causa y enlace a la siguiente página.
_Avoid_: plantilla de sección

**Primera mención**:
La primera aparición de un término del glosario en una página; debe enlazar a su definición o explicarlo ahí mismo.
_Avoid_: primer uso

**Vigencia**:
La correspondencia entre una página y los archivos que documenta, en el commit en que se verificó; se pierde cuando esos archivos cambian después.
_Avoid_: actualización, frescura
