# DISCORD REVOLUTION — TRIPOL 0.28 development

## SIMPLE R2 bots / Bots con SIMPLE R2

Use `bot.slash("info", "Informacion del bot", info)` to infer the signature of `func info(ctx)`. End the callback with `ctx.reply(value)` to return the payload. The accepted values are strings, a single embed, an embed batch and Components V2. Only a returned payload or an explicit `ctx.respond(payload)` sends the initial response. An earlier expression whose result is discarded builds a payload without sending it.

Usa `bot.slash("info", "Informacion del bot", info)` para inferir la firma de `func info(ctx)`. Termina el callback con `ctx.reply(valor)` para devolver el payload. Puedes pasar texto, un embed, un lote de embeds o un panel de Components V2. El envío inicial ocurre al devolver el payload o llamar explícitamente a `ctx.respond(payload)`. Una expresión anterior cuyo resultado se descarta solo construye el payload.

- [Bot con un embed / Embed bot](../examples/discord_simple_embeds.tri)
- [Bot con panel y botón inline / Components V2 bot](../examples/discord_simple_components.tri)

`discord.container("#FFFFFF")` wraps the layout in an accented container. Use button styles `"primary"`, `"secondary"`, `"success"` and `"danger"` or their numeric equivalents. `panel.premium(sku_id)` adds a purchase button for an application SKU. Discord handles the purchase. Use the typed button callback for ordinary buttons.

`discord.container("#FFFFFF")` envuelve el panel en un contenedor con color. Usa los estilos `"primary"`, `"secondary"`, `"success"` y `"danger"` o sus equivalentes numéricos. `panel.premium(sku_id)` añade un botón de compra de un SKU de la aplicación. Discord gestiona la compra. Los botones comunes admiten los callbacks tipados.

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

## Permissions and decimal choices / Permisos y elecciones decimales

Specify a permission name such as `"manage_messages"` when registering a command. Use `ctx.can("manage_messages")` in the callback to check the effective member permissions Discord supplied, including Administrator. Numeric bitsets still work for registration. This check does not inspect a target’s role hierarchy.

Puedes indicar un permiso como `"manage_messages"` al registrar un comando. Usa `ctx.can("manage_messages")` en el callback para comprobar los permisos efectivos del miembro informados por Discord, incluido Administrador. El registro también admite máscaras numéricas. Para actuar sobre un objetivo, la jerarquía de sus roles se comprueba por separado.

Build fixed decimal choices with `choices.number("Half", "0.5")`. The decimal text becomes a JSON number. `ctx.target_username()` and `ctx.target_content()` read the resolved selection in context menus. `server.move_role(role_id, position)` reorders a role above @everyone.

Añade elecciones decimales con `choices.number("Mitad", "0.5")`. El texto decimal se convierte en un número JSON. `ctx.target_username()` y `ctx.target_content()` leen la selección de un menú contextual. `server.move_role(role_id, position)` reordena un rol por encima de @everyone.

## Connect / Conecta

Set `DISCORD_TOKEN` before starting a bot in an authorized private server. For text commands and message callbacks, enable Message Content Intent. Use `tripol new MiBot --template discord` and `tripol run` on Windows or Linux x86-64.

Configura `DISCORD_TOKEN` antes de iniciar el bot en un servidor privado autorizado. Para comandos de texto y callbacks de mensajes, activa Message Content Intent. Usa `tripol new MiBot --template discord` y `tripol run` en Windows o Linux x86-64.

On Linux x86-64, `tripol build` packages the native program together with the TRIPOL launcher when it uses Discord networking. The launcher extracts the inner ELF into a private temporary directory and uses two inherited pipes for transport. It uses system trust roots for TLS. A raw Linux ELF cross-compiled from Windows can be launched on Linux with `tripol discord-run FILE`. A Linux bundle includes the compiler launcher, so it is larger than the inner ELF.

En Linux x86-64, `tripol build` empaqueta el programa nativo con el lanzador de TRIPOL si usa la red de Discord. El lanzador extrae el ELF a un directorio temporal privado y usa dos tuberías heredadas para el transporte. TLS usa los certificados de confianza del sistema. Un ELF generado desde Windows se ejecuta en Linux con `tripol discord-run FILE`. El paquete de Linux incluye el lanzador del compilador y ocupa más que el ELF interior.

On Windows x64, `tripol run` and `tripol discord-run FILE.exe` use the same TLS launcher as Linux. It negotiates zlib-stream, resumes sessions, checks heartbeat acknowledgements, reuses HTTPS connections and handles global and per-bucket rate limits. The directly launched PE uses the standalone WinHTTP fallback. See the [live validation guide](discord-validation.md).

En Windows x64, `tripol run` y `tripol discord-run FILE.exe` usan el mismo lanzador TLS que Linux. Negocia zlib-stream, reanuda sesiones, verifica las confirmaciones de heartbeat, reutiliza HTTPS y controla los límites globales y por bucket. El PE iniciado directamente usa WinHTTP. Consulta la [guía de validación conectada](discord-validation.md).

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

## Gateway events / Eventos del Gateway

```tripol
import tripol.discord

func joined(event:discord.GatewayEvent)->bool {
    show "Nuevo miembro", event.guild_id, event.user_id
    return true
}

registry = discord.events()
registry.intent("GUILD_MEMBERS")
registry.on("member_join", &joined)
show registry.preview("{\"op\":0,\"t\":\"GUILD_MEMBER_ADD\",\"s\":9,\"d\":{\"guild_id\":\"321\",\"user\":{\"id\":\"567\"}}}")
```

Call `bot.start_with_events(registry)` to connect with the same callbacks while keeping existing text and slash commands. Each event uses a `func(discord.GatewayEvent)->bool` handler. The context exposes `name`, `sequence`, `guild_id`, `channel_id`, `user_id` and `data_json`. Use `event.field(name)` to read any event property as JSON. Callbacks execute in arrival order and `registry.preview(event_json)` returns the handler's boolean result for testing.

Usa `bot.start_with_events(registry)` para conectar los mismos callbacks y conservar los comandos de texto y slash. Cada evento usa un manejador `func(discord.GatewayEvent)->bool`. El contexto expone `name`, `sequence`, `guild_id`, `channel_id`, `user_id` y `data_json`. Usa `event.field(nombre)` para leer cualquier propiedad como JSON. Los callbacks se ejecutan en orden de llegada y `registry.preview(event_json)` devuelve el booleano del manejador durante las pruebas.

Register official dispatch names such as `READY`, `GUILD_ROLE_UPDATE`, `MESSAGE_DELETE`, `MESSAGE_REACTION_ADD`, `VOICE_STATE_UPDATE` and `RESUMED`. The three member aliases are `member_join`, `member_leave` and `member_update`. The registry requests `GUILDS` by default. Use `registry.intent("GUILD_MESSAGE_REACTIONS")` for guild reactions and enable `GUILD_MEMBERS` or `GUILD_PRESENCES` in the registry and Developer Portal before registering privileged callbacks. The dedicated slash, button, select and modal methods handle interactions. Voice server credentials remain reserved for voice transport.

Registra nombres oficiales como `READY`, `GUILD_ROLE_UPDATE`, `MESSAGE_DELETE`, `MESSAGE_REACTION_ADD`, `VOICE_STATE_UPDATE` y `RESUMED`. Los alias de miembros son `member_join`, `member_leave` y `member_update`. El registro solicita `GUILDS` de forma predeterminada. Usa `registry.intent("GUILD_MESSAGE_REACTIONS")` para las reacciones del servidor y activa `GUILD_MEMBERS` o `GUILD_PRESENCES` en el registro y en el portal de desarrolladores antes de registrar callbacks privilegiados. Los métodos slash, de botones, selectores y formularios manejan las interacciones. Las credenciales del servidor de voz quedan reservadas para el transporte de voz.

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
show bot.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"type\":2,\"data\":{\"type\":1,\"name\":\"hello\",\"guild_id\":\"123\",\"options\":[{\"type\":3,\"name\":\"name\",\"value\":\"Ana\"}]}}}")
```

Use `bot.slash(name, description, &handler, guild_id="")` for a command without options. An empty guild ID registers globally. A guild ID limits registration to one server. The same name can have a different callback in each registered scope. At startup the bot reads its application ID, compares remote commands with the local definitions and creates or updates the changed ones. Registration leaves unrelated commands in place. This version supports ASCII lowercase command names and string, integer, boolean, user, channel, role, mentionable, number and attachment options.

Usa `bot.slash(nombre, descripción, &manejador, guild_id="")` para un comando sin opciones. Un ID de servidor vacío hace el registro global. Con un ID, el registro se limita a ese servidor. Un mismo nombre puede tener un callback diferente en cada ámbito. Al iniciar, el bot consulta su aplicación, compara los comandos y crea o actualiza los que cambiaron. El registro conserva los demás comandos. Se admiten nombres ASCII en minúsculas y opciones de texto, entero, booleano, usuario, canal, rol, mencionable, número y adjunto.

Register component and modal callbacks with `bot.on_button(id, &handler)`, `bot.on_select(id, &handler)` and `bot.on_modal(id, &handler)`. The handler can return an initial response such as `ctx.reply(text, true)`, `ctx.defer(true)`, `ctx.ack()`, `ctx.reply_components(layout)`, `ctx.update_components(layout)`, `ctx.reply_embeds(batch)` or `ctx.open_modal(form)`. Discord requires an initial response within three seconds.

Registra callbacks con `bot.on_button(id, &manejador)`, `bot.on_select(id, &manejador)` y `bot.on_modal(id, &manejador)`. El manejador puede devolver la respuesta inicial. `ctx.reply(texto, true)` responde en privado, `ctx.defer(true)` difiere la respuesta privada, y `ctx.ack()` confirma una interacción sobre un mensaje. Los métodos de componentes, embeds y formularios devuelven sus respectivos payloads. Discord exige una respuesta inicial en tres segundos.

```tripol
import tripol.discord

func report(ctx:discord.Interaction)->string {
    if ctx.respond(ctx.defer(true)) {
        if ctx.edit_reply("Informe listo") {
            ctx.followup("Consulta el resultado arriba", true)
        }
    }
    return ""
}

bot = discord.bot()
bot.slash("informe", "Prepara un informe", &report)
bot.start()
```

Connected native handlers can use `ctx.respond(payload)` to send the initial response before doing more work. It attempts delivery once and reports whether Discord acknowledged it. Use `ctx.edit_reply(text)` to fill a deferred response or edit its original message. `ctx.edit_response(body)` accepts a Discord edit payload, including `layout.payload()` for an original Components V2 message. `ctx.followup(text, true)` sends a private follow-up and returns its message ID, or an empty string if the request fails. For a response deferred with `ctx.defer()`, edit the original before posting a separate follow-up. The interaction token stays internal. Linux end-to-end tests use a local TLS server. Windows builds use WinHTTP and still need live validation.

Los callbacks nativos conectados pueden usar `ctx.respond(payload)` para enviar la respuesta inicial antes de continuar. Intenta entregarla una vez e indica si Discord la confirmó. Usa `ctx.edit_reply(texto)` para completar una respuesta diferida o editar su mensaje original. `ctx.edit_response(body)` acepta un payload de edición de Discord, incluido `layout.payload()` para un mensaje original Components V2. `ctx.followup(texto, true)` envía un seguimiento privado y devuelve el ID del mensaje, o una cadena vacía si falla la solicitud. Para una respuesta diferida con `ctx.defer()`, edita el original antes de enviar un seguimiento separado. El token de interacción queda en el contexto interno. Las pruebas de extremo a extremo de Linux usan un servidor TLS local. Las compilaciones de Windows usan WinHTTP y siguen pendientes de validación en vivo.

`ctx.option(name)` returns raw JSON, `ctx.option_text(name)` decodes a JSON string, `ctx.values()` returns selected values and `ctx.field(id)` reads a submitted text input. `ctx.field_data(id)` returns the JSON array for select and file-upload inputs. The original interaction data is available as `ctx.data_json`, including resolved attachment metadata.

`ctx.option(nombre)` devuelve JSON, `ctx.option_text(nombre)` decodifica una cadena, `ctx.values()` devuelve los valores seleccionados y `ctx.field(id)` lee un texto enviado. `ctx.field_data(id)` devuelve el arreglo JSON de los selectores y campos de archivos. `ctx.data_json` contiene los datos originales de la interacción y los metadatos resueltos de adjuntos.

After `id = ctx.followup(text)`, use `ctx.edit_followup(id, text)` or `ctx.delete_followup(id)`. Use `ctx.edit_followup_payload(id, layout.payload())` to update an existing Components V2 follow-up. These methods return false on a failed request and require a successful initial response in the same handler.

Después de `id = ctx.followup(texto)`, usa `ctx.edit_followup(id, texto)` o `ctx.delete_followup(id)`. Usa `ctx.edit_followup_payload(id, layout.payload())` para actualizar un seguimiento existente de Components V2. Estos métodos devuelven false si falla la solicitud y requieren una respuesta inicial confirmada dentro del mismo callback.

Use `tripol new Feedback --template discord-slash` for a runnable example. Set `DISCORD_TOKEN`, optionally set `DISCORD_GUILD_ID`, and install the bot application with the `applications.commands` scope on a private test server.

Usa `tripol new Feedback --template discord-slash` para generar un ejemplo ejecutable. Configura `DISCORD_TOKEN`, opcionalmente `DISCORD_GUILD_ID`, e instala la aplicación con el alcance `applications.commands` en un servidor privado.

### Context menus / Menús contextuales

```tripol
import tripol.discord

func inspect(ctx) {
    ctx.reply("ID: " + ctx.target_id(), true)
}

bot = discord.bot()
bot.user_menu("Ver perfil", inspect, "123")
bot.message_menu("Ver mensaje", inspect, "123")
bot.start()
```

User and message menus can share a name because Discord assigns them different command types. `ctx.target_id()` returns the selected ID, and `ctx.target()` returns the user or message object from Discord's resolved data. Use the last argument of `user_menu`, `message_menu`, `slash` or `slash_with_options` for default permissions. For example, `"8192"` requires Manage Messages, `"0"` disables the command by default except for administrators, and `"everyone"` clears that default requirement. An empty value omits the permission field.

Los menús de usuario y mensaje pueden compartir nombre porque Discord les asigna tipos distintos. `ctx.target_id()` devuelve el ID seleccionado y `ctx.target()` devuelve el objeto de usuario o mensaje que resolvió Discord. Usa el último argumento de `user_menu`, `message_menu`, `slash` o `slash_with_options` para configurar los permisos predeterminados. Por ejemplo, `"8192"` requiere Manage Messages, `"0"` desactiva el comando de forma predeterminada salvo para administradores y `"everyone"` elimina ese requisito predeterminado. Una cadena vacía omite el campo.

### Subcommands, choices and autocomplete / Subcomandos, elecciones y autocompletado

```tripol
import tripol.discord

func search(ctx:discord.Interaction)->string {
    return ctx.reply("Busca " + ctx.option_text("query"), true)
}

func suggest(ctx:discord.Interaction)->string {
    choices = discord.command_choices()
    if ctx.focused()=="query" and ctx.option_text("query")=="ca" {
        choices.add("Café", "cafe")
    }
    return ctx.autocomplete(choices)
}

arguments = discord.slash_options()
arguments.add("query", "Qué buscar", "string", true, true)
limits = discord.command_choices()
limits.integer("Diez", 10)
arguments.choose("limit", "Resultados máximos", limits)
commands = discord.slash_options()
commands.subcommand("find", "Busca una bebida", arguments)
root = discord.slash_options()
root.group("drinks", "Bebidas", commands)
bot = discord.bot()
bot.slash_with_options("order", "Haz tu pedido", &search, root)
bot.on_autocomplete("order", &suggest)
bot.start()
```

The example registers `/order drinks find`. `subcommand(name, description, options)` creates a branch and `group(name, description, subcommands)` collects one or more subcommands. Each command or subcommand accepts up to 25 options at its level. Keep required arguments first. The callback reads `ctx.command_group()` and `ctx.subcommand()` when it needs to select a branch. Pass `true` as the fifth argument of `options.add` to enable autocomplete. Use `ctx.focused()` to find the active option and `ctx.option(name)` to read a raw numeric value. Return `ctx.autocomplete(choices)` even when the list is empty. `options.choose` defines a fixed choice list. The builder supports up to 25 text or integer choices, including integer suggestions for number options. Fractional choice values are pending.

El ejemplo registra `/order drinks find`. `subcommand(nombre, descripción, opciones)` crea un subcomando y `group(nombre, descripción, subcomandos)` agrupa uno o más. Cada comando o subcomando admite hasta 25 opciones por nivel. Coloca primero los argumentos obligatorios. El callback lee `ctx.command_group()` y `ctx.subcommand()` para identificar la rama elegida. Pasa `true` como quinto argumento de `options.add` para activar el autocompletado. Usa `ctx.focused()` para identificar la opción activa y `ctx.option(nombre)` para leer un valor numérico en JSON. Devuelve `ctx.autocomplete(elecciones)` aunque la lista esté vacía. `options.choose` define una lista de elecciones fijas. El constructor admite hasta 25 elecciones de texto o enteras, con sugerencias enteras para opciones de tipo número. Quedan pendientes los valores fraccionarios.

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

The payload includes flag 32768 and disables mentions. The builder also has `button(label, custom_id, style=1)`, `section(content, thumbnail_url, description="")` and `gallery(urls)`. It validates the 40-component limit, duplicate button IDs, container nesting and field lengths. Text Display content can total 4,000 characters across the message. The encoded body must fit in the TRIPOL string capacity of 65,535 bytes. Use `message.send(channel_id)` on Windows or Linux. Use `discord.button_row()` to group up to five buttons. `discord.select(id, kind)` supports string, user, role, mentionable and channel menus. `discord.modal(id, title)` supports text inputs, selects, help text, file-upload inputs, radio groups, checkbox groups and individual checkboxes. The Linux launcher can send multipart attachments with `send_files`.

El JSON incluye la bandera 32768 y desactiva las menciones. El constructor también ofrece `button(label, custom_id, style=1)`, `section(content, thumbnail_url, description="")` y `gallery(urls)`. Valida el límite de 40 componentes, los IDs de botones duplicados, el anidamiento y las longitudes de los campos. El contenido de Text Display admite 4.000 caracteres por mensaje. El JSON debe caber en la capacidad de cadena de TRIPOL de 65.535 bytes. Usa `message.send(channel_id)` en Windows o Linux. Usa `discord.button_row()` para agrupar hasta cinco botones. `discord.select(id, tipo)` admite texto, usuario, rol, mencionable y canal. `discord.modal(id, título)` admite texto, selectores, ayuda, campos de archivos, grupos de opción única, grupos de casillas y casillas individuales. Los lanzadores de Windows y Linux envían adjuntos multipart con `send_files`.

For modal choices, create a String Select with the options and pass it to `form.radio(label, menu)` or `form.checkbox_group(label, menu)`. Use `form.checkbox(id, label)` for a single checkbox. Read a radio selection with `ctx.field(id)`, a group or selector with `ctx.field_values(id)` and a checkbox with `ctx.checked(id)`.

Para los campos de opciones, crea un selector de texto con las opciones y pásalo a `form.radio(etiqueta, selector)` o `form.checkbox_group(etiqueta, selector)`. Usa `form.checkbox(id, etiqueta)` para una casilla individual. Lee la selección única con `ctx.field(id)`, el grupo o selector con `ctx.field_values(id)` y la casilla con `ctx.checked(id)`.

#Use `ctx.has_field(id)` to distinguish a submitted empty field from an absent one. `ctx.field_attachment(id, index=0)` returns the selected file’s `id`, `filename`, `url`, `size` and `content_type`. An unresolved file has an empty ID. The method reads the submission without opening or downloading the file.

Usa `ctx.has_field(id)` para distinguir un campo enviado vacío de uno ausente. `ctx.field_attachment(id, index=0)` devuelve `id`, `filename`, `url`, `size` y `content_type`. Un archivo sin resolver tiene ID vacío. El método lee la respuesta del formulario sin abrir ni descargar el archivo.

## Button and select customization / Personalización de botones y selectores

`panel.button("Aceptar", "accept", "success", "✅")` combines text and emoji. Use an empty label for an emoji-only button. Add `true` after the emoji to disable it. For a row, the existing disabled argument comes first, as in `row.button("Aceptar", "accept", "success", false, "✅")`. `row.disable("accept")` disables that button before the row is added to a new layout, and `row.disable("accept", false)` enables it. To update a sent Components V2 message, build the new layout and use `ctx.update_components(panel)` or `room.edit_components(message_id, panel)`.

`panel.button("Aceptar", "accept", "success", "✅")` combina texto y emoji. Usa una etiqueta vacía para mostrar solo el emoji. Añade `true` después del emoji para deshabilitarlo. En una fila se conserva primero el argumento de deshabilitación, como en `row.button("Aceptar", "accept", "success", false, "✅")`. `row.disable("accept")` deshabilita ese botón antes de añadir la fila a un panel nuevo. Usa `row.disable("accept", false)` para habilitarlo. Para actualizar un mensaje Components V2 ya enviado, crea el nuevo panel y usa `ctx.update_components(panel)` o `room.edit_components(message_id, panel)`.

Use Unicode, a custom emoji ID, `"<:check:123>"` or `"<a:party:123>"` for an animated custom emoji. Pass it as the last argument of `menu.option(label, value, description, selected, emoji)`. `menu.disable()` applies to a select in a message. Discord handles link and premium buttons directly. See the executable example for all six styles. [Buttons example](../examples/discord_buttons.tri)

Usa Unicode, un ID de emoji personalizado, `"<:check:123>"` o `"<a:party:123>"` para un emoji personalizado animado. Pásalo como último argumento de `menu.option(etiqueta, valor, descripción, seleccionado, emoji)`. `menu.disable()` se aplica a los selectores de mensajes. Discord gestiona los botones de enlace y de compra. El ejemplo ejecutable reúne los seis estilos. [Ejemplo de botones](../examples/discord_buttons.tri)

### Attachments / Adjuntos

```tripol
import tripol.discord

panel = discord.container()
panel.file("report.txt")
message_id = panel.send_files("123", list(["./report.txt"]))
show message_id
```

With the Linux or Windows launcher, `send_files` streams 1–10 local files into a multipart channel message and returns its ID or an empty string on failure. `ctx.reply_files(content, paths, true)` sends a private initial response immediately and returns whether it was accepted. It accepts text, an embed, a batch or Components V2. Use distinct ASCII filenames and use the exact basename in `discord.attachment_url("photo.png")` or `panel.file("report.txt")`. Follow-up uploads and replacement of existing attachments remain pending.

Con el lanzador de Linux o Windows, `send_files` adjunta entre 1 y 10 archivos locales a un canal y devuelve el ID o una cadena vacía si falla. `ctx.reply_files(contenido, rutas, true)` envía de inmediato la respuesta inicial privada e indica si se aceptó. Admite texto, un embed, un lote o Components V2. Usa nombres ASCII distintos y el nombre exacto del archivo en `discord.attachment_url("photo.png")` o `panel.file("report.txt")`. Quedan pendientes los adjuntos en seguimientos y la sustitución de archivos.

For uploads submitted by users, call `form.file("files", "Upload", list(["image", ".pdf"]), 0, 2, false)`. Discord applies up to ten extension filters. Check the actual contents before trusting an upload.

Para archivos enviados por usuarios, llama a `form.file("files", "Adjuntar", list(["image", ".pdf"]), 0, 2, false)`. Discord aplica hasta diez filtros por extensión. Comprueba el contenido del archivo antes de confiar en él.

## Discord Studio

```sh
tripol discord studio main.tri --metrics studio.jsonl
```

On Linux x86-64, Studio runs the bot and samples its RSS, CPU ticks, connection state, heartbeat latency, uptime, guild count at READY, queued dispatches, REST requests, errors and rate limits. The launcher and bot process are reported separately. The output file must be new. Without `--metrics`, snapshots are sent to standard error. The export contains counters and process measurements. Live command and interaction metrics, log filtering and Windows monitoring are pending.

En Linux x86-64, Studio ejecuta el bot y toma muestras de RSS, ticks de CPU, conexión, latencia de heartbeat, tiempo activo, servidores en READY, cola de eventos, peticiones REST, errores y límites. El lanzador y el bot se miden por separado. Usa un archivo de salida nuevo. Sin `--metrics`, las muestras se escriben en la salida de error. Quedan pendientes las métricas por comando e interacción, el filtrado de logs y el monitoreo en Windows.

Studio adds `events_by_type`, `interactions_by_type` and `commands_by_name` to its Linux samples. They count events admitted to the Gateway queue, not callback results. Command keys include their type, such as `user/Ver perfil`. The counters contain neither interaction tokens nor component IDs. Names beyond the tracked set are counted under `other`.

Studio añade `events_by_type`, `interactions_by_type` y `commands_by_name` a sus muestras de Linux. Cuentan los eventos admitidos en la cola del Gateway. Los resultados de los callbacks se registran por separado si el programa los mide. Las claves incluyen el tipo del comando, como `user/Ver perfil`. Los contadores mantienen fuera los tokens de interacción y los IDs de componentes. Los nombres que exceden el conjunto registrado se agrupan en `other`.

## Member moderation and welcome bot / Moderación de miembros y bienvenida

```tripol
import tripol.discord

server = discord.guild("123")
member = server.member("456")
show member.read()
```

Use the numeric IDs from an authorized server. `member.kick()`, `member.ban(delete_message_seconds=0)` and `member.unban()` remove a member, ban them with optional removal of up to seven days of messages, or revoke a ban. `member.nickname(text)` sets a nickname, and an empty string clears it. `member.timeout(until)` sets an ISO 8601 deadline, and `member.timeout()` clears it. Discord enforces the 28-day maximum. `member.role(id, false)` removes a role and `member.role(id)` adds it. These methods report success as a boolean so the bot can continue after a failed request. Discord enforces bot permissions and role hierarchy. `member.dm(text)` and `discord.post_message(channel_id, text)` return the sent message ID or an empty string if delivery fails.

Usa los IDs numéricos de un servidor autorizado. `member.kick()`, `member.ban(delete_message_seconds=0)` y `member.unban()` expulsan, banean con borrado opcional de hasta siete días de mensajes o revocan el baneo. `member.nickname(texto)` cambia el apodo y una cadena vacía lo elimina. `member.timeout(hasta)` recibe una fecha ISO 8601 y `member.timeout()` retira el aislamiento. Discord aplica el máximo de 28 días. `member.role(id, false)` quita un rol y `member.role(id)` lo añade. Estos métodos devuelven un booleano para que el bot pueda continuar si falla la solicitud. Discord aplica los permisos y la jerarquía de roles del bot. `member.dm(texto)` y `discord.post_message(channel_id, texto)` devuelven el ID del mensaje o una cadena vacía si falla el envío.

Use `tripol new Welcome --template discord-welcome` to send a welcome to a channel when a human joins the configured server. The generated README covers the required environment variables, Server Members Intent and channel permissions. The template skips bots and other servers, and continues receiving events if a send fails.

Usa `tripol new Bienvenida --template discord-welcome` para enviar un saludo a un canal cuando se une una persona al servidor configurado. El README generado indica las variables de entorno, Server Members Intent y los permisos del canal. La plantilla omite bots y otros servidores y continúa recibiendo eventos si falla un envío.

Create a restricted bot with `tripol new Moderator --template discord-moderation`. Set `DISCORD_TOKEN`, `DISCORD_GUILD_ID` and `DISCORD_OPERATOR_ID`. The `/kick`, `/ban`, `/unban`, `/timeout` and `/warn` commands check the operator, guild and effective permissions returned by Discord. Verified callers receive a private acknowledgment. `/warn` appends a JSON line to `discord-warnings.jsonl` or the path in `DISCORD_WARNINGS_PATH`. Keep it in a directory writable only by the bot account.

Crea un bot restringido con `tripol new Moderador --template discord-moderation`. Configura `DISCORD_TOKEN`, `DISCORD_GUILD_ID` y `DISCORD_OPERATOR_ID`. Los comandos `/kick`, `/ban`, `/unban`, `/timeout` y `/warn` comprueban operador, servidor y permisos efectivos devueltos por Discord. Se confirma en privado a los usuarios autorizados. `/warn` guarda una línea JSON en `discord-warnings.jsonl` o en la ruta de `DISCORD_WARNINGS_PATH`. Usa un directorio donde solo pueda escribir la cuenta del bot.

## Channel moderation / Moderación de canales

```tripol
import tripol.discord

room = discord.channel_ref("123")
show room.read()
```

`room.slowmode(15)` sets a 15-second delay and `room.slowmode(0)` clears it. Use `room.pin(message_id)` to pin and `room.pin(message_id, false)` to unpin. `room.delete_message(id)` removes one message. `room.purge(ids)` requires a list of 2–100 distinct message IDs. Discord enforces bulk deletion eligibility, including the 14-day age limit. `room.react(id, "🔥")` adds the bot's reaction. Use `false` as the last argument to remove it. Unicode emoji and custom `name:id` values are percent-encoded automatically.

`room.slowmode(15)` fija una espera de 15 segundos y `room.slowmode(0)` la retira. Usa `room.pin(message_id)` para fijar un mensaje y `room.pin(message_id, false)` para desfijarlo. `room.delete_message(id)` borra un mensaje. `room.purge(ids)` requiere entre 2 y 100 IDs distintos. Discord aplica los requisitos del borrado masivo, incluido el límite de 14 días de antigüedad. `room.react(id, "🔥")` añade la reacción del bot. Usa `false` como último argumento para retirarla. Los emoji Unicode y los valores personalizados `nombre:id` se codifican automáticamente.

Configure the token of an authorized bot. Discord checks the required permissions. These methods return a boolean on mutations and keep the native program running if the request fails.

Configura el token de un bot autorizado. Discord comprueba los permisos requeridos. Los métodos que modifican el canal o un mensaje devuelven un booleano y permiten que el programa siga ejecutándose si falla la solicitud.

### Channels, categories and audit / Canales, categorías y auditoría

Use `server.create_channel(name, kind, parent_id, topic)` with `"text"`, `"voice"`, `"category"` or `"announcement"`. The result is the new ID or an empty string on failure. `server.channels()` returns the REST response as JSON. Read the audit log with `server.audit(limit, before, action)`. On a channel, use `rename`, `topic`, `move` and `delete`. An empty parent in `move()` detaches the channel from its category. `edit_message` and `edit_components` update messages owned by the bot.

Usa `server.create_channel(nombre, tipo, parent_id, tema)` con `"text"`, `"voice"`, `"category"` o `"announcement"`. Devuelve el ID nuevo o una cadena vacía si falla. `server.channels()` devuelve la respuesta REST como JSON. Consulta auditoría con `server.audit(límite, antes, acción)`. En un canal están disponibles `rename`, `topic`, `move` y `delete`. Un padre vacío en `move()` saca el canal de su categoría. `edit_message` y `edit_components` actualizan mensajes del propio bot.

### Guild roles / Roles de servidor

Use `server.create_role("Reviewer", "0", "#FFFFFF")` to create a role with no permissions of its own. The method returns its ID. `server.roles()` lists the roles. Use `role = server.role_ref(id)` with `role.read()`, `role.edit(name, permissions, color)` and `role.delete()`. Empty edit arguments leave those properties as they are. `role.permissions()` reads the combined guild permissions of the role and @everyone as a decimal string. Effective permissions for a member in a channel also depend on other roles and overrides.

Usa `server.create_role("Revisor", "0", "#FFFFFF")` para crear un rol sin permisos propios. El método devuelve su ID. `server.roles()` lista los roles. Usa `role = server.role_ref(id)` con `role.read()`, `role.edit(nombre, permisos, color)` y `role.delete()`. Los argumentos vacíos de edición conservan esas propiedades. `role.permissions()` devuelve como cadena decimal los permisos combinados del rol y @everyone en el servidor. Los permisos efectivos de un miembro en un canal también dependen de los demás roles y las excepciones.

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

The standalone Windows PE uses WinHTTP and the system certificate store. It accepts uncompressed text frames up to 65,535 bytes and dispatches callbacks in order. Use the TRIPOL launcher on Windows for Resume, transport compression, HTTPS reuse, rate-limit handling and multipart uploads.

Explicit fallback: `discord.run_rest()` restores ECLIPSE's 2.5-second REST polling, which may miss messages. Legacy `listen`, `send`, `latest`, `matches`, `message_id`, `gateway` and token/URL helpers remain available. See `tripol docs discord` for the complete module, including the protocol helpers.

Official protocol: https://discord.com/developers/docs/events/gateway ; WinHTTP WebSocket API: https://learn.microsoft.com/windows/win32/api/winhttp/nf-winhttp-winhttpwebsocketcompleteupgrade ; rate limits: https://discord.com/developers/docs/topics/rate-limits .

### ORIGIN routing

`configured()` now means a token is configured, rather than a token and channel. `has_channel()` checks whether a channel filter exists. `all_channels()` clears the filter. `send_to(channel, message)` uses the configured bot token; `say()` and a nonempty `welcome()` require a selected channel. `run_rest()` still requires a channel. `event_channel(event)` validates the numeric channel ID and optional filter; `gateway_reply(event)` ignores bot/webhook authors as before. Guild message intents are unchanged; DMs are not enabled. Tests verify routing on Linux in all four profiles, but live Windows Gateway is not validated here.

## Cooldowns and command permissions / Cooldowns y permisos de comandos

Register the slash command before configuring it. `bot.cooldown("work", 60)` reserves one execution per user every 60 seconds. Choose `"channel"`, `"guild"` or `"global"` for a shared cooldown within the command’s registration scope. An optional message is followed by the remaining seconds. Use zero seconds to disable the rule. A command that needs a channel, server or user ID is denied when that context is absent.

Registra el comando slash antes de configurarlo. `bot.cooldown("work", 60)` reserva una ejecución por usuario cada 60 segundos. Elige `"channel"`, `"guild"` o `"global"` para compartir el cooldown dentro del ámbito de registro del comando. Un mensaje opcional precede a los segundos restantes. Usa cero segundos para desactivar la regla. Si falta el ID requerido por el ámbito, el comando se rechaza.

```tripol
import tripol.discord

bot=discord.bot()
bot.slash("work","Earn points",func(ctx) {ctx.reply("Done",true)})
show bot.cooldown("work",60,"user","Try again in")
```

`ctx.cooldown_remaining()` returns rounded-up seconds, or `-1` when unavailable. Use `ctx.reset_cooldown()` for the current bucket. For administrative tools, use `discord.cooldown_bucket("work","user",user_id,guild_id)` with `.remaining()` and `.reset()`. A reservation is made before the handler and remains after an error unless the handler resets it. `bot.persist_cooldowns("cooldowns.json")` keeps active cooldowns across restarts. Use a separate file per app in a directory writable by its account. In persistent mode each successful reservation or reset saves a snapshot. A file lock prevents two launchers from writing the same file. In-memory entries expire automatically.

`ctx.cooldown_remaining()` devuelve los segundos redondeados hacia arriba o `-1` cuando la consulta no está disponible. Usa `ctx.reset_cooldown()` para restablecer el registro actual. Para herramientas administrativas, usa `discord.cooldown_bucket("work","user",user_id,guild_id)` con `.remaining()` y `.reset()`. La reserva se hace antes del callback y se mantiene si ocurre un error. El callback puede restablecerla. `bot.persist_cooldowns("cooldowns.json")` conserva los cooldowns activos al reiniciar. Usa un archivo distinto por aplicación en un directorio donde pueda escribir su cuenta. Cuando la persistencia está activa, cada reserva o restablecimiento guarda una instantánea. Un candado impide la escritura por dos lanzadores. Los registros en memoria caducan automáticamente.

`bot.permission("ban","BAN_MEMBERS")` sets the registration default and checks the effective permissions Discord sends for the invoking member and bot in the current channel. For kick, ban and timeout, a single top-level user option is inferred as the member to protect. Specify the option as the third argument when the target cannot be inferred. Pass the registration guild ID as the fourth argument to `permission` and the fifth to `cooldown`. Run protected moderation and cooldowns through `tripol run` or `tripol discord-run`.

`bot.permission("ban","BAN_MEMBERS")` configura el permiso predeterminado y comprueba los permisos efectivos del miembro y del bot que Discord comunica para el canal. Para expulsar, banear y aislar se infiere el objetivo cuando hay una sola opción de usuario en el nivel superior. Indica la opción como tercer argumento si el objetivo no se puede deducir. El ID de servidor de registro es el cuarto argumento de `permission` y el quinto de `cooldown`. Ejecuta la moderación protegida y los cooldowns con `tripol run` o `tripol discord-run`.

The hierarchy check protects the owner, the operator and the bot itself. It verifies the highest roles for the member, target and bot, and denies timeout of an administrator. `ctx.can()`, `ctx.bot_can()` and `ctx.can_moderate()` are available for explicit checks. Lower-level REST methods use the bot’s authorization on Discord. Use a protected callback when accepting administrative requests from users.

La comprobación de jerarquía protege al propietario, al operador y al propio bot. Verifica los roles superiores del miembro, el objetivo y el bot, y rechaza el aislamiento de un administrador. `ctx.can()`, `ctx.bot_can()` y `ctx.can_moderate()` permiten comprobaciones explícitas. Los métodos REST de bajo nivel usan la autorización del bot en Discord. Usa un callback protegido para aceptar solicitudes administrativas de usuarios.

Use `bot.restrict("ban","role",role_id)` to restrict a command to a role. Use `"user"`, `"channel"` or `"guild"` for those scopes. IDs within a scope are alternatives. A command that configures several scopes must satisfy each of them. `ctx.has_role(role_id)` supports an explicit role check. The launcher evaluates restrictions before reserving a cooldown or starting a hierarchy check.

Usa `bot.restrict("ban","role",role_id)` para restringir un comando a un rol. Usa `"user"`, `"channel"` o `"guild"` para esos ámbitos. Los IDs de un mismo ámbito son alternativas. Un comando con varios ámbitos configurados debe satisfacer cada uno. `ctx.has_role(role_id)` permite comprobar un rol expresamente. El lanzador evalúa las restricciones antes de reservar el cooldown o verificar la jerarquía.

## Channel and resource management / Gestión de canales y recursos

`channel.lock()` saves the @everyone overwrite and denies Send Messages for that role. Explicit user and role overrides continue to apply. `channel.unlock()` restores exactly the saved overwrite in the same launcher, or removes the temporary overwrite if the channel did not have one. If the overwrite has changed externally, the operation returns false and preserves the external edit. Use `overwrite(id, allow, deny, "user")` or `overwrite(id, allow, deny, "role")` for explicit changes, and `clear_overwrite(id)` to remove one. `guild.private_channel(name, user_id)` creates a text channel that grants that member access and denies viewing to @everyone. `guild.move_channel(id, position)` updates the channel order.

`channel.lock()` guarda el permiso de @everyone y le deniega Enviar mensajes. Los permisos explícitos de otros roles y usuarios siguen aplicándose. `channel.unlock()` restaura la configuración guardada dentro del mismo lanzador o quita el permiso temporal si antes no existía. Si el permiso cambió externamente, devuelve false y conserva el cambio. Usa `overwrite(id, allow, deny, "user")` o `overwrite(id, allow, deny, "role")` para cambios explícitos, y `clear_overwrite(id)` para eliminar uno. `guild.private_channel(name, user_id)` crea un canal de texto que permite el acceso a ese miembro y deniega la vista a @everyone. `guild.move_channel(id, position)` actualiza el orden.

`channel.create_webhook(name)` returns a typed webhook. Use `.read()`, `.edit(name)` and `.delete()` to administer it. The returned read information omits credentials. `.send()` accepts text, embeds or Components V2. `.send_files()` streams local attachments through the launcher. Discord controls which components each webhook can use. `channel.invite(max_age, max_uses)` returns the code; list invitations with `.invites()`, inspect one with `discord.invite_info(code)` and remove it with `discord.revoke_invite(code)`.

`channel.create_webhook(name)` devuelve un webhook tipado. Usa `.read()`, `.edit(name)` y `.delete()` para administrarlo. La información devuelta al consultar omite las credenciales. `.send()` recibe texto, embeds o Components V2. `.send_files()` transmite adjuntos locales mediante el lanzador. Discord controla qué componentes puede usar cada webhook. `channel.invite(max_age, max_uses)` devuelve el código. Enumera invitaciones con `.invites()`, consulta una con `discord.invite_info(code)` y elimínala con `discord.revoke_invite(code)`.

Create a thread with `channel.thread(name)` or `channel.thread_from(message_id, name)`. Use `thread.thread_member(id)`, `thread.archive()` and `thread.lock_thread()`. For forums, use `guild.create_channel(name, "forum")`, `forum.post(title, text, tags)` and `post.apply_tags(tags)`. Use `discord.forum_tags()` to build the full set supplied to `forum.set_forum_tags(tags)`, including the IDs of existing labels that should stay. The [live validation guide](discord-validation.md) lists the required permissions and expected results.

Crea un hilo con `channel.thread(name)` o `channel.thread_from(message_id, name)`. Usa `thread.thread_member(id)`, `thread.archive()` y `thread.lock_thread()`. Para foros, usa `guild.create_channel(name, "forum")`, `forum.post(title, text, tags)` y `post.apply_tags(tags)`. Usa `discord.forum_tags()` para construir el conjunto completo que recibe `forum.set_forum_tags(tags)), incluidos los IDs de las etiquetas que quieras conservar. La [guía de validación conectada](discord-validation.md) recoge los permisos y resultados esperados.

## Deferred external HTTP / HTTP externo diferido

Use `ctx.defer_http(request, "profile")` with an HTTP request and register `events.on_http("profile", handler)` before starting the bot. The result arrives in `event.http_response()`. Use `event.complete_http(text)` to edit the original deferred reply once. [Guide and sample](http.md#https-errores-y-discord--https-errors-and-discord)

Llama a `ctx.defer_http(request, "profile")` y registra `events.on_http("profile", handler)` antes de iniciar el bot. Lee el resultado con `event.http_response()` y edita la respuesta diferida una vez con `event.complete_http(text)`. [Guía y ejemplo](http.md#https-errores-y-discord--https-errors-and-discord)

## HTTP diferido con paneles / Deferred HTTP panels

`ctx.defer_http_task(request,"weather")` devuelve un recibo. `discord.cancel_http(receipt)` solicita la cancelación y deja que el callback reciba `response.error == "canceled"`. En un botón o formulario con mensaje asociado, `ctx.defer_http_update(request,"weather")` reconoce la interacción y permite editar ese mismo mensaje.

`ctx.defer_http_task(request,"weather")` returns a receipt. Call `discord.cancel_http(receipt)` to cancel the upstream request and collect a canceled result. Use `ctx.defer_http_update` from a component or modal interaction that has an original message to refresh.

Registra `events.on_http("weather", handler)` para tratar el resultado. `event.queue_http(panel)` pone la edición en cola y permite seguir recibiendo eventos mientras Discord responde. Admite una cadena, un embed o Components V2. `events.on_http_delivery("weather", delivered)` recibe el estado de la entrega y el mismo recibo. Usa `event.http_delivered()` para comprobarlo. `event.complete_http` permite esperar el resultado de la edición de forma síncrona.

Use `events.on_http` for the upstream result and `event.queue_http` to enqueue the original reply edit. The optional `events.on_http_delivery` handler receives its delivery status and the same receipt. `event.complete_http` waits for the edit. The launcher expires continuations and cancels their outstanding work when a bot explicitly closes its Gateway. [Executable examples and authorized validation](http.md)
