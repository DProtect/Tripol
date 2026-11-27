package tests

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

func TestPublicCompilerCLI(t *testing.T) {
	if runtime.GOARCH != "amd64" || (runtime.GOOS != "linux" && runtime.GOOS != "windows") {
		t.Skip("native x86-64 Windows or Linux required")
	}
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	dir := t.TempDir()
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	compiler := filepath.Join(dir, "tripol"+ext)
	cmd := exec.CommandContext(ctx, "go", "build", "-o", compiler, ".")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build public CLI: %v\n%s", err, out)
	}
	run := func(exe string, args ...string) string {
		t.Helper()
		cmd := exec.CommandContext(ctx, exe, args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "TRIPOL_NO_UPDATE=1")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s %v: %v\n%s", exe, args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	if v := run(compiler, "--version"); !strings.HasPrefix(v, "TRIPOL ") {
		t.Fatalf("unexpected version: %q", v)
	}
	src := filepath.Join(dir, "hello.tri")
	if err := os.WriteFile(src, []byte("import tripol.math\nshow math.square(7)\n"), 0600); err != nil {
		t.Fatal(err)
	}
	program := filepath.Join(dir, "hello"+ext)
	run(compiler, "build", src, "--target", runtime.GOOS, "-o", program)
	if got := run(program); got != "49" {
		t.Fatalf("native output: %q", got)
	}
	if runtime.GOOS == "linux" {
		src = filepath.Join(dir, "bot.tri")
		if err := os.WriteFile(src, []byte("import tripol.discord\nimport tripol.env\nshow env.get(\"TRIPOL_DISCORD_PIPE\")\nb=discord.bot()\nb.start()\n"), 0600); err != nil {
			t.Fatal(err)
		}
		bot := filepath.Join(dir, "bot")
		run(compiler, "build", src, "--target", "linux", "-o", bot)
		cmd := exec.CommandContext(ctx, bot)
		cmd.Env = append(os.Environ(), "DISCORD_TOKEN=", "TRIPOL_DISCORD_PIPE=")
		output, err := cmd.CombinedOutput()
		if err == nil || !strings.HasPrefix(string(output), "1\n") || !strings.Contains(string(output), "TP2003") || strings.Contains(string(output), "requires the TRIPOL launcher") {
			t.Fatalf("standalone Discord executable: %v\n%s", err, output)
		}
	}

}
