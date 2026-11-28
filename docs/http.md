# TRIPOL v0.29 DIPO HTTP 2

## Primeros pasos / Getting started

Importa `tripol.http` para crear peticiones HTTP y HTTPS. Usa `tripol.https` si el cliente debe rechazar URLs sin TLS. `tripol run` activa el transporte del lanzador en Windows y Linux x86-64. En Windows puedes ejecutar un PE compilado con `tripol net-run archivo.exe`. En Linux, las compilaciones que enlazan el puente se empaquetan con el lanzador.

Import `tripol.http` for HTTP and HTTPS. Use `tripol.https` to require TLS. `tripol run` starts the launcher on x86-64 Windows and Linux. On Windows, run a compiled PE with `tripol net-run file.exe`. Linux builds that link the bridge include the launcher.

```tripol
import tripol.http

result = http.get("http://127.0.0.1:8765/health")
show result.status
show result.text
show result.json()

result = http.post("http://127.0.0.1:8765/echo", {
    "name": "Alex",
    "active": true
})
show result.status
```

El servidor local de prueba está en `tests/manual/http-fixture.py`. Ejecútalo con Python 3 y luego ejecuta `tripol run examples/http_requests.tri`. Las respuestas esperadas incluyen `200`, `{"ok":true}` y `201`.

Start `tests/manual/http-fixture.py` with Python 3 and run `tripol run examples/http_requests.tri`. The output includes `200`, `{"ok":true}` and `201`.

## Peticiones, respuestas y métodos / Requests, responses and methods

Usa `http.get`, `post`, `put`, `patch`, `delete`, `head` y `options`, o construye una petición con `http.request(method,url)`. `Response.status` conserva el código HTTP incluso para 4xx y 5xx. `Response.error` contiene un código de error del transporte como `timeout`, `tls_certificate`, `connection`, `invalid_header` o `body_too_large`. `Response.ok()` exige 2xx y ausencia de error. Los antiguos `http.get_raw`, `http.status`, `http.body` y `https.get(host,path)` siguen disponibles.

Use the seven method helpers or `http.request(method,url)`. HTTP errors preserve their status code. Transport failures use `Response.error`; `ok()` requires a 2xx result without a transport error. The legacy raw HTTP functions and two-argument Windows HTTPS function remain available.

## Encabezados, autenticación y formularios / Headers, authentication and forms

`request.header(name,value)` añade un encabezado. `request.query(name,value)` codifica un parámetro de consulta. `bearer(token)`, `basic(user,password)` y `api_key(name,value)` configuran credenciales. En un cliente autenticado configura primero `base_url`. El transporte exige el mismo origen en las peticiones del cliente. Lee las credenciales del entorno.

Use `header`, `query`, `bearer`, `basic` and `api_key` on a request or the supported authentication methods on a client. Authenticated clients require a base URL and enforce same-origin requests. Supply credentials through the process environment.

`request.json` acepta una cadena JSON validada o un `json.Object`. `request.form` crea un formulario codificado como URL. Al añadir `request.file(name,path)`, los campos y archivos se envían como multipart. `Response.json()` devuelve el texto JSON validado. Usa `json.get_raw_large` o las funciones de lectura de `tripol.json` para consultar los campos.

Use `request.json` for validated JSON, `request.form` for URL-encoded fields and `request.file` for a multipart body. `Response.json()` returns validated JSON text. Its fields can be read using the `tripol.json` helpers.

## Clientes, cookies, archivos y tareas / Clients, cookies, files and tasks

`http.client()` reutiliza conexiones y mantiene cookies aisladas durante la sesión. `cookies(false)` debe configurarse antes de la primera petición. `close()` libera las conexiones y cancela las tareas del cliente. Hay hasta 32 clientes y 16 tareas pendientes por programa lanzado.

A client pools connections and keeps an isolated cookie jar for its lifetime. Use `cookies(false)` before its first request. Call `close()` to release connections and cancel its jobs. A launcher supports 32 clients and 16 outstanding jobs.

`request.start()` ejecuta el trabajo en segundo plano. Usa `task.poll()` para una consulta sin bloqueo, `task.await()` para esperar y `task.cancel()` para cancelar. El resultado se recoge una vez. Las descargas y subidas multipart transmiten los bytes sin cargar el archivo entero. Pasa un callback a `request.send(callback)` para recibir el progreso. `progress.uploaded` y `progress.upload_total` cuentan los bytes de archivo, excluidos los metadatos multipart. `http.download` admite un callback con `bytes`, `total` y `percent`. Si el total se desconoce devuelve `-1` mientras dura la transferencia. La descarga conserva el destino si existe y solo lo publica cuando termina.

Use `start`, `poll`, `await` and `cancel` to manage a background operation. A result is collected once. File downloads and multipart uploads stream bytes. Pass a callback to `request.send(callback)` and read `progress.uploaded` and `progress.upload_total` for multipart file bytes. Downloads can report progress and publish the completed file without replacing an existing destination. Use binary downloads for large or binary responses. Inline text is limited to 48 KiB and the launcher response envelope to 65,535 bytes.

## HTTPS, errores y Discord / HTTPS, errors and Discord

HTTPS usa certificados de confianza del sistema y TLS 1.2 o posterior. El transporte negocia HTTP/2 por TLS y decodifica Gzip, Brotli y Zstd. Las redirecciones están desactivadas de forma predeterminada. `redirects("same-origin")` permite que GET y HEAD sigan redirecciones del mismo origen. Los reintentos configurables se aplican a GET, HEAD y OPTIONS sin cuerpo en fallos de conexión y estados transitorios. Las escrituras se envían una sola vez.

HTTPS validates certificates with system trust roots and requires TLS 1.2 or later. The transport negotiates HTTP/2 over TLS and handles Gzip, Brotli and Zstd. By default it returns redirects to the caller. The same-origin option applies to GET and HEAD. Configured retries apply to bodyless GET, HEAD and OPTIONS on transient failures.

Para integrar una API con Discord, usa `ctx.defer_http(request,"nombre")`, registra el resultado con `events.on_http` y termina con `event.complete_http`. El Gateway puede seguir atendiendo eventos mientras el lanzador ejecuta la petición. `examples/discord_http_apps.tri` muestra GitHub, clima, monitorización y descarga de un informe a un canal privado. `examples/discord_http_services.tri` conecta OpenAI Responses, LibreTranslate, GNews y Steam mediante sus credenciales de entorno. Usa un bot y servidor privado autorizados para verificarlo en Discord.

Use `ctx.defer_http`, `events.on_http` and `event.complete_http` for deferred interaction replies. The Gateway continues to poll while the launcher executes HTTP. The examples cover GitHub profiles, weather, site monitoring and report attachments. The services example adds OpenAI Responses, LibreTranslate, GNews and Steam. Connected validation requires an authorized app and private Discord server.

## Validación de las integraciones / Integration validation

En Windows o Linux, configura `DISCORD_TOKEN` y `DISCORD_GUILD_ID` para un bot y servidor de pruebas autorizados. Ejecuta `tripol run examples/discord_http_apps.tri` y prueba `/perfil`, `/tiempo` y `/estado`. Configura `TRIPOL_MONITOR_URL` con una URL HTTPS autorizada para este último. Durante una respuesta lenta, invoca otro comando y comprueba que el bot puede responder. Para `/informe` configura `TRIPOL_REPORT_URL` y `DISCORD_REPORT_CHANNEL_ID`. Verifica los bytes del archivo descargado y el adjunto del canal.

Run the applications example on Windows and Linux with an authorized bot token and private test guild. Check the profile, weather and configured monitor commands. Run a second command during a slow response. For reports, configure `TRIPOL_REPORT_URL` and `DISCORD_REPORT_CHANNEL_ID` and compare the downloaded file with the Discord attachment.

Ejecuta `tripol run examples/discord_http_services.tri`. Configura `TRIPOL_OPENAI_API_KEY` y `TRIPOL_AI_MODEL` para `/ia`, `TRIPOL_TRANSLATION_URL` y la clave opcional `TRIPOL_TRANSLATION_KEY` para `/traducir`, `TRIPOL_GNEWS_KEY` para `/titular` y `TRIPOL_STEAM_KEY` para `/jugador`. Usa cuentas de prueba y verifica la respuesta privada contra el endpoint correspondiente.

The services example reads the explicit provider variables above. Supply only authorized credentials and verify each private reply with its provider. Keep tokens, API keys and private endpoints out of logs.

## Persistencia y telemetría / Persistence and telemetry

Configura `client.base_url(url)` y `client.cookie_file(path)` antes de la primera petición. Usa un archivo distinto para cada cliente y un directorio local privado. El archivo guarda únicamente cookies con caducidad, conserva dominio y ruta y está ligado al esquema, host y puerto. El bloqueo impide abrirlo a la vez en dos clientes. En Linux usa permisos privados del usuario. En Windows el directorio debe heredar una ACL privada. Los datos se guardan en texto plano, así que protege el archivo como una credencial. `response.error` indica los errores de apertura y de escritura.

Set `client.base_url(url)` and `client.cookie_file(path)` before the first request. Give each client its own file in a private local directory. The locked file saves cookies with an explicit expiry, preserves domain and path, and is tied to the base scheme, host and port. On Windows choose a directory with a private ACL. Treat the file as a credential because it stores cookie values in plaintext.

`client.metrics()` devuelve peticiones activas y terminadas, errores, cancelaciones, intentos, conexiones nuevas y reutilizadas, bytes de respuesta, bytes de archivos multipart y latencia acumulada en microsegundos. Las conexiones cuentan asignaciones del transporte, incluidas las reutilizaciones HTTP/2. Los contadores consultados durante peticiones activas pueden avanzar entre campos.

`client.metrics()` reports in-flight and completed requests, HTTP or transport failures, cancellations, attempts, new and reused connections, downloaded bytes, multipart file bytes and cumulative elapsed microseconds. Connection counts come from transport assignments. A snapshot during active requests can contain counters observed at slightly different times.

## Pruebas manuales de Discord / Manual Discord checks

Usa una aplicación y un servidor privado autorizados. Configura las variables de la sección de integración y ejecuta ambos ejemplos con el lanzador en Windows y Linux. En `discord_http_apps.tri`, verifica `/github usuario:Tripollas`, `/pokemon nombre:ditto`, `/clima`, `/estado`, `/servidor` y `/informe`. Comprueba que el botón **Consultar Pikachu** modifica el panel existente y que el informe solo se envía al canal privado configurado. En `discord_http_services.tri`, verifica `/ia`, `/traducir`, `/noticias` y `/steam` con las claves correspondientes. Durante una API lenta, invoca `/servidor` y provoca una desconexión temporal en un entorno controlado. Comprueba que vuelve a recibir comandos, que la respuesta diferida se entrega como máximo una vez y que el registro no incluye tokens.

On both Windows and Linux, run the two example bots in an authorized private guild. Exercise the commands listed above, run a second command during a slow upstream request, and check the Pokémon refresh button edits the same message. Simulate a brief Gateway interruption and reconnect, then verify that dispatch resumes and the deferred response is edited at most once. Supply provider and Discord credentials through the environment. Keep tokens and cookie files out of logs.

## Tareas asíncronas / Asynchronous tasks

`task = async { http.get(url) }` envía el GET al lanzador. `response = await task` recoge una sola vez el resultado. `await http.get(url)` inicia la operación y espera inmediatamente en el hilo de TRIPOL. El mismo protocolo admite objetos que implementan `start()` y devuelven una tarea con `await()`. La preparación de la operación ocurre en el hilo que la solicita.

`task = async { http.get(url) }` submits a GET to the launcher. `response = await task` collects the result once. `await http.get(url)` submits and then waits on the TRIPOL thread. A library can expose a startable object with `start()` returning a task with `await()`. The calling thread still prepares the operation.

`http.group()` agrupa tareas ya iniciadas. `group.add(task)` rechaza duplicados, `group.cancel()` solicita su cancelación y `await group` devuelve las respuestas en el orden de alta. Configura los timeouts en cada petición. Los errores de red o cancelación se consultan en `Response.error`. Ejecuta [el ejemplo](../examples/http_async_group.tri) junto con el servidor local de esta guía.

An HTTP group collects submitted tasks in insertion order. Set timeouts on the requests and inspect `Response.error` for transport or cancellation failures. In Discord callbacks, use `ctx.defer_http` with a queued reply to return control to the Gateway dispatcher while the request runs. Suspending and resuming a callback at an arbitrary `await` is pending. [Discord continuation guide](discord.md#http-diferido-con-paneles--deferred-http-panels)
