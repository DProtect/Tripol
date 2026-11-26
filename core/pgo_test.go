package core

import (
 "bytes"
 "os"
 "os/exec"
 "path/filepath"
 "reflect"
 "runtime"
 "testing"
)

func TestPGOOrder(t *testing.T){
 fs:=[]fn{{name:"cold"},{name:"hot"},{name:"warm"},{name:"tie"}}
 order:=pgoOrder(fs,map[string]uint64{"hot":8,"warm":4,"tie":4})
 got:=[]string{};for _,f:=range order{got=append(got,f.name)}
 if !reflect.DeepEqual(got,[]string{"hot","warm","tie","cold"})||fs[0].name!="cold"{t.Fatal(got)}
}
func TestPGOTrainingAndLayout(t *testing.T){
 if runtime.GOARCH!="amd64"||(runtime.GOOS!="linux"&&runtime.GOOS!="windows"){t.Skip("native host required")}
 source:="func cold(x:int)->int {show x;return x};func hot(x:int)->int {show x;return x};show cold(1);show hot(2);show hot(3)"
 fs,ss:=parse(source)
 ext:="";if runtime.GOOS=="windows"{ext=".exe"}
 execute:=func(image []byte)(string,string){
  p:=filepath.Join(t.TempDir(),"program"+ext)
  if e:=os.WriteFile(p,image,0755);e!=nil{t.Fatal(e)}
  cmd:=exec.Command(p);var stderr bytes.Buffer;cmd.Stderr=&stderr
  stdout,e:=cmd.Output();if e!=nil{t.Fatal(e,stderr.String())}
  return string(stdout),stderr.String()
 }
 training:=linkCode(generateCodePGO(fs,ss,runtime.GOOS,false,nil,false,nil,true))
 stdout,log:=execute(training)
 counts:=pgoCounts(log,fs)
 if counts["hot"]!=2||counts["cold"]!=1{t.Fatal(counts,log)}
 for _,target:=range []string{"linux","windows"}{
  a:=generateCodePGO(fs,ss,target,true,reachable(fs,ss),false,counts,false,true,true)
  if a.labels["fn:hot"]>=a.labels["fn:cold"]{t.Fatal("hot layout missing",target)}
  p:=compileParsedPGO(fs,ss,target,"turbo",counts)
  if !bytes.Equal(p,compileParsedPGO(fs,ss,target,"turbo",counts)){t.Fatal("nondeterministic PGO",target)}
  if target==runtime.GOOS{
   out,errout:=execute(p)
   if out!=stdout||errout!=""{t.Fatal("PGO changed output",out,errout,stdout)}
  }
 }
}
