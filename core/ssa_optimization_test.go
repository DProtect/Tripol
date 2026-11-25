package core
import "testing"
func TestQuantumSSARedundancy(t *testing.T){
 fs,_:=parse("func f(a:int)->int {x=a+3;y=a+3;k=4;z=k*5;return x+y+z}")
 p,ok:=buildScalarSSA(fs[0]);if !ok||p.reused<1||p.propagated<1{t.Fatal(p)}
 for _,profile:=range []string{"nano","balanced","turbo","debug"}{
  t.Run(profile,func(t *testing.T){
   got,bad:=hyperRun(t,`func f(a:int)->int {x=a+3;y=a+3;k=4;z=k*5;return x+y+z}
 show f(7)
 try {f(9223372036854775807)} catch err {show err}`,profile)
   if bad||got!="40\n10\n"{t.Fatalf("%q %v",got,bad)}
  })
 }
}
