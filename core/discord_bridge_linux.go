package core

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/binary"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"tripol/internal/discordtransport"
)

func runNativeCommand(cmd *exec.Cmd) error {
	if !needsDiscordLauncher(cmd.Path) {
		return cmd.Run()
	}
	return runDiscordCommand(cmd, discordtransport.New())
}

func needsDiscordLauncher(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	marker := discordBridgeMarker
	if _, err := file.Seek(-int64(len(marker)), io.SeekEnd); err != nil {
		return false
	}
	data := make([]byte, len(marker))
	_, err = io.ReadFull(file, data)
	return err == nil && string(data) == marker
}

func runDiscordCommand(cmd *exec.Cmd, client *discordtransport.Client) error {
	defer client.Close()
	return runDiscordPipes(cmd, func(ctx context.Context, requests, replies *os.File) error {
		return discordtransport.Serve(ctx, client, requests, replies, os.Stderr)
	})
}

func runDiscordPipes(cmd *exec.Cmd, serve func(context.Context, *os.File, *os.File) error) error {
	requests, childWrite, err := os.Pipe()
	if err != nil {
		return err
	}
	defer requests.Close()
	defer childWrite.Close()
	childRead, replies, err := os.Pipe()
	if err != nil {
		return err
	}
	defer childRead.Close()
	defer replies.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	cmd.ExtraFiles = []*os.File{childWrite, childRead}
	cmd.Env = append(cmd.Environ(), "TRIPOL_DISCORD_PIPE=1")
	// The bridge and cmd.Wait can write diagnostics concurrently; keep the
	// bridge on the process stderr rather than the command's diagnostic buffer.
	if err := cmd.Start(); err != nil {
		return err
	}
	childWrite.Close()
	childRead.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = serve(ctx, requests, replies)
		replies.Close()
	}()
	finished := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = cmd.Process.Kill()
		case <-finished:
		}
	}()
	err = cmd.Wait()
	close(finished)
	cancel()
	requests.Close()
	replies.Close()
	<-done
	return err
}

//go:embed discord-websocket-LICENSE.txt
var discordWebSocketLicense string

const bundleMarker = "\x00TRIPOL-DISCORD-ELF-V1\x00"
const maxBundleProgram = 32 << 20

func packageDiscordExecutable(program []byte, target string) []byte {
	if target != "linux" || runtime.GOARCH != "amd64" || !bytes.HasSuffix(program, []byte(discordBridgeMarker)) {
		return program
	}
	if len(program) > maxBundleProgram {
		die("Discord ELF exceeds the launcher packaging limit")
	}
	source, err := os.Executable()
	if err != nil {
		die("%v", err)
	}
	launcher, err := os.ReadFile(source)
	if err != nil {
		die("%v", err)
	}
	launcher = append(launcher, discordWebSocketLicense...)
	launcher = append(launcher, program...)
	launcher = binary.LittleEndian.AppendUint64(launcher, uint64(len(program)))
	return append(launcher, bundleMarker...)
}

func runBundledDiscord() bool {
	source, err := os.Open("/proc/self/exe")
	if err != nil {
		return false
	}
	defer source.Close()
	footer := make([]byte, 8+len(bundleMarker))
	if _, err := source.Seek(-int64(len(footer)), io.SeekEnd); err != nil {
		return false
	}
	if _, err := io.ReadFull(source, footer); err != nil || string(footer[8:]) != bundleMarker {
		return false
	}
	length := binary.LittleEndian.Uint64(footer[:8])
	if length == 0 || length > maxBundleProgram {
		die("invalid bundled Discord ELF length")
	}
	if _, err := source.Seek(-int64(len(footer))-int64(length), io.SeekEnd); err != nil {
		die("invalid bundled Discord ELF")
	}
	program := make([]byte, int(length))
	if _, err := io.ReadFull(source, program); err != nil || !bytes.HasPrefix(program, []byte("\x7fELF")) || !bytes.HasSuffix(program, []byte(discordBridgeMarker)) {
		die("invalid bundled Discord ELF")
	}
	dir, err := os.MkdirTemp("", "tripol-discord-")
	if err != nil {
		die("%v", err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "bot")
	write(path, program)
	runNative(path)
	return true
}

func runDiscordExecutable(args []string) {
	if len(args) != 1 {
		die("usage: tripol discord-run FILE")
	}
	path, err := filepath.Abs(args[0])
	if err != nil {
		die("%v", err)
	}
	if !needsDiscordLauncher(path) {
		die("expected a generated Discord ELF with transport enabled")
	}
	runNative(path)
}
