package core

import (
 "fmt"
 "strings"
)

// Straight-line typed SSA analysis. Every definition has a unique value ID.
// Checked operations are retained, including those whose results are unused.
// CFG joins/phi insertion remain unimplemented; eligible blocks use linear-scan allocation.
type ssaValue struct{ op string; args []int; literal int64; name string; uses int }
type ssaProgram struct{ values []ssaValue; result int; reused, propagated int }
func buildScalarSSA(f fn)(ssaProgram,bool){
 p:=ssaProgram{result:-1}
 if f.record||f.ret!="int"{return p,false}
 bindings:=map[string]int{}
 // CSE is restricted to a single dominating block of pure integer
 // operations. A repeated checked operation cannot newly trap after the
 // identical earlier operation succeeded.
 seen:=map[string]int{}
 add:=func(v ssaValue)int{
  if v.op!="param.i64"{
   if len(v.args)>0{
    args:=[]*expr{};all:=true
    for _,id:=range v.args{prior:=p.values[id];if prior.op!="const.i64"{all=false;break};args=append(args,&expr{op:"int",n:prior.literal})}
    if all{
     op:=map[string]string{"neg.checked.i64":"-","add.checked.i64":"+","sub.checked.i64":"-","mul.checked.i64":"*","div.checked.i64":"/","rem.checked.i64":"%"}[v.op]
     kind:="bin";if len(args)==1{kind="un"}
     if k,n,ok:=fold(&expr{op:kind,s:op,args:args});ok&&k=="int"{v=ssaValue{op:"const.i64",literal:n};p.propagated++}
    }
   }
   key:=fmt.Sprintf("%s:%d:%v",v.op,v.literal,v.args)
   if i,ok:=seen[key];ok{p.reused++;return i}
   seen[key]=len(p.values)
  }
  for _,a:=range v.args{p.values[a].uses++}
  p.values=append(p.values,v);return len(p.values)-1
 }
 for _,a:=range f.args{if a.typ!="int"{return p,false};bindings[a.name]=add(ssaValue{op:"param.i64",name:a.name})}
 var lower func(*expr)int
 lower=func(e *expr)int{
  if e==nil{return -1}
  if k,n,ok:=fold(e);ok&&k=="int"{return add(ssaValue{op:"const.i64",literal:n})}
  if e.op=="var"{if v,ok:=bindings[e.s];ok{return v};return -1}
  if e.op=="un"&&e.s=="-"{x:=lower(e.args[0]);if x<0{return -1};return add(ssaValue{op:"neg.checked.i64",args:[]int{x}})}
  if e.op!="bin"{return -1}
  op,ok:=map[string]string{"+":"add","-":"sub","*":"mul","/":"div","%":"rem"}[e.s];if !ok{return -1}
  a,b:=lower(e.args[0]),lower(e.args[1]);if a<0||b<0{return -1}
  return add(ssaValue{op:op+".checked.i64",args:[]int{a,b}})
 }
 for i,s:=range f.body{
  switch s.op{
  case "decl","assign":
   if s.typ!=""&&s.typ!="int"{return p,false}
   x:=lower(s.e);if x<0{return p,false};bindings[s.name]=x
  case "return":
   if i!=len(f.body)-1{return p,false}
   x:=lower(s.e);if x<0{return p,false};p.result=x;p.values[x].uses++
  default:return p,false
  }
 }
 return p,p.result>=0
}
func printSSA(f fn){
 p,ok:=buildScalarSSA(f);if !ok{return}
 fmt.Printf("ssa %s (single block, inspection):\n",f.name)
 for i,v:=range p.values{
  args:=[]string{}
  for _,a:=range v.args{args=append(args,fmt.Sprintf("v%d",a))}
  if v.op=="const.i64"{args=append(args,fmt.Sprint(v.literal))}
  if v.name!=""{args=append(args,v.name)}
  fmt.Printf("  v%d = %s %s ; uses=%d\n",i,v.op,strings.Join(args,", "),v.uses)
 }
 fmt.Printf("  return v%d\n",p.result)
 fmt.Println("  "+ssaLocationReport(p))
 fmt.Printf("  redundant values reused=%d; constants propagated=%d\n",p.reused,p.propagated)
}
