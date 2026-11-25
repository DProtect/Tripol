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

func TestAtlasLibraries(t *testing.T) {
	cases := []struct{ src, want string }{
		{`import tripol.math
show math.square(7), math.cube(3), math.power(2,10), math.gcd(36,60), math.clamp(5,0,4)
show math.power(0,0), math.root(81.0), math.percent(200.0,15.0), math.pi()
show math.absolute(-3), math.minimum(2,4), math.maximum(2,4), math.sqrt_value(16.0)`, "49 27 1024 12 4\n1 9.000000 30.000000 3.141593\n3 2 4 4.000000\n"},
		{`import tripol.stats
n=[10,20,30,40]
show stats.mean(n),stats.minimum(n),stats.maximum(n),stats.variance(n),stats.deviation(n),stats.total(n)`, "25.000000 10 40 125.000000 11.180340 100\n"},
		{`import tripol.arrays
n=[1,2,2,3]
show arrays.total(n),arrays.count(n,2),arrays.contains(n,3),arrays.contains(n,4),arrays.length(n)
b=arrays.reversed(n); for v in b {show v}
c=arrays.copy(n); c[0]=10; show n[0], c[0]
d=arrays.make(3,7);show arrays.total(d)`, "8 2 true false 4\n3\n2\n2\n1\n1 10\n21\n"},
		{`import tripol.text
show text.size("abc"), text.empty(""), text.from_int(42), text.to_int("-42"), text.part("abc",1,3)
show text.starts_with("abc","ab"),text.starts_with("a","ab"),text.ends_with("abc","bc"),text.ends_with("a","ab"),text.ends_with("abc","")`, "3 true 42 -42 bc\ntrue false true false true\n"},
		{`import tripol.console
console.line("Hello")`, "Hello\n"},
		{`import tripol.arrays
import tripol.stats
import tripol.stats
n=[1,2,3]
// More than 8 MiB of copies would be needed if this borrowed array allocated.
result=0
repeat 400000 { result+=arrays.total(n) }
show result`, "2400000\n"},
		{`import tripol.math
import tripol.text
func square(n){return 999}
func len(s:string)->int{return 999}
show square(2),math.square(2),text.size("abc")`, "999 4 3\n"},
		{`import tripol.arrays
func go(n:int)->ints { a=[1,2,3]; return arrays.copy(a) }
x=go(0);show sum(x)`, "6\n"},
	}
	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			runCase(t, c.src, c.want, false)
			p := filepath.Join(t.TempDir(), "debug.exe")
			if e := os.WriteFile(p, compileProfile(c.src, runtime.GOOS, "debug"), 0755); e != nil {
				t.Fatal(e)
			}
			out, e := exec.Command(p).CombinedOutput()
			if e != nil || string(out) != c.want {
				t.Fatalf("debug %q %v", out, e)
			}
		})
	}
}
func TestAtlasFilesConsole(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "a.exe")
	src := `import tripol.files
import tripol.console
name=console.ask("Name?")
age=console.ask_int("Age?")
files.write_line("note.txt",name)
files.append_line("note.txt",text(age))
show files.read("note.txt")`
	if e := os.WriteFile(p, compile(src, runtime.GOOS), 0755); e != nil {
		t.Fatal(e)
	}
	cmd := exec.Command(p)
	cmd.Dir = d
	cmd.Stdin = strings.NewReader("Ana\n16\n")
	out, e := cmd.CombinedOutput()
	if e != nil || string(out) != "Name?\nAge?\nAna\n16\n\n" {
		t.Fatalf("%q %v", out, e)
	}
	b, e := os.ReadFile(filepath.Join(d, "note.txt"))
	if e != nil || string(b) != "Ana\n16\n" {
		t.Fatalf("%q %v", b, e)
	}
}
func TestAtlasReject(t *testing.T) {
	for _, s := range []string{`import tripol.missing`, `show math.square(3)`, "import tripol.math\nshow math.square(true)", "import tripol.text\nshow text.unknown(1)", "import other.math\nshow 1", "import tripol.math\nshow math.square()", "import tripol.math\nshow math.square(1,2)"} {
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
func TestAtlasTraps(t *testing.T) {
	runCase(t, "import tripol.math\nshow math.power(2,-1)", "TRIPOL assertion failed at line 8: Negative integer exponent\n", true)
	runCase(t, "import tripol.stats\na=ints(0);show stats.mean(a)", "TRIPOL assertion failed at line 3: Mean requires nonempty array\n", true)
}
func TestAtlasElimination(t *testing.T) {
	s := "show \"Hello\""
	imp := "import tripol.math\nimport tripol.stats\nimport tripol.text\nimport tripol.files\nimport tripol.console\n" + s
	for _, target := range []string{"windows", "linux"} {
		if !bytes.Equal(compile(s, target), compile(imp, target)) {
			t.Fatalf("unused libraries change %s binary", target)
		}
	}
	b := compile("import tripol.math\nshow math.square(7)", "windows")
	f, e := pecheck.NewFile(bytes.NewReader(b))
	if e != nil {
		t.Fatal(e)
	}
	if _, e = f.ImportedSymbols(); e != nil {
		t.Fatal(e)
	}
}
func TestAtlasMixedModules(t *testing.T) {
	d := t.TempDir()
	main := filepath.Join(d, "main.tri")
	os.WriteFile(main, []byte("import tripol.math\nimport \"util.tri\"\nshow twice(5),math.square(3)"), 0644)
	os.WriteFile(filepath.Join(d, "util.tri"), []byte("import tripol.math\nfunc twice(n){return math.square(n)+math.square(n)}"), 0644)
	fs, ss, paths := loadProject(main)
	if len(paths) != 2 {
		t.Fatal(paths)
	}
	p := filepath.Join(d, "app.exe")
	os.WriteFile(p, compileParsed(fs, ss, runtime.GOOS, "nano"), 0755)
	out, e := exec.Command(p).CombinedOutput()
	if e != nil || string(out) != "50 9\n" {
		t.Fatalf("%q %v", out, e)
	}
}
