package core

import "testing"

func TestNativeSSALinearScan(t *testing.T){
 src:=`func calculate(a,b,c,d,e,f,g,h:int)->int {
  x=a+b
  y=c*d
  z=e-f
  q=g+h
  return (x+y)*(z+q)
 }
 func dead(x:int)->int { unused=x+1;return 7 }
 func neg(x:int)->int { y=-x;return y+2 }
 show calculate(1,2,3,4,10,5,6,7)
 show neg(3)
 try {dead(9223372036854775807)} catch e {show e}
 try {neg(-9223372036854775807-1)} catch e {show e}`
 for _,profile:=range []string{"nano","balanced","turbo","debug"}{
  t.Run(profile,func(t *testing.T){got,bad:=hyperRun(t,src,profile);if bad||got!="270\n-1\n10\n10\n"{t.Fatalf("%q %v",got,bad)}})
 }
 fs,_:=parse(`func f(a:int,b:int,c:int,d:int,e:int,f:int,g:int,h:int)->int {x=a+b;return x+c+d+e+f+g+h}`)
 p,ok:=buildScalarSSA(fs[0]);if !ok{t.Fatal("SSA lowering failed")}
 plan:=allocateSSA(p);if plan.spills==0{t.Fatal("expected pressure and spills")}
}
