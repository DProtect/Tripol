'use strict';
const fs=require('node:fs'),assert=require('node:assert/strict'),cp=require('node:child_process'),os=require('node:os'),path=require('node:path');
require('./wasm_exec.js');const host=require('./host-r3.js');
(async()=>{
 const go=new Go();const {instance}=await WebAssembly.instantiate(fs.readFileSync(path.join(__dirname,'tripol-compiler.wasm')),go.importObject);go.run(instance);
 const cases=[
 ['import tripol.files\nfiles.write("demo.txt","hello")\nfiles.append("demo.txt","!")\nshow files.read("demo.txt")','hello!\n'],
 ['import tripol.io\na=ints(2)\na[0]=255\nio.write("demo.bin",a)\nshow io.read("demo.bin")','[255, 0]\n'],
 ['import tripol.atomic\na=ints(1)\nshow atomic.add(a,0,4),atomic.compare_exchange(a,0,4,9),atomic.load(a,0)','4 true 9\n'],
 ['import tripol.system\nshow system.platform(),system.architecture(),system.is_linux(),system.is_windows()','web wasm32 false false\n'],
 ['import tripol.memory\na=ints(2)\nshow memory.capacity()==6291456, memory.used()>=16, memory.available()==memory.capacity()-memory.used()','true true true\n'],
 ['import tripol.audio\na=audio.square_wave(440,1)\nshow len(audio.wav_mono(a))','88\n'],
 ['import tripol.datetime\nshow datetime.leap(2024),datetime.leap(1900),datetime.days_in_month(2024,2)','true false 29\n'],
 ['import tripol.text\nshow text.size("é😀"),text.index_of("éabc","a"),text.part("é",0,2),text.to_int("-42")','6 2 é -42\n'],
 ['import tripol.json\nshow json.quote("é😀"),json.valid("{\"a\":1}")','"é😀" true\n'],
 ['import tripol.arrays\nvalues = arrays.make(4, 7)\nshow arrays.total(values)\ncopy = ints(4)\nassert arrays.copy_into(copy, values)\nshow copy','28\n[7, 7, 7, 7]\n'],
 ['import tripol.arrays\na=arrays.make(3,9)\nb=arrays.copy(a)\nb[0]=3\nshow arrays.contains(a,9),a[0],b[0]','true 9 3\n'],
 ['import tripol.collections\na=ints(2)\na[0]=5\na[1]=6\nshow collections.append(a,7)','[5, 6, 7]\n'],
 ['struct Player {name:string="Alex";score:int=0}\np=Player()\np.score=100\nshow f"{p.name}: {p.score}"','Alex: 100\n'],
 ['p={"name":"Alex","health":100}\np["health"]-=20\nshow p["name"],p["health"],len(p)','Alex 80 2\n'],
 ['show "Hello, TRIPOL!"','Hello, TRIPOL!\n'],
 ['show 42','42\n'],
 ['show true and not false','true\n'],
 ['show 1.25+2.5','3.750000\n'],
 ['show float(9)/2.0','4.500000\n'],
 ['show 1.5<2.0, 2.0==2.0','true true\n'],
 ['name="TRIPOL"\nshow f"Hello {name}: {2+3}"','Hello TRIPOL: 5\n'],
 ['name="TRIPOL"\nshow "Hello "+name','Hello TRIPOL\n'],
 ['show len("á")','2\n'],
 ['a=[1,2,3]\na[1]+=4\nshow a[1],len(a),sum(a)','6 3 10\n'],
 ['a=ints(3)\nfor i in range(3) {a[i]=i*2}\nshow a[2]','4\n'],
 ['a=[1,2,3]\nfor i,v in a {show i,v}','0 1\n1 2\n2 3\n'],
 ['a=list([1,2,3])\na.add(4)\nshow a.get(3),len(a)','4 4\n'],
 ['a=list(["a","b"])\na.add("c")\nshow a[2]','c\n'],
 ['func greet(name="Guest")->string {return name}\nshow greet()','Guest\n'],
 ['name=input("Name: ")\nshow name','Name: \nAda\n','Ada'],
 ['n=0\nwhile n<3 {n+=1}\nshow n','3\n'],
 ['for i in range(5) {if i==2 {continue}; if i==4 {break}; show i}','0\n1\n3\n'],
 ['func twice(n) => n*2\nshow twice(21)','42\n'],
 ['import tripol.math\nshow math.square(7)','49\n']
 ];
 function execute(src,input){const r=tripolCompileWeb(src);if(!r.ok)throw Error(r.error);return host.execute(r.wasm,{input}).output}
 for(const [s,w,i] of cases)assert.equal(execute(s,i),w,s);
 assert.throws(()=>execute('show 9223372036854775807 + 1'),/overflow/i);
 assert.throws(()=>execute('a=ints(1)\nshow a[1]'),/bounds/i);
 assert.equal(tripolCompileWeb('show missing').ok,false);
 assert.throws(()=>execute('import tripol.files\nshow files.read("demo.txt")'),/not found/i);
 assert.throws(()=>execute('import tripol.files\nfiles.write("../escape.txt","x")'),/relative virtual/i);
 assert.throws(()=>execute('import tripol.io\na=ints(1)\na[0]=256\nio.write("bad.bin",a)'),/invalid binary/i);
 console.log('PASS required '+(cases.length+3)+' R3 checks');
 // Native oracle: compare stdout from native x86-64 with actual WASM.
 // Report every disagreement (including unsupported examples), never substitute
 // the educational interpreter's result for the native compiler.
 const exampleFile=process.argv[2],native=process.argv[3];
 if(exampleFile&&native){
  const examples=require(path.resolve(exampleFile));const dir=fs.mkdtempSync(path.join(os.tmpdir(),'tripol-parity-'));let rows=[];
  try{for(const [name,es,source,want,input] of examples){
   const file=path.join(dir,'main.tri'),bin=path.join(dir,'program');fs.writeFileSync(file,source);
   const build=cp.spawnSync(native,['build',file,'--target','linux','-o',bin],{encoding:'utf8',timeout:10000});
   if(build.status!==0){rows.push({name,status:'native-rejected',detail:String(build.stderr||build.error||'').slice(0,500)});continue}
   const actual=cp.spawnSync(bin,[],{encoding:'utf8',input:(input||'')+'\n',timeout:2000});
   if(actual.status!==0){rows.push({name,status:'native-runtime-error',detail:String(actual.stderr||actual.error||'').slice(0,500)});continue}
   try{const wasm=execute(source,input);rows.push({name,status:wasm===actual.stdout?'match':'mismatch',native:actual.stdout,wasm})}catch(e){rows.push({name,status:'wasm-rejected',detail:e.message})}
  }}finally{fs.rmSync(dir,{recursive:true,force:true})}
  const counts={};for(const r of rows)counts[r.status]=(counts[r.status]||0)+1;
  fs.writeFileSync(path.join(__dirname,'parity-report.json'),JSON.stringify({total:rows.length,counts,rows},null,2));console.log('PARITY',counts);assert.equal(rows.length,100,'Expected all Playground examples');assert.equal(counts.match,100,'Native/WebAssembly parity must be 100/100 before publishing')

  // Library probes cover the candidate modules, not every function. Keep the
  // full report even when a module has partial/unsupported browser operations.
  const probes={"math":"show math.square(7)","arrays":"a=arrays.make(4,7)\nshow arrays.total(a)","collections":"a=ints(2)\na[0]=7\nshow collections.append(a,9)","datetime":"show datetime.leap(2024),datetime.leap(1900),datetime.days_in_month(2024,2)","bits":"show bits.set(0,5),bits.count(31)","bitset":"a=bitset.create(100)\nshow bitset.add(a,65),bitset.contains(a,65),bitset.count(a)","binary":"show binary.u16_le(513)","buffer":"b=buffer.create(1)\nb=buffer.push(b,10)\nb=buffer.push(b,20)\nshow buffer.values(b)","checksum":"show checksum.adler32(\"abc\")","colors":"show colors.rgb(255,0,0),colors.invert(0)","compression":"a=ints(3)\na[0]=7\na[1]=7\na[2]=7\nshow compression.rle_decode(compression.rle_encode(a))","console":"show console.line(\"hello\")","convert":"show convert.celsius_to_fahrenheit(20.0)","dictionary":"d=dictionary.create(2)\nshow dictionary.set(d,\"key\",42),dictionary.get(d,\"key\")","encoding":"show encoding.hex_text(\"é\")","geometry":"show geometry.distance(0.0,0.0,3.0,4.0)","json":"show json.quote(\"é😀\"),json.valid(\"{\\\"a\\\":1}\")","matrix":"show matrix.identity(2)","path":"show path.basename(\"a/b.txt\"),path.extension(\"a/b.txt\")","physics":"show physics.integrate(1.0,2.0,3.0)","random":"show random.next(1)","search":"a=ints(2)\na[1]=7\nshow search.linear(a,7)","sequences":"show sequences.fibonacci(8)","simd":"a=ints(2)\nb=ints(2)\nshow simd.equal(a,b)","sort":"a=ints(3)\na[0]=9\na[1]=2\na[2]=1\nshow sort.ascending(a)","stats":"a=ints(2)\na[0]=2\na[1]=4\nshow stats.mean(a),stats.variance(a)","text":"show text.size(\"é😀\"),text.index_of(\"éabc\",\"a\"),text.part(\"é\",0,2),text.to_int(\"-42\")","validate":"show validate.integer(\"-12\"),validate.ascii(\"é\")","vector":"show vector.length2(3.0,4.0)","stack":"s=stack.create(2)\nshow stack.push(s,7),stack.pop(s)","queue":"q=queue.create(2)\nshow queue.push(q,9),queue.pop(q)","animation":"show animation.smoothstep(0.5)","camera":"show camera.screen(12,2,3)","particles":"p=particles.create(2)\nshow particles.burst(p,1,2),particles.step(p)","canvas":"p=canvas.create(3,2)\nshow canvas.clear(p,0),canvas.rectangle(p,3,2,1,0,2,1,255)\nshow p","sprite":"p=sprite.solid(2,2,255)\nshow sprite.rotate90(p,2,2)","testing":"show testing.equal_int(4,4)","errors":"show errors.int_or(\"42\",0),errors.int_or(\"bad\",9)","network":"show network.ipv4(127,0,0,1)","http":"show http.status(\"HTTP/1.0 200 OK\"),http.body(\"HTTP/1.0 200 OK\\r\\n\\r\\nhello\")","image":"p=ints(2)\np[0]=1\np[1]=2\nshow image.flip_x(p,2,1)","time":"show time.seconds(2500),time.now_ms()>=0","profiler":"s=profiler.start()\nshow profiler.elapsed(s)>=0","heap":"h=heap.create(4)\nheap.push(h,8)\nheap.push(h,2)\nshow heap.pop(h)","utf8":"show utf8.count(\"é😀\"),utf8.scalar_at(\"é😀\",2)","csv":"show csv.fields(\"a,b,c\"),csv.escape(\"é\")","stream":"c=stream.create()\na=ints(2)\na[0]=1\na[1]=2\nshow stream.read_u16(c,a)","regex":"show regex.matches(\"a*b\",\"aaab\")","cli":"show cli.integer(\"--n=7\",\"n\",0)","crypto":"show crypto.sha256_text(\"abc\")"};
  const libDir=fs.mkdtempSync(path.join(os.tmpdir(),'tripol-library-'));const libRows=[];
  try{for(const [name,body] of Object.entries(probes)){
   const source='import tripol.'+name+'\n'+body,file=path.join(libDir,'main.tri'),bin=path.join(libDir,'program');fs.writeFileSync(file,source);
   const build=cp.spawnSync(native,['build',file,'--target','linux','-o',bin],{encoding:'utf8',timeout:10000});
   if(build.status!==0){libRows.push({name,status:'native-rejected',detail:String(build.stderr||build.stdout||build.error||'').slice(0,500)});continue}
   const actual=cp.spawnSync(bin,[],{encoding:'utf8',timeout:2000});
   if(actual.status!==0){libRows.push({name,status:'native-runtime-error',detail:String(actual.stderr||actual.error||'').slice(0,500)});continue}
   try{const wasm=execute(source);libRows.push({name,status:wasm===actual.stdout?'match':'mismatch',native:actual.stdout,wasm})}catch(e){libRows.push({name,status:'wasm-rejected',detail:e.message})}
  }}finally{fs.rmSync(libDir,{recursive:true,force:true})}
  const libraryCounts={};for(const r of libRows)libraryCounts[r.status]=(libraryCounts[r.status]||0)+1;
  fs.writeFileSync(path.join(__dirname,'library-parity-report.json'),JSON.stringify({total:libRows.length,counts:libraryCounts,rows:libRows},null,2));console.log('LIBRARY PARITY',libraryCounts);

 }
 process.exit(0)
})().catch(e=>{console.error(e);process.exit(1)});
