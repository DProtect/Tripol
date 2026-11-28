package core

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"tripol/internal/discordtransport"
)

func runNativeCommand(cmd *exec.Cmd, studio *discordStudioOptions) error {
	if !needsDiscordLauncher(cmd.Path) {
		if studio != nil {
			return errors.New("Discord Studio requires a program that uses Discord networking")
		}
		return cmd.Run()
	}
	return runDiscordCommand(cmd, discordtransport.New(), studio)
}

func runDiscordCommand(cmd *exec.Cmd, client *discordtransport.Client, studio *discordStudioOptions) error {
	defer client.Close()
	var monitor func(context.Context, int)
	if studio != nil {
		client.EnableDispatchBreakdown()
		output := os.Stderr
		if studio.metricsFile != "" {
			file, err := os.OpenFile(studio.metricsFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
			if err != nil {
				return err
			}
			defer file.Close()
			output = file
		}
		monitor = func(ctx context.Context, pid int) { monitorDiscord(ctx, pid, client, output) }
	}
	return runDiscordPipes(cmd, func(ctx context.Context, requests, replies *os.File) error {
		return discordtransport.Serve(ctx, client, requests, replies, os.Stderr)
	}, monitor)
}

func inheritDiscordPipes(cmd *exec.Cmd, write, read *os.File) {
	cmd.ExtraFiles = []*os.File{write, read}
	cmd.Env = append(cmd.Environ(), "TRIPOL_DISCORD_PIPE=1")
}

//go:embed discord-websocket-LICENSE.txt
var discordWebSocketLicense string

//go:embed http-dependencies-LICENSE.txt
var httpDependenciesLicense string

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
	launcher = append(launcher, httpDependenciesLicense...)
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
	runNative(path, nil)
	return true
}

func runDiscordExecutable(args []string) {
	if len(args) != 1 {
		die("usage: tripol net-run FILE")
	}
	path, err := filepath.Abs(args[0])
	if err != nil {
		die("%v", err)
	}
	if !needsDiscordLauncher(path) {
		die("expected a generated ELF with network transport enabled")
	}
	runNative(path, nil)
}
