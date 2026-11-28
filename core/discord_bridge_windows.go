package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"tripol/internal/discordtransport"
)

func runNativeCommand(cmd *exec.Cmd, studio *discordStudioOptions) error {
	if studio != nil {
		return errors.New("Discord Studio is available on Linux x86-64")
	}
	if !needsDiscordLauncher(cmd.Path) {
		return cmd.Run()
	}
	client := discordtransport.New()
	defer client.Close()
	return runDiscordPipes(cmd, func(ctx context.Context, requests, replies *os.File) error {
		return discordtransport.Serve(ctx, client, requests, replies, os.Stderr)
	}, nil)
}

func inheritDiscordPipes(cmd *exec.Cmd, write, read *os.File) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	// os.Pipe creates inheritable handles. The explicit process handle list
	// restricts inheritance to these child ends and the standard streams.
	cmd.SysProcAttr.AdditionalInheritedHandles = append(cmd.SysProcAttr.AdditionalInheritedHandles, syscall.Handle(write.Fd()), syscall.Handle(read.Fd()))
	cmd.Env = append(cmd.Environ(),
		fmt.Sprintf("TRIPOL_DISCORD_WRITE=%016x", write.Fd()),
		fmt.Sprintf("TRIPOL_DISCORD_READ=%016x", read.Fd()))
}

func packageDiscordExecutable(binary []byte, _ string) []byte { return binary }
func runBundledDiscord() bool                                 { return false }

func runDiscordExecutable(args []string) {
	if len(args) != 1 {
		die("usage: tripol discord-run FILE")
	}
	path, err := filepath.Abs(args[0])
	if err != nil {
		die("%v", err)
	}
	if !needsDiscordLauncher(path) {
		die("expected a generated Discord PE with transport enabled")
	}
	runNative(path, nil)
}
