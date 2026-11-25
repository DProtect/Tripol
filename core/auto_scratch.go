package core

// A deliberately narrow escape proof for loop-local integer aggregates.
// No user functions, reference stores into outer bindings, control transfer,
// IO, strings or opaque operations are eligible. Validation remains independent.
func (c *compiler) localLoopRegion(ss []stmt)bool {
 if len(c.scratchScopes)>0||len(ss)==0{return false}
 local:=map[string]string{}
 allocates:=false
 get:=func(name string)(string,bool){
  if t,ok:=local[name];ok{return t,true}
  if v,ok:=c.lookup(name);ok&&(v.typ=="int"||v.typ=="bool"){return v.typ,true}
  return "",false
 }
 var exprType func(*expr)string
 exprType=func(e *expr)string{
  if e==nil{return ""}
  switch e.op{
  case "int","bool":return e.op
  case "var":t,_:=get(e.s);return t
  case "index":
   if local[e.s]=="ints"&&exprType(e.args[0])=="int"{return "int"}
  case "un":
   t:=exprType(e.args[0]);if (e.s=="-"&&t=="int")||(e.s=="!"&&t=="bool"){return t}
  case "bin":
   a,z:=exprType(e.args[0]),exprType(e.args[1]);if a==""||a!=z{return ""}
   switch e.s{case "+","-","*","/","%":if a=="int"{return "int"};case "==","!=":if a=="int"||a=="bool"{return "bool"};case "<","<=",">",">=":if a=="int"{return "bool"};case "&&","||":if a=="bool"{return "bool"}}
  case "call":
   if _,exists:=c.fs[e.s];exists{return ""}
   switch e.s{
   case "ints":
    if len(e.args)==1&&exprType(e.args[0])=="int"{allocates=true;return "ints"}
   case "len":
    if len(e.args)==1&&exprType(e.args[0])=="ints"{return "int"}
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
   if old,ok:=get(s.name);ok{if old!=t{return false};if s.op=="decl"{return false}
   }else{
    if _,exists:=c.lookup(s.name);exists{return false}
    local[s.name]=t
   }
  case "indexassign":
   if local[s.name]!="ints"||exprType(s.bound)!="int"||exprType(s.e)!="int"{return false}
  case "assert":
   if exprType(s.e)!="bool"{return false}
  case "show":
   if exprType(s.e)==""{return false}
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
