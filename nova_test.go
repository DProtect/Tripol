package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNovaHD(t *testing.T) {
	runCase(t, `import tripol.canvas
import tripol.memory
p=canvas.create(1920,1080)
show len(p),memory.used()
show canvas.clear(p,1)
show canvas.rectangle(p,1920,1080,-2,1078,5,4,9)
show canvas.get(p,1920,1080,0,1079),canvas.get(p,1920,1080,3,1079),canvas.get(p,1920,1080,1919,1079)
show array_fill_range(p,2073599,1,7),p[2073599]
try {array_fill_range(p,2073600,1,0)} catch e {show e}
try {array_fill_range(p,-1,1,0)} catch e {show e}
try {array_fill_range(p,0,-1,0)} catch e {show e}
show array_fill_range(p,2073600,0,0)`, "2073600 16588808\n2073600\n6\n9 1 1\n1 7\n11\n11\n11\n0\n", false)
}
func TestNovaDrawing(t *testing.T) {
	runCase(t, `import tripol.canvas
import tripol.sprite
p=canvas.create(5,5)
show canvas.line(p,5,5,0,0,4,4,10)
show canvas.circle(p,5,5,2,2,1,20)
show canvas.get(p,5,5,2,1),canvas.get(p,5,5,0,0)
s=[0,3,4,5];show sprite.draw_key(p,5,5,s,2,2,3,3,0)
show canvas.get(p,5,5,3,3),canvas.get(p,5,5,4,3),canvas.get(p,5,5,3,4),canvas.get(p,5,5,4,4)
show canvas.blend(p,5,5,0,1,16777215,50),canvas.get(p,5,5,0,1)
show canvas.line(p,5,5,-2,0,2,0,42)
show sprite.draw(p,5,5,s,2,2,-1,-1),p[0]`, "5\n5\n20 10\n3\n10 3 4 5\ntrue 8355711\n3\n1 5\n", false)
}
func TestNovaDictionary(t *testing.T) {
	runCase(t, `import tripol.dictionary
m=dictionary.create(2)
show dictionary.set(m,"vida",100),dictionary.set(m,"nivel",5)
show dictionary.set(m,"vida",80),dictionary.get(m,"vida"),dictionary.size(m)
show dictionary.set(m,"otro",7),dictionary.contains(m,"otro")
show dictionary.remove(m,"nivel"),dictionary.remove(m,"nivel")
show dictionary.set(m,"otro",7),dictionary.get(m,"nivel",-1),dictionary.size(m)
try {dictionary.set(m,"",1)} catch e {show e}`, "true true\ntrue 80 2\nfalse false\ntrue false\ntrue -1 2\n12\n", false)
}
func TestNovaParticles(t *testing.T) {
	runCase(t, `import tripol.particles
import tripol.camera
import tripol.testing
p=particles.create(1);p[0]=1;p[1]=2;p[2]=3;p[3]=-1;p[4]=2
show particles.step(p),p
show particles.step(p),p
show particles.step(p),p
show camera.screen(10,2,2),camera.world(16,2,2),camera.follow(90,20,100),camera.visible(-2,-2,3,3,5,5)
show testing.equal_int(2,2),testing.equal_text("a","a"),testing.equal_bool(true,true),testing.near(1.0,1.1,0.2)`, "1 [4, 1, 3, -1, 1]\n0 [7, 0, 3, -1, 0]\n0 [7, 0, 3, -1, 0]\n16 10 80 true\ntrue true true true\n", false)
}
func TestNovaTurbo(t *testing.T) {
	s := `func calc(a,b,c,d) => (a+b)*(c-d)
func times(a,b) => a*b
show calc(2,3,9,2),times(-8,7)
try {show times(9223372036854775807,2)} catch e {show e}`
	for _, profile := range []string{"nano", "turbo", "debug"} {
		if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
			continue
		}
		p := filepath.Join(t.TempDir(), "prog")
		if e := os.WriteFile(p, compileProfile(s, "linux", profile), 0755); e != nil {
			t.Fatal(e)
		}
		out, e := exec.Command(p).CombinedOutput()
		if e != nil || string(out) != "35 -56\n10\n" {
			t.Fatalf("%s %q %v", profile, out, e)
		}
	}
}

func TestNovaParser(t *testing.T) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		t.Skip("Linux native execution")
	}
	fs, ss, _ := loadProject("bootstrap/parser_demo.tri")
	p := filepath.Join(t.TempDir(), "parser")
	if e := os.WriteFile(p, compileParsed(fs, ss, "linux", "nano"), 0755); e != nil {
		t.Fatal(e)
	}
	out, e := exec.Command(p).CombinedOutput()
	if e != nil || string(out) != "14 20 -4\n12\n10\n" {
		t.Fatalf("%q %v", out, e)
	}
}
