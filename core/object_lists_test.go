package core

import "testing"

const objectListsSource = `struct Point {x:int;y:int}
struct Map {rows:list[list[Point]]}
func make_map() {
    let empty:list[Point]=list()
    empty.add(Point(3,4))
    other=list([Point(1,2),Point(5,6)])
    return Map(list([empty,other]))
}
m=make_map()
rows=m.rows
alias=rows.get(0)
alias.add(Point(9,10))
show rows.get(0).get(1).x
show m.rows.get(1).get(0).y
for row in rows {
    for point in row {show point.x}
}
`

func TestSimpleR3ObjectAndNestedLists(t *testing.T) {
	runCase(t, objectListsSource, "9\n2\n3\n9\n1\n5\n", false)
	if _, err := CompileWeb(objectListsSource); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleR3DecodeRecord(t *testing.T) {
	const source = `import tripol.json
struct Point {x:int;ok:bool}
struct Map {name:string;rows:list[list[Point]]}
m=json.decode(Map,"{\"name\":\"world\",\"rows\":[[{\"\\u0078\":42,\"ok\":true}],[]]}")
show m.name,m.rows.get(0).get(0).x,m.rows.get(0).get(0).ok
show m.rows.get(1).size()
`
	runCase(t, source, "world 42 true\n0\n", false)
	if _, err := CompileWeb(source); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleR3DecodeRejectsDuplicateMissingAndWrongFields(t *testing.T) {
	const source = `import tripol.json
struct Point {x:int}
bad=list(["{}","{\"x\":1,\"\\u0078\":2}","{\"x\":null}","{\"x\":9223372036854775808}","{\"x\":1,}"])
for raw in bad {
    try {
        point=json.decode(Point,raw)
        show point.x
    } catch code {show code}
}
`
	runCase(t, source, "12\n12\n12\n12\n12\n", false)
}
