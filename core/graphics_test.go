package core

import (
 "bytes"
 pecheck "debug/pe"
 "runtime"
 "strings"
 "testing"
)

func TestQuantumGraphicsImports(t *testing.T){
 src:=`import tripol.game
 import tripol.events
 a=[-800,-800,16711680,800,-800,65280,0,800,255]
 show game.gl_frame(a,0)
 show events.pad_connected()`
 file,e:=pecheck.NewFile(bytes.NewReader(compile(src,"windows")));if e!=nil{t.Fatal(e)};defer file.Close()
 syms,e:=file.ImportedSymbols();if e!=nil{t.Fatal(e)}
 for _,name:=range []string{"wglCreateContext","wglMakeCurrent","wglDeleteContext","glVertex2f","SwapBuffers","XInputGetState"}{if !strings.Contains(strings.Join(syms," "),name){t.Fatal(name,syms)}}
 // No window opened: valid vertices return false. Invalid vertices trap
 // before calling Windows, independent of graphics hardware.
 if runtime.GOOS=="windows"{runCase(t,`import tripol.game
 a=[0,0,0,100,100,1,200,200,2]
 show game.gl_frame(a)
 b=[1001,0,0,0,0,0,0,0,0]
 try {game.gl_frame(b)} catch e {show e}
 c=ints(1)
 try {game.gl_frame(c)} catch e {show e}`,"false\n11\n11\n",false)}
 if runtime.GOOS=="linux"{runCase(t,src,"false\nfalse\n",false)}
}
