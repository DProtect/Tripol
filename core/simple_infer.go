package core

import "strings"

// Monomorphic, bounded inference for single-return functions in one source
// unit. Explicit annotations and defaults win. Ambiguous parameters retain
// the legacy int default. The ordinary compiler independently checks every
// call and body, so an inference guess cannot bypass type safety.
func inferSimpleFunctions(fs []fn,main []stmt){
 byName:=map[string]int{}
 eligible:=make([]bool,len(fs))
 for i,f:=range fs{
  byName[f.name]=i
  if !f.record&&f.inferRet&&len(f.body)==1&&f.body[0].op=="return" {
   eligible[i]=true
   for j:=range f.args{if j<len(f.implicit)&&f.implicit[j]{fs[i].args[j].typ=""}}
   fs[i].ret=""
  }
 }
 changed:=false
 scalar:=func(t string)bool{return listType(t)!=""||isList(t)||t=="ints"}
 var guess func(*expr,map[string]string,int)string
 constrain:=func(e *expr,t string,env map[string]string,owner int){
  if owner<0||e==nil||e.op!="var"||!scalar(t)||env[e.s]!=""{return}
  for j,p:=range fs[owner].args{if p.name==e.s&&p.typ==""{fs[owner].args[j].typ=t;env[e.s]=t;changed=true;break}}
 }
 guess=func(e *expr,env map[string]string,owner int)string{
  if e==nil{return ""}
  switch e.op{
  case "int","bool","float","string":return e.op
  case "var":return env[e.s]
  case "format":
   for _,x:=range e.args{guess(x,env,owner)}
   return "string"
  case "index":
   t:=env[e.s];if isList(t){return listElement(t)};if t=="ints"||strings.HasPrefix(t,"array:"){return "int"}
  case "un":
   t:=guess(e.args[0],env,owner)
   if e.s=="!"{constrain(e.args[0],"bool",env,owner);return "bool"}
   if t=="float"{return t}
   constrain(e.args[0],"int",env,owner);return "int"
  case "bin":
   a,b:=guess(e.args[0],env,owner),guess(e.args[1],env,owner)
   switch e.s{
   case "&&","||":constrain(e.args[0],"bool",env,owner);constrain(e.args[1],"bool",env,owner);return "bool"
   case "==","!=","<","<=",">",">=":
    if a!=""{constrain(e.args[1],a,env,owner)};if b!=""{constrain(e.args[0],b,env,owner)};return "bool"
   case "+","-","*","/","%":
    if a!=""{constrain(e.args[1],a,env,owner);return a}
    if b!=""{constrain(e.args[0],b,env,owner);return b}
   }
  case "array":return "array:literal"
  case "call":
   if i,ok:=byName[e.s];ok{
    for j,arg:=range e.args{
     t:=guess(arg,env,owner)
     if j>=len(fs[i].args){continue}
     want:=fs[i].args[j].typ
     if eligible[i]&&want==""&&scalar(t){fs[i].args[j].typ=t;changed=true}
     if want!=""{constrain(arg,want,env,owner)}
    }
    return fs[i].ret
   }
   if e.s=="list"&&len(e.args)==1&&e.args[0].op=="array"&&len(e.args[0].args)>0{return listType(guess(e.args[0].args[0],env,owner))}
   if isList(e.s){return e.s}
   if name,method,ok:=strings.Cut(e.s,".");ok&&isList(env[name]){
    for _,x:=range e.args{guess(x,env,owner)}
    switch method{case "get","remove_at","set":return listElement(env[name]);case "add","size","capacity","clear","sort":return "int";case "contains","remove":return "bool"}
   }
   for _,x:=range e.args{guess(x,env,owner)}
   switch e.s{case "input","text","slice":return "string";case "len","abs","min","max","input_int","parse_int","sum","@array_len":return "int";case "float","sqrt":return "float";case "ints":return "ints"}
  }
  return ""
 }
 copyEnv:=func(env map[string]string)map[string]string{out:=map[string]string{};for k,v:=range env{out[k]=v};return out}
 var walk func([]stmt,map[string]string)
 walk=func(body []stmt,env map[string]string){
  for _,s:=range body{
   t:=guess(s.e,env,-1);guess(s.bound,env,-1)
   if s.op=="decl"||s.op=="assign"{if s.typ!=""{t=s.typ};if t!=""{env[s.name]=t}}
   if s.op=="showmany"{for _,x:=range s.e.args{guess(x,env,-1)}}
   if s.op=="for"{local:=copyEnv(env);local[s.name]="int";walk(s.yes,local)}else{walk(s.yes,copyEnv(env))}
   walk(s.no,copyEnv(env));walk(s.final,copyEnv(env))
  }
 }
 for pass:=0;pass<16;pass++{
  changed=false
  walk(main,map[string]string{})
  for i,f:=range fs{
   if !eligible[i]{continue}
   env:=map[string]string{};for _,p:=range f.args{env[p.name]=p.typ}
   t:=guess(f.body[0].e,env,i)
   if f.ret==""&&scalar(t){fs[i].ret=t;changed=true}
  }
  if !changed{break}
 }
 for i:=range fs{if eligible[i]{for j:=range fs[i].args{if fs[i].args[j].typ==""{fs[i].args[j].typ="int"}};if fs[i].ret==""{fs[i].ret="int"}}}
}
