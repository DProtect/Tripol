package core

// Explicit temporary arena region. Outer scalar updates are safe, but references
// cannot escape: outer aggregates are inaccessible and outer text is read-only.
func (c *compiler) outerScratchBinding(name string) bool {
	boundary := c.scratchScopes[len(c.scratchScopes)-1]
	for i := len(c.scopes) - 1; i >= 0; i-- {
		if _, ok := c.scopes[i][name]; ok {
			return i < boundary
		}
	}
	return false
}
func (c *compiler) scratchBlock(s stmt) {
	a := c.a
	old := c.slots
	for i := 0; i < 4; i++ {
		c.alloc("int", true)
	}
	if a.trace {
		c.alloc("int", true)
	}
	off := -8 * c.slots
	for i := 0; i < 3; i++ {
		c.alloc("int", true)
	}
	mark := -8 * c.slots
	pending := c.alloc("int", true)
	// Ensure the arena mapping exists even if this is the first allocating scope.
	a.imm(0)
	a.call("alloc")
	a.rel("arena_state", 0x48, 0x8d, 0x15)
	for i := 0; i < 3; i++ {
		a.emit(0x48, 0x8b, 0x42, byte(i*8))
		a.store(mark + i*8)
	}
	a.imm(0)
	a.store(pending.off)
	fail, cleanup, done := a.fresh(), a.fresh(), a.fresh()
	a.installHandler(off, fail)
	c.scratchScopes = append(c.scratchScopes, len(c.scopes))
	c.protectedBlock(s.yes)
	c.scratchScopes = c.scratchScopes[:len(c.scratchScopes)-1]
	a.removeHandler(off)
	a.jmp(cleanup)
	a.label(fail)
	a.store(pending.off)
	a.label(cleanup)
	a.rel("arena_state", 0x48, 0x8d, 0x15)
	for i := 0; i < 3; i++ {
		a.load(mark + i*8)
		a.emit(0x48, 0x89, 0x42, byte(i*8))
	}
	a.load(pending.off)
	a.emit(0x48, 0x85, 0xc0)
	a.j(0x84, done)
	a.jmp("guard_raise")
	a.label(done)
	if a.opt {
		c.slots = old
	}
}
