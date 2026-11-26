(function(root){
'use strict';
// R3 ABI: string handles are interned per execution, with literal data carried
// in the program's tripol.strings custom section. No external I/O.
function execute(binary,options={}){
 const module=new WebAssembly.Module(binary), sections=WebAssembly.Module.customSections(module,'tripol.strings');
 if(sections.length!==1)throw Error('Missing or duplicate string table');
 const literals=JSON.parse(new TextDecoder('utf-8',{fatal:true}).decode(sections[0]));
 if(!Array.isArray(literals)||literals.some(v=>typeof v!=='string'))throw Error('Invalid string table');
 const strings=[],ids=new Map(),encoder=new TextEncoder();let allocated=0,output='',steps=0,cells=0;const collections=[];const lines=String(options.input||'').split(/\r?\n/);let cursor=0;const bits=new DataView(new ArrayBuffer(8));
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
   if(op===21)return intern(scalar(a,Number(b)));
   if(op===22||op===57){emit(decoded(get(a))+'\n');const text=bytes(lines[cursor++]??'');if(text.length>4096)fail(line,'Input limit reached');if(op===57){if(!/^[+-]?\d+$/.test(text))fail(line,'Invalid integer input');return checked(BigInt(text),line)}return intern(text)}
   if(op===50)return intern(String(a));
   if(op===51){const s=get(a);if(!/^[+-]?\d+$/.test(s))fail(line,'Invalid integer text');return checked(BigInt(s),line)}
   if(op===52){const s=get(a);if(b<0n||b>=BigInt(s.length))fail(line,'Byte index out of bounds');return BigInt(s.charCodeAt(Number(b)))}
   if(op===53){const s=get(a);if(b<0n||c<b||c>BigInt(s.length))fail(line,'String slice out of bounds');return intern(s.slice(Number(b),Number(c)))}
   if(op===54)return get(a).startsWith(get(b))?1n:0n;
   if(op===55)return get(a).endsWith(get(b))?1n:0n;
   if(op===56)return BigInt(get(a).indexOf(get(b)));
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
 try{new WebAssembly.Instance(module,imports).exports.run();return{output,steps}}catch(e){e.output=output;throw e}
}
root.TripolWebHost={execute};if(typeof module!=='undefined')module.exports={execute};
})(globalThis);
