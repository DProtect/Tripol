package core

// Binary file operations keep bytes in checked integer arrays. Native I/O
// errors remain fatal, like the existing text operations.
func (c *compiler) byteFileBuiltin(e *expr) string {
 if len(e.args)!=2 {c.err(e.s+" expects path and limit or byte array")}
 a:=c.a
 c.require(c.expression(e.args[0]),"string");a.emit(0x50)
 if e.s=="read_bytes" {
  c.require(c.expression(e.args[1]),"int")
  a.emit(0x48,0x89,0xc1,0x58);a.call("bytes_read");return "ints"
 }
 if !c.promoteFixed(e.args[1]) {c.require(c.expression(e.args[1]),"ints")}
 a.emit(0x48,0x89,0xc1,0x58);a.call("bytes_write");return "int"
}
