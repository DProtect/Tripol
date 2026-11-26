package core
import "testing"
func TestSimpleListMethods(t *testing.T){
 src:=`names=list(["Carlos","Ana","Pedro","Ana"])
 show names.contains("Ana"),names.remove("Ana"),names.remove("Luis")
 names.sort()
 show names
 a=list([5,-1,3,3,0])
 a.sort()
 show a
 show a.remove(3),a
 flags=list([true,false,true])
 flags.sort()
 show flags
 empty=string_list()
 show empty.sort(),empty.contains("x")
 `
 want:="true true false\n[Ana, Carlos, Pedro]\n[-1, 0, 3, 3, 5]\ntrue [-1, 0, 3, 5]\n[false, true, true]\n0 false\n"
 for _,profile:=range []string{"nano","balanced","turbo","debug"}{t.Run(profile,func(t *testing.T){got,bad:=hyperRun(t,src,profile);if bad||got!=want{t.Fatalf("%q %v",got,bad)}})}
}
