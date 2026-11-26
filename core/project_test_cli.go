package core

import (
 "bytes"
 "context"
 "fmt"
 "io"
 "os"
 "os/exec"
 "path/filepath"
 "runtime"
 "strings"
 "time"
)

// The runner executes trusted project code, not a sandbox. Only explicitly
// selected files or immediate *.test.tri directory members are executed.
func projectTestFiles(path string)([]string,error){
 info,e:=os.Lstat(path);if e!=nil{return nil,e}
 if info.Mode().IsRegular(){
  if !strings.HasSuffix(path,".tri"){return nil,fmt.Errorf("expected a .tri file")}
  return []string{path},nil
 }
 if !info.IsDir(){return nil,fmt.Errorf("expected a regular file or directory")}
 entries,e:=os.ReadDir(path);if e!=nil{return nil,e}
 files:=[]string{}
 for _,entry:=range entries{
  if strings.HasSuffix(entry.Name(),".test.tri")&&entry.Type().IsRegular(){files=append(files,filepath.Join(path,entry.Name()))}
  if len(files)>64{return nil,fmt.Errorf("select at most 64 tests per run")}
 }
 if len(files)==0{return nil,fmt.Errorf("no *.test.tri files found in %s",path)}
 return files,nil
}
type testOutput struct{ buf bytes.Buffer; truncated bool }
func(w *testOutput) Write(p []byte)(int,error){
 n:=len(p);remain:=65536-w.buf.Len()
 if len(p)>remain{p=p[:remain];w.truncated=true}
 _,_ = w.buf.Write(p);return n,nil
}
func runProjectTest(ctx context.Context,path,target string)(detail string,err error){
 defer func(){if p:=recover();p!=nil{d:=normalizeDiagnostic(p);err=fmt.Errorf("%s: %s",d.Code,d.Message)}}()
 fs,ss,_:=loadProject(path)
 bin:=compileParsed(fs,ss,target,"balanced")
 dir,e:=os.MkdirTemp("","tripol-test-");if e!=nil{return "",e};defer os.RemoveAll(dir)
 exe:=filepath.Join(dir,"test");if target=="windows"{exe+=".exe"}
 if e=os.WriteFile(exe,bin,0700);e!=nil{return "",e}
 child,cancel:=context.WithTimeout(ctx,10*time.Second);defer cancel()
 cmd:=exec.CommandContext(child,exe);cmd.Dir=dir
 // No interactive stdin; tests that need files can create their own.
 output:=&testOutput{};cmd.Stdout=output;cmd.Stderr=output
 cmd.WaitDelay=2*time.Second
 e=cmd.Run()
 detail=output.buf.String();if output.truncated{detail+="\n[output truncated]\n"}
 if child.Err()!=nil{return detail,fmt.Errorf("test timed out or was cancelled")}
 return detail,e
}
func projectTestCLI(args []string){
 if len(args)>1{die("usage: tripol test [FILE.tri|DIRECTORY]")}
 if runtime.GOARCH!="amd64"||(runtime.GOOS!="windows"&&runtime.GOOS!="linux"){die("native tests require Windows or Linux x86-64")}
 path:=".";if len(args)==1{path=args[0]}
 files,e:=projectTestFiles(path);if e!=nil{die("%v",e)}
 ctx,cancel:=context.WithTimeout(context.Background(),2*time.Minute);defer cancel()
 passed,failed:=0,0
 for _,file:=range files{
  if ctx.Err()!=nil{failed+=len(files)-passed-failed;break}
  detail,e:=runProjectTest(ctx,file,runtime.GOOS)
  if e!=nil{
   fmt.Printf("FAIL %s: %v\n",file,e)
   if detail!=""{_,_=io.WriteString(os.Stdout,detail);fmt.Println()}
   failed++
  }else{fmt.Printf("PASS %s\n",file);passed++}
 }
 fmt.Printf("%d passed, %d failed\n",passed,failed)
 if failed>0{die("project tests failed")}
}
