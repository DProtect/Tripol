'use strict';
let ready,build='';
async function compiler(){
 if(!ready)ready=(async()=>{
  const base='engine-r3/',fresh=Date.now().toString(36),r=await fetch(base+'manifest.json?refresh='+fresh,{cache:'no-store'});if(!r.ok)throw Error('Web Engine R3 is not published yet.');
  const m=await r.json();if(m.verified!==true||m.abi!==3||!/^[0-9a-f]{40}$/.test(m.commit))throw Error('Web Engine R3 build is not verified.');
  build=m.commit;
  const suffix='?v='+build+'&refresh='+fresh;
  importScripts(base+'wasm_exec.js'+suffix,base+'host-r3.js'+suffix);
  const go=new Go(),response=await fetch(base+'tripol-compiler.wasm'+suffix,{cache:'no-store'});if(!response.ok)throw Error('Compiler missing');
  const {instance}=await WebAssembly.instantiate(await response.arrayBuffer(),go.importObject);go.run(instance).catch(console.error);
  if(typeof self.tripolCompileWeb!=='function')throw Error('Compiler failed to initialize');
  const probe=self.tripolCompileWeb('import tripol.arrays\nshow arrays.total(arrays.make(4,7))');
  if(!probe.ok)throw Error('Outdated Web Engine build: '+build.slice(0,7)+'. Arrays-enabled build has not reached this Pages server. Reload the page.');
  if(TripolWebHost.execute(probe.wasm).output!=='28\n')throw Error('Web Engine asset mismatch; reload after deployment.');
 })();return ready;
}
self.onmessage=async({data})=>{try{
 await compiler();let t=performance.now();const r=self.tripolCompileWeb(data.source),compileMs=performance.now()-t;
 if(!r.ok)throw Error(r.error);
 const binary=r.wasm;if(!WebAssembly.validate(binary))throw Error('Invalid compiler output');
 if(data.compileOnly){self.postMessage({ok:true,output:'WebAssembly compiled: '+binary.length+' bytes\n',steps:0,compileMs,runMs:0,binary,build});return}
 t=performance.now();const result=TripolWebHost.execute(binary,{input:data.input});self.postMessage({ok:true,...result,compileMs,runMs:performance.now()-t,binary,build});
}catch(e){self.postMessage({ok:false,output:e.output||'',error:String(e.message||e),frames:Array.isArray(e.frames)?e.frames:[],build})}};
