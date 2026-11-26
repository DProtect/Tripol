package core

import (
 "bytes"
 "testing"
)

func TestWebEngineR1(t *testing.T){
 for _,source:=range []string{
  "show 42",
  "n=0\nfor i in range(1,6) {n+=i}\nshow n",
  "func double(n:int)->int {return n*2}\nshow double(21)",
  "show true and not false",
  "import tripol.math\nshow math.square(7)",
 }{
  b,err:=CompileWeb(source);if err!=nil{t.Fatalf("%s: %v",source,err)}
  if len(b)<8||!bytes.Equal(b[:8],[]byte{0,97,115,109,1,0,0,0}){t.Fatal("missing WASM header")}
 }
 for _,source:=range []string{"show missing","show \"not supported\"","import tripol.files\nshow 1"}{if _,err:=CompileWeb(source);err==nil{t.Fatalf("accepted %s",source)}}
}
