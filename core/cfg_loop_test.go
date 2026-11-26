package core
import "testing"

func TestCFGLoopInvariantFacts(t *testing.T){
 src:=`factor=3
 count=4
 total=0
 for i in range(count) {
  total+=i*factor
  if i==2 {count=99}
 }
 show total,count,factor
 i=3
 while i>0 {show i*factor;i-=1}
 show i
 repeat 2 {
  try {factor+=1;assert false} catch err {total+=err}
 }
 show factor,total
 `
 for _,p:=range []string{"nano","balanced","turbo","debug"}{
  t.Run(p,func(t *testing.T){got,bad:=hyperRun(t,src,p);if bad||got!="18 99 3\n9\n6\n3\n0\n5 42\n"{t.Fatalf("%q %v",got,bad)}})
 }
 fs,ss:=parse("factor=3;total=0;for i in range(4){total+=i*factor};show factor")
 _,out,stats:=optimizeCFG(fs,ss)
 if stats.LoopInvariants==0||out[2].yes[0].e.args[1].args[1].op!="int"||out[3].e.op!="int"{t.Fatal(stats)}
}
func TestCFGLoopDoesNotMoveTraps(t *testing.T){
 src:=`x=0
 n=0
 repeat n {show 1/x}
 show 9
 n=1
 try {repeat n {show 1/x}} catch err {show err}
 `
 for _,p:=range []string{"nano","balanced","turbo","debug"}{got,bad:=hyperRun(t,src,p);if bad||got!="9\n10\n"{t.Fatalf("%s %q %v",p,got,bad)}}
}
