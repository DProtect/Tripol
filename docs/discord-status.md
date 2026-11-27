# TRIPOL v0.28 — DISCORD REVOLUTION

Estado de desarrollo frente al prompt maestro. «Completo» indica una capacidad delimitada con implementación y pruebas. «Experimental» indica una ruta implementada con validación local o con alcance reducido. Las pruebas TLS y Gateway usan un servidor local. Las pruebas reales de administración, mensajes y voz requieren un servidor privado autorizado.

| Área | Funcionalidad | Estado | Evidencia |
| --- | --- | --- | --- |
| Lenguaje | Métodos de estructuras, inferencia básica, valores predeterminados y callbacks tipados con `&nombre` | Completo | [Pruebas de métodos](../core/methods_test.go), [pruebas de callbacks](../core/callbacks_test.go) |
| Lenguaje | Contexto básico de mensaje y errores recuperables de validación | Completo | [Pruebas de contexto y validación](../core/discord_components_test.go) |
| Lenguaje | Estructuras anidadas generales, colecciones de objetos y concurrencia general de callbacks | Pendiente | [Referencia del lenguaje](language.md) |
| API simple | `bot()`, `start()`, prefijo, estado, actividad, comandos exactos y vista previa | Completo en su alcance | [API de Discord](discord.md), [pruebas de estado](../core/methods_test.go) |
| API simple | Stop, restart, uptime, selección general de intents y gestión pública del ciclo de vida | Pendiente | [API de Discord](discord.md) |
| API simple | Indicadores `connected()` y `ping()` de bajo nivel en Windows | Experimental | [Configuración y diagnóstico Gateway](../core/gateway_test.go) |
| Conexión | Gateway TLS, Identify, Hello, heartbeat, ACK y mensajes fragmentados sin compresión | Experimental en Windows y Linux | [WinHTTP](../native/liberty.S), [pruebas WebSocket TLS](../internal/discordtransport/transport_test.go) |
| Conexión | ELF Linux con entorno, lanzador integrado y backend HTTPS nativo | Experimental | [Prueba ELF y tuberías](../core/discord_bridge_linux_test.go), [prueba de paquete](../tests/cli_test.go) |
| Conexión | Descubrimiento autenticado, cuotas Identify, secuencia, Resume y reconexión con jitter | Experimental en Linux | [Pruebas de descubrimiento y Resume](../internal/discordtransport/transport_test.go) |
| Conexión | Replay tras Invalid Session y cierre en errores de autenticación o configuración | Experimental en Linux | [Supervisor Gateway](../internal/discordtransport/gateway.go), [prueba de cierre 4014](../internal/discordtransport/transport_test.go) |
| Conexión | Reconexión y Resume en Windows, cancelación desde la API y compresión oficial | Pendiente | [Ruta WinHTTP](discord.md) |
| REST | Reutilización HTTPS, 429, Retry-After, límites globales y buckets, reintentos acotados | Experimental en Linux | [Pruebas REST](../internal/discordtransport/transport_test.go) |
| REST | Colas de prioridad, cancelación pública, coordinación con sharding y rate limits en Windows | Pendiente | [Transporte REST](../internal/discordtransport/rest.go) |
| Comandos | Texto por prefijo, registro de respuestas fijas, comparación UTF-8 y callbacks reales de mensajes | Experimental con transporte y pruebas locales | [Prueba ELF con callback y solicitud REST](../core/discord_bridge_linux_test.go), [prueba de escapes Unicode](../core/gateway_test.go) |
| Comandos | Argumentos, alias, subcomandos, ayuda, cooldowns, permisos y restricciones | Pendiente | [Guía de comandos y callbacks](discord.md) |
| Slash | Registro y sincronización, comandos globales y por servidor, opciones, subcomandos, context menus y autocompletado | Pendiente | [Guía de Discord](discord.md) |
| Interacciones | Respuestas efímeras, diferidas, editadas y de seguimiento, callbacks de botones y selectores | Pendiente | [Guía de Discord](discord.md) |
| Eventos | READY, MESSAGE_CREATE, RESUMED en Linux, reenvío del JSON original de mensajes | Experimental | [Pruebas Gateway](../internal/discordtransport/transport_test.go), [contexto de mensaje](../core/discord_components_test.go) |
| Eventos | Registro general de eventos, miembros, roles, canales, reacciones, voz y auditoría | Pendiente | [Eventos oficiales](https://docs.discord.com/developers/events/gateway-events) |
| Embeds | Título, descripción, color, pie, campos, enlace, imagen, miniatura y JSON | Completo para serialización | [Pruebas de embeds](../core/methods_test.go), [pruebas de medios](../core/discord_components_test.go) |
| Embeds | Envío nativo | Experimental | [Prueba de importaciones PE](../core/methods_test.go), [API de Discord](discord.md) |
| Embeds | Autor, timestamp, múltiples embeds y adjuntos por mensaje | Pendiente | [Guía de Discord](discord.md) |
| Components V2 | Text Display, separador, sección, miniatura, galería, contenedor, fila con botón y flag 32768 | Completo para los constructores y límites documentados | [JSON y validación de Components V2](../core/discord_components_test.go) |
| Components V2 | Envío nativo de Components V2 | Experimental | [Constructor y envío REST](../stdlib/discord.tri), [guía bilingüe](discord.md) |
| Components V2 | Filas arbitrarias, archivos, todos los selectores oficiales, formularios y sus componentes, edición dinámica | Pendiente | [Referencia oficial](https://docs.discord.com/developers/components/reference) |
| Servidores | Moderación, usuarios, miembros, roles, permisos, canales, categorías, webhooks, invitaciones, hilos y foros | Pendiente | [Guía de Discord](discord.md) |
| Servidores | Administración de mensajes, archivos, reacciones y registros de auditoría | Pendiente | [Guía de Discord](discord.md) |
| Rendimiento | Comparación JSON de texto sin construir una copia decodificada, límites de memoria y reutilización HTTPS | Experimental | [Comparación de texto y memoria](../core/gateway_test.go), [pruebas REST](../internal/discordtransport/transport_test.go) |
| Rendimiento | Benchmarks equivalentes frente a Discord.js, Discord.py y Serenity, latencia real, RSS y Ryzen 5 3400G | Pendiente | [Benchmarks generales del compilador](../benchmarks/) |
| Rendimiento | Sharding, despachador concurrente, cachés de entidades y reducción del tamaño del paquete Linux | Pendiente | [Arquitectura del lanzador](../core/discord_bridge_linux.go) |
| Studio | Estado del Gateway, latencia heartbeat, CPU ticks, RSS, tiempo activo, servidores en READY, contador de eventos, cola, REST, errores, Resume y 429 | Experimental en Linux | [Métricas](../internal/discordtransport/metrics.go), [pruebas](../internal/discordtransport/transport_test.go) |
| Studio | Exportación JSON Lines de muestras y separación de recursos de bot y lanzador | Completo para esa exportación en Linux | [Prueba de CLI](../tests/cli_test.go) |
| Studio | Monitoreo en Windows, contadores por comando/interacción y filtrado INFO/DEBUG/WARNING/ERROR | Pendiente | [Guía](discord.md) |
| Plantillas | `discord`, `discord-embeds`, `discord-components` con preview y rutas de envío | Experimental para operación conectada | [Generador](../core/nexus.go), [compilación en Windows y Linux](../core/nexus_test.go) |
| Plantillas | Moderación, tickets, niveles, economía, bienvenida, utilidades y almacenamiento persistente asociado | Pendiente | [Generador](../core/nexus.go) |
| Seguridad | Tokens desde entorno, TLS con certificados, validación de destino, límites de cola y payloads, menciones desactivadas y filtro de bots/webhooks | Completo en las rutas verificadas | [Pruebas TLS y Gateway](../internal/discordtransport/transport_test.go), [pruebas de filtrado](../core/discord_components_test.go) |
| Seguridad | Auditoría de registros y exportaciones de todo el ecosistema, permisos avanzados de servidores | Pendiente | [Guía de Discord](discord.md) |
| Voz | Conexión, sesiones, Opus, DAVE, control de reproducción, colas y recuperación | Pendiente | [Protocolo oficial de voz](https://docs.discord.com/developers/topics/voice-connections) |
| Compatibilidad | API anterior, PE32+, ELF64, perfiles nativos y previews WebAssembly | Completo para los casos de regresión | [CI del compilador](../.github/workflows/ci.yml), [pruebas web](../cmd/tripol-web/test-r3.cjs) |
| Ecosistema | Autocompletado y ejemplos del Playground, VS Code, guía bilingüe y referencia de IA | Completo para la API implementada | [Playground](https://tripollas.github.io/Tripol/playground/?library=discord), [VS Code](../vscode/), [TRIPOL-AI](../TRIPOL-AI.md) |
| Validación | Regresiones y artefactos Linux, Windows y VSIX | Evidencia por commit | [Workflows](https://github.com/Tripollas/Tripol/actions), [empaquetado](../.github/workflows/public-ecosystem.yml) |
| Validación | Prueba contra Discord, pruebas prolongadas y demostración integral con slash, botones y bienvenida | Pendiente | [Guía de Discord](discord.md) |
| Release | Release experimental final | Pendiente | [Public ecosystem](../.github/workflows/public-ecosystem.yml) |
