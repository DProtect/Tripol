package core

import(
 "bufio"
 "context"
 "fmt"
 "io"
 "os"
 "path/filepath"
 "runtime"
 "strings"
)

// Native compile-and-run teaching shell. Only deterministic scalar expressions
// and assignments are replayed; no IO, mutation of aggregates, arbitrary
// calls or control flow can be silently repeated.
func shellPure(e *expr)bool{
 if e==nil{return false}
 switch e.op{
 case "int","float","bool","string","var":return true
 case "un","bin","format":
  for _,a:=range e.args{if !shellPure(a){return false}};return true
 case "call":
  switch e.s{case "len","abs","min","max","text","slice","float","sqrt","parse_int":
   for _,a:=range e.args{if !shellPure(a){return false}};return true
  }
 }
 return false
}
func shellLine(history,line,target string)(next,output string,err error){
 next=history
 defer func(){if p:=recover();p!=nil{d:=normalizeDiagnostic(p);err=fmt.Errorf("%s",d.Message);next=history}}()
 imports,fs,ss:=parseUnit(line)
 if len(imports)!=0||len(fs)!=0||len(ss)!=1{return next,"",fmt.Errorf("enter one scalar expression, assignment or show statement")}
 s:=ss[0]
 if s.op!="assign"&&s.op!="decl"&&s.op!="show"&&s.op!="showmany"&&s.op!="expr"{return next,"",fmt.Errorf("statement unsupported in teaching shell")}
 if s.op=="showmany"{for _,x:=range s.e.args{if !shellPure(x){return next,"",fmt.Errorf("only pure scalar expressions are supported")}}}else if !shellPure(s.e){return next,"",fmt.Errorf("only pure scalar expressions are supported")}
 source:=history+line+"\n"
 if s.op=="expr"{source=history+"show "+line+"\n"}
 if len(source)>65536{return next,"",fmt.Errorf("session limit; use :reset")}
 dir,e:=os.MkdirTemp("","tripol-shell-");if e!=nil{return next,"",e};defer os.RemoveAll(dir)
 path:=filepath.Join(dir,"shell.tri")
 if e=os.WriteFile(path,[]byte(source),0600);e!=nil{return next,"",e}
 output,e=runProjectTest(context.Background(),path,target)
 if e!=nil{return next,output,e}
 if s.op=="assign"||s.op=="decl"{next=source}
 return next,output,nil
}
func shellSession(in io.Reader,out io.Writer,target string){
 scanner:=bufio.NewScanner(in);scanner.Buffer(make([]byte,4096),65536)
 history:=""
 fmt.Fprintln(out,"TRIPOL Shell R1 — native scalar playground")
 fmt.Fprintln(out,":help  :reset  :code  :exit")
 for{
  fmt.Fprint(out,">>> ")
  if !scanner.Scan(){break}
  line:=strings.TrimSpace(scanner.Text())
  switch line{
  case "":continue
  case ":exit":return
  case ":reset":history="";fmt.Fprintln(out,"Session reset.");continue
  case ":code":fmt.Fprint(out,history);continue
  case ":help":fmt.Fprintln(out,"name = \"Alex\" | show f\"Hello {name}\" | 25 * 4\nScalar assignments persist; only pure expressions are accepted.");continue
  }
  next,result,e:=shellLine(history,line,target)
  if result!=""{fmt.Fprint(out,result)}
  if e!=nil{fmt.Fprintln(out,"Error:",e)}else{history=next}
 }
 if e:=scanner.Err();e!=nil{fmt.Fprintln(out,"Error:",e)}
}
func shellCLI(args []string){
 if len(args)!=0{die("usage: tripol shell")}
 if runtime.GOARCH!="amd64"||(runtime.GOOS!="linux"&&runtime.GOOS!="windows"){die("shell requires Windows or Linux x86-64")}
 shellSession(os.Stdin,os.Stdout,runtime.GOOS)
}
