//go:build linux || windows

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
	"time"

	"tripol/internal/discordtransport"
)

func TestDiscordNativeBridgeUploadsInteraction(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("requires native x86-64")
	}
	const source = `import tripol.discord
func send_file(ctx) {
 paths=string_list()
 paths.add("report.txt")
 panel=discord.components()
 panel.file("report.txt")
 show ctx.reply_files(panel,paths,true)
 show ctx.reply_files("duplicate",paths,true)
 ""
}
bot=discord.bot()
bot.on_button("café",send_file)
bot.start()
`
	path := filepath.Join(t.TempDir(), "bot.exe")
	if err := os.WriteFile(path, compileProfile(source, runtime.GOOS, "turbo"), 0700); err != nil {
		t.Fatal(err)
	}
	if !needsDiscordLauncher(path) {
		t.Fatal("generated program is missing the transport marker")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Env = append(cmd.Environ(), "DISCORD_TOKEN=local-test-token")
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	var fixtureErr error
	err := runDiscordPipes(cmd, func(_ context.Context, requests, replies *os.File) error {
		fixtureErr = interactionUploadFixture(requests, replies)
		return fixtureErr
	}, nil)
	if err != nil || fixtureErr != nil || !strings.Contains(strings.ReplaceAll(output.String(), "\r\n", "\n"), "true\nfalse\n") {
		t.Fatalf("exit %v, bridge %v, output %q", err, fixtureErr, output.String())
	}
}

func interactionUploadFixture(input io.Reader, output io.Writer) error {
	steps := []struct {
		op    uint32
		reply string
	}{
		{discordtransport.OpStart, "1"},
		{discordtransport.OpPoll, buttonEvent},
		{discordtransport.OpInteractionUpload, "{}"},
		{discordtransport.OpPoll, "!"},
		{discordtransport.OpClose, "1"},
	}
	for _, step := range steps {
		var header [16]byte
		if _, err := io.ReadFull(input, header[:]); err != nil {
			return err
		}
		if op := binary.LittleEndian.Uint32(header[:]); op != step.op {
			return fmt.Errorf("expected operation %d, got %d", step.op, op)
		}
		var args [3]string
		for i := range args {
			n := binary.LittleEndian.Uint32(header[4+i*4:])
			if n > 65535 {
				return fmt.Errorf("invalid argument length %d", n)
			}
			value := make([]byte, n)
			if _, err := io.ReadFull(input, value); err != nil {
				return err
			}
			args[i] = string(value)
		}
		if step.op == discordtransport.OpInteractionUpload {
			var envelope struct {
				Payload struct {
					Type int `json:"type"`
					Data struct {
						Flags      int             `json:"flags"`
						Components json.RawMessage `json:"components"`
					} `json:"data"`
				} `json:"payload"`
				Files []string `json:"files"`
			}
			if args[0] != "/api/v10/interactions/711/button-fixture/callback" || args[1] != "" ||
				json.Unmarshal([]byte(args[2]), &envelope) != nil || envelope.Payload.Type != 4 ||
				envelope.Payload.Data.Flags != 32832 || !bytes.Contains(envelope.Payload.Data.Components, []byte("attachment://report.txt")) ||
				len(envelope.Files) != 1 || envelope.Files[0] != "report.txt" {
				return fmt.Errorf("unexpected interaction upload: %+v", args)
			}
		}
		var length [4]byte
		binary.LittleEndian.PutUint32(length[:], uint32(len(step.reply)))
		if _, err := output.Write(length[:]); err != nil {
			return err
		}
		if _, err := io.WriteString(output, step.reply); err != nil {
			return err
		}
	}
	return nil
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

func TestDiscordNativeCommandPolicy(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("requires native x86-64")
	}
	const source = `import tripol.discord
func work(ctx) {
 show ctx.cooldown_remaining()
 if ctx.id=="74" {show ctx.reset_cooldown()}
 ctx.reply("Hecho",true)
}
b=discord.bot()
opts=discord.slash_options()
opts.add("miembro","Miembro","user",true)
b.slash_with_options("work","Trabajo",work,opts,"321")
b.permission("work","KICK_MEMBERS","","321")
b.cooldown("work",60,"channel","En espera","321")
b.persist_cooldowns("cooldowns.json")
b.start()
`
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.exe")
	if err := os.WriteFile(path, compileProfile(source, runtime.GOOS, "turbo"), 0700); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Env = append(cmd.Environ(), "DISCORD_TOKEN=fixture-token")
	cmd.Dir = dir
	event := func(id, permissions, appPermissions string) string {
		return fmt.Sprintf(`{"op":0,"t":"INTERACTION_CREATE","d":{"id":%q,"token":"callback","type":2,"application_id":"890","guild_id":"321","channel_id":"444","app_permissions":%q,"member":{"permissions":%q,"user":{"id":"555"}},"data":{"type":1,"guild_id":"321","name":"work","options":[{"name":"miembro","type":6,"value":"999"}]}}}`, id, appPermissions, permissions)
	}
	callback := func(id string) string { return "/api/v10/interactions/" + id + "/callback/callback" }
	commands := "/api/v10/applications/890/guilds/321/commands"
	const cooldown = `{"command":"swork/321","scope":"channel","subject":"444","seconds":60}`
	const moderation = `{"guild":"321","actor":"555","target":"999","permission":1}`
	steps := []bridgeStep{
		{discordtransport.OpCooldown, "load", "", `"cooldowns.json"`, "0"},
		{discordtransport.OpREST, "/api/v10/oauth2/applications/@me", "fixture-token", "", `{"id":"890"}`},
		{discordtransport.OpREST, commands, "fixture-token", "", `[]`},
		{discordtransport.OpREST, commands, "fixture-token", `"default_member_permissions":"2"`, `{"id":"31"}`},
		{discordtransport.OpStart, "", "", `"intents":1`, "1"},
		{discordtransport.OpPoll, "", "", "", event("70", "4", "2")},
		{discordtransport.OpInteraction, callback("70"), "", "Te falta el permiso", "{}"},
		{discordtransport.OpPoll, "", "", "", event("71", "2", "0")},
		{discordtransport.OpInteraction, callback("71"), "", "El bot necesita el permiso", "{}"},
		{discordtransport.OpPoll, "", "", "", event("72", "2", "2")},
		{discordtransport.OpAuthorizeModeration, "", "fixture-token", moderation, "owner"},
		{discordtransport.OpInteraction, callback("72"), "", "propietario", "{}"},
		{discordtransport.OpPoll, "", "", "", event("73", "2", "2")},
		{discordtransport.OpAuthorizeModeration, "", "fixture-token", moderation, "1"},
		{discordtransport.OpCooldown, "use", "", cooldown, "6"},
		{discordtransport.OpInteraction, callback("73"), "", "En espera 6 s.", "{}"},
		{discordtransport.OpPoll, "", "", "", event("74", "2", "2")},
		{discordtransport.OpAuthorizeModeration, "", "fixture-token", moderation, "1"},
		{discordtransport.OpCooldown, "use", "", cooldown, "0"},
		{discordtransport.OpCooldown, "remaining", "", cooldown, "60"},
		{discordtransport.OpCooldown, "reset", "", cooldown, "0"},
		{discordtransport.OpInteraction, callback("74"), "", "Hecho", "{}"},
		{discordtransport.OpPoll, "", "", "", "!"},
		{discordtransport.OpClose, "", "", "", "1"},
	}
	var stdout bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stdout
	var fixtureErr error
	err := runDiscordPipes(cmd, func(_ context.Context, input, output *os.File) error {
		fixtureErr = serveBridgeFixture(input, output, steps)
		return fixtureErr
	}, nil)
	if err != nil || fixtureErr != nil || strings.Count(stdout.String(), "60\ntrue\n") != 1 {
		t.Fatalf("process error %v, bridge error %v, output %q", err, fixtureErr, stdout.String())
	}
}

func TestDiscordNativeAdministrativeOperations(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("requires native x86-64")
	}
	const source = `import tripol.discord
g=discord.guild("321")
room=discord.channel_ref("123")
assert room.lock()
assert room.unlock()
assert room.overwrite("678","3072","0","user")
assert room.clear_overwrite("678")
assert g.private_channel("support","678","99")=="124"
assert g.move_channel("124",2)
assert room.invite(600,1)=="abc"
assert len(room.invites())>0
assert len(g.invites())>0
assert len(discord.invite_info("abc",true))>0
assert discord.revoke_invite("abc")
assert room.thread("Planning","hour",true)=="235"
assert room.thread_from("200","Follow up")=="236"
thread=discord.channel_ref("235")
assert thread.thread_member("678")
assert thread.thread_member("678",false)
assert thread.archive()
assert thread.lock_thread()
assert g.create_channel("ideas","forum")=="300"
forum=discord.channel_ref("300")
labels=discord.forum_tags()
labels.tag("Idea","100",true,"✅")
assert forum.set_forum_tags(labels)
assert forum.post("Build","First post",list(["100"]))=="301"
post=discord.channel_ref("301")
assert post.apply_tags(list(["100"]))
hook=room.create_webhook("News")
assert hook.id()=="901"
info=hook.read()
assert info.name=="News"
assert hook.edit("Updates")
assert hook.send("Available")=="960"
card=discord.embed()
card.image(discord.attachment_url("report.txt"))
assert hook.send_files(card,list(["report.txt"]))=="961"
assert hook.delete()
show "administration OK"
`
	path := filepath.Join(t.TempDir(), "admin.exe")
	if err := os.WriteFile(path, compileProfile(source, runtime.GOOS, "turbo"), 0700); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Env = append(cmd.Environ(), "DISCORD_TOKEN=fixture-token")
	const key = "fixture-token"
	steps := []bridgeStep{
		{discordtransport.OpChannelLock, "123", key, "lock", "1"},
		{discordtransport.OpChannelLock, "123", key, "unlock", "1"},
		{discordtransport.OpRESTPut, "/api/v10/channels/123/permissions/678", key, `"type":1,"allow":"3072"`, "{}"},
		{discordtransport.OpRESTDelete, "/api/v10/channels/123/permissions/678", key, "", "{}"},
		{discordtransport.OpRESTPost, "/api/v10/guilds/321/channels", key, `"id":"321","type":0,"deny":"1024"`, `{"id":"124"}`},
		{discordtransport.OpRESTPatch, "/api/v10/guilds/321/channels", key, `[{"id":"124","position":2}]`, "[]"},
		{discordtransport.OpRESTPost, "/api/v10/channels/123/invites", key, `"max_age":600,"max_uses":1`, `{"code":"abc"}`},
		{discordtransport.OpRESTGet, "/api/v10/channels/123/invites", key, "", "[]"},
		{discordtransport.OpRESTGet, "/api/v10/guilds/321/invites", key, "", "[]"},
		{discordtransport.OpRESTGet, "/api/v10/invites/abc?with_counts=true&with_expiration=true", key, "", "{}"},
		{discordtransport.OpRESTDelete, "/api/v10/invites/abc", key, "", "{}"},
		{discordtransport.OpRESTPost, "/api/v10/channels/123/threads", key, `"type":12,"auto_archive_duration":60`, `{"id":"235"}`},
		{discordtransport.OpRESTPost, "/api/v10/channels/123/messages/200/threads", key, `"auto_archive_duration":1440`, `{"id":"236"}`},
		{discordtransport.OpRESTPut, "/api/v10/channels/235/thread-members/678", key, "", "{}"},
		{discordtransport.OpRESTDelete, "/api/v10/channels/235/thread-members/678", key, "", "{}"},
		{discordtransport.OpRESTPatch, "/api/v10/channels/235", key, `{"archived":true}`, "{}"},
		{discordtransport.OpRESTPatch, "/api/v10/channels/235", key, `{"locked":true}`, "{}"},
		{discordtransport.OpRESTPost, "/api/v10/guilds/321/channels", key, `"type":15`, `{"id":"300"}`},
		{discordtransport.OpRESTPatch, "/api/v10/channels/300", key, `"emoji_name":"✅"`, "{}"},
		{discordtransport.OpRESTPost, "/api/v10/channels/300/threads", key, `"applied_tags":["100"]`, `{"id":"301"}`},
		{discordtransport.OpRESTPatch, "/api/v10/channels/301", key, `{"applied_tags":["100"]}`, "{}"},
		{discordtransport.OpRESTPost, "/api/v10/channels/123/webhooks", key, `"name":"News"`, `{"id":"901","token":"local-hook-fixture"}`},
		{discordtransport.OpRESTGet, "/api/v10/webhooks/901", key, "", `{"id":"901","name":"News","channel_id":"123","guild_id":"321","token":"local-hook-fixture"}`},
		{discordtransport.OpRESTPatch, "/api/v10/webhooks/901", key, `"name":"Updates"`, "{}"},
		{discordtransport.OpWebhook, "901", "local-hook-fixture", `"content":"Available"`, `{"id":"960"}`},
		{discordtransport.OpWebhookUpload, "901", "local-hook-fixture", `"files":["report.txt"]`, `{"id":"961"}`},
		{discordtransport.OpRESTDelete, "/api/v10/webhooks/901", key, "", "{}"},
	}
	var stdout bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stdout
	var fixtureErr error
	err := runDiscordPipes(cmd, func(_ context.Context, input, output *os.File) error {
		fixtureErr = serveBridgeFixture(input, output, steps)
		return fixtureErr
	}, nil)
	if err != nil || fixtureErr != nil || !strings.Contains(stdout.String(), "administration OK") ||
		strings.Contains(stdout.String(), "local-hook-fixture") {
		t.Fatalf("process error %v, bridge error %v, output %q", err, fixtureErr, stdout.String())
	}
}
