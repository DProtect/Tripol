package core

import (
 "net/url"
 "path"
 "regexp"
 "sort"
 "strconv"
 "strings"
)

// Open buffers only: no filesystem, subprocess, network, or Vault access.
func editorFileURI(raw string)*url.URL{
 u,e:=url.Parse(raw)
 if e!=nil||u.Scheme!="file"||u.User!=nil||u.RawQuery!=""||u.Fragment!=""||!strings.HasPrefix(u.Path,"/"){return nil}
 return u
}
func resolveEditorImport(root,importer,name string,docs map[string]lspDocument)(lspDocument,string){
 if name==""||path.IsAbs(name)||path.Ext(name)!=".tri"||strings.ContainsAny(name,":\\\x00"){return lspDocument{},"Invalid local import: "+name}
 r,u:=editorFileURI(root),editorFileURI(importer)
 if r==nil||u==nil||!strings.EqualFold(r.Host,u.Host){return lspDocument{},"Open local file URIs to resolve "+name}
 dest:=path.Clean(path.Join(path.Dir(u.Path),name))
 dir:=path.Dir(path.Clean(r.Path))
 if dir!="/"&&!strings.HasPrefix(dest,dir+"/"){return lspDocument{},"Import outside entry directory: "+name}
 keys:=make([]string,0,len(docs));for k:=range docs{keys=append(keys,k)};sort.Strings(keys)
 for _,k:=range keys{
  candidate:=editorFileURI(k)
  if candidate!=nil&&strings.EqualFold(candidate.Host,r.Host)&&path.Clean(candidate.Path)==dest{return docs[k],""}
 }
 return lspDocument{},"Open imported file to validate: "+name
}
var editorLocalImport=regexp.MustCompile(`(?m)^[ \t]*import[ \t]+("(?:[^"\\\r\n]|\\.)*")[ \t]*(?:;|//|#|\r?$)`)
func editorVisibleDocs(root lspDocument,docs map[string]lspDocument)[]lspDocument{
 seen:=map[string]bool{}
 var result []lspDocument
 var visit func(lspDocument)
 visit=func(doc lspDocument){
  if seen[doc.URI]||len(seen)>=32{return}
  seen[doc.URI]=true;result=append(result,doc)
  // Navigation can work even in an incomplete buffer. Validation uses the
  // real parser below instead of this best-effort import scan.
  for _,m:=range editorLocalImport.FindAllStringSubmatch(doc.Text,-1){
   name,e:=strconv.Unquote(m[1]);if e!=nil{continue}
   child,warning:=resolveEditorImport(root.URI,doc.URI,name,docs);if warning==""{visit(child)}
  }
 }
 visit(root);return result
}
func editorProject(root lspDocument,docs map[string]lspDocument)([]fn,[]stmt,string){
 libs:=standardLoader()
 states:=map[string]int{}
 var fs []fn;var main []stmt
 total:=0;warning:=""
 var visit func(lspDocument,bool)
 visit=func(doc lspDocument,entry bool){
  if states[doc.URI]==1{panic(diagnose("compiler","cyclic import: "+doc.URI,doc.URI,1,1))}
  if states[doc.URI]==2{return}
  states[doc.URI]=1
  total+=len(doc.Text);if total>editorSourceLimit{panic(diagnose("compiler","project source too large",doc.URI,1,1))}
  imports,local,ss:=parseNamed(doc.Text,doc.URI)
  if !entry&&len(ss)!=0{panic(diagnose("compiler","imported module may contain functions only",doc.URI,1,1))}
  for _,name:=range imports{
   if strings.HasPrefix(name,"@std:"){libs.add(strings.TrimPrefix(name,"@std:"));continue}
   if strings.HasPrefix(name,"@pkg:"){if warning==""{warning="Vault import not resolved in open-buffer mode"};continue}
   child,problem:=resolveEditorImport(root.URI,doc.URI,name,docs)
   if problem!=""{if warning==""{warning=problem};continue}
   visit(child,false)
  }
  fs=append(fs,local...);if entry{main=ss}
  states[doc.URI]=2
 }
 visit(root,true)
 return append(libs.functions,fs...),main,warning
}
func editorDefinitions(doc lspDocument,docs map[string]lspDocument,word string)[]map[string]any{
 result:=[]map[string]any{}
 for _,d:=range editorVisibleDocs(doc,docs){
  for _,sym:=range editorSymbols(d.Text){if sym["name"]==word{result=append(result,map[string]any{"uri":d.URI,"range":sym["range"]})}}
 }
 return result
}
