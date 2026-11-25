package core

import (
 "bytes"
 "context"
 pecheck "debug/pe"
 "os"
 "os/exec"
 "path/filepath"
 "runtime"
 "strings"
 "testing"
 "time"
)
func TestGenesisWindowsNative(t *testing.T){
 if runtime.GOARCH!="amd64"||(runtime.GOOS!="linux"&&runtime.GOOS!="windows"){t.Skip("native host required")}
 fs,ss,_:=loadProject("compiler/windows.tri")
 stage:=filepath.Join(t.TempDir(),"stage.exe")
 if e:=os.WriteFile(stage,compileParsed(fs,ss,runtime.GOOS,"nano"),0755);e!=nil{t.Fatal(e)}
 for _,tc:=range []struct{src,want string;status int}{
 {"func multiply(a:int,b:int)->int {return a*b}\nfor i in range(3) {show multiply(i,5)}","0\n5\n10\n",0},
 {"a=9223372036854775807;show -a-1","-9223372036854775808\n",0},
 {"a=7;show a/0","",10},
 {"a=9223372036854775807;show a+1","",10},
 {"a=-7;show a%3","-1\n",0},
 }{
  t.Run(tc.src,func(t *testing.T){
   dir:=t.TempDir();src:=filepath.Join(dir,"input.tri");out:=filepath.Join(dir,"genesis.exe")
   if e:=os.WriteFile(src,[]byte(tc.src),0600);e!=nil{t.Fatal(e)}
   ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second);defer cancel()
   emit:=func()[]byte{
    cmd:=exec.CommandContext(ctx,stage);cmd.Stdin=strings.NewReader(src+"\n"+out+"\n")
    if b,e:=cmd.CombinedOutput();e!=nil{t.Fatalf("stage: %q %v",b,e)}
    b,e:=os.ReadFile(out);if e!=nil{t.Fatal(e)};return b
   }
   image:=emit();if !bytes.Equal(image,emit()){t.Fatal("non-deterministic PE")}
   file,e:=pecheck.NewFile(bytes.NewReader(image));if e!=nil{t.Fatal(e)};defer file.Close()
   header,ok:=file.OptionalHeader.(*pecheck.OptionalHeader64);if !ok||file.Machine!=pecheck.IMAGE_FILE_MACHINE_AMD64||header.AddressOfEntryPoint!=4096||len(file.Sections)!=2{t.Fatal("invalid PE")}
   symbols,e:=file.ImportedSymbols();if e!=nil{t.Fatal(e)}
   for _,name:=range []string{"GetStdHandle","WriteFile","ExitProcess"}{found:=false;for _,x:=range symbols{if strings.HasPrefix(x,name+":"){found=true}};if !found{t.Fatal("missing import",name,symbols)}}
   if runtime.GOOS=="windows"{
    b,e:=exec.CommandContext(ctx,out).CombinedOutput()
    if tc.status==0{if e!=nil||string(b)!=tc.want{t.Fatalf("%q %v",b,e)}}else{exit,ok:=e.(*exec.ExitError);if !ok||exit.ExitCode()!=tc.status{t.Fatalf("expected status %d: %q %v",tc.status,b,e)}}
   }
  })
 }
}
