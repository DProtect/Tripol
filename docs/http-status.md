# Estado de DIPO HTTP 2 / DIPO HTTP 2 status

| Área | Estado | Evidencia |
| --- | --- | --- |
| HTTP y HTTPS | Verificados siete verbos, autenticación, encabezados, JSON, formularios, certificados inválidos, HTTP/2 y compresión Gzip, Brotli y Zstd | [Transporte local](../internal/httptransport/http_test.go), [ejemplo nativo](../core/http_transport_test.go) |
| Persistencia de cookies | Verificada en la reapertura, con aislamiento, bloqueo, dominio, ruta, caducidad y eliminación de valores de sesión | [Pruebas de cookies y métricas](../internal/httptransport/cookies_test.go), [guía bilingüe](http.md#persistencia-y-telemetría--persistence-and-telemetry) |
| Telemetría de cliente | Peticiones, errores, cancelaciones, intentos, reutilización, bytes y latencia acumulada. CPU y RAM se miden para el proceso completo en el benchmark | [Métricas](../internal/httptransport/metrics.go) |
| Archivos y tareas HTTP | Verificadas carga multipart y descarga incremental, cancelación, destino intacto en fallos y recolección de trabajos | [Pruebas de transferencia y cancelación](../internal/httptransport/http_test.go) |
| Continuaciones de Discord | Verificadas localmente la cancelación y la edición única, encolada, de un panel. Validación conectada pendiente | [Prueba TLS local](../internal/discordtransport/http_continuation_test.go), [reconexión y Resume simulados](../internal/discordtransport/transport_test.go), [guía](discord.md#http-diferido-con-paneles--deferred-http-panels) |
| Diez comandos de ejemplo | Código verificado por el compilador. Pendientes las respuestas de proveedores y las interacciones reales en Discord | [Aplicaciones](../examples/discord_http_apps.tri), [servicios](../examples/discord_http_services.tri), [pruebas manuales](http.md#pruebas-manuales-de-discord--manual-discord-checks) |
| Colecciones de SIMPLE R3 | Verificadas listas de objetos, listas anidadas, alias tipados, comprobación de tipos y protección de la arena temporal | [Pruebas nativas y WebAssembly](../core/object_lists_test.go), [guía del lenguaje](language.md#simple-r3-collection-and-json-types--colecciones-y-tipos-json) |
| Decodificación tipada | Verificados structs, listas, cadenas, booleanos, enteros, claves duplicadas y entradas malformadas. Pendientes decimales y valores nulos | [Pruebas](../core/object_lists_test.go), [ejemplo](../examples/simple_r3_objects.tri) |
| Callbacks generales | Continúan disponibles inferencia y referencias estáticas. Pendientes capturas dinámicas con cierre y ejecución concurrente general de callbacks TRIPOL | [Guía del lenguaje](language.md) |
| Rendimiento y estabilidad | Mediciones locales de seis clientes y la base 08063b3. Prueba sostenida de 150.000 peticiones. Medición Windows y Linux en GitHub Actions | [Metodología y datos](../benchmarks/http_r2/) |
| Validación conectada | Pendiente con acceso autorizado a Discord y proveedores. Las pruebas con servidores locales cubren fallos controlados | [Pasos de validación](http.md#pruebas-manuales-de-discord--manual-discord-checks) |

## Final Polish

| Área / Area | Estado / Status | Evidencia / Evidence |
| --- | --- | --- |
| Endpoints and HTTP transport | Verified private response format and grouped I/O. The native fixture verifies newlines, quotes, headers and Unicode round-trip. Legacy frame actions remain available. | [Native transport test](../core/http_transport_test.go), [performance data](../benchmarks/http_r2/) |
| Named arguments and nullables | Verified named arguments on functions, overloads, methods and struct constructors. Tests check source-order evaluation, optional defaults and fallback. | [Compiler tests](../core/final_polish_test.go) |
| Scalar closures | Verified independent snapshots, nested captures, returning callbacks, and protection against retaining a closure allocated in a temporary arena. Mutable aggregates and general parallel callbacks remain pending. | [Lifetime tests](../core/final_polish_test.go) |
| Typed JSON | Verified defaults, optional values, nested lists and records, field validation and exact finite binary64 round trips. | [Decimal and typed JSON tests](../core/final_polish_test.go), [official example](../examples/simple_r3_polish.tri) |
| Async task protocol | Verified `async` and `await` on launcher-backed HTTP, task grouping and response collection. Arbitrary parallel execution of TRIPOL blocks and suspension inside a Discord handler remain pending. | [Native task test](../core/http_transport_test.go), [HTTP group example](../examples/http_async_group.tri) |
| Sustained workloads | One million real local GETs in about 208 seconds. Multi-hour runs and live Discord validation remain pending. | [One million results](../benchmarks/http_r2/results-polish-million-linux.json), [manual validation](http.md#pruebas-manuales-de-discord--manual-discord-checks) |
