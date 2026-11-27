'use strict';
// Audit the actual public documentation, not handwritten substitute probes.
const fs=require('node:fs'),cp=require('node:child_process'),path=require('node:path');
const decode=s=>s.replace(/&(?:amp|lt|gt|quot|apos|#\d+|#x[\da-f]+);/gi,m=>{const e=m.slice(1,-1);if(e[0]==='#'){const n=e[1].toLowerCase()==='x'?parseInt(e.slice(2),16):Number(e.slice(1));return n<=0x10ffff?String.fromCodePoint(n):m}return({amp:'&',lt:'<',gt:'>',quot:'"',apos:"'"})[e]??m});
const ref=process.argv[2]||'origin/gh-pages';
const files=cp.execFileSync('git',['ls-tree','-r','--name-only',ref,'docs/libraries'],{encoding:'utf8'}).trim().split('\n').filter(p=>p.endsWith('.html'));
const snippets=[];
for(const file of files){
 const html=cp.execFileSync('git',['show',ref+':'+file],{encoding:'utf8',maxBuffer:4*1024*1024});
 // Match explicit TRIPOL blocks, excluding function signatures and implementation source.
 const main=html.split(/<section\s+class="function-card"/i)[0];
 let n=0;
 for(const m of main.matchAll(/<pre\b([^>]*)>\s*<code\b[^>]*class=["'][^"']*\blanguage-tripol\b[^"']*["'][^>]*>([\s\S]*?)<\/code>\s*<\/pre>/gi)){
  const source=decode(m[2]);if(source.trim())snippets.push({file,index:++n,source,native:/\bdata-native-only\b/i.test(m[1])})
 }
}
require('./wasm_exec.js');const host=require('./host-r3.js');
(async()=>{
 const go=new Go();const {instance}=await WebAssembly.instantiate(fs.readFileSync(path.join(__dirname,'tripol-compiler.wasm')),go.importObject);go.run(instance);
 const rows=[];
 for(const {file,index,source,native} of snippets){
  const row={file,index,source};
  if(native){rows.push({...row,status:'native-only'});continue}
  try{
   const compiled=tripolCompileWeb(source);
   if(!compiled.ok){rows.push({...row,status:'compiler-rejected',diagnostic:compiled.error});continue}
   if(!WebAssembly.validate(compiled.wasm)){rows.push({...row,status:'invalid-wasm'});continue}
   try{const r=host.execute(compiled.wasm);rows.push({...row,status:'executed',output:r.output})}
   catch(e){rows.push({...row,status:'runtime-error',diagnostic:e.message})}
  }catch(e){rows.push({...row,status:'audit-error',diagnostic:e.message})}
 }
 const counts={};for(const r of rows)counts[r.status]=(counts[r.status]||0)+1;
 const report={commit:process.env.GITHUB_SHA||'',pagesRef:ref,pages:files.length,total:rows.length,counts,rows};
 fs.writeFileSync(path.join(__dirname,'docs-parity-report.json'),JSON.stringify(report,null,2));
 console.log('DOCS AUDIT',JSON.stringify({pages:report.pages,total:report.total,counts}));process.exit(rows.some(row=>!['executed','native-only'].includes(row.status))?1:0)
})().catch(e=>{console.error(e);process.exit(1)});
