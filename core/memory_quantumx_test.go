package core

import "testing"

func TestAutomaticLoopReclamation(t *testing.T){
 src:=`import tripol.memory
before=memory.used()
total=0
repeat 1000 {
 tmp=ints(65535)
 tmp[0]=2
 total+=tmp[0]
}
show total,before==memory.used()
while total>0 {
 tmp=ints(16)
 tmp[0]=1
 total-=tmp[0]
}
show total,before==memory.used()`
 for _,p:=range []string{"nano","balanced","turbo","debug"}{t.Run(p,func(t *testing.T){got,bad:=hyperRun(t,src,p);if bad||got!="2000 true\n0 true\n"{t.Fatalf("%q %v",got,bad)}})}
}
func TestAutomaticLoopEscapeFallback(t *testing.T){
 runCase(t,`a=ints(1)
repeat 2 {a=ints(2);a[0]=3}
show a`,"[3, 0]\n",false)
}
func TestQuantumXAtomic(t *testing.T){
 runCase(t,`import tripol.atomic
a=ints(2)
show atomic.load(a,0),atomic.exchange(a,0,5)
show atomic.compare_exchange(a,0,4,9),atomic.compare_exchange(a,0,5,9)
show atomic.add(a,0,-3),atomic.load(a,0)
try {atomic.add(a,2,1)} catch e {show e}
atomic.exchange(a,1,9223372036854775807)
try {atomic.add(a,1,1)} catch e {show e}
show atomic.load(a,1)`,"0 0\nfalse true\n6 6\n11\n10\n9223372036854775807\n",false)
}
func TestQuantumXConstantBranch(t *testing.T){
 src:=`x=1
 if 2+3==5 {x+=1} else {show "unreachable"}
 while false {show "unreachable"}
 show x`
 for _,p:=range []string{"nano","balanced","turbo","debug"}{t.Run(p,func(t *testing.T){got,bad:=hyperRun(t,src,p);if bad||got!="2\n"{t.Fatalf("%q %v",got,bad)}})}
 opt:=compileProfile(src,"linux","nano");debug:=compileProfile(src,"linux","debug")
 if len(opt)>=len(debug){t.Fatal("constant branch optimization did not reduce binary")}
}
