package main

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

func TestEvolutionCollections(t *testing.T) {
 src := `import tripol.stack
import tripol.queue
import tripol.memory
s=stack.create(2);q=queue.create(2)
before=memory.used()
show stack.push(s,10),stack.push(s,20),stack.push(s,30)
show stack.peek(s),stack.pop(s),stack.pop(s)
show queue.push(q,10),queue.push(q,20),queue.push(q,30)
show queue.pop(q),queue.push(q,30),queue.peek(q),queue.pop(q),queue.pop(q)
repeat 1000 {assert stack.push(s,1);assert stack.pop(s)==1;assert queue.push(q,2);assert queue.pop(q)==2}
show before==memory.used(),stack.size(s),queue.size(q)
try {stack.pop(s)} catch e {show e}
try {queue.pop(q)} catch e {show e}`
 for _, profile := range []string{"nano","balanced","turbo","debug"} {
  t.Run(profile,func(t *testing.T){
   got,bad:=hyperRun(t,src,profile)
   want:="true true false\n20 20 10\ntrue true false\n10 true 20 20 30\ntrue 0 0\n12\n12\n"
   if bad || got!=want {t.Fatalf("%q %v",got,bad)}
  })
 }
}

func TestEvolutionTIRIdentity(t *testing.T) {
 src := `func add(n:int)->int {return 0+(n+1)}
func mul(n:int)->int {return 1*(n*2)}
show add(41),mul(-4)
try {show add(9223372036854775807)} catch e {show e}
try {show mul(9223372036854775807)} catch e {show e}`
 for _,profile:=range []string{"nano","balanced","turbo","debug"} {
  t.Run(profile,func(t *testing.T){
   got,bad:=hyperRun(t,src,profile)
   if bad || got!="42 -8\n10\n10\n" {t.Fatalf("%q %v",got,bad)}
  })
 }
}

func TestGenesisNativeELF(t *testing.T) {
 if runtime.GOARCH!="amd64" || (runtime.GOOS!="linux" && runtime.GOOS!="windows") {t.Skip("native x86-64 required")}
 fs,ss,_:=loadProject("compiler/main.tri")
 stage:=filepath.Join(t.TempDir(),"genesis.exe")
 if err:=os.WriteFile(stage,compileParsed(fs,ss,runtime.GOOS,"nano"),0755);err!=nil {t.Fatal(err)}
 run:=func(input string)([]byte,error){
  ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel()
  cmd:=exec.CommandContext(ctx,stage);cmd.Stdin=strings.NewReader(input+"\n")
  return cmd.CombinedOutput()
 }
 for _,tc:=range []struct{src,want string}{{"show (6+1)*6","42\n"},{"show -7+2","-5\n"},{"show 17%5","2\n"},{"show 0","0\n"}} {
  t.Run(tc.src,func(t *testing.T){
   out,err:=run(tc.src);if err!=nil {t.Fatalf("GENESIS: %q %v",out,err)}
   image,err:=hex.DecodeString(strings.TrimSpace(string(out)));if err!=nil {t.Fatal(err)}
   out2,err:=run(tc.src);if err!=nil||!bytes.Equal(out,out2){t.Fatal("non-deterministic GENESIS output")}
   file,err:=elfcheck.NewFile(bytes.NewReader(image));if err!=nil {t.Fatal(err)}
   defer file.Close()
   if file.Machine!=elfcheck.EM_X86_64||file.Type!=elfcheck.ET_EXEC||file.Entry!=0x400078||len(file.Progs)!=1 {t.Fatal("unexpected ELF metadata")}
   if !bytes.HasSuffix(image,[]byte(tc.want)){t.Fatal("wrong emitted message")}
   if runtime.GOOS=="linux" {
    p:=filepath.Join(t.TempDir(),"program")
    if err:=os.WriteFile(p,image,0755);err!=nil {t.Fatal(err)}
    ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second);defer cancel()
    output,err:=exec.CommandContext(ctx,p).CombinedOutput()
    if err!=nil||string(output)!=tc.want {t.Fatalf("ELF run: %q %v",output,err)}
   }
  })
 }
 for _,bad:=range []string{"print 2","show 1/0","show 2+","show x","show 9223372036854775807+1"}{
  if out,err:=run(bad);err==nil {t.Fatalf("accepted invalid source %q: %q",bad,out)}
 }
}

func TestGenesisRuntimeR2(t *testing.T) {
 if runtime.GOARCH!="amd64" || (runtime.GOOS!="linux" && runtime.GOOS!="windows") {t.Skip("native x86-64 required")}
 fs,ss,_:=loadProject("compiler/main.tri")
 stage:=filepath.Join(t.TempDir(),"genesis.exe")
 if e:=os.WriteFile(stage,compileParsed(fs,ss,runtime.GOOS,"nano"),0755);e!=nil{t.Fatal(e)}
 for _,tc:=range []struct{src,want string;status int}{
 {"a=20;b=30;show a*b","600\n",0},
 {"a=7;b=10;a=a-b;show a*2","-6\n",0},
 {"a=0;show a","0\n",0},
 {"a=9223372036854775807;show -a-1","-9223372036854775808\n",0},
 {"a=9223372036854775807;show a+1","",10},
 }{
  t.Run(tc.src,func(t *testing.T){
   ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel()
   cmd:=exec.CommandContext(ctx,stage);cmd.Stdin=strings.NewReader(tc.src+"\n")
   out,e:=cmd.CombinedOutput();if e!=nil{t.Fatalf("stage: %q %v",out,e)}
   image,e:=hex.DecodeString(strings.TrimSpace(string(out)));if e!=nil{t.Fatal(e)}
   file,e:=elfcheck.NewFile(bytes.NewReader(image));if e!=nil{t.Fatal(e)};defer file.Close()
   if file.Entry!=0x400078||file.Machine!=elfcheck.EM_X86_64{t.Fatal("invalid ELF")}
   if runtime.GOOS!="linux"{return}
   p:=filepath.Join(t.TempDir(),"generated")
   if e=os.WriteFile(p,image,0755);e!=nil{t.Fatal(e)}
   got,e:=exec.CommandContext(ctx,p).CombinedOutput()
   if tc.status==0 {if e!=nil||string(got)!=tc.want{t.Fatalf("%q %v",got,e)}}else{exit,ok:=e.(*exec.ExitError);if !ok||exit.ExitCode()!=tc.status{t.Fatalf("expected exit %d: %q %v",tc.status,got,e)}}
  })
 }
}
func TestOverdriveBitset(t *testing.T){
 runCase(t,`import tripol.bitset
import tripol.profiler
b=bitset.create(130)
show bitset.add(b,0),bitset.add(b,63),bitset.add(b,129),bitset.add(b,63)
show bitset.count(b),bitset.contains(b,129),bitset.remove(b,63),bitset.count(b)
try {bitset.add(b,130)} catch e {show e}
start=profiler.start()
show profiler.elapsed(start)>=0`,"true true true false\n3 true true 2\n12\ntrue\n",false)
}
