package core

func ssaFitsImmediate(n int64)bool{return n>=-2147483648&&n<=2147483647}
func ssaImmediateOp(op string)bool{
 return op=="add.checked.i64"||op=="sub.checked.i64"||op=="mul.checked.i64"
}

// A literal used exclusively as the encodable RHS of checked arithmetic
// need not occupy a value register or spill slot. Other uses retain ordinary
// allocation, including division and return values.
func ssaImmediateOnly(p ssaProgram)[]bool{
 eligible:=make([]bool,len(p.values))
 for i,v:=range p.values{eligible[i]=v.op=="const.i64"&&ssaFitsImmediate(v.literal)}
 for _,v:=range p.values{
  for j,id:=range v.args{if j!=1||!ssaImmediateOp(v.op){eligible[id]=false}}
 }
 if p.result>=0{eligible[p.result]=false}
 return eligible
}
func(a *asm) ssaArithmeticImmediate(op string,n int64)bool{
 if !ssaFitsImmediate(n){return false}
 short:=n>=-128&&n<=127
 switch op{
 case "+","-":
  mod:=byte(0xc0);if op=="-"{mod=0xe8}
  if short{a.emit(0x48,0x83,mod,byte(n))}else{a.emit(0x48,0x81,mod);a.u32(int(n))}
 case "*":
  if short{a.emit(0x48,0x6b,0xc0,byte(n))}else{a.emit(0x48,0x69,0xc0);a.u32(int(n))}
 default:return false
 }
 a.j(0x80,"err")
 return true
}
