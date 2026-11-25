package main

// Velocity R3: conservative frame elision and direct memory operand loads.
func needsFrame(ss []stmt) bool {
	for _, s := range ss {
		// A bare assignment might introduce a local. Keep a frame in that case.
		if s.op == "decl" || s.op == "assign" || s.op == "for" || needsFrame(s.yes) || needsFrame(s.no) {
			return true
		}
	}
	return false
}
func (a *asm) loadRCX(off int) {
	if a.opt && off >= -128 && off <= 127 {
		a.emit(0x48, 0x8b, 0x4d, byte(off))
		return
	}
	a.emit(0x48, 0x8b, 0x8d)
	a.u32(off)
}
func (c *compiler) fileBuiltin(e *expr) string {
	a := c.a
	if e.s == "read_text" {
		if len(e.args) != 1 {
			c.err("read_text expects one path")
		}
		c.require(c.expression(e.args[0]), "string")
		a.call("file_read")
		return "string"
	}
	if len(e.args) != 2 {
		c.err(e.s + " expects path and text")
	}
	c.require(c.expression(e.args[0]), "string")
	a.emit(0x50)
	c.require(c.expression(e.args[1]), "string")
	a.emit(0x48, 0x89, 0xc1, 0x58)
	if e.s == "append_text" {
		a.call("file_append")
	} else {
		a.call("file_write")
	}
	return "int"
}
