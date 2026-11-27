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

func TestGenesisDirectProtectsExistingFiles(t *testing.T) {
	if runtime.GOARCH != "amd64" || (runtime.GOOS != "linux" && runtime.GOOS != "windows") {
		t.Skip("native x86-64 host required")
	}
	for _, driver := range []string{"direct", "windows"} {
		t.Run(driver, func(t *testing.T) {
			dir := t.TempDir()
			stage := filepath.Join(dir, "stage.exe")
			fs, ss, _ := loadProject("compiler/" + driver + ".tri")
			if err := os.WriteFile(stage, compileParsed(fs, ss, runtime.GOOS, "nano"), 0755); err != nil {
				t.Fatal(err)
			}
			source := filepath.Join(dir, "input.tri")
			const program = "show 42\n"
			if err := os.WriteFile(source, []byte(program), 0600); err != nil {
				t.Fatal(err)
			}
			check := func(t *testing.T, output, original string) {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				cmd := exec.CommandContext(ctx, stage)
				cmd.Stdin = strings.NewReader(source + "\n" + output + "\n")
				log, err := cmd.CombinedOutput()
				if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 14 {
					t.Fatalf("expected file error, got %v (%q)", err, log)
				}
				for path, want := range map[string]string{source: program, output: original} {
					data, err := os.ReadFile(path)
					if err != nil || string(data) != want {
						t.Fatalf("%s changed after refusal (%v)", path, err)
					}
				}
			}
			t.Run("source", func(t *testing.T) { check(t, source, program) })
			t.Run("existing_output", func(t *testing.T) {
				output := filepath.Join(dir, "output.bin")
				const original = "previous output"
				if err := os.WriteFile(output, []byte(original), 0600); err != nil {
					t.Fatal(err)
				}
				check(t, output, original)
			})
			t.Run("alias_of_source", func(t *testing.T) {
				output := filepath.Join(dir, "alias.bin")
				if err := os.Link(source, output); err != nil {
					t.Skipf("hard links unavailable: %v", err)
				}
				check(t, output, program)
			})
		})
	}
}
