package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"tripol/internal/discordtransport"
)

type discordProcessMetrics struct {
	CPUTicks uint64 `json:"cpu_ticks"`
	RSSBytes uint64 `json:"rss_bytes"`
}

func discordProcess(pid int) *discordProcessMetrics {
	prefix := "/proc/" + strconv.Itoa(pid) + "/"
	stat, err := os.ReadFile(prefix + "stat")
	if err != nil {
		return nil
	}
	line := string(stat)
	end := strings.LastIndexByte(line, ')')
	if end < 0 {
		return nil
	}
	fields := strings.Fields(line[end+1:])
	if len(fields) < 13 {
		return nil
	}
	user, err1 := strconv.ParseUint(fields[11], 10, 64)
	system, err2 := strconv.ParseUint(fields[12], 10, 64)
	if err1 != nil || err2 != nil {
		return nil
	}
	status, err := os.ReadFile(prefix + "status")
	if err != nil {
		return nil
	}
	for _, line := range strings.Split(string(status), "\n") {
		if !strings.HasPrefix(line, "VmRSS:") {
			continue
		}
		data := strings.Fields(line)
		if len(data) != 3 || data[2] != "kB" {
			return nil
		}
		value, err := strconv.ParseUint(data[1], 10, 64)
		if err != nil {
			return nil
		}
		return &discordProcessMetrics{CPUTicks: user + system, RSSBytes: value * 1024}
	}
	return nil
}

func monitorDiscord(ctx context.Context, pid int, client *discordtransport.Client, output io.Writer) {
	encoder := json.NewEncoder(output)
	emit := func(level string, final bool) bool {
		sample := struct {
			Level     string                   `json:"level"`
			Final     bool                     `json:"final"`
			Bot       *discordProcessMetrics   `json:"bot,omitempty"`
			Launcher  *discordProcessMetrics   `json:"launcher,omitempty"`
			Transport discordtransport.Metrics `json:"transport"`
		}{Level: level, Final: final, Launcher: discordProcess(os.Getpid()), Transport: client.Metrics()}
		if !final {
			sample.Bot = discordProcess(pid)
		}
		if err := encoder.Encode(sample); err != nil {
			fmt.Fprintln(os.Stderr, "Discord Studio: unable to write metrics")
			return false
		}
		return true
	}
	if !emit("INFO", false) {
		return
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if !emit("INFO", false) {
				return
			}
		case <-ctx.Done():
			emit("INFO", true)
			return
		}
	}
}
