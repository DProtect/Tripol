package core

import (
 "bytes"
 "context"
 "fmt"
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
 for i,tc:=range []struct{src,want string;status int}{
 {"a=0\nfor i in range(6) {if i==1 {continue};if i==4 {break};a+=i}\nwhile a>0 {a-=1;if a==1 {break};if a==4 {continue};show a}\nshow a\n","3\n2\n1\n",0},
 {"func multiply(a:int,b:int)->int {return a*b}\nfor i in range(3) {show multiply(i,5)}","0\n5\n10\n",0},
 {"a=9223372036854775807;show -a-1","-9223372036854775808\n",0},
 {"a=7;show a/0","",10},
 {"a=9223372036854775807;show a+1","",10},
 {"a=-7;show a%3","-1\n",0},
 {"func even(n:int)->bool {return n%2==0}\nfunc choose(flag:bool,a:int,b:int)->int {if flag {return a};return b}\nfunc both(a:bool,b:bool)->bool {return a and b}\nflag=even(8)\nother=not even(5)\nflag=both(flag,other)\nshow choose(flag,10,20)\nshow choose(false and 1/0==0,1,2)\nif flag==true {show 7}\nfor i in range(2) {if even(i) {show i}}\n","10\n2\n7\n0\n",0},
 {"func combine(a:int,b:int)->int {return a*10+b}\nx=3\nx+=combine(1+2,4)\nx*=2\nx-=4\nx/=2\nx%=10\nshow x\nshow combine(2,3+combine(1,2))\nshow -(2+3)*4\n","5\n35\n-20\n",0},
 {"a=9223372036854775807;a+=1","",10},
 {"func first(n:int)->int {\n return second(n)+1\n}\nfunc second(n:int)->int {\n if n==0 {return 0}\n return first(n-1)\n}\nx=2\nif false and 1/0==0 {show 99} else if (x+1)>2 and not (x==0 or false) {show first(3)} else {show 88}\nif true or 1/0==0 {show 42}\nwhile x>0 and true {show x;x=x-1}\n","4\n42\n2\n1\n",0},

 }{
  t.Run(fmt.Sprintf("case_%02d",i),func(t *testing.T){
   t.Log(tc.src)
   dir:=t.TempDir();src:=filepath.Join(dir,"input.tri");out:=filepath.Join(dir,"genesis.exe")
   if e:=os.WriteFile(src,[]byte(tc.src),0600);e!=nil{t.Fatal(e)}
   ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second);defer cancel()
   emit:=func(destination string)[]byte{
    cmd:=exec.CommandContext(ctx,stage);cmd.Stdin=strings.NewReader(src+"\n"+destination+"\n")
    if b,e:=cmd.CombinedOutput();e!=nil{t.Fatalf("stage: %q %v",b,e)}
    b,e:=os.ReadFile(destination);if e!=nil{t.Fatal(e)};return b
   }
   image:=emit(out);if !bytes.Equal(image,emit(filepath.Join(dir,"second.exe"))){t.Fatal("non-deterministic PE")}
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
