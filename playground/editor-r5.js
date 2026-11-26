(function editorTools(root){
'use strict';
const keywords=new Set('show let const func return if else while for in repeat break continue assert and or not import struct try catch finally throw scratch'.split(' '));
const types=new Set('int bool float string ints list string_list bool_list float_list'.split(' '));
const builtins=new Set('input input_int len sum range ints abs min max sqrt float text parse_int slice'.split(' '));
const reserved=new Set([...keywords,...types,...builtins,'true','false']);
function tokenize(src){const pattern=/(\/\/[^\n]*|#[^\n]*|f?"(?:\\.|[^"\\\n])*"?|\b\d+(?:\.\d+)?\b|[A-Za-z_][A-Za-z_0-9]*|==|!=|<=|>=|=>|->|\+=|-=|\*=|\/=|%=|&&|\|\||[+\-*/%=<>!]|[{}()[\],.:;]|\s+|.)/gu;let out=[];for(const m of src.matchAll(pattern)){const v=m[0];let kind='plain';if(v.startsWith('//')||v.startsWith('#'))kind='comment';else if(v.startsWith('"')||v.startsWith('f"'))kind='string';else if(/^\d/.test(v))kind='number';else if(v==='true'||v==='false')kind='literal';else if(keywords.has(v))kind='keyword';else if(types.has(v))kind='type';else if(builtins.has(v)||/^[A-Za-z_]/.test(v)&&/^\s*\(/.test(src.slice(m.index+v.length)))kind='function';else if(/^[A-Za-z_]/.test(v)&&/^\s*\./.test(src.slice(m.index+v.length)))kind='namespace';else if(/^[+\-*/%=<>!&|]/.test(v))kind='operator';else if(/^[{}()[\],.:;]$/.test(v))kind='punct';out.push({v,kind})}return out}
function suggestions(src,pos){
 const before=src.slice(0,pos),line=before.split('\n').at(-1),libs=root.TripolLibraries||[];
 if(/^\s*(?:\/\/|#)/.test(line)||tokenize(line).at(-1)?.kind==='string')return {prefix:'',items:[],details:{}};
 const imp=line.match(/^\s*import\s+tripol\.([A-Za-z_0-9]*)$/);
 if(imp){const prefix=imp[1];return {prefix,items:libs.map(l=>l.name).filter(n=>n.startsWith(prefix)&&n!==prefix).slice(0,12),details:{}}}
 const member=line.match(/\b([A-Za-z_]\w*)\.([A-Za-z_0-9]*)$/);
 if(member){const l=libs.find(x=>x.name===member[1]);if(l){const prefix=member[2],fs=l.functions.filter(f=>f.name.startsWith(prefix)&&f.name!==prefix).slice(0,12);return {prefix,items:fs.map(f=>f.name),details:Object.fromEntries(fs.map(f=>[f.name,f])),library:l.name,member:true}}}
 const prefix=(before.match(/[A-Za-z_][A-Za-z_0-9]*$/)||[''])[0];
 const known=[...keywords,...types,...builtins,'true','false',...libs.map(l=>l.name)];
 const local=tokenize(src).filter(t=>t.kind==='plain'||t.kind==='function').map(t=>t.v).filter(v=>/^[A-Za-z_][A-Za-z_0-9]*$/.test(v)&&!reserved.has(v));
 return {prefix,items:[...new Set([...known,...local])].filter(v=>v.startsWith(prefix)&&v!==prefix).slice(0,12),details:{}}
}
function format(src){let depth=0;return src.split('\n').map(line=>{const trimmed=line.trim();if(!trimmed)return '';const tokens=tokenize(trimmed).filter(t=>t.kind!=='string'&&t.kind!=='comment'&&t.v.trim());let lead=0;while(tokens[lead]?.v==='}')lead++;const indent=Math.max(0,depth-lead);for(const t of tokens){if(t.v==='{')depth++;if(t.v==='}')depth=Math.max(0,depth-1)}return '    '.repeat(Math.min(indent,128))+trimmed}).join('\n')}
function location(src,offset){const a=src.slice(0,offset).split('\n');return {line:a.length,column:a.at(-1).length+1}}
function lineRange(src,line){const lines=src.split('\n');const i=Math.max(0,Math.min(lines.length-1,line-1));const start=lines.slice(0,i).reduce((n,s)=>n+s.length+1,0);return[start,start+lines[i].length]}
const api={tokenize,suggestions,format,location,lineRange};root.TripolEditor=api;if(typeof module!=='undefined')module.exports=api;
})(globalThis);
