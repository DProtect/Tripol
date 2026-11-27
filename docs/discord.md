# DISCORD NEW — TRIPOL 0.28 development

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

## Connect from Windows / Conecta desde Windows

Replace `show bot.preview("hola")` with `bot.start()`, or start from `tripol new MiBot --template discord`. Set `DISCORD_TOKEN` in the environment of the Windows process using the instructions below. Enable Message Content Intent for guild text commands and test in a private server with View Channel and Send Messages. `embed.send("CHANNEL_ID")` uses the same environment variable and requires Embed Links in the target channel.

Cambia `show bot.preview("hola")` por `bot.start()`, o usa `tripol new MiBot --template discord`. Configura `DISCORD_TOKEN` en el entorno del proceso de Windows como se indica abajo. Activa Message Content Intent para comandos de texto y prueba en un servidor privado con permisos para ver el canal y enviar mensajes. `embed.send("CHANNEL_ID")` usa la misma variable y requiere el permiso para insertar enlaces.

This development build compiles the Windows Gateway and REST paths and tests offline behavior on Linux and WebAssembly. Live Discord validation requires an authorized private server. A process runs one Gateway session at a time. Automatic reconnect, slash commands and Components V2 remain pending.

Esta compilación de desarrollo genera las rutas Gateway y REST para Windows y prueba el comportamiento sin conexión en Linux y WebAssembly. La validación contra Discord requiere un servidor privado autorizado. Cada proceso ejecuta una sesión Gateway a la vez. Quedan pendientes la reconexión automática, los comandos slash y Components V2.

## API

| Method / Método | Behavior / Comportamiento |
| --- | --- |
| `discord.bot()` | Create a bot configuration / Crea una configuración de bot |
| `bot.prefix(value)` | Set a prefix of up to 32 UTF-8 bytes / Define un prefijo de hasta 32 bytes UTF-8 |
| `bot.status(value)` | Set `online`, `idle`, `dnd` or `invisible` / Define el estado indicado |
| `bot.activity(value)` | Set the activity text, up to 128 UTF-8 bytes / Define una actividad de hasta 128 bytes UTF-8 |
| `bot.command(name, response)` | Set one of up to 16 exact text commands / Configura uno de hasta 16 comandos de texto exactos |
| `bot.preview(content)` | Return the exact reply or an empty string / Devuelve la respuesta exacta o texto vacío |
| `bot.start()` | Open a blocking Windows Gateway session / Abre una sesión Gateway bloqueante en Windows |
| `discord.embed()` | Create an empty embed / Crea un embed vacío |
| `embed.title(value)` | Set the title / Define el título |
| `embed.description(value)` | Set the description / Define la descripción |
| `embed.footer(value)` | Set the footer text / Define el pie |
| `embed.color(value)` | Set an RGB color as `#RRGGBB` / Define un color RGB como `#RRGGBB` |
| `embed.field(label, value, inline=false)` | Append one of up to 25 fields / Añade uno de hasta 25 campos |
| `embed.payload()` | Build JSON including disabled mentions / Genera el JSON con las menciones desactivadas |
| `embed.send(channel_id)` | Send the embed over Discord REST from Windows / Envía el embed por Discord REST desde Windows |

Text commands accept names of 1–96 bytes and replies of 1–512 bytes, both UTF-8. An embed requires text and accepts up to 6,000 text bytes, with a 256-byte title or field label and a 512-byte description, footer or field value. The encoded JSON limit is 16,384 bytes. The optional `intents` parameter on legacy `discord.run()` accepts `33281` (default for guild prefix commands) or `1` (GUILDS only).

Los comandos aceptan nombres de 1–96 bytes y respuestas de 1–512 bytes, ambos en UTF-8. El embed requiere texto y acepta hasta 6.000 bytes de texto. El título y la etiqueta de un campo aceptan hasta 256 bytes, y la descripción, el pie y el valor de un campo hasta 512 bytes. El JSON codificado admite hasta 16.384 bytes. El parámetro opcional `intents` de `discord.run()` acepta `33281` (predeterminado para comandos de texto) o `1` (solo GUILDS).

## Discord R5 — ORIGIN Gateway preview (Windows x64)

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

`setup` simply combines `login` and `channel`. These older methods continue to work. The optional `welcome` sends a greeting once on READY; `say` sends immediately by REST. `status("online")` is optional (default online); `idle`, `dnd`, `invisible` also supported in the initial Identify payload. `activity` sets a playing activity, up to 128 UTF-8 bytes. Up to 16 exact commands across accessible guild channels (or one selected channel); duplicate registration replaces the response. Bots and webhooks are ignored. Prefix command matching uses exact JSON string spelling; ASCII commands recommended.

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

This is an experimental first Gateway client for a small test server. Windows native live execution, TLS and actual presence/latency have not been verified in the Linux release environment. The Gateway uses a fixed official host and OS certificate validation without redirects. One receive worker at a time; no separate message queue or concurrent event callbacks. Max complete UTF-8 payload 65535 bytes, bounded JSON nesting 32. Oversized or unsupported frames stop reception; complex/oversized READY payloads may prevent startup. No compression, binary frames, sharding, slash commands, interaction callbacks, Resume, automatic reconnect or production delivery guarantee. Reconnect (7), Invalid Session (9) and closed/failed receive stop safely; missed ACK, timeout waiting for READY, or send failure terminates the process to avoid an unsafe close while a receive is pending. Restart manually only after checking the cause; respect Discord identify and rate limits. No guaranteed delivery on disconnect. REST POST still creates its own WinHTTP session; delays depend on network/TLS/Discord, and failure (including HTTP 429) stops the process rather than retrying blindly. No measured Windows end-to-end latency claim. Ctrl+C exits. Do not deploy for critical workloads.

Explicit fallback: `discord.run_rest()` restores ECLIPSE's 2.5-second REST polling, which may miss messages. Legacy `listen`, `send`, `latest`, `matches`, `message_id`, `gateway` and token/URL helpers remain available. `discord.run()` on Linux cannot open the native Windows Gateway. See `tripol docs discord` for the complete module, including the protocol helpers.

Official protocol: https://discord.com/developers/docs/events/gateway ; WinHTTP WebSocket API: https://learn.microsoft.com/windows/win32/api/winhttp/nf-winhttp-winhttpwebsocketcompleteupgrade ; rate limits: https://discord.com/developers/docs/topics/rate-limits .

### ORIGIN routing

`configured()` now means a token is configured, rather than a token and channel. `has_channel()` checks whether a channel filter exists. `all_channels()` clears the filter. `send_to(channel, message)` uses the configured bot token; `say()` and a nonempty `welcome()` require a selected channel. `run_rest()` still requires a channel. `event_channel(event)` validates the numeric channel ID and optional filter; `gateway_reply(event)` ignores bot/webhook authors as before. Guild message intents are unchanged; DMs are not enabled. Tests verify routing on Linux in all four profiles, but live Windows Gateway is not validated here.
