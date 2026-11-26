package core
import "testing"
func TestSimpleTypedCollections(t *testing.T){
 src:=`names=list(["Ana","Carlos"])
 alias=names
 names.add("Pedro")
 names[0]="Eva"
 show alias
 show names.remove_at(1),names.get(0),len(names)
 for i,name in names {show f"{i}: {name}"}
 flags=bool_list()
 flags.add(true);flags.add(false)
 show flags,flags[0]
 values=list([1.5,2.5])
 show len(values),values[1]==2.5
 player={"name":"Alex","health":100,"active":true}
 other=player
 player["health"]-=20
 other["name"]="Eva"
 show player
 show player["health"],player["active"],len(player)
 `
 want:="[Eva, Carlos, Pedro]\nCarlos Eva 2\n0: Eva\n1: Pedro\n[true, false] true\n2 true\n{name: Eva, health: 80, active: true}\n80 true 3\n"
 for _,profile:=range []string{"nano","balanced","turbo","debug"}{t.Run(profile,func(t *testing.T){got,bad:=hyperRun(t,src,profile);if bad||got!=want{t.Fatalf("%q %v",got,bad)}})}
}
func TestSimpleCollectionSafety(t *testing.T){
 for _,s:=range []string{
 `a=list(["a"]);a.add(3)`,`a=bool_list();a.add(1)`,`a=list(["a"]);show sum(a)`,`a=list(["a"]);a[0]=1`,
 `a={"x":1,"x":2}`,`a={"x":1};a["x"]="bad"`,`a={"x":1};show a["z"]`,
 `a={"x":1};scratch {a["x"]=2}`,`a=string_list();scratch {a.add("bad")}`,
 }{t.Run(s,func(t *testing.T){defer func(){if recover()==nil{t.Fatal("accepted unsafe collection")}}();compile(s,"linux")})}
}
