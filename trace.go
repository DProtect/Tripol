package main

// Optional language-frame trace. r15 is reserved for the trace chain (and is
// nonvolatile under Win64); the native helper templates do not modify it.
func (c *compiler) traceEnter(f fn) {
	a := c.a
	for i := 0; i < 3; i++ {
		c.alloc("int", true)
	}
	off := -8 * c.slots
	a.emit(0x4c, 0x89, 0xbd)
	a.u32(off)
	msg := "  at " + f.name + "\n"
	a.str(msg)
	a.store(off + 8)
	a.imm(int64(len(msg)))
	a.store(off + 16)
	a.emit(0x4c, 0x8d, 0xbd)
	a.u32(off)
}
func (a *asm) traceRuntime() {
	if !a.trace {
		return
	}
	a.label("trace_dump")
	loop, done := a.fresh(), a.fresh()
	a.label(loop)
	a.emit(0x4d, 0x85, 0xff)
	a.j(0x84, done)
	a.emit(0x49, 0x8b, 0x77, 0x08, 0x49, 0x8b, 0x57, 0x10)
	a.call("writeerr")
	a.emit(0x4d, 0x8b, 0x3f)
	a.jmp(loop)
	a.label(done)
	a.emit(0xc3)
}
