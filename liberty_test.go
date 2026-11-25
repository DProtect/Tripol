package main

import (
	"bytes"
	pecheck "debug/pe"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLibertyPrograms(t *testing.T) {
	for i, c := range []struct{ s, w string }{
		{`show 1.25+2.5
show 8.0/4.0
show -2.125
show sqrt(9.0)
show float(7)`, "3.750000\n2.000000\n-2.125000\n3.000000\n7.000000\n"},
		{`func area(a:float,b:float)->float{return a*b}
show area(2.5,4.0)
x=0.5
while x<2.0 {show x; x=x+0.5}`, "10.000000\n0.500000\n1.000000\n1.500000\n"},
		{`show 1.0<2.0
show 1.0<=1.0
show 2.0>1.0
show 1.0>=2.0
show 1.0==1.0
show 1.0!=2.0`, "true\ntrue\ntrue\nfalse\ntrue\ntrue\n"},
		{`x=0.0/0.0
show x==x
show x!=x
show x<1.0
show x<=1.0
show x>1.0
show x>=1.0`, "false\ntrue\nfalse\nfalse\nfalse\nfalse\n"},
		{`x="Hola"
y=x+" mundo"
show y
show x
show len(y)
show y=="Hola mundo"`, "Hola mundo\nHola\n10\ntrue\n"},
		{`func join(a:string,b:string)->string{return a+b}
show join("a","b")+"c"
show len("é")
show len("")`, "abc\n2\n0\n"},
		{`show parse_int("-9223372036854775808")
show parse_int("9223372036854775807")
show parse_int("+0")`, "-9223372036854775808\n9223372036854775807\n0\n"},
		{`show -0.0000004
show 0.9999999
show 123456.123456`, "0.000000\n1.000000\n123456.123456\n"},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) { runCase(t, c.s, c.w, false) })
	}
}
func TestLibertyInput(t *testing.T) {
	if runtime.GOARCH != "amd64" || (runtime.GOOS != "linux" && runtime.GOOS != "windows") {
		t.Skip()
	}
	for i, c := range []struct {
		input, src, w string
		bad           bool
	}{
		{"Ana\r\n42\r\n", "name=input()\nshow \"Hola \"+name\nshow parse_int(input())*2", "Hola Ana\n84\n", false},
		{"uno\ndos", "show input()\nshow input()\nshow len(input())", "uno\ndos\n0\n", false},
		{strings.Repeat("a", 4096) + "\n", "show len(input())", "4096\n", false},
		{strings.Repeat("a", 4096) + "\r\n", "show len(input())", "4096\n", false},
		{strings.Repeat("a", 4097) + "\n", "show input()", "TRIPOL runtime error: invalid input, text too long or allocation failed\n", true},
		{"a\x00b\n", "show input()", "TRIPOL runtime error: invalid input, text too long or allocation failed\n", true},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			p := filepath.Join(t.TempDir(), "p.exe")
			if e := os.WriteFile(p, compile(c.src, runtime.GOOS), 0755); e != nil {
				t.Fatal(e)
			}
			cmd := exec.Command(p)
			cmd.Stdin = strings.NewReader(c.input)
			out, e := cmd.CombinedOutput()
			if string(out) != c.w || (e != nil) != c.bad {
				t.Fatalf("%q %v want %q", out, e, c.w)
			}
		})
	}
}
func TestLibertyTraps(t *testing.T) {
	for _, s := range []string{`show parse_int("9223372036854775808")`, `show parse_int("-9223372036854775809")`, `show parse_int("12a")`, `show parse_int("-")`, `show parse_int("")`} {
		t.Run(s, func(t *testing.T) {
			runCase(t, s, "TRIPOL runtime error: invalid input, text too long or allocation failed\n", true)
		})
	}
	for _, s := range []string{`show 1.0/0.0`, `show sqrt(-1.0)`, `show 10000000000000.0`} {
		t.Run(s, func(t *testing.T) { runCase(t, s, "TRIPOL runtime error: float outside printable range\n", true) })
	}
}
func TestLibertyReject(t *testing.T) {
	for _, s := range []string{`show 1.0+2`, `show 1.0%2.0`, `show sqrt(2)`, `show input(1)`, `show parse_int(1)`, `show len(1)`, `show 1.0+"a"`} {
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
func TestLibertyPE(t *testing.T) {
	p, e := pecheck.NewFile(bytes.NewReader(compile("show input()+\"!\"\nshow 1.5", "windows")))
	if e != nil {
		t.Fatal(e)
	}
	names, e := p.ImportedSymbols()
	if e != nil {
		t.Fatal(e)
	}
	s := strings.Join(names, " ")
	for _, n := range []string{"ReadFile", "VirtualAlloc", "GetLastError", "GetStdHandle", "WriteFile", "ExitProcess"} {
		if !strings.Contains(s, n) {
			t.Fatalf("missing %s: %v", n, names)
		}
	}
}

func TestLibertyArenaLimit(t *testing.T) {
	runCase(t, `s="a"
for i in range(0,17){s=s+s}
show s`, "TRIPOL runtime error: invalid input, text too long or allocation failed\n", true)
	runCase(t, `s="abcdefghijklmnopqrstuvwxyz"
for i in range(0,400000){s="abcdefghijklmnopqrstuvwxyz"+"x"}
show s`, "TRIPOL runtime error: invalid input, text too long or allocation failed\n", true)
}
func TestLocalModules(t *testing.T) {
	d := t.TempDir()
	mk := func(name, src string) {
		p := filepath.Join(d, name)
		if e := os.MkdirAll(filepath.Dir(p), 0755); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(p, []byte(src), 0644); e != nil {
			t.Fatal(e)
		}
	}
	mk("main.tri", "import \"lib/math.tri\"\nimport \"lib/math.tri\"\nshow double(12)\nshow area(2.0,3.0)")
	mk("lib/math.tri", "import \"base.tri\"\nfunc area(a:float,b:float)->float{return a*b}")
	mk("lib/base.tri", "func double(x){return x*2}")
	fs, ss, files := loadProject(filepath.Join(d, "main.tri"))
	if len(files) != 3 || len(fs) != 2 {
		t.Fatal(files, len(fs))
	}
	b := compileParsed(fs, ss, runtime.GOOS, "nano")
	p := filepath.Join(d, "p.exe")
	if e := os.WriteFile(p, b, 0755); e != nil {
		t.Fatal(e)
	}
	o, e := exec.Command(p).CombinedOutput()
	if e != nil || string(o) != "24\n6.000000\n" {
		t.Fatalf("%q %v", o, e)
	}
	for _, c := range []struct{ name, src string }{{"cycle.tri", `import "cycle.tri"`}, {"bad.tri", `import "lib/body.tri"`}, {"outside.tri", `import "../outside.tri"`}, {"absolute.tri", `import "/tmp/x.tri"`}, {"notfirst.tri", "show 1\nimport \"lib/base.tri\""}} {
		t.Run(c.name, func(t *testing.T) {
			mk(c.name, c.src)
			mk("lib/body.tri", "show 1")
			defer func() {
				if recover() == nil {
					t.Fatal("accepted invalid import")
				}
			}()
			loadProject(filepath.Join(d, c.name))
		})
	}
}
func TestLibertyDebug(t *testing.T) {
	if runtime.GOARCH != "amd64" || (runtime.GOOS != "linux" && runtime.GOOS != "windows") {
		t.Skip()
	}
	p := filepath.Join(t.TempDir(), "p.exe")
	s := `show "ab"+input()
show 2.0*3.25
show len("foo")`
	if e := os.WriteFile(p, compileProfile(s, runtime.GOOS, "debug"), 0755); e != nil {
		t.Fatal(e)
	}
	c := exec.Command(p)
	c.Stdin = strings.NewReader("cd\n")
	b, e := c.CombinedOutput()
	if e != nil || string(b) != "abcd\n6.500000\n3\n" {
		t.Fatalf("%q %v", b, e)
	}
}
