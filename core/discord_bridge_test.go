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
