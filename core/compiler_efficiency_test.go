package core

import (
 "bytes"
 "fmt"
 "strings"
 "testing"
)

func TestValidationBufferReuse(t *testing.T){
 var source strings.Builder
 for i:=0;i<100;i++{fmt.Fprintf(&source,"func unused%d()->int {show \"literal %d\";return %d}\n",i,i,i)}
 source.WriteString("show 42")
 fs,ss:=parse(source.String())
 old:=generateCode(fs,ss,"linux",false,nil)
 fresh:=validateCode(fs,ss,"linux")
 if len(fresh.b)!=0||len(fresh.fixes)!=0||len(fresh.labels)!=0||len(fresh.data)!=0||cap(fresh.b)>=len(old.b){t.Fatal("validation image not discarded",cap(fresh.b),len(old.b))}
}
func TestCompilerPrunesConstantBranches(t *testing.T){
 src:=`func dead()->int {show "dead marker";return 0}
 func live()->int {return 42}
 if 1>2 {show dead()} else {show live()}
 while false {show dead()}`
 fs,ss:=parse(src);live:=reachable(fs,ss)
 if live["dead"]||!live["live"]{t.Fatal(live)}
 for _,target:=range []string{"linux","windows"}{
  for _,profile:=range []string{"nano","balanced","turbo"}{
   bin:=compileParsed(fs,ss,target,profile)
   if bytes.Contains(bin,[]byte("dead marker")){t.Fatal("retained unreachable function",target,profile)}
   if !bytes.Equal(bin,compileParsed(fs,ss,target,profile)){t.Fatal("non-deterministic")}
  }
 }
 for _,profile:=range []string{"nano","balanced","turbo","debug"}{
  got,bad:=hyperRun(t,src,profile);if bad||got!="42\n"{t.Fatalf("%q %v",got,bad)}
 }
}
func TestCompilerStillValidatesDeadCode(t *testing.T){
 defer func(){if recover()==nil{t.Fatal("accepted invalid dead code")}}()
 compile(`func bad()->int {return true};if false {show bad()}`,"linux")
}
func TestSSAReusableSpillSlots(t *testing.T){
 // Two disjoint groups with register pressure. Both groups must spill,
 // but the second can use the first group's expired frame slots.
 p:=ssaProgram{result:-1}
 for group:=0;group<2;group++{
  base:=len(p.values)
  for i:=0;i<10;i++{p.values=append(p.values,ssaValue{op:"const.i64",literal:int64(i)+1<<40})}
  prev:=base
  for i:=1;i<10;i++{p.values=append(p.values,ssaValue{op:"add.checked.i64",args:[]int{prev,base+i}});prev=len(p.values)-1}
  if group==1{p.result=prev}
 }
 plan:=allocateSSA(p)
 if plan.reusedSlots==0||plan.frame/8>=plan.spills{t.Fatal("expected slot reuse",plan)}
}
func BenchmarkValidationBuffer(b *testing.B){
 var source strings.Builder
 for i:=0;i<300;i++{fmt.Fprintf(&source,"func unused%d(n:int)->int {x=n+%d;show \"item %d\";return x}\n",i,i,i)}
 source.WriteString("show 42")
 fs,ss:=parse(source.String())
 b.Run("accumulated",func(b *testing.B){b.ReportAllocs();for i:=0;i<b.N;i++{generateCode(fs,ss,"linux",false,nil)}})
 b.Run("reused",func(b *testing.B){b.ReportAllocs();for i:=0;i<b.N;i++{validateCode(fs,ss,"linux")}})
}
