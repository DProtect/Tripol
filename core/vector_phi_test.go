package core

import (
 "bytes"
 "testing"
)

func TestVectorFillProfiles(t *testing.T){
 for _,s:=range []string{
  "a=ints(0);for i in range(0){a[i]=7};show 1",
  "a=ints(5);for i in range(5){a[i]=-4294967296};show a[0],a[3],a[4]",
  "a=ints(5);for i in range(4){a[i]=9};show a[3],a[4]",
  "a=ints(1);for i in range(-1){a[i]=8};show a[0]",
  "a=ints(3);for i in range(4){a[i]=7}",
  "a=ints(3);try{for i in range(4){a[i]=7}}catch code{show code};show a[2]",
  "func n(){show 7;return 3};a=ints(3);for i in range(n()){a[i]=9};show a[2]",
 }{
  w,b:=hyperRun(t,s,"debug")
  for _,p:=range []string{"nano","balanced","turbo"}{o,e:=hyperRun(t,s,p);if o!=w||e!=b{t.Fatal(p,s,o,e,w,b)}}
 }
 s:="a=ints(5);for i in range(5){a[i]=7};show a[4]"
 for _,target:=range []string{"linux","windows"}{
  if !bytes.Contains(compileProfile(s,target,"turbo"),[]byte{0xf3,0x0f,0x7f,0x02}){t.Fatal("missing SSE2 store",target)}
 }
}

func TestReturnPhiProfiles(t *testing.T){
 s:="func f(a:int,b:bool)->int {x=a;if b{x=a+1}else{x=a-1};return x};show f(8,true),f(8,false)"
 fs,ss:=parse(s)
 validateCode(fs,ss,"linux")
 body:=fs[len(fs)-1].body
 out:=eliminateReturnPhi(body)
 if len(out)!=len(body)-1||out[len(out)-1].yes[0].op!="return"||body[len(body)-1].op!="return"{t.Fatal("phi not eliminated or AST mutated")}
 for _,source:=range []string{s,"func f(a:int,b:bool)->int{x=a;if b{x=a+1}else{x=a-1};return x};show f(9223372036854775807,true)"}{
  want,bad:=hyperRun(t,source,"debug")
  for _,p:=range []string{"balanced","turbo"}{got,e:=hyperRun(t,source,p);if got!=want||e!=bad{t.Fatal(p,got,e,want,bad)}}
 }
}
