package core

import (
	"encoding/json"
	"strings"
	"testing"
)

const slashEvent = `{"op":0,"t":"INTERACTION_CREATE","s":4,"d":{"id":"710","application_id":"890","guild_id":"321","channel_id":"123","type":2,"token":"interaction-fixture","data":{"type":1,"name":"hola","guild_id":"321","options":[{"type":1,"name":"say","options":[{"type":3,"name":"text","value":"caf\u00e9 \ud83d\ude00"}]}]},"member":{"user":{"id":"567"}}}}`
const buttonEvent = `{"op":0,"t":"INTERACTION_CREATE","s":5,"d":{"id":"711","application_id":"890","type":3,"token":"button-fixture","message":{"id":"99"},"data":{"component_type":2,"custom_id":"caf\u00e9"}}}`
const selectEvent = `{"op":0,"t":"INTERACTION_CREATE","s":6,"d":{"id":"712","application_id":"890","type":3,"token":"select-fixture","message":{"id":"99"},"data":{"component_type":6,"custom_id":"roles","values":["101","102"]}}}`
const modalEvent = `{"op":0,"t":"INTERACTION_CREATE","s":7,"d":{"id":"713","application_id":"890","type":5,"token":"modal-fixture","data":{"custom_id":"ticket","components":[{"type":18,"component":{"type":4,"custom_id":"topic","value":"Revisi\u00f3n"}}]}}}`

const interactionBot = `import tripol.discord
func greet(c:discord.Interaction)->string {
 return c.reply(c.user_id+" "+c.option_text("text"),true)
}
func defer_reply(c:discord.Interaction)->string {return c.defer(true)}
func button(c:discord.Interaction)->string {
 form=discord.modal("ticket","Cuéntanos")
 form.help_text("Describe el problema")
 form.text("topic","Tema",2,true,2,4000)
 return c.open_modal(form)
}
func choose_role(c:discord.Interaction)->string {
 list=c.values()
 return c.reply(list[0]+" "+list[1],true)
}
func submit(c:discord.Interaction)->string { return c.reply(c.field("topic"),true) }
b=discord.bot()
opts=discord.slash_options()
opts.add("text","La opción","string",true)
b.slash_with_options("hola","La descripción",&greet,opts,"321")
b.slash("wait","Espera",&defer_reply,"321")
b.on_button("café",&button)
b.on_select("roles",&choose_role)
b.on_modal("ticket",&submit)
`

func TestDiscordTypedInteractionPreviews(t *testing.T) {
	events := []string{slashEvent, buttonEvent, selectEvent, modalEvent, strings.Replace(slashEvent, `"name":"hola"`, `"name":"wait"`, 1)}
	src := interactionBot
	for _, event := range events {
		value, _ := json.Marshal(event)
		src += "show b.preview_interaction(" + string(value) + ")\n"
	}
	src += `other=discord.bot()
show other.preview_interaction("{}")==""`
	for _, profile := range []string{"debug", "nano", "balanced", "turbo"} {
		t.Run(profile, func(t *testing.T) {
			out, failed := hyperRun(t, src, profile)
			if failed {
				t.Fatal(out)
			}
			lines := strings.Split(strings.TrimSpace(out), "\n")
			if len(lines) != 6 || lines[5] != "true" {
				t.Fatal(out)
			}
			for i, expected := range []struct {
				typ     int
				content string
			}{{4, "567 café 😀"}, {9, ""}, {4, "101 102"}, {4, "Revisión"}, {5, ""}} {
				var callback struct {
					Type int
					Data struct {
						Content    string
						Flags      int
						CustomID   string `json:"custom_id"`
						Components []json.RawMessage
					}
				}
				if err := json.Unmarshal([]byte(lines[i]), &callback); err != nil || callback.Type != expected.typ || callback.Data.Content != expected.content {
					t.Fatalf("callback %d: %s (%v)", i, lines[i], err)
				}
				if expected.typ != 9 && callback.Data.Flags != 64 {
					t.Fatal(lines[i])
				}
				if expected.typ == 9 && (callback.Data.CustomID != "ticket" || len(callback.Data.Components) != 2 || !strings.Contains(lines[i], `"type":18,"label":"Tema","component":{"type":4`)) {
					t.Fatal(lines[i])
				}
			}
		})
	}
	if _, err := CompileWeb(src); err != nil {
		t.Fatal(err)
	}
}

func TestDiscordInteractionBuildersAndOfficialTextLimits(t *testing.T) {
	const source = `import tripol.discord
import tripol.json
menu=discord.select("team","string",0,2)
menu.option("Rojo","r")
menu.option("Azul","b","Opción",true)
groups=discord.select("roles","role",0,2)
groups.default_id("101","role")
channels=discord.select("channels","channel")
channels.channels(ints(0))
layout=discord.components()
layout.select(menu)
layout.select(groups)
layout.select(discord.select("people","user"))
layout.select(discord.select("mention","mentionable"))
row=discord.button_row()
row.button("Aceptar","confirm",3,true)
row.link("Docs","https://example.com")
layout.row(row)
form=discord.modal("choose","Configurar")
form.select("Canal",channels)
form.file("upload","Un archivo",0,2,false)
a=discord.embed()
a.title("Una tarjeta")
a.author("Ada","https://example.com","https://example.com/ada.png")
a.timestamp("2028-02-29T03:12:48.001-05:00")
a.footer("Fin","https://example.com/footer.png")
b=discord.embed()
b.description("Segundo embed")
batch=discord.embeds()
batch.add(a)
batch.add(b)
show layout.payload()
show form.payload()
show batch.payload()
show json.unquote("\"\\u00e9\\u20ac\\ud83d\\ude00\"")
text_value=""
repeat 4000 {text_value+="x"}
large=discord.components()
show large.text(text_value)
try {large.text("x")} catch code {show code}
try {a.timestamp("2027-02-29T03:12:48Z")} catch code {show code}
try {layout.select(menu)} catch code {show code}
`
	out, failed := hyperRun(t, source, "turbo")
	if failed {
		t.Fatal(out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 8 || strings.Join(lines[3:], "\n") != "é€😀\ntrue\n12\n12\n12" {
		t.Fatal(out)
	}
	if !strings.Contains(lines[0], `"type":5,"custom_id":"people"`) || !strings.Contains(lines[0], `"type":7,"custom_id":"mention"`) ||
		!strings.Contains(lines[0], `"custom_id":"confirm","disabled":true`) || !strings.Contains(lines[1], `"type":19,"custom_id":"upload"`) {
		t.Fatal(out)
	}
	for i, key := range []string{`"default_values":[{"id":"101","type":"role"}]`, `"channel_types":[]`, `"timestamp":"2028-02-29T03:12:48.001-05:00"`} {
		if !json.Valid([]byte(lines[i])) || !strings.Contains(lines[i], key) {
			t.Fatal(lines[i])
		}
	}
}

func TestCallbackListsReadOnlyAcrossScratch(t *testing.T) {
	runCase(t, `func first(n:int)->int {return n+1}
func second(n:int)->int {return n*2}
let callbacks:list[func(int)->int]=list([&first,&second])
scratch {
 h=callbacks.get(1)
 show h(21),callbacks.size()
}`, "42 2\n", false)
	bad := `func f(n:int)->int {return n}
a=list([&f])
scratch {a.add(&f)}`
	if _, err := CompileWeb(bad); err == nil || !strings.Contains(err.Error(), "scratch cannot access outer aggregate") {
		t.Fatalf("unexpected diagnostic: %v", err)
	}
}

func TestNativeQuoteAllocationAndInvalidUnicode(t *testing.T) {
	src := `import tripol.json
import tripol.memory
s=REPLACE
before=memory.used()
encoded=json.quote(s)
after=memory.used()
show after-before<=len(encoded)+8
show json.unquote(encoded)==s
show json.unquote("\"\\ud800\"")=="" and json.unquote("\"\\ud83d\\u0041\"")=="" and json.unquote("\"\\u0000\"")==""
`
	literal, _ := json.Marshal(strings.Repeat("á😀\"\\\n", 400))
	runCase(t, strings.Replace(src, "REPLACE", string(literal), 1), "true\ntrue\ntrue\n", false)
}

func TestDiscordModalChoiceGroups(t *testing.T) {
	src := `import tripol.discord
colors=discord.select("colors")
colors.option("Rojo","red","",true)
colors.option("Azul","blue")
days=discord.select("days","string",0,2)
days.option("Lunes","mo")
days.option("Martes","tu")
form=discord.modal("feedback","Comentarios")
form.radio("Color",colors,true,"Un color")
form.checkbox_group("Días",days,false)
form.checkbox("agree","Acepto",true)
show form.payload()
ctx=discord.Interaction(5,"feedback","","","","","","{\"components\":[{\"type\":18,\"component\":{\"type\":21,\"custom_id\":\"colors\",\"value\":\"blue\"}},{\"type\":18,\"component\":{\"type\":22,\"custom_id\":\"days\",\"values\":[\"mo\",\"tu\"]}},{\"type\":18,\"component\":{\"type\":23,\"custom_id\":\"agree\",\"value\":true}}]}",false)
show ctx.field("colors"),ctx.field_data("days"),ctx.checked("agree"),ctx.checked("missing")`
	out, failed := hyperRun(t, src, "nano")
	if failed {
		t.Fatal(out)
	}
	payload, suffix, _ := strings.Cut(out, "\n")
	if !json.Valid([]byte(payload)) || !strings.Contains(payload, `"type":21,"custom_id":"colors"`) ||
		!strings.Contains(payload, `"type":22,"custom_id":"days"`) || !strings.Contains(payload, `"type":23,"custom_id":"agree","default":true`) ||
		suffix != "blue [\"mo\",\"tu\"] true false\n" {
		t.Fatal(out)
	}
	if _, err := CompileWeb(src); err != nil {
		t.Fatal(err)
	}
}

func TestDiscordSlashRoutesAreScopedByRegistration(t *testing.T) {
	runCase(t, `import tripol.discord
func global(c:discord.Interaction)->string {return "global "+c.name}
func local(c:discord.Interaction)->string {return "local "+c.name}
func another(c:discord.Interaction)->string {return "another "+c.name}
b=discord.bot()
b.slash("hola","Global",&global)
b.slash("hola","Local",&local,"321")
b.slash("hola","Another",&another,"654")
show b.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"type\":2,\"guild_id\":\"321\",\"data\":{\"type\":1,\"name\":\"hola\"}}}")
show b.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"type\":2,\"guild_id\":\"321\",\"data\":{\"type\":1,\"name\":\"hola\",\"guild_id\":\"321\"}}}")
show b.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"type\":2,\"guild_id\":\"654\",\"data\":{\"type\":1,\"name\":\"hola\",\"guild_id\":\"654\"}}}")
show b.preview_interaction("{\"op\":0,\"t\":\"INTERACTION_CREATE\",\"d\":{\"type\":2,\"data\":{\"type\":1,\"name\":\"hola\",\"guild_id\":\"111\"}}}")==""
`, "global hola\nlocal hola\nanother hola\ntrue\n", false)
}

func TestDiscordLiveInteractionStateIsPrivate(t *testing.T) {
	context := `import tripol.discord
ctx=discord.Interaction(2,"hola","1","2","","","","{}",false`
	for _, src := range []string{
		context + ")\nshow ctx.__interaction_token",
		context + ")\nctx.__response_state=2",
		context + `,"literal-credential")`,
	} {
		if _, err := CompileWeb(src); err == nil || !strings.Contains(err.Error(), "private field discord.Interaction") {
			t.Fatalf("unexpected diagnostic: %v", err)
		}
	}
}
