package core

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
)

func TestTitanRegisterBranches(t *testing.T) {
	cases := []string{
		`a=0;b=0;c=0;d=0;for i in range(100){a+=i;b+=a;if i%2==0{c+=b}else{d+=a}};show a,b,c,d`,
		`a=2;b=3;c=4;d=5;e=6;for i in range(50){a+=i;b+=a;c+=b;d+=c;e+=d};show a,b,c,d,e`,
		`a=1;b=2;for i in range(30){if i==5{continue};a+=i;if a>b{b+=a}else{b-=i};if i==15{break}};show a,b`,
		`func f(a,b){for i in range(30){if a>b{a-=1}else{b-=2};a+=i};return a+b};show f(9,8)`,
		`a=2;b=3;for i in range(4,2){a+=i;b+=a};show a,b`,
		`a=0;b=0;for i in range(5){if i==2{x=i; a+=x}else{b+=i}};show a,b`,
		`a=0;b=0;for i in range(10){x=i+1;if x>5{a+=x}else{b+=x};x+=1;a+=x};show a,b`,
		`a=0;b=9223372036854775807;for i in range(10){a+=1;b+=1};show a,b`,
		`a=0;b=0;try{for i in range(5){a+=i;b+=i;if i==2{a+=1/0}}}catch code{show code};show a,b`,
		`a=0;b=0;for i in range(100){if i<2 or i>95{a+=i}else{b+=i};if not (a<b){a+=1}};show a,b`,
	}
	for i, src := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			w, err := hyperRun(t, src, "debug")
			for _, p := range []string{"nano", "balanced", "turbo"} {
				o, b := hyperRun(t, src, p)
				if o != w || b != err {
					t.Fatalf("%s: %q %v != %q %v", p, o, b, w, err)
				}
			}
		})
	}
}
func TestTitanRandomSignedDivision(t *testing.T) {
	rng := rand.New(rand.NewSource(1801))
	var src, w strings.Builder
	for i := 0; i < 150; i++ {
		x := int64(rng.Uint64())
		if i == 0 {
			x = math.MinInt64
		}
		if i == 1 {
			x = math.MaxInt64
		}
		lit := fmt.Sprint(x)
		if x == math.MinInt64 {
			lit = "(-9223372036854775807-1)"
		}
		n := int64(1) << uint(rng.Intn(31))
		fmt.Fprintf(&src, "x=%s\nshow x/%d,x%%%d\n", lit, n, n)
		fmt.Fprintf(&w, "%d %d\n", x/n, x%n)
	}
	for _, p := range []string{"nano", "balanced", "turbo"} {
		out, bad := hyperRun(t, src.String(), p)
		if bad || out != w.String() {
			t.Fatalf("%s signed division mismatch", p)
		}
	}
}
func TestTitanRegisterMachineCode(t *testing.T) {
	src := `a=0;b=0;c=0;d=0;for i in range(100){a+=i;b+=a;c+=b;d+=c};show a,b,c,d`
	for _, target := range []string{"linux", "windows"} {
		b := compileProfile(src, target, "turbo")
		for _, op := range [][]byte{{0x41, 0x54}, {0x41, 0x55}, {0x41, 0x56}, {0x41, 0x5e}, {0x41, 0x5d}, {0x41, 0x5c}} {
			if !bytes.Contains(b, op) {
				t.Fatalf("missing register save/restore %x", op)
			}
		}
		if !bytes.Equal(b, compileProfile(src, target, "turbo")) {
			t.Fatal("nondeterministic")
		}
	}
}
func TestTitanArenaOutput(t *testing.T) {
	for _, v := range []string{"42\nTRIPOL_NS: 123\n", "42\r\nTRIPOL_NS: 123\r\n"} {
		if ns, ok := arenaOutput(v, "42"); !ok || ns != 123 {
			t.Fatal(v)
		}
	}
	for _, v := range []string{"43\nTRIPOL_NS: 123", "42\nTRIPOL_NS: 123\nTRIPOL_NS: 3", "42\nTRIPOL_NS: -1", "42\nerror\nTRIPOL_NS: 1"} {
		if _, ok := arenaOutput(v, "42"); ok {
			t.Fatal(v)
		}
	}
}
