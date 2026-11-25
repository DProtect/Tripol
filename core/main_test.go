package core

import (
	"bytes"
	elfcheck "debug/elf"
	pecheck "debug/pe"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func runCase(t *testing.T, s, w string, bad bool) {
	t.Helper()
	if runtime.GOARCH != "amd64" || (runtime.GOOS != "linux" && runtime.GOOS != "windows") {
		t.Skip("requires supported host")
	}
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	p := filepath.Join(t.TempDir(), "test"+ext)
	profile := os.Getenv("TRIPOL_TEST_PROFILE")
	if profile == "" {
		profile = "nano"
	}
	if e := os.WriteFile(p, compileProfile(s, runtime.GOOS, profile), 0755); e != nil {
		t.Fatal(e)
	}
	o, e := exec.Command(p).CombinedOutput()
	if (e != nil) != bad || string(o) != w {
		t.Fatalf("output %q, error %v, want %q (bad=%v)", o, e, w, bad)
	}
}
func TestPrograms(t *testing.T) {
	for _, c := range []struct{ s, w string }{
		{"show \"Hola\"", "Hola\n"}, {"show 2 + 3 * 4\nshow (2 + 3) * 4", "14\n20\n"}, {"x=1\nx=x+3\nshow x", "4\n"}, {"let x:bool=true\nconst s:string=\"ok\"\nshow x\nshow s", "true\nok\n"}, {"func f(a,b){return a-b}\nshow f(9,4)", "5\n"}, {"func f(x:string)->string{return x}\nshow f(\"yes\")", "yes\n"}, {"func fib(n){if n<=1{return n}\nreturn fib(n-1)+fib(n-2)}\nshow fib(15)", "610\n"}, {"show foo(5)\nfunc foo(x){return x*2}", "10\n"}, {"i=0\nwhile i<3 {show i;i=i+1}", "0\n1\n2\n"}, {"if false {show 1}else{show 2}", "2\n"}, {"show !false && true\nshow false||true", "true\ntrue\n"}, {"show false && (1/0==1)\nshow true || (1/0==1)", "false\ntrue\n"}, {"show \"a\"==\"a\"\nshow \"a\"!=\"b\"", "true\ntrue\n"}, {`show "a\nb\t\"c"`, "a\nb\t\"c\n"}, {"# comment\nshow 5 // comment", "5\n"}, {"show -7/3\nshow -7%3", "-2\n-1\n"}, {"func f(x){if x==0{return 1}else{return 2}}\nshow f(0)", "1\n"}, {"x=0\nif true{let y=3;x=y}\nshow x", "3\n"}, {"show 9223372036854775807\nshow -9223372036854775807-1\nshow 0", "9223372036854775807\n-9223372036854775808\n0\n"}, {"show 1<2\nshow 2>=2\nshow 2>3", "true\ntrue\nfalse\n"}, {"show \"\"", "\n"}, {"show \"¡Hola, TRIPOL!\"", "¡Hola, TRIPOL!\n"}} {
		t.Run(c.s, func(t *testing.T) { runCase(t, c.s, c.w, false) })
	}
}
func TestTraps(t *testing.T) {
	for _, s := range []string{"show 1/0", "show 9223372036854775807+1", "show (-9223372036854775807-1)/-1", "show 9223372036854775807*2", "show -(-9223372036854775807-1)", "show -9223372036854775807-2"} {
		t.Run(s, func(t *testing.T) { runCase(t, s, "TRIPOL runtime error: invalid integer arithmetic\n", true) })
	}
}
func TestReject(t *testing.T) {
	for _, s := range []string{"const x=1\nx=2", "let x:int=\"a\"", "show x", "show f(1)", "func f(x){return x}\nshow f()", "func f(x){if x==0{return x}}", "if 1{show 1}", "x=1\nif true{let x=2}", "if true{show 1", "show @\nshow 1", "show 9223372036854775808", "if true{x=1}\nshow x", "return 1", "func f(x,x){return x}", "func f(x){return x}\nfunc f(x){return x}", "show \"\\u0000\""} {
		t.Run(s, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("accepted invalid program")
				}
			}()
			compile(s, "linux")
		})
	}
}
func TestFormats(t *testing.T) {
	p, e := pecheck.NewFile(bytes.NewReader(compile("show 1", "windows")))
	if e != nil {
		t.Fatal(e)
	}
	libs, e := p.ImportedSymbols()
	if e != nil || len(libs) != 3 || !strings.Contains(strings.Join(libs, " "), "KERNEL32.dll") {
		t.Fatalf("imports %v %v", libs, e)
	}
	f, e := elfcheck.NewFile(bytes.NewReader(compile("show 1", "linux")))
	if e != nil || f.Machine != elfcheck.EM_X86_64 {
		t.Fatalf("bad ELF %v", e)
	}
}
func TestDeterminism(t *testing.T) {
	s := "show \"a\"\nshow \"b\""
	for _, target := range []string{"windows", "linux"} {
		if !bytes.Equal(compile(s, target), compile(s, target)) {
			t.Fatal("non-deterministic")
		}
	}
}
