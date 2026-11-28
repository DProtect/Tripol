//go:build !linux && !windows

package core

import (
	"errors"
	"os/exec"
)

func runNativeCommand(cmd *exec.Cmd, studio *discordStudioOptions) error {
	if studio != nil {
		return errors.New("Discord Studio is available on Linux x86-64")
	}
	return cmd.Run()
}

func packageDiscordExecutable(binary []byte, _ string) []byte { return binary }
func runBundledDiscord() bool                                 { return false }
func runDiscordExecutable(_ []string)                         { die("tripol discord-run is available on Linux and Windows") }
