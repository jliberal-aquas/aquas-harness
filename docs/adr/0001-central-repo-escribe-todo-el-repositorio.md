# El repositorio central escribe en todo el repositorio

La regla 8 limitaba la escritura a `product/` y citaba excepciones del central que no existían (`aquas.install.json`, `canon/central/`, ADR 0003 y 0004). La documentación de capacidades del gobierno no pertenece al producto, así que decidimos declarar el rol con `central_repo` en `.aquas-governance/config/settings.json`: con `true` el agente escribe en todo el repositorio, salvo lo que prohíben las reglas 9 (secretos) y 10 (publicar, push); con `false` rige solo `product/`.

## Status

accepted

## Considered Options

- `product/docs/`: no cambiaba reglas, pero mezclaba la documentación del gobierno con la del producto.
- Ampliar solo la excepción a `docs/`: resolvía hoy, pero dejaba el rol central sin respaldo en la configuración.
- Recrear los ADR 0003 y 0004: descartado, porque no se sabe qué decidían y inventarlos violaría la regla 5. La numeración empieza en 0001.

## Consequences

- La documentación vive en `docs/` de la raíz y nunca entra en la compilación: `canonc` solo recorre `canon/`, `config/` y `skills/` de cada capa.
- Hoy la regla se cumple solo si el agente la lee; no hay un gancho que la haga cumplir.
- Al instalar el arnés en una instancia, el valor deberá ser `false`; esa instalación aún no existe.
