"use strict";
const fs=require("fs"),path=require("path");
const root=path.resolve(__dirname,"../stdlib"),catalog={};
for(const file of fs.readdirSync(root).filter(n=>n.endsWith(".tri")).sort()){
 const module=file.slice(0,-4);catalog[module]={};
 const source=fs.readFileSync(path.join(root,file),"utf8");
 for(const match of source.matchAll(/^\s*func\s+(\w+)\s*(\([^\n]*?\)\s*(?:->\s*\w+)?)\s*\{/gm)){
  catalog[module][match[1]]="func "+module+"."+match[1]+match[2];
 }
}
fs.writeFileSync(path.join(__dirname,"catalog.json"),JSON.stringify(catalog,null,2)+"\n");
