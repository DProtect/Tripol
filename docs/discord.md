# DISCORD REVOLUTION — TRIPOL 0.28 development

## Run an offline example / Prueba un ejemplo sin conexión

```tripol
import tripol.discord

bot = discord.bot()
bot.command("hola", "¡Hola, compañero!")
show bot.preview("hola")

embed = discord.embed()
embed.title("Hola")
embed.description("Creado con TRIPOL")
embed.color("#00A0FF")
embed.field("Versión", "0.28.0", true)
show embed.payload()
```

Run `tripol run examples/discord_new.tri` on Linux or Windows, or open the Discord library examples in the Playground. Neither offline operation uses credentials. On a new bot the prefix is empty. Set `bot.prefix("!")` to match `!hola`. Duplicate command names replace their previous response.

Ejecuta `tripol run examples/discord_new.tri` en Linux o Windows, o abre los ejemplos de la biblioteca Discord en el Playground. Ambas operaciones funcionan sin credenciales. El prefijo inicial es vacío. Usa `bot.prefix("!")` para responder a `!hola`. Un comando repetido sustituye la respuesta anterior.

## Connect / Conecta

Set `DISCORD_TOKEN` before starting a bot in an authorized private server. For text commands and message callbacks, enable Message Content Intent. Use `tripol new MiBot --template discord` and `tripol run` on Windows or Linux x86-64.

Configura `DISCORD_TOKEN` antes de iniciar el bot en un servidor privado autorizado. Para comandos de texto y callbacks de mensajes, activa Message Content Intent. Usa `tripol new MiBot --template discord` y `tripol run` en Windows o Linux x86-64.

On Linux x86-64, `tripol build` packages the native program together with the TRIPOL launcher when it uses Discord networking. The launcher extracts the inner ELF into a private temporary directory and uses two inherited pipes for transport. It uses system trust roots for TLS. A raw Linux ELF cross-compiled from Windows can be launched on Linux with `tripol discord-run FILE`. A Linux bundle includes the compiler launcher, so it is larger than the inner ELF.

En Linux x86-64, `tripol build` empaqueta el programa nativo con el lanzador de TRIPOL si usa la red de Discord. El lanzador extrae el ELF a un directorio temporal privado y usa dos tuberías heredadas para el transporte. TLS usa los certificados de confianza del sistema. Un ELF generado desde Windows se ejecuta en Linux con `tripol discord-run FILE`. El paquete de Linux incluye el lanzador del compilador y ocupa más que el ELF interior.

Linux TLS integration tests cover HTTPS reuse, HTTP 429, discovery, Resume, heartbeats and fatal Gateway close codes. Windows uses its WinHTTP backend. Live behavior against Discord and automatic recovery on Windows remain pending.

En Linux, las pruebas de integración TLS cubren reutilización HTTPS, HTTP 429, descubrimiento, Resume, heartbeats y cierres definitivos del Gateway. Windows usa WinHTTP. Quedan pendientes la validación contra Discord y la recuperación automática en Windows.

## Message callbacks / Callbacks de mensajes

```tripol
import tripol.discord

func answer(message:discord.Message)->string {
    if message.matches("hola") {
        return message.reply("¡Hola!")
    }
    return ""
}

bot = discord.bot()
bot.on_message(&answer)
show bot.preview_event("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"hola\",\"author\":{\"id\":\"456\"}}}")
```

Use `bot.start()` for the Gateway. Fixed text commands have precedence over the callback. Returning an empty string skips the reply. The context exposes `author_id`, `channel_id`, `content_json` and `event_json`. Bot and webhook messages are ignored.

Usa `bot.start()` para conectarte al Gateway. Los comandos de texto fijos tienen prioridad. Devuelve un texto vacío para omitir la respuesta. El contexto expone `author_id`, `channel_id`, `content_json` y `event_json`. Se omiten los mensajes de bots y webhooks.

## Slash commands and interactions / Comandos slash e interacciones

```tripol
import tripol.discord

func hello(ctx:discord.Interaction)->string {
    return ctx.reply("Hola, " + ctx.option_text("name"), true)
}

options=discord.slash_options()
options.add("name", "Tu nombre", "string", true)
bot=discord.bot()
bot.slash_with_options("hello", "Saluda", &hello, options, "123")
show bot.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"type\":2,\"data\":{\"type\":1,\"name\":\"hello\",\"options\":[{\"type\":3,\"name\":\"name\",\"value\":\"Ana\"}]}}}")
```

Use `bot.slash(name, description, &handler, guild_id="")` for a command without options. An empty guild ID registers globally. A guild ID limits registration to a test server. At startup the bot reads its application ID, compares remote commands with the local definitions and creates or updates the changed ones. Registration leaves unrelated commands in place. This version supports ASCII lowercase command names and flat string, integer, boolean, user, channel, role, mentionable, number and attachment options.

Usa `bot.slash(nombre, descripción, &manejador, guild_id="")` para un comando sin opciones. Un ID de servidor vacío hace el registro global. Con un ID, el registro se limita al servidor de pruebas. Al iniciar, el bot consulta su aplicación, compara los comandos y crea o actualiza los que cambiaron. El registro conserva los demás comandos. Se admiten nombres ASCII en minúsculas y opciones planas de texto, entero, booleano, usuario, canal, rol, mencionable, número y adjunto.

Register component and modal callbacks with `bot.on_button(id, &handler)`, `bot.on_select(id, &handler)` and `bot.on_modal(id, &handler)`. The handler returns an initial response such as `ctx.reply(text, true)`, `ctx.defer(true)`, `ctx.ack()`, `ctx.reply_components(layout)`, `ctx.update_components(layout)`, `ctx.reply_embeds(batch)` or `ctx.open_modal(form)`. Discord requires an initial response within three seconds. Edit and follow-up endpoints after a deferred response are pending.

Registra callbacks con `bot.on_button(id, &manejador)`, `bot.on_select(id, &manejador)` y `bot.on_modal(id, &manejador)`. El manejador devuelve la respuesta inicial. `ctx.reply(texto, true)` responde en privado, `ctx.defer(true)` difiere la respuesta privada, y `ctx.ack()` confirma una interacción sobre un mensaje. Los métodos de componentes, embeds y formularios devuelven sus respectivos payloads. Discord exige una respuesta inicial en tres segundos. Quedan pendientes la edición y los mensajes de seguimiento posteriores.

`ctx.option(name)` returns raw JSON, `ctx.option_text(name)` decodes a JSON string, `ctx.values()` returns selected values and `ctx.field(id)` reads a submitted text input. `ctx.field_data(id)` returns the JSON array for select and file-upload inputs. The original interaction data is available as `ctx.data_json`, including resolved attachment metadata.

`ctx.option(nombre)` devuelve JSON, `ctx.option_text(nombre)` decodifica una cadena, `ctx.values()` devuelve los valores seleccionados y `ctx.field(id)` lee un texto enviado. `ctx.field_data(id)` devuelve el arreglo JSON de los selectores y campos de archivos. `ctx.data_json` contiene los datos originales de la interacción y los metadatos resueltos de adjuntos.

Use `tripol new Feedback --template discord-slash` for a runnable example. Set `DISCORD_TOKEN`, optionally set `DISCORD_GUILD_ID`, and install the bot application with the `applications.commands` scope on a private test server.

Usa `tripol new Feedback --template discord-slash` para generar un ejemplo ejecutable. Configura `DISCORD_TOKEN`, opcionalmente `DISCORD_GUILD_ID`, e instala la aplicación con el alcance `applications.commands` en un servidor privado.

## Components V2

```tripol
import tripol.discord

panel = discord.components()
panel.text("## Bienvenido")
panel.separator()
panel.link("Documentación", "https://tripollas.github.io/Tripol/docs/")
message = discord.components()
message.container(panel, 41215)
show message.payload()
```

The payload includes flag 32768 and disables mentions. The builder also has `button(label, custom_id, style=1)`, `section(content, thumbnail_url, description="")` and `gallery(urls)`. It validates the 40-component limit, duplicate button IDs, container nesting and field lengths. Text Display content can total 4,000 characters across the message. The encoded body must fit in the TRIPOL string capacity of 65,535 bytes. Use `message.send(channel_id)` on Windows or Linux. Use `discord.button_row()` to group up to five buttons. `discord.select(id, kind)` supports string, user, role, mentionable and channel menus. `discord.modal(id, title)` supports text inputs, selects, help text and file-upload inputs. Multipart file sending, radio groups and checkbox components are pending.

El JSON incluye la bandera 32768 y desactiva las menciones. El constructor también ofrece `button(label, custom_id, style=1)`, `section(content, thumbnail_url, description="")` y `gallery(urls)`. Valida el límite de 40 componentes, los IDs de botones duplicados, el anidamiento y las longitudes de los campos. El contenido de Text Display admite 4.000 caracteres por mensaje. El JSON debe caber en la capacidad de cadena de TRIPOL de 65.535 bytes. Usa `message.send(channel_id)` en Windows o Linux. Usa `discord.button_row()` para agrupar hasta cinco botones. `discord.select(id, tipo)` admite texto, usuario, rol, mencionable y canal. `discord.modal(id, título)` admite texto, selectores, ayuda y campos de archivos. Quedan pendientes el envío multipart, radio groups y casillas.

## Discord Studio

```sh
tripol discord studio main.tri --metrics studio.jsonl
```

On Linux x86-64, Studio runs the bot and samples its RSS, CPU ticks, connection state, heartbeat latency, uptime, guild count at READY, queued dispatches, REST requests, errors and rate limits. The launcher and bot process are reported separately. The output file must be new. Without `--metrics`, snapshots are sent to standard error. The export contains counters and process measurements. Live command and interaction metrics, log filtering and Windows monitoring are pending.

En Linux x86-64, Studio ejecuta el bot y toma muestras de RSS, ticks de CPU, conexión, latencia de heartbeat, tiempo activo, servidores en READY, cola de eventos, peticiones REST, errores y límites. El lanzador y el bot se miden por separado. Usa un archivo de salida nuevo. Sin `--metrics`, las muestras se escriben en la salida de error. Quedan pendientes las métricas por comando e interacción, el filtrado de logs y el monitoreo en Windows.

## API

| Method / Método | Behavior / Comportamiento |
| --- | --- |
| `discord.bot()` | Create a bot configuration / Crea una configuración de bot |
| `bot.prefix(value)` | Set a prefix of up to 32 UTF-8 bytes / Define un prefijo de hasta 32 bytes UTF-8 |
| `bot.status(value)` | Set `online`, `idle`, `dnd` or `invisible` / Define el estado indicado |
| `bot.activity(value)` | Set the activity text, up to 128 UTF-8 bytes / Define una actividad de hasta 128 bytes UTF-8 |
| `bot.command(name, response)` | Set one of up to 16 exact text commands / Configura uno de hasta 16 comandos de texto exactos |
| `bot.preview(content)` | Return the exact reply or an empty string / Devuelve la respuesta exacta o texto vacío |
| `bot.start()` | Open a blocking Gateway session / Abre una sesión Gateway bloqueante |
| `discord.embed()` | Create an empty embed / Crea un embed vacío |
| `embed.title(value)` | Set the title / Define el título |
| `embed.description(value)` | Set the description / Define la descripción |
| `embed.footer(value)` | Set the footer text / Define el pie |
| `embed.url(value)`, `embed.image(value)`, `embed.thumbnail(value)` | Set the title link, image or thumbnail / Define el enlace, la imagen o la miniatura |
| `embed.color(value)` | Set an RGB color as `#RRGGBB` / Define un color RGB como `#RRGGBB` |
| `embed.field(label, value, inline=false)` | Append one of up to 25 fields / Añade uno de hasta 25 campos |
| `embed.payload()` | Build JSON including disabled mentions / Genera el JSON con las menciones desactivadas |
| `embed.send(channel_id)` | Send the embed over Discord REST from Windows or Linux / Envía el embed por Discord REST desde Windows o Linux |

Text commands accept names of 1–96 bytes and replies of 1–512 bytes, both UTF-8. Use `embed.author(name, url="", icon_url="")` and `embed.timestamp(rfc3339)` for authors and timestamps. `discord.embeds()` groups up to 10 embeds with 6,000 characters of combined text. Individual limits are 256 for title, author and field name, 4,096 for description, 2,048 for footer, and 1,024 for field value. The optional `intents` parameter on legacy `discord.run()` accepts `33281` (default for guild prefix commands) or `1` (GUILDS only).

Los comandos aceptan nombres de 1–96 bytes y respuestas de 1–512 bytes, ambos en UTF-8. Usa `embed.author(nombre, url="", icon_url="")` y `embed.timestamp(rfc3339)` para autor y fecha. `discord.embeds()` agrupa hasta 10 embeds con 6.000 caracteres de texto combinado. Los límites son 256 para título, autor y nombre de campo, 4.096 para descripción, 2.048 para pie y 1.024 para valor de campo. El parámetro opcional `intents` de `discord.run()` acepta `33281` (predeterminado para comandos de texto) o `1` (solo GUILDS).

## Discord R5 — Windows WinHTTP backend

`discord.run()` now uses a persistent TLS WebSocket to `gateway.discord.gg/?v=10&encoding=json`, not REST polling. The native WinHTTP upgrade, send and receive calls are compiled directly into the executable. A native receiver thread reads one complete UTF-8 message (including fragments) into bounded private storage. The main thread waits locally, wakes on message arrival and schedules Gateway heartbeats; no 2.5-second network polling. Sending messages still uses Discord HTTPS REST POST. No Python, Node.js or external Discord SDK is required by the generated executable.

### Create a legacy bot

Create a bot at https://discord.com/developers/applications, enable Message Content Intent for prefix commands, install it on a small private test server and give it View Channel and Send Messages permissions. Read Message History is needed for the explicit REST fallback. A channel ID is optional: only select one if you want to restrict replies or send a greeting before any command. Use a bot token, not a user token.

```tripol
import tripol.discord

discord.login_env()
discord.activity("Programando en TRIPOL")
discord.reply("!hola", "Hola compa!")
discord.reply("!ping", "Pong!")
discord.run()
```

`setup` simply combines `login` and `channel`. These older methods continue to work. The optional `welcome` sends a greeting once on READY; `say` sends immediately by REST. `status("online")` is optional (default online); `idle`, `dnd`, `invisible` also supported in the initial Identify payload. `activity` sets a playing activity, up to 128 UTF-8 bytes. Up to 16 exact commands across accessible guild channels (or one selected channel); duplicate registration replaces the response. Bots and webhooks are ignored. Prefix commands match decoded text exactly, including UTF-8 characters represented with JSON escapes.

For sharing code or binaries, do not embed a token. Instead use:

```powershell
$secure = Read-Host 'Bot token' -AsSecureString
$env:DISCORD_TOKEN = [System.Net.NetworkCredential]::new('', $secure).Password
tripol new MiBot --template discord
cd MiBot
notepad main.tri
tripol run
```

The generated template does not need a channel ID. `login_env()` reads `DISCORD_TOKEN`; optional `channel(...)` restricts replies. Never publish or share the token; a literal is extractable from both source and compiled binary. Reset a leaked token in the developer portal. Remove the environment variable after stopping: `Remove-Item Env:DISCORD_TOKEN`.

### Protocol and limitations

Initial Hello (10), Identify (2), Dispatch (0), sequence tracking, periodic Heartbeat (1) with null before the first sequence, Heartbeat ACK (11), server heartbeat requests, READY and MESSAGE_CREATE are implemented. `connected()` reports local READY state; `ping()` reports last heartbeat round trip in milliseconds or -1 before ACK, not message latency. Both are primarily diagnostics while running; the simple `run` API blocks. `preview` and `gateway_reply` support offline matching tests. `identify_payload` is a protocol helper containing the secret; do not print or log it.

This is an experimental first Gateway client for a small test server. Windows native live execution, TLS and actual presence/latency have not been verified in the Linux release environment. The Gateway uses a fixed official host and OS certificate validation without redirects. One receive worker at a time; no separate message queue or concurrent event callbacks. For this Windows backend, max complete UTF-8 payload 65535 bytes, bounded JSON nesting 32. Oversized or unsupported frames stop reception; complex/oversized READY payloads may prevent startup. No compression, binary frames, sharding, slash commands, interaction callbacks, Resume, automatic reconnect or production delivery guarantee. Reconnect (7), Invalid Session (9) and closed/failed receive stop safely; missed ACK, timeout waiting for READY, or send failure terminates the process to avoid an unsafe close while a receive is pending. Restart manually only after checking the cause; respect Discord identify and rate limits. No guaranteed delivery on disconnect. REST POST still creates its own WinHTTP session; delays depend on network/TLS/Discord, and failure (including HTTP 429) stops the process rather than retrying blindly. No measured Windows end-to-end latency claim. Ctrl+C exits. Do not deploy for critical workloads.

Explicit fallback: `discord.run_rest()` restores ECLIPSE's 2.5-second REST polling, which may miss messages. Legacy `listen`, `send`, `latest`, `matches`, `message_id`, `gateway` and token/URL helpers remain available. See `tripol docs discord` for the complete module, including the protocol helpers.

Official protocol: https://discord.com/developers/docs/events/gateway ; WinHTTP WebSocket API: https://learn.microsoft.com/windows/win32/api/winhttp/nf-winhttp-winhttpwebsocketcompleteupgrade ; rate limits: https://discord.com/developers/docs/topics/rate-limits .

### ORIGIN routing

`configured()` now means a token is configured, rather than a token and channel. `has_channel()` checks whether a channel filter exists. `all_channels()` clears the filter. `send_to(channel, message)` uses the configured bot token; `say()` and a nonempty `welcome()` require a selected channel. `run_rest()` still requires a channel. `event_channel(event)` validates the numeric channel ID and optional filter; `gateway_reply(event)` ignores bot/webhook authors as before. Guild message intents are unchanged; DMs are not enabled. Tests verify routing on Linux in all four profiles, but live Windows Gateway is not validated here.
