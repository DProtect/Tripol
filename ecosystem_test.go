package main

import (
	"bytes"
	pecheck "debug/pe"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestEcosystemLibraries(t *testing.T) {
	cases := []struct{ src, want string }{
		{`import tripol.sort
import tripol.search
n=[90,20,5,100,15,5,-1]
b=sort.ascending(n)
show b
show n
show sort.descending(n)
show sort.sorted(b),sort.sorted(n),search.binary(b,5),search.linear(n,100),search.last(n,5),search.lower_bound(b,99)`, "[-1, 5, 5, 15, 20, 90, 100]\n[90, 20, 5, 100, 15, 5, -1]\n[100, 90, 20, 15, 5, 5, -1]\ntrue false 1 3 5 6\n"},
		{`import tripol.sort
import tripol.search
empty=ints(0)
show sort.ascending(empty)
show sort.sorted(empty),search.binary(empty,9),search.lower_bound(empty,1),search.last(empty,1)
a=[7];show sort.ascending(a)`, "[]\ntrue -1 0 -1\n[7]\n"},
		{`import tripol.validate
show validate.integer("+42"),validate.integer("-"),validate.integer(""),validate.integer("1.2"),validate.integer("123"),validate.integer(" 1")
show validate.identifier("hello_2"),validate.identifier("2hello"),validate.identifier(""),validate.identifier("_"),validate.ascii("abc"),validate.ascii("é")
show validate.whitespace(9),validate.letter(90),validate.digit(48)`, "true false false false true false\ntrue false false true true false\ntrue true true\n"},
		{`import tripol.convert
show convert.celsius_to_fahrenheit(100.0),convert.fahrenheit_to_celsius(32.0),convert.km_to_miles(1.609344),convert.miles_to_km(1.0)
show convert.kg_to_lb(0.45359237),convert.lb_to_kg(1.0),convert.celsius_to_kelvin(0.0),convert.kelvin_to_celsius(273.15)`, "212.000000 0.000000 1.000000 1.609344\n1.000000 0.453592 273.150000 0.000000\n"},
		{`import tripol.encoding
show encoding.hex_byte(255),encoding.hex_int(0),encoding.hex_int(4096),encoding.hex_text("Aé"),encoding.hex_digit(10)`, "ff 0 1000 41c3a9 a\n"},
		{`import tripol.random
show random.next(1),random.next(48271),random.bounded(1,10),random.unit(1)`, "48271 182605794 0 0.000022\n"},
		{`import tripol.json
show json.pair("name","Ana\n\"\\")
a=[1,-2,3];show json.integers(a)
show json.integer(-5),json.boolean(true),json.boolean(false)
show json.quote("é")`, "{\"name\":\"Ana\\u000a\\\"\\\\\"}\n[1,-2,3]\n-5 true false\n\"é\"\n"},
		{`a=[1,2,3]
show a
show ints(0)
func get(n:int)->ints{return ints(n)}
show get(2)
show byte_at("abc",1),byte_at("é",0)`, "[1, 2, 3]\n[]\n[0, 0]\n98 195\n"},
	}
	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			runCase(t, c.src, c.want, false)
			p := filepath.Join(t.TempDir(), "debug")
			os.WriteFile(p, compileProfile(c.src, runtime.GOOS, "debug"), 0755)
			out, e := exec.Command(p).CombinedOutput()
			if e != nil || string(out) != c.want {
				t.Fatalf("debug %q %v", out, e)
			}
		})
	}
}
func TestEcosystemBounds(t *testing.T) {
	for _, s := range []string{`show byte_at("abc",-1)`, `show byte_at("abc",3)`, `show byte_at("",0)`} {
		t.Run(s, func(t *testing.T) { runCase(t, s, "TRIPOL runtime error: array index out of bounds\n", true) })
	}
}
func TestEcosystemNoAllocRead(t *testing.T) {
	runCase(t, `import tripol.validate
n=0
repeat 1000000 {if validate.integer("123") {n+=1}}
show n`, "1000000\n", false)
}
func TestEcosystemNoUnused(t *testing.T) {
	src := "show 42"
	imports := ""
	for _, n := range standardNames {
		imports += "import tripol." + n + "\n"
	}
	for _, target := range []string{"windows", "linux"} {
		if !bytes.Equal(compile(src, target), compile(imports+src, target)) {
			t.Fatal(target)
		}
	}
}
func mustReject(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected rejection")
		}
	}()
	f()
}
func TestVaultPackages(t *testing.T) {
	d := t.TempDir()
	src := filepath.Join(d, "original.tri")
	os.WriteFile(src, []byte("import tripol.math\nfunc double(n){return n*2}\nfunc result(n){return math.square(double(n))}"), 0644)
	addVault(d, "demo", src, "1.0.0")
	l := readLock(d)
	if len(l.Packages) != 1 {
		t.Fatal(l)
	}
	os.Remove(src)
	main := filepath.Join(d, "main.tri")
	os.WriteFile(main, []byte("import vault.demo\nshow demo.result(7)"), 0644)
	fs, ss, _ := loadProject(main)
	p := filepath.Join(d, "app")
	os.WriteFile(p, compileParsed(fs, ss, runtime.GOOS, "nano"), 0755)
	out, e := exec.Command(p).CombinedOutput()
	if e != nil || string(out) != "196\n" {
		t.Fatalf("%q %v", out, e)
	}
	os.WriteFile(vaultPath(d, l.Packages["demo"].SHA256), []byte("func result(n){return 1}"), 0644)
	mustReject(t, func() { loadProject(main) })
}
func TestVaultDependencies(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "source.tri")
	os.WriteFile(p, []byte("func square(n){return n*n}"), 0644)
	addVault(d, "foo", p, "1.0.0")
	os.WriteFile(p, []byte("import vault.foo\nfunc square(n){return foo.square(n)+1}"), 0644)
	addVault(d, "bar", p, "1.0.0")
	main := filepath.Join(d, "main.tri")
	os.WriteFile(main, []byte("import vault.bar\nshow bar.square(5)"), 0644)
	fs, ss, _ := loadProject(main)
	out := filepath.Join(d, "app")
	os.WriteFile(out, compileParsed(fs, ss, runtime.GOOS, "nano"), 0755)
	b, e := exec.Command(out).CombinedOutput()
	if e != nil || string(b) != "26\n" {
		t.Fatalf("%q %v", b, e)
	}
	os.WriteFile(p, []byte("import vault.bar\nfunc square(n){return bar.square(n)}"), 0644)
	mustReject(t, func() { addVault(d, "foo", p, "2.0.0") })
	if readLock(d).Packages["foo"].Version != "1.0.0" {
		t.Fatal("failed add altered lock")
	}
}
func TestVaultReject(t *testing.T) {
	for i, c := range []struct{ name, src, version string }{{"math", "func f(n){return n}", "1.0.0"}, {"../x", "func f(n){return n}", "1.0.0"}, {"foo", "show 1", "1.0.0"}, {"foo", "func f(n){return true}", "1.0.0"}, {"foo", "import vault.missing\nfunc f(n){return n}", "1.0.0"}, {"foo", "import \"../evil.tri\"\nfunc f(n){return n}", "1.0.0"}, {"foo", "func f(n){return n}", "latest"}} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			d := t.TempDir()
			p := filepath.Join(d, "p.tri")
			os.WriteFile(p, []byte(c.src), 0644)
			mustReject(t, func() { addVault(d, c.name, p, c.version) })
			if len(readLock(d).Packages) != 0 {
				t.Fatal("invalid lock")
			}
		})
	}
}
func TestVaultNoBorrow(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "p.tri")
	os.WriteFile(p, []byte("func identity(a:ints)->ints{return a}"), 0644)
	addVault(d, "foo", p, "1.0.0")
	main := filepath.Join(d, "main.tri")
	os.WriteFile(main, []byte("import vault.foo\na=[1,2]\nb=foo.identity(a)\nshow b"), 0644)
	fs, ss, _ := loadProject(main)
	mustReject(t, func() { compileParsed(fs, ss, "linux", "nano") })
}
func TestVaultLock(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "p.tri")
	os.WriteFile(p, []byte("func f(n){return n}"), 0644)
	addVault(d, "foo", p, "1.0.0")
	old, _ := os.ReadFile(filepath.Join(d, "tripol.lock"))
	addVault(d, "foo", p, "1.0.0")
	now, _ := os.ReadFile(filepath.Join(d, "tripol.lock"))
	if !bytes.Equal(old, now) {
		t.Fatal("non-deterministic lock")
	}
	os.WriteFile(p, []byte("func f(n){return n+1}"), 0644)
	mustReject(t, func() { addVault(d, "foo", p, "1.0.0") })
	addVault(d, "foo", p, "1.0.1")
	if readLock(d).Packages["foo"].Version != "1.0.1" {
		t.Fatal("update failed")
	}
}
func TestEcosystemJSON(t *testing.T) {
	src := `import tripol.json
show json.quote("\u0001\t\r\n\"\\é")`
	p := filepath.Join(t.TempDir(), "app")
	os.WriteFile(p, compile(src, runtime.GOOS), 0755)
	b, e := exec.Command(p).CombinedOutput()
	var s string
	if e != nil || json.Unmarshal(b, &s) != nil || s != "\x01\t\r\n\"\\é" {
		t.Fatalf("%q %v", b, e)
	}
}
func TestEcosystemSortLarge(t *testing.T) {
	runCase(t, `import tripol.sort
a=ints(4096)
for i in range(len(a)){a[i]=len(a)-i}
b=sort.ascending(a)
show sort.sorted(b),b[0],b[4095],a[0],sum(b)`, "true 1 4096 4096 8390656\n", false)
}
func TestVaultRejectQualifiedAccess(t *testing.T) { // keep built-in names distinct
	if strings.Contains(strings.Join(standardNames, ","), "vault") {
		t.Fatal("reserved namespace")
	}
}

func TestEcosystemTimeSystem(t *testing.T) {
	src := `import tripol.time
import tripol.system
start=time.now_ms()
time.sleep(15)
assert time.elapsed(start)>=10, "sleep elapsed too short"
show system.platform(),system.architecture(),system.compiler()
show time.seconds(1500),system.is_linux(),system.is_windows()`
	want := runtime.GOOS + " x86_64 0.11.0\n1.500000 "
	if runtime.GOOS == "linux" {
		want += "true false\n"
	} else {
		want += "false true\n"
	}
	runCase(t, src, want, false)
	// Cross-build using every optional OS service, then inspect all PE imports.
	full := `import tripol.time
import tripol.system
show time.now_ms();time.sleep(0);show system.platform()
x=input();write_text("x.txt",x);show read_text("x.txt")`
	b := compile(full, "windows")
	f, e := pecheck.NewFile(bytes.NewReader(b))
	if e != nil {
		t.Fatal(e)
	}
	im, e := f.ImportedSymbols()
	if e != nil {
		t.Fatal(e)
	}
	joined := strings.Join(im, " ")
	for _, n := range []string{"GetTickCount64", "Sleep", "CreateFileA", "VirtualAlloc", "ReadFile"} {
		if !strings.Contains(joined, n) {
			t.Fatal("missing", n, im)
		}
	}
}
func TestEcosystemTextNoAlloc(t *testing.T) {
	runCase(t, `import tripol.text
n=0
repeat 1000000 {if text.starts_with("TRIPOL","TRI") and text.ends_with("TRIPOL","POL") {n+=1}}
show n`, "1000000\n", false)
}
func TestEcosystemShowValues(t *testing.T) {
	runCase(t, `struct Data { a:ints }
x=Data(ints(2));show "a",x.a
n=[1,2];show "fixed",n,10`, "a [0, 0]\nfixed [1, 2] 10\n", false)
}
func TestVaultExplicitDependencies(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "p.tri")
	os.WriteFile(p, []byte("func f(n){return n}"), 0644)
	addVault(d, "foo", p, "1.0.0")
	os.WriteFile(p, []byte("func f(n){return foo.f(n)}"), 0644)
	mustReject(t, func() { addVault(d, "bar", p, "1.0.0") })
}
