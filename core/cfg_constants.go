package core

// Forward scalar facts over structured, acyclic control flow. Only int/bool
// values are tracked. Calls cannot mutate scalar caller slots; nevertheless
// call arguments and aggregate access are left unchanged (some intrinsics
// require a variable expression). Loops and protected scopes are barriers.
// This runs only after full validation and never mutates the parsed AST.
type scalarFact struct{kind string; value int64}
type scalarFacts map[string]scalarFact
type cfgConstantStats struct{ Substitutions,Joins,ConstantBranches int }

func copyFacts(in scalarFacts)scalarFacts{
 out:=make(scalarFacts,len(in));for k,v:=range in{out[k]=v};return out
}
func scalarSubstitute(e *expr,env scalarFacts,stats *cfgConstantStats)*expr{
 if e==nil{return nil}
 if e.op=="var"{
  if v,ok:=env[e.s];ok{stats.Substitutions++;return &expr{op:v.kind,n:v.value}}
  return e
 }
 if e.op!="un"&&e.op!="bin"&&e.op!="format"{return e}
 var args []*expr
 for i,x:=range e.args{
  y:=scalarSubstitute(x,env,stats)
  if y!=x{if args==nil{args=append([]*expr(nil),e.args...)};args[i]=y}
 }
 if args==nil{return e}
 result:=*e;result.args=args;return &result
}
func propagateBlock(body []stmt,env scalarFacts,stats *cfgConstantStats)[]stmt{
 if len(body)==0{return body}
 var result []stmt
 for i,original:=range body{
  s:=original
  switch s.op{
  case "assign","decl":
   s.e=scalarSubstitute(s.e,env,stats)
   if k,n,ok:=fold(s.e);ok&&(k=="int"||k=="bool"){env[s.name]=scalarFact{k,n}}else{delete(env,s.name)}
  case "if":
   s.e=scalarSubstitute(s.e,env,stats)
   yes,no:=copyFacts(env),copyFacts(env)
   s.yes=propagateBlock(s.yes,yes,stats)
   s.no=propagateBlock(s.no,no,stats)
   k,n,constant:=fold(s.e)
   if constant&&k=="bool"{stats.ConstantBranches++}
   // Only incoming known bindings can cross the join. Branch-local
   // declarations cannot escape their lexical scope.
   for name:=range env{
    a,aok:=yes[name];b,bok:=no[name]
    if constant&&k=="bool"{
     if n==0{a,aok=b,bok}
     if aok{env[name]=a}else{delete(env,name)}
    }else if aok&&bok&&a==b{env[name]=a;stats.Joins++}else{delete(env,name)}
   }
  case "while","for","scratch","try":
   // No speculation across a backedge, exception or arena lifetime.
   // Keep the original body and existing loop optimization patterns.
   clear(env)
  case "show","return","assert","throw","expr":
   s.e=scalarSubstitute(s.e,env,stats)
  case "showmany":
   if s.e!=nil{
    x:=*s.e;x.args=append([]*expr(nil),s.e.args...)
    for j,e:=range x.args{x.args[j]=scalarSubstitute(e,env,stats)}
    s.e=&x
   }
  case "indexassign","fieldassign":
   // Do not replace aggregate lvalues or assume anything about aliases.
   s.e=scalarSubstitute(s.e,env,stats)
  default:
   clear(env)
  }
  if result==nil{result=append([]stmt(nil),body...)}
  result[i]=s
 }
 return result
}
func optimizeCFG(fs []fn,ss []stmt)([]fn,[]stmt,cfgConstantStats){
 stats:=cfgConstantStats{}
 if diagnosticTrace{return fs,ss,stats}
 out:=append([]fn(nil),fs...)
 for i,f:=range fs{if !f.record{out[i].body=propagateBlock(f.body,scalarFacts{},&stats)}}
 return out,propagateBlock(ss,scalarFacts{},&stats),stats
}
