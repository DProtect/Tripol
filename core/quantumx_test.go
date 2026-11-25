package core

import (
 "bytes"
 "context"
 elfcheck "debug/elf"
 "encoding/hex"
 "os"
 "os/exec"
 "path/filepath"
 "runtime"
 "strings"
 "testing"
 "time"
)

func TestQuantumXCompoundArray(t *testing.T) {
 src:=`func index(counter:ints)->int {counter[0]=counter[0]+1;return 0}
 fixed=[10,20];a=ints(2);a[0]=10;counter=ints(1)
 fixed[0]+=5;fixed[1]/=4
 a[index(counter)]+=5
 a[0]*=2;a[0]-=3;a[0]%=10
 show fixed[0],fixed[1],a[0],counter[0]
 try {a[index(counter)]/=0} catch e {show e}
 show a[0],counter[0]
 try {a[5]+=1} catch e {show e}
 a[0]=9223372036854775807
 try {a[0]+=1} catch e {show e}
 show a[0]`
 for _,p:=range []string{"nano","balanced","turbo","debug"}{
  t.Run(p,func(t *testing.T){got,bad:=hyperRun(t,src,p);if bad||got!="15 5 7 1\n10\n7 2\n11\n10\n9223372036854775807\n"{t.Fatalf("%q %v",got,bad)}})
 }
}
func TestQuantumXLibraries(t *testing.T) {
 runCase(t,`import tripol.heap
import tripol.utf8
import tripol.csv
import tripol.memory
h=heap.create(4)
before=memory.used()
show heap.push(h,9),heap.push(h,-2),heap.push(h,4),heap.push(h,0),heap.push(h,8)
show heap.peek(h),heap.pop(h),heap.pop(h),heap.pop(h),heap.pop(h)
show before==memory.used()
try {heap.pop(h)} catch e {show e}
show utf8.valid("Aé中"),utf8.count("Aé中"),utf8.scalar_at("Aé中",1)
show utf8.width_at("Aé中",2),utf8.scalar_at("Aé中",3)
show csv.fields("a,\"b,c\",d"),csv.valid("a,\"b\"x"),csv.escape("a\"b")
show csv.fields("a,,b,"),csv.fields("\"unfinished")`,"true true true true false\n-2 -2 0 4 9\ntrue\n12\ntrue 3 233\n0 20013\n3 false \"a\"\"b\"\n4 -1\n",false)
}
func TestGenesisQuantumX(t *testing.T) {
 if runtime.GOARCH!="amd64"||(runtime.GOOS!="linux"&&runtime.GOOS!="windows"){t.Skip("native x86-64 required")}
 fs,ss,_:=loadProject("compiler/file.tri")
 stage:=filepath.Join(t.TempDir(),"genesis.exe")
 if e:=os.WriteFile(stage,compileParsed(fs,ss,runtime.GOOS,"nano"),0755);e!=nil{t.Fatal(e)}
 for _,tc:=range []struct{src,want string;status int}{
 {"i=0\nsum=0\nwhile i<10 {\nsum=sum+i\ni=i+1\n}\nif sum==45 {sum=sum*2} else {sum=0}\nshow sum/3\n","30\n",0},
 {"x=9;if x<2 {x=0} else {x=x%4};show x","1\n",0},
 {"x=0;while x<0 {x=x+1};show x","0\n",0},
 {"a=7;b=0;show a/b","",10},
 {"a=9223372036854775807;b=-a-1;show b/-1","",10},
 {"a=-7;show a%3","-1\n",0},
 }{
  t.Run(tc.src,func(t *testing.T){
   ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second);defer cancel()
   dir:=t.TempDir();source:=filepath.Join(dir,"source.tri")
   if e:=os.WriteFile(source,[]byte(tc.src),0600);e!=nil{t.Fatal(e)}
   emit:=func()[]byte{cmd:=exec.CommandContext(ctx,stage);cmd.Stdin=strings.NewReader(source+"\n");out,e:=cmd.CombinedOutput();if e!=nil{t.Fatalf("stage: %q %v",out,e)};image,e:=hex.DecodeString(strings.TrimSpace(string(out)));if e!=nil{t.Fatal(e)};return image}
   image:=emit();if !bytes.Equal(image,emit()){t.Fatal("non-deterministic ELF")}
   elf,e:=elfcheck.NewFile(bytes.NewReader(image));if e!=nil{t.Fatal(e)};defer elf.Close()
   if elf.Entry!=0x400078||elf.Machine!=elfcheck.EM_X86_64{t.Fatal("invalid ELF")}
   if runtime.GOOS!="linux"{return}
   exe:=filepath.Join(dir,"native");if e=os.WriteFile(exe,image,0755);e!=nil{t.Fatal(e)}
   out,e:=exec.CommandContext(ctx,exe).CombinedOutput()
   if tc.status==0 {if e!=nil||string(out)!=tc.want{t.Fatalf("%q %v",out,e)}}else{exit,ok:=e.(*exec.ExitError);if !ok||exit.ExitCode()!=tc.status{t.Fatalf("expected %d: %q %v",tc.status,out,e)}}
  })
 }
 // Declarations inside conditional blocks cannot escape or remain uninitialized.
 dir:=t.TempDir();source:=filepath.Join(dir,"invalid.tri")
 if e:=os.WriteFile(source,[]byte("a=1;if a==0 {b=3};show b"),0600);e!=nil{t.Fatal(e)}
 ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel()
 cmd:=exec.CommandContext(ctx,stage);cmd.Stdin=strings.NewReader(source+"\n");if e:=cmd.Run();e==nil||ctx.Err()!=nil{t.Fatalf("invalid source result: %v",e)}
}
