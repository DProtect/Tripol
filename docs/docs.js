'use strict';
(()=>{
const $=s=>document.querySelector(s), $$=s=>[...document.querySelectorAll(s)];
let lang='en';try{if(localStorage.getItem('tripol-language')==='es')lang='es'}catch{}
const chosen=new URLSearchParams(location.search).get('lang');if(chosen==='en'||chosen==='es')lang=chosen;
const root=document.body.dataset.docRoot||'./';
// A freshly loaded script must also load its matching stylesheet.
for(const link of document.querySelectorAll('link[rel="stylesheet"]')){const u=new URL(link.href,location.href);if(u.pathname.endsWith('/docs/docs.css')&&u.searchParams.get('v')!=='library-categories-r4'){u.searchParams.set('v','library-categories-r4');link.href=u.href}}
function localLinks(){for(const a of $$('a[href]')){const raw=a.getAttribute('href');if(!raw||raw.startsWith('#')||/^(https?:|mailto:)/i.test(raw)&&a.origin!==location.origin)continue;const u=new URL(raw,location.href);if(!u.pathname.endsWith('.html'))continue;u.searchParams.set('lang',lang);a.href=u.href;}}
// Shared navigation and current release label on all reference pages.
const side=document.querySelector('#doc-side');
side.querySelector('.side-kicker').textContent='DISCORD REVOLUTION / 0.28.0';
const crumb=document.querySelector('.crumb > span:last-child');if(crumb&&/^0\.(24|25)$/.test(crumb.textContent))crumb.textContent='0.26';
const meta=document.querySelector('.lib-meta > span:first-child');if(meta&&/^(ASCENSION 0\.24|EVOLUTION 0\.25)$/.test(meta.textContent))meta.textContent='OVERDRIVE 0.26';
const guide=side.querySelector('nav[aria-label="Guides"]');
for(const [file,label] of [['genesis','GENESIS R3'],['update','UPDATE R2']]){let a=guide.querySelector('a[href$="'+file+'.html"]');if(!a){a=document.createElement('a');a.href=root+file+'.html';guide.append(a);}a.textContent=label;}
const libs=side.querySelector('.side-libs');
for(const [name,before,count] of [['queue','random',7],['stack','stats',7],['bitset','buffer',6],['profiler','queue',2]]){if(!libs.querySelector('a[href$="/'+name+'.html"]')){const a=document.createElement('a');a.href=root+'libraries/'+name+'.html';const label=document.createElement('span');label.textContent=name;const n=document.createElement('small');n.textContent=String(count);a.append(label,n);libs.insertBefore(a,libs.querySelector('a[href$="/'+before+'.html"]'));}}
const LIBRARY_LEVELS=[{"en":"Beginner","es":"Principiante","names":["console","math","text","convert","random","datetime","arrays","collections","search","sort","validate","errors","files","path","cli","csv"]},{"en":"Intermediate","es":"Intermedio","names":["dictionary","buffer","stack","queue","sequences","stats","testing","time","json","colors","geometry","animation","vector","matrix","physics","canvas","sprite","camera","particles","image","heap","io","regex","stream","utf8"]},{"en":"Advanced","es":"Avanzado","names":["binary","bits","bitset","encoding","checksum","compression","simd","memory","system","profiler","env","ui","game","events","network","http","https","discord","atomic","audio","crypto","process","thread"]}];
const LIBRARY_COUNTS={"arrays":8,"cli":5,"collections":6,"console":3,"convert":8,"csv":3,"datetime":7,"errors":7,"files":5,"math":12,"path":6,"random":3,"search":4,"sort":4,"text":9,"validate":6,"animation":6,"buffer":8,"camera":4,"canvas":10,"colors":8,"dictionary":8,"geometry":6,"heap":6,"image":6,"io":2,"json":18,"matrix":6,"particles":6,"physics":5,"queue":7,"regex":2,"sequences":4,"sprite":6,"stack":7,"stats":6,"stream":6,"testing":4,"time":6,"utf8":5,"vector":7,"atomic":4,"audio":6,"binary":8,"bits":6,"bitset":6,"checksum":3,"compression":3,"crypto":3,"discord":98,"encoding":4,"env":1,"events":14,"game":7,"http":16,"https":3,"memory":4,"network":3,"process":3,"profiler":2,"simd":1,"system":5,"thread":3,"ui":7};

function libraryLearningOrder(){
 const groups=LIBRARY_LEVELS;
 const nameOf=a=>a.getAttribute('href')?.match(/\/?([^/]+)\.html(?:[?#].*)?$/)?.[1];
 function reorder(container,selector,kind){
  if(!container)return;
  const links=new Map([...container.querySelectorAll(selector)].map(a=>[nameOf(a),a]));
  const fragment=document.createDocumentFragment();
  groups.forEach((g,i)=>{
   const header=document.createElement(kind==='side'?'div':'h3');
   header.className=kind==='side'?'side-lib-level':'catalog-lib-level';
   header.dataset.es=g.es;header.textContent=g.en;header.dataset.level=String(i);fragment.append(header);
   for(const name of g.names){
    let a=links.get(name);
    if(!a&&kind==='side'){a=document.createElement('a');a.href=root+'libraries/'+name+'.html';const label=document.createElement('span'),n=document.createElement('small');label.textContent=name;a.append(label,n)}
    if(!a)continue;
    links.delete(name);a.dataset.level=String(i);
    const n=a.querySelector('small');if(n){if(kind==='side')n.textContent=String(LIBRARY_COUNTS[name]);else if(n.firstChild?.nodeType===3)n.firstChild.nodeValue=String(LIBRARY_COUNTS[name])+' '}
    if(kind==='catalog'&&typeof a.dataset.filter==='string'){a.dataset.filter+=' '+g.en+' '+g.es+' '+(typeof DOC_SEARCH==='undefined'?[]:DOC_SEARCH).filter(x=>x.url.startsWith('libraries/'+name+'.html#')).map(x=>x.title).join(' ')}
    fragment.append(a)
   }
  });
  // Preserve future modules even before they receive a difficulty classification.
  for(const a of links.values())fragment.append(a);
  container.replaceChildren(fragment)
 }
 reorder(document.querySelector('.side-libs'),'a','side');
 reorder(document.querySelector('#library-grid'),'.library-card','catalog');
}
function refreshLevelHeadings(container){
 if(!container)return;
 for(const h of container.querySelectorAll('.side-lib-level,.catalog-lib-level'))h.hidden=![...container.querySelectorAll('[data-level="'+h.dataset.level+'"]:is(a,.side-library-group)')].some(a=>!a.hidden)
}

libraryLearningOrder();

/* Library category navigation: register additional modules here as their pages are grouped. */
const LIBRARY_CATEGORY_NAV={http:[{"id":"http-guide","en":"HTTP guide","es":"Guía HTTP","count":0},{"id":"http-start","en":"Getting started","es":"Primeros pasos","count":1},{"id":"http-requests","en":"Requests","es":"Peticiones","count":4},{"id":"http-responses","en":"Responses","es":"Respuestas","count":6},{"id":"http-methods","en":"HTTP methods","es":"Métodos HTTP","count":15},{"id":"http-auth","en":"Headers & authentication","es":"Encabezados y autenticación","count":7},{"id":"http-json","en":"JSON & forms","es":"JSON y formularios","count":2},{"id":"http-clients","en":"Clients & sessions","es":"Clientes y sesiones","count":3},{"id":"http-cookies","en":"Cookies","es":"Cookies","count":1},{"id":"http-files","en":"Files & downloads","es":"Archivos y descargas","count":2},{"id":"http-streaming","en":"Streaming","es":"Transmisión por flujo","count":2},{"id":"http-concurrent","en":"Concurrent requests","es":"Peticiones concurrentes","count":5},{"id":"http-security","en":"HTTPS & security","es":"HTTPS y seguridad","count":4},{"id":"http-errors","en":"Errors & diagnostics","es":"Errores y diagnóstico","count":4},{"id":"http-discord","en":"Discord integration","es":"Integración con Discord","count":0}],discord:[{"id":"api-start","en":"Getting started & connection","es":"Primeros pasos y conexión","count":19},{"id":"api-commands","en":"Commands & context menus","es":"Comandos y menús contextuales","count":55},{"id":"api-events","en":"Events & messages","es":"Eventos y mensajes","count":22},{"id":"api-embeds","en":"Embeds","es":"Embeds","count":23},{"id":"api-components","en":"Components V2","es":"Components V2","count":16},{"id":"api-buttons","en":"Buttons & emojis","es":"Botones y emojis","count":11},{"id":"api-selects","en":"Select menus","es":"Selectores","count":11},{"id":"api-forms","en":"Forms & modals","es":"Formularios","count":21},{"id":"api-files","en":"Files & attachments","es":"Archivos y adjuntos","count":10},{"id":"api-moderation","en":"Moderation & permissions","es":"Moderación y permisos","count":25},{"id":"api-economy","en":"Cooldowns & economy","es":"Cooldowns y economía","count":7},{"id":"api-servers","en":"Servers, channels & roles","es":"Servidores, canales y roles","count":27},{"id":"api-webhooks","en":"Webhooks & invitations","es":"Webhooks e invitaciones","count":14},{"id":"api-threads","en":"Threads & forums","es":"Hilos y foros","count":12},{"id":"api-helpers","en":"Low-level helpers & legacy API","es":"Utilidades de bajo nivel y API clásica","count":21},{"id":"api-http","en":"Deferred HTTP integration","es":"Integración HTTP diferida","count":4}]};
function installLibraryCategories(){
 for(const [name,sections] of Object.entries(LIBRARY_CATEGORY_NAV)){
  const original=[...libs.querySelectorAll(':scope > a')].find(a=>new URL(a.href,location.href).pathname.endsWith('/libraries/'+name+'.html'));
  if(!original)continue;
  const url=new URL(original.href,location.href),group=document.createElement('details');
  group.className='side-library-group';group.dataset.library=name;group.dataset.level=original.dataset.level||'';
  const summary=document.createElement('summary'),label=document.createElement('span'),count=document.createElement('small');
  label.textContent=name;count.textContent=String(LIBRARY_COUNTS[name]);summary.append(label,count);
  group.append(summary);
  const list=document.createElement('div');list.className='side-library-categories';
  function add(en,es,hash='',href=url.href){
   const a=document.createElement('a'),title=document.createElement('span');
   title.textContent=en;title.dataset.es=es;a.append(title);a.href=href+(hash?'#'+hash:'');
   if(location.pathname===url.pathname&&hash){a.addEventListener('click',()=>{side.classList.remove('open');$('#doc-menu').setAttribute('aria-expanded','false')})}
   list.append(a);
  }
  add('Overview / all functions','Introducción / todas las funciones');
  if(name==='discord')add('Bot guide','Guía de bots','',new URL(root+'discord.html',location.href).href);
  for(const category of sections)add(category.en,category.es,category.id);
  group.append(list);group.open=[...list.querySelectorAll('a')].some(a=>new URL(a.href,location.href).pathname===location.pathname);
  if(group.open)summary.classList.add('current-library');
  original.replaceWith(group);
 }
}
installLibraryCategories();

/* Stable anchor navigation and scroll tracking for categorized libraries. */
function installLibraryCategoryTracking(){
 const group=[...libs.querySelectorAll('.side-library-group')].find(g=>[...g.querySelectorAll('.side-library-categories a')].some(a=>new URL(a.href,location.href).pathname===location.pathname));
 if(!group)return;
 const links=[...group.querySelectorAll('.side-library-categories a')].filter(a=>new URL(a.href,location.href).pathname===location.pathname);
 const overview=links.find(a=>!new URL(a.href,location.href).hash);
 const sections=links.map(a=>({a,id:new URL(a.href,location.href).hash.slice(1)})).filter(x=>x.id).map(x=>({...x,heading:document.getElementById(x.id)})).filter(x=>x.heading);
 if(!overview||!sections.length)return;
 document.documentElement.classList.add('categorized-doc');
 let selected=null,scheduled=false,lockedY=null;
 const header=document.querySelector('.docs-header');
 const anchorOffset=()=>Math.ceil((header?.getBoundingClientRect().bottom||0)+24);
 function mark(link){
  if(!link||selected===link)return;
  for(const a of links){a.removeAttribute('aria-current');a.classList.remove('category-active')}
  link.setAttribute('aria-current',link===overview?'page':'location');
  link.classList.add('category-active');selected=link;
 }
 function resolve(hash){
  if(!hash)return null;
  try{return document.getElementById(decodeURIComponent(hash.slice(1)))}catch{return null}
 }
 function linkFor(target){
  const category=target?.closest('.library-api-category')?.querySelector('h2[id]')?.id||target?.id;
  return sections.find(x=>x.id===category)?.a||overview;
 }
 function jump(target){
  mark(linkFor(target));
  const y=target?Math.max(0,window.scrollY+target.getBoundingClientRect().top-anchorOffset()):0;
  window.scrollTo({top:y,left:0,behavior:'instant'});
  // Capture the clamped position after the synchronous jump. No scrollspy
  // updates may replace the clicked category until the viewport really moves.
  lockedY=window.scrollY;
 }
 function track(){
  scheduled=false;
  if(lockedY!==null){
   if(Math.abs(window.scrollY-lockedY)<=2)return;
   lockedY=null;
  }
  const offset=anchorOffset()+3;
  let current=overview;
  for(const x of sections){if(x.heading.getBoundingClientRect().top<=offset)current=x.a;else break}
  if(window.scrollY>0&&window.scrollY+window.innerHeight>=document.documentElement.scrollHeight-4)current=sections.at(-1).a;
  mark(current);
 }
 function schedule(){if(!scheduled){scheduled=true;window.requestAnimationFrame(track)}}
 document.addEventListener('click',event=>{
  if(event.defaultPrevented||event.button!==0||event.ctrlKey||event.metaKey||event.shiftKey||event.altKey)return;
  const a=event.target.closest?.('a[href]');
  if(!a||a.hasAttribute('download')||(a.target&&a.target!=='_self'))return;
  const url=new URL(a.href,location.href);
  if(url.origin!==location.origin||url.pathname!==location.pathname)return;
  const target=resolve(url.hash);
  if(!target&&a!==overview)return;
  if(target&&!target.closest('.docs-article'))return;
  event.preventDefault();
  if(location.hash!==url.hash){history.pushState(null,'',location.pathname+location.search+url.hash)}
  if(side.classList.contains('open')){side.classList.remove('open');$('#doc-menu').setAttribute('aria-expanded','false')}
  jump(target);
  // Keyboard activation moves focus to the destination without scrolling again.
  if(event.detail===0&&target){if(!target.hasAttribute('tabindex'))target.setAttribute('tabindex','-1');target.focus({preventScroll:true})}
 });
 window.addEventListener('hashchange',()=>{const target=resolve(location.hash);if(target)jump(target);else{lockedY=null;schedule()}});
 window.addEventListener('scroll',schedule,{passive:true});
 window.addEventListener('resize',()=>{lockedY=null;schedule()},{passive:true});
 window.addEventListener('pageshow',event=>{if(event.persisted){lockedY=null;schedule()}});
 const initial=resolve(location.hash);
 if(initial)jump(initial);else{mark(overview);schedule()}
}

function filterLibraryNavigation(q){
 const query=normalizeSearch(q);
 for(const el of libs.querySelectorAll(':scope > a,:scope > .side-library-group')){
  const hit=normalizeSearch(el.textContent+' '+(el.dataset.library||'')).includes(query)||[...el.querySelectorAll('[data-es]')].some(x=>normalizeSearch(x.dataset.es).includes(query));
  el.hidden=!hit;
  if(query&&hit&&el.tagName==='DETAILS')el.open=true;
 }
 refreshLevelHeadings(libs);
}

const LEARN_LESSONS=[{"file":"learn-hello.html","en":"Hello world","es":"Hola mundo","num":1},{"file":"learn-variables.html","en":"Variables","es":"Variables","num":2},{"file":"learn-types.html","en":"Data types","es":"Tipos de datos","num":3},{"file":"learn-operators.html","en":"Operators","es":"Operadores","num":4},{"file":"learn-conditions.html","en":"Conditionals","es":"Condicionales","num":5},{"file":"learn-loops.html","en":"Loops","es":"Bucles","num":6},{"file":"learn-functions.html","en":"Functions","es":"Funciones","num":7},{"file":"learn-arrays.html","en":"Arrays","es":"Arreglos","num":8},{"file":"learn-lists.html","en":"Lists","es":"Listas","num":9},{"file":"learn-strings.html","en":"Strings","es":"Strings","num":10},{"file":"learn-structs.html","en":"Structures","es":"Estructuras","num":11},{"file":"learn-exceptions.html","en":"Error handling","es":"Manejo de errores","num":12},{"file":"learn-memory.html","en":"Memory","es":"Memoria","num":13},{"file":"learn-modules.html","en":"Modules","es":"Módulos","num":14}];
function installLearningNavigation(){
 if(side.querySelector('.learn-nav'))return;
 const heading=document.createElement('h2');heading.className='learn-nav-title';const label=document.createElement('span');label.dataset.es='Aprende TRIPOL';label.textContent='Learn TRIPOL';heading.append(label);
 const nav=document.createElement('nav');nav.className='learn-nav';nav.setAttribute('aria-label','Learning guide');
 const add=(file,en,es,num)=>{const a=document.createElement('a');a.href=root+file;const label=document.createElement('span');label.textContent=en;label.dataset.es=es;a.append(label);if(num){const small=document.createElement('small');small.textContent=String(num).padStart(2,'0');a.append(small)}nav.append(a)};
 add('learn.html','Learning path','Ruta de aprendizaje');for(const l of LEARN_LESSONS)add(l.file,l.en,l.es,l.num);
 guide.after(heading,nav);
}
installLearningNavigation();

for(const a of side.querySelectorAll('nav a')){a.removeAttribute('aria-current');if(new URL(a.href,location.href).pathname===location.pathname&&!new URL(a.href,location.href).hash)a.setAttribute('aria-current','page');}
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
function liveRun(p){liveFiles(p);p.preview?.reset();const previewToken=p.previewToken=(p.previewToken||0)+1;
 if(liveJob)liveStop(liveJob.panel);
 const source=p.editor.value;
 if(new TextEncoder().encode(source).length>32768||new TextEncoder().encode(p.input.value).length>32768){liveStatus(p,'Input exceeds 32 KiB','La entrada supera 32 KiB','error');return}
 p.output.textContent='';p.run.disabled=true;p.stop.hidden=false;liveStatus(p,'Compiling with Web Engine R3…','Compilando con Web Engine R3…');
 let worker;try{worker=new Worker(new URL(root+'../playground/web-worker-r3.2.js',location.href))}catch(e){p.run.disabled=false;p.stop.hidden=true;liveStatus(p,'Compiler worker unavailable','Worker del compilador no disponible','error');p.output.textContent=String(e.message||e);return}
 const job={panel:p,worker,timer:null};liveJob=job;
 worker.onmessage=({data})=>{if(liveJob!==job)return;liveStop(p,false);p.output.textContent=(data.output||'')+(data.ok?'':'\n'+translate('Error: ','Error: ')+String(data.error||'Unknown error'));if(Array.isArray(data.frames)&&data.frames.length){ensureGraphics(p).then(view=>{if(p.previewToken===previewToken&&!p.panel.hidden)view.show(data.frames)}).catch(e=>{if(p.previewToken===previewToken){p.graphics.hidden=false;p.graphics.textContent=String(e.message||e)}})}if(data.ok){if(Array.isArray(data.files))liveFiles(p,data.files);const time=Number(data.compileMs||0).toFixed(1)+' ms'+(data.build?' · '+data.build.slice(0,7):'');liveStatus(p,'WASM R3 · compiled in '+time,'WASM R3 · compilado en '+time,'ok')}else liveStatus(p,'Compiler/runtime diagnostic · no fallback','Diagnóstico del compilador/runtime · sin sustitución','error')};
 worker.onerror=()=>{if(liveJob!==job)return;liveStop(p,false);liveStatus(p,'Web Engine unavailable','Web Engine no disponible','error');p.output.textContent=translate('The compiler worker could not start. Open in Playground or try again.','No se pudo iniciar el compilador. Abre en Playground o vuelve a intentar.')};
 job.timer=setTimeout(()=>{if(liveJob!==job)return;liveStop(p,false);liveStatus(p,'Execution time limit (15 s)','Límite de ejecución (15 s)','error')},15000);
 try{worker.postMessage({source,input:p.input.value,compileOnly:false})}catch(e){liveStop(p,false);liveStatus(p,'Could not send source','No se pudo enviar el código','error');p.output.textContent=String(e.message||e)}
}
let graphicsLoad;
function ensureGraphics(p){if(p.preview)return Promise.resolve(p.preview);if(!graphicsLoad)graphicsLoad=import(new URL(root+'../playground/graphics-preview-r1.js',location.href).href).catch(e=>{graphicsLoad=null;throw e});return graphicsLoad.then(()=>p.preview||(p.preview=globalThis.TripolGraphics.create(p.graphics,()=>lang)))}
function resetGraphics(p){p.previewToken=(p.previewToken||0)+1;p.preview?.reset();if(p.graphics)p.graphics.hidden=true}
function makeLivePanel(source,original,trigger){
 const p={original,source};const panel=liveElement('section','doc-live-panel');p.panel=panel;panel.hidden=true;
 const toolbar=liveElement('div','doc-live-toolbar');const filename=liveElement('span','doc-live-file','example.tri');const actions=liveElement('div','doc-live-tools');toolbar.append(filename,actions);
 p.run=liveButton('▶ Run','▶ Ejecutar',()=>liveRun(p),'doc-live-run');
 p.stop=liveButton('Stop','Detener',()=>liveStop(p));p.stop.hidden=true;
 const reset=liveButton('Reset','Restaurar',()=>{liveStop(p,false);p.editor.value=p.original;p.input.value=p.source.dataset.input||'';resetGraphics(p);liveFiles(p);p.output.textContent='';livePaint(p);liveStatus(p,'Ready','Listo')});
 const open=liveButton('Playground ↗','Playground ↗',()=>livePlayground(p.editor.value));
 const close=liveButton('Close','Cerrar',()=>{liveStop(p);resetGraphics(p);panel.hidden=true;source.hidden=false;trigger.setAttribute('aria-expanded','false');trigger.focus()});
 actions.append(p.run,p.stop,reset,open,close);
 p.stack=liveElement('div','doc-live-stack');p.highlight=liveElement('pre','doc-live-highlight');p.highlight.setAttribute('aria-hidden','true');p.editor=liveElement('textarea','doc-live-editor');p.editor.spellcheck=false;p.editor.autocapitalize='off';p.editor.setAttribute('autocorrect','off');p.editor.setAttribute('wrap','off');p.editor.maxLength=32768;p.editor.value=original;p.editor.setAttribute('aria-label',translate('Editable TRIPOL example','Ejemplo TRIPOL editable'));p.stack.append(p.highlight,p.editor);
 const inputBox=liveElement('details','doc-live-input');inputBox.append(liveLabel(liveElement('summary','doc-live-label'),'Optional input','Entrada opcional'));p.input=liveElement('textarea','doc-live-stdin');p.input.rows=2;p.input.maxLength=32768;inputBox.append(p.input);inputBox.open=/\binput(?:_int)?\s*\(/.test(original);
 const consoleBox=liveElement('div','doc-live-console');p.status=liveElement('div','doc-live-status');p.status.setAttribute('role','status');p.output=liveElement('pre','doc-live-output');p.output.setAttribute('aria-live','polite');p.files=liveElement('div','doc-live-files');p.files.hidden=true;p.graphics=liveElement('section','tripol-preview');p.graphics.hidden=true;consoleBox.append(p.status,p.output,p.graphics,p.files);
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
 let p;const button=liveButton('Try it','Probar',()=>{if(!p){p=makeLivePanel(source,original,button);if(source.dataset.input){p.input.value=source.dataset.input;p.input.parentElement.open=true;}}if(!p.panel.hidden){liveStop(p);resetGraphics(p);p.panel.hidden=true;source.hidden=false;button.setAttribute('aria-expanded','false');return}source.hidden=true;p.panel.hidden=false;button.setAttribute('aria-expanded','true');livePaint(p);p.editor.focus()},'doc-copy doc-live-action');
 button.setAttribute('aria-expanded','false');host.append(button)
}
function installAllLive(){
 for(const pre of $$('.docs-article pre')){
  // Function signatures and library implementation source are reference-only.
  if(pre.classList.contains('signature')||pre.hasAttribute('data-native-only'))continue;
  const lib=location.pathname.match(/\/docs\/libraries\/([^/]+)\.html$/)?.[1];const card=pre.closest('.function-card');if(card&&!pre.classList.contains('doc-usage')&&(!['game','canvas','sprite','image','particles','animation','camera','physics','colors'].includes(lib)||lib==='game'&&card.id==='key'||lib==='image'&&card.id==='save_ppm'))continue;
  if(!card&&document.querySelector('.doc-lib-guide[data-web-status="compiler-rejected"]'))continue;
  if(codeMode(pre.querySelector('code')||pre)!=='tripol')continue;
  const head=pre.closest('.code-shell')?.querySelector('.code-head');if(!head)continue;
  let actions=head.querySelector('.doc-code-actions');if(!actions){actions=liveElement('div','doc-code-actions');const copy=head.querySelector('.doc-copy');if(copy)actions.append(copy);head.append(actions)}
  installLive(pre,actions,pre.textContent);
 }
}
window.addEventListener('pagehide',()=>{if(liveJob)liveStop(liveJob.panel,false)});

function apply(){document.documentElement.lang=lang;labels.forEach(({el,en,es})=>el.textContent=translate(en,es));placeholders.forEach(({el,en,es})=>el.placeholder=translate(en,es));$('#doc-lang').textContent=lang==='en'?'ES':'EN';$('#doc-lang').setAttribute('aria-label',translate('Switch to Spanish','Cambiar a inglés'));try{localStorage.setItem('tripol-language',lang)}catch{};try{const u=new URL(location.href);u.searchParams.set('lang',lang);history.replaceState(null,'',u.href)}catch{};localLinks();highlightDocs();liveRefresh();buildToc();if($('#doc-search')?.open)renderSearch();}
let timer;function toast(msg){const t=$('#doc-toast');t.textContent=msg;t.classList.add('visible');clearTimeout(timer);timer=setTimeout(()=>t.classList.remove('visible'),1800)}
const LEARNING_GUIDES = Object.fromEntries((globalThis.TripolLibraries || []).flatMap(lib =>
 [...lib.functions, ...(lib.methods || [])].map(f => [lib.name + '.' + (f.receiver ? f.receiver + '.' : '') + f.name, f])));
const LEARNING_STATUS={"animation":"executed","arrays":"executed","binary":"executed","bits":"executed","bitset":"executed","buffer":"executed","camera":"executed","canvas":"executed","checksum":"executed","collections":"executed","colors":"executed","compression":"executed","console":"executed","convert":"executed","datetime":"executed","dictionary":"executed","discord":"executed","encoding":"executed","env":"compiler-rejected","errors":"executed","events":"compiler-rejected","files":"executed","game":"graphics-preview","geometry":"executed","http":"compiled-only","https":"compiler-rejected","image":"executed","json":"executed","math":"executed","matrix":"executed","memory":"executed","network":"compiled-only","particles":"executed","path":"executed","physics":"executed","profiler":"compiled-only","queue":"executed","random":"executed","search":"executed","sequences":"executed","simd":"executed","sort":"executed","sprite":"executed","stack":"executed","stats":"executed","system":"executed","testing":"executed","text":"executed","time":"compiled-only","ui":"compiler-rejected","validate":"executed","vector":"executed"};

function learningReference(){
 const match=location.pathname.match(/\/docs\/libraries\/([^/]+)\.html$/);if(!match)return;
 const lib=match[1],article=$('.docs-article');if(!article)return;
 const data=LEARNING_GUIDES,module=globalThis.TripolLibraries?.find(item=>item.name===lib);
 const status=LEARNING_STATUS[lib]||(module?.functions.some(f=>!f.native)?'executed':'unknown');
 const make=(tag,cls,text)=>{const el=document.createElement(tag);if(cls)el.className=cls;if(text!==undefined)el.textContent=text;return el};
 const bilingual=(en,es)=>{const el=make('span','',en);el.dataset.es=es;return el};
 const title=(en,es)=>{const el=make('strong');el.append(bilingual(en,es));return el};
 const types={int:['64-bit integer','Entero de 64 bits'],bool:['Boolean (true or false)','Booleano (true o false)'],string:['Text (UTF-8)','Texto (UTF-8)'],float:['Decimal number','Número decimal'],ints:['Integer array','Arreglo de enteros'],list:['List','Lista'],void:['No value','Sin valor'],string_list:['Text list','Lista de textos'],bool_list:['Boolean list','Lista de booleanos'],float_list:['Decimal list','Lista de decimales']};
 function splitParams(s){let out=[],part='',depth=0,quote='',esc=false;for(const c of s){if(quote){part+=c;if(esc)esc=false;else if(c==='\\')esc=true;else if(c===quote)quote='';continue}if(c==='"'||c==="'"){quote=c;part+=c;continue}if('([{'.includes(c))depth++;if(')]}'.includes(c))depth--;if(c===','&&depth===0){out.push(part.trim());part=''}else part+=c}if(part.trim())out.push(part.trim());return out}
 const guide=make('div','doc-lib-guide');guide.dataset.webStatus=status;
 const labels={'graphics-preview':['Web graphics preview: framebuffer and triangle playback. Live keyboard, mouse and controller polling requires native TRIPOL.','Vista gráfica web: reproducción de píxeles y triángulos. El teclado, ratón y mando en tiempo real requieren TRIPOL nativo.'],executed:['The main library example ran in Web Engine R3. Individual functions may differ.','El ejemplo principal se ejecutó en Web Engine R3. Algunas funciones pueden diferir.'],'compiler-rejected':['The main example requires native TRIPOL; browser execution is unavailable.','El ejemplo principal requiere TRIPOL nativo; no se puede ejecutar en el navegador.'],'compiled-only':['The main example compiled for WASM; execution was not verified.','El ejemplo principal compiló para WASM; la ejecución no está verificada.'],unknown:['Web compatibility has not been verified.','La compatibilidad web no está verificada.']};
 const note=make('p');note.append(bilingual(...labels[status]||labels.unknown));
 const source=make('a');source.href='https://github.com/Tripollas/Tripol/blob/main/stdlib/'+lib+'.tri';source.target='_blank';source.rel='noopener noreferrer';source.append(bilingual('View library source on GitHub ↗','Ver código fuente de la biblioteca en GitHub ↗'));guide.append(note,source);article.querySelector('.lib-meta')?.after(guide);
 for(const card of article.querySelectorAll('.function-card')){
  const signature=card.querySelector('pre.signature');if(!signature)continue;
  const parsed=signature.textContent.trim().match(/^func\s+[\w.]+\s*\(([\s\S]*)\)\s*(?:->\s*([\w]+))?\s*$/);
  for(const details of card.querySelectorAll('details')){if(details.querySelector('pre'))details.remove()}
  const row=data[lib+'.'+card.id];if(row){const p=card.querySelector('p');if(p)p.replaceChildren(bilingual(row.en,row.es))}
  if(parsed){
   const box=make('div','doc-function-info'),params=make('div'),returns=make('div');params.append(title('Parameters','Parámetros'));returns.append(title('Returns','Devuelve'));
   const args=splitParams(parsed[1]);if(!args.length){const p=make('p','doc-empty');p.append(bilingual('No arguments.','No recibe argumentos.'));params.append(p)}else{
    const dl=make('dl','doc-params');for(const arg of args){const m=arg.match(/^(\w+)(?:\s*:\s*(\w+))?(?:\s*=\s*([\s\S]+))?$/);const r=make('div'),dt=make('dt'),dd=make('dd');dt.append(make('code','',m?m[1]:arg));if(m){dd.append(bilingual(...(m[2]?(types[m[2]]||[m[2],m[2]]):['Inferred type','Tipo inferido'])));if(m[3]!==undefined){dd.append(' · ',bilingual('Default','Predeterminado'),': ',make('code','',m[3]))}}r.append(dt,dd);dl.append(r)}params.append(dl)
   }
   const ret=parsed[2]||'void',p=make('p','doc-empty');p.append(make('code','',ret),' · ',bilingual(...(types[ret]||[ret,ret])));returns.append(p);if(row?.returns){const meaning=make('p','doc-empty');meaning.append(bilingual(...row.returns));returns.append(meaning)}box.append(params,returns);card.append(box)
  }
  if(row?.source){const h=make('h4','doc-example-title');h.append(bilingual('Usage example','Ejemplo de uso'));const pre=make('pre','doc-usage');pre.append(make('code','language-tripol',row.source));if(row.native)pre.dataset.nativeOnly='';if(row.input)pre.dataset.input=row.input;card.append(h,pre);
   if(row.native){const p=make('p','doc-function-note');p.append(bilingual('Native TRIPOL example','Ejemplo de TRIPOL nativo'));card.append(p)}
   if(row.expected!==null&&row.expected!==undefined){const out=make('div','doc-expected'),p=make('pre');p.append(make('code','language-output',row.expected));out.append(title('Expected output','Resultado esperado'),p);card.append(out)}
  }
  if(row?.noteEn){const p=make('p','doc-function-note');p.append(bilingual(row.noteEn,row.noteEs));card.append(p)}
 }
}

learningReference();
$$('.docs-article pre').forEach(pre=>{if(pre.classList.contains('signature'))return;const shell=document.createElement('div');shell.className='code-shell';const head=document.createElement('div');head.className='code-head';const title=document.createElement('span');const mode=codeMode(pre.querySelector('code')||pre);title.textContent=mode==='tripol'?'TRIPOL':['cmd','sh','bash','shell','powershell','ps1'].includes(mode)?(mode==='powershell'?'POWERSHELL':'TERMINAL'):mode==='text'?'CODE':mode.toUpperCase();const copy=document.createElement('button');copy.className='doc-copy';copy.dataset.es='Copiar';copy.textContent='Copy';copy.addEventListener('click',async()=>{try{await navigator.clipboard.writeText(pre.textContent);toast(translate('Copied to clipboard','Copiado al portapapeles'))}catch{const range=document.createRange();range.selectNodeContents(pre);const selection=window.getSelection();selection.removeAllRanges();selection.addRange(range);toast(translate('Selected — copy manually','Seleccionado — copia manualmente'))}});head.append(title,copy);pre.before(shell);shell.append(head,pre)});
function buildToc(){const target=$('#toc-links');target.replaceChildren();$$('.docs-article h2,.docs-article h3[id],.docs-article .function-card h3').filter(h=>!h.closest('.library-api-category .function-card')).forEach((h,i)=>{const id=h.id||h.closest('.function-card')?.id||('section-'+i);if(!h.id&&!h.closest('.function-card'))h.id=id;const a=document.createElement('a');a.href='#'+id;a.textContent=h.textContent.replace(/\s*#\s*$/,'');if(h.tagName==='H3')a.className='sub';target.append(a)});}
const dialog=$('#doc-search'),input=$('#global-search');
function normalizeSearch(value){return String(value??'').normalize('NFD').replace(/[\u0300-\u036f]/g,'').toLowerCase().trim()}
function findDocs(query,entries){
 const q=normalizeSearch(query),terms=q.split(/\s+/).filter(Boolean);
 if(!q)return entries.filter(x=>!x.url.includes('#')).slice(0,30);
 return entries.map((item,i)=>{
  const title=normalizeSearch(item.title),es=normalizeSearch(item.es),detail=normalizeSearch(item.detail),url=normalizeSearch(item.url),hay=[title,es,detail,url].join(' ');
  if(!terms.every(t=>hay.includes(t)))return null;
  const score=(title===q||es===q?100:0)+(title.startsWith(q)||es.startsWith(q)?40:0)+(title.includes(q)||es.includes(q)?20:0)+(detail.includes(q)?5:0);
  return {item,i,score}
 }).filter(Boolean).sort((a,b)=>b.score-a.score||a.i-b.i).slice(0,40).map(x=>x.item)
}
function renderSearch(){
 const results=$('#search-results');if(!results||!input)return;
 const list=findDocs(input.value,typeof DOC_SEARCH==='undefined'?[]:DOC_SEARCH);
 results.replaceChildren();
 if(!list.length){const p=document.createElement('p');p.textContent=typeof DOC_SEARCH==='undefined'?translate('Search index unavailable. Reload the page.','Índice de búsqueda no disponible. Recarga la página.'):translate('No results. Try a library, topic or function name.','Sin resultados. Prueba una biblioteca, tema o función.');results.append(p);return}
 for(const item of list){const a=document.createElement('a');a.className='search-item';const u=new URL(root+item.url,location.href);u.searchParams.set('lang',lang);a.href=u.href;const title=document.createElement('strong');title.textContent=item.title;const detail=document.createElement('small');detail.textContent=lang==='es'&&item.es?item.es:item.detail;a.append(title,detail);results.append(a)}
}
function openSearch(){if(!dialog||!input)return;if(!dialog.open)dialog.showModal();renderSearch();input.focus()}
if(dialog&&input){
 $('#search-open')?.addEventListener('click',openSearch);
 $('#search-close')?.addEventListener('click',()=>dialog.close());
 input.addEventListener('input',renderSearch);
 dialog.addEventListener('click',e=>{if(e.target===dialog)dialog.close()});
 input.addEventListener('keydown',e=>{if(e.key==='ArrowDown'){e.preventDefault();$('#search-results a')?.focus()}else if(e.key==='Enter'){e.preventDefault();$('#search-results a')?.click()}});
 document.addEventListener('keydown',e=>{if((e.ctrlKey||e.metaKey)&&e.key.toLowerCase()==='k'){e.preventDefault();openSearch()}if(e.key==='Escape'&&side.classList.contains('open')){side.classList.remove('open');$('#doc-menu').setAttribute('aria-expanded','false')}});
}


/* AI reference download: shared across documentation and library pages. */
function installAIReference(){
 const key='tripol-ai-reference-dismissed-v0271';
 const notice=document.createElement('aside');notice.className='doc-ai-reference';
 const panel=document.createElement('div');panel.className='doc-ai-panel';
 const makeLabel=(tag,cls,en,es)=>{const el=document.createElement(tag);el.className=cls;el.textContent=en;el.dataset.es=es;return el};
 const toggle=makeLabel('button','doc-ai-toggle','AI reference','Referencia para IA');toggle.type='button';toggle.setAttribute('aria-controls','doc-ai-panel');
 panel.id='doc-ai-panel';
 const top=document.createElement('div');top.className='doc-ai-head';
 const kicker=makeLabel('span','doc-ai-kicker','TRIPOL / AI reference','TRIPOL / Referencia para IA');
 const close=makeLabel('button','doc-ai-close','Close','Cerrar');close.type='button';top.append(kicker,close);
 const title=makeLabel('strong','doc-ai-title','New to TRIPOL?','¿No sabes nada de TRIPOL?');
 const description=makeLabel('p','doc-ai-description','Don’t want to read all the docs? Download this file and send it to your favorite AI. Ask it to explain TRIPOL or help you write code.','¿No quieres leer toda la documentación? Descarga este archivo y envíaselo a tu IA de confianza. Pídele que te explique TRIPOL o te ayude a programar.');
 const download=makeLabel('a','doc-ai-download','Download AI reference ↓','Descargar archivo para IA ↓');
 download.href=new URL(root+'../TRIPOL-AI.md',location.href).href;download.download='TRIPOL-AI.md';
 const meta=makeLabel('small','doc-ai-meta','Markdown · Language and library reference','Markdown · Referencia del lenguaje y bibliotecas');
 panel.append(top,title,description,download,meta);notice.append(toggle,panel);
 let dismissed=false;try{dismissed=sessionStorage.getItem(key)==='1'}catch{}
 const setOpen=(open,remember=false)=>{panel.hidden=!open;toggle.hidden=open;toggle.setAttribute('aria-expanded',String(open));if(remember)try{sessionStorage.setItem(key,open?'0':'1')}catch{}};
 setOpen(!dismissed&&!window.matchMedia('(max-width: 700px)').matches);
 close.addEventListener('click',()=>{setOpen(false,true);toggle.focus()});
 toggle.addEventListener('click',()=>{setOpen(true,true);close.focus()});
 panel.addEventListener('keydown',e=>{if(e.key==='Escape'){e.preventDefault();setOpen(false,true);toggle.focus()}});
 document.body.append(notice);
}
installAIReference();

installAllLive();registerLabels();apply();installLibraryCategoryTracking();$('#doc-lang').addEventListener('click',()=>{lang=lang==='en'?'es':'en';apply()});
$('#doc-menu').addEventListener('click',()=>{$('#doc-menu').setAttribute('aria-expanded',String(side.classList.toggle('open')))});
$('#side-filter').addEventListener('input',e=>{const q=e.target.value.trim().toLowerCase();filterLibraryNavigation(q)});
const lf=$('#library-filter');if(lf)lf.addEventListener('input',()=>{let n=0;const q=lf.value.trim().toLowerCase();$$('.library-card').forEach(a=>{a.hidden=!a.dataset.filter.toLowerCase().includes(q);if(!a.hidden)n++});$('#library-empty').hidden=!!n;refreshLevelHeadings($('#library-grid'))});

const active=side.querySelector('[aria-current=page]');if(active&&active.offsetTop>side.clientHeight-90)side.scrollTop=active.offsetTop-100;
})();
