package core
import(
 "context"
 "os"
 "os/exec"
 "path/filepath"
 "runtime"
 "strings"
 "testing"
 "time"
)
func TestGenesisTypeRejection(t *testing.T){
 if runtime.GOARCH!="amd64"||(runtime.GOOS!="linux"&&runtime.GOOS!="windows"){t.Skip("native host required")}
 fs,ss,_:=loadProject("compiler/direct.tri")
 dir:=t.TempDir();stage:=filepath.Join(dir,"stage.exe")
 if e:=os.WriteFile(stage,compileParsed(fs,ss,runtime.GOOS,"nano"),0755);e!=nil{t.Fatal(e)}
 for _,src:=range []string{
  "x=true;x=1",
  "break",
  "if true {continue}",
  "x=true;x+=1",
  "x=1;if x {show 1}",
  "x=true;show x+1",
  "func f(a:bool)->int {return a};show f(true)",
  "func f(a:bool)->int {return 1};show f(1)",
  "func f(a:int)->bool {return a<2};show f(1)",
  "if false and 1 {show 1}",
  "for i in range(true) {show i}",
 }{
  t.Run(src,func(t *testing.T){
   input:=filepath.Join(dir,"source.tri");output:=filepath.Join(dir,"output")
   _=os.Remove(output)
   if e:=os.WriteFile(input,[]byte(src),0600);e!=nil{t.Fatal(e)}
   ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second);defer cancel()
   cmd:=exec.CommandContext(ctx,stage);cmd.Stdin=strings.NewReader(input+"\n"+output+"\n")
   b,e:=cmd.CombinedOutput();if e==nil||ctx.Err()!=nil{t.Fatalf("expected compile-time rejection: %q %v",b,e)}
   if _,e=os.Stat(output);!os.IsNotExist(e){t.Fatal("invalid source created output")}
  })
 }
}
