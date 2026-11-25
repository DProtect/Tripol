package main

import "encoding/binary"

// Register-backed straight-line integer loops. No calls, nested control flow,
// guards or tracing: no register value survives an unknown clobber.
func (a *asm) regMove(dst, src byte) {
	if dst == src {
		return
	}
	rex := byte(0x48)
	if src >= 8 {
		rex |= 4
	}
	if dst >= 8 {
		rex |= 1
	}
	a.emit(rex, 0x89, 0xc0|((src&7)<<3)|(dst&7))
}
func (a *asm) regArith(op string, dst, src byte) bool {
	switch op {
	case "+", "-":
		rex := byte(0x48)
		if src >= 8 {
			rex |= 4
		}
		if dst >= 8 {
			rex |= 1
		}
		code := byte(0x01)
		if op == "-" {
			code = 0x29
		}
		a.emit(rex, code, 0xc0|((src&7)<<3)|(dst&7))
	case "*":
		rex := byte(0x48)
		if dst >= 8 {
			rex |= 4
		}
		if src >= 8 {
			rex |= 1
		}
		a.emit(rex, 0x0f, 0xaf, 0xc0|((dst&7)<<3)|(src&7))
	default:
		return false
	}
	a.j(0x80, "err")
	return true
}
func (a *asm) popReg(reg byte) {
	if reg >= 8 {
		a.emit(0x41)
	}
	a.emit(0x58 | (reg & 7))
}

// Promote multiple outer scalars. Branch-local declarations fall back.
func (c *compiler) loopPlan(s stmt) ([]variable, bool) {
	a := c.a
	if !a.hyper || a.guard || a.trace || a.loopRegs != nil || len(s.yes) == 0 {
		return nil, false
	}
	locals := map[string]bool{s.name: true}
	var scalar func(*expr) bool
	scalar = func(e *expr) bool {
		if e == nil {
			return false
		}
		switch e.op {
		case "int":
			return true
		case "var":
			if locals[e.s] {
				return true
			}
			v, ok := c.lookup(e.s)
			return ok && v.typ == "int"
		case "un":
			return e.s == "-" && scalar(e.args[0])
		case "bin":
			switch e.s {
			case "+", "-", "*", "/", "%":
				return scalar(e.args[0]) && scalar(e.args[1])
			}
		}
		return false
	}
	var condition func(*expr) bool
	condition = func(e *expr) bool {
		if e == nil {
			return false
		}
		switch e.op {
		case "bool":
			return true
		case "un":
			return e.s == "!" && condition(e.args[0])
		case "bin":
			switch e.s {
			case "&&", "||":
				return condition(e.args[0]) && condition(e.args[1])
			case "==", "!=", "<", ">", "<=", ">=":
				return scalar(e.args[0]) && scalar(e.args[1])
			}
		}
		return false
	}
	var acc []variable
	seen := map[int]bool{}
	var visit func([]stmt, bool) bool
	visit = func(ss []stmt, nested bool) bool {
		for _, st := range ss {
			switch st.op {
			case "if":
				if !condition(st.e) || !visit(st.yes, true) || !visit(st.no, true) {
					return false
				}
			case "break", "continue":
			case "assign", "decl":
				if !scalar(st.e) || (st.typ != "" && st.typ != "int") {
					return false
				}
				if locals[st.name] {
					if st.op == "decl" || st.name == s.name {
						return false
					}
					continue
				}
				if v, ok := c.lookup(st.name); ok {
					if v.typ != "int" || !v.mut || st.op == "decl" {
						return false
					}
					if !seen[v.off] {
						seen[v.off] = true
						acc = append(acc, v)
					}
				} else {
					if nested {
						return false
					}
					locals[st.name] = true
				}
			default:
				return false
			}
		}
		return true
	}
	if !visit(s.yes, false) {
		return nil, false
	}
	if len(acc) > 4 {
		acc = acc[:4]
	}
	return acc, true
}

// Private register ABI R8..R11 for pure checked int64 leaf expressions only.
func (c *compiler) nativeTIR(f fn) (tirFunction, bool) {
	if !c.a.hyper || c.a.trace || len(f.args) == 0 {
		return tirFunction{}, false
	}
	ir, ok := c.tirCache[f.name]
	if !ok {
		ir, _ = lowerTIR(f)
		c.tirCache[f.name] = ir
	}
	return ir, ir.root != nil
}
func (c *compiler) registerArgs(args []*expr) {
	direct := true
	for _, x := range args {
		if x.op != "var" && x.op != "int" {
			direct = false
		}
	}
	if direct {
		for i, x := range args {
			c.require(c.expression(x), "int")
			c.a.regMove(byte(8+i), 0)
		}
		return
	}
	for _, x := range args {
		c.require(c.expression(x), "int")
		c.a.emit(0x50)
	}
	for i := len(args) - 1; i >= 0; i-- {
		c.a.popReg(byte(8 + i))
	}
}

// Three expression registers, immutable parameter registers, bounded spills.
func (a *asm) emitRegTIR(n *tirNode, dst byte, free []byte) {
	switch n.op {
	case tirConst:
		if dst == 0 {
			a.imm(n.value)
		} else {
			a.emit(0x48, 0xb8|dst)
			a.b = binary.LittleEndian.AppendUint64(a.b, uint64(n.value))
		}
	case tirParam:
		a.regMove(dst, byte(8+n.value))
	default:
		a.emitRegTIR(n.left, dst, free)
		op := map[tirOp]string{tirAdd: "+", tirSub: "-", tirMul: "*"}[n.op]
		if n.right.op == tirConst && n.right.value >= -128 && n.right.value <= 127 {
			if n.op == tirMul {
				a.emit(0x48, 0x6b, 0xc0|(dst<<3)|dst, byte(n.right.value))
			} else {
				code := byte(0xc0)
				if n.op == tirSub {
					code = 0xe8
				}
				a.emit(0x48, 0x83, code|dst, byte(n.right.value))
			}
			a.j(0x80, "err")
			return
		}
		if n.right.op == tirParam {
			a.regArith(op, dst, byte(8+n.right.value))
			return
		}
		if len(free) > 0 {
			r := free[0]
			a.emitRegTIR(n.right, r, free[1:])
			a.regArith(op, dst, r)
			return
		}
		scratch := byte(1)
		if dst == 1 {
			scratch = 2
		}
		a.emit(0x50|scratch, 0x50|dst)
		a.emitRegTIR(n.right, dst, nil)
		a.regMove(scratch, dst)
		a.popReg(dst)
		a.regArith(op, dst, scratch)
		a.popReg(scratch)
	}
}
