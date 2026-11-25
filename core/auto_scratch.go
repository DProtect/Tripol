package core

import "strings"

// Conservative region proof: scalar values, immutable strings and locally
// owned integer arrays/lists. No calls to arbitrary user functions, reference
// stores into outer bindings, control transfer, IO or opaque operations.
func (c *compiler) localLoopRegion(ss []stmt)bool {
 if len(c.scratchScopes)>0||len(ss)==0{return false}
 local:=map[string]string{}
 allocates:=false
 get:=func(name string)(string,bool){
  if t,ok:=local[name];ok{return t,true}
  if v,ok:=c.lookup(name);ok&&(v.typ=="int"||v.typ=="bool"||v.typ=="string"){return v.typ,true}
  return "",false
 }
 var exprType func(*expr)string
 exprType=func(e *expr)string{
  if e==nil{return ""}
  switch e.op{
  case "int","bool","string":return e.op
  case "var":t,_:=get(e.s);return t
  case "index":
   if (local[e.s]=="ints"||local[e.s]=="list")&&exprType(e.args[0])=="int"{return "int"}
  case "format":
   for _,arg:=range e.args{t:=exprType(arg);if t!="int"&&t!="bool"&&t!="string"{return ""}}
   allocates=true;return "string"
  case "un":
   t:=exprType(e.args[0]);if (e.s=="-"&&t=="int")||(e.s=="!"&&t=="bool"){return t}
  case "bin":
   a,z:=exprType(e.args[0]),exprType(e.args[1]);if a==""||a!=z{return ""}
   switch e.s{
   case "+":if a=="string"{allocates=true;return "string"};if a=="int"{return "int"}
   case "-","*","/","%":if a=="int"{return "int"}
   case "==","!=":if a=="int"||a=="bool"||a=="string"{return "bool"}
   case "<","<=",">",">=":if a=="int"{return "bool"}
   case "&&","||":if a=="bool"{return "bool"}
   }
  case "call":
   if _,exists:=c.fs[e.s];exists{return ""}
   if name,method,ok:=strings.Cut(e.s,".");ok&&local[name]=="list"{
    count,known:=map[string]int{"add":1,"remove_at":1,"get":1,"set":2,"size":0,"capacity":0,"clear":0}[method]
    if !known||len(e.args)!=count{return ""}
    for _,arg:=range e.args{if exprType(arg)!="int"{return ""}}
    if method=="add"{allocates=true}
    return "int"
   }
   switch e.s{
   case "ints":
    if len(e.args)==1&&exprType(e.args[0])=="int"{allocates=true;return "ints"}
   case "list":
    if len(e.args)==0||(len(e.args)==1&&exprType(e.args[0])=="int"){allocates=true;return "list"}
   case "len":
    if len(e.args)==1{t:=exprType(e.args[0]);if t=="ints"||t=="list"||t=="string"{return "int"}}
   case "text":
    if len(e.args)==1&&exprType(e.args[0])=="int"{allocates=true;return "string"}
   case "slice":
    if len(e.args)==3&&exprType(e.args[0])=="string"&&exprType(e.args[1])=="int"&&exprType(e.args[2])=="int"{allocates=true;return "string"}
   case "byte_at":
    if len(e.args)==2&&exprType(e.args[0])=="string"&&exprType(e.args[1])=="int"{return "int"}
   case "abs","min","max":
    n:=2;if e.s=="abs"{n=1};if len(e.args)!=n{return ""}
    for _,arg:=range e.args{if exprType(arg)!="int"{return ""}};return "int"
   }
  }
  return ""
 }
 for _,s:=range ss{
  switch s.op{
  case "decl","assign":
   t:=exprType(s.e);if t==""{return false}
   if old,ok:=get(s.name);ok{
    if old!=t||s.op=="decl"{return false}
    // Outer strings are readable but cannot be reassigned within a region.
    if _,inside:=local[s.name];!inside&&t=="string"{return false}
   }else{
    if _,exists:=c.lookup(s.name);exists{return false}
    local[s.name]=t
   }
  case "indexassign":
   if (local[s.name]!="ints"&&local[s.name]!="list")||exprType(s.bound)!="int"||exprType(s.e)!="int"{return false}
  case "assert":
   if exprType(s.e)!="bool"{return false}
  case "show","expr":
   if exprType(s.e)==""{return false}
  case "showmany":
   for _,e:=range s.e.args{if exprType(e)==""{return false}}
  default:return false
  }
 }
 return allocates
}
func (c *compiler) loopBody(ss []stmt){
 if c.localLoopRegion(ss){
  c.a.guard=true
  c.scratchBlock(stmt{op:"scratch",yes:ss})
  return
 }
 c.block(ss,true)
}
