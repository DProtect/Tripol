package core

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDiscordChannelModerationValidation(t *testing.T) {
	runCase(t, `import tripol.discord
room=discord.channel_ref("123")
show room.id()
try {room.purge(list(["5"]))} catch code {show code}
try {room.purge(list(["5","5"]))} catch code {show code}
try {room.delete_message("../5")} catch code {show code}
try {room.slowmode(21601)} catch code {show code}
try {room.react("5","")} catch code {show code}
`, "123\n12\n12\n12\n12\n12\n", false)
}

const discordComponentsDemo = `import tripol.discord
panel=discord.components()
panel.text("¡Hola! \"TRIPOL\"")
panel.separator(2,false)
panel.link("Read docs","https://example.org/docs")
panel.button("Aceptar","accept",3)
panel.section("Thumbnail","https://example.org/photo.png","Descripción")
panel.gallery(string_list(["https://example.org/a.png","https://example.org/b.png"]))
message=discord.components()
message.container(panel,255,true)
show message.payload()
embed=discord.embed()
embed.title("Awesome")
embed.url("https://example.org/")
embed.image("https://example.org/a.png")
embed.thumbnail("https://example.org/b.png")
show embed.payload()`

func TestDiscordComponentsPayload(t *testing.T) {
	for _, profile := range []string{"debug", "nano", "balanced", "turbo"} {
		t.Run(profile, func(t *testing.T) {
			out, failed := hyperRun(t, discordComponentsDemo, profile)
			if failed {
				t.Fatal(out)
			}
			first, second, _ := strings.Cut(out, "\n")
			var payload struct {
				Flags      int
				Components []struct {
					Type        int
					AccentColor int `json:"accent_color"`
					Spoiler     bool
					Components  []struct {
						Type, Spacing int
						Divider       bool
						Content       string
						Accessory     *struct {
							Type        int
							Description string
						}
						Items      []struct{ Media struct{ URL string } }
						Components []struct {
							Type, Style int
							URL         string
							ID          string `json:"custom_id"`
						}
					}
				}
				Allowed struct{ Parse []string } `json:"allowed_mentions"`
			}
			if err := json.Unmarshal([]byte(first), &payload); err != nil {
				t.Fatal(err)
			}
			if payload.Flags != 32768 || payload.Allowed.Parse == nil || len(payload.Components) != 1 {
				t.Fatal(first)
			}
			panel := payload.Components[0]
			if panel.Type != 17 || panel.AccentColor != 255 || !panel.Spoiler || len(panel.Components) != 6 {
				t.Fatal(first)
			}
			children := panel.Components
			if children[0].Type != 10 || children[0].Content != "¡Hola! \"TRIPOL\"" ||
				children[1].Type != 14 || children[1].Spacing != 2 || children[1].Divider ||
				children[2].Type != 1 || children[2].Components[0].Style != 5 || children[2].Components[0].URL != "https://example.org/docs" ||
				children[3].Type != 1 || children[3].Components[0].Style != 3 || children[3].Components[0].ID != "accept" ||
				children[4].Type != 9 || children[4].Accessory.Type != 11 || children[4].Accessory.Description != "Descripción" ||
				children[5].Type != 12 || len(children[5].Items) != 2 || children[5].Items[1].Media.URL != "https://example.org/b.png" {
				t.Fatal(first)
			}
			if !json.Valid([]byte(second)) || !strings.Contains(second, `"image":{"url":"https://example.org/a.png"}`) ||
				!strings.Contains(second, `"thumbnail":{"url":"https://example.org/b.png"}`) {
				t.Fatal(second)
			}
		})
	}
}

func TestDiscordComponentsValidation(t *testing.T) {
	runCase(t, `import tripol.discord
m=discord.components()
try {m.payload()} catch code {show code}
try {m.link("link","javascript:alert(1)")} catch code {show code}
try {m.separator(9)} catch code {show code}
m.button("A","a")
try {m.button("B","a")} catch code {show code}
p=discord.components()
p.button("C","a")
try {m.container(p)} catch code {show code}
c=discord.components()
c.container(m)
try {p.container(c)} catch code {show code}
repeat 38 {m.separator()}
try {m.text("over")} catch code {show code}
show m.total
`, "12\n12\n12\n12\n12\n12\n12\n40\n", false)
}

func TestDiscordCallbackPreview(t *testing.T) {
	source := `import tripol.discord
func answer(message:discord.Message)->string {
 if message.matches("ping") {return message.reply(message.author_id+" Pong!")}
 return ""
}
bot=discord.bot()
bot.on_message(&answer)
show bot.preview_event("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"ping\",\"author\":{\"id\":\"456\"}}}")
show bot.preview_event("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"ping\",\"author\":{\"id\":\"456\",\"bot\":true}}}")==""
show bot.preview_event("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"ping\",\"webhook_id\":\"789\",\"author\":{\"id\":\"456\"}}}")==""
`
	runCase(t, source, "456 Pong!\ntrue\ntrue\n", false)
	if _, err := CompileWeb(source); err != nil {
		t.Fatal(err)
	}
}

func TestDiscordButtonAndSelectEmoji(t *testing.T) {
	const source = `import tripol.discord
panel=discord.container()
row=discord.button_row()
row.button("Go","go","primary",false,"✅")
row.button("Later","later","secondary",true)
row.button("Accept","approve","success",false,"<:check:123>")
row.button("","remove","danger",false,"<a:party:456>")
row.link("","https://example.org",false,"🔗")
assert row.disable("go")
assert row.disable("later",false)
assert not row.disable("missing")
panel.row(row)
panel.premium("789",true)
panel.button("Other","other","success","456",true)
menu=discord.select("colors")
menu.option("Blue","blue","Available",false,"💙")
menu.disable()
panel.select(menu)
show panel.payload()
form=discord.modal("upload","Archivos")
form.file("files","Documento",list(["image",".pdf"]),0,2,false,"Opcional")
show form.payload()
`
	out, failed := hyperRun(t, source, "turbo")
	if failed {
		t.Fatal(out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 || !json.Valid([]byte(lines[0])) || !json.Valid([]byte(lines[1])) {
		t.Fatal(out)
	}
	for _, check := range []string{
		`"style":1,"label":"Go","custom_id":"go","emoji":{"name":"✅"},"disabled":true`,
		`"style":2,"label":"Later","custom_id":"later"}`,
		`"style":3,"label":"Accept","custom_id":"approve","emoji":{"id":"123","name":"check","animated":false}`,
		`"style":4,"custom_id":"remove","emoji":{"id":"456","name":"party","animated":true}`,
		`"style":5,"url":"https://example.org","emoji":{"name":"🔗"}`,
		`"style":6,"sku_id":"789","disabled":true`,
		`"emoji":{"id":"456"},"disabled":true`,
		`"emoji":{"name":"💙"}}],"disabled":true`,
	} {
		if !strings.Contains(lines[0], check) {
			t.Fatalf("missing %s in %s", check, lines[0])
		}
	}
	if !strings.Contains(lines[1], `"description":"Opcional","component":{"type":19`) || !strings.Contains(lines[1], `"file_types":["image",".pdf"]`) {
		t.Fatal(lines[1])
	}
	if _, err := CompileWeb(source); err != nil {
		t.Fatal(err)
	}
	compileProfile(source, "windows", "turbo")
}

func TestDiscordRejectsInvalidComponentCustomization(t *testing.T) {
	runCase(t, `import tripol.discord
row=discord.button_row()
try {row.button("","empty")} catch code {show code}
try {row.button("OK","x","primary",false,"<:ok:abc>")} catch code {show code}
try {row.link("","https://example.org")} catch code {show code}
form=discord.modal("upload","Archivos")
try {form.file("files","Documento",list(["pdf"]))} catch code {show code}
menu=discord.select("colors")
menu.option("Blue","blue","",false,"💙")
menu.disable()
try {form.select("Color",menu)} catch code {show code}
menu.disable(false)
try {form.checkbox_group("Color",menu)} catch code {show code}
`, "12\n12\n12\n12\n12\n12\n", false)
}
