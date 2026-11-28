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
const autocompleteEvent = `{"op":0,"t":"INTERACTION_CREATE","s":8,"d":{"id":"714","application_id":"890","type":4,"token":"autocomplete-fixture","guild_id":"321","data":{"type":1,"name":"inspect","guild_id":"321","options":[{"type":2,"name":"stock","options":[{"type":1,"name":"find","options":[{"name":"query","type":3,"value":"ca","focused":true}]}]}]}}}`

const commandTree = `
func inspect(c:discord.Interaction)->string {return c.command_group()+"/"+c.subcommand()+"/"+c.option_text("query")}
func suggest(c:discord.Interaction)->string {
 choices=discord.command_choices()
 if c.command_group()=="stock" and c.subcommand()=="find" and c.option_text(c.focused())=="ca" {choices.add("Café", "cafe")}
 return c.autocomplete(choices)
}
arguments=discord.slash_options()
arguments.add("query","Qué buscar","string",true,true)
limits=discord.command_choices()
limits.integer("Década",10)
arguments.choose("limit","Límite",limits)
subcommands=discord.slash_options()
subcommands.subcommand("find","Busca",arguments)
root=discord.slash_options()
root.group("stock","Inventario",subcommands)
b.slash_with_options("inspect","Inspecciona",&inspect,root,"321")
b.on_autocomplete("inspect",&suggest,"321")
`

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

func TestDiscordSlashTreesAndAutocomplete(t *testing.T) {
	preview, _ := json.Marshal(autocompleteEvent)
	command, _ := json.Marshal(strings.Replace(autocompleteEvent, `"type":4,"token"`, `"type":2,"token"`, 1))
	outside, _ := json.Marshal(strings.Replace(autocompleteEvent, `"name":"inspect","guild_id":"321"`, `"name":"inspect","guild_id":"123"`, 1))
	portable := "import tripol.discord\nb=discord.bot()" + commandTree + `
show b.preview_interaction(` + string(command) + `)
show b.preview_interaction(` + string(preview) + `)
show b.preview_interaction(` + string(outside) + `)==""
`
	src := portable + `
try {root.add("wrong","Cannot mix options")} catch code {show code}
try {subcommands.group("wrong","Only ordinary subcommands",root)} catch code {show code}
try {arguments.add("flag","Cannot complete booleans","boolean",false,true)} catch code {show code}
try {subcommands.subcommand("find","Duplicate",arguments)} catch code {show code}
choices=discord.command_choices()
choices.integer("Answer",42)
show choices.payload()
try {choices.add("Mixed types","bad")} catch code {show code}
numbers=discord.command_choices()
try {numbers.integer("Outside safe range",9007199254740993)} catch code {show code}
`
	runCase(t, src, "stock/find/ca\n{\"type\":8,\"data\":{\"choices\":[{\"name\":\"Café\",\"value\":\"cafe\"}]}}\ntrue\n12\n12\n12\n12\n[{\"name\":\"Answer\",\"value\":42}]\n12\n12\n", false)
	if _, err := CompileWeb(portable); err != nil {
		t.Fatal(err)
	}
	const bypass = `import tripol.discord
options=discord.slash_options()
options.__append("raw","Raw",1,"")`
	if _, err := CompileWeb(bypass); err == nil || !strings.Contains(err.Error(), "private method discord.SlashOptions") {
		t.Fatalf("unexpected diagnostic: %v", err)
	}
}

const userMenuEvent = `{"op":0,"t":"INTERACTION_CREATE","d":{"id":"741","token":"menu-fixture","type":2,"guild_id":"321","channel_id":"987","data":{"type":2,"name":"Perfil / ficha","guild_id":"321","target_id":"42","resolved":{"users":{"42":{"id":"42","username":"Ada"}}}}}}`
const messageMenuEvent = `{"op":0,"t":"INTERACTION_CREATE","d":{"id":"742","token":"menu-fixture","type":2,"guild_id":"321","channel_id":"987","data":{"type":3,"name":"Perfil / ficha","guild_id":"321","target_id":"43","resolved":{"messages":{"43":{"id":"43","content":"Un mensaje"}}}}}}`

const contextMenuHandlers = `
func profile(ctx) {
    ctx.reply(ctx.name + ": " + ctx.target_id())
}
func quote(ctx) {
    ctx.reply(ctx.target())
}
b.user_menu("Perfil / ficha",profile,"321","8192")
b.message_menu("Perfil / ficha",quote,"321","everyone")
`

func TestDiscordContextMenuDispatchAndPermissions(t *testing.T) {
	src := "import tripol.discord\nb=discord.bot()\n" + contextMenuHandlers
	for _, event := range []string{userMenuEvent, messageMenuEvent, strings.Replace(userMenuEvent, `"guild_id":"321","target_id"`, `"guild_id":"111","target_id"`, 1)} {
		encoded, _ := json.Marshal(event)
		src += "show b.preview_interaction(" + string(encoded) + ")\n"
	}
	out, failed := hyperRun(t, src, "turbo")
	if failed || !strings.Contains(out, `"content":"Perfil / ficha: 42"`) || !strings.Contains(out, `Un mensaje`) || !strings.HasSuffix(out, "\n\n") {
		t.Fatal(out)
	}
	if _, err := CompileWeb(src); err != nil {
		t.Fatal(err)
	}
	compileProfile(src, "windows", "turbo")
	runCase(t, `import tripol.discord
func handle(ctx:discord.Interaction)->string {return ""}
b=discord.bot()
try {b.user_menu("bad",handle,"321","08")} catch code {show code}
try {b.message_menu("bad",handle,"321","WRITE")} catch code {show code}
try {b.user_menu("",handle)} catch code {show code}
show discord.slash_is_current("[{\"type\":2,\"name\":\"Perfil\",\"default_member_permissions\":\"8\"}]","{\"type\":2,\"name\":\"Perfil\",\"default_member_permissions\":null}")
show discord.slash_is_current("[{\"type\":2,\"name\":\"Perfil\",\"default_member_permissions\":null}]","{\"type\":2,\"name\":\"Perfil\",\"default_member_permissions\":null}")
`, "12\n12\n12\nfalse\ntrue\n", false)
}

func TestDiscordTypedFormAnswersAndResolvedTargets(t *testing.T) {
	const source = `import tripol.discord

data="{\"type\":5,\"data\":{\"components\":[{\"type\":18,\"component\":{\"type\":4,\"custom_id\":\"empty\",\"value\":\"\"}},{\"type\":18,\"component\":{\"type\":6,\"custom_id\":\"roles\",\"values\":[\"101\",\"102\"]}},{\"type\":18,\"component\":{\"type\":19,\"custom_id\":\"files\",\"values\":[\"123\"]}}],\"resolved\":{\"attachments\":{\"123\":{\"id\":\"123\",\"filename\":\"report.pdf\",\"url\":\"https://cdn.example/report.pdf\",\"size\":42,\"content_type\":\"application/pdf\"}}}}}"
ctx=discord.interaction_context(data,"mreport")
roles=ctx.field_values("roles")
report=ctx.field_attachment("files")
missing=ctx.field_attachment("files",1)
other=ctx.field_attachment("roles")
no_values=ctx.field_values("missing")
show ctx.has_field("empty"),ctx.has_field("missing"),len(ctx.field("empty")),no_values.size()
show roles[0],roles[1]
show report.id,report.filename,report.url,report.size,report.content_type
show missing.id=="",other.id==""
u=discord.interaction_context("{\"type\":2,\"guild_id\":\"321\",\"member\":{\"permissions\":\"18446744073709551624\"},\"data\":{\"type\":2,\"target_id\":\"4\",\"resolved\":{\"users\":{\"4\":{\"username\":\"Ada\"}}}}}","uProfile")
m=discord.interaction_context("{\"type\":2,\"data\":{\"type\":3,\"target_id\":\"9\",\"resolved\":{\"messages\":{\"9\":{\"content\":\"Hola\"}}}}}","pQuote")
show u.target_username(),m.target_content(),u.target_content()=="",m.target_username()==""
show u.can("moderate_members"),u.can("unknown"),m.can("kick_members")
`
	runCase(t, source, "true false 0 0\n101 102\n123 report.pdf https://cdn.example/report.pdf 42 application/pdf\ntrue true\nAda Hola true true\ntrue false false\n", false)
	if _, err := CompileWeb(source); err != nil {
		t.Fatal(err)
	}
}

func TestDiscordFractionalChoicesAndSync(t *testing.T) {
	const portable = `import tripol.discord
choices=discord.command_choices()
choices.number("Fraction","1.250")
choices.number("Negative","-1.25e-1")
o=discord.slash_options()
o.choose("amount","Amount",choices,true)
show o.payload()
show discord.slash_number_signature("0.00123")==discord.slash_number_signature("123e-5")
show discord.slash_number_signature("-0.0e-2")==discord.slash_number_signature("0")
show discord.slash_number_signature("9007199254740992.00")!=""
show discord.slash_number_signature("9007199254740992.01")==""
show discord.slash_number_signature("1e+50")==""
show discord.slash_number_signature("1.0,0")==""
show discord.command_permission_payload("moderate_members")
`
	const expected = `[{"name":"amount","description":"Amount","type":10,"required":true,"choices":[{"name":"Fraction","value":1.250},{"name":"Negative","value":-1.25e-1}]}]
true
true
true
true
true
true
,"default_member_permissions":"1099511627776"
`
	runCase(t, portable, expected, false)
	if _, err := CompileWeb(portable); err != nil {
		t.Fatal(err)
	}
	runCase(t, `import tripol.discord
show discord.slash_number_signature("5e-324")!=""
show discord.slash_number_signature("1e-325")==""
c=discord.command_choices()
try {c.number("Invalid","[1.2]")} catch code {show code}
c.number("Real","1.5")
try {c.integer("Mixed",3)} catch code {show code}
`, "true\ntrue\n12\n12\n", false)
}

func TestDiscordJSONLookupValidatesEntireObject(t *testing.T) {
	runCase(t, `import tripol.json
show json.get_raw_large("  {\"id\":1,\"id\":2,\"other\":[true, {\"nested\":null}]}  ","id")
show json.get_raw_large("{\"id\":1,\"other\":[true,]}","id")==""
show json.get_raw_large("{\"id\":1,}","id")==""
show json.get_raw_large("{\"id\":1} {}","id")==""
show json.get_raw_large("{\"id\":1,\"other\":null}","missing")==""
`, "1\ntrue\ntrue\ntrue\ntrue\n", false)
}
