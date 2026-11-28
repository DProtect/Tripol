"use strict";
const vscode = require("vscode");
const cp = require("child_process");
const path = require("path");
const catalog = require("./catalog.json");
catalog.json.decode = "json.decode(YourStruct, jsonText) -> YourStruct";
catalog.json.encode = "json.encode(YourStruct, value) -> string";
const gatewayChoices = require("./discord-gateway.json");
function activate(context) {
 const output = vscode.window.createOutputChannel("TRIPOL");
 const diagnostics = vscode.languages.createDiagnosticCollection("tripol");
 const selector = {language:"tripol",scheme:"file"};
 const keyword = "import if else while repeat for in break continue return show assert try catch finally throw scratch func struct let const int bool string float ints list true false null some async await and or not".split(" ");
 const recordType=(source,name)=>{
  const annotated=[...source.matchAll(/\b([A-Za-z_]\w*)\s*:\s*(\w+\.\w+)\b/g)].reverse().find(m=>m[1]===name);
  if(annotated&&catalog[annotated[2]])return annotated[2];
  for(const binding of [...source.matchAll(/^\s*(?:let\s+)?([A-Za-z_]\w*)\s*=\s*(\w+)\.(\w+)\s*\(/gm)].reverse()){
   if(binding[1]!==name)continue;
   const [, , receiver, factory]=binding;
   const owner=catalog[receiver]?receiver:recordType(source.slice(0,binding.index),receiver);
   const result=catalog[owner]?.[factory]?.match(/->\s*([\w.]+)$/)?.[1];
   if(result)return result.includes(".")?result:owner.split(".")[0]+"."+result;
  }
 };
 const members=(doc,pos,name)=>catalog[name]||catalog[recordType(doc.getText(new vscode.Range(0,0,pos.line,pos.character)),name)];
 const localNames = doc => [...new Set([...doc.getText().matchAll(/\b(?:func\s+)?([A-Za-z_]\w*)\s*(?:\(|=(?!=))/g)].map(x=>x[1]))];
 context.subscriptions.push(output,diagnostics,vscode.languages.registerCompletionItemProvider(selector,{
  provideCompletionItems(doc,pos){
   const line=doc.lineAt(pos.line).text.slice(0,pos.character);
   const gatewayArgument=line.match(/^\s*(\w+)\.(on|intent)\(\s*"(\w*)$/);
   if(gatewayArgument){
    const [,receiver,method,prefix]=gatewayArgument;
    if(members(doc,pos,receiver)?.[method]?.startsWith("GatewayEvents.")){
     const range=new vscode.Range(pos.line,pos.character-prefix.length,pos.line,pos.character);
     return gatewayChoices[method].filter(value=>value.startsWith(prefix)).map(value=>{
      const item=new vscode.CompletionItem(value,vscode.CompletionItemKind.EnumMember);
      item.range=range;
      return item;
     });
    }
   }
   const moduleMatch=line.match(/\bimport\s+tripol\.[\w]*$/);
   const member=line.match(/\b(\w+)\.\w*$/);
   if(moduleMatch)return Object.keys(catalog).filter(n=>!n.includes(".")).map(n=>new vscode.CompletionItem(n,vscode.CompletionItemKind.Module));
   const entries=member&&members(doc,pos,member[1]);
   if(entries)return Object.entries(entries).map(([name,detail])=>{const item=new vscode.CompletionItem(name,vscode.CompletionItemKind.Function);item.detail=detail;return item;});
   return [...keyword.map(n=>new vscode.CompletionItem(n,vscode.CompletionItemKind.Keyword)),...localNames(doc).map(n=>new vscode.CompletionItem(n,vscode.CompletionItemKind.Variable))];
  }
 },".","\""),vscode.languages.registerHoverProvider(selector,{
  provideHover(doc,pos){const range=doc.getWordRangeAtPosition(pos,/\w+\.\w+/);if(!range)return;const [m,n]=doc.getText(range).split(".");const signature=members(doc,pos,m)?.[n];if(signature){const md=new vscode.MarkdownString();md.appendCodeblock(signature,"tripol");return new vscode.Hover(md,range);}}
 }),vscode.languages.registerDefinitionProvider(selector,{
  provideDefinition(doc,pos){const r=doc.getWordRangeAtPosition(pos);if(!r)return;const name=doc.getText(r);for(let i=0;i<doc.lineCount;i++){const m=doc.lineAt(i).text.match(/^\s*func\s+(\w+)\s*\(/);if(m&&m[1]===name)return new vscode.Location(doc.uri,new vscode.Position(i,0));}}
 }));
 context.subscriptions.push(vscode.languages.registerDocumentFormattingEditProvider(selector,{
  provideDocumentFormattingEdits(doc){
   if(!vscode.workspace.isTrusted)return [];
   const exe=vscode.workspace.getConfiguration("tripol",doc.uri).get("compilerPath","tripol");
   return new Promise((resolve,reject)=>{
    const child=cp.execFile(exe,["fmt","--stdin"],{cwd:path.dirname(doc.uri.fsPath),timeout:30000,maxBuffer:4*1024*1024,windowsHide:true,env:{...process.env,TRIPOL_NO_UPDATE:"1"}},(err,stdout,stderr)=>{
     if(err){reject(new Error(stderr||err.message));return;}
     const end=doc.lineAt(doc.lineCount-1).range.end;
     resolve(stdout===doc.getText()?[]:[vscode.TextEdit.replace(new vscode.Range(new vscode.Position(0,0),end),stdout)]);
    });
    child.stdin.on("error",()=>{});
    child.stdin.end(doc.getText(),"utf8");
   });
  }
 }));
 const invoke=(doc,mode)=>{
  if(!vscode.workspace.isTrusted){vscode.window.showWarningMessage("TRIPOL requires a trusted workspace for compiler commands.");return;}
  if(!doc||doc.languageId!=="tripol"||doc.uri.scheme!=="file")return;
  if(doc.isDirty){vscode.window.showWarningMessage("Save the TRIPOL file first.");return;}
  const exe=vscode.workspace.getConfiguration("tripol",doc.uri).get("compilerPath","tripol");
  const args=[mode,doc.uri.fsPath];if(mode==="check")args.push("--diagnostics=json");
  // execFile uses an argument vector, not shell interpolation.
  cp.execFile(exe,args,{cwd:path.dirname(doc.uri.fsPath),timeout:30000,maxBuffer:1024*1024,windowsHide:true,env:{...process.env,TRIPOL_NO_UPDATE:"1"}},(err,stdout,stderr)=>{
   output.appendLine("TRIPOL "+mode+" "+doc.uri.fsPath);output.append(stdout||"");output.append(stderr||"");
   if(mode==="check"){
    const list=[];if(err){let found=false;for(const line of (stderr||"").split(/\r?\n/)){try{const d=JSON.parse(line);if(d.message){const i=Math.min(doc.lineCount-1,Math.max(0,(d.line||1)-1));const col=Math.min(doc.lineAt(i).text.length,Math.max(0,(d.column||1)-1));const item=new vscode.Diagnostic(new vscode.Range(i,col,i,Math.min(doc.lineAt(i).text.length,col+1)),[d.message,d.hint,d.suggestion&&"Did you mean "+d.suggestion+"?"].filter(Boolean).join("\n"),vscode.DiagnosticSeverity.Error);item.code=d.code;list.push(item);found=true;}}catch{}}
    if(!found)list.push(new vscode.Diagnostic(new vscode.Range(0,0,0,1),err.message,vscode.DiagnosticSeverity.Error));}
    diagnostics.set(doc.uri,list);
   }else{output.show(true);if(err)vscode.window.showErrorMessage("TRIPOL "+mode+" failed; see output.");}
  });
 };
 for(const mode of ["check","build","run"])context.subscriptions.push(vscode.commands.registerCommand("tripol."+mode,()=>invoke(vscode.window.activeTextEditor?.document,mode)));
 context.subscriptions.push(vscode.workspace.onDidSaveTextDocument(doc=>{if(vscode.workspace.getConfiguration("tripol",doc.uri).get("checkOnSave",false))invoke(doc,"check");}),vscode.workspace.onDidCloseTextDocument(doc=>diagnostics.delete(doc.uri)));
}
module.exports={activate,deactivate(){}};
