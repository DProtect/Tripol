package core
import(
 "bytes"
 pecheck "debug/pe"
 "strings"
 "testing"
)
func TestQuantum2DTransforms(t *testing.T){
 runCase(t,`import tripol.sprite
 import tripol.canvas
 import tripol.particles
 a=[1,2,3,4,5,6]
 show sprite.rotate90(a,3,2)
 show sprite.scale_nearest(a,3,2,6,2)
 d=canvas.create(4,4)
 s=canvas.create(1,1);s[0]=16711680
 show sprite.draw_scaled(d,4,4,s,1,1,-1,1,3,2,-1,50)
 show d
 tiles=ints(4);tiles[0]=1;tiles[3]=2
 pal=ints(3);pal[1]=7;pal[2]=9
 canvas.clear(d,0)
 show canvas.tile_map(d,4,4,tiles,2,2,pal,2,0,0)
 show d
 p=particles.create(4);show particles.burst(p,2,2,4)
 show particles.draw(p,d,4,4,0,0,99)
 show d[10]
 show particles.step(p)`,"[4, 1, 5, 2, 6, 3]\n[1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6]\n4\n[0, 0, 0, 0, 8323072, 8323072, 0, 0, 8323072, 8323072, 0, 0, 0, 0, 0, 0]\n2\n[7, 7, 0, 0, 7, 7, 0, 0, 0, 0, 9, 9, 0, 0, 9, 9]\n4\n4\n99\n4\n",false)
}
func TestQuantumMouseAndDemo(t *testing.T){
 runCase(t,`import tripol.events
 show events.mouse_x(320),events.mouse_y(180)`,"-1 -1\n",false)
 src:=`import tripol.audio
 import tripol.events
 show events.mouse_x(320)
 show audio.beep(440,20)`
 file,e:=pecheck.NewFile(bytes.NewReader(compile(src,"windows")));if e!=nil{t.Fatal(e)};defer file.Close()
 syms,e:=file.ImportedSymbols();if e!=nil{t.Fatal(e)}
 for _,name:=range []string{"GetCursorPos","ScreenToClient","Beep"}{if !strings.Contains(strings.Join(syms," "),name){t.Fatal(name,syms)}}
 for _,p:=range []string{"examples/projects/quantum_x/main.tri","examples/projects/quantum_x/opengl.tri"}{
  fs,ss,_:=loadProject(p);for _,profile:=range []string{"nano","balanced","turbo","debug"}{if len(compileParsed(fs,ss,"windows",profile))==0{t.Fatal(p,profile)}}
 }
}
