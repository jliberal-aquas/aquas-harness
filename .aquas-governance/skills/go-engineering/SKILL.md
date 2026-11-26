---
name: go-engineering
description: Reglas y flujo para escribir, modificar, revisar o probar código Go en este arnés (archivos .go, go.mod, módulos en tools/). Cárgala antes de tocar código Go o de planificar una herramienta, hook o paquete en Go.
when_to_use: Crear un comando en cmd/, un paquete en internal/, un hook de Claude Code en Go, corregir un fallo en Go, agregar pruebas Go, refactorizar o revisar un diff de Go, subir la versión de Go del go.mod.
---

# Ingeniería en Go · v0.2
Misión: código Go simple, explícito y probado, que otro agente o humano entienda sin contexto.
No aplica a: scripts de un solo uso que no se guardan en el repositorio.
Fuentes: go.dev/doc/modules/layout · go.dev/wiki/CodeReviewComments · google.github.io/styleguide/go.

## Prioridad si dos reglas chocan
1. Correcto y probado.
2. Simple y explícito.
3. Consistente con el código que ya existe.
4. Rendimiento, solo con una medición que lo justifique.

## Flujo de trabajo
1. **Ubicar el módulo.** Buscar el `go.mod` más cercano (en este arnés: `tools/go.mod`). Leer su línea `go`: es la versión de Go permitida. No usar funciones de una versión posterior.
2. **Leer antes de escribir.** Revisar el paquete que se va a tocar y sus pruebas. Seguir sus nombres y su estilo.
3. **Ubicar el cambio** según «Estructura» y «Paquetes». Si el cambio pide una carpeta nueva, justificar en una línea por qué no cabe en una existente.
4. **Prueba primero cuando se corrige un fallo:** escribir el caso que falla, comprobar que falla, y recién entonces corregir.
5. **Implementar** con las reglas de abajo.
6. **Comprobar** con la lista final, en orden, desde la carpeta del módulo.
7. **Reportar** el resultado real de cada comprobación. Si un paso no se pudo ejecutar, decirlo con la etiqueta NO CONFIRMADO y el motivo.

## Estructura
- Un módulo por repositorio de herramientas. Crear otro módulo solo si el código se publica o versiona aparte.
- Cada ejecutable en `cmd/<nombre>/main.go`. `main` lee opciones, arma dependencias y llama al paquete; no contiene lógica.
- La lógica va en `internal/<capacidad>/`. Ningún paquete fuera de `internal/` sin un consumidor externo real.
- Empezar con la estructura mínima. Agregar carpetas cuando el código lo pida, no por plantilla. No usar `pkg/` ni copiar `golang-standards/project-layout`.
- Tras `go mod init`, comprobar la línea `go` del `go.mod`; corregir con `go get go@<versión>`.

## Paquetes y nombres
- Nombrar el paquete por lo que provee: `canon`, `hook`, `guard`. Prohibidos: `util`, `common`, `helpers`, `misc`, `types`, `models`.
- Organizar por capacidad del producto, no por capa técnica (`controllers/`, `services/`).
- No repetir el nombre del paquete en sus identificadores: `canon.Run`, no `canon.CanonRun`.
- Mayúscula = público. Exportar solo lo que otro paquete usa hoy.
- Sin getters ni setters vacíos: el campo se usa directo, o la función se llama `X()`, no `GetX()`.
- Identificadores en inglés, salvo los nombres de prueba, que describen la conducta en castellano.

## Arquitectura
- El núcleo decide con datos en memoria; el disco, la red y la terminal quedan en el borde (archivo `io.go` o paquete propio). Así el núcleo se prueba sin archivos.
- Las dependencias entran como parámetros: `Run(root string, w io.Writer)`, no escribir directo en `os.Stdout` dentro del paquete.
- Interfaces: se definen en el paquete que las usa, con solo las funciones que necesita. Quien implementa devuelve tipos concretos.
- No crear una interfaz sin un segundo uso real o una frontera real (disco, red). No crearla solo para simular en pruebas.
- Estado global: solo constantes o configuración fija. Si un valor cambia según el producto o la ejecución, pasarlo como parámetro.
- `init()` solo para registro trivial; nunca para leer archivos, red ni entorno.
- Genéricos solo cuando ya existen dos o más tipos reales que comparten el mismo código.
- Acceso a disco confinado a una carpeta: usar `os.Root` en vez de unir rutas a mano.

## Errores
- Toda función que puede fallar devuelve `error`. No descartar errores con `_`, salvo que ignorarlo sea intencional; en ese caso, un comentario en la misma línea dice por qué.
- Manejar el error primero y salir temprano. El camino normal va sin sangría extra, sin `else` tras `return`.
- Agregar contexto al envolver: `fmt.Errorf("leyendo %s: %w", ruta, err)`. Sin «failed to» ni «error al» redundantes.
- `%w` cuando quien llama necesite `errors.Is` o `errors.As`; `%v` en fronteras donde el error original no debe filtrarse.
- Nunca distinguir errores comparando su texto. Usar errores centinela (`var ErrX = errors.New(...)`) o tipos de error.
- Un paquete devuelve el error; no lo registra en el log además de devolverlo. Solo `main` decide terminar el programa.
- `panic` solo para errores de programación imposibles de manejar. Todo borde con otro proceso (por ejemplo, un hook) recupera el pánico y responde de forma controlada. Motivo: un pánico sin recuperar sale con código 2, que Claude Code interpreta como bloqueo.

## Concurrencia
- Escribir funciones síncronas que devuelven su resultado. Quien llama agrega goroutines si las necesita.
- Cada goroutine lanzada tiene un final claro: cancelación por contexto, cierre de canal o `sync.WaitGroup`. Documentar cuál.
- `context.Context` va como primer parámetro (`ctx`) de toda función que hace entrada o salida, o tarda. No guardarlo en un struct.

## Pruebas
- Toda función con lógica tiene prueba en `_test.go` del mismo paquete.
- Preferir pruebas en tabla: lista de casos con nombre, entrada y resultado esperado.
- Comparar dentro de la función `Test` con `t.Errorf("got %v, want %v", got, want)`. Sin librerías de aserciones.
- `t.Helper()` en toda función ayudante; `t.TempDir()` para archivos temporales.
- Probar con implementaciones reales (carpetas temporales, servidor de prueba). Simular solo lo que no se puede ejecutar en la prueba.
- Cada corrección de fallo agrega un caso que fallaba antes del arreglo y pasa después.

## Documentación
- Todo paquete tiene comentario `// Package x ...` que dice qué resuelve.
- Todo identificador exportado tiene comentario que empieza con su nombre.
- Comentar el porqué de lo no obvio (una protección, un límite, una decisión). No narrar lo que el código ya dice.
- Comentarios y mensajes al humano en castellano.

## Comprobación final
Ejecutar desde la carpeta del módulo, en este orden:
1. `gofmt -l .` → sin salida.
2. `go vet ./...` → sin salida.
3. `go test ./...` → todo `ok`. Con concurrencia: `go test -race ./...`.
4. `go build -o bin/ ./cmd/...` → sin errores.
5. `govulncheck ./...` → solo si cambiaron dependencias.
6. `go fix ./...` → solo tras subir la versión de Go del `go.mod`; revisar el diff antes de aceptarlo.
- No declarar «probado» sin el paso 3 ejecutado en esta sesión.
- En Windows, usar PowerShell: `;` separa comandos; `&&` solo existe desde PowerShell 7.
