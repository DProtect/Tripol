package core

import (
 "bufio"
 "encoding/json"
 "fmt"
 "io"
 "regexp"
 "runtime"
 "strconv"
 "strings"
 "unicode/utf16"
)

type lspPosition struct{ Line int `json:"line"`; Character int `json:"character"` }
type lspRange struct{ Start lspPosition `json:"start"`; End lspPosition `json:"end"` }
type lspDocument struct{ URI string `json:"uri"`; Text string `json:"text"`; Version int `json:"version"` }
type lspParams struct{
 TextDocument lspDocument `json:"textDocument"`
 Position lspPosition `json:"position"`
 ContentChanges []struct{Text string `json:"text"`;Range *lspRange `json:"range"`} `json:"contentChanges"`
}
type lspRequest struct{ID json.RawMessage `json:"id"`; Method string `json:"method"`; Params json.RawMessage `json:"params"`}
type lspServer struct{out io.Writer;docs map[string]lspDocument;shutdown bool;initialized bool}
var lspDecl=regexp.MustCompile(`^\s*(func|struct)\s+([A-Za-z_][A-Za-z_0-9]*)\b`)
var lspSignature=regexp.MustCompile(`(?m)^\s*func\s+([A-Za-z_][A-Za-z_0-9]*)\s*(\([^\n]*?\)\s*(?:->\s*\w+)?)\s*\{`)
var lspWord=regexp.MustCompile(`[A-Za-z_][A-Za-z_0-9]*(?:\.[A-Za-z_0-9]*)?`)

func readLSP(r *bufio.Reader)([]byte,error){
 n:=-1;headers:=0
 for{
  line,e:=r.ReadString('\n');if e!=nil{return nil,e}
  headers+=len(line);if headers>8192{return nil,fmt.Errorf("LSP header limit")}
  line=strings.TrimRight(line,"\r\n");if line==""{break}
  key,value,ok:=strings.Cut(line,":");if !ok{return nil,fmt.Errorf("invalid LSP header")}
  if strings.EqualFold(key,"Content-Length"){
   if n>=0{return nil,fmt.Errorf("duplicate content length")}
   n,e=strconv.Atoi(strings.TrimSpace(value));if e!=nil||n<0||n>editorSourceLimit{return nil,fmt.Errorf("invalid LSP content length")}
  }
 }
 if n<0{return nil,fmt.Errorf("missing content length")}
 b:=make([]byte,n);_,e:=io.ReadFull(r,b);return b,e
}
func (s *lspServer) send(value any){
 b,e:=json.Marshal(value);if e!=nil{die("LSP response: %v",e)}
 if _,e=fmt.Fprintf(s.out,"Content-Length: %d\r\n\r\n",len(b));e!=nil{die("%v",e)}
 if _,e=s.out.Write(b);e!=nil{die("%v",e)}
}
func(s *lspServer) reply(id json.RawMessage,result any){if len(id)!=0{s.send(map[string]any{"jsonrpc":"2.0","id":id,"result":result})}}
func(s *lspServer) fail(id json.RawMessage,code int,msg string){if len(id)!=0{s.send(map[string]any{"jsonrpc":"2.0","id":id,"error":map[string]any{"code":code,"message":msg}})}}
func utf16len(s string)int{return len(utf16.Encode([]rune(s)))}
func editorLine(text string,line int)string{lines:=strings.Split(text,"\n");if line<0||line>=len(lines){return ""};return strings.TrimSuffix(lines[line],"\r")}
func editorWord(text string,pos lspPosition)string{
 line:=editorLine(text,pos.Line);if pos.Character<0{return ""}
 for _,m:=range lspWord.FindAllStringIndex(line,-1){
  start,end:=utf16len(line[:m[0]]),utf16len(line[:m[1]])
  if pos.Character>=start&&pos.Character<=end{return line[m[0]:m[1]]}
 }
 return ""
}
func moduleSignatures(name string)map[string]string{
 signatures:=map[string]string{}
 for _,mod:=range standardNames{if name==mod{for _,m:=range lspSignature.FindAllStringSubmatch(standardSource(name),-1){signatures[m[1]]="func "+name+"."+m[1]+m[2]};break}}
 return signatures
}
func editorSymbols(text string)[]map[string]any{
 result:=[]map[string]any{}
 for i,line:=range strings.Split(text,"\n"){
  if m:=lspDecl.FindStringSubmatchIndex(line);m!=nil{
   kind:=12;if line[m[2]:m[3]]=="struct"{kind=23}
   r:=lspRange{lspPosition{i,utf16len(line[:m[4]])},lspPosition{i,utf16len(line[:m[5]])}}
   result=append(result,map[string]any{"name":line[m[4]:m[5]],"kind":kind,"range":r,"selectionRange":r})
  }
 }
 return result
}
func validateEditorText(doc lspDocument)(out []map[string]any){
 out=[]map[string]any{}
 defer func(){if p:=recover();p!=nil{
  d:=normalizeDiagnostic(p);line:=d.Line-1;if line<0{line=0}
  lines:=strings.Split(doc.Text,"\n");if line>=len(lines){line=len(lines)-1}
  raw:=editorLine(doc.Text,line);col:=d.Column-1;if col<0{col=0};if col>len(raw){col=len(raw)}
  start:=utf16len(raw[:col]);end:=start;if col<len(raw){end++}
  out=append(out,map[string]any{"range":lspRange{lspPosition{line,start},lspPosition{line,end}},"severity":1,"source":"TRIPOL","code":d.Code,"message":d.Message})
 }}()
 imports,fs,ss:=parseNamed(doc.Text,doc.URI)
 libs:=standardLoader()
 for _,name:=range imports{
  // Local/Vault resolution is intentionally not performed for unsaved buffers:
  // do not read arbitrary files or execute compiler subprocesses from LSP.
  if !strings.HasPrefix(name,"@std:"){return out}
  libs.add(strings.TrimPrefix(name,"@std:"))
 }
 target:=runtime.GOOS;if target!="windows"{target="linux"}
 generateCode(append(libs.functions,fs...),ss,target,false,nil)
 return out
}
func(s *lspServer) publish(doc lspDocument){
 s.send(map[string]any{"jsonrpc":"2.0","method":"textDocument/publishDiagnostics","params":map[string]any{"uri":doc.URI,"version":doc.Version,"diagnostics":validateEditorText(doc)}})
}
func(s *lspServer) dispatch(req lspRequest)(stop bool){
 defer func(){if p:=recover();p!=nil{s.fail(req.ID,-32603,fmt.Sprint(p))}}()
 if req.Method=="exit"{return true}
 if req.Method=="initialize"{
  if s.initialized{s.fail(req.ID,-32600,"already initialized");return false}
  s.initialized=true
  s.reply(req.ID,map[string]any{"serverInfo":map[string]any{"name":"TRIPOL","version":version},"capabilities":map[string]any{"positionEncoding":"utf-16","textDocumentSync":map[string]any{"openClose":true,"change":1},"completionProvider":map[string]any{"triggerCharacters":[]string{"."}},"hoverProvider":true,"definitionProvider":true,"documentSymbolProvider":true,"documentFormattingProvider":true}})
  return false
 }
 if !s.initialized{s.fail(req.ID,-32002,"server not initialized");return false}
 if req.Method=="shutdown"{s.shutdown=true;s.reply(req.ID,nil);return false}
 if s.shutdown{s.fail(req.ID,-32600,"server shut down");return false}
 var p lspParams
 if len(req.Params)>0&&json.Unmarshal(req.Params,&p)!=nil{s.fail(req.ID,-32602,"invalid params");return false}
 uri:=p.TextDocument.URI
 doc,exists:=s.docs[uri]
 switch req.Method{
 case "initialized","$/cancelRequest","textDocument/didSave":return false
 case "textDocument/didOpen":
  if len(p.TextDocument.Text)>editorSourceLimit||len(uri)>8192||uri==""||(!exists&&len(s.docs)>=32){return false}
  s.docs[uri]=p.TextDocument;s.publish(p.TextDocument)
 case "textDocument/didChange":
  if !exists||len(p.ContentChanges)!=1||p.ContentChanges[0].Range!=nil||len(p.ContentChanges[0].Text)>editorSourceLimit||p.TextDocument.Version<doc.Version{return false}
  doc.Text=p.ContentChanges[0].Text;doc.Version=p.TextDocument.Version;s.docs[uri]=doc;s.publish(doc)
 case "textDocument/didClose":
  delete(s.docs,uri);s.send(map[string]any{"jsonrpc":"2.0","method":"textDocument/publishDiagnostics","params":map[string]any{"uri":uri,"diagnostics":[]any{}}})
 case "textDocument/completion":
  items:=[]map[string]any{}
  word:=editorWord(doc.Text,p.Position)
  if module,_,ok:=strings.Cut(word,".");ok{
   if module=="tripol"{for _,name:=range standardNames{items=append(items,map[string]any{"label":name,"kind":9})}}else{for name,sig:=range moduleSignatures(module){items=append(items,map[string]any{"label":name,"kind":3,"detail":sig})}}
  }else{
   for _,k:=range strings.Fields("import func struct int bool float string ints list let const if else while for in range repeat return show assert try catch finally scratch throw break continue true false and or not"){items=append(items,map[string]any{"label":k,"kind":14})}
   for _,sym:=range editorSymbols(doc.Text){items=append(items,map[string]any{"label":sym["name"],"kind":3})}
  }
  s.reply(req.ID,items)
 case "textDocument/hover":
  word:=editorWord(doc.Text,p.Position);module,name,ok:=strings.Cut(word,".");if ok{
   if sig:=moduleSignatures(module)[name];sig!=""{s.reply(req.ID,map[string]any{"contents":map[string]string{"kind":"plaintext","value":sig}});break}
  }
  s.reply(req.ID,nil)
 case "textDocument/definition":
  word:=editorWord(doc.Text,p.Position);var result any
  for _,sym:=range editorSymbols(doc.Text){if sym["name"]==word{result=map[string]any{"uri":uri,"range":sym["range"]};break}}
  s.reply(req.ID,result)
 case "textDocument/documentSymbol":s.reply(req.ID,editorSymbols(doc.Text))
 case "textDocument/formatting":
  if !exists{s.reply(req.ID,[]any{});break}
  formatted:=formatTripol(doc.Text);if formatted==doc.Text{s.reply(req.ID,[]any{});break}
  lines:=strings.Split(doc.Text,"\n")
  r:=lspRange{lspPosition{0,0},lspPosition{len(lines)-1,utf16len(lines[len(lines)-1])}}
  s.reply(req.ID,[]any{map[string]any{"range":r,"newText":formatted}})
 default:s.fail(req.ID,-32601,"method not supported")
 }
 return false
}
func serveLSP(in io.Reader,out io.Writer)error{
 r:=bufio.NewReader(in);s:=lspServer{out:out,docs:map[string]lspDocument{}}
 for{
  b,e:=readLSP(r);if e==io.EOF{return nil};if e!=nil{return e}
  var req lspRequest
  if json.Unmarshal(b,&req)!=nil{s.send(map[string]any{"jsonrpc":"2.0","id":nil,"error":map[string]any{"code":-32700,"message":"parse error"}});continue}
  if s.dispatch(req){return nil}
 }
}
