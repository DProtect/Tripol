package core

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"tripol/internal/discordtransport"
)

const hostedBot = `import tripol.discord
func respond(m:discord.Message)->string {
 if m.matches("hola") {return m.reply("Hola, "+m.author_id)}
 return ""
}
bot=discord.bot()
bot.on_message(&respond)
bot.start()`

const fixtureMessage = `{"op":0,"t":"MESSAGE_CREATE","s":4,"d":{"channel_id":"123","content":"hola","author":{"id":"567"}}}`

func TestLinuxDiscordLauncherAndNativeCallback(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("requires native x86-64")
	}
	path := filepath.Join(t.TempDir(), "bot")
	if err := os.WriteFile(path, compileProfile(hostedBot, "linux", "turbo"), 0700); err != nil {
		t.Fatal(err)
	}
	if !needsDiscordLauncher(path) {
		t.Fatal("compiled bot lacks its launcher marker")
	}

	t.Setenv("TRIPOL_DISCORD_PIPE", "")
	t.Setenv("DISCORD_TOKEN", "fixture-token")
	if output, err := exec.Command(path).CombinedOutput(); err == nil || !bytes.Contains(output, []byte("requires the TRIPOL launcher")) {
		t.Fatalf("missing launcher diagnostic: %q %v", output, err)
	}

	cmd := exec.Command(path)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	fixture := make(chan error, 1)
	if err := runDiscordPipes(cmd, func(ctx context.Context, requests, replies *os.File) error {
		err := botFixture(ctx, requests, replies)
		fixture <- err
		return err
	}, nil); err != nil {
		t.Fatalf("bot failed: %v\n%s\n%v", err, output.String(), <-fixture)
	}
	if err := <-fixture; err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Gateway READY") {
		t.Fatal(output.String())
	}
}

func botFixture(_ context.Context, input io.Reader, output io.Writer) error {
	sequence := []struct {
		op    uint32
		reply string
	}{
		{discordtransport.OpStart, "1"},
		{discordtransport.OpPoll, `{"op":0,"t":"READY","s":3,"d":{}}`},
		{discordtransport.OpPoll, fixtureMessage},
		{discordtransport.OpREST, `{"id":"890"}`},
		{discordtransport.OpPoll, "!"},
		{discordtransport.OpClose, "1"},
	}
	for _, step := range sequence {
		var hdr [16]byte
		if _, err := io.ReadFull(input, hdr[:]); err != nil {
			return err
		}
		if op := binary.LittleEndian.Uint32(hdr[:4]); op != step.op {
			return fmt.Errorf("expected operation %d, got %d", step.op, op)
		}
		var args [3]string
		for i := range args {
			n := binary.LittleEndian.Uint32(hdr[4+i*4:])
			if n > 65535 {
				return fmt.Errorf("invalid argument length %d", n)
			}
			b := make([]byte, n)
			if _, err := io.ReadFull(input, b); err != nil {
				return err
			}
			args[i] = string(b)
		}
		switch step.op {
		case discordtransport.OpStart:
			var req struct {
				Op   int `json:"op"`
				Data struct {
					Token   string `json:"token"`
					Intents int    `json:"intents"`
				} `json:"d"`
			}
			if json.Unmarshal([]byte(args[0]), &req) != nil || req.Op != 2 || req.Data.Token != "fixture-token" || req.Data.Intents != 33281 {
				return fmt.Errorf("invalid Identify from executable")
			}
		case discordtransport.OpREST:
			if args[0] != "/api/v10/channels/123/messages" || args[1] != "fixture-token" ||
				args[2] != `{"content":"Hola, 567","allowed_mentions":{"parse":[]}}` {
				return fmt.Errorf("invalid callback message: %q %q", args[0], args[2])
			}
		}
		var n [4]byte
		binary.LittleEndian.PutUint32(n[:], uint32(len(step.reply)))
		if _, err := output.Write(n[:]); err != nil {
			return err
		}
		if _, err := io.WriteString(output, step.reply); err != nil {
			return err
		}
	}
	return nil
}

func TestLinuxSlashSynchronizationAndInteractionCallback(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("requires native x86-64")
	}
	path := filepath.Join(t.TempDir(), "bot")
	source := strings.Replace(interactionBot,
		`func defer_reply(c:discord.Interaction)->string {return c.defer(true)}`,
		`func defer_reply(c:discord.Interaction)->string {
 if not c.respond(c.defer(true)) {return c.reply("Should not retry")}
 assert not c.respond(c.defer(true))
 assert c.edit_reply("Listo")
 followup=c.followup("Más",true)
 assert followup=="999"
 assert c.edit_followup(followup,"Actualizado")
 assert not c.delete_followup("../999")
 assert c.delete_followup(followup)
 return c.reply("Should not acknowledge twice")
}`, 1)
	if err := os.WriteFile(path, compileProfile(source+commandTree+contextMenuHandlers+eventRegistry+"b.start_with_events(registry)", "linux", "turbo"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DISCORD_TOKEN", "fixture-token")
	cmd := exec.Command(path)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	fixture := make(chan error, 1)
	if err := runDiscordPipes(cmd, func(_ context.Context, requests, replies *os.File) error {
		err := interactionFixture(requests, replies)
		fixture <- err
		return err
	}, nil); err != nil {
		t.Fatalf("bot failed: %v\n%s\n%v", err, output.String(), <-fixture)
	}
	if err := <-fixture; err != nil {
		t.Fatal(err)
	}
	if out := output.String(); !strings.Contains(out, "GUILD_MEMBER_ADD 9 321 18446744073709551614\n") ||
		!strings.Contains(out, "GUILD_ROLE_UPDATE 321 {\"id\":\"812\",\"name\":\"Maintainer\"}\n") {
		t.Fatal(out)
	}
}

func interactionFixture(input io.Reader, output io.Writer) error {
	deferred := strings.ReplaceAll(strings.ReplaceAll(slashEvent, `"name":"hola"`, `"name":"wait"`), `"id":"710"`, `"id":"715"`)
	expired := strings.Replace(deferred, `"id":"715"`, `"id":"716"`, 1)
	first := `[{"type":1,"name":"hola","description":"La descripci\u00f3n","id":"919","options":[{"required":true,"type":3,"description":"La opci\u00f3n","name":"text"}]},{"type":1,"name":"unmanaged","description":"Leave it alone","id":"920"},{"type":2,"name":"Perfil / ficha","default_member_permissions":"8192","id":"870"},{"type":1,"name":"inspect","description":"Inspecciona","id":"921","options":[{"type":2,"name":"stock","description":"Inventario","options":[{"type":1,"name":"find","description":"Busca","options":[{"type":3,"name":"query","description":"Qu\u00e9 buscar","autocomplete":true,"required":true},{"choices":[{"value":10,"name":"D\u00e9cada"}],"type":4,"name":"limit","description":"L\u00edmite"}]}]}]}]`
	steps := []bridgeStep{
		{discordtransport.OpREST, "/api/v10/oauth2/applications/@me", "fixture-token", "", `{"id":"890"}`},
		{discordtransport.OpREST, "/api/v10/applications/890/guilds/321/commands", "fixture-token", "", first},
		{discordtransport.OpREST, "/api/v10/applications/890/guilds/321/commands", "fixture-token", `"name":"wait"`, `{"id":"922"}`},
		{discordtransport.OpREST, "/api/v10/applications/890/guilds/321/commands", "fixture-token", `"type":3,"default_member_permissions":null`, `{"id":"871"}`},
		{discordtransport.OpStart, "", "", `"intents":3,`, "1"},
		{discordtransport.OpPoll, "", "", "", `{"op":0,"t":"READY","s":3,"d":{}}`},
		{discordtransport.OpPoll, "", "", "", slashEvent},
		{discordtransport.OpInteraction, "/api/v10/interactions/710/interaction-fixture/callback", "", `"content":"567 café 😀","flags":64`, "{}"},
		{discordtransport.OpPoll, "", "", "", buttonEvent},
		{discordtransport.OpInteraction, "/api/v10/interactions/711/button-fixture/callback", "", `"type":9`, "{}"},
		{discordtransport.OpPoll, "", "", "", modalEvent},
		// An empty reply models a rejected or expired interaction. The bot must
		// continue polling the same Gateway and avoid replaying the callback.
		{discordtransport.OpInteraction, "/api/v10/interactions/713/modal-fixture/callback", "", `"content":"Revisión"`, ""},
		{discordtransport.OpPoll, "", "", "", selectEvent},
		{discordtransport.OpInteraction, "/api/v10/interactions/712/select-fixture/callback", "", `"content":"101 102"`, "{}"},
		{discordtransport.OpPoll, "", "", "", autocompleteEvent},
		{discordtransport.OpInteraction, "/api/v10/interactions/714/autocomplete-fixture/callback", "", `"type":8,"data":{"choices":[{"name":"Café","value":"cafe"}]}`, "{}"},
		{discordtransport.OpPoll, "", "", "", deferred},
		{discordtransport.OpInteraction, "/api/v10/interactions/715/interaction-fixture/callback", "", `"type":5,"data":{"flags":64}`, "{}"},
		{discordtransport.OpInteractionEdit, "/api/v10/webhooks/890/interaction-fixture/messages/@original", "", `"content":"Listo"`, `{"id":"998"}`},
		{discordtransport.OpInteraction, "/api/v10/webhooks/890/interaction-fixture", "", `"content":"Más","flags":64`, `{"id":"999"}`},
		{discordtransport.OpInteractionEdit, "/api/v10/webhooks/890/interaction-fixture/messages/999", "", `"content":"Actualizado"`, `{"id":"999"}`},
		{discordtransport.OpInteractionEdit, "/api/v10/webhooks/890/interaction-fixture/messages/999", "", "", `{}`},
		{discordtransport.OpPoll, "", "", "", expired},
		{discordtransport.OpInteraction, "/api/v10/interactions/716/interaction-fixture/callback", "", `"type":5`, ""},
		{discordtransport.OpPoll, "", "", "", userMenuEvent},
		{discordtransport.OpInteraction, "/api/v10/interactions/741/menu-fixture/callback", "", `"content":"Perfil / ficha: 42"`, "{}"},
		{discordtransport.OpPoll, "", "", "", messageMenuEvent},
		{discordtransport.OpInteraction, "/api/v10/interactions/742/menu-fixture/callback", "", `Un mensaje`, "{}"},
		{discordtransport.OpPoll, "", "", "", memberJoined},
		{discordtransport.OpPoll, "", "", "", roleChanged},
		{discordtransport.OpPoll, "", "", "", "!"},
		{discordtransport.OpClose, "", "", "", "1"},
	}
	return serveBridgeFixture(input, output, steps)
}

type bridgeStep struct {
	op                          uint32
	path, auth, contains, reply string
}

func serveBridgeFixture(input io.Reader, output io.Writer, steps []bridgeStep) error {
	for index, step := range steps {
		var hdr [16]byte
		if _, err := io.ReadFull(input, hdr[:]); err != nil {
			return err
		}
		if op := binary.LittleEndian.Uint32(hdr[:4]); op != step.op {
			return fmt.Errorf("step %d: want operation %d, got %d", index, step.op, op)
		}
		var args [3]string
		for i := range args {
			n := binary.LittleEndian.Uint32(hdr[4+i*4:])
			if n > 65535 {
				return fmt.Errorf("invalid argument length")
			}
			b := make([]byte, n)
			if _, err := io.ReadFull(input, b); err != nil {
				return err
			}
			args[i] = string(b)
		}
		if step.op == discordtransport.OpStart {
			if !strings.Contains(args[0], step.contains) {
				return fmt.Errorf("incorrect Identify in step %d", index)
			}
		} else if args[0] != step.path || args[1] != step.auth || !strings.Contains(args[2], step.contains) ||
			(step.contains == "" && args[2] != "") {
			return fmt.Errorf("invalid operation arguments in step %d", index)
		}
		var size [4]byte
		binary.LittleEndian.PutUint32(size[:], uint32(len(step.reply)))
		if _, err := output.Write(size[:]); err != nil {
			return err
		}
		if _, err := io.WriteString(output, step.reply); err != nil {
			return err
		}
	}
	return nil
}

func TestLinuxDiscordModeration(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("requires native x86-64")
	}
	const source = `import tripol.discord
guild=discord.guild("321")
member=guild.member("567")
show member.read()
assert not member.kick()
assert member.ban(3600)
assert member.unban()
assert member.nickname("Niña")
assert member.nickname("")
assert member.timeout("2026-12-01T00:00:00Z")
assert member.timeout()
assert member.role("876")
assert member.role("876",false)
assert member.dm("Bienvenido")=="42"
room=discord.channel_ref("345")
assert room.read()=="{}"
assert room.slowmode(15)
assert room.pin("567")
assert room.pin("567",false)
assert room.react("567","🔥")
assert room.react("567","test:890",false)
assert room.purge(list(["567","568"]))
assert room.delete_message("569")
show "moderation OK"
`
	path := filepath.Join(t.TempDir(), "moderation")
	if err := os.WriteFile(path, compileProfile(source, "linux", "turbo"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DISCORD_TOKEN", "fixture-token")
	cmd := exec.Command(path)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	steps := []bridgeStep{
		{discordtransport.OpRESTGet, "/api/v10/guilds/321/members/567", "fixture-token", "", `{"user":{"id":"567"}}`},
		{discordtransport.OpRESTDelete, "/api/v10/guilds/321/members/567", "fixture-token", "", ""},
		{discordtransport.OpRESTPut, "/api/v10/guilds/321/bans/567", "fixture-token", `{"delete_message_seconds":3600}`, `{}`},
		{discordtransport.OpRESTDelete, "/api/v10/guilds/321/bans/567", "fixture-token", "", `{}`},
		{discordtransport.OpRESTPatch, "/api/v10/guilds/321/members/567", "fixture-token", `{"nick":"Niña"}`, `{}`},
		{discordtransport.OpRESTPatch, "/api/v10/guilds/321/members/567", "fixture-token", `{"nick":null}`, `{}`},
		{discordtransport.OpRESTPatch, "/api/v10/guilds/321/members/567", "fixture-token", `{"communication_disabled_until":"2026-12-01T00:00:00Z"}`, `{}`},
		{discordtransport.OpRESTPatch, "/api/v10/guilds/321/members/567", "fixture-token", `{"communication_disabled_until":null}`, `{}`},
		{discordtransport.OpRESTPut, "/api/v10/guilds/321/members/567/roles/876", "fixture-token", "", `{}`},
		{discordtransport.OpRESTDelete, "/api/v10/guilds/321/members/567/roles/876", "fixture-token", "", `{}`},
		{discordtransport.OpRESTPost, "/api/v10/users/@me/channels", "fixture-token", `{"recipient_id":"567"}`, `{"id":"555"}`},
		{discordtransport.OpRESTPost, "/api/v10/channels/555/messages", "fixture-token", `{"content":"Bienvenido","allowed_mentions":{"parse":[]}}`, `{"id":"42"}`},
		{discordtransport.OpRESTGet, "/api/v10/channels/345", "fixture-token", "", `{}`},
		{discordtransport.OpRESTPatch, "/api/v10/channels/345", "fixture-token", `{"rate_limit_per_user":15}`, `{}`},
		{discordtransport.OpRESTPut, "/api/v10/channels/345/messages/pins/567", "fixture-token", "", `{}`},
		{discordtransport.OpRESTDelete, "/api/v10/channels/345/messages/pins/567", "fixture-token", "", `{}`},
		{discordtransport.OpRESTPut, "/api/v10/channels/345/messages/567/reactions/%F0%9F%94%A5/@me", "fixture-token", "", `{}`},
		{discordtransport.OpRESTDelete, "/api/v10/channels/345/messages/567/reactions/test%3A890/@me", "fixture-token", "", `{}`},
		{discordtransport.OpRESTPost, "/api/v10/channels/345/messages/bulk-delete", "fixture-token", `{"messages":["567","568"]}`, `{}`},
		{discordtransport.OpRESTDelete, "/api/v10/channels/345/messages/569", "fixture-token", "", `{}`},
	}
	fixture := make(chan error, 1)
	err := runDiscordPipes(cmd, func(_ context.Context, requests, replies *os.File) error {
		err := serveBridgeFixture(requests, replies, steps)
		fixture <- err
		return err
	}, nil)
	if err != nil {
		t.Fatalf("moderation executable: %v (%s)", err, output.String())
	}
	if fixtureErr := <-fixture; fixtureErr != nil || output.String() != "{\"user\":{\"id\":\"567\"}}\nmoderation OK\n" {
		t.Fatalf("fixture %v, output %q", fixtureErr, output.String())
	}
}

func TestLinuxDiscordWelcomeTemplateFiltersJoinsAndRecovers(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("requires native x86-64")
	}
	dir := filepath.Join(t.TempDir(), "welcome")
	newProjectTemplate(dir, "discord-welcome")
	fs, ss, _ := loadProject(filepath.Join(dir, "main.tri"))
	tool := filepath.Join(dir, "welcome-bot")
	if err := os.WriteFile(tool, compileParsed(fs, ss, "linux", "turbo"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DISCORD_TOKEN", "fixture-token")
	t.Setenv("DISCORD_GUILD_ID", "321")
	t.Setenv("DISCORD_CHANNEL_ID", "345")
	cmd := exec.Command(tool)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	body := `{"content":"¡Bienvenido al servidor!","allowed_mentions":{"parse":[]}}`
	steps := []bridgeStep{
		{discordtransport.OpStart, "", "", `"intents":3,`, "1"},
		{discordtransport.OpPoll, "", "", "", `{"op":0,"t":"READY","s":3,"d":{}}`},
		{discordtransport.OpPoll, "", "", "", strings.Replace(memberJoined, `"guild_id":"321"`, `"guild_id":"999"`, 1)},
		{discordtransport.OpPoll, "", "", "", strings.Replace(memberJoined, `"user":{"id"`, `"user":{"bot":true,"id"`, 1)},
		{discordtransport.OpPoll, "", "", "", memberJoined},
		{discordtransport.OpRESTPost, "/api/v10/channels/345/messages", "fixture-token", body, ""},
		{discordtransport.OpPoll, "", "", "", strings.Replace(memberJoined, `"s":9`, `"s":10`, 1)},
		{discordtransport.OpRESTPost, "/api/v10/channels/345/messages", "fixture-token", body, `{"id":"999"}`},
		{discordtransport.OpPoll, "", "", "", "!"},
		{discordtransport.OpClose, "", "", "", "1"},
	}
	fixture := make(chan error, 1)
	if err := runDiscordPipes(cmd, func(_ context.Context, in, out *os.File) error {
		err := serveBridgeFixture(in, out, steps)
		fixture <- err
		return err
	}, nil); err != nil {
		t.Fatalf("welcome executable: %v (%s)", err, output.String())
	}
	if err := <-fixture; err != nil || strings.Count(output.String(), "No se pudo enviar el mensaje de bienvenida.") != 1 {
		t.Fatalf("fixture %v, output %q", err, output.String())
	}
}

func TestLinuxSimpleDiscordBotReceivesGatewayInteractions(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("requires native x86-64")
	}
	source, err := os.ReadFile(filepath.Join("examples", "discord_simple_components.tri"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "simple-bot")
	if err := os.WriteFile(path, compileProfile(string(source), "linux", "turbo"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DISCORD_TOKEN", "fixture-token")
	steps := []bridgeStep{
		{discordtransport.OpREST, "/api/v10/oauth2/applications/@me", "fixture-token", "", `{"id":"890"}`},
		{discordtransport.OpREST, "/api/v10/applications/890/commands", "fixture-token", "", `[]`},
		{discordtransport.OpREST, "/api/v10/applications/890/commands", "fixture-token", `"name":"info"`, `{"id":"922"}`},
		{discordtransport.OpStart, "", "", `"intents":1,`, "1"},
		{discordtransport.OpPoll, "", "", "", `{"op":0,"t":"READY","s":1,"d":{}}`},
		{discordtransport.OpPoll, "", "", "", `{"op":0,"t":"INTERACTION_CREATE","s":2,"d":{"type":2,"id":"710","application_id":"890","token":"simple-fixture","data":{"type":1,"name":"info"}}}`},
		{discordtransport.OpInteraction, "/api/v10/interactions/710/simple-fixture/callback", "", `"type":17,"accent_color":16777215`, `{}`},
		{discordtransport.OpPoll, "", "", "", `{"op":0,"t":"INTERACTION_CREATE","s":3,"d":{"type":3,"id":"711","application_id":"890","token":"simple-fixture","data":{"component_type":2,"custom_id":"accept"}}}`},
		{discordtransport.OpInteraction, "/api/v10/interactions/711/simple-fixture/callback", "", `"content":"Gracias!","flags":64`, `{}`},
		{discordtransport.OpPoll, "", "", "", "!"},
		{discordtransport.OpClose, "", "", "", "1"},
	}
	cmd := exec.Command(path)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	fixture := make(chan error, 1)
	if err := runDiscordPipes(cmd, func(_ context.Context, in, out *os.File) error {
		err := serveBridgeFixture(in, out, steps)
		fixture <- err
		return err
	}, nil); err != nil {
		t.Fatalf("bot: %v (%s)", err, output.String())
	}
	if err := <-fixture; err != nil {
		t.Fatal(err)
	}
}

func TestLinuxDiscordUploadsAndGuildManagement(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("requires native x86-64")
	}
	const source = `import tripol.discord
server=discord.guild("321")
assert server.channels()=="[]"
category=server.create_channel("Reports","category")
assert category=="450"
id=server.create_channel("daily","text",category,"Daily reports")
assert id=="451"
room=discord.channel_ref(id)
assert room.rename("status")
assert room.topic("Updates")
assert room.move()
panel=discord.container()
assert panel.file("report.txt",true)
assert panel.send_files(id,list(["work/report.txt"]))=="901"
embed=discord.embed()
embed.image(discord.attachment_url("photo.png"))
assert embed.send_files(id,list(["photo.png"]))=="902"
assert room.edit_message("333","Actualizado")
row=discord.button_row()
row.button("Confirm","confirm","success",false,"✅")
row.disable("confirm")
other=discord.container()
other.row(row)
assert room.edit_components("901",other)
assert server.audit(5,"900",10)=="{}"
assert room.delete()
assert server.roles()=="[]"
role_id=server.create_role("Reviewer","0","#FFFFFF")
assert role_id=="952"
role=server.role_ref(role_id)
assert role.read()=="{}"
assert role.permissions()=="8192"
assert role.edit("Review","8192","#000000")
assert server.move_role(role_id,2)
try {server.move_role("321",1)} catch code {assert code==12}
assert role.delete()
`
	program := filepath.Join(t.TempDir(), "guild")
	if err := os.WriteFile(program, compileProfile(source, "linux", "turbo"), 0700); err != nil {
		t.Fatal(err)
	}
	compileProfile(source, "windows", "turbo")
	t.Setenv("DISCORD_TOKEN", "fixture-token")
	steps := []bridgeStep{
		{discordtransport.OpRESTGet, "/api/v10/guilds/321/channels", "fixture-token", "", `[]`},
		{discordtransport.OpRESTPost, "/api/v10/guilds/321/channels", "fixture-token", `{"name":"Reports","type":4}`, `{"id":"450"}`},
		{discordtransport.OpRESTPost, "/api/v10/guilds/321/channels", "fixture-token", `{"name":"daily","type":0,"parent_id":"450","topic":"Daily reports"}`, `{"id":"451"}`},
		{discordtransport.OpRESTPatch, "/api/v10/channels/451", "fixture-token", `{"name":"status"}`, `{}`},
		{discordtransport.OpRESTPatch, "/api/v10/channels/451", "fixture-token", `{"topic":"Updates"}`, `{}`},
		{discordtransport.OpRESTPatch, "/api/v10/channels/451", "fixture-token", `{"parent_id":null}`, `{}`},
		{discordtransport.OpRESTUpload, "/api/v10/channels/451/messages", "fixture-token", `"type":13,"file":{"url":"attachment://report.txt"},"spoiler":true`, `{"id":"901"}`},
		{discordtransport.OpRESTUpload, "/api/v10/channels/451/messages", "fixture-token", `"image":{"url":"attachment://photo.png"}`, `{"id":"902"}`},
		{discordtransport.OpRESTPatch, "/api/v10/channels/451/messages/333", "fixture-token", `{"content":"Actualizado","allowed_mentions":{"parse":[]}}`, `{}`},
		{discordtransport.OpRESTPatch, "/api/v10/channels/451/messages/901", "fixture-token", `"emoji":{"name":"✅"},"disabled":true`, `{}`},
		{discordtransport.OpRESTGet, "/api/v10/guilds/321/audit-logs?limit=5&before=900&action_type=10", "fixture-token", "", `{}`},
		{discordtransport.OpRESTDelete, "/api/v10/channels/451", "fixture-token", "", `{}`},
		{discordtransport.OpRESTGet, "/api/v10/guilds/321/roles", "fixture-token", "", `[]`},
		{discordtransport.OpRESTPost, "/api/v10/guilds/321/roles", "fixture-token", `{"name":"Reviewer","permissions":"0","colors":{"primary_color":16777215}}`, `{"id":"952"}`},
		{discordtransport.OpRESTGet, "/api/v10/guilds/321/roles/952", "fixture-token", "", `{}`},
		{discordtransport.OpRESTGet, "/api/v10/guilds/321/permissions/role/952", "fixture-token", "", `{"permissions":"8192"}`},
		{discordtransport.OpRESTPatch, "/api/v10/guilds/321/roles/952", "fixture-token", `{"name":"Review","permissions":"8192","colors":{"primary_color":0}}`, `{}`},
		{discordtransport.OpRESTPatch, "/api/v10/guilds/321/roles", "fixture-token", `[{"id":"952","position":2}]`, `[]`},
		{discordtransport.OpRESTDelete, "/api/v10/guilds/321/roles/952", "fixture-token", "", `{}`},
	}
	cmd := exec.Command(program)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	fixture := make(chan error, 1)
	err := runDiscordPipes(cmd, func(_ context.Context, input, output *os.File) error {
		err := serveBridgeFixture(input, output, steps)
		fixture <- err
		return err
	}, nil)
	if err != nil {
		t.Fatalf("guild executable: %v (%s)", err, output.String())
	}
	if err := <-fixture; err != nil {
		t.Fatal(err)
	}
}

func TestLinuxDiscordModerationTemplateEnforcesOperator(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("requires native x86-64")
	}
	dir := filepath.Join(t.TempDir(), "moderator")
	newProjectTemplate(dir, "discord-moderation")
	fs, ss, _ := loadProject(filepath.Join(dir, "main.tri"))
	program := filepath.Join(dir, "moderator")
	if err := os.WriteFile(program, compileParsed(fs, ss, "linux", "turbo"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DISCORD_TOKEN", "fixture-token")
	t.Setenv("DISCORD_GUILD_ID", "321")
	t.Setenv("DISCORD_OPERATOR_ID", "22")
	event := func(id, name, operator, permissions string) string {
		return fmt.Sprintf(`{"op":0,"t":"INTERACTION_CREATE","d":{"id":%q,"token":"moderation-fixture","type":2,"application_id":"890","guild_id":"321","member":{"user":{"id":%q},"permissions":%q},"data":{"type":1,"name":%q,"guild_id":"321","options":[{"name":"miembro","type":6,"value":"45"},{"name":"motivo","type":3,"value":"Prueba de advertencia"}]}}}`, id, operator, permissions, name)
	}
	paths := "/api/v10/applications/890/guilds/321/commands"
	steps := []bridgeStep{
		{discordtransport.OpREST, "/api/v10/oauth2/applications/@me", "fixture-token", "", `{"id":"890"}`},
		{discordtransport.OpREST, paths, "fixture-token", "", `[]`},
		{discordtransport.OpREST, paths, "fixture-token", `"default_member_permissions":"2"`, `{"id":"911"}`},
		{discordtransport.OpREST, paths, "fixture-token", `"default_member_permissions":"4"`, `{"id":"912"}`},
		{discordtransport.OpREST, paths, "fixture-token", `"name":"unban"`, `{"id":"913"}`},
		{discordtransport.OpREST, paths, "fixture-token", `"name":"warn"`, `{"id":"914"}`},
		{discordtransport.OpREST, paths, "fixture-token", `"name":"timeout"`, `{"id":"915"}`},
		{discordtransport.OpStart, "", "", `"intents":1`, "1"},
		{discordtransport.OpPoll, "", "", "", `{"op":0,"t":"READY","s":1,"d":{}}`},
		{discordtransport.OpPoll, "", "", "", event("50", "kick", "99", "8")},
		{discordtransport.OpInteraction, "/api/v10/interactions/50/moderation-fixture/callback", "", `"content":"Operación no autorizada.","flags":64`, `{}`},
		{discordtransport.OpPoll, "", "", "", event("51", "kick", "22", "2")},
		{discordtransport.OpInteraction, "/api/v10/interactions/51/moderation-fixture/callback", "", `"type":5,"data":{"flags":64}`, `{}`},
		{discordtransport.OpRESTDelete, "/api/v10/guilds/321/members/45", "fixture-token", "", `{}`},
		{discordtransport.OpInteractionEdit, "/api/v10/webhooks/890/moderation-fixture/messages/@original", "", `"content":"Operación aplicada."`, `{}`},
		{discordtransport.OpPoll, "", "", "", event("52", "ban", "22", "4")},
		{discordtransport.OpInteraction, "/api/v10/interactions/52/moderation-fixture/callback", "", `"type":5,"data":{"flags":64}`, `{}`},
		{discordtransport.OpRESTPut, "/api/v10/guilds/321/bans/45", "fixture-token", `"delete_message_seconds":0`, ""},
		{discordtransport.OpInteractionEdit, "/api/v10/webhooks/890/moderation-fixture/messages/@original", "", `"content":"Discord no pudo aplicar la operación."`, `{}`},
		{discordtransport.OpPoll, "", "", "", event("53", "warn", "22", "4")},
		{discordtransport.OpInteraction, "/api/v10/interactions/53/moderation-fixture/callback", "", `"content":"Te falta el permiso requerido.","flags":64`, `{}`},
		{discordtransport.OpPoll, "", "", "", event("54", "warn", "22", "1099511627776")},
		{discordtransport.OpInteraction, "/api/v10/interactions/54/moderation-fixture/callback", "", `"type":5,"data":{"flags":64}`, `{}`},
		{discordtransport.OpInteractionEdit, "/api/v10/webhooks/890/moderation-fixture/messages/@original", "", `"content":"Operación aplicada."`, `{}`},
		{discordtransport.OpPoll, "", "", "", "!"},
		{discordtransport.OpClose, "", "", "", "1"},
	}
	cmd := exec.Command(program)
	cmd.Dir = dir
	t.Setenv("DISCORD_WARNINGS_PATH", "")
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	fixture := make(chan error, 1)
	if err := runDiscordPipes(cmd, func(_ context.Context, in, out *os.File) error {
		err := serveBridgeFixture(in, out, steps)
		fixture <- err
		return err
	}, nil); err != nil {
		t.Fatalf("moderator executable: %v (%s)", err, output.String())
	}
	if err := <-fixture; err != nil {
		t.Fatal(err)
	}
	warnings, err := os.ReadFile(filepath.Join(dir, "discord-warnings.jsonl"))
	const want = `{"interaction":"54","guild":"321","moderator":"22","target":"45","reason":"Prueba de advertencia"}` + "\n"
	if err != nil || string(warnings) != want {
		t.Fatalf("warnings %q, error %v", warnings, err)
	}
}
