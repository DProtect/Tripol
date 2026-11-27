package core

import (
 "bytes"
 "testing"
)

func TestWebEngineR1(t *testing.T){
 for _,source:=range []string{
  "import tripol.game\nimport tripol.canvas\n\ngame.window(\"Preview\", 2, 2)\np = canvas.create(2, 2)\ncanvas.clear(p, 16711680)\ncanvas.pixel(p, 2, 2, 1, 1, 255)\nshow game.present(p, 2, 2)\ngame.close()",
  "import tripol.game\n\ngame.window(\"Triangle\", 32, 32)\nv = ints(9)\nv[0] = -800\nv[1] = -800\nv[2] = 16711680\nv[3] = 800\nv[4] = -800\nv[5] = 65280\nv[6] = 0\nv[7] = 800\nv[8] = 255\nshow game.gl_frame(v)\ngame.close()",
  "struct Counter {value:int=0}\nfunc Counter.add(self:Counter,step=1)->int {self.value=self.value+step;return self.value}\nshow Counter().add()",
  "import tripol.discord\nb=discord.bot()\nb.command(\"hola\",\"Hola!\")\nshow b.preview(\"hola\")",
  "import tripol.discord\ne=discord.embed()\ne.title(\"Hola\")\nshow e.payload()",
  "show 42",
  "import tripol.collections\nitems=[1,2]\nvalues=collections.append(items,3)\nshow values",
  "import tripol.arrays\nitems=[1,2,3]\nshow arrays.total(items)",
  "import tripol.collections\na=[1,2]\nb=[3,4]\nshow collections.concat(a,b)",
  "import tripol.files\nfiles.write(\"demo.txt\",\"hello\")\nshow files.read(\"demo.txt\")",
  "import tripol.io\na=ints(2)\na[0]=255\nio.write(\"demo.bin\",a)\nshow io.read(\"demo.bin\")",
  "import tripol.atomic\na=ints(1)\nshow atomic.add(a,0,4),atomic.compare_exchange(a,0,4,9),atomic.load(a,0)",
  "import tripol.system\nshow system.platform(),system.architecture()",
  "import tripol.memory\nshow memory.capacity(),memory.used()",
  "import tripol.audio\na=audio.square_wave(440,1)\nshow len(audio.wav_mono(a))",
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
 for _,source:=range []string{"show missing","show missing2","import tripol.env\nshow 1"}{if _,err:=CompileWeb(source);err==nil{t.Fatalf("accepted %s",source)}}
}
