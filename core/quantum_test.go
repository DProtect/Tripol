package core

import (
	"bytes"
	pecheck "debug/pe"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestQuantumMemory(t *testing.T) {
	runCase(t, `import tripol.memory
show memory.used(),memory.committed(),memory.capacity()
a=ints(65535)
a[0]=77
a[65534]=99
b=a
repeat 31 {b=ints(65535);b[65534]=1}
show memory.used(),memory.committed()
show a[0],a[65534]
before=memory.used()
repeat 100 {scratch {temp=ints(65535);temp[65534]=8}}
show memory.used()==before,memory.committed()>=before
try {scratch {temp=ints(65535);throw 101}} catch e {show e}
show memory.used()==before`, "0 0 268435456\n16777216 16777216\n77 99\ntrue true\n101\ntrue\n", false)
}
func TestQuantumCanvas(t *testing.T) {
	runCase(t, `import tripol.canvas
import tripol.image
p=canvas.create(4,3)
show canvas.clear(p,0)
show canvas.rectangle(p,4,3,-1,1,3,4,7)
show canvas.pixel(p,4,3,3,0,9)
show canvas.pixel(p,4,3,4,0,9)
show p
show image.flip_x(p,4,3)
show image.flip_y(p,4,3)
show image.crop(p,4,3,0,1,2,2)
show canvas.get(p,4,3,3,0)
gray=[16777215,0,16711680];show image.grayscale(gray)
two=[16711680,255];show image.ppm(two,2,1)`, "12\n4\ntrue\nfalse\n[0, 0, 0, 9, 7, 7, 0, 0, 7, 7, 0, 0]\n[9, 0, 0, 0, 0, 0, 7, 7, 0, 0, 7, 7]\n[7, 7, 0, 0, 7, 7, 0, 0, 0, 0, 0, 9]\n[7, 7, 7, 7]\n9\n[16777215, 0, 5000268]\nP3\n2 1\n255\n255 0 0\n0 0 255\n\n", false)
}
func TestQuantumMath(t *testing.T) {
	runCase(t, `import tripol.physics
import tripol.animation
import tripol.vector
import tripol.path
import tripol.events
show physics.overlap(0,0,5,5,4,4,2,2),physics.overlap(0,0,5,5,5,0,1,1)
show physics.clamp_position(90,10,50),physics.contains(0,0,5,5,4,4)
show animation.linear(0.0,10.0,0.25),animation.smoothstep(0.5),animation.frame(190,50,3)
v=[1,2,3];show vector.length2(3.0,4.0),vector.cross2(1.0,2.0,3.0,4.0),vector.scale(v,2)
show path.basename("a/b.txt"),path.dirname("a/b.txt"),path.extension("a/b.txt"),path.join("a","b")
show events.left(),events.escape()`, "true false\n40 true\n2.500000 0.500000 0\n5.000000 -2.000000 [2, 4, 6]\nb.txt a .txt a/b\n37 27\n", false)
}
func TestQuantumReject(t *testing.T) {
	for _, s := range []string{`import tripol.canvas
const a=[1,2];canvas.clear(a,3)`, `func bad()->bool => 42
show bad()`, `func f(n) => n+1 show f(1)`} {
		t.Run(s, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("accepted invalid program")
				}
			}()
			compile(s, "linux")
		})
	}
	runCase(t, `import tripol.canvas
try {p=canvas.create(2049,1)} catch e {show e}
try {p=canvas.create(2,2);canvas.clear(p,-1)} catch e {show e}`, "12\n12\n", false)
}
func TestQuantumPE(t *testing.T) {
	p, e := pecheck.NewFile(bytes.NewReader(compile(quantumGameExample, "windows")))
	if e != nil {
		t.Fatal(e)
	}
	syms, e := p.ImportedSymbols()
	if e != nil {
		t.Fatal(e)
	}
	all := strings.Join(syms, " ")
	for _, n := range []string{"CreateWindowExA:USER32.dll", "PeekMessageA:USER32.dll", "GetAsyncKeyState:USER32.dll", "StretchDIBits:GDI32.dll", "VirtualFree:KERNEL32.dll", "VirtualAlloc:KERNEL32.dll", "GetModuleHandleA:KERNEL32.dll"} {
		if !strings.Contains(all, n) {
			t.Fatal(n, all)
		}
	}
	if runtime.GOOS == "linux" {
		runCase(t, `import tripol.game
show game.window("Test",320,180)`, "TRIPOL runtime error: game window requires Windows or graphics operation failed\n", true)
	}
}
func TestQuantumTurbo(t *testing.T) {
	s := `func add(a,b) => a+b
func calc(a,b,c,d) => (a+b)*(c-d)
show add(4,5),calc(2,3,9,2)
try {show add(9223372036854775807,1)} catch e {show e}`
	for _, profile := range []string{"nano", "turbo", "debug"} {
		if runtime.GOARCH != "amd64" || runtime.GOOS != "linux" {
			continue
		}
		p := filepath.Join(t.TempDir(), "prog")
		if e := os.WriteFile(p, compileProfile(s, "linux", profile), 0755); e != nil {
			t.Fatal(e)
		}
		out, e := exec.Command(p).CombinedOutput()
		if e != nil || string(out) != "9 35\n10\n" {
			t.Fatalf("%s: %s %v", profile, out, e)
		}
	}
}

func TestQuantumBuffer(t *testing.T) {
	runCase(t, `import tripol.buffer
b=buffer.create(0)
for i in range(1000) {b=buffer.push(b,i)}
show buffer.size(b),buffer.capacity(b),buffer.get(b,999)
show buffer.set(b,0,42),buffer.pop(b),buffer.size(b),buffer.get(b,0)
small=buffer.create();small=buffer.push(small,7);small=buffer.push(small,8)
show buffer.values(small)
try {show buffer.get(small,2)} catch e {show e}`, "1000 1024 999\n42 999 999 42\n[7, 8]\n12\n", false)
}

func TestQuantumFixPreview(t *testing.T) {
	d := diagnostic{Line: 2, Message: "unknown variable numro", Suggestion: "numero"}
	a, b, ok := replacementPreview(d, "numero=1\nshow numro // numro")
	if !ok || a != "show numro // numro" || b != "show numero // numro" {
		t.Fatal(a, b, ok)
	}
	_, _, ok = replacementPreview(d, "numero=1\nshow numro,numro")
	if ok {
		t.Fatal("ambiguous replacement")
	}
}

func TestQuantumBootstrap(t *testing.T) {
	fs, ss, _ := loadProject("compiler/demo.tri")
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		t.Skip("Linux execution test")
	}
	p := filepath.Join(t.TempDir(), "bootstrap")
	if e := os.WriteFile(p, compileParsed(fs, ss, "linux", "nano"), 0755); e != nil {
		t.Fatal(e)
	}
	out, e := exec.Command(p).CombinedOutput()
	want := "1 func\n1 doble\n4 (\n1 n\n4 )\n4 =>\n1 n\n4 *\n2 2\n5 \n\n1 show\n1 doble\n4 (\n2 5\n4 )\n"
	if e != nil || string(out) != want {
		t.Fatalf("%q %v", out, e)
	}
}
