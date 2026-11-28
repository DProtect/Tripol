//go:build linux || windows

package core

import (
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

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

func runDiscordPipes(cmd *exec.Cmd, serve func(context.Context, *os.File, *os.File) error, monitor func(context.Context, int)) error {
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
	inheritDiscordPipes(cmd, childWrite, childRead)
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
	var monitored chan struct{}
	var stopMonitor context.CancelFunc
	if monitor != nil {
		monitored = make(chan struct{})
		monitorCtx, cancelMonitor := context.WithCancel(context.Background())
		stopMonitor = cancelMonitor
		go func() { defer close(monitored); monitor(monitorCtx, cmd.Process.Pid) }()
	}
	err = cmd.Wait()
	close(finished)
	cancel()
	requests.Close()
	replies.Close()
	<-done
	if monitored != nil {
		stopMonitor()
		<-monitored
	}
	return err
}
