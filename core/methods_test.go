package core

import (
	"bytes"
	pecheck "debug/pe"
	"encoding/json"
	"strings"
	"testing"
)

func TestNativeTypedMethods(t *testing.T) {
	source := `struct Counter {value:int=0}
struct Other {value:int}
func Counter.add(self:Counter,step=1)->int {
 self.value+=step
 return self.value
}
func Other.add(self:Other,step:int)->int {return self.value+step}
func fresh()->Counter {show "new";return Counter()}
counter=Counter()
alias=counter
show counter.add(),alias.add(3),counter.value
show fresh().add()
show Other(7).add(2)`
	const want = "1 4 4\nnew\n1\n9\n"
	for _, profile := range []string{"debug", "nano", "balanced", "turbo"} {
		t.Run(profile, func(t *testing.T) {
			got, failed := hyperRun(t, source, profile)
			if failed || got != want {
				t.Fatalf("got %q (failed=%v), want %q", got, failed, want)
			}
		})
	}
}

func TestDiscordNewObjects(t *testing.T) {
	const source = `import tripol.discord
a=discord.bot()
b=discord.bot()
a.command("hola","Primera")
a.command("hola","Hola")
a.command("ping","Pong")
b.prefix("?")
b.command("hola","Otra")
show a.preview("hola"),b.preview("?hola")
a.prefix("!")
show a.preview("!hola"),a.preview("hola")=="",a.preview("!ping")
embed=discord.embed()
embed.title("Hola, Ada")
embed.description("Una \"cita\" y café")
embed.color("#Aa11Ff")
embed.footer("DISCORD REVOLUTION")
embed.field("Versión","0.28.0",true)
embed.field("Estado","Experimental")
show embed.payload()`
	for _, profile := range []string{"debug", "nano", "balanced", "turbo"} {
		t.Run(profile, func(t *testing.T) {
			got, failed := hyperRun(t, source, profile)
			before := "Hola Otra\nHola true Pong\n"
			if failed || !strings.HasPrefix(got, before) {
				t.Fatalf("unexpected output: %q", got)
			}
			var body struct {
				Embeds []struct {
					Description string
					Color       int
					Fields      []struct {
						Name, Value string
						Inline      bool
					}
					Footer struct{ Text string }
				}
				Mentions struct{ Parse []string } `json:"allowed_mentions"`
			}
			if err := json.Unmarshal([]byte(strings.TrimPrefix(got, before)), &body); err != nil {
				t.Fatal(err, got)
			}
			if len(body.Embeds) != 1 || body.Mentions.Parse == nil || len(body.Mentions.Parse) != 0 {
				t.Fatal(got)
			}
			embed := body.Embeds[0]
			if embed.Description != "Una \"cita\" y café" || embed.Color != 0xaa11ff || embed.Footer.Text != "DISCORD REVOLUTION" ||
				len(embed.Fields) != 2 || embed.Fields[0].Name != "Versión" || embed.Fields[0].Value != "0.28.0" ||
				!embed.Fields[0].Inline || embed.Fields[1].Inline {
				t.Fatal(got)
			}
		})
	}
}

func TestDiscordNewNativeImports(t *testing.T) {
	for _, tc := range []struct {
		source, symbol string
	}{
		{`bot=discord.bot();bot.command("hola","hola");bot.start()`, "WinHttpWebSocketReceive:WINHTTP.dll"},
		{`embed=discord.embed();embed.title("hola");embed.send("123")`, "WinHttpSendRequest:WINHTTP.dll"},
	} {
		pe, err := pecheck.NewFile(bytes.NewReader(compileProfile("import tripol.discord\n"+tc.source, "windows", "turbo")))
		if err != nil {
			t.Fatal(err)
		}
		defer pe.Close()
		symbols, err := pe.ImportedSymbols()
		if err != nil || !strings.Contains(strings.Join(symbols, " "), tc.symbol) {
			t.Fatal(symbols, err)
		}
	}
}

func TestDiscordNewRecoverableLimits(t *testing.T) {
	runCase(t, `import tripol.discord
bot=discord.bot()
for i in range(16) {bot.command(text(i),"Hi")}
try {bot.command("extra","Hi")} catch code {show code}
show bot.preview("15"),bot.preview("extra")==""
embed=discord.embed()
try {embed.color("#zzzzzz")} catch code {show code}
try {embed.payload()} catch code {show code}
`, "12\nHi true\n12\n12\n", false)
}
