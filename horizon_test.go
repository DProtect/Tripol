package main

import (
	"bytes"
	pecheck "debug/pe"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestHorizonRecords(t *testing.T) {
	cases := []struct{ s, w string }{
		{`struct Player { name:string; score:int; alive:bool; height:float }
p=Player("Ana",10,true,1.5)
p.score+=5
show p.name,p.score,p.alive,p.height`, "Ana 15 true 1.500000\n"},
		{`struct Point { x:int, y:int }
func move(p:Point, dx:int)->Point {p.x+=dx;return p}
func create()->Point {return Point(2,3)}
p=create()
q=move(p,4)
show p.x,q.y`, "6 3\n"},
		{`struct Data { a:ints }
a=ints(2);a[0]=7;d=Data(a);b=d.a; b[1]=8;show sum(b)`, "15\n"},
		{`struct X { x:int }
show X(7).x
p=X(3);q=p;q.x*=2;show p.x`, "7\n6\n"},
		{`struct X { s:string }
p=X("abc"+"def");p.s+="!";show p.s`, "abcdef!\n"},
		{`struct T { v:int }
func f(p:T,n:int)->T {if n==0{return p};p.v+=1;return f(p,n-1)}
p=f(T(0),10000);show p.v`, "10000\n"},
		{`s=0;for i in range(-3,4){s+=i};show s;for i in range(9223372036854775806,9223372036854775807){show i}`, "0\n9223372036854775806\n"},
		{`assert 1+2==3
show "ok"`, "ok\n"},
	}
	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			runCase(t, c.s, c.w, false)
			if i == 5 {
				return
			} // debug disables tail-call elimination
			p := filepath.Join(t.TempDir(), "debug.exe")
			if e := os.WriteFile(p, compileProfile(c.s, runtime.GOOS, "debug"), 0755); e != nil {
				t.Fatal(e)
			}
			o, e := exec.Command(p).CombinedOutput()
			if e != nil || string(o) != c.w {
				t.Fatalf("debug %q %v", o, e)
			}
		})
	}
}
func TestHorizonReject(t *testing.T) {
	for _, s := range []string{`struct X {}`, `struct X { a:int; a:int }`, `struct X { a:Unknown }`, `struct X { a:X }`, `struct X { a:int };p=X(true)`, `struct X { a:int };p=X()`, `struct X { a:int };p=X(1);show p.b`, `struct X { a:int };const p=X(1);p.a=2`, `struct X { a:int };p=X(1);p.a="x"`, `struct X { a:int };p=X(1);show p`, `struct X { a:int };p=X(1);show p==p`, `func f(x:Unknown){return 1}`, `struct X { a:int };X(1).a=2`, `assert 1`, `assert false, 1`, `struct X { a:int };func X(){return 0}`} {
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
func TestHorizonAssertion(t *testing.T) {
	runCase(t, `assert false, "expected failure"`, "TRIPOL assertion failed at line 1: expected failure\n", true)
}
func TestHorizonModulesPE(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "main.tri")
	os.WriteFile(filepath.Join(d, "data.tri"), []byte(`struct Item { value:int }
func make(x:int)->Item {return Item(x)}`), 0644)
	os.WriteFile(p, []byte(`import "data.tri"
i=make(5);i.value+=1;assert i.value==6;show i.value`), 0644)
	fs, ss, _ := loadProject(p)
	b := compileParsed(fs, ss, "windows", "nano")
	f, e := pecheck.NewFile(bytes.NewReader(b))
	if e != nil {
		t.Fatal(e)
	}
	syms, e := f.ImportedSymbols()
	if e != nil || len(syms) < 3 {
		t.Fatalf("imports %v %v", syms, e)
	}
	p = filepath.Join(d, "out.exe")
	os.WriteFile(p, compileParsed(fs, ss, runtime.GOOS, "nano"), 0755)
	o, e := exec.Command(p).CombinedOutput()
	if e != nil || string(o) != "6\n" {
		t.Fatalf("%q %v", o, e)
	}
}
func TestHorizonTemplate(t *testing.T) {
	d := filepath.Join(t.TempDir(), "records")
	newProjectTemplate(d, "records")
	fs, ss, _ := loadProject(filepath.Join(d, "main.tri"))
	compileParsed(fs, ss, "linux", "nano")
	compileParsed(fs, ss, "windows", "nano")
}
