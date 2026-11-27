# Documentación del gobierno AQUAS

Aquí se explica, de lo abstracto a lo detallado, cada [capacidad](../GLOSSARY.md) del gobierno: para qué sirve, cómo funciona y cómo se construyó. Los términos del proyecto se definen una sola vez en el [glosario](../GLOSSARY.md); las decisiones difíciles de revertir están en [`adr/`](adr/).

## Capacidades

| Capacidad | Qué hace |
|---|---|
| [Compilación del canon](compilacion-del-canon/README.md) | Une las reglas, la configuración y las skills de cada [capa](../GLOSSARY.md) y las escribe donde las lee cada [arnés](../GLOSSARY.md). |

## Cómo se mantiene
Estas páginas las escribe la skill `capability-docs` y las revisa `go run ./tools/cmd/docscheck` desde la raíz: enlaces, cajas de diagramas sin enlace, términos usados antes de explicarse y páginas desactualizadas respecto de su código.
