package main

import "fmt"

// TIR R1 is a deliberately narrow typed intermediate representation for
// side-effect-free int64 return expressions. All other language constructs
// continue through the existing checked backend. The debug validation pass
// runs before this representation is used. Checked arithmetic is retained.
type tirOp uint8

const (
	tirConst tirOp = iota
	tirParam
	tirAdd
	tirSub
	tirMul
)

type tirNode struct {
	op          tirOp
	value       int64
	left, right *tirNode
}
type tirFunction struct {
	root  *tirNode
	arity int
}

func lowerTIR(f fn) (tirFunction, bool) {
	if f.record || f.ret != "int" || len(f.args) > 4 || len(f.body) != 1 || f.body[0].op != "return" {
		return tirFunction{}, false
	}
	indices := map[string]int{}
	for i, p := range f.args {
		if p.typ != "int" {
			return tirFunction{}, false
		}
		indices[p.name] = i
	}
	budget := 9
	var lower func(*expr) *tirNode
	lower = func(e *expr) *tirNode {
		if e == nil || budget <= 0 {
			return nil
		}
		budget--
		if t, n, ok := fold(e); ok && t == "int" {
			return &tirNode{op: tirConst, value: n}
		}
		switch e.op {
		case "var":
			if i, ok := indices[e.s]; ok {
				return &tirNode{op: tirParam, value: int64(i)}
			}
		case "bin":
			op, ok := map[string]tirOp{"+": tirAdd, "-": tirSub, "*": tirMul}[e.s]
			if !ok {
				return nil
			}
			l, r := lower(e.args[0]), lower(e.args[1])
			if l == nil || r == nil {
				return nil
			}
			if r.op == tirConst && ((r.value == 0 && (op == tirAdd || op == tirSub)) || (r.value == 1 && op == tirMul)) {
				return l
			}
			return &tirNode{op: op, left: l, right: r}
		}
		return nil
	}
	root := lower(f.body[0].e)
	return tirFunction{root, len(f.args)}, root != nil
}
func (a *asm) emitTIR(n *tirNode, arity int) {
	switch n.op {
	case tirConst:
		a.imm(n.value)
	case tirParam:
		if arity == 1 {
			a.emit(0x4c, 0x89, 0xd8)
			return
		}
		// r11 holds the original argument stack top; expression pushes do not change it.
		off := 8 * (arity - 1 - int(n.value))
		if off == 0 {
			a.emit(0x49, 0x8b, 0x03)
		} else {
			a.emit(0x49, 0x8b, 0x43, byte(off))
		}
	default:
		a.emitTIR(n.left, arity)
		if n.right.op == tirParam && arity == 1 {
			switch n.op {
			case tirAdd:
				a.emit(0x4c, 0x01, 0xd8)
			case tirSub:
				a.emit(0x4c, 0x29, 0xd8)
			case tirMul:
				a.emit(0x49, 0x0f, 0xaf, 0xc3)
			}
			a.j(0x80, "err")
			return
		}
		// For multi-argument scalar expressions, operate directly on the
		// original argument stack slot. No temporary expression push/pop.
		if n.right.op == tirParam && arity > 1 {
			off := byte(8 * (arity - 1 - int(n.right.value)))
			switch n.op {
			case tirAdd:
				a.emit(0x49, 0x03, 0x43, off)
			case tirSub:
				a.emit(0x49, 0x2b, 0x43, off)
			case tirMul:
				a.emit(0x49, 0x0f, 0xaf, 0x43, off)
			}
			a.j(0x80, "err")
			return
		}
		if n.right.op == tirConst && n.right.value >= -128 && n.right.value <= 127 {
			v := byte(n.right.value)
			switch n.op {
			case tirAdd:
				a.emit(0x48, 0x83, 0xc0, v)
			case tirSub:
				a.emit(0x48, 0x83, 0xe8, v)
			case tirMul:
				a.emit(0x48, 0x6b, 0xc0, v)
			}
			a.j(0x80, "err")
			return
		}
		a.emit(0x50)
		a.emitTIR(n.right, arity)
		a.emit(0x48, 0x89, 0xc1, 0x58)
		switch n.op {
		case tirAdd:
			a.emit(0x48, 0x01, 0xc8)
		case tirSub:
			a.emit(0x48, 0x29, 0xc8)
		case tirMul:
			a.emit(0x48, 0x0f, 0xaf, 0xc1)
		}
		a.j(0x80, "err")
	}
}
func (c *compiler) inlineTIR(f fn, args []*expr) bool {
	a := c.a
	if !a.opt || a.trace {
		return false
	}
	if len(args) > 1 && !a.turbo {
		return false
	}
	ir, cached := c.tirCache[f.name]
	if !cached {
		ir, _ = lowerTIR(f)
		c.tirCache[f.name] = ir
	}
	if ir.root == nil {
		return false
	}
	if a.hyper && len(args) > 1 {
		c.registerArgs(args)
		a.emitRegTIR(ir.root, 0, []byte{1, 2})
		return true
	}
	// Evaluation remains left-to-right and exactly once, even for unused parameters.
	if len(args) == 1 {
		c.require(c.expression(args[0]), "int")
		root := ir.root
		if root.op == tirParam {
			return true
		}
		if root.op == tirConst {
			a.imm(root.value)
			return true
		}
		if root.left.op == tirParam && root.right.op == tirParam {
			switch root.op {
			case tirAdd:
				a.emit(0x48, 0x01, 0xc0)
			case tirSub:
				a.emit(0x48, 0x29, 0xc0)
			case tirMul:
				a.emit(0x48, 0x0f, 0xaf, 0xc0)
			}
			a.j(0x80, "err")
			return true
		}
		if root.left.op == tirParam && root.right.op == tirConst && root.right.value >= -128 && root.right.value <= 127 {
			v := byte(root.right.value)
			switch root.op {
			case tirAdd:
				a.emit(0x48, 0x83, 0xc0, v)
			case tirSub:
				a.emit(0x48, 0x83, 0xe8, v)
			case tirMul:
				a.emit(0x48, 0x6b, 0xc0, v)
			}
			a.j(0x80, "err")
			return true
		}
		a.emit(0x49, 0x89, 0xc3) // mov r11,rax
		a.emitTIR(ir.root, 1)
		return true
	}
	for _, arg := range args {
		c.require(c.expression(arg), "int")
		a.emit(0x50)
	}
	a.emit(0x49, 0x89, 0xe3)
	a.emitTIR(ir.root, ir.arity)
	if len(args) > 0 {
		a.emit(0x48, 0x83, 0xc4, byte(len(args)*8))
	}
	return true
}

// Human-readable view of the subset eligible for TIR inlining.
func (n *tirNode) String() string {
	switch n.op {
	case tirConst:
		return fmt.Sprintf("const.i64 %d", n.value)
	case tirParam:
		return fmt.Sprintf("param.i64 %d", n.value)
	default:
		return fmt.Sprintf("(%s.checked.i64 %s %s)", map[tirOp]string{tirAdd: "add", tirSub: "sub", tirMul: "mul"}[n.op], n.left, n.right)
	}
}
func inspectTIR(file string) {
	fs, ss, _ := loadProject(file)
	generate(fs, ss, "linux", false, nil)
	fmt.Println("TRIPOL TIR R4 — lowered scalar functions")
	n := 0
	for _, f := range fs {
		if ir, ok := lowerTIR(f); ok {
			fmt.Printf("%s(%d): %s\n", f.name, ir.arity, ir.root)
			n++
		}
	}
	fmt.Printf("%d lowered functions; nano: zero/one argument; balanced: register leaf ABI; turbo: up to four arguments may inline.\n", n)
}
