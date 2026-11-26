package core

import (
 "bytes"
 "testing"
)

func TestWebEngineR1(t *testing.T){
 for _,source:=range []string{
  "show 42",
  "struct Player {name:string=\"Alex\";score:int=0}\np=Player()\np.score=100\nshow f\"{p.name}: {p.score}\"",
  "p={\"name\":\"Alex\",\"health\":100}\np[\"health\"]-=20\nshow p[\"name\"],p[\"health\"],len(p)",
  "show 1.25+2.5",
  "show float(9)/2.0",
  "name=\"TRIPOL\"\nshow f\"Hello {name}: {2+3}\"",
  "a=[1,2,3]\na[1]+=4\nshow a[1],len(a),sum(a)",
  "a=ints(3)\nfor i in range(3) {a[i]=i*2}\nshow a[2]",
  "a=[1,2,3]\nfor i,v in a {show i,v}",
  "a=list([1,2,3])\na.add(4)\nshow a.get(3),len(a)",
  "a=list([\"a\",\"b\"])\na.add(\"c\")\nshow a[2]",
  "func greet(name=\"Guest\")->string {return name}\nshow greet()",
  "name=input(\"Name: \")\nshow name",
  "show \"Hello, TRIPOL!\"",
  "name=\"TRIPOL\"\nshow \"Hello, \"+name",
  "func greet(name:string)->string {return \"Hi \"+name}\nshow greet(\"Ada\")",
  "show len(\"TRIPOL\")",
  "n=0\nfor i in range(1,6) {n+=i}\nshow n",
  "func double(n:int)->int {return n*2}\nshow double(21)",
  "show true and not false",
  "import tripol.math\nshow math.square(7)",
 }{
  b,err:=CompileWeb(source);if err!=nil{t.Fatalf("%s: %v",source,err)}
  if len(b)<8||!bytes.Equal(b[:8],[]byte{0,97,115,109,1,0,0,0}){t.Fatal("missing WASM header")}
 }
 for _,source:=range []string{"show missing","show missing2","import tripol.files\nshow 1"}{if _,err:=CompileWeb(source);err==nil{t.Fatalf("accepted %s",source)}}
}
