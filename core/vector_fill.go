package core

// TIR VECTOR R1: vectorize a zero-based constant fill of a mutable dynamic
// int64 array. SSE2 is baseline on x86-64. No checked arithmetic is reordered;
// out-of-bounds fills write the same valid prefix before the original trap.
// Protected scopes and register-promoted loops use the scalar fallback.
func (c *compiler) vectorFill(s stmt) bool {
 a:=c.a
 if !a.hyper||a.guard||a.trace||a.loopRegs!=nil||s.e==nil||s.e.op!="int"||s.e.n!=0||len(s.yes)!=1{return false}
 st:=s.yes[0]
 if st.op!="indexassign"||st.typ!=""||st.bound==nil||st.bound.op!="var"||st.bound.s!=s.name{return false}
 k,value,ok:=fold(st.e)
 if !ok||k!="int"{return false}
 arr,ok:=c.lookup(st.name)
 if !ok||arr.typ!="ints"||!arr.mut{return false}
 c.require(c.expression(s.bound),"int")
 done,limit,loop,tail,finish:=a.fresh(),a.fresh(),a.fresh(),a.fresh(),a.fresh()
 a.emit(0x48,0x85,0xc0);a.j(0x8e,done)
 a.regMove(1,0) // rcx=count
 a.load(arr.off)
 a.emit(0x48,0x8b,0x10,0x45,0x31,0xc0) // rdx=length, r8d=0
 a.emit(0x48,0x39,0xd1);a.j(0x86,limit)
 a.regMove(1,2);a.emit(0x41,0xb8,1,0,0,0)
 a.label(limit)
 a.emit(0x48,0x8d,0x50,8) // rdx=data
 a.imm(value)
 if a.turbo {
  a.call("vector_fill")
 } else {
  a.emit(0x66,0x48,0x0f,0x6e,0xc0) // movq xmm0,rax
  a.emit(0x66,0x0f,0x6c,0xc0) // punpcklqdq xmm0,xmm0
  a.label(loop)
  a.emit(0x48,0x83,0xf9,2);a.j(0x82,tail)
  a.emit(0xf3,0x0f,0x7f,0x02) // movdqu [rdx],xmm0
  a.emit(0x48,0x83,0xc2,16,0x48,0x83,0xe9,2);a.jmp(loop)
  a.label(tail)
  a.emit(0x48,0x85,0xc9);a.j(0x84,finish)
  a.emit(0x48,0x89,0x02)
 }
 a.label(finish)
 a.emit(0x45,0x85,0xc0);a.j(0x85,"bounds")
 a.label(done)
 return true
}
