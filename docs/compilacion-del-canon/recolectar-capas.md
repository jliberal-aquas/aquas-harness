---
sources:
  - tools/internal/canon/compile.go
verified: 018d405
---

# Recolectar capas

> **En una frase.** `canonc` recorre cada [capa](../../GLOSSARY.md) en orden y junta en un grupo los archivos que terminarán en el mismo lugar de cada [arnés](../../GLOSSARY.md), con lo del gobierno siempre antes que lo del producto.

## Estás aquí

```mermaid
flowchart LR
  A[Recolectar capas] --> B[Compilar archivos] --> C[Elegir arneses] --> D[Sincronizar arnés]
  B -.-> T[Unir texto]
  B -.-> J[Fusionar JSON]
  B -.-> K[Copiar intacto]
  style A fill:#ffd54f,stroke:#333,stroke-width:3px
```

- [Recolectar capas](recolectar-capas.md): esta página.
- [Compilar archivos](compilar-archivos.md): recibe los grupos que salen de aquí.
- [Unir texto](unir-texto.md): une reglas en Markdown.
- [Fusionar JSON](fusionar-json.md): fusiona configuración.
- [Copiar intacto](compilar-archivos.md#copiar-intacto): copia skills.
- [Elegir arneses](elegir-arneses.md): lee a qué arneses escribir.
- [Sincronizar arnés](sincronizar-arnes.md): escribe en `.claude/`.

## Lo que ya sabes

En el [caso conductor](caso-conductor.md) hay dos `rules.md` en carpetas distintas que deben terminar en un solo archivo. Antes de unirlos, `canonc` necesita saber que son «el mismo archivo» y en qué orden van. Esta es la primera caja del mapa: todavía no une nada, solo ordena.

## Entrada

Cada capa tiene hasta tres subcarpetas. A cada una la llamamos **[sección](../../GLOSSARY.md)**: `canon` (las reglas del [canon](../../GLOSSARY.md)), `config` (configuración) y `skills`. En el caso conductor:

```mermaid
flowchart LR
  subgraph G["Capa de gobierno · se lee primero"]
    direction TB
    g1[/"canon/human-communication/rules.md"/] ~~~ g0[/"canon/.borrador.md"/] ~~~ g2[/"config/settings.json"/] ~~~ g3[/"skills/saludo/SKILL.md"/]
  end
  subgraph P["Capa de producto · se lee después"]
    direction TB
    p1[/"canon/human-communication/rules.md"/] ~~~ p2[/"config/settings.json"/]
  end
  G ~~~ P
```

## Resultado

Tres grupos, uno por archivo de salida. Lo que cada capa trae para un mismo grupo es un **[aporte](../../GLOSSARY.md)**; dentro de cada grupo, el aporte del gobierno va primero. El archivo oculto no entra en ninguno:

```mermaid
flowchart TD
  subgraph K1["Grupo rules/human-communication/rules.md · dos aportes"]
    direction LR
    a1[/"1 · gobierno"/] --> c1[/"grupo"/]
    b1[/"2 · producto"/] --> c1
  end
  subgraph K2["Grupo config/settings.json · dos aportes"]
    direction LR
    a2[/"1 · gobierno"/] --> c2[/"grupo"/]
    b2[/"2 · producto"/] --> c2
  end
  subgraph K3["Grupo skills/saludo/SKILL.md · un aporte"]
    direction LR
    a3[/"1 · gobierno"/] --> c3[/"grupo"/]
  end
  subgraph K4["canon/.borrador.md · se ignora"]
    direction LR
    a4[/"empieza con punto"/] -. fuera .-> c4[/"ningún grupo"/]
  end
  K1 ~~~ K2 ~~~ K3 ~~~ K4
  style c1 fill:#ffd54f
  style c2 fill:#ffd54f
  style c3 fill:#ffd54f
```

En el paso 1 del caso conductor, la salida `(governance + product)` nombra justo esos aportes, en ese orden.

## Agrupar por clave

El mecanismo se llama **agrupar por clave**. Cada archivo recibe una clave: el tipo de su sección más su ruta dentro de ella. La sección `canon` da el tipo `rules`, así que `canon/human-communication/rules.md` de cualquier capa da la clave `rules/human-communication/rules.md`:

```mermaid
flowchart TD
  F([Un archivo de una sección]) --> O{¿Su nombre empieza con punto?}
  O -->|sí| X([Se ignora])
  O -->|no| K{¿Ya existe un grupo con su clave?}
  K -->|sí| A([Se agrega al final de ese grupo])
  K -->|no| N([Abre un grupo nuevo])
```

Se diseñó así para que el orden de los aportes salga del orden de lectura, sin ordenar nada después: como el gobierno se lee primero, su aporte siempre queda antes que el del producto.

## Cuándo se detiene

Si al terminar no hay ningún grupo, `canonc` sale con código 2. Casi siempre significa que `-root` apunta a la carpeta equivocada. Paso 8 del caso conductor:

```text
> canonc -root vacio
error: no encontré archivos en "vacio": revisa -root
código de salida: 2
```

## Síntesis

```mermaid
flowchart TD
  Q{¿Cómo es el archivo?} -->|mismo destino en dos capas| X([Un grupo con dos aportes, gobierno primero])
  Q -->|está en una sola capa| Y([Un grupo con un aporte])
  Q -->|su nombre empieza con punto| Z([No entra en ningún grupo])
```

Recolectar no combina nada: solo decide qué va junto y en qué orden. Una consecuencia directa: para que el producto agregue reglas a un archivo del gobierno, basta con crear un archivo con la misma ruta dentro de su sección.

<details>
<summary><strong>Detalle completo</strong> (para quien va a cambiar el código)</summary>

**Función:** `collect` en `tools/internal/canon/compile.go`. Las tablas `layers` y `sections` del mismo archivo fijan el orden.

**Orden de lectura:** gobierno y luego producto; dentro de cada capa, `canon`, `config` y `skills`.

**Casos límite**
- Una sección que no existe se salta sin error.
- Un archivo o una carpeta cuyo nombre empieza con punto se ignora, para no recoger archivos ocultos del editor. La regla aplica al nombre de cada archivo, no a la ruta completa.
- Cualquier otro error al leer una sección detiene todo con su ruta.
- `Run` es quien da el error de «no encontré archivos», cuando `collect` devuelve cero grupos.

**Recorrido completo**

```mermaid
flowchart TD
  I([Cada capa, cada sección]) --> D{¿Existe la carpeta?}
  D -->|no| I
  D -->|sí| W[Recorrer archivos]
  W --> O{¿Empieza con punto?}
  O -->|sí| W
  O -->|no| K[Agrupar por clave]
  K --> F([Grupos de aportes])
```

- [Recorrer archivos](#entrada): `fs.WalkDir` sobre cada sección.
- [Agrupar por clave](#agrupar-por-clave): la sección de arriba.

</details>

**Siguiente:** [Compilar archivos](compilar-archivos.md). Recibe estos tres grupos y decide cómo combinar cada uno.
