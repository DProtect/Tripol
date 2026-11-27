package core

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestAutomaticScratchFrameForWhile(t *testing.T) {
	if runtime.GOARCH != "amd64" || (runtime.GOOS != "linux" && runtime.GOOS != "windows") {
		t.Skip("native x86-64 host required")
	}
	const source = `func show_temporaries()->int {
		while input()=="more" {show text(17)}
		return 42
	}
	show show_temporaries()`
	for _, profile := range []string{"debug", "nano", "balanced", "turbo"} {
		t.Run(profile, func(t *testing.T) {
			exe := filepath.Join(t.TempDir(), "scratch.exe")
			if err := os.WriteFile(exe, compileProfile(source, runtime.GOOS, profile), 0755); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, exe)
			cmd.Stdin = strings.NewReader("more\nmore\nstop\n")
			out, err := cmd.CombinedOutput()
			if err != nil || string(out) != "17\n17\n42\n" {
				t.Fatalf("output %q, error %v", out, err)
			}
		})
	}
}
