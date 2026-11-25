package core

import (
 "fmt"
 "io"
 "os"
 "strings"
)

const editorSourceLimit=4<<20

// Preserve token spelling, line separators and comments. Only line indentation
// changes; preserve BOM, CRLF and interpolation adjacency. Never reserialize
// strings or edit source files implicitly.
func formatTripol(source string)string{
 if len(source)>editorSourceLimit{die("formatter source limit")}
 bom:="";if strings.HasPrefix(source,"\ufeff"){bom="\ufeff";source=strings.TrimPrefix(source,bom)}
 ts:=lex(source)
 lines:=strings.Split(source,"\n")
 depth:=0;pos:=0
 for i,line:=range lines{
  first:=pos
  for pos<len(ts)&&ts[pos].line==i+1&&ts[pos].k!="eof"{pos++}
  head:=depth
  for k:=first;k<pos;k++{
   t:=ts[k];if t.k=="op"&&t.s=="}"{head--;continue};break
  }
  if head<0{head=0}
  trimmed:=strings.TrimLeft(line," \t")
  if trimmed!=""&&trimmed!="\r"{lines[i]=strings.Repeat("    ",head)+trimmed}else{lines[i]=trimmed}
  for k:=first;k<pos;k++{t:=ts[k];if t.k!="op"{continue};if t.s=="{"{depth++};if t.s=="}"{depth--;if depth<0{depth=0}}}
 }
 result:=strings.Join(lines,"\n")
 // Lexical equivalence, including significant f"..." adjacency. Line/column
 // offsets may differ, but token spelling and within-line distances may not.
 after:=lex(result)
 if len(ts)!=len(after){die("formatter changed token count")}
 for i,t:=range ts{
  u:=after[i];if t.k!=u.k||t.s!=u.s||t.line!=u.line{die("formatter changed tokens")}
  if i>0&&ts[i-1].line==t.line&&t.k!="eof"&&t.col-ts[i-1].col!=u.col-after[i-1].col{die("formatter changed token adjacency")}
 }
 return bom+result
}
func formatCLI(args []string){
 if len(args)>1{die("usage: tripol fmt [FILE.tri|--stdin]")}
 path:="main.tri";if len(args)==1{path=args[0]}
 var reader io.Reader=os.Stdin
 if path!="--stdin"{f,e:=os.Open(path);if e!=nil{die("%v",e)};defer f.Close();reader=f}
 b,e:=io.ReadAll(io.LimitReader(reader,editorSourceLimit+1));if e!=nil{die("%v",e)}
 fmt.Print(formatTripol(string(b)))
}
