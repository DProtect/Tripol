package core

import "encoding/binary"

// RAX is the old element; RCX is the already evaluated RHS. The index is
// preserved beneath both operands, so it is evaluated exactly once.
func (c *compiler) compoundInteger(op string) {
 a:=c.a
 switch op {
 case "+": a.emit(0x48,0x01,0xc8);a.j(0x80,"err")
 case "-": a.emit(0x48,0x29,0xc8);a.j(0x80,"err")
 case "*": a.emit(0x48,0x0f,0xaf,0xc1);a.j(0x80,"err")
 case "/","%":
  a.emit(0x48,0x85,0xc9);a.j(0x84,"err")
  ok:=a.fresh();a.emit(0x48,0x83,0xf9,0xff);a.j(0x85,ok)
  a.emit(0x48,0xba);a.b=binary.LittleEndian.AppendUint64(a.b,1<<63)
  a.emit(0x48,0x39,0xd0);a.j(0x84,"err");a.label(ok)
  a.emit(0x48,0x99,0x48,0xf7,0xf9)
  if op=="%" {a.emit(0x48,0x89,0xd0)}
 default:c.err("unsupported compound array operation")
 }
}
