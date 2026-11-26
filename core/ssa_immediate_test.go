package core

import "testing"

func TestSSAImmediateArithmetic(t *testing.T){
 src:=`func compute(x:int)->int {
  a=x+127
  b=a-128
  c=b*300
  return c+2147483647
 }
 func negative(x:int)->int {
  a=x*-128
  return a-2147483648
 }
 func wide(x:int)->int {a=x+2147483648;return a*2}
 func checked(x:int)->int {a=x+1;return a*2}
 show compute(5)
 show negative(2)
 show wide(1)
 try {checked(9223372036854775807)} catch err {show err}
 try {checked(4611686018427387903)} catch err {show err}
 `
 for _,profile:=range []string{"nano","balanced","turbo","debug"}{
  t.Run(profile,func(t *testing.T){got,bad:=hyperRun(t,src,profile);if bad||got!="2147484847\n-2147483904\n4294967298\n10\n10\n"{t.Fatalf("%q %v",got,bad)}})
 }
 fs,_:=parse("func f(x:int)->int {a=x+3;b=a*7;return b-8}")
 p,ok:=buildScalarSSA(fs[0]);if !ok{t.Fatal("SSA failed")}
 plan:=allocateSSA(p);if plan.immediateValues!=3{t.Fatal(plan)}
 for i,v:=range p.values{if v.op=="const.i64"&&plan.locations[i].present{t.Fatal("literal unnecessarily allocated",i)}}
}
