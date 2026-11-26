(function(root){
'use strict';
// R2 ABI: string handles are interned per execution, with literal data carried
// in the program's tripol.strings custom section. No external I/O.
function execute(binary,options={}){
 const module=new WebAssembly.Module(binary), sections=WebAssembly.Module.customSections(module,'tripol.strings');
 if(sections.length!==1)throw Error('Missing or duplicate string table');
 const literals=JSON.parse(new TextDecoder('utf-8',{fatal:true}).decode(sections[0]));
 if(!Array.isArray(literals)||literals.some(v=>typeof v!=='string'))throw Error('Invalid string table');
 const strings=[],ids=new Map(),encoder=new TextEncoder();let allocated=0,output='',steps=0;
 const limit=options.steps||100000,MIN=-(1n<<63n),MAX=(1n<<63n)-1n;
 function intern(s){if(ids.has(s))return ids.get(s);const n=encoder.encode(s).length;if(n>65535||allocated+n>4*1024*1024||strings.length>=50000)throw Error('String allocation limit reached');const id=BigInt(strings.length);allocated+=n;strings.push(s);ids.set(s,id);return id}
 for(const s of literals){if(intern(s)!==BigInt(strings.length-1))throw Error('Duplicate literal')}
 function get(id){if(id<0n||id>=BigInt(strings.length))throw Error('Invalid string handle');return strings[Number(id)]}
 const imports={tripol:{
  tick(line){if(++steps>limit)throw Error('Line '+line+': Execution limit reached')},
  arithmetic(op,a,b,line){if(op===5)return intern(get(a)+get(b));if(op===6)return BigInt(encoder.encode(get(a)).length);
   let v;if(op===0)v=a+b;else if(op===1)v=a-b;else if(op===2)v=a*b;else if(op===3||op===4){if(b===0n||a===MIN&&b===-1n)throw Error('Line '+line+': Invalid integer arithmetic');v=op===3?a/b:a%b}else throw Error('Unknown arithmetic operation');if(v<MIN||v>MAX)throw Error('Line '+line+': Integer overflow');return v},
  show(v,kind,last){const t=(kind===2?get(v):kind===1?String(v!==0n):String(v))+(last?'\n':' ');if(output.length+t.length>32768)throw Error('Console output limit reached');output+=t},
  assertFail(line){throw Error('Line '+line+': Assertion failed')}
 }};
 try{new WebAssembly.Instance(module,imports).exports.run();return{output,steps}}catch(e){e.output=output;throw e}
}
root.TripolWebHost={execute};if(typeof module!=='undefined')module.exports={execute};
})(globalThis);
