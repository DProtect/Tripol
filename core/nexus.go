package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Nexus conveniences lower to existing native operations. No interpreter or extra runtime.
func (c *compiler) inputPrompt(e *expr) {
	if len(e.args) > 1 {
		c.err(e.s + " expects zero arguments or one text prompt")
	}
	if len(e.args) == 1 {
		c.require(c.expression(e.args[0]), "string")
		c.a.call("print:string")
	}
}
func (c *compiler) nexusBuiltin(e *expr) (string, bool) {
	a := c.a
	switch e.s {
	case "input_int":
		c.inputPrompt(e)
		a.call("input")
		a.call("parseint")
		return "int", true
	case "@array_len":
		v := c.arrayVar(e.args[0].s)
		if v.typ == "ints" || isList(v.typ) {
			a.load(v.off)
			if isList(v.typ) {
				a.emit(0x48, 0x8b, 0x40, 0x08)
			}
			a.emit(0x48, 0x8b, 0x00)
		} else {
			a.imm(int64(arrayLen(v)))
		}
		return "int", true
	case "abs":
		if len(e.args) != 1 {
			c.err("abs expects one integer")
		}
		c.require(c.expression(e.args[0]), "int")
		done := a.fresh()
		a.emit(0x48, 0x85, 0xc0)
		a.j(0x89, done)
		a.emit(0x48, 0xf7, 0xd8)
		a.j(0x80, "err")
		a.label(done)
		return "int", true
	case "min", "max":
		if len(e.args) != 2 {
			c.err(e.s + " expects two integers")
		}
		c.require(c.expression(e.args[0]), "int")
		a.emit(0x50)
		c.require(c.expression(e.args[1]), "int")
		a.emit(0x48, 0x89, 0xc1, 0x58, 0x48, 0x39, 0xc8)
		cc := byte(0x4f)
		if e.s == "max" {
			cc = 0x4c
		}
		a.emit(0x48, 0x0f, cc, 0xc1)
		return "int", true
	case "sum":
		if len(e.args) != 1 || e.args[0].op != "var" {
			c.err("sum expects a named integer array")
		}
		v := c.arrayVar(e.args[0].s)
		if isList(v.typ) && v.typ != "list" {
			c.err("sum requires integer elements")
		}
		if v.typ == "ints" || v.typ == "list" {
			a.load(v.off)
			if v.typ == "list" {
				a.emit(0x48, 0x8b, 0x40, 0x08)
			}
			a.emit(0x48, 0x8b, 0x08, 0x48, 0x8d, 0x50, 0x08) // rcx=len, rdx=data
		} else {
			a.emit(0xb9)
			a.u32(arrayLen(v))
			a.emit(0x48, 0x8d, 0x95)
			a.u32(v.off)
		}
		a.emit(0x31, 0xc0)
		loop, end := a.fresh(), a.fresh()
		a.label(loop)
		a.emit(0x48, 0x85, 0xc9)
		a.j(0x84, end)
		a.emit(0x48, 0x03, 0x02)
		a.j(0x80, "err")
		a.emit(0x48, 0x83, 0xc2, 8, 0x48, 0xff, 0xc9)
		a.jmp(loop)
		a.label(end)
		return "int", true
	}
	return "", false
}
func learn() {
	fmt.Print("TRIPOL " + version + ` — Learn in five minutes

1. Show values:
   show "Hello, TRIPOL!"
   show "Total:", 10 + 20

2. Read input:
   name = input("Your name?")
   age = input_int("Your age?")
   show "Hello", name

3. Decide:
   if age >= 18 { show "Adult" } else { show "Under 18" }

4. Repeat:
   repeat 3 { show "Hi" }
   for i in range(5) { show i }
   numbers = [10, 20, 30]
   for n in numbers { show n }
   show sum(numbers)
   show numbers

5. Save data:
   write_text("note.txt", "Hello")
   show read_text("note.txt")

6. Group data:
   struct Player { name: string; score: int }
   p = Player("Ana", 10)
   p.score += 5
   show p.name, p.score

7. Libraries: import tripol.math
   show math.square(7)
   tripol libs lists available modules.
   tripol docs sort shows sorting routines.
   tripol packages lists local Vault packages.

Save code as main.tri, then: tripol run main.tri
Build a standalone executable: tripol build main.tri
`)
}
func doctor() {
	exe, err := os.Executable()
	if err != nil {
		die("%v", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		die("%v", err)
	}
	fmt.Printf("TRIPOL %s\nCompiler: %s\nDirectory: %s\nHost: %s/%s\nTargets: windows/amd64, linux/amd64\nExternal compiler required: no\n", version, exe, cwd, runtime.GOOS, runtime.GOARCH)
	if runtime.GOARCH != "amd64" || !strings.Contains("|linux|windows|", "|"+runtime.GOOS+"|") {
		fmt.Println("Native run: unsupported host; cross-build is available")
	} else {
		fmt.Println("Native run: supported host")
	}
}

func newProjectTemplate(dir, template string) {
	if template == "basic" {
		newProject(dir)
		return
	}
	templates := map[string]string{
		"game": novaGameExample,
		"discord": `import tripol.discord

bot = discord.bot()
bot.prefix("!")
bot.activity("Programando en TRIPOL")
bot.command("hola", "Hola desde TRIPOL!")
bot.command("ping", "Pong!")
bot.start()
`,
		"discord-embeds": `import tripol.discord
import tripol.env

card = discord.embed()
card.title("TRIPOL")
card.description("Bienvenido al servidor")
card.color("#5865F2")
card.field("Versión", "0.28.0", true)

channel = env.get("DISCORD_CHANNEL_ID")
if channel == "" {show card.payload()} else {assert card.send(channel)}
`,
		"discord-components": `import tripol.discord
import tripol.env

panel = discord.components()
panel.text("## Bienvenido")
panel.separator()
panel.link("Documentación", "https://tripollas.github.io/Tripol/docs/")
message = discord.components()
message.container(panel, 5793266)

channel = env.get("DISCORD_CHANNEL_ID")
if channel == "" {show message.payload()} else {assert message.send(channel)}
`,
		"discord-slash": `import tripol.discord
import tripol.env

func open_feedback(ctx:discord.Interaction)->string {
    form=discord.modal("feedback_form","Cuéntanos")
    form.text("feedback","Sugerencia",2,true,4)
    return ctx.open_modal(form)
}

func receive_feedback(ctx:discord.Interaction)->string {
    if ctx.field("feedback")=="" {return ctx.reply("No se recibió la sugerencia.",true)}
    return ctx.reply("Gracias por tu sugerencia.",true)
}

bot=discord.bot()
guild=env.get("DISCORD_GUILD_ID")
bot.slash("feedback","Comparte una sugerencia",&open_feedback,guild)
bot.on_modal("feedback_form",&receive_feedback)
bot.start()
`,
		"discord-welcome": `import tripol.discord
import tripol.env
import tripol.json

func welcome(event:discord.GatewayEvent)->bool {
    if event.guild_id!=env.get("DISCORD_GUILD_ID") {return true}
    if event.user_id=="" or json.get_raw_large(event.field("user"),"bot")=="true" {return true}
    if discord.post_message(env.get("DISCORD_CHANNEL_ID"),"¡Bienvenido al servidor!")=="" {
        show "No se pudo enviar el mensaje de bienvenida."
        return false
    }
    return true
}

discord.guild(env.get("DISCORD_GUILD_ID"))
discord.channel_path(env.get("DISCORD_CHANNEL_ID"))
registry=discord.events()
registry.intent("GUILD_MEMBERS")
registry.on("member_join",&welcome)
bot=discord.bot()
bot.start_with_events(registry)
`,
		"ecosystem": `import tripol.sort
import tripol.time
import tripol.system
nums=[5,9,1,3]
start=time.now_ms()
show "Sorted:",sort.ascending(nums)
show "Platform:",system.platform()
show "Elapsed ms:",time.elapsed(start)
`,
		"libraries": `import tripol.math
import tripol.stats
numbers = [10, 20, 30]
show "Square:", math.square(7)
show "Mean:", stats.mean(numbers)
`,
		"calculator": `a = input_int("First number?")
b = input_int("Second number?")
show "Sum:", a + b
show "Product:", a * b
if b != 0 { show "Division:", float(a) / float(b) }
`,
		"records": `struct Player { name: string; score: int }
p = Player(input("Your name?"), 0)
p.score += 10
show "Player:", p.name, "Score:", p.score
`,
		"files": `name = input("Your name?")
write_text("profile.txt", "Name: " + name + "\n")
show "Saved:", read_text("profile.txt")
`,
	}
	source, ok := templates[template]
	if !ok {
		die("unknown project template %s; choose basic, calculator, files, records, libraries, ecosystem, game, discord, discord-embeds, discord-components, discord-slash or discord-welcome", template)
	}
	if err := os.Mkdir(dir, 0755); err != nil {
		die("cannot create project: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.tri"), []byte(source), 0644); err != nil {
		die("%v", err)
	}
	if template == "discord-embeds" || template == "discord-components" {
		guide := "# TRIPOL Discord\n\n" +
			"Run with `tripol run` on Windows or Linux x86-64. For sending, set DISCORD_CHANNEL_ID and DISCORD_TOKEN in the process environment. The bot needs View Channel and Send Messages. Embeds also need Embed Links. An unset channel ID opens the local JSON preview.\n\n" +
			"Ejecuta `tripol run` en Windows o Linux x86-64. Para enviar, configura DISCORD_CHANNEL_ID y DISCORD_TOKEN en el entorno del proceso. El bot necesita permisos para ver el canal y enviar mensajes. Los embeds también requieren insertar enlaces. Con el ID de canal vacío se muestra el JSON local.\n"
		if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(guide), 0644); err != nil {
			die("%v", err)
		}
	}
	if template == "discord-slash" {
		guide := "# Discord slash bot\n\n" +
			"Use an authorized bot application with the applications.commands scope. Set DISCORD_TOKEN and DISCORD_GUILD_ID for a test server, then run `tripol run`. Invoke `/feedback` to open the form and receive a private acknowledgment. An empty guild ID registers the command globally.\n\n" +
			"Usa una aplicación de bot autorizada con el alcance applications.commands. Configura DISCORD_TOKEN y DISCORD_GUILD_ID para un servidor de pruebas y ejecuta `tripol run`. Usa `/feedback` para abrir el formulario y recibir confirmación privada. Si el ID de servidor está vacío, el registro es global.\n"
		if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(guide), 0644); err != nil {
			die("%v", err)
		}
	}
	if template == "discord-welcome" {
		guide := "# Discord welcome bot\n\n" +
			"Set DISCORD_TOKEN, DISCORD_GUILD_ID and DISCORD_CHANNEL_ID in the process environment. Use a channel in that server. Enable Server Members Intent in the Discord Developer Portal and grant View Channel and Send Messages in the welcome channel. Start with `tripol run` on Linux x86-64 or Windows x64. The bot sends a welcome message when a human joins the configured server and logs a delivery failure while continuing to handle events.\n\n" +
			"Configura DISCORD_TOKEN, DISCORD_GUILD_ID y DISCORD_CHANNEL_ID en el entorno del proceso. Usa un canal de ese servidor. Activa Server Members Intent en el portal de desarrolladores de Discord y concede permisos para ver el canal y enviar mensajes. Ejecuta `tripol run` en Linux x86-64 o Windows x64. El bot da la bienvenida cuando una persona se une al servidor configurado y sigue procesando eventos si falla el envío.\n"
		if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(guide), 0644); err != nil {
			die("%v", err)
		}
	}
	fmt.Printf("Created %s (%s)\ncd %s\ntripol run\n", dir, template, dir)
}
