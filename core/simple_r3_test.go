package core

import "testing"

func TestSimpleR3OverloadsNestedRecordsAndVoidHandlers(t *testing.T) {
	const source = `import tripol.json
struct Point {x:int}
struct Payload {point:Point}
func pick(x:int)->int {return x+1}
func pick(x:string)->string {return x+"!"}
func report(item:Payload,fn:func(Payload)->void)->void {
 fn(item)
}
func serialize(object:json.Object)->string {return object.encode()}
x=Payload(Point(42))
report(x,func(p) {show p.point.x})
show pick(3),pick("yes")
show serialize({"ok":true,"meta":{"point":x.point.x}})
`
	runCase(t, source, "42\n4 yes!\n{\"ok\":true,\"meta\":{\"point\":42}}\n", false)
	if _, err := CompileWeb(source); err != nil {
		t.Fatal(err)
	}
}
