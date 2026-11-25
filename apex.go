package main

// APEX: register calling convention for a narrow, structurally verified
// binary recurrence. No name or benchmark-specific recognition. Debug,
// tracing and programs with exception unwinding use the generic ABI.
type recurrence struct {
	threshold, base int64
	paramBase       bool
	cc              byte
	step1, step2    byte
	op              string
}

func (c *compiler) recurrence(f fn) (r recurrence, ok bool) {
	if !c.a.hyper || c.a.trace || c.a.guard || f.record || f.ret != "int" || len(f.args) != 1 || f.args[0].typ != "int" || len(f.body) != 2 {
		return r, false
	}
	p := f.args[0].name
	a, b := f.body[0], f.body[1]
	if a.op != "if" || len(a.yes) != 1 || len(a.no) != 0 || a.yes[0].op != "return" || b.op != "return" {
		return r, false
	}
	e := a.e
	if e.op != "bin" || (e.s != "<" && e.s != "<=") || e.args[0].op != "var" || e.args[0].s != p || e.args[1].op != "int" || e.args[1].n < 0 || e.args[1].n > 127 {
		return r, false
	}
	r.threshold = e.args[1].n
	r.cc = (compareCode[e.s] - 0x10) ^ 1
	base := a.yes[0].e
	if base.op == "var" && base.s == p {
		r.paramBase = true
	} else if t, v, k := fold(base); k && t == "int" {
		r.base = v
	} else {
		return r, false
	}
	e = b.e
	if e.op != "bin" || (e.s != "+" && e.s != "-" && e.s != "*") {
		return r, false
	}
	r.op = e.s
	step := func(x *expr) (byte, bool) {
		if x.op != "call" || x.s != f.name || len(x.args) != 1 {
			return 0, false
		}
		y := x.args[0]
		if y.op != "bin" || y.s != "-" || y.args[0].op != "var" || y.args[0].s != p || y.args[1].op != "int" || y.args[1].n < 1 || y.args[1].n > 127 {
			return 0, false
		}
		return byte(y.args[1].n), true
	}
	r.step1, ok = step(e.args[0])
	if !ok {
		return r, false
	}
	r.step2, ok = step(e.args[1])
	return r, ok
}
func (a *asm) emitRecurrence(f fn, r recurrence) {
	recur := a.fresh()
	a.emit(0x49, 0x83, 0xf8, byte(r.threshold)) // cmp r8, threshold
	a.j(r.cc, recur)
	if r.paramBase {
		a.regMove(0, 8)
	} else {
		a.imm(r.base)
	}
	a.emit(0xc3)
	a.label(recur)
	// In this branch n >= 0, so subtracting either 1..127 cannot overflow.
	a.emit(0x41, 0x50) // save original n
	a.emit(0x49, 0x83, 0xe8, r.step1)
	a.call("fn:" + f.name)
	a.emit(0x50)                         // first result
	a.emit(0x4c, 0x8b, 0x44, 0x24, 0x08) // reload original n
	a.emit(0x49, 0x83, 0xe8, r.step2)
	a.call("fn:" + f.name)
	a.emit(0x59) // left result
	if r.op == "-" {
		a.regArith("-", 1, 0)
		a.regMove(0, 1)
	} else {
		a.regArith(r.op, 0, 1)
	}
	a.emit(0x48, 0x83, 0xc4, 0x08, 0xc3)
}

// One safe loop invariant may live in R15. Only total positive-power-of-two
// div/rem of an outer integer unmodified anywhere in this call-free loop.
func (c *compiler) invariantPlan(s stmt) (*expr, map[*expr]bool) {
	if !c.a.turbo {
		return nil, nil
	}
	written := map[string]bool{s.name: true}
	var writes func([]stmt)
	writes = func(ss []stmt) {
		for _, st := range ss {
			if st.op == "assign" || st.op == "decl" {
				written[st.name] = true
			}
			writes(st.yes)
			writes(st.no)
		}
	}
	writes(s.yes)
	var chosen *expr
	matches := map[*expr]bool{}
	var visit func(*expr)
	visit = func(e *expr) {
		if e == nil {
			return
		}
		if e.op == "bin" && (e.s == "/" || e.s == "%") && e.args[0].op == "var" && !written[e.args[0].s] && e.args[1].op == "int" {
			n := e.args[1].n
			if n > 1 && n <= 1<<30 && n&(n-1) == 0 {
				if v, ok := c.lookup(e.args[0].s); ok && v.typ == "int" {
					if chosen == nil {
						chosen = e
					}
					if chosen.s == e.s && chosen.args[0].s == e.args[0].s && chosen.args[1].n == n {
						matches[e] = true
					}
				}
			}
		}
		for _, x := range e.args {
			visit(x)
		}
	}
	var walk func([]stmt)
	walk = func(ss []stmt) {
		for _, st := range ss {
			visit(st.e)
			walk(st.yes)
			walk(st.no)
		}
	}
	walk(s.yes)
	return chosen, matches
}
func (c *compiler) positiveIndex(e *expr) bool {
	if !c.a.loopIndexNonnegative || e.op != "var" {
		return false
	}
	v, ok := c.lookup(e.s)
	return ok && c.a.loopRegs[v.off] == 8
}

// (x % 2^k) == 0 iff low k bits are zero, even for negative x.
func (c *compiler) branchModulo(e *expr, label string) bool {
	if !c.a.opt || e.op != "bin" || (e.s != "==" && e.s != "!=") || e.args[1].op != "int" || e.args[1].n != 0 {
		return false
	}
	x := e.args[0]
	if x.op != "bin" || x.s != "%" || x.args[1].op != "int" {
		return false
	}
	n := x.args[1].n
	if n <= 0 || n > 1<<30 || n&(n-1) != 0 {
		return false
	}
	c.require(c.expression(x.args[0]), "int")
	c.a.emit(0xa9)
	c.a.u32(int(n - 1))
	cc := byte(0x85)
	if e.s == "!=" {
		cc = 0x84
	}
	c.a.j(cc, label)
	return true
}

// Fuse adjacent checked accumulator updates using quotient and remainder of
// the same immutable-between-updates operand. Preserve both overflow checks
// and their order; never apply to general expressions or aliased accumulator.
func (c *compiler) divisionPair(s, t stmt) bool {
	a := c.a
	if !a.turbo || a.loopRegs == nil || s.op != "assign" || t.op != "assign" || s.name != t.name {
		return false
	}
	match := func(st stmt, op string) *expr {
		e := st.e
		if e == nil || e.op != "bin" || e.s != "+" || e.args[0].op != "var" || e.args[0].s != st.name {
			return nil
		}
		x := e.args[1]
		if x.op != "bin" || x.s != op || x.args[0].op != "var" || x.args[1].op != "int" {
			return nil
		}
		return x
	}
	q, m := match(s, "/"), match(t, "%")
	if q == nil || m == nil || q.args[0].s != m.args[0].s || q.args[0].s == s.name || q.args[1].n != m.args[1].n {
		return false
	}
	n := q.args[1].n
	if n <= 1 || n > 1<<30 || n&(n-1) != 0 {
		return false
	}
	v, ok := c.lookup(s.name)
	if !ok || v.typ != "int" || !v.mut {
		return false
	}
	r, ok := a.loopRegs[v.off]
	if !ok {
		return false
	}
	x, ok := c.lookup(q.args[0].s)
	if !ok || x.typ != "int" {
		return false
	}
	a.load(x.off)
	a.emit(0x48, 0x99, 0x81, 0xe2)
	a.u32(int(n - 1))
	a.emit(0x48, 0x01, 0xd0)
	a.regMove(1, 0)
	shift := byte(0)
	for d := n; d > 1; d >>= 1 {
		shift++
	}
	a.emit(0x48, 0xc1, 0xf8, shift)
	a.regArith("+", r, 0)
	a.emit(0x81, 0xe1)
	a.u32(int(n - 1))
	a.emit(0x48, 0x29, 0xd1)
	a.regArith("+", r, 1)
	return true
}
