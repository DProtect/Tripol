package core

import (
 "bytes"
 "context"
 elfcheck "debug/elf"
 "fmt"
 "os"
 "os/exec"
 "path/filepath"
 "runtime"
 "strconv"
 "strings"
 "testing"
 "time"
)

func TestQuantumXLists(t *testing.T){
 src:=`func append_to(a:list,n:int)->int {return a.add(n)}
 a=list(0);b=a
 for i in range(20) {a.add(i)}
 show b.size(),len(a),a.capacity(),b.get(19)
 show append_to(b,50),a.remove_at(0)
 show a.get(0),a.get(19),a.set(1,88)
 show a.get(1)
 try {a.get(-1)} catch e {show e}
 try {a.remove_at(100)} catch e {show e}
 b.clear();show a
 a.add(9);show b`
 for _,p:=range []string{"nano","balanced","turbo","debug"}{
  t.Run(p,func(t *testing.T){got,bad:=hyperRun(t,src,p);if bad||got!="20 20 32 19\n21 0\n1 50 88\n88\n11\n11\n[]\n[9]\n"{t.Fatalf("%q %v",got,bad)}})
 }
}
func TestQuantumXStreamRegex(t *testing.T){
 runCase(t,`import tripol.regex
import tripol.stream
data=[0,255,17]
c=stream.create()
show stream.read_byte(c,data),stream.read_u16(c,data),stream.remaining(c,data)
try {stream.read_byte(c,data)} catch e {show e}
show stream.seek(c,data,1),stream.position(c)
show regex.matches("a*b.","aaabx"),regex.matches("a*b.","b"),regex.matches("a*",""),regex.matches(".*","abc"),regex.valid("**")
`,"0 4607 0\n12\n1 1\ntrue false true true false\n",false)
}
func TestQuantumXBinaryIO(t *testing.T){
 if runtime.GOARCH!="amd64"||(runtime.GOOS!="windows"&&runtime.GOOS!="linux"){t.Skip("native host required")}
 dir:=t.TempDir()
 path:=filepath.Join(dir,"bytes.bin")
 src:=fmt.Sprintf(`import tripol.io
 p=%s
 data=ints(513)
 for i in range(len(data)) {data[i]=i%%256}
 show io.write(p,data)
 result=io.read(p,513)
 show len(result),result[0],result[255],result[512]
`,strconv.Quote(path))
 exe:=filepath.Join(dir,"io.exe");if e:=os.WriteFile(exe,compile(src,runtime.GOOS),0755);e!=nil{t.Fatal(e)}
 ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel()
 out,e:=exec.CommandContext(ctx,exe).CombinedOutput()
 if e!=nil||strings.ReplaceAll(string(out),"\r\n","\n")!="513\n513 0 255 0\n"{t.Fatalf("%q %v",out,e)}
 data,e:=os.ReadFile(path);if e!=nil||len(data)!=513{t.Fatal(len(data),e)}
 for i,b:=range data{if b!=byte(i%256){t.Fatal(i,b)}}
 // Oversized input must fail, never silently truncate.
 bad:=fmt.Sprintf(`import tripol.io
 show io.read(%s,512)`,strconv.Quote(path))
 exe2:=filepath.Join(dir,"bad.exe");if e=os.WriteFile(exe2,compile(bad,runtime.GOOS),0755);e!=nil{t.Fatal(e)}
 if e=exec.CommandContext(ctx,exe2).Run();e==nil{t.Fatal("oversized file accepted")}
}
func TestGenesisDirectFunctions(t *testing.T){
 if runtime.GOARCH!="amd64"||(runtime.GOOS!="linux"&&runtime.GOOS!="windows"){t.Skip("native host required")}
 fs,ss,_:=loadProject("compiler/direct.tri")
 dir:=t.TempDir();stage:=filepath.Join(dir,"stage.exe")
 if e:=os.WriteFile(stage,compileParsed(fs,ss,runtime.GOOS,"nano"),0755);e!=nil{t.Fatal(e)}
 for i,tc:=range []struct{source,want string}{
 {"a=0\nfor i in range(6) {if i==1 {continue};if i==4 {break};a+=i}\nwhile a>0 {a-=1;if a==1 {break};if a==4 {continue};show a}\nshow a\n","3\n2\n1\n"},
 {"func multiply(a:int,b:int)->int {return a*b}\nfor i in range(10) {show multiply(i,5)}","0\n5\n10\n15\n20\n25\n30\n35\n40\n45\n"},
 {"func sum(n:int)->int {if n<=0 {return 0};return n+sum(n-1)}\nshow sum(10)","55\n"},
 {"func twice(x:int)->int {return x*2}\nfunc plus(a:int,b:int)->int {return a+b}\nshow plus(twice(3),twice(4))","14\n"},
 {"s=0;for i in range(2,6) {s=s+i};for j in range(4,2) {s=s+100};show s","14\n"},
 {"func even(n:int)->bool {return n%2==0}\nfunc choose(flag:bool,a:int,b:int)->int {if flag {return a};return b}\nfunc both(a:bool,b:bool)->bool {return a and b}\nflag=even(8)\nother=not even(5)\nflag=both(flag,other)\nshow choose(flag,10,20)\nshow choose(false and 1/0==0,1,2)\nif flag==true {show 7}\nfor i in range(2) {if even(i) {show i}}\n","10\n2\n7\n0\n"},
 {"func combine(a:int,b:int)->int {return a*10+b}\nx=3\nx+=combine(1+2,4)\nx*=2\nx-=4\nx/=2\nx%=10\nshow x\nshow combine(2,3+combine(1,2))\nshow -(2+3)*4\n","5\n35\n-20\n"},
 {"func first(n:int)->int {\n return second(n)+1\n}\nfunc second(n:int)->int {\n if n==0 {return 0}\n return first(n-1)\n}\nx=2\nif false and 1/0==0 {show 99} else if (x+1)>2 and not (x==0 or false) {show first(3)} else {show 88}\nif true or 1/0==0 {show 42}\nwhile x>0 and true {show x;x=x-1}\n","4\n42\n2\n1\n"},
 }{
  t.Run(fmt.Sprintf("case_%02d",i),func(t *testing.T){
   t.Log(tc.source)
   ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second);defer cancel()
   dir:=t.TempDir()
   source:=filepath.Join(dir,"input.tri");out:=filepath.Join(dir,"output")
   if e:=os.WriteFile(source,[]byte(tc.source),0600);e!=nil{t.Fatal(e)}
   cmd:=exec.CommandContext(ctx,stage);cmd.Stdin=strings.NewReader(source+"\n"+out+"\n")
   if b,e:=cmd.CombinedOutput();e!=nil{t.Fatalf("stage: %q %v",b,e)}
   image,e:=os.ReadFile(out);if e!=nil{t.Fatal(e)}
   elf,e:=elfcheck.NewFile(bytes.NewReader(image));if e!=nil{t.Fatal(e)};defer elf.Close()
   if elf.Machine!=elfcheck.EM_X86_64||elf.Type!=elfcheck.ET_EXEC{t.Fatal("wrong ELF")}
   if runtime.GOOS=="linux"{
    if e=os.Chmod(out,0755);e!=nil{t.Fatal(e)}
    b,e:=exec.CommandContext(ctx,out).CombinedOutput();if e!=nil||string(b)!=tc.want{t.Fatalf("%q %v",b,e)}
   }
  })
 }
}
