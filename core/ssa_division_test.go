package core

import "testing"

func TestNativeSSACheckedDivision(t *testing.T){
 src:=`func calculate(a:int,b:int,c:int,d:int,e:int,f:int,g:int,h:int)->int {
  q=a/b
  r=a%b
  return q+r+c+d+e+f+g+h
 }
 func divide(a:int,b:int)->int {q=a/b;return q}
 func modulo(a:int,b:int)->int {r=a%b;return r}
 func power(a:int)->int {q=a/8;r=a%8;return q+r}
 func unused(a:int,b:int)->int {q=a/b;return 1}
 show calculate(-17,5,1,2,3,4,5,6)
 show divide(17,-5),modulo(17,-5),power(-19)
 show divide(-9223372036854775807-1,1),modulo(-9223372036854775807-1,8)
 try {divide(1,0)} catch err {show err}
 try {modulo(-9223372036854775807-1,-1)} catch err {show err}
 try {unused(1,0)} catch err {show err}
 `
 for _,profile:=range []string{"nano","balanced","turbo","debug"}{
  t.Run(profile,func(t *testing.T){got,bad:=hyperRun(t,src,profile);want:="16\n-3 2 -5\n-9223372036854775808 0\n10\n10\n10\n";if bad||got!=want{t.Fatalf("%q %v",got,bad)}})
 }
 fs,_:=parse(`func f(a:int,b:int,c:int,d:int,e:int,f:int)->int {x=a/b;return x+c+d+e+f}`)
 p,ok:=buildScalarSSA(fs[0]);if !ok{t.Fatal("SSA lowering failed")}
 plan:=allocateSSA(p);if plan.registers!=4||plan.spills==0{t.Fatal("division register reservation/spills",plan)}
 for _,l:=range plan.locations{if l.present&&l.spill==0&&(l.reg==1||l.reg==2){t.Fatal("allocated division scratch register")}}
}
