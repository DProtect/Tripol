package main

import (
	"bytes"
	pecheck "debug/pe"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestVisionDiagnostics(t *testing.T) {
	for _, x := range []struct {
		src, code string
		line, col int
	}{{"show 1\nshow @", "TP1001", 2, 6}, {"show (1", "TP1002", 1, 8}, {"show missing", "TP1003", 1, 0}, {"let x:int=\"abc\"", "TP1004", 1, 0}} {
		t.Run(x.code, func(t *testing.T) {
			defer func() {
				d, ok := recover().(diagnostic)
				if !ok || d.Code != x.code || d.Line != x.line || d.Column != x.col {
					t.Fatalf("bad diagnostic %#v", d)
				}
			}()
			im, fs, ss := parseNamed(x.src, "demo.tri")
			_ = im
			compileParsed(fs, ss, "linux", "nano")
		})
	}
}
func TestVisionModuleLocation(t *testing.T) {
	d := t.TempDir()
	main := filepath.Join(d, "main.tri")
	mod := filepath.Join(d, "lib.tri")
	os.WriteFile(main, []byte("import \"lib.tri\"\nshow f()"), 0600)
	os.WriteFile(mod, []byte("func f(){\nreturn missing\n}"), 0600)
	defer func() {
		v, ok := recover().(diagnostic)
		if !ok || v.File != mod || v.Line != 2 {
			t.Fatalf("bad location %#v", v)
		}
	}()
	fs, ss, _ := loadProject(main)
	compileParsed(fs, ss, "linux", "nano")
}
func TestVisionSafeInt(t *testing.T) {
	runCase(t, `import tripol.errors
show errors.int_or("42",-1)
show errors.int_or("abc",-1)
show errors.int_or("9223372036854775808",-1)
show errors.int_or("-9223372036854775808",0)
show errors.is_int("-9223372036854775809")
show errors.is_int("+")
show errors.is_int("+9223372036854775807")`, "42\n-1\n-1\n-9223372036854775808\nfalse\nfalse\ntrue\n", false)
}
func TestVisionImports(t *testing.T) {
	for _, s := range []string{`import tripol.ui
show ui.message("Title","Hi")`, `import tripol.ui
import tripol.network
show ui.confirm("Title","Continue?")
show network.exchange(network.loopback(),8765,"hi")`} {
		p, e := pecheck.NewFile(bytes.NewReader(compile(s, "windows")))
		if e != nil {
			t.Fatal(e)
		}
		syms, e := p.ImportedSymbols()
		if e != nil || !strings.Contains(strings.Join(syms, " "), "MessageBoxA:USER32.dll") {
			t.Fatalf("imports %v %v", syms, e)
		}
	}
	if runtime.GOOS == "linux" {
		runCase(t, `import tripol.ui
show ui.message("T","M")`, "TRIPOL runtime error: native dialogs require Windows or UI operation failed\n", true)
	}
}
func TestVisionExitCodes(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip()
	}
	for _, x := range []struct {
		s    string
		code int
	}{{"show 1/0", 10}, {"a=[1];show a[1]", 11}, {"assert false", 12}, {`show parse_int("bad")`, 13}} {
		p := filepath.Join(t.TempDir(), "program.exe")
		os.WriteFile(p, compile(x.s, runtime.GOOS), 0755)
		e := exec.Command(p).Run()
		if ex, ok := e.(*exec.ExitError); !ok || ex.ExitCode() != x.code {
			t.Fatalf("%s: %v", x.s, e)
		}
	}
}
func TestVisionRecoverableOps(t *testing.T) {
	runCase(t, `import tripol.errors
show errors.div_or(9,0,-1),errors.div_or(8,2,-1),errors.div_or(-9223372036854775807-1,-1,7)
a=[3,4]
show errors.get_or(a,0,99),errors.get_or(a,-1,99),errors.get_or(a,2,99)`, "-1 4 7\n3 99 99\n", false)
}
