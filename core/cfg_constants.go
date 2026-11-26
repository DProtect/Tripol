package core

import (
 "os"
 "runtime"
 "strconv"
 "sync"
)

// Forward scalar facts over structured, acyclic control flow. Only int/bool
// values are tracked. Calls cannot mutate scalar caller slots; nevertheless
// call arguments and aggregate access are left unchanged (some intrinsics
// require a variable expression). Loops retain only invariant scalar facts;
// protected scopes are barriers.
// This runs only after full validation and never mutates the parsed AST.
type scalarFact struct{kind string; value int64}
type scalarFacts map[string]scalarFact
type cfgConstantStats struct{ Substitutions,Joins,ConstantBranches,LoopInvariants int }

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
// Recursively collect scalar assignments. A protected region or opaque
// statement disables loop propagation, rather than assuming purity.
func loopWrites(body []stmt,writes map[string]bool)bool{
 for _,s:=range body{
  switch s.op{
  case "assign","decl":writes[s.name]=true
  case "for":writes[s.name]=true
  case "if","while","show","showmany","return","assert","expr","break","continue","indexassign","fieldassign":
  default:return false
  }
  if !loopWrites(s.yes,writes)||!loopWrites(s.no,writes)||!loopWrites(s.final,writes){return false}
 }
 return true
}
// Refine facts on the edge where a validated scalar comparison proves
// equality. Facts are local to the branch; existing join rules govern exit.
func refineCondition(e *expr,yes,no scalarFacts){
 if e==nil{return}
 if e.op=="var"{yes[e.s]=scalarFact{"bool",1};no[e.s]=scalarFact{"bool",0};return}
 if e.op!="bin"||len(e.args)!=2||(e.s!="=="&&e.s!="!="){return}
 a,b:=e.args[0],e.args[1]
 if a.op!="var"{a,b=b,a}
 if a.op!="var"||(b.op!="int"&&b.op!="bool"){return}
 if e.s=="=="{yes[a.s]=scalarFact{b.op,b.n}}else{no[a.s]=scalarFact{b.op,b.n}}
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
   refineCondition(s.e,yes,no)
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
  case "while","for":
   // Range bounds are evaluated once. While conditions are reevaluated;
   // substitute only variables untouched anywhere in the loop.
   if s.op=="for"{s.e=scalarSubstitute(s.e,env,stats);s.bound=scalarSubstitute(s.bound,env,stats)}
   writes:=map[string]bool{}
   if s.op=="for"{writes[s.name]=true}
   if !loopWrites(s.yes,writes){clear(env);break}
   for name:=range env{if writes[name]{delete(env,name)}else{stats.LoopInvariants++}}
   if s.op=="while"{s.e=scalarSubstitute(s.e,env,stats)}
   // The body starts with invariant facts, not values from the previous
   // iteration. Body-local facts cannot escape the backedge.
   s.yes=propagateBlock(s.yes,copyFacts(env),stats)
  case "scratch","try":
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
  // All statements have already been validated. A direct terminator
  // makes the remaining statements in this lexical block unreachable.
  // Do not move or discard the terminator's own checked expression.
  if s.op=="return"||s.op=="break"||s.op=="continue"{return result[:i+1]}
 }
 return result
}
// Bounded worker pool: independently optimize function bodies. Workers do
// not emit code, modify the input AST, share mutable facts or publish output.
// The original function order is retained for deterministic native images.
func cfgWorkerCount(n int)int{
 jobs:=runtime.GOMAXPROCS(0)
 if jobs>4{jobs=4}
 if value:=os.Getenv("TRIPOL_JOBS");value!=""{
  if requested,e:=strconv.Atoi(value);e==nil&&requested>=1&&requested<=4{jobs=requested}
 }
 if n<64||jobs<2{return 1}
 return jobs
}
func optimizeCFG(fs []fn,ss []stmt)([]fn,[]stmt,cfgConstantStats){
 return optimizeCFGWorkers(fs,ss,cfgWorkerCount(len(fs)))
}
func optimizeCFGWorkers(fs []fn,ss []stmt,workers int)([]fn,[]stmt,cfgConstantStats){
 stats:=cfgConstantStats{}
 if diagnosticTrace{return fs,ss,stats}
 out:=append([]fn(nil),fs...)
 if workers<1{workers=1};if workers>4{workers=4}
 if len(fs)<workers{workers=len(fs)}
 local:=make([]cfgConstantStats,len(fs))
 work:=func(i int){
  f:=fs[i]
  if !f.record{out[i].body=propagateBlock(f.body,scalarFacts{},&local[i])}
 }
 if workers<=1{
  for i:=range fs{work(i)}
 }else{
  var wg sync.WaitGroup
  // Static disjoint ranges avoid task-channel allocation and shared writes.
  for worker:=0;worker<workers;worker++{
   start,end:=len(fs)*worker/workers,len(fs)*(worker+1)/workers
   wg.Add(1)
   go func(start,end int){defer wg.Done();for i:=start;i<end;i++{work(i)}}(start,end)
  }
  wg.Wait()
 }
 for _,r:=range local{
  stats.Substitutions+=r.Substitutions
  stats.Joins+=r.Joins
  stats.ConstantBranches+=r.ConstantBranches
  stats.LoopInvariants+=r.LoopInvariants
 }
 return out,propagateBlock(ss,scalarFacts{},&stats),stats
}
