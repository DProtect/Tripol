package main

// ArrayCore: specialize an immutable zero-based traversal with one checked
// scalar reduction. Preserve left-to-right addition and fault precedence.
// No speculative reads, no reassociation, no elided arithmetic checks.
func (c *compiler) arrayReduction(s stmt) bool {
	a := c.a
	if !a.hyper || a.guard || a.trace || a.loopRegs != nil || len(s.yes) != 1 || s.e.op != "int" || s.e.n != 0 {
		return false
	}
	st := s.yes[0]
	e := st.e
	if st.op != "assign" || e == nil || e.op != "bin" || e.s != "+" || e.args[0].op != "var" || e.args[0].s != st.name {
		return false
	}
	x := e.args[1]
	if x.op != "index" || len(x.args) != 1 || x.args[0].op != "var" || x.args[0].s != s.name {
		return false
	}
	acc, ok := c.lookup(st.name)
	if !ok || acc.typ != "int" || !acc.mut {
		return false
	}
	arr, ok := c.lookup(x.s)
	if !ok || arr.typ != "ints" {
		return false
	}
	// Bounds expression executes once, before any array access, as before.
	c.require(c.expression(s.bound), "int")
	done, lim, loop, tail, finish := a.fresh(), a.fresh(), a.fresh(), a.fresh(), a.fresh()
	a.emit(0x48, 0x85, 0xc0)
	a.j(0x8e, done)                            // empty signed range
	a.regMove(1, 0)                            // rcx requested count
	a.load(arr.off)                            // rax header
	a.emit(0x48, 0x8b, 0x10, 0x45, 0x31, 0xc0) // rdx actual len, r8d=0
	a.emit(0x48, 0x39, 0xd1)
	a.j(0x86, lim) // unsigned n <= len
	a.regMove(1, 2)
	a.emit(0x41, 0xb8, 1, 0, 0, 0) // cap, record deferred bounds error
	a.label(lim)
	a.emit(0x48, 0x8d, 0x50, 0x08) // rdx data
	a.load(acc.off)
	// Four sequential checked adds per iteration; preserve original order.
	a.label(loop)
	a.emit(0x48, 0x83, 0xf9, 4)
	a.j(0x82, tail)
	a.emit(0x48, 0x03, 0x02)
	a.j(0x80, "err")
	for _, off := range []byte{8, 16, 24} {
		a.emit(0x48, 0x03, 0x42, off)
		a.j(0x80, "err")
	}
	a.emit(0x48, 0x83, 0xc2, 32, 0x48, 0x83, 0xe9, 4)
	a.jmp(loop)
	a.label(tail)
	a.emit(0x48, 0x85, 0xc9)
	a.j(0x84, finish)
	a.emit(0x48, 0x03, 0x02)
	a.j(0x80, "err")
	a.emit(0x48, 0x83, 0xc2, 8, 0x48, 0xff, 0xc9)
	a.jmp(tail)
	a.label(finish)
	a.emit(0x45, 0x85, 0xc0)
	a.j(0x85, "bounds")
	a.store(acc.off)
	a.label(done)
	return true
}

func (c *compiler) omegaBuiltin(e *expr) (string, bool) {
	a := c.a
	switch e.s {
	case "array_equal":
		if len(e.args) != 2 {
			c.err("array_equal expects two arrays")
		}
		c.require(c.expression(e.args[0]), "ints")
		a.emit(0x50)
		c.require(c.expression(e.args[1]), "ints")
		a.emit(0x48, 0x89, 0xc1, 0x58)
		a.call("array_equal")
		return "bool", true
	case "env_get":
		if len(e.args) != 1 {
			c.err("env_get expects one name")
		}
		c.require(c.expression(e.args[0]), "string")
		a.call("env_get")
		return "string", true
	case "discord_request":
		if len(e.args) != 3 {
			c.err("discord_request expects path, token, body")
		}
		for _, x := range e.args {
			c.require(c.expression(x), "string")
			a.emit(0x50)
		}
		a.emit(0x5a, 0x59, 0x58)
		a.call("discord_request")
		return "string", true
	}
	return "", false
}
