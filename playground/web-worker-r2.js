'use strict';
let ready;
async function compiler(){
 if(!ready)ready=(async()=>{
  const base='engine-r2/',r=await fetch(base+'manifest.json',{cache:'no-store'});if(!r.ok)throw Error('Web Engine R2 is not published yet.');
  const m=await r.json();if(m.verified!==true||m.abi!==2||!/^[0-9a-f]{40}$/.test(m.commit))throw Error('Web Engine R2 build is not verified.');
  const suffix='?v='+m.commit;
  importScripts(base+'wasm_exec.js'+suffix,base+'host-r2.js'+suffix);
  const go=new Go(),response=await fetch(base+'tripol-compiler.wasm'+suffix);if(!response.ok)throw Error('Compiler missing');
  const {instance}=await WebAssembly.instantiate(await response.arrayBuffer(),go.importObject);go.run(instance).catch(console.error);
  if(typeof self.tripolCompileWeb!=='function')throw Error('Compiler failed to initialize');
 })();return ready;
}
self.onmessage=async({data})=>{try{
 await compiler();let t=performance.now();const r=self.tripolCompileWeb(data.source),compileMs=performance.now()-t;
 if(!r.ok)throw Error(r.error);
 const binary=r.wasm;if(!WebAssembly.validate(binary))throw Error('Invalid compiler output');
 if(data.compileOnly){self.postMessage({ok:true,output:'WebAssembly compiled: '+binary.length+' bytes\n',steps:0,compileMs,runMs:0,binary});return}
 t=performance.now();const result=TripolWebHost.execute(binary);self.postMessage({ok:true,...result,compileMs,runMs:performance.now()-t,binary});
}catch(e){self.postMessage({ok:false,output:e.output||'',error:String(e.message||e)})}};
