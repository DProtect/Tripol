package main

// Guard R2: stack-resident exception handlers. Only safe language traps are
// recoverable; native I/O and allocator failures remain fatal until cleanup exists.
func containsGuard(ss []stmt) bool {
	for _, s := range ss {
		if s.op == "try" || s.op == "scratch" || s.op == "throw" || containsGuard(s.yes) || containsGuard(s.no) || containsGuard(s.final) {
			return true
		}
	}
	return false
}
func (a *asm) guardPtr() { a.rel("guard_state", 0x48, 0x8d, 0x15) }
func (a *asm) installHandler(off int, target string) {
	a.guardPtr()
	a.emit(0x48, 0x8b, 0x0a) // previous handler
	a.emit(0x48, 0x89, 0x8d)
	a.u32(off)
	a.emit(0x48, 0x89, 0xa5)
	a.u32(off + 8) // rsp
	a.emit(0x48, 0x89, 0xad)
	a.u32(off + 16) // rbp
	a.rel(target, 0x48, 0x8d, 0x0d)
	a.emit(0x48, 0x89, 0x8d)
	a.u32(off + 24)
	if a.trace {
		a.emit(0x4c, 0x89, 0xbd)
		a.u32(off + 32)
	}
	a.emit(0x48, 0x8d, 0x8d)
	a.u32(off)
	a.emit(0x48, 0x89, 0x0a)
}
func (a *asm) removeHandler(off int) {
	a.guardPtr()
	a.emit(0x48, 0x8b, 0x8d)
	a.u32(off)
	a.emit(0x48, 0x89, 0x0a)
}
func (c *compiler) protectedBlock(ss []stmt) {
	c.protected = append(c.protected, len(c.loops))
	c.block(ss, true)
	c.protected = c.protected[:len(c.protected)-1]
}
func (c *compiler) guardBlock(s stmt) {
	a := c.a
	old := c.slots
	if _, ok := c.lookup(s.name); ok {
		c.err("duplicate catch variable " + s.name)
	}
	for i := 0; i < 4; i++ {
		c.alloc("int", true)
	}
	if a.trace {
		c.alloc("int", true)
	}
	off := -8 * c.slots
	pending := c.alloc("int", true)
	caught := c.alloc("int", false)
	handler, catchFail, final, done := a.fresh(), a.fresh(), a.fresh(), a.fresh()
	a.imm(0)
	a.store(pending.off)
	a.installHandler(off, handler)
	c.protectedBlock(s.yes)
	a.removeHandler(off)
	a.jmp(final)
	a.label(handler)
	a.store(caught.off)
	a.installHandler(off, catchFail)
	c.scopes = append(c.scopes, map[string]variable{s.name: caught})
	c.protectedBlock(s.no)
	c.scopes = c.scopes[:len(c.scopes)-1]
	a.removeHandler(off)
	a.jmp(final)
	a.label(catchFail)
	a.store(pending.off)
	a.label(final)
	c.protectedBlock(s.final)
	a.load(pending.off)
	a.emit(0x48, 0x85, 0xc0)
	a.j(0x84, done)
	a.jmp("guard_raise")
	a.label(done)
	if a.opt {
		c.slots = old
	}
}
func (a *asm) guardTrap(code int) {
	if !a.guard {
		return
	}
	a.guardPtr()
	a.emit(0x48, 0x83, 0x3a, 0)
	skip := a.fresh()
	a.j(0x84, skip)
	a.imm(int64(code))
	a.jmp("guard_raise")
	a.label(skip)
}
func (a *asm) guardRuntime() {
	if !a.guard {
		return
	}
	a.label("guard_raise")
	a.guardPtr()
	a.emit(0x48, 0x8b, 0x0a, 0x48, 0x85, 0xc9)
	fatal := a.fresh()
	a.j(0x84, fatal)
	// Pop handler, restore stack/frame, jump to its landing pad, keeping RAX=code.
	a.emit(0x4c, 0x8b, 0x01, 0x4c, 0x89, 0x02, 0x48, 0x8b, 0x61, 0x08, 0x48, 0x8b, 0x69, 0x10)
	if a.trace {
		a.emit(0x4c, 0x8b, 0x79, 0x20)
	}
	a.emit(0xff, 0x61, 0x18)
	a.label(fatal)
	a.emit(0x50)
	msg := "TRIPOL runtime error: unhandled exception\n"
	a.str(msg)
	a.emit(0x48, 0x89, 0xc6, 0xba)
	a.u32(len(msg))
	a.call("writeerr")
	a.emit(0x58)
	a.jmp("exit")
}
