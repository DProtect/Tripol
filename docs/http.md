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
