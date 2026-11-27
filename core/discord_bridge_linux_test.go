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
