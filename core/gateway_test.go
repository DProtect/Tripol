package core

import (
	"bytes"
	pecheck "debug/pe"
	"strings"
	"testing"
)

func TestNexusGatewayConfiguration(t *testing.T) {
	src := `import tripol.discord
import tripol.json
import tripol.memory
show discord.setup("dummy.token","123")
show discord.status("idle"),discord.activity("TRIPOL"),discord.welcome("Ready!")
show discord.connected(),discord.ping()
scratch {
 raw=discord.identify_payload()
 show json.valid(raw),json.get_raw(raw,"op")
 d=json.get_raw(raw,"d")
 show json.get_raw(d,"intents")
 p=json.get_raw(d,"presence")
 show json.get_raw(p,"status"),json.get_raw(p,"activities")
}
show discord.heartbeat_payload(-1),discord.heartbeat_payload(42)
before=memory.used()
repeat 1000 {discord.connected();discord.ping();discord.status("online")}
show before==memory.used()
discord.reset()
show discord.configured(),discord.connected(),discord.ping()`
	want := "true\ntrue true true\nfalse -1\ntrue 2\n33281\n\"idle\" [{\"name\":\"TRIPOL\",\"type\":0}]\n{\"op\":1,\"d\":null} {\"op\":1,\"d\":42}\ntrue\nfalse false -1\n"
	for _, p := range []string{"nano", "balanced", "turbo", "debug"} {
		t.Run(p, func(t *testing.T) {
			got, bad := hyperRun(t, src, p)
			if bad || got != want {
				t.Fatalf("%q %v", got, bad)
			}
		})
	}
}
func TestNexusGatewayDispatch(t *testing.T) {
	src := `import tripol.discord
import tripol.json
 discord.channel("123")
 discord.reply("!ping","Pong!")
 show discord.gateway_reply("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"s\":1,\"d\":{\"channel_id\":\"123\",\"content\":\"!ping\",\"author\":{\"id\":\"45\"}}}")
 show discord.gateway_reply("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"999\",\"content\":\"!ping\",\"author\":{}}}")==""
 show discord.gateway_reply("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"!ping\",\"author\":{\"bot\":true}}}")==""
 show discord.gateway_reply("{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"!ping\",\"author\":{},\"webhook_id\":\"42\"}}") == ""
 show discord.gateway_reply("{bad")==""
 large="{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"!ping\",\"author\":{},\"extra\":\""
 repeat 5000 {large+="a"}
 large+="\"}}"
 show json.valid(large),json.get_raw_large(large,"op"),discord.gateway_reply(large)`
	runCase(t, src, "Pong!\ntrue\ntrue\ntrue\ntrue\nfalse 0 Pong!\n", false)
}
func TestDiscordGatewayDecodesCommandText(t *testing.T) {
	source := `import tripol.discord
import tripol.json
import tripol.memory

command="!café/😀"
discord.reply(command,"Pong!")
message="{\"content\":\"!caf\\u00e9\\/\\ud83d\\ude00\",\"author\":{}}"
event="{\"op\":0,\"t\":\"MESSAGE_CREATE\",\"d\":{\"channel_id\":\"123\",\"content\":\"!caf\\u00E9\\/\\uD83D\\uDe00\",\"author\":{}}}"
show discord.matches(message,command),discord.dispatch(message),discord.gateway_reply(event)
show json.string_equals("\"\\n\\t\\r\"", "\n\t\r")
show json.string_equals("\"\\uD800\"", "�"),json.string_equals("\"\\uDC00\"", "�")
show json.string_equals("\"\\uD800\\u0041\"", "𐀀"),json.string_equals("\"\\u00x0\"", "Z")
show json.string_equals("\"\\u0000\"", ""),json.string_equals("\"a\" \"", "a")
show json.string_equals("\"ab\"", "a"),json.string_equals("\"a\"", "ab"),json.string_equals("\"\"", "")
before=memory.used()
repeat 2000 {assert json.string_equals("\"!caf\\u00E9\\/\\uD83D\\uDe00\"",command)}
show before==memory.used()`
	const want = "true Pong! Pong!\ntrue\nfalse false\nfalse false\nfalse false\nfalse false true\ntrue\n"
	for _, profile := range []string{"debug", "nano", "balanced", "turbo"} {
		t.Run(profile, func(t *testing.T) {
			got, failed := hyperRun(t, source, profile)
			if failed || got != want {
				t.Fatalf("got %q (failed=%v), want %q", got, failed, want)
			}
		})
	}
}
func TestNexusGatewayImports(t *testing.T) {
	src := `import tripol.discord
 discord.setup("dummy.token","123")
 discord.welcome("Ready!")
 discord.reply("!ping","Pong!")
 discord.run()`
	for _, profile := range []string{"nano", "balanced", "turbo", "debug"} {
		t.Run(profile, func(t *testing.T) {
			p, err := pecheck.NewFile(bytes.NewReader(compileProfile(src, "windows", profile)))
			if err != nil {
				t.Fatal(err)
			}
			defer p.Close()
			names, err := p.ImportedSymbols()
			if err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{"WinHttpWebSocketCompleteUpgrade:WINHTTP.dll", "WinHttpWebSocketSend:WINHTTP.dll", "WinHttpWebSocketReceive:WINHTTP.dll", "CreateThread:KERNEL32.dll", "WaitForSingleObject:KERNEL32.dll"} {
				if !strings.Contains(strings.Join(names, " "), name) {
					t.Fatal("missing", name)
				}
			}
		})
	}
}
func TestNexusVectorSearch(t *testing.T) {
	src := `import tripol.search
 for n in range(0,80) {
  a=ints(n)
  for i in range(n) {a[i]=i*4294967296+7}
  assert search.linear(a,-1)==-1
  for i in range(n) {assert search.linear(a,a[i])==i}
  if n>32 {a[29]=-9223372036854775807-1;a[n-1]=a[29];assert search.linear(a,a[29])==29}
 }
 show "ok"`
	for _, profile := range []string{"nano", "balanced", "turbo", "debug"} {
		t.Run(profile, func(t *testing.T) {
			got, bad := hyperRun(t, src, profile)
			if bad || got != "ok\n" {
				t.Fatalf("%q %v", got, bad)
			}
		})
	}
}
