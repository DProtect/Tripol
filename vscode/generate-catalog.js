"use strict";
const fs=require("fs"),path=require("path");
const root=path.resolve(__dirname,"../stdlib"),catalog={};
for(const file of fs.readdirSync(root).filter(n=>n.endsWith(".tri")).sort()){
 const module=file.slice(0,-4);catalog[module]={};
 const source=fs.readFileSync(path.join(root,file),"utf8");
 for(const match of source.matchAll(/^\s*func\s+(\w+(?:\.\w+)?)\s*(\([^\n]*?\)\s*(?:->\s*(?:list\[[^\n\]]+\]|\w+))?)\s*\{/gm)){
  const [receiver,method]=match[1].split(".");
  if(method){
   const type=module+"."+receiver;
   (catalog[type]??={})[method]=receiver+"."+method+match[2].replace(/^\(\s*\w+\s*:\s*\w+\s*,?\s*/,"(");
  }else{
   catalog[module][receiver]="func "+module+"."+receiver+match[2];
  }
 }
}
fs.writeFileSync(path.join(__dirname,"catalog.json"),JSON.stringify(catalog,null,2)+"\n");
