package core

import (
	"bytes"
	pecheck "debug/pe"
	"fmt"
	"strings"
	"testing"
)

func TestEclipseBotState(t *testing.T) {
	src := `import tripol.discord
import tripol.memory
show discord.configured()
scratch {discord.login("test_token"+".ABC");discord.channel("123"+"456");discord.reply("!hola","Ho"+"la")}
show discord.configured(),discord.preview("!hola"),discord.preview("!missing")==""
old=discord.preview("!hola")
discord.reply("!hola","Actualizado")
discord.reply("!ping","Pong!")
show old,discord.preview("!hola"),discord.preview("!ping")
before=memory.used()
repeat 2000 {discord.reply("!ping","Pong!");discord.configured()}
show before==memory.used()
discord.reset()
show discord.configured(),discord.preview("!hola")==""
show discord.message_id("{\"id\":\"9999999999999999999\"}")`
	for _, p := range []string{"nano", "balanced", "turbo", "debug"} {
		t.Run(p, func(t *testing.T) {
			got, bad := hyperRun(t, src, p)
			want := "false\ntrue Hola true\nHola Actualizado Pong!\ntrue\nfalse true\n0\n"
			if bad || got != want {
				t.Fatalf("%q %v", got, bad)
			}
		})
	}
}
func TestEclipseBotCapacity(t *testing.T) {
	var src strings.Builder
	src.WriteString("import tripol.discord\n")
	for i := 0; i < 16; i++ {
		fmt.Fprintf(&src, "discord.reply(\"!%d\",\"%d\")\n", i, i)
	}
	src.WriteString(`show discord.preview("!15")
try {discord.reply("!extra","extra")} catch code {show code}
show discord.preview("!0"),discord.preview("!extra")==""`)
	runCase(t, src.String(), "15\n12\n0 true\n", false)
}
func TestEclipseBotInvalid(t *testing.T) {
	src := `import tripol.discord
try {discord.login("bad token")} catch code {show code}
try {discord.channel("123?x")} catch code {show code}
try {discord.reply("","a")} catch code {show code}
try {discord.run()} catch code {show code}
show discord.configured()`
	runCase(t, src, "12\n12\n12\n12\nfalse\n", false)
}
func TestEclipseBotPE(t *testing.T) {
	src := `import tripol.discord
discord.login("dummy.test")
discord.channel("123")
discord.say("Ready!")
discord.reply("!hola","Hola")
discord.reply("!ping","Pong!")
discord.run()`
	bin := compileProfile(src, "windows", "turbo")
	p, err := pecheck.NewFile(bytes.NewReader(bin))
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	s, err := p.ImportedSymbols()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(s, " "), "WinHttpSendRequest:WINHTTP.dll") {
		t.Fatal(s)
	}
	for _, sec := range p.Sections {
		if sec.Name == ".idata" {
			if sec.VirtualSize-sec.Size != botStateSize+65584 {
				t.Fatalf("wrong bot storage: %d %d", sec.VirtualSize, sec.Size)
			}
			return
		}
	}
	t.Fatal("missing writable storage")
}
func TestEclipseArrays(t *testing.T) {
	src := `import tripol.arrays
import tripol.search
a=[-9223372036854775807-1,0,9,9,9223372036854775807]
b=arrays.copy(a)
b[2]=4
show a,b
show search.linear(a,9),search.linear(b,9),search.linear(a,-9223372036854775807-1),search.linear(a,9223372036854775807),search.linear(a,7)
show arrays.contains(a,9),arrays.contains(a,7)
e=ints(0);show arrays.copy(e),search.linear(e,0),arrays.contains(e,0)`
	for _, p := range []string{"nano", "balanced", "turbo", "debug"} {
		t.Run(p, func(t *testing.T) {
			got, bad := hyperRun(t, src, p)
			want := "[-9223372036854775808, 0, 9, 9, 9223372036854775807] [-9223372036854775808, 0, 4, 9, 9223372036854775807]\n2 3 0 4 -1\ntrue false\n[] -1 false\n"
			if bad || got != want {
				t.Fatalf("%q %v", got, bad)
			}
		})
	}
}

func TestEclipseDispatch(t *testing.T) {
	runCase(t, `import tripol.discord
 discord.reply("!ping","Pong!")
 show discord.dispatch("{\"id\":\"12\",\"content\":\"!ping\",\"author\":{\"id\":\"45\"}}")
 show discord.dispatch("{\"content\":\"!ping\",\"author\":{\"bot\":true}}") == ""
 show discord.dispatch("{\"content\":\"!ping\",\"author\":{},\"webhook_id\":\"42\"}") == ""
 show discord.dispatch("{}") == ""`, "Pong!\ntrue\ntrue\ntrue\n", false)
}
