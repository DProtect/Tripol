package core

func scalarCapture(typ string) bool {
	switch typ {
	case "int", "bool", "float", "string":
		return true
	default:
		return false
	}
}

func (r *simpleResolver) captureSlot(function, name, typ string) int {
	captures := r.capturedValues[function]
	for i, capture := range captures {
		if capture.name == name {
			return i
		}
	}
	r.capturedValues[function] = append(captures, param{name, typ})
	r.changed = true
	return len(captures)
}

func (c *compiler) emitClosure(e *expr) string {
	f, ok := c.fs[e.s]
	if !ok {
		c.err("unknown closure " + e.s)
	}
	a := c.a
	for _, capture := range e.args {
		if !scalarCapture(c.expression(capture)) {
			c.err("closures capture scalar snapshots")
		}
		a.emit(0x50)
	}
	a.imm(int64((len(e.args) + 1) * 8))
	a.call("alloc")
	for i := len(e.args) - 1; i >= 0; i-- {
		a.emit(0x59, 0x48, 0x89, 0x88)
		a.u32((i + 1) * 8)
	}
	a.rel("fn:"+f.name, 0x48, 0x8d, 0x0d)
	a.emit(0x48, 0x89, 0x08, 0x48, 0x83, 0xc8, 1)
	return functionCallbackType(f)
}

func (e *webEmitter) webClosure(x *expr) string {
	f, ok := e.module.definitions[x.s]
	if !ok {
		webReject("unknown closure %s", x.s)
	}
	slot := e.module.callbackSlot(e.module.declare(f).index)
	z := e.constant(0)
	e.val(10, e.constant(int64(len(x.args)+1)), e.constant(4), z)
	env := e.temp()
	e.code.inst(0x21, env)
	e.val(12, env, z, e.constant(int64(slot)))
	e.code.raw(0x1a)
	for i, capture := range x.args {
		value, _ := e.save(capture)
		e.val(12, env, e.constant(int64(i+1)), value)
		e.code.raw(0x1a)
	}
	e.code.i64(-1)
	e.code.inst(0x20, env)
	e.code.raw(0x7d)
	return functionCallbackType(f)
}
