'use strict';
(()=>{
const $=s=>document.querySelector(s), $$=s=>[...document.querySelectorAll(s)];
let lang='en';try{if(localStorage.getItem('tripol-language')==='es')lang='es'}catch{}
const chosen=new URLSearchParams(location.search).get('lang');if(chosen==='en'||chosen==='es')lang=chosen;
const root=document.body.dataset.docRoot||'./';
// A freshly loaded script must also load its matching stylesheet.
for(const link of document.querySelectorAll('link[rel="stylesheet"]')){const u=new URL(link.href,location.href);if(u.pathname.endsWith('/docs/docs.css')&&u.searchParams.get('v')!=='live-r1'){u.searchParams.set('v','live-r1');link.href=u.href}}
function localLinks(){for(const a of $$('a[href]')){const raw=a.getAttribute('href');if(!raw||raw.startsWith('#')||/^(https?:|mailto:)/i.test(raw)&&a.origin!==location.origin)continue;const u=new URL(raw,location.href);if(!u.pathname.endsWith('.html'))continue;u.searchParams.set('lang',lang);a.href=u.href;}}
// Shared navigation and current release label on all reference pages.
const side=document.querySelector('#doc-side');
side.querySelector('.side-kicker').textContent='OVERDRIVE / 0.26.0';
const crumb=document.querySelector('.crumb > span:last-child');if(crumb&&/^0\.(24|25)$/.test(crumb.textContent))crumb.textContent='0.26';
const meta=document.querySelector('.lib-meta > span:first-child');if(meta&&/^(ASCENSION 0\.24|EVOLUTION 0\.25)$/.test(meta.textContent))meta.textContent='OVERDRIVE 0.26';
const guide=side.querySelector('nav[aria-label="Guides"]');
for(const [file,label] of [['genesis','GENESIS R2'],['update','UPDATE R1']]){let a=guide.querySelector('a[href$="'+file+'.html"]');if(!a){a=document.createElement('a');a.href=root+file+'.html';guide.append(a);}a.textContent=label;}
const libs=side.querySelector('.side-libs');
for(const [name,before,count] of [['queue','random',7],['stack','stats',7],['bitset','buffer',6],['profiler','queue',2]]){if(!libs.querySelector('a[href$="/'+name+'.html"]')){const a=document.createElement('a');a.href=root+'libraries/'+name+'.html';const label=document.createElement('span');label.textContent=name;const n=document.createElement('small');n.textContent=String(count);a.append(label,n);libs.insertBefore(a,libs.querySelector('a[href$="/'+before+'.html"]'));}}
for(const a of side.querySelectorAll('nav a')){a.removeAttribute('aria-current');if(new URL(a.href,location.href).pathname===location.pathname)a.setAttribute('aria-current','page');}
const translate=(en,es)=>lang==='es'?es:en;
let labels=[], placeholders=[];
function registerLabels(){labels=$$('[data-es]').map(el=>({el,en:el.textContent,es:el.dataset.es}));placeholders=$$('[data-placeholder-es]').map(el=>({el,en:el.placeholder,es:el.dataset.placeholderEs}));}
function docsTokens(source,mode='tripol'){
 const keywords=new Set('show let const func return if else while for in repeat break continue assert and or not import struct try catch finally throw scratch'.split(' '));
 const types=new Set('int float bool string ints list string_list bool_list float_list'.split(' '));
 const builtins=new Set('input input_int len sum range abs min max sqrt text slice parse_int'.split(' '));
 const terminal=['cmd','sh','bash','shell','powershell','ps1'].includes(mode);
 const pattern=terminal?/(?:https?:\/\/[^\s"'<>]+|#[^\r\n]*|"(?:\\.|[^"\\\r\n])*"?|'(?:\\.|[^'\\\r\n])*'?|\$[\w:]+|--?[\w][\w-]*|\b[\w.-]+(?:[\/\\][\w./\\-]+)+|\b\d+(?:\.\d+)?\b|[A-Za-z_][A-Za-z_0-9.-]*|[=|&<>]+|\s+|.)/gu:/(?:\/\/[^\r\n]*|#[^\r\n]*|f?"(?:\\.|[^"\\\r\n])*"?|'(?:\\.|[^'\\\r\n])*'?|\b(?:0[xX][\da-fA-F]+|\d+(?:\.\d+)?(?:[eE][+-]?\d+)?)\b|[A-Za-z_][A-Za-z_0-9]*|==|!=|<=|>=|=>|->|\+=|-=|\*=|\/=|%=|&&|\|\||[+\-*/%=<>!]|[{}()[\],.:;]|\s+|.)/gu;
 const out=[];const emit=(v,k='plain')=>{if(v)out.push({v,k})};
 for(const m of source.matchAll(pattern)){
  const v=m[0],after=source.slice(m.index+v.length);let k='plain';
  if(v.startsWith('#')||(!terminal&&v.startsWith('//')))k='comment';
  else if(!terminal&&v.startsWith('f"')){
   let start=0,i=2;
   while(i<v.length){
    if(v[i]==='\\'){i+=2;continue}
    if(v.slice(i,i+2)==='{{'||v.slice(i,i+2)==='}}'){i+=2;continue}
    if(v[i]==='{'){const end=v.indexOf('}',i+1);if(end<0)break;const field=v.slice(i+1,end);if(!field.includes('{')&&!field.includes('\\')&&!field.includes('"')){emit(v.slice(start,i),'string');emit('{','interpolation');out.push(...docsTokens(field,'tripol'));emit('}','interpolation');start=end+1;i=start;continue}}
    i++;
   }
   emit(v.slice(start),'string');continue;
  }
  else if(v.startsWith('"')||v.startsWith("'"))k='string';
  else if(/^\d/.test(v))k='number';
  else if(terminal){
   if(/^https?:\/\//.test(v))k='path';
   else if(v.startsWith('$')||/^[A-Z_][A-Z_0-9]+$/.test(v)&&after.startsWith('='))k='variable';
   else if(/^--?[\w]/.test(v))k='option';
   else if(/^(tripol|git|go|python|python3|cd|mkdir|cp|curl|node|rustc|cargo|npm|cmake|notepad|powershell|cmd|Read-Host|Get-ChildItem|Set-Location|echo|chmod|make)$/i.test(v))k='function';
   else if(/[\/\\]/.test(v))k='path';
   else if(/^[=|&<>]+$/.test(v))k='operator';
  }else{
   if(v==='true'||v==='false'||v==='null')k='literal';
   else if(keywords.has(v))k='keyword';
   else if(types.has(v))k='type';
   else if(builtins.has(v)||/^[A-Za-z_]/.test(v)&&/^\s*\(/.test(after))k='function';
   else if(v==='tripol'||v==='vault'||/^[A-Za-z_]/.test(v)&&/^\s*\./.test(after))k='namespace';
   else if(/^[+\-*/%=<>!&|]/.test(v))k='operator';
   else if(/^[{}()[\],.:;]$/.test(v))k='punct';
  }
  emit(v,k);
 }
 return out;
}
function codeMode(code){
 const pre=code.closest('pre');const explicit=[...code.classList,...(pre?[...pre.classList]:[])].find(c=>c.startsWith('language-'));
 if(explicit)return explicit.slice(9);
 if(pre?.classList.contains('signature'))return 'tripol';
 const s=code.textContent.trim();
 if(/^(?:tripol|git|go|python3?|cd|rustc|cargo|npm|cmake|curl|mkdir|cp|notepad|powershell)\s|^--?[\w-]|^\$[\w:]+|^[A-Z_][A-Z_0-9]+=/.test(s))return 'sh';
 if(/^(?:func|import|show|let|const|struct|if|for|while|try|scratch)\b|^[A-Za-z_]\w*\.[A-Za-z_]\w*|[=()[\]{}]/.test(s))return 'tripol';
 return pre?'text':'tripol';
}
function highlightDocs(){
 for(const code of [...document.querySelectorAll('.docs-article code, .docs-article pre:not(:has(code))')]){
  if(code.closest('.doc-live-panel'))continue;
  const raw=code.textContent,mode=codeMode(code);
  if(raw.length>100000||!raw||['text','output','plaintext'].includes(mode))continue;
  if(code.dataset.syntaxReady==='yes'&&code._tripolSource===raw&&code.querySelector('.doc-tok'))continue;
  if([...code.children].some(c=>!c.classList.contains('doc-tok')))continue;
  const frag=document.createDocumentFragment();
  for(const t of docsTokens(raw,mode)){if(t.k==='plain'){frag.append(document.createTextNode(t.v));continue}const span=document.createElement('span');span.className='doc-tok doc-tok-'+t.k;span.textContent=t.v;frag.append(span)}
  code.replaceChildren(frag);code.dataset.syntaxReady='yes';code._tripolSource=raw;
 }
}


// DOCS LIVE R1 — lazy editors, real compiler workers, no automatic execution.
const livePanels=[];let liveJob=null;
function liveElement(tag,cls,text){const e=document.createElement(tag);if(cls)e.className=cls;if(text!==undefined)e.textContent=text;return e}
function liveLabel(e,en,es){e.textContent=translate(en,es);e._liveLabel=[en,es];return e}
function liveButton(en,es,action,cls='doc-copy'){const b=liveLabel(liveElement('button',cls),en,es);b.type='button';b.addEventListener('click',action);return b}
function liveStatus(p,en,es,kind=''){p.message=[en,es];p.status.textContent=translate(en,es);p.status.dataset.state=kind}
function liveRefresh(){for(const e of $$('.doc-live-action,.doc-live-panel button,.doc-live-label,.doc-live-help'))if(e._liveLabel)e.textContent=translate(...e._liveLabel);for(const p of livePanels){if(p.message)p.status.textContent=translate(...p.message);p.input.placeholder=translate('One line per input()','Una línea por input()');p.editor.setAttribute('aria-label',translate('Editable TRIPOL example','Ejemplo TRIPOL editable'));p.input.setAttribute('aria-label',translate('Program input','Entrada del programa'));p.output.setAttribute('aria-label',translate('Program output','Salida del programa'))}}
function livePlayground(source){if(source.length>32768){toast(translate('Source exceeds 32 KiB','El código supera 32 KiB'));return}const u=new URL(root+'../playground/index.html',location.href);u.searchParams.set('lang',lang);u.hash='code='+encodeURIComponent(source);window.open(u.href,'_blank','noopener,noreferrer')}
function livePaint(p){const frag=document.createDocumentFragment();for(const t of docsTokens(p.editor.value)){if(t.k==='plain')frag.append(document.createTextNode(t.v));else{const el=liveElement('span','doc-tok doc-tok-'+t.k,t.v);frag.append(el)}}frag.append(document.createTextNode('\n'));p.highlight.replaceChildren(frag);p.stack.classList.add('painted');liveScroll(p)}
function liveScroll(p){p.highlight.scrollTop=p.editor.scrollTop;p.highlight.scrollLeft=p.editor.scrollLeft}
function liveStop(p,message=true){if(!liveJob||liveJob.panel!==p)return;const j=liveJob;liveJob=null;clearTimeout(j.timer);j.worker.terminate();p.run.disabled=false;p.stop.hidden=true;if(message)liveStatus(p,'Stopped','Detenido')}
function liveFiles(p,files=[]){
 p.files.replaceChildren();p.files.hidden=true;
 for(const file of files.slice(0,64)){
  if(typeof file.name!=='string'||!(file.bytes instanceof Uint8Array))continue;
  const row=liveElement('div','doc-live-file-row');row.append(liveElement('span','',file.name+' · '+file.bytes.length+' B'));
  row.append(liveButton('Download','Descargar',()=>{const u=URL.createObjectURL(new Blob([file.bytes],{type:file.name.endsWith('.wav')?'audio/wav':'application/octet-stream'}));const a=document.createElement('a');a.href=u;a.download=file.name.split('/').at(-1);a.click();setTimeout(()=>URL.revokeObjectURL(u),1000)}));p.files.append(row)
 }
 p.files.hidden=!p.files.children.length
}
function liveRun(p){liveFiles(p);
 if(liveJob)liveStop(liveJob.panel);
 const source=p.editor.value;
 if(new TextEncoder().encode(source).length>32768||new TextEncoder().encode(p.input.value).length>32768){liveStatus(p,'Input exceeds 32 KiB','La entrada supera 32 KiB','error');return}
 p.output.textContent='';p.run.disabled=true;p.stop.hidden=false;liveStatus(p,'Compiling with Web Engine R3…','Compilando con Web Engine R3…');
 let worker;try{worker=new Worker(new URL(root+'../playground/web-worker-r3.1.js',location.href))}catch(e){p.run.disabled=false;p.stop.hidden=true;liveStatus(p,'Compiler worker unavailable','Worker del compilador no disponible','error');p.output.textContent=String(e.message||e);return}
 const job={panel:p,worker,timer:null};liveJob=job;
 worker.onmessage=({data})=>{if(liveJob!==job)return;liveStop(p,false);p.output.textContent=(data.output||'')+(data.ok?'':'\n'+translate('Error: ','Error: ')+String(data.error||'Unknown error'));if(data.ok){if(Array.isArray(data.files))liveFiles(p,data.files);const time=Number(data.compileMs||0).toFixed(1)+' ms'+(data.build?' · '+data.build.slice(0,7):'');liveStatus(p,'WASM R3 · compiled in '+time,'WASM R3 · compilado en '+time,'ok')}else liveStatus(p,'Compiler/runtime diagnostic · no fallback','Diagnóstico del compilador/runtime · sin sustitución','error')};
 worker.onerror=()=>{if(liveJob!==job)return;liveStop(p,false);liveStatus(p,'Web Engine unavailable','Web Engine no disponible','error');p.output.textContent=translate('The compiler worker could not start. Open in Playground or try again.','No se pudo iniciar el compilador. Abre en Playground o vuelve a intentar.')};
 job.timer=setTimeout(()=>{if(liveJob!==job)return;liveStop(p,false);liveStatus(p,'Execution time limit (15 s)','Límite de ejecución (15 s)','error')},15000);
 try{worker.postMessage({source,input:p.input.value,compileOnly:false})}catch(e){liveStop(p,false);liveStatus(p,'Could not send source','No se pudo enviar el código','error');p.output.textContent=String(e.message||e)}
}
function makeLivePanel(source,original,trigger){
 const p={original,source};const panel=liveElement('section','doc-live-panel');p.panel=panel;panel.hidden=true;
 const toolbar=liveElement('div','doc-live-toolbar');const filename=liveElement('span','doc-live-file','example.tri');const actions=liveElement('div','doc-live-tools');toolbar.append(filename,actions);
 p.run=liveButton('▶ Run','▶ Ejecutar',()=>liveRun(p),'doc-live-run');
 p.stop=liveButton('Stop','Detener',()=>liveStop(p));p.stop.hidden=true;
 const reset=liveButton('Reset','Restaurar',()=>{liveStop(p,false);p.editor.value=p.original;p.input.value='';liveFiles(p);p.output.textContent='';livePaint(p);liveStatus(p,'Ready','Listo')});
 const open=liveButton('Playground ↗','Playground ↗',()=>livePlayground(p.editor.value));
 const close=liveButton('Close','Cerrar',()=>{liveStop(p);panel.hidden=true;source.hidden=false;trigger.setAttribute('aria-expanded','false');trigger.focus()});
 actions.append(p.run,p.stop,reset,open,close);
 p.stack=liveElement('div','doc-live-stack');p.highlight=liveElement('pre','doc-live-highlight');p.highlight.setAttribute('aria-hidden','true');p.editor=liveElement('textarea','doc-live-editor');p.editor.spellcheck=false;p.editor.autocapitalize='off';p.editor.setAttribute('autocorrect','off');p.editor.setAttribute('wrap','off');p.editor.maxLength=32768;p.editor.value=original;p.editor.setAttribute('aria-label',translate('Editable TRIPOL example','Ejemplo TRIPOL editable'));p.stack.append(p.highlight,p.editor);
 const inputBox=liveElement('details','doc-live-input');inputBox.append(liveLabel(liveElement('summary','doc-live-label'),'Optional input','Entrada opcional'));p.input=liveElement('textarea','doc-live-stdin');p.input.rows=2;p.input.maxLength=32768;inputBox.append(p.input);inputBox.open=/\binput(?:_int)?\s*\(/.test(original);
 const consoleBox=liveElement('div','doc-live-console');p.status=liveElement('div','doc-live-status');p.status.setAttribute('role','status');p.output=liveElement('pre','doc-live-output');p.output.setAttribute('aria-live','polite');p.files=liveElement('div','doc-live-files');p.files.hidden=true;consoleBox.append(p.status,p.output,p.files);
 const help=liveLabel(liveElement('p','doc-live-help'),'Runs with the real Web Engine R3. Unsupported browser features report a compiler error.','Ejecuta con Web Engine R3 real. Las funciones no compatibles muestran el error del compilador.');
 panel.append(toolbar,p.stack,inputBox,consoleBox,help);source.after(panel);livePanels.push(p);
 p.editor.addEventListener('input',()=>{livePaint(p);if(/\binput(?:_int)?\s*\(/.test(p.editor.value))inputBox.open=true});
 p.editor.addEventListener('scroll',()=>liveScroll(p));
 p.editor.addEventListener('keydown',e=>{if((e.ctrlKey||e.metaKey)&&e.key==='Enter'){e.preventDefault();liveRun(p)}else if(e.key==='Tab'&&!e.shiftKey){e.preventDefault();p.editor.setRangeText('    ',p.editor.selectionStart,p.editor.selectionEnd,'end');livePaint(p)}else if(e.key==='Escape'&&liveJob?.panel===p)liveStop(p)});
 liveStatus(p,'Ready · Ctrl + Enter','Listo · Ctrl + Enter');liveRefresh();
 return p
}
function installLive(source,host,original){
 if(!original||original.length>32768)return;
 let p;const button=liveButton('Try it','Probar',()=>{if(!p)p=makeLivePanel(source,original,button);if(!p.panel.hidden){liveStop(p);p.panel.hidden=true;source.hidden=false;button.setAttribute('aria-expanded','false');return}source.hidden=true;p.panel.hidden=false;button.setAttribute('aria-expanded','true');livePaint(p);p.editor.focus()},'doc-copy doc-live-action');
 button.setAttribute('aria-expanded','false');host.append(button)
}
function installAllLive(){
 for(const pre of $$('.docs-article pre')){
  // Function signatures and library implementation source are reference-only.
  if(pre.classList.contains('signature')||pre.closest('.function-card'))continue;
  if(codeMode(pre.querySelector('code')||pre)!=='tripol')continue;
  const head=pre.closest('.code-shell')?.querySelector('.code-head');if(!head)continue;
  let actions=head.querySelector('.doc-code-actions');if(!actions){actions=liveElement('div','doc-code-actions');const copy=head.querySelector('.doc-copy');if(copy)actions.append(copy);head.append(actions)}
  installLive(pre,actions,pre.textContent);
 }
}
window.addEventListener('pagehide',()=>{if(liveJob)liveStop(liveJob.panel,false)});

function apply(){document.documentElement.lang=lang;labels.forEach(({el,en,es})=>el.textContent=translate(en,es));placeholders.forEach(({el,en,es})=>el.placeholder=translate(en,es));$('#doc-lang').textContent=lang==='en'?'ES':'EN';$('#doc-lang').setAttribute('aria-label',translate('Switch to Spanish','Cambiar a inglés'));try{localStorage.setItem('tripol-language',lang)}catch{};try{const u=new URL(location.href);u.searchParams.set('lang',lang);history.replaceState(null,'',u.href)}catch{};localLinks();highlightDocs();liveRefresh();buildToc();if($('#doc-search').open)renderSearch();}
let timer;function toast(msg){const t=$('#doc-toast');t.textContent=msg;t.classList.add('visible');clearTimeout(timer);timer=setTimeout(()=>t.classList.remove('visible'),1800)}
$$('.docs-article pre').forEach(pre=>{if(pre.classList.contains('signature'))return;const shell=document.createElement('div');shell.className='code-shell';const head=document.createElement('div');head.className='code-head';const title=document.createElement('span');const mode=codeMode(pre.querySelector('code')||pre);title.textContent=mode==='tripol'?'TRIPOL':['cmd','sh','bash','shell','powershell','ps1'].includes(mode)?(mode==='powershell'?'POWERSHELL':'TERMINAL'):mode==='text'?'CODE':mode.toUpperCase();const copy=document.createElement('button');copy.className='doc-copy';copy.dataset.es='Copiar';copy.textContent='Copy';copy.addEventListener('click',async()=>{try{await navigator.clipboard.writeText(pre.textContent);toast(translate('Copied to clipboard','Copiado al portapapeles'))}catch{const range=document.createRange();range.selectNodeContents(pre);const selection=window.getSelection();selection.removeAllRanges();selection.addRange(range);toast(translate('Selected — copy manually','Seleccionado — copia manualmente'))}});head.append(title,copy);pre.before(shell);shell.append(head,pre)});
function buildToc(){const target=$('#toc-links');target.replaceChildren();$$('.docs-article h2,.docs-article .function-card h3').forEach((h,i)=>{const id=h.id||h.closest('.function-card')?.id||('section-'+i);if(!h.id&&!h.closest('.function-card'))h.id=id;const a=document.createElement('a');a.href='#'+id;a.textContent=h.textContent.replace(/\s*#\s*$/,'');if(h.tagName==='H3')a.className='sub';target.append(a)});}
installAllLive();registerLabels();apply();$('#doc-lang').addEventListener('click',()=>{lang=lang==='en'?'es':'en';apply()});
$('#doc-menu').addEventListener('click',()=>{$('#doc-menu').setAttribute('aria-expanded',String(side.classList.toggle('open')))});
$('#side-filter').addEventListener('input',e=>{const q=e.target.value.trim().toLowerCase();$$('.side-libs a').forEach(a=>a.hidden=!a.textContent.toLowerCase().includes(q))});
const lf=$('#library-filter');if(lf)lf.addEventListener('input',()=>{let n=0;const q=lf.value.trim().toLowerCase();$$('.library-card').forEach(a=>{a.hidden=!a.dataset.filter.toLowerCase().includes(q);if(!a.hidden)n++});$('#library-empty').hidden=!!n});
const dialog=$('#doc-search'),input=$('#global-search');
function renderSearch(){const q=input.value.trim().toLowerCase();const list=DOC_SEARCH.filter(item=>[item.title,item.detail,item.es||''].join(' ').toLowerCase().includes(q)).slice(0,30);const results=$('#search-results');results.replaceChildren();if(!list.length){const p=document.createElement('p');p.textContent=translate('No results. Try a library or function name.','Sin resultados. Prueba una biblioteca o función.');results.append(p);return}list.forEach(item=>{const a=document.createElement('a');a.className='search-item';const u=new URL(root+item.url,location.href);u.searchParams.set('lang',lang);a.href=u.href;const title=document.createElement('strong');title.textContent=item.title;const detail=document.createElement('small');detail.textContent=lang==='es'&&item.es?item.es:item.detail;a.append(title,detail);results.append(a)})}
function openSearch(){if(!dialog.open)dialog.showModal();renderSearch();input.focus()}
$('#search-open').addEventListener('click',openSearch);$('#search-close').addEventListener('click',()=>dialog.close());input.addEventListener('input',renderSearch);dialog.addEventListener('click',e=>{if(e.target===dialog)dialog.close()});input.addEventListener('keydown',e=>{if(e.key==='ArrowDown'){e.preventDefault();$('#search-results a')?.focus()}else if(e.key==='Enter')$('#search-results a')?.click()});document.addEventListener('keydown',e=>{if((e.ctrlKey||e.metaKey)&&e.key.toLowerCase()==='k'){e.preventDefault();openSearch()}if(e.key==='Escape'&&side.classList.contains('open')){side.classList.remove('open');$('#doc-menu').setAttribute('aria-expanded','false')}});
const active=side.querySelector('[aria-current=page]');if(active&&active.offsetTop>side.clientHeight-90)side.scrollTop=active.offsetTop-100;
})();
