package core

import "fmt"

// Experimental statement-level CFG. Each node is a basic block containing one
// operation. It does not rewrite code or claim SSA/register allocation.
type flowBlock struct { op string; next []int }
type flowGraph struct { blocks []flowBlock; entry,exit int }
type flowLoop struct { end,again int }
func inspectFlow(ss []stmt)(flowGraph,bool){
 g:=flowGraph{}
 add:=func(op string,edges ...int)int{g.blocks=append(g.blocks,flowBlock{op,edges});return len(g.blocks)-1}
 g.exit=add("exit")
 valid:=true
 var lower func([]stmt,int,*flowLoop)int
 lower=func(stmts []stmt,next int,loop *flowLoop)int{
  for i:=len(stmts)-1;i>=0;i--{
   s:=stmts[i]
   switch s.op{
   case "if":
    yes:=lower(s.yes,next,loop);no:=lower(s.no,next,loop)
    next=add("if",yes,no)
   case "while","for":
    head:=add(s.op+".condition")
    step:=head
    if s.op=="for"{step=add("for.step",head)}
    body:=lower(s.yes,step,&flowLoop{end:next,again:step})
    g.blocks[head].next=[]int{body,next}
    next=head
    if s.op=="for"{next=add("for.init",head)}
   case "return":next=add("return",g.exit)
   case "break","continue":
    if loop==nil{valid=false;return next}
    dest:=loop.end;if s.op=="continue"{dest=loop.again}
    next=add(s.op,dest)
   case "try","scratch","throw":valid=false;return next
   default:next=add(s.op,next)
   }
  }
  return next
 }
 g.entry=lower(ss,g.exit,nil)
 return g,valid
}
func printFlow(name string,ss []stmt){
 g,ok:=inspectFlow(ss);if !ok{fmt.Printf("cfg %s: skipped (protected or unsupported control flow)\n",name);return}
 fmt.Printf("cfg %s entry=b%d exit=b%d\n",name,g.entry,g.exit)
 for i,b:=range g.blocks{fmt.Printf("  b%d %s ->",i,b.op);for _,n:=range b.next{fmt.Printf(" b%d",n)};fmt.Println()}
}
