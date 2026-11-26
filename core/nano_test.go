package core

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNanoPrograms(t *testing.T) {
	for i, c := range []struct{ s, w string }{
		{"a=[10,20,30]\nshow a[0]\na[1]=a[0]+5\nshow a[1]\nshow len(a)", "10\n15\n3\n"},
		{"const a=[1,2]\nshow a[1]", "2\n"},
		{"sum=0\nfor i in range(0,10){if i==3{continue}\nif i==6{break}\nsum=sum+i}\nshow sum", "12\n"},
		{"i=0\nwhile true {i=i+1\nif i==2{continue}\nif i==4{break}\nshow i}", "1\n3\n"},
		{"n=3\nfor i in range(0,n){show i\nn=0}", "0\n1\n2\n"},
		{"for i in range(5,2){show i}\nshow 99", "99\n"},
		{"func loop(n,a){if n==0{return a}\nreturn loop(n-1,a+1)}\nshow loop(1000000,0)", "1000000\n"},
		{"func swap(n,a,b){if n==0{return a-b}\nreturn swap(n-1,b,a)}\nshow swap(3,10,3)", "-7\n"},
		{"show (2+3)*4\nshow 5>2\nshow !(true && false)\nshow -5*8", "20\ntrue\ntrue\n-40\n"},
		{"x=3\nshow x+128\nshow x-129\nshow x*1000\nif x>=3{show 1}\nif x<3{show 2}else{show 3}", "131\n-126\n3000\n1\n3\n"},
		{"\ufeffshow 7", "7\n"},
		{"func f(n){if n<=1{return n}\nreturn f(n-1)+f(n-2)}\nshow f(35)", "9227465\n"},
		{"x=0\nif true{a=[1,2];x=a[1]}\nif true{b=[8,9];show b[0]}\nshow x", "8\n2\n"},
		{"a=[1,2]\nfor i in range(0,len(a)){a[i]=a[i]*2}\nshow a[1]", "4\n"},
		{"a=[1]\nfunc len(x){return x+1}\nshow len(7)", "8\n"},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) { runCase(t, c.s, c.w, false) })
	}
}
func TestNanoErrors(t *testing.T) {
	for i, s := range []string{"break", "continue", "const a=[1]\na[0]=2", "a=[true]", "a=[]", "a=[1]\nshow a[true]", "a=[1]\na[0]=true", "for i in range(true,3){}", "for i in range(0,3){i=2}", "func dead()->int{return false}\nshow 1", "a=[1]\na=[2]", "show len(3)"} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("accepted invalid input")
				}
			}()
			compile(s, "linux")
		})
	}
}
func TestBounds(t *testing.T) {
	for _, s := range []string{"a=[1]\nshow a[-1]", "a=[1]\nshow a[1]", "a=[1]\na[999]=2"} {
		t.Run(s, func(t *testing.T) { runCase(t, s, "TRIPOL runtime error: array index out of bounds\n", true) })
	}
}
func TestNanoSize(t *testing.T) {
	s := "func dead(x){return x*12}\nshow 2+3"
	for _, target := range []string{"linux", "windows"} {
		b := compile(s, target)
		d := compileProfile(s, target, "debug")
		if len(b) >= len(d) {
			t.Fatalf("%s nano %d debug %d", target, len(b), len(d))
		}
		if !bytes.Equal(b, compile(s, target)) {
			t.Fatal("non-deterministic")
		}
	}
}
func TestProfileEquivalence(t *testing.T) {
	if runtime.GOARCH != "amd64" || (runtime.GOOS != "linux" && runtime.GOOS != "windows") {
		t.Skip()
	}
	s := "a=[4,8,12]\nfor i in range(0,len(a)){show a[i]+2}\nshow 4*5"
	for _, profile := range []string{"nano", "debug"} {
		p := filepath.Join(t.TempDir(), "program.exe")
		if e := os.WriteFile(p, compileProfile(s, runtime.GOOS, profile), 0755); e != nil {
			t.Fatal(e)
		}
		b, e := exec.Command(p).CombinedOutput()
		if e != nil || string(b) != "6\n10\n14\n20\n" {
			t.Fatalf("%s %q %v", profile, b, e)
		}
	}
}
