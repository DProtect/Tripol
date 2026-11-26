package core
import(
 "bytes"
 "testing"
)
func TestCFGConstantJoins(t *testing.T){
 src:=`func unused()->int {show "unreachable CFG marker";return 1}
 func choice(flag:bool)->int {
  x=0
  if flag {x=7} else {x=7}
  if x!=7 {return unused()}
  return x*6
 }
 show choice(true),choice(false)
 x=3
 if x==3 {x=4} else {x=9}
 show x
 repeat 3 {x+=1}
 show x
 try {x=1;show 1/0} catch e {x=8}
 show x
 `
 for _,profile:=range []string{"nano","balanced","turbo","debug"}{
  t.Run(profile,func(t *testing.T){got,bad:=hyperRun(t,src,profile);if bad||got!="42 42\n4\n7\n8\n"{t.Fatalf("%q %v",got,bad)}})
 }
 fs,ss:=parse(src)
 optimized,body,stats:=optimizeCFG(fs,ss)
 if stats.Joins==0||stats.Substitutions==0||stats.ConstantBranches==0||reachable(optimized,body)["unused"]{t.Fatal(stats)}
 for _,target:=range []string{"linux","windows"}{
  for _,p:=range []string{"balanced","turbo"}{
   if bytes.Contains(compileParsed(fs,ss,target,p),[]byte("unreachable CFG marker")){t.Fatal("dead CFG function retained",target,p)}
  }
 }
}
func TestCFGCheckedSemantics(t *testing.T){
 src:=`func divergent(flag:bool)->int {
  x=1
  if flag {x=2} else {x=3}
  return x
 }
 show divergent(true),divergent(false)
 x=9223372036854775807
 try {show x+1} catch e {show e}
 y=2
 while y>0 {show y;y-=1}
 show y
 `
 for _,p:=range []string{"nano","balanced","turbo","debug"}{got,bad:=hyperRun(t,src,p);if bad||got!="2 3\n10\n2\n1\n0\n"{t.Fatalf("%s %q %v",p,got,bad)}}
}
func TestCFGDoesNotMutateSource(t *testing.T){
 fs,ss:=parse("x=3;show x+1")
 a:=ss[1].e.args[0]
 _,out,stats:=optimizeCFG(fs,ss)
 if a.op!="var"||ss[1].e.args[0]!=a||out[1].e.args[0].op!="int"||stats.Substitutions!=1{t.Fatal("unexpected mutation or substitution")}
}
