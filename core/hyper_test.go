package core

import (
	"bytes"
	pecheck "debug/pe"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func hyperRun(t *testing.T, src, profile string) (string, bool) {
	t.Helper()
	if runtime.GOARCH != "amd64" || (runtime.GOOS != "linux" && runtime.GOOS != "windows") {
		t.Skip("native x64 required")
	}
	p := filepath.Join(t.TempDir(), "run.exe")
	if err := os.WriteFile(p, compileProfile(src, runtime.GOOS, profile), 0755); err != nil {
		t.Fatal(err)
	}
	out, e := exec.Command(p).CombinedOutput()
	return string(out), e != nil
}
func TestHyperRegisterLoops(t *testing.T) {
	cases := []string{
		`s=0;for i in range(-100,101){x=i-7;s+=x/8;s+=x%8};show s`,
		`s=100;for i in range(4,2){s+=i};show s`,
		`s=0;for i in range(9223372036854775805,9223372036854775807){s+=1};show s`,
		`s=10;b=4;for i in range(10){x=i+s;b+=x;s=b-i};show s,b`,
		`func f(n){s=0;for i in range(n){s+=i};return s};show f(100)`,
		`func f(n){for i in range(5){n+=i};return n};show f(10)`,
		`s=0;for i in range(4){for j in range(5){s+=j+i}};show s`,
		`s=0;for i in range(5){if i==2{continue};s+=i};show s`,
		`s=9223372036854775807;for i in range(1){s+=1};show s`,
		`s=0;try {for i in range(3){s+=1;show 1/(2-i)}} catch code {show code};show s`,
		`s=5;for i in range(s){s-=1};show s`,
		`s=1;for i in range(4){s=i*s+i};show s`,
	}
	for i, src := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			w, bad := hyperRun(t, src, "debug")
			for _, p := range []string{"nano", "balanced", "turbo"} {
				out, b := hyperRun(t, src, p)
				if w != out || bad != b {
					t.Fatalf("%s: got %q %v want %q %v", p, out, b, w, bad)
				}
			}
		})
	}
}
func TestHyperRegisterABI(t *testing.T) {
	src := `func f(a,b,c,d){return (a-b)*(c+d)}
func g(a,b=2){return a*b+1}
func deep(a,b,c,d){return a+(b+(c+(d+1)))}
func side(n){show n;return n}
show f(9,2,3,4),g(5),deep(1,2,3,4)
show f(side(1),side(2),side(3),side(4))
show f(g(1),g(2),g(3),g(4))
try {show f(9223372036854775807,-1,1,1)} catch code {show code}
show g(2,3)`
	w, bad := hyperRun(t, src, "debug")
	for _, p := range []string{"balanced", "turbo"} {
		o, b := hyperRun(t, src, p)
		if w != o || bad != b {
			t.Fatalf("%s %q %v vs %q %v", p, o, b, w, bad)
		}
	}
}
func TestHyperDifferential(t *testing.T) {
	rng := rand.New(rand.NewSource(1701))
	var s, w strings.Builder
	ops := []string{"+", "-", "*"}
	for i := 0; i < 100; i++ {
		a, b, c, d := int64(rng.Intn(201)-100), int64(rng.Intn(201)-100), int64(rng.Intn(201)-100), int64(rng.Intn(201)-100)
		op := ops[rng.Intn(3)]
		fmt.Fprintf(&s, "func f%d(a,b,c,d){return (a+b)%s(c-d)}\nshow f%d(%d,%d,%d,%d)\n", i, op, i, a, b, c, d)
		v := a + b
		switch op {
		case "+":
			v += c - d
		case "-":
			v -= c - d
		case "*":
			v *= c - d
		}
		fmt.Fprintln(&w, v)
	}
	for _, p := range []string{"debug", "balanced", "turbo"} {
		o, b := hyperRun(t, s.String(), p)
		if b || o != w.String() {
			t.Fatalf("%s mismatch %q", p, o)
		}
	}
}
func TestHyperClock(t *testing.T) {
	src := `import tripol.time
start=time.now_ns()
time.sleep(15)
assert time.elapsed_ns(start)>=10000000
show true`
	for _, p := range []string{"nano", "balanced", "turbo"} {
		o, b := hyperRun(t, src, p)
		if b || o != "true\n" {
			t.Fatalf("%s %q %v", p, o, b)
		}
	}
	f, e := pecheck.NewFile(bytes.NewReader(compile(src, "windows")))
	if e != nil {
		t.Fatal(e)
	}
	im, e := f.ImportedSymbols()
	if e != nil {
		t.Fatal(e)
	}
	for _, n := range []string{"QueryPerformanceCounter", "QueryPerformanceFrequency"} {
		if !strings.Contains(strings.Join(im, " "), n) {
			t.Fatal(n, im)
		}
	}
}
func TestHyperInternalMarker(t *testing.T) {
	for _, s := range []string{"", "TRIPOL_NS: -1", "TRIPOL_NS: 1\nTRIPOL_NS: 2", "TRIPOL_NS: hi", "TRIPOL_NS: 9223372036854775808"} {
		if _, ok := parseInternal(s); ok {
			t.Fatal(s)
		}
	}
	if n, ok := parseInternal("result 1\r\nTRIPOL_NS: 234\r\n"); !ok || n != 234 {
		t.Fatal(n, ok)
	}
}
func TestHyperProfileCode(t *testing.T) {
	src := `s=0;for i in range(100){s+=i};show s`
	for _, target := range []string{"linux", "windows"} {
		a, b := compileProfile(src, target, "nano"), compileProfile(src, target, "turbo")
		if bytes.Equal(a, b) {
			t.Fatal("register loop not emitted")
		}
		if !bytes.Contains(b, []byte{0x49, 0xff, 0xc0}) {
			t.Fatal("no register increment")
		}
		if !bytes.Equal(b, compileProfile(src, target, "turbo")) {
			t.Fatal("nondeterministic")
		}
	}
}
