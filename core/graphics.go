package core

// Separate graphics backend, emitted only for programs that request it.
// Software presentation is unchanged. OpenGL can be backed by a GPU or by
// the Windows software driver; a successful call is not a hardware claim.
func (c *compiler) graphicsBuiltin(e *expr)(string,bool){
 a:=c.a
 switch e.s{
 case "native_gl_frame":
  if len(e.args)!=2{c.err("gl_frame expects vertices and background")}
  c.require(c.expression(e.args[0]),"ints");a.emit(0x50)
  c.require(c.expression(e.args[1]),"int")
  a.emit(0x48,0x89,0xc1,0x58)
  if a.target=="windows"{a.call("gl_frame")}else{a.imm(0)}
  return "bool",true
 case "native_mouse_axis":
  if len(e.args)!=2{c.err("mouse_axis expects axis and logical size")}
  c.require(c.expression(e.args[0]),"int");a.emit(0x50)
  c.require(c.expression(e.args[1]),"int");a.emit(0x48,0x89,0xc1,0x58)
  if a.target=="windows"{c.emitMouseAxis()}else{a.imm(-1)}
  return "int",true
 case "native_beep":
  if len(e.args)!=2{c.err("beep expects frequency and duration")}
  c.require(c.expression(e.args[0]),"int");a.emit(0x50)
  c.require(c.expression(e.args[1]),"int");a.emit(0x48,0x89,0xc2,0x59)
  if a.target!="windows"{a.imm(0);return "bool",true}
  a.emit(0x55,0x48,0x89,0xe5,0x48,0x83,0xe4,0xf0,0x48,0x83,0xec,0x20)
  a.rel("iat:90",0xff,0x15);a.emit(0xc9,0x85,0xc0,0x0f,0x95,0xc0,0x0f,0xb6,0xc0)
  return "bool",true
 case "pad_read":
  if len(e.args)!=2{c.err("pad_read expects controller and field")}
  c.require(c.expression(e.args[0]),"int");a.emit(0x50)
  c.require(c.expression(e.args[1]),"int")
  a.emit(0x48,0x89,0xc1,0x58)
  if a.target=="windows"{a.call("pad_read")}else{a.imm(-1)}
  return "int",true
 }
 return "",false
}
func(a *asm) graphicsRuntime(){
 if a.target!="windows"{return}
 emit:=func(name string){
  t:=graphicsCode[name];a.label(name);base:=len(a.b);a.b=append(a.b,t.code...)
  for _,r:=range t.reloc{a.addFix(fix{base+r.at,r.label})}
 }
 if a.needed("gl_frame"){emit("gl_close");emit("gl_frame")}
 if a.needed("pad_read"){emit("pad_read")}
}

// Mouse is scoped to the focused game window, converted into logical pixels.
// GetCursorPos, ScreenToClient and GetClientRect use stack-local structures.
func(c *compiler) emitMouseAxis(){
 a:=c.a
 fail,selectY,scale,done:=a.fresh(),a.fresh(),a.fresh(),a.fresh()
 a.emit(0x55,0x48,0x89,0xe5,0x48,0x83,0xe4,0xf0,0x48,0x83,0xec,0x50)
 a.emit(0x48,0x89,0x44,0x24,56,0x48,0x89,0x4c,0x24,64)
 a.rel("iat:60",0xff,0x15)
 a.emit(0x48,0x85,0xc0);a.j(0x84,fail)
 a.rel("game_state",0x48,0x3b,0x05);a.j(0x85,fail)
 a.emit(0x48,0x8d,0x4c,0x24,32);a.rel("iat:91",0xff,0x15)
 a.emit(0x85,0xc0);a.j(0x84,fail)
 a.rel("game_state",0x48,0x8b,0x0d);a.emit(0x48,0x8d,0x54,0x24,32);a.rel("iat:92",0xff,0x15)
 a.emit(0x85,0xc0);a.j(0x84,fail)
 a.rel("game_state",0x48,0x8b,0x0d);a.emit(0x48,0x8d,0x54,0x24,40);a.rel("iat:58",0xff,0x15)
 a.emit(0x85,0xc0);a.j(0x84,fail)
 // Both coordinates must be in the client rectangle.
 a.emit(0x8b,0x44,0x24,32,0x3b,0x44,0x24,48);a.j(0x83,fail)
 a.emit(0x8b,0x44,0x24,36,0x3b,0x44,0x24,52);a.j(0x83,fail)
 a.emit(0x48,0x83,0x7c,0x24,56,0);a.j(0x85,selectY)
 a.emit(0x8b,0x44,0x24,32,0x8b,0x4c,0x24,48);a.jmp(scale)
 a.label(selectY);a.emit(0x8b,0x44,0x24,36,0x8b,0x4c,0x24,52)
 a.label(scale)
 a.emit(0x48,0x0f,0xaf,0x44,0x24,64,0x31,0xd2,0x48,0xf7,0xf1);a.jmp(done)
 a.label(fail);a.imm(-1)
 a.label(done);a.emit(0xc9)
}
