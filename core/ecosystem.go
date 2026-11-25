package core

// byte_at returns a UTF-8 byte without allocating a substring. Bounds are checked.
func (c *compiler) byteAt(e *expr) string {
	if len(e.args) != 2 {
		c.err("byte_at expects text and index")
	}
	a := c.a
	c.require(c.expression(e.args[0]), "string")
	a.emit(0x50)
	c.require(c.expression(e.args[1]), "int")
	a.emit(0x50)
	a.emit(0x48, 0x8b, 0x44, 0x24, 0x08)
	a.call("strlen")
	a.emit(0x59, 0x5a, 0x48, 0x39, 0xc1)
	a.j(0x83, "bounds") // unsigned index >= length also catches negatives
	a.emit(0x0f, 0xb6, 0x04, 0x0a)
	return "int"
}

// Print a fixed or dynamic integer array without allocating or evaluating twice.
func (c *compiler) showValue(e *expr, newline bool) {
	a := c.a
	t := "ints"
	if !c.promoteFixed(e) {
		t = c.expression(e)
	}
	if t != "ints" {
		if c.isRecord(t) {
			c.err("show requires scalar, string or integer array")
		}
		if newline {
			a.call("print:" + t)
		} else {
			a.call("raw_" + t)
		}
		return
	}
	a.emit(0x50)
	a.str("[")
	a.call("raw_string")
	a.emit(0x6a, 0x00) // index, pointer below
	loop, end, first := a.fresh(), a.fresh(), a.fresh()
	a.label(loop)
	a.emit(0x48, 0x8b, 0x0c, 0x24, 0x48, 0x8b, 0x54, 0x24, 0x08, 0x48, 0x3b, 0x0a)
	a.j(0x83, end)
	a.emit(0x48, 0x85, 0xc9)
	a.j(0x84, first)
	a.str(", ")
	a.call("raw_string")
	a.label(first)
	a.emit(0x48, 0x8b, 0x0c, 0x24, 0x48, 0x8b, 0x54, 0x24, 0x08, 0x48, 0x8b, 0x44, 0xca, 0x08)
	a.call("raw_int")
	a.emit(0x48, 0xff, 0x04, 0x24)
	a.jmp(loop)
	a.label(end)
	a.emit(0x48, 0x83, 0xc4, 0x10)
	if newline {
		a.str("]\n")
	} else {
		a.str("]")
	}
	a.call("raw_string")
}

func (c *compiler) ecosystemBuiltin(e *expr) string {
	a := c.a
	switch e.s {
	case "platform_name", "compiler_version", "monotonic_ms", "monotonic_ns":
		if len(e.args) != 0 {
			c.err(e.s + " expects no arguments")
		}
		if e.s == "monotonic_ns" {
			a.call("clock_ns")
			return "int"
		}
		if e.s == "monotonic_ms" {
			a.call("clock_ms")
			return "int"
		}
		if e.s == "platform_name" {
			a.str(a.target)
		} else {
			a.str(version)
		}
		return "string"
	case "sleep_ms":
		if len(e.args) != 1 {
			c.err("sleep_ms expects one integer")
		}
		c.require(c.expression(e.args[0]), "int")
		a.call("sleep_ms")
		return "int"
	case "has_prefix", "has_suffix":
		if len(e.args) != 2 {
			c.err(e.s + " expects two strings")
		}
		c.require(c.expression(e.args[0]), "string")
		a.emit(0x50)
		c.require(c.expression(e.args[1]), "string")
		a.emit(0x48, 0x89, 0xc1, 0x58)
		a.call(e.s)
		return "bool"
	}
	c.err("unknown ecosystem intrinsic")
	return ""
}
