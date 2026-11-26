package core

import (
 "strings"
 "testing"
)
func TestLSPOpenBufferImports(t *testing.T){
 root:=lspDocument{URI:"file:///project/main.tri",Text:"import \"lib/math.tri\"\nshow twice(3)",Version:1}
 lib:=lspDocument{URI:"file:///project/lib/math.tri",Text:"func twice(x:int)->int {return x*2}",Version:1}
 docs:=map[string]lspDocument{root.URI:root}
 missing:=validateEditorProject(root,docs)
 if len(missing)!=1||missing[0]["severity"]!=2{t.Fatal(missing)}
 docs[lib.URI]=lib
 if d:=validateEditorProject(root,docs);len(d)!=0{t.Fatal(d)}
 defs:=editorDefinitions(root,docs,"twice")
 if len(defs)!=1||defs[0]["uri"]!=lib.URI{t.Fatal(defs)}
 // Unsaved imported edits must change the root program's diagnostics.
 lib.Text="func twice(x:string)->string {return x}"
 docs[lib.URI]=lib
 d:=validateEditorProject(root,docs);if len(d)!=1||d[0]["severity"]!=1{t.Fatal(d)}
 lib.Text="import \"../main.tri\"\nfunc twice(x:int)->int {return x}"
 docs[lib.URI]=lib
 d=validateEditorProject(root,docs);if len(d)!=1||!strings.Contains(d[0]["message"].(string),"cyclic import"){t.Fatal(d)}
 outside:=lspDocument{URI:"file:///outside.tri",Text:"func secret()->int {return 1}"}
 docs[outside.URI]=outside
 if _,warning:=resolveEditorImport(root.URI,root.URI,"../outside.tri",docs);warning==""{t.Fatal("resolved outside entry directory")}
 if _,warning:=resolveEditorImport(root.URI,root.URI,"https://example.com/math.tri",docs);warning==""{t.Fatal("resolved external URL")}
}
