package core

import "testing"

func TestQuantumXIterationAndDefaults(t *testing.T){
 src:=`struct Player {name:string="Guest";score:int=2+3;active:bool=true}
 p=Player()
 q=Player("Ana",12)
 show p.name,p.score,p.active,q.name,q.score,q.active
 a=[10,20,30]
 b=ints(2);b[0]=4;b[1]=5
 c=list(0);c.add(9);c.add(8);c[0]+=3
 for i,v in a {show i,v}
 for i,v in b {show i,v}
 for i,v in c {show i,v}
 show c[0],sum(c)
 func shrink(x:list)->int {x.clear();return 7}
 try {c[0]=shrink(c)} catch err {show err}
 show c.size()`
 for _,profile:=range []string{"nano","balanced","turbo","debug"}{
  t.Run(profile,func(t *testing.T){got,bad:=hyperRun(t,src,profile);want:="Guest 5 true Ana 12 true\n0 10\n1 20\n2 30\n0 4\n1 5\n0 12\n1 8\n12 20\n11\n0\n";if bad||got!=want{t.Fatalf("%q %v",got,bad)}})
 }
}
func TestQuantumXDefaultsReject(t *testing.T){
 for _,s:=range []string{`struct Bad {x:int="str"}`,`struct Bad {x:int=1;y:int}`,`struct Bad {x:int=1/0}`,`for i,v in range(3) {show v}`,`a=[1];for i,i in a {show i}`,`a=list();const b=a;b[0]=1`}{
  t.Run(s,func(t *testing.T){defer func(){if recover()==nil{t.Fatal("accepted invalid program")}}();compile(s,"linux")})
 }
}
