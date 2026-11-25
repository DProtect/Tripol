package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNexusSyntax(t *testing.T) {
	cases := []struct{ s, w string }{
		{`show "Total:", 20+5, true, 1.25, "ok"`, "Total: 25 true 1.250000 ok\n"},
		{`show "", -9223372036854775807-1, false, -1.25`, " -9223372036854775808 false -1.250000\n"},
		{`total=0
repeat 5 { total += 2 }
show total
for i in range(3) { show i }
repeat 0 { show 999 }`, "10\n0\n1\n2\n"},
		{`x=4
x*=3
x-=2
x/=2
x%=3
show x
s="a"
s+="b"
show s`, "2\nab\n"},
		{`x=5
if x==0 {show 0} else if x==5 {show 5} else {show 9}
show not x==0 and true or false
show false and (1/0==0)`, "5\ntrue\nfalse\n"},
		{`a=[10,20,30]
total=0
for v in a { total+=v }
show total, sum(a)
b=ints(3)
b[0]=1
b[1]=2
b[2]=3
for v in b { if v==2 {continue}; show v }
show sum(b), sum_empty() 
func sum_empty(){a=ints(0);return sum(a)}`, "60 60\n1\n3\n6 0\n"},
		{`show min(-10,2), max(-10,2), abs(-45)
show min(5,5), max(-3,-9), abs(0)`, "-10 2 45\n5 -3 0\n"},
		{`a=[1,2,3]
for x in a {repeat 2 {show x}; if x==2 {break}}`, "1\n1\n2\n2\n"},
		{`func pick(x){if x<0{return -1}else if x==0{return 0}else{return 1}}
show pick(-2),pick(0),pick(8)`, "-1 0 1\n"},
	}
	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) { runCase(t, c.s, c.w, false) })
	}
}
func TestNexusInput(t *testing.T) {
	for _, profile := range []string{"nano", "debug"} {
		t.Run(profile, func(t *testing.T) {
			p := filepath.Join(t.TempDir(), "p.exe")
			src := `name=input("Name?")
n=input_int("Number?")
show "Hello",name,n*2`
			if e := os.WriteFile(p, compileProfile(src, runtime.GOOS, profile), 0755); e != nil {
				t.Fatal(e)
			}
			cmd := exec.Command(p)
			cmd.Stdin = strings.NewReader("Ana\r\n25\r\n")
			o, e := cmd.CombinedOutput()
			if e != nil || string(o) != "Name?\nNumber?\nHello Ana 50\n" {
				t.Fatalf("%q %v", o, e)
			}
		})
	}
}
func TestNexusReject(t *testing.T) {
	for _, src := range []string{`x+=1`, `const x=1; x+=2`, `for x in "a" {show x}`, `x="a";for y in x {show y}`, `repeat true {show 1}`, `show min(1,true)`, `show sum(1)`, `show input_int(1)`, `a=[1]; show "a",a`, `show 1,`, `x=1; for x in range(3){show x}`} {
		t.Run(src, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("accepted invalid program")
				}
			}()
			compile(src, "linux")
		})
	}
}
func TestNexusTraps(t *testing.T) {
	for _, src := range []string{`show abs(-9223372036854775807-1)`, `a=[9223372036854775807,1];show sum(a)`, `x=9223372036854775807;x+=1`} {
		t.Run(src, func(t *testing.T) { runCase(t, src, "TRIPOL runtime error: invalid integer arithmetic\n", true) })
	}
}
func TestNexusPE(t *testing.T) {
	src := `a=[1,2];for n in a {show "Value",n,2.5,true};show input_int("Number?")`
	b := compile(src, "windows")
	if len(b) < 1024 || string(b[:2]) != "MZ" {
		t.Fatal("bad PE")
	}
}

func TestNexusTemplates(t *testing.T) {
	for _, name := range []string{"basic", "calculator", "files"} {
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join(t.TempDir(), name)
			newProjectTemplate(dir, name)
			fs, ss, _ := loadProject(filepath.Join(dir, "main.tri"))
			compileParsed(fs, ss, "windows", "nano")
			compileParsed(fs, ss, "linux", "nano")
		})
	}
}
