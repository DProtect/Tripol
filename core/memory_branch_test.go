package core
import "testing"

func TestQuantumBranchRegions(t *testing.T){
 src:=`import tripol.memory
before=memory.used()
total=0
repeat 1000 {
 temp=ints(65535)
 temp[0]=2
 if total%4==0 {
  branch=ints(65535)
  branch[0]=3
  total+=branch[0]
 } else {
  branch=text(17)
  total+=len(branch)
 }
 total+=temp[0]
}
show total,before==memory.used()
try {
 repeat 2 {
  temp=ints(65535)
  if len(temp)>0 {assert false}
 }
} catch err {show err}
show before==memory.used()
outer="a"
repeat 3 {
 if len(outer)<5 {outer=outer+"b"}
}
show outer
`
 for _,p:=range []string{"nano","balanced","turbo","debug"}{
  t.Run(p,func(t *testing.T){got,bad:=hyperRun(t,src,p);if bad||got!="4001 true\n12\ntrue\nabbb\n"{t.Fatalf("%q %v",got,bad)}})
 }
}
