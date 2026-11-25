package core

import "fmt"

// Linear scan on single-block scalar SSA. RAX is reserved for computation;
// six volatile registers hold values across operations. Excess live values
// receive frame slots. No calls, joins or unsafe speculative evaluation.
type ssaLocation struct { reg byte; spill int; present bool }
type ssaAllocation struct { locations []ssaLocation; frame int; spills int }
func allocateSSA(p ssaProgram) ssaAllocation {
 a:=ssaAllocation{locations:make([]ssaLocation,len(p.values))}
 last:=make([]int,len(p.values));for i:=range last{last[i]=-1}
 for i,v:=range p.values{for _,operand:=range v.args{last[operand]=i}}
 last[p.result]=len(p.values)
 regs:=[]byte{8,9,10,11,1,2}
 owners:=map[byte]int{}
 for i:=range p.values {
  for r,owner:=range owners {if last[owner]<=i{delete(owners,r)}}
  if last[i]<0{continue}
  found:=false
  for _,r:=range regs{if _,ok:=owners[r];!ok{owners[r]=i;a.locations[i]=ssaLocation{reg:r,present:true};found=true;break}}
  if !found{a.frame+=8;a.spills++;a.locations[i]=ssaLocation{spill:-a.frame,present:true}}
 }
 return a
}
func (c *compiler) emitNativeSSA(f fn)bool{
 a:=c.a
 if !a.hyper||a.trace||len(f.body)<2{return false}
 p,ok:=buildScalarSSA(f);if !ok||len(p.values)>1024{return false}
 for _,v:=range p.values {
  switch v.op {case "param.i64","const.i64","neg.checked.i64","add.checked.i64","sub.checked.i64","mul.checked.i64":default:return false}
 }
 plan:=allocateSSA(p)
 load:=func(i int){
  l:=plan.locations[i]
  if !l.present{die("internal SSA missing operand")}
  if l.spill!=0{a.load(l.spill)}else{a.regMove(0,l.reg)}
 }
 save:=func(i int){
  l:=plan.locations[i];if !l.present{return}
  if l.spill!=0{a.store(l.spill)}else{a.regMove(l.reg,0)}
 }
 a.emit(0x55,0x48,0x89,0xe5)
 if plan.frame>0{a.emit(0x48,0x81,0xec);a.u32(plan.frame)}
 parameter:=0
 for i,v:=range p.values{
  switch v.op{
  case "param.i64":
   if plan.locations[i].present{a.load(16+8*(len(f.args)-1-parameter));save(i)}
   parameter++
  case "const.i64":
   if plan.locations[i].present{a.imm(v.literal);save(i)}
  case "neg.checked.i64":
   load(v.args[0]);a.emit(0x48,0xf7,0xd8);a.j(0x80,"err");save(i)
  default:
   // Even a dead checked result is evaluated; error order is unchanged.
   load(v.args[0]);op:=map[string]string{"add.checked.i64":"+","sub.checked.i64":"-","mul.checked.i64":"*"}[v.op]
   right:=plan.locations[v.args[1]]
   if right.spill!=0{a.integerMemoryOp(op,right.spill)}else{a.regArith(op,0,right.reg)}
   save(i)
  }
 }
 load(p.result);a.emit(0xc9,0xc3)
 return true
}
func ssaLocationReport(p ssaProgram)string{
 a:=allocateSSA(p)
 return fmt.Sprintf("linear scan: 6 value registers; %d spill slots",a.spills)
}
