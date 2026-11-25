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
