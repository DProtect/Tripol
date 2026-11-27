(function(root){
'use strict';
// R3 ABI: string handles are interned per execution, with literal data carried
// in the program's tripol.strings custom section. Files are per-run virtual storage only.
function execute(binary,options={}){
 const module=new WebAssembly.Module(binary), sections=WebAssembly.Module.customSections(module,'tripol.strings');
 if(sections.length!==1)throw Error('Missing or duplicate string table');
 const literals=JSON.parse(new TextDecoder('utf-8',{fatal:true}).decode(sections[0]));
 if(!Array.isArray(literals)||literals.some(v=>typeof v!=='string'))throw Error('Invalid string table');
 const strings=[],ids=new Map(),encoder=new TextEncoder();let allocated=0,output='',steps=0,cells=0;const collections=[];const lines=String(options.input||'').split(/\r?\n/);let cursor=0;const virtualFiles=new Map();let fileBytes=0,atomicExpected=null;let windowState=null,frames=[],frameBytes=0;const bits=new DataView(new ArrayBuffer(8));
 const limit=options.steps||100000,MIN=-(1n<<63n),MAX=(1n<<63n)-1n;
 function bytes(s){let out='';for(const b of encoder.encode(s))out+=String.fromCharCode(b);return out}
 function decoded(s){return new TextDecoder().decode(Uint8Array.from(s,c=>c.charCodeAt(0)))}
 function intern(s){if(ids.has(s))return ids.get(s);const n=s.length;if(n>65535||allocated+n>4*1024*1024||strings.length>=50000)throw Error('String allocation limit reached');const id=BigInt(strings.length);allocated+=n;strings.push(s);ids.set(s,id);return id}
 for(const s of literals){if(intern(bytes(s))!==BigInt(strings.length-1))throw Error('Duplicate literal')}
 function get(id){if(id<0n||id>=BigInt(strings.length))throw Error('Invalid string handle');return strings[Number(id)]}
 function fail(line,msg){throw Error('Line '+line+': '+msg)}
 function checked(v,line){if(v<MIN||v>MAX)fail(line,'Integer overflow');return v}
 function emit(s){if(output.length+s.length>32768)throw Error('Console output limit reached');output+=s}
 function f64(v){bits.setBigInt64(0,v,true);return bits.getFloat64(0,true)}
 function i64(v){bits.setFloat64(0,v,true);return bits.getBigInt64(0,true)}
 function numeric(n,line){if(!Number.isFinite(n))fail(line,'Invalid floating point result');return i64(n)}
 function collection(id,line){if(id<0n||id>=BigInt(collections.length))fail(line,'Invalid collection handle');return collections[Number(id)]}
 function index(c,i,line){if(i<0n||i>=BigInt(c.items.length))fail(line,'Array index out of bounds');return Number(i)}
 function filePath(v,line){const p=decoded(get(v));if(!p||p.length>256||p.startsWith('/')||p.includes('\\')||p.includes(':')||p.includes('\0')||p.split('/').some(x=>!x||x==='.'||x==='..'))fail(line,'Only relative virtual file paths are allowed');return p}
 function storeFile(p,s,line){const old=virtualFiles.get(p)||'';if(s.length>65536||fileBytes-old.length+s.length>524288||(!virtualFiles.has(p)&&virtualFiles.size>=64))fail(line,'Virtual filesystem limit reached');fileBytes+=s.length-old.length;virtualFiles.set(p,s);return BigInt(s.length)}
 function loadFile(p,line){if(!virtualFiles.has(p))fail(line,'Virtual file not found: '+p);return virtualFiles.get(p)}
 function newArray(values,line){reserve(values.length,line);if(collections.length>=50000)fail(line,'Collection count limit');const id=BigInt(collections.length);collections.push({kind:4,capacity:values.length,items:values});return id}
 function reserve(n,line){if(n<0||cells+n>262144)fail(line,'Collection allocation limit reached');cells+=n}
 function scalar(v,k){if(k===2)return get(v);if(k===1)return String(v!==0n);if(k===3){let n=f64(v);if(!Number.isFinite(n))throw Error('Invalid floating point display');return n.toFixed(6)}return String(v)}
 function cmpText(a,b){const x=get(a),y=get(b);return x<y?-1:x>y?1:0}
 function display(v,kind,line){if(kind<4)return kind===2?decoded(get(v)):scalar(v,kind);const c=collection(v,line);let result='[';for(let i=0;i<c.items.length;i++){if(i)result+=', ';result+=c.kind===7?decoded(get(c.items[i])):scalar(c.items[i],c.kind===4?0:c.kind-5);if(result.length>32768)fail(line,'Console output limit reached')}return result+']'}
 const imports={tripol:{
  tick(line){if(++steps>limit)throw Error('Line '+line+': Execution limit reached')},
  arithmetic(op,a,b,line){if(op===5)return intern(get(a)+get(b));if(op===6)return BigInt(get(a).length);
   let v;if(op===0)v=a+b;else if(op===1)v=a-b;else if(op===2)v=a*b;else if(op===3||op===4){if(b===0n||a===MIN&&b===-1n)throw Error('Line '+line+': Invalid integer arithmetic');v=op===3?a/b:a%b}else throw Error('Unknown arithmetic operation');if(v<MIN||v>MAX)throw Error('Line '+line+': Integer overflow');return v},
  show(v,kind,last){emit(display(v,kind,0)+(last?'\n':' '))},
  value(op,a,b,c,line){
   if(op===70)return BigInt(allocated+cells*8);
   if(op>=71&&op<=75){const p=filePath(a,line);
    if(op===71){const s=loadFile(p,line);if(s.length>65535||s.includes('\0'))fail(line,'Invalid text file');return intern(s)}
    if(op===72||op===73){const s=get(b);if(s.includes('\0'))fail(line,'NUL is not supported by text files');const old=op===73?(virtualFiles.get(p)||''):'';storeFile(p,old+s,line);return BigInt(s.length)}
    if(op===74){if(b<0n||b>65536n)fail(line,'Invalid read limit');const s=loadFile(p,line);if(BigInt(s.length)>b)fail(line,'File exceeds read limit');return newArray(Array.from(s,x=>BigInt(x.charCodeAt(0))),line)}
    const data=collection(b,line).items;if(data.length>65536||data.some(v=>v<0n||v>255n))fail(line,'Invalid binary file');let content='';for(const v of data)content+=String.fromCharCode(Number(v));return storeFile(p,content,line)
   }
   if(op===82){if(windowState)return 0n;const w=Number(b),h=Number(c);if(b<1n||c<1n||b>512n||c>512n||w*h>262144)fail(line,'Web graphics dimensions: 1–512 pixels');windowState={title:decoded(get(a)).slice(0,128),width:w,height:h,mode:null};return 1n}
   if(op===83)return windowState?1n:0n;
   if(op===85){const wasOpen=!!windowState;windowState=null;return wasOpen?1n:0n}
   if(op===84||op===86){
    if(!windowState)return 0n;
    const w=windowState.width,h=windowState.height,mode=op===84?'rgb':'triangles';
    if(windowState.mode&&windowState.mode!==mode)fail(line,'Do not mix graphics presenters');
    const src=collection(a,line);if(src.kind!==4)fail(line,'Graphics requires an integer array');
    let frame,size;
    if(op===84){
     if(b!==BigInt(w)||c!==BigInt(h)||src.items.length!==w*h)fail(line,'Invalid graphics framebuffer dimensions');
     size=w*h*4;if(src.items.some(v=>v<0n||v>0xffffffn))fail(line,'Invalid RGB24 pixel');
     const rgba=new Uint8Array(size);for(let i=0;i<src.items.length;i++){const v=Number(src.items[i]);rgba[i*4]=v>>>16;rgba[i*4+1]=(v>>>8)&255;rgba[i*4+2]=v&255;rgba[i*4+3]=255}
     frame={kind:mode,width:w,height:h,rgba,title:windowState.title,delay:16};
    }else{
     if(b<0n||b>0xffffffn||src.items.length%9!==0||src.items.length>9000)fail(line,'Invalid triangle array or background');
     for(let i=0;i<src.items.length;i++){const v=src.items[i];if(i%3===2?(v<0n||v>0xffffffn):(v< -1000n||v>1000n))fail(line,'Invalid triangle coordinate or RGB')}
     const vertices=Int32Array.from(src.items,Number);size=vertices.byteLength;frame={kind:mode,width:w,height:h,vertices,background:Number(b),title:windowState.title,delay:16};
    }
    if(frames.length>=120||frameBytes+size>16*1024*1024)fail(line,'Graphics preview limit: 120 frames or 16 MiB');
    frameBytes+=size;windowState.mode=mode;frames.push(frame);return 1n;
   }
   if(op===81){if(a<0n||a>1000n)fail(line,'Web sleep limit: 1000 ms');if(windowState){if(frames.length)frames[frames.length-1].delay=Math.max(16,Number(a));return 0n}if(a!==0n){if(typeof SharedArrayBuffer==='undefined'||typeof Atomics==='undefined')fail(line,'Blocking sleep requires an isolated worker');Atomics.wait(new Int32Array(new SharedArrayBuffer(4)),0,0,Number(a))}return 0n}
   if(op===21)return intern(scalar(a,Number(b)));
   if(op===22||op===57){if(b===0n)emit(decoded(get(a))+'\n');const text=bytes(lines[cursor++]??'');if(text.length>4096)fail(line,'Input limit reached');if(op===57){if(!/^[+-]?\d+$/.test(text))fail(line,'Invalid integer input');return checked(BigInt(text),line)}return intern(text)}
   if(op===50)return intern(String(a));
   if(op===51){const s=get(a);if(!/^[+-]?\d+$/.test(s))fail(line,'Invalid integer text');return checked(BigInt(s),line)}
   if(op===52){const s=get(a);if(b<0n||b>=BigInt(s.length))fail(line,'Byte index out of bounds');return BigInt(s.charCodeAt(Number(b)))}
   if(op===53){const s=get(a);if(b<0n||c<b||c>BigInt(s.length))fail(line,'String slice out of bounds');return intern(s.slice(Number(b),Number(c)))}
   if(op===54)return get(a).startsWith(get(b))?1n:0n;
   if(op===55)return get(a).endsWith(get(b))?1n:0n;
   if(op===56)return BigInt(get(a).indexOf(get(b)));
   if(op===88){
    if(a<=0n||a>0x10ffffn||a>=0xd800n&&a<=0xdfffn)fail(line,'Invalid Unicode scalar');
    return intern(bytes(String.fromCodePoint(Number(a))));
   }
   if(op===87){
    const s=get(a).replace(/["\\\x00-\x1f]/g,c=>c==='"'||c==='\\'?'\\'+c:'\\u00'+c.charCodeAt(0).toString(16).padStart(2,'0'));
    return intern('"'+s+'"');
   }
   if(op>=60&&op<=63){if(a<0n||a>0xffffffffn||b<0n||b>(op>=62?31n:0xffffffffn))fail(line,'Invalid u32 argument');if(op===60)return a^b;if(op===61)return a&b;if(op===62)return ((a>>b)|(a<<(32n-b)))&0xffffffffn;return a>>b}
   if(op===58||op===59)return BigInt(Math.floor(performance.now()*(op===59?1000000:1)));
   if(op===30){const x=f64(a),y=f64(b);return numeric([()=>x+y,()=>x-y,()=>x*y,()=>x/y][Number(c)](),line)}
   if(op===31){const x=f64(a),y=f64(b);return [x===y,x!==y,x<y,x>y,x<=y,x>=y][Number(c)]?1n:0n}
   if(op===32)return numeric(Number(a),line);
   if(op===33)return numeric(-f64(a),line);
   if(op===34)return numeric(Math.sqrt(f64(a)),line);
   if(op===35)return checked(a<0n?-a:a,line);
   if(op===36)return a<b?a:b;
   if(op===37)return a>b?a:b;
   if(op===10){if(a<0n||a>262144n||b<4n||b>8n)fail(line,'Collection allocation limit reached');const n=Number(a),kind=Number(b);reserve(n,line);if(collections.length>=50000)fail(line,'Collection count limit');const id=BigInt(collections.length);collections.push({kind,capacity:n,items:kind===4?Array(n).fill(0n):[]});return id}
   const list=collection(a,line),items=list.items;
   switch(op){
    case 76:return items[index(list,b,line)];
    case 77:{const i=index(list,b,line),old=items[i];items[i]=c;return old}
    case 78:{const i=index(list,b,line),v=checked(items[i]+c,line);items[i]=v;return v}
    case 79:{const i=index(list,b,line);atomicExpected={a,i,value:c};return 0n}
    case 80:{const i=index(list,b,line),expected=atomicExpected;atomicExpected=null;if(!expected||expected.a!==a||expected.i!==i)fail(line,'Invalid compare_exchange state');if(items[i]!==expected.value)return 0n;items[i]=c;return 1n}
    case 43:{const other=collection(b,line).items;return items.length===other.length&&items.every((x,i)=>x===other[i])?1n:0n}
    case 44:items.fill(b);return BigInt(items.length);
    case 45:if(b<0n||c<0n||b+c>BigInt(items.length))fail(line,'Array fill range out of bounds');return (b<<32n)|c;
    case 46:{const start=Number(b>>32n),count=Number(b&0xffffffffn);if(start<0||count<0||start+count>items.length)fail(line,'Array fill range out of bounds');items.fill(c,start,start+count);return BigInt(count)}
    case 40:return BigInt(items.indexOf(b));
    case 41:{reserve(items.length,line);if(collections.length>=50000)fail(line,'Collection count limit');const id=BigInt(collections.length);collections.push({kind:4,capacity:items.length,items:items.slice()});return id}
    case 42:{const src=collection(b,line).items;if(src.length!==items.length)return 0n;for(let i=0;i<src.length;i++)items[i]=src[i];return 1n}
    case 11:return items[index(list,b,line)];
    case 12:items[index(list,b,line)]=c;return c;
    case 13:return BigInt(items.length);
    case 14:if(list.kind===4)fail(line,'Fixed array cannot grow');if(items.length>=list.capacity){const cap=Math.max(1,list.capacity*2);reserve(cap,line);list.capacity=cap}items.push(b);return BigInt(items.length);
    case 15:return items.splice(index(list,b,line),1)[0];
    case 16:items.length=0;return 0n;
    case 17:return BigInt(list.capacity);
    case 18:return items.includes(b)?1n:0n;
    case 19:{const i=items.indexOf(b);if(i<0)return 0n;items.splice(i,1);return 1n}
    case 20:{let n=0n;for(const v of items)n=checked(n+v,line);return n}
    case 23:if(items.length>4096)fail(line,'Sort limit reached');items.sort(list.kind===7?cmpText:(x,y)=>x<y?-1:x>y?1:0);return BigInt(items.length);
   }
   fail(line,'Unknown value operation')
  },
  assertFail(line){throw Error('Line '+line+': Assertion failed')}
 }};
 try{new WebAssembly.Instance(module,imports).exports.run();return{output,steps,frames,files:[...virtualFiles].map(([name,data])=>({name,bytes:Uint8Array.from(data,c=>c.charCodeAt(0))}))}}catch(e){e.output=output;e.frames=frames;throw e}
}
root.TripolWebHost={execute};if(typeof module!=='undefined')module.exports={execute};
})(globalThis);
