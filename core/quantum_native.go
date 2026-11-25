package core

// Narrow native primitives for QUANTUM X. The ordinary checked signed-int
// operators are unchanged. The u32 helpers deliberately operate on low bits.
func (c *compiler) quantumNative(e *expr)(string,bool){
 a:=c.a
 switch e.s{
 case "u32_xor","u32_and","u32_rotr","u32_shr":
  if len(e.args)!=2{c.err("expected two integer operands")}
  c.require(c.expression(e.args[0]),"int");a.emit(0x50)
  c.require(c.expression(e.args[1]),"int");a.emit(0x48,0x89,0xc1,0x58)
  switch e.s{
  case "u32_xor":a.emit(0x31,0xc8)
  case "u32_and":a.emit(0x21,0xc8)
  default:
   a.emit(0x48,0x83,0xf9,31);a.j(0x87,"bounds")
   if e.s=="u32_rotr"{a.emit(0xd3,0xc8)}else{a.emit(0xd3,0xe8)}
  }
  return "int",true
 case "process_id":
  if len(e.args)!=0{c.err("process_id expects no arguments")}
  if a.target=="linux"{a.emit(0xb8,39,0,0,0,0x0f,0x05)}else{
   a.emit(0x55,0x48,0x89,0xe5,0x48,0x83,0xe4,0xf0,0x48,0x83,0xec,0x20)
   a.rel("iat:72",0xff,0x15);a.emit(0x89,0xc0,0xc9)
  }
  return "int",true
 case "process_exit":
  if len(e.args)!=1{c.err("process_exit expects one code")}
  c.require(c.expression(e.args[0]),"int")
  a.emit(0x48,0x3d,255,0,0,0);a.j(0x87,"bounds");a.jmp("exit")
  return "int",true
 case "audio_play_file":
  if len(e.args)!=1{c.err("audio_play_file expects a WAV path")}
  c.require(c.expression(e.args[0]),"string")
  if a.target!="windows"{a.imm(0);return "bool",true}
  // PlaySoundA(path,NULL,SND_FILENAME|SND_NODEFAULT), synchronous.
  a.emit(0x48,0x89,0xc1,0x31,0xd2,0x41,0xb8,2,0,2,0)
  a.emit(0x55,0x48,0x89,0xe5,0x48,0x83,0xe4,0xf0,0x48,0x83,0xec,0x20)
  a.rel("iat:74",0xff,0x15);a.emit(0xc9,0x85,0xc0,0x0f,0x95,0xc0,0x0f,0xb6,0xc0)
  return "bool",true
 case "thread_sum":
  if len(e.args)!=1{c.err("thread_sum expects an integer array")}
  c.require(c.expression(e.args[0]),"ints")
  if a.target!="windows"{a.jmp("thread_unsupported");return "int",true}
  c.emitThreadSum()
  return "int",true
 }
 return "",false
}

// Run a read-only checked reduction on a real Windows worker and join it.
// The callback uses no allocator, no other TRIPOL code and no shared writes
// except its private result/status. The stack context stays live until join.
func (c *compiler) emitThreadSum(){
 a:=c.a
 worker,body,loop,done,overflow,joined,fail:=a.fresh(),a.fresh(),a.fresh(),a.fresh(),a.fresh(),a.fresh(),a.fresh()
 a.jmp(body)
 a.label(worker)
 // Win64 callback RCX = context; [0] array, [8] result, [16] overflow flag.
 a.emit(0x48,0x8b,0x11,0x4c,0x8b,0x02,0x48,0x83,0xc2,8,0x31,0xc0)
 a.label(loop);a.emit(0x4d,0x85,0xc0);a.j(0x84,done)
 a.emit(0x48,0x03,0x02);a.j(0x80,overflow)
 a.emit(0x48,0x83,0xc2,8,0x49,0xff,0xc8);a.jmp(loop)
 a.label(overflow);a.emit(0x48,0xc7,0x41,16,1,0,0,0)
 a.label(done);a.emit(0x48,0x89,0x41,8,0x31,0xc0,0xc3)
 a.label(body)
 // 32-byte shadow area, two stack args, context at +48, handle at +72.
 a.emit(0x55,0x48,0x89,0xe5,0x48,0x83,0xe4,0xf0,0x48,0x83,0xec,0x60)
 a.emit(0x48,0x89,0x44,0x24,48,0x48,0xc7,0x44,0x24,56,0,0,0,0,0x48,0xc7,0x44,0x24,64,0,0,0,0)
 a.emit(0x31,0xc9,0x31,0xd2);a.rel(worker,0x4c,0x8d,0x05)
 a.emit(0x4c,0x8d,0x4c,0x24,48,0x48,0xc7,0x44,0x24,32,0,0,0,0,0x48,0xc7,0x44,0x24,40,0,0,0,0)
 a.rel("iat:70",0xff,0x15)
 a.emit(0x48,0x85,0xc0);a.j(0x84,fail)
 a.emit(0x48,0x89,0x44,0x24,72,0x48,0x89,0xc1,0xba,255,255,255,255)
 a.rel("iat:71",0xff,0x15)
 a.emit(0x85,0xc0);a.j(0x85,fail)
 a.emit(0x48,0x8b,0x4c,0x24,72);a.rel("iat:7",0xff,0x15)
 a.emit(0x48,0x8b,0x44,0x24,56,0x48,0x83,0x7c,0x24,64,0,0xc9)
 a.j(0x85,"err");a.jmp(joined)
 // On failed wait the worker may still be live: terminate process, never
 // return and invalidate its context.
 a.label(fail);a.emit(0xc9);a.jmp("thread_unsupported")
 a.label(joined)
}
func (a *asm) quantumNativeRuntime(){
 if !a.needed("thread_unsupported"){return}
 a.label("thread_unsupported");msg:="TRIPOL runtime error: native thread unavailable or failed\n"
 a.str(msg);a.emit(0x48,0x89,0xc6,0xba);a.u32(len(msg));a.call("writeerr");a.imm(16);a.jmp("exit")
}
