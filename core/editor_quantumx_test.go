package core

import (
 "bufio"
 "bytes"
 "encoding/json"
 "fmt"
 "io"
 "strings"
 "testing"
)
func TestQuantumFormatter(t *testing.T){
 source:="\ufefffunc f(x:int)->string {\r\n\tif x>0 {\r\n return f\"{x} // { {{}}\"\r\n}\r\n return \"}\" // {\r\n}\r\n"
 want:="\ufefffunc f(x:int)->string {\r\n    if x>0 {\r\n        return f\"{x} // { {{}}\"\r\n    }\r\n    return \"}\" // {\r\n}\r\n"
 got:=formatTripol(source);if got!=want||formatTripol(got)!=got{t.Fatalf("%q",got)}
}
func TestQuantumLSP(t *testing.T){
 var in,out bytes.Buffer
 send:=func(m string){fmt.Fprintf(&in,"Content-Length: %d\r\n\r\n%s",len(m),m)}
 send(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
 send(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///sample.tri","version":1,"text":"import tripol.math\nfunc f(x:int)->int {\nreturn x+1\n}\nshow math.square(f(2))"}}}`)
 send(`{"jsonrpc":"2.0","id":2,"method":"textDocument/hover","params":{"textDocument":{"uri":"file:///sample.tri"},"position":{"line":4,"character":12}}}`)
 send(`{"jsonrpc":"2.0","id":3,"method":"textDocument/formatting","params":{"textDocument":{"uri":"file:///sample.tri"}}}`)
 send(`{"jsonrpc":"2.0","method":"textDocument/didChange","params":{"textDocument":{"uri":"file:///sample.tri","version":2},"contentChanges":[{"text":"show missing_variable"}]}}`)
 send(`{"jsonrpc":"2.0","id":4,"method":"shutdown"}`)
 send(`{"jsonrpc":"2.0","method":"exit"}`)
 if e:=serveLSP(&in,&out);e!=nil{t.Fatal(e)}
 r:=bufio.NewReader(&out);seen:=map[string]bool{}
 for{b,e:=readLSP(r);if e==io.EOF{break};if e!=nil{t.Fatal(e)};var m map[string]any;if e=json.Unmarshal(b,&m);e!=nil{t.Fatal(e)};if m["error"]!=nil{t.Fatal(string(b))};for _,key:=range []string{"documentFormattingProvider","math.square","newText","missing_variable"}{if strings.Contains(string(b),key){seen[key]=true}}}
 for _,key:=range []string{"documentFormattingProvider","math.square","newText","missing_variable"}{if !seen[key]{t.Fatal("missing LSP result",key)}}
}
