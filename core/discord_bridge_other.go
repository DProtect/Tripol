//go:build !linux

package core

import "os/exec"

func runNativeCommand(cmd *exec.Cmd) error { return cmd.Run() }

func packageDiscordExecutable(binary []byte, _ string) []byte { return binary }
func runBundledDiscord() bool                                 { return false }
func runDiscordExecutable(_ []string)                         { die("tripol discord-run is available on Linux") }
