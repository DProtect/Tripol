package core

import "strings"

func optionalElement(typ string) string {
	if strings.HasPrefix(typ, callbackPrefix) {
		return ""
	}
	element, ok := strings.CutSuffix(typ, "?")
	if !ok {
		return ""
	}
	return element
}

func optionalType(element string) string {
	if element == "" || element == "void" || optionalElement(element) != "" ||
		strings.HasPrefix(element, callbackPrefix) || strings.HasPrefix(element, "@") || strings.HasPrefix(element, "array:") {
		return ""
	}
	return element + "?"
}

func matchDefault(typ string, value *expr) bool {
	if value.op == "null" && optionalElement(typ) != "" {
		value.op, value.s = "none", typ
		return true
	}
	return value.op == typ
}

func (r *simpleResolver) optionalMethod(element, method string, args []*expr, scope simpleScope) string {
	if method == "or_else" && len(args) == 1 {
		r.expression(args[0], element, scope)
	}
	if method == "has" {
		return "bool"
	}
	return element
}

func (c *compiler) emitOptionalMethod(element, method string, args []*expr) string {
	a := c.a
	switch method {
	case "has":
		if len(args) != 0 {
			c.err("has() takes no arguments")
		}
		a.emit(0x48, 0x85, 0xc0, 0x0f, 0x95, 0xc0, 0x48, 0x0f, 0xb6, 0xc0)
		return "bool"
	case "unwrap":
		if len(args) != 0 {
			c.err("unwrap() takes no arguments")
		}
		a.emit(0x48, 0x85, 0xc0)
		a.j(0x84, "bounds")
		a.emit(0x48, 0x8b, 0)
	case "or_else":
		if len(args) != 1 {
			c.err("or_else() needs one fallback value")
		}
		absent, done := a.fresh(), a.fresh()
		a.emit(0x48, 0x85, 0xc0)
		a.j(0x84, absent)
		a.emit(0x48, 0x8b, 0)
		a.jmp(done)
		a.label(absent)
		c.require(c.expression(args[0]), element)
		a.label(done)
	default:
		c.err("unknown optional method " + method + "; try has(), unwrap() or or_else(fallback)")
	}
	return element
}

func (e *webEmitter) webOptional(receiver uint32, element, method string, args []*expr) string {
	e.code.inst(0x20, receiver)
	e.code.raw(0x50) // i64.eqz
	if method == "has" {
		e.code.raw(0x45, 0xad) // i32.eqz, i64.extend_i32_u
		return "bool"
	}
	read := func() {
		e.code.inst(0x20, receiver)
		e.code.i64(1)
		e.code.raw(0x7d)
		slot := e.temp()
		e.code.inst(0x21, slot)
		z := e.constant(0)
		e.val(11, slot, z, z)
	}
	if method == "or_else" {
		e.code.raw(0x04, 0x7e)
		if e.expression(args[0]) != element {
			webReject("optional fallback must be %s", element)
		}
		e.code.raw(0x05)
		read()
		e.code.raw(0x0b)
		return element
	}
	e.code.raw(0x04, 0x40)
	e.code.i32(int64(e.line))
	e.code.inst(0x10, 3)
	e.code.raw(0x0b)
	read()
	return element
}
