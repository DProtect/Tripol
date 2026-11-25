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

func TestVelocityText(t *testing.T) {
	for _, c := range []struct {
		s, w string
		bad  bool
	}{
		{`show text(0)
show text(-123)
show text(9223372036854775807)
show text(-9223372036854775807-1)`, "0\n-123\n9223372036854775807\n-9223372036854775808\n", false},
		{`show slice("abcdef",1,4)
show len(slice("a",1,1))
show slice("hello",0,5)`, "bcd\n0\nhello\n", false},
		{`show slice("abc",-1,2)`, "TRIPOL runtime error: invalid input, text too long or allocation failed\n", true},
		{`show slice("abc",0,4)`, "TRIPOL runtime error: invalid input, text too long or allocation failed\n", true},
		{`show slice("abc",2,1)`, "TRIPOL runtime error: invalid input, text too long or allocation failed\n", true},
	} {
		t.Run(c.s, func(t *testing.T) { runCase(t, c.s, c.w, c.bad) })
	}
}
func TestVelocityFiles(t *testing.T) {
	if runtime.GOARCH != "amd64" || (runtime.GOOS != "linux" && runtime.GOOS != "windows") {
		t.Skip()
	}
	for _, profile := range []string{"nano", "debug"} {
		t.Run(profile, func(t *testing.T) {
			d := t.TempDir()
			p := filepath.Join(d, "program.exe")
			src := `show write_text("test.txt","Hello")
show append_text("test.txt",", TRIPOL!"+text(5))
show read_text("test.txt")
show write_text("test.txt","")
show len(read_text("test.txt"))
show append_text("new.txt","new")
show read_text("new.txt")`
			if err := os.WriteFile(p, compileProfile(src, runtime.GOOS, profile), 0755); err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(p)
			cmd.Dir = d
			o, e := cmd.CombinedOutput()
			if e != nil || string(o) != "5\n10\nHello, TRIPOL!5\n0\n0\n3\nnew\n" {
				t.Fatalf("%q %v", o, e)
			}
			b, e := os.ReadFile(filepath.Join(d, "test.txt"))
			if e != nil || len(b) != 0 {
				t.Fatal(b, e)
			}
			b, e = os.ReadFile(filepath.Join(d, "new.txt"))
			if e != nil || string(b) != "new" {
				t.Fatal(b, e)
			}
		})
	}
	for i, c := range []struct {
		data, src, w string
		missing, bad bool
	}{
		{strings.Repeat("a", 65535), `show len(read_text("file.txt"))`, "65535\n", false, false},
		{strings.Repeat("a", 65536), `show read_text("file.txt")`, "", false, true},
		{"a\x00b", `show read_text("file.txt")`, "", false, true},
		{"", `show read_text("file.txt")`, "", true, true},
		{"", `show write_text("missing/file.txt","x")`, "", false, true},
		{"a", `for i in range(0,1000){s=read_text("file.txt"); if s!="a" {show "bad"}}
show "ok"`, "ok\n", false, false},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			d := t.TempDir()
			p := filepath.Join(d, "p.exe")
			if !c.missing {
				if e := os.WriteFile(filepath.Join(d, "file.txt"), []byte(c.data), 0644); e != nil {
					t.Fatal(e)
				}
			}
			if e := os.WriteFile(p, compile(c.src, runtime.GOOS), 0755); e != nil {
				t.Fatal(e)
			}
			cmd := exec.Command(p)
			cmd.Dir = d
			o, e := cmd.CombinedOutput()
			if c.bad {
				if e == nil || !strings.Contains(string(o), "TRIPOL runtime error: file operation failed") {
					t.Fatalf("%q %v", o, e)
				}
			} else if e != nil || string(o) != c.w {
				t.Fatalf("%q %v", o, e)
			}
		})
	}
}
func TestVelocityPE(t *testing.T) {
	p, e := pecheck.NewFile(bytes.NewReader(compile(`show read_text("a")
show write_text("b","c")
show append_text("b","d")
show input()`, "windows")))
	if e != nil {
		t.Fatal(e)
	}
	s, e := p.ImportedSymbols()
	if e != nil {
		t.Fatal(e)
	}
	for _, n := range []string{"CreateFileA", "CloseHandle", "ReadFile", "WriteFile", "VirtualAlloc"} {
		if !strings.Contains(strings.Join(s, " "), n) {
			t.Fatal(n, s)
		}
	}
}
func TestVelocityReject(t *testing.T) {
	for _, s := range []string{`show text("a")`, `show slice("a",0)`, `show slice("a",0,false)`, `show read_text(1)`, `show write_text("a",1)`, `show append_text("a")`} {
		t.Run(s, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("accepted invalid source")
				}
			}()
			compile(s, "linux")
		})
	}
}
func TestVelocityFrameAndRHS(t *testing.T) {
	for _, c := range []struct{ s, w string }{
		{`func sub(a,b){return a-b}
show sub(9,3)`, "6\n"},
		{`func f(a,b){if a==b{return a*b};return a+b}
show f(4,4)
show f(4,5)`, "16\n9\n"},
		{`func f(a:float,b:float)->float{return a/b}
show f(9.0,2.0)`, "4.500000\n"},
		{`func f(a:string,b:string)->string{return a+b}
show f("a","b")`, "ab\n"},
		{`func f(a){if a>0{let b=2;return b};return 0}
show f(1)`, "2\n"},
	} {
		t.Run(c.s, func(t *testing.T) { runCase(t, c.s, c.w, false) })
	}
}
func TestVelocityArrays(t *testing.T) {
	for _, c := range []struct {
		s, w string
		bad  bool
	}{
		{`n=5
v=ints(n)
for i in range(0,len(v)){v[i]=i*10}
show v[0]
show v[4]
show len(v)`, "0\n40\n5\n", false},
		{`func make(n)->ints{return ints(n)}
func update(v:ints)->int{v[1]=99;return v[1]}
a=make(3)
b=a
show update(b)
show a[1]
show a[2]`, "99\n99\n0\n", false},
		{`a=ints(0)
show len(a)`, "0\n", false},
		{`a=ints(2)
show a[-1]`, "TRIPOL runtime error: array index out of bounds\n", true},
		{`a=ints(2)
a[2]=1`, "TRIPOL runtime error: array index out of bounds\n", true},
		{`a=ints(-1)`, "TRIPOL runtime error: invalid input, text too long or allocation failed\n", true},
		{`a=ints(65536)`, "TRIPOL runtime error: invalid input, text too long or allocation failed\n", true},
		{`for i in range(0,513){v=ints(65535)}`, "TRIPOL runtime error: invalid input, text too long or allocation failed\n", true},
	} {
		t.Run(c.s, func(t *testing.T) { runCase(t, c.s, c.w, c.bad) })
	}
	for _, s := range []string{`a=ints(2);show a==a`, `const a=ints(2);a[0]=1`, `a=ints(2);a[0]="s"`, `a=ints(2.0)`} {
		t.Run(s, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("accepted invalid source")
				}
			}()
			compile(s, "linux")
		})
	}
}
