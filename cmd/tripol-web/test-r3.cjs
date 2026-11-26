'use strict';
const fs=require('node:fs'),assert=require('node:assert/strict'),cp=require('node:child_process'),os=require('node:os'),path=require('node:path');
require('./wasm_exec.js');const host=require('./host-r3.js');
(async()=>{
 const go=new Go();const {instance}=await WebAssembly.instantiate(fs.readFileSync(path.join(__dirname,'tripol-compiler.wasm')),go.importObject);go.run(instance);
 const cases=[
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
  fs.writeFileSync(path.join(__dirname,'parity-report.json'),JSON.stringify({total:rows.length,counts,rows},null,2));console.log('PARITY',counts)
 }
 process.exit(0)
})().catch(e=>{console.error(e);process.exit(1)});
