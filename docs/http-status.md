# Estado de DIPO HTTP 2 / DIPO HTTP 2 status

| Área | Estado | Evidencia |
| --- | --- | --- |
| HTTP R2 | Siete verbos, cliente, petición, respuesta, autenticación, formularios, cookies de sesión y conexión reutilizada | [Transporte local](../internal/httptransport/http_test.go), [programa nativo](../core/http_transport_test.go) |
| HTTPS R2 | Certificados válidos e inválidos, HTTP/2 sobre TLS, Gzip, Brotli, Zstd, redirects y reintentos | [Servidor TLS y compresión local](../internal/httptransport/http_test.go) |
| Archivos | Multipart y descargas binarias, archivos interrumpidos y eventos de progreso de descarga y subida | [Pruebas de archivos y descarga nativa](../internal/httptransport/http_test.go) |
| Concurrencia | Tareas HTTP en el lanzador, cancelación y recolección de resultados | [Dos peticiones simultáneas y cancelación independiente](../internal/httptransport/http_test.go) |
| Discord | Respuestas diferidas desde callbacks tipados sin bloquear el trabajo HTTP en el ciclo de eventos | [Continuación TLS local](../internal/discordtransport/http_continuation_test.go), [bot de ejemplo](../examples/discord_http_apps.tri) |
| SIMPLE R3 | Sobrecargas de funciones, registros anidados, callbacks sin resultado y literales con contexto JSON | [Casos nativo y WebAssembly](../core/simple_r3_test.go) |
| Rendimiento | GET secuencial, RAM, CPU, transferencia y reutilización de conexiones en seis clientes, con perfil previo y optimizado | [Medición, scripts y resultados](../benchmarks/http_r2/) |
| Pendiente | Comparativa de rendimiento concurrente y transferencia prolongada con todos los clientes | [Metodología secuencial](../benchmarks/http_r2/) |
| Pendiente | Colecciones generales de objetos, deserialización a tipos del usuario, captura dinámica y concurrencia de callbacks del lenguaje | [Referencia](language.md) |
| Pendiente | Cookies en disco y telemetría de recursos por cliente | [Guía](http.md) |
| Pendiente | Validación conectada en Discord y proveedores de ejemplo con credenciales autorizadas | [Integración HTTP y Discord](http.md) |
