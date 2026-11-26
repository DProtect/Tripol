'use strict';
// The Go compiler itself is a WASM module. Its output is a separate program module.
let ready;
async function compiler(){
 if(!ready)ready=(async()=>{
  const manifest=await fetch('engine-r1/manifest.json',{cache:'no-store'});
  if(!manifest.ok)throw Error('Web Engine has not been published. Use the interpreter.');
  const info=await manifest.json();if(info.verified!==true||info.abi!==1)throw Error('Web Engine build is not verified.');
  try{importScripts('engine-r1/wasm_exec.js')}catch{throw Error('Web Engine runtime unavailable. Use the interpreter until the build is published.')}
  const go=new Go();
  const response=await fetch('engine-r1/tripol-compiler.wasm');
  if(!response.ok)throw Error('Web Engine build is not published yet; use the educational interpreter.');
  const {instance}=await WebAssembly.instantiate(await response.arrayBuffer(),go.importObject);
  go.run(instance).catch(e=>{console.error(e)});
  if(typeof self.tripolCompileWeb!=='function')throw Error('Compiler failed to initialize');
 })();
 return ready;
}
self.onmessage=async({data})=>{
 let output='',steps=0,compileMs=0,runMs=0;
 try{
  await compiler();let t=performance.now();
  const result=self.tripolCompileWeb(data.source);
  compileMs=performance.now()-t;
  if(!result.ok)throw Error(result.error);
  const binary=result.wasm;
  if(!WebAssembly.validate(binary))throw Error('Compiler emitted invalid WebAssembly');
  if(data.compileOnly){self.postMessage({ok:true,output:'WebAssembly compiled: '+binary.length+' bytes\n',steps,compileMs,runMs,binary});return}
  const limit=100000,MIN=-(1n<<63n),MAX=(1n<<63n)-1n;
  const imports={tripol:{
   tick(line){if(++steps>limit)throw Error('Line '+line+': Execution limit reached')},
   arithmetic(op,a,b,line){let v;if(op===0)v=a+b;else if(op===1)v=a-b;else if(op===2)v=a*b;else if(op===3||op===4){if(b===0n||a===MIN&&b===-1n)throw Error('Line '+line+': Invalid integer arithmetic');v=op===3?a/b:a%b}else throw Error('Unknown arithmetic operation');if(v<MIN||v>MAX)throw Error('Line '+line+': Integer overflow');return v},
   show(v,kind,last){const text=(kind?String(v!==0n):String(v))+(last?'\n':' ');if(output.length+text.length>32768)throw Error('Console output limit reached');output+=text},
   assertFail(line){throw Error('Line '+line+': Assertion failed')}
  }};
  t=performance.now();const program=new WebAssembly.Instance(new WebAssembly.Module(binary),imports);program.exports.run();runMs=performance.now()-t;
  self.postMessage({ok:true,output,steps,compileMs,runMs,binary});
 }catch(e){self.postMessage({ok:false,output,error:String(e.message||e)})}
};
