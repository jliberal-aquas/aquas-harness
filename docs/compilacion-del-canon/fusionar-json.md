---
sources:
  - tools/internal/canon/compile.go
  - tools/internal/canon/io.go
verified: 018d405
---

# Fusionar JSON

> **En una frase.** Cuando dos [capas](../../GLOSSARY.md) traen el mismo `settings.json`, `canonc` los combina en el único archivo de configuración que recibe cada [arnés](../../GLOSSARY.md): el producto gana donde chocan, salvo dentro de objetos, que se mezclan clave por clave.

## Estás aquí

```mermaid
flowchart LR
  A[Recolectar capas] --> B[Compilar archivos] --> C[Elegir arneses] --> D[Sincronizar arnés]
  B -.-> T[Unir texto]
  B -.-> J[Fusionar JSON]
  B -.-> K[Copiar intacto]
  style J fill:#ffd54f,stroke:#333,stroke-width:3px,color:#000
```

- [Recolectar capas](recolectar-capas.md): arma los grupos.
- [Compilar archivos](compilar-archivos.md): manda aquí los grupos `.json`.
- [Unir texto](unir-texto.md): el camino de los `.md` y `.txt`.
- [Fusionar JSON](fusionar-json.md): esta página.
- [Copiar intacto](compilar-archivos.md#copiar-intacto): el camino de las skills.
- [Elegir arneses](elegir-arneses.md): lee el JSON que sale de aquí.
- [Sincronizar arnés](sincronizar-arnes.md): escribe en `.claude/`.

## Lo que ya sabes

[Compilar archivos](compilar-archivos.md) mandó aquí el grupo `config/settings.json`, con dos [aportes](../../GLOSSARY.md): uno del gobierno y otro del producto. Pegarlos uno debajo del otro daría un JSON roto, así que hay que combinarlos.

## Entrada

Los dos archivos del [caso conductor](caso-conductor.md). El gobierno va primero y el producto después:

```mermaid
flowchart LR
  subgraph G["Gobierno · va primero"]
    direction TB
    g1[/"central_repo: true"/] ~~~ g2[/"shell: powershell"/] ~~~ g3[/"paths.root: ."/] ~~~ g4[/"paths.product: product"/] ~~~ g5[/"tags: [gobierno]"/] ~~~ g6[/"harness: [claude-code]"/]
  end
  subgraph P["Producto · va después"]
    direction TB
    p1[/"central_repo: false"/] ~~~ p3[/"paths.product: app"/] ~~~ p5[/"tags: [producto]"/]
  end
  G ~~~ P
```

## Resultado

Cada clave de `.claude/aquas/settings.json`, una fila por caso. A la izquierda, lo que trae cada capa; a la derecha, lo que queda:

```mermaid
flowchart TD
  subgraph K1["central_repo · valor simple: gana el producto"]
    direction LR
    a1[/"Gobierno<br/>true"/] -. pierde .-> c1[/"Resultado<br/>false"/]
    b1[/"Producto<br/>false"/] == gana ==> c1
  end
  subgraph K2["tags · lista: se reemplaza entera"]
    direction LR
    a2[/"Gobierno<br/>[gobierno]"/] -. pierde .-> c2[/"Resultado<br/>[producto]"/]
    b2[/"Producto<br/>[producto]"/] == gana ==> c2
  end
  subgraph K3["paths · objeto: se mezcla por dentro"]
    direction LR
    a3[/"Gobierno<br/>root: .<br/>product: product"/] -- aporta root --> c3[/"Resultado<br/>root: .<br/>product: app"/]
    b3[/"Producto<br/>product: app"/] == gana product ==> c3
  end
  subgraph K4["shell y harness · solo en una capa: pasan igual"]
    direction LR
    a4[/"Gobierno<br/>powershell<br/>[claude-code]"/] --> c4[/"Resultado<br/>powershell<br/>[claude-code]"/]
  end
  K1 ~~~ K2 ~~~ K3 ~~~ K4
  style c1 fill:#ffd54f,color:#000
  style c2 fill:#ffd54f,color:#000
  style c3 fill:#ffd54f,color:#000
  style c4 fill:#ffd54f,color:#000
```

El producto gana porque va después. La excepción es `paths`: es un objeto en las dos capas, así que se mezcla por dentro y `root` se conserva. Archivo real del paso 1:

```json
{
  "central_repo": false,
  "harness": [
    "claude-code"
  ],
  "paths": {
    "product": "app",
    "root": "."
  },
  "shell": "powershell",
  "tags": [
    "producto"
  ]
}
```

## Fusionar en profundidad

La regla que produce ese resultado se llama **fusionar en profundidad**. Para cada clave del producto, `canonc` mira si los dos valores son objetos:

```mermaid
flowchart TD
  K([Una clave del producto]) --> Q{¿Ambos valores son objetos?}
  Q -->|sí, como paths| M([Se mezclan por dentro, clave por clave])
  Q -->|no, como central_repo o tags| W([Gana el valor del producto entero])
```

Se diseñó así para que el producto pueda **cambiar un detalle sin reescribir todo**. En el caso conductor, el producto cambió solo `paths.product` y no tuvo que copiar `paths.root`.

## Cuándo se detiene

Si un archivo repite una clave, `canonc` no elige uno de los dos valores. Nombra el archivo culpable, sale con código 2 y no escribe nada en `.claude/`. Paso 6 del caso conductor:

```text
> canonc -root compilacion-del-canon
error: .product-governance/config/settings.json: JSON inválido: jsontext: duplicate object member name "central_repo"
código de salida: 2
```

## Síntesis

```mermaid
flowchart TD
  Q{¿Qué valor trae la clave?} -->|simple o lista| X([La última capa gana entero])
  Q -->|objeto en ambas capas| Y([Se mezcla clave por clave])
  Q -->|clave repetida en un archivo| Z([Error, no se escribe nada])
```

La última capa gana, salvo dentro de objetos, donde se mezcla. Si el JSON es ambiguo, `canonc` se detiene. Una consecuencia directa: si el producto declara `"harness"`, su lista reemplaza entera la del gobierno; el paso 7 del caso conductor lo muestra.

<details>
<summary><strong>Detalle completo</strong> (para quien va a cambiar el código)</summary>

**Funciones:** `mergeJSON` y `deepMerge` en `tools/internal/canon/compile.go`; `decodeObject` y `encodeObject` en `tools/internal/canon/io.go`.

**Qué cuenta como entrada válida**
- Un archivo vacío o con `null` cuenta como objeto vacío.
- Se quita la marca BOM antes de leer, porque algunos editores de Windows la agregan.
- Una clave repetida es error, porque `encoding/json/v2` la rechaza.
- El mensaje de error trae la ruta del archivo culpable.

**Cómo se escribe el resultado**
- Las claves van en orden alfabético (`json.Deterministic(true)`), con sangría de dos espacios y un salto de línea al final.
- Ordenar importa porque `canonc -check` compara bytes. En Go, un mapa no tiene orden fijo, así que sin ordenar el mismo JSON podría salir distinto en cada ejecución y parecer «desactualizado».

**Recorrido completo**

```mermaid
flowchart TD
  I([Aportes JSON en orden de capa]) --> L[Quitar BOM y leer objeto]
  L --> V{¿JSON válido?}
  V -->|no| E([Error con la ruta, código 2])
  V -->|sí| M[Fusionar sobre el acumulado]
  M --> N{¿Quedan aportes?}
  N -->|sí| L
  N -->|no| S[Escribir ordenado]
  S --> F([Bytes estables])
```

- [Quitar BOM y leer objeto](#cuándo-se-detiene): `decodeObject`; es donde falla una clave repetida.
- [Fusionar sobre el acumulado](#fusionar-en-profundidad): `deepMerge`.
- [Escribir ordenado](#resultado): `encodeObject`; produce el JSON que se ve en Resultado.

</details>

**Siguiente:** [Elegir arneses](elegir-arneses.md). Lee la lista `harness` de este JSON ya fusionado.
