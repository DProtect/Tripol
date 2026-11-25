package core

// Aligned int64 operations on arena arrays. All stores/RMWs use locked x86
// instructions; no user thread API or allocator thread-safety is implied.
func (c *compiler) atomicBuiltin(e *expr)(string,bool){
 counts:=map[string]int{"atomic_load":2,"atomic_exchange":3,"atomic_compare_exchange":4,"atomic_add":3}
 count,ok:=counts[e.s];if !ok{return "",false}
 if len(e.args)!=count{c.err("invalid atomic argument count")}
 a:=c.a
 for i,arg:=range e.args{
  typ:="int";if i==0{typ="ints"}
  c.require(c.expression(arg),typ);a.emit(0x50)
 }
 // Save new/delta in R8, expected in RAX, index in RCX, pointer in RDX.
 if count==4{a.emit(0x41,0x58,0x58,0x59,0x5a)
 }else if count==3{a.emit(0x41,0x58,0x59,0x5a)
 }else{a.emit(0x59,0x5a)}
 a.emit(0x48,0x3b,0x0a);a.j(0x83,"bounds")
 switch e.s{
 case "atomic_load":a.emit(0x48,0x8b,0x44,0xca,0x08)
 case "atomic_exchange":
  a.emit(0x4c,0x89,0xc0,0x48,0x87,0x44,0xca,0x08)
 case "atomic_compare_exchange":
  a.emit(0xf0,0x4c,0x0f,0xb1,0x44,0xca,0x08,0x0f,0x94,0xc0,0x0f,0xb6,0xc0)
  return "bool",true
 case "atomic_add":
  a.emit(0x48,0x8b,0x44,0xca,0x08)
  retry:=a.fresh();a.label(retry)
  a.emit(0x49,0x89,0xc1,0x4d,0x01,0xc1);a.j(0x80,"err")
  a.emit(0xf0,0x4c,0x0f,0xb1,0x4c,0xca,0x08);a.j(0x85,retry)
  a.emit(0x4c,0x89,0xc8)
 }
 return "int",true
}
