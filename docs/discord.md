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

The payload includes flag 32768 and disables mentions. The builder also has `button(label, custom_id, style=1)`, `section(content, thumbnail_url, description="")` and `gallery(urls)`. It validates the 40-component limit, duplicate button IDs, container nesting and field lengths. The current byte limits are 512 per text display and 16,384 per payload. Use `message.send(channel_id)` on Windows or Linux. Button interactions, select menus, modals and file uploads are pending.

El JSON incluye la bandera 32768 y desactiva las menciones. El constructor también ofrece `button(label, custom_id, style=1)`, `section(content, thumbnail_url, description="")` y `gallery(urls)`. Valida el límite de 40 componentes, los IDs de botones duplicados, el anidamiento y las longitudes de los campos. Los límites actuales son 512 bytes por texto y 16.384 bytes por mensaje codificado. Usa `message.send(channel_id)` en Windows o Linux. Quedan pendientes las interacciones de botones, los selectores, los formularios y los archivos adjuntos.

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

Text commands accept names of 1–96 bytes and replies of 1–512 bytes, both UTF-8. An embed requires text or media and accepts up to 6,000 text bytes, with a 256-byte title or field label and a 512-byte description, footer or field value. The encoded JSON limit is 16,384 bytes. The optional `intents` parameter on legacy `discord.run()` accepts `33281` (default for guild prefix commands) or `1` (GUILDS only).

Los comandos aceptan nombres de 1–96 bytes y respuestas de 1–512 bytes, ambos en UTF-8. El embed requiere texto o imágenes y acepta hasta 6.000 bytes de texto. El título y la etiqueta de un campo aceptan hasta 256 bytes, y la descripción, el pie y el valor de un campo hasta 512 bytes. El JSON codificado admite hasta 16.384 bytes. El parámetro opcional `intents` de `discord.run()` acepta `33281` (predeterminado para comandos de texto) o `1` (solo GUILDS).

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
