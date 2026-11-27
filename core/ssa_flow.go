package core

// Flow SSA uses block parameters for values read before their first local
// definition. Resolve incoming values only after all backedges are known.
const maxFlowValues = 512
const maxFlowBlocks = 96

type flowValue struct {
	op, typ string
	x, y    int
	number  int64
	inputs  []int
}

type flowSSABlock struct {
	code, phis, predecessors []int
	entry, definitions       map[int]int
	term                     string
	value, yes, no           int
}

type flowBinding struct {
	id  int
	typ string
}

type flowSSA struct {
	blocks []flowSSABlock
	values []flowValue
}

type flowBuilder struct {
	flowSSA
	current, bindings int
	scope             map[string]flowBinding
	failed            bool
}

func (b *flowBuilder) newBlock() int {
	id := len(b.blocks)
	b.blocks = append(b.blocks, flowSSABlock{
		entry: make(map[int]int), definitions: make(map[int]int),
	})
	b.failed = b.failed || len(b.blocks) > maxFlowBlocks
	return id
}

func (b *flowBuilder) newValue(v flowValue) int {
	id := len(b.values)
	b.values = append(b.values, v)
	b.failed = b.failed || len(b.values) > maxFlowValues
	return id
}

func (b *flowBuilder) emit(v flowValue) int {
	id := b.newValue(v)
	n := b.current
	b.blocks[n].code = append(b.blocks[n].code, id)
	return id
}

func (b *flowBuilder) literal(kind string, n int64) int {
	return b.emit(flowValue{op: "literal", typ: kind, number: n, x: -1, y: -1})
}

func (b *flowBuilder) binding(typ string) flowBinding {
	v := flowBinding{id: b.bindings, typ: typ}
	b.bindings++
	return v
}

func (b *flowBuilder) write(v flowBinding, id int) {
	b.blocks[b.current].definitions[v.id] = id
}

func (b *flowBuilder) read(block int, v flowBinding) int {
	if id, ok := b.blocks[block].definitions[v.id]; ok {
		return id
	}
	return b.readEntry(block, v)
}

func (b *flowBuilder) readEntry(block int, v flowBinding) int {
	if id, ok := b.blocks[block].entry[v.id]; ok {
		return id
	}
	id := b.newValue(flowValue{op: "phi", typ: v.typ, number: int64(v.id), x: -1, y: -1})
	b.blocks[block].entry[v.id] = id
	b.blocks[block].phis = append(b.blocks[block].phis, id)
	return id
}

func (b *flowBuilder) terminate(op string, value, yes, no int) {
	block := b.current
	dst := &b.blocks[block]
	dst.term, dst.value, dst.yes, dst.no = op, value, yes, no
	if yes >= 0 {
		b.blocks[yes].predecessors = append(b.blocks[yes].predecessors, block)
	}
	if no >= 0 {
		b.blocks[no].predecessors = append(b.blocks[no].predecessors, block)
	}
	b.current = -1
}

func (b *flowBuilder) jump(dst int) {
	if b.current >= 0 {
		b.terminate("jump", -1, dst, -1)
	}
}

func (b *flowBuilder) expression(e *expr) int {
	if e == nil || b.failed {
		b.failed = true
		return -1
	}
	if typ, n, ok := fold(e); ok {
		return b.literal(typ, n)
	}
	switch e.op {
	case "var":
		if v, ok := b.scope[e.s]; ok {
			return b.read(b.current, v)
		}
	case "un":
		x := b.expression(e.args[0])
		if x >= 0 {
			kind := b.values[x].typ
			if (e.s == "-" && kind == "int") || (e.s == "!" && kind == "bool") {
				return b.emit(flowValue{op: "un" + e.s, typ: kind, x: x, y: -1})
			}
		}
	case "bin":
		x := b.expression(e.args[0])
		if x < 0 {
			return -1
		}
		if e.s == "&&" || e.s == "||" {
			if b.values[x].typ != "bool" {
				break
			}
			return b.logical(e, x)
		}
		y := b.expression(e.args[1])
		if y < 0 {
			return -1
		}
		kind := b.values[x].typ
		if kind != b.values[y].typ {
			break
		}
		if e.s == "==" || e.s == "!=" {
			return b.emit(flowValue{op: e.s, typ: "bool", x: x, y: y})
		}
		if kind != "int" {
			break
		}
		switch e.s {
		case "<", "<=", ">", ">=":
			return b.emit(flowValue{op: e.s, typ: "bool", x: x, y: y})
		case "+", "-", "*", "/", "%":
			return b.emit(flowValue{op: e.s, typ: "int", x: x, y: y})
		}
	}
	b.failed = true
	return -1
}

func (b *flowBuilder) logical(e *expr, x int) int {
	rhs, join := b.newBlock(), b.newBlock()
	result := b.binding("bool")
	// The condition itself is the value on the short-circuit edge.
	b.write(result, x)
	if e.s == "&&" {
		b.terminate("branch", x, rhs, join)
	} else {
		b.terminate("branch", x, join, rhs)
	}
	b.current = rhs
	y := b.expression(e.args[1])
	if y < 0 || b.values[y].typ != "bool" {
		b.failed = true
		return -1
	}
	b.write(result, y)
	b.jump(join)
	b.current = join
	return b.read(b.current, result)
}

func (b *flowBuilder) block(body []stmt, loop *flowLoop, nested bool) {
	// The validated language forbids shadowing, so only newly introduced
	// bindings need to be removed when a lexical scope closes.
	var locals []string
	for _, s := range body {
		if b.failed || b.current < 0 {
			break
		}
		switch s.op {
		case "assign", "decl":
			id := b.expression(s.e)
			if id < 0 {
				break
			}
			v, exists := b.scope[s.name]
			if !exists {
				v = b.binding(b.values[id].typ)
				b.scope[s.name] = v
				if nested {
					locals = append(locals, s.name)
				}
			}
			b.write(v, id)
		case "return":
			if id := b.expression(s.e); id >= 0 {
				b.terminate("return", id, -1, -1)
			}
		case "break", "continue":
			if loop == nil {
				b.failed = true
			} else if s.op == "break" {
				b.jump(loop.end)
			} else {
				b.jump(loop.again)
			}
		case "if":
			b.ifStatement(s, loop)
		case "while", "for":
			b.loopStatement(s)
		default:
			b.failed = true
		}
	}
	for _, name := range locals {
		delete(b.scope, name)
	}
}

func (b *flowBuilder) ifStatement(s stmt, loop *flowLoop) {
	if kind, n, ok := fold(s.e); ok && kind == "bool" {
		if n == 0 {
			b.block(s.no, loop, true)
		} else {
			b.block(s.yes, loop, true)
		}
		return
	}
	id := b.expression(s.e)
	if id < 0 {
		return
	}
	yes, no, join := b.newBlock(), b.newBlock(), b.newBlock()
	b.terminate("branch", id, yes, no)
	b.current = yes
	b.block(s.yes, loop, true)
	b.jump(join)
	b.current = no
	b.block(s.no, loop, true)
	b.jump(join)
	if len(b.blocks[join].predecessors) != 0 {
		b.current = join
	}
}

func (b *flowBuilder) loopStatement(s stmt) {
	head, body, exit := b.newBlock(), b.newBlock(), b.newBlock()
	again := head
	var index, bound flowBinding
	if s.op == "for" {
		start, stop := b.expression(s.e), b.expression(s.bound)
		if start < 0 || stop < 0 {
			return
		}
		index, bound = b.binding("int"), b.binding("int")
		b.write(index, start)
		b.write(bound, stop)
		b.scope[s.name] = index
		defer delete(b.scope, s.name)
		again = b.newBlock()
	}
	b.jump(head)
	b.current = head
	var cond int
	if s.op == "for" {
		x, y := b.read(head, index), b.read(head, bound)
		cond = b.emit(flowValue{op: "<", typ: "bool", x: x, y: y})
	} else {
		cond = b.expression(s.e)
		if cond < 0 {
			return
		}
	}
	b.terminate("branch", cond, body, exit)
	b.current = body
	b.block(s.yes, &flowLoop{end: exit, again: again}, true)
	b.jump(again)
	if s.op == "for" && len(b.blocks[again].predecessors) != 0 {
		b.current = again
		x, y := b.read(again, index), b.literal("int", 1)
		b.write(index, b.emit(flowValue{op: "+", typ: "int", x: x, y: y}))
		b.jump(head)
	}
	b.current = exit
}

func buildFlowSSA(f fn) (flowSSA, bool) {
	if f.record || f.ret != "int" && f.ret != "bool" {
		return flowSSA{}, false
	}
	b := flowBuilder{scope: make(map[string]flowBinding)}
	b.current = b.newBlock()
	for i, arg := range f.args {
		if arg.typ != "int" && arg.typ != "bool" {
			return flowSSA{}, false
		}
		v := b.binding(arg.typ)
		b.scope[arg.name] = v
		b.write(v, b.emit(flowValue{op: "param", typ: arg.typ, number: int64(i), x: -1, y: -1}))
	}
	b.block(f.body, nil, false)
	if b.failed || b.current >= 0 || len(b.blocks) == 1 {
		return flowSSA{}, false
	}
	// Each unresolved entry read may introduce reads in predecessors.
	// Iterate until every phi has materialized its incoming values.
	for changed := true; changed && !b.failed; {
		changed = false
		for block := range b.blocks {
			for _, id := range b.blocks[block].phis {
				v := b.values[id]
				if v.inputs != nil {
					continue
				}
				preds := b.blocks[block].predecessors
				if len(preds) == 0 {
					return flowSSA{}, false
				}
				inputs := make([]int, len(preds))
				for i, pred := range preds {
					inputs[i] = b.read(pred, flowBinding{id: int(v.number), typ: v.typ})
				}
				b.values[id].inputs = inputs
				changed = true
			}
		}
	}
	if b.failed {
		return flowSSA{}, false
	}
	b.flowSSA.removeTrivialPhis()
	return b.flowSSA, true
}

func (p *flowSSA) removeTrivialPhis() {
	alias := make([]int, len(p.values))
	for i := range alias {
		alias[i] = i
	}
	root := func(id int) int {
		for id != alias[id] {
			id = alias[id]
		}
		return id
	}
	for changed := true; changed; {
		changed = false
		for id, v := range p.values {
			if v.op != "phi" || alias[id] != id {
				continue
			}
			common := -1
			for _, in := range v.inputs {
				in = root(in)
				if in != id {
					if common >= 0 && common != in {
						common = -1
						break
					}
					common = in
				}
			}
			if common >= 0 {
				alias[id] = common
				changed = true
			}
		}
	}
	for i := range p.values {
		v := &p.values[i]
		if v.x >= 0 {
			v.x = root(v.x)
		}
		if v.y >= 0 {
			v.y = root(v.y)
		}
		for j, in := range v.inputs {
			v.inputs[j] = root(in)
		}
	}
	for i := range p.blocks {
		block := &p.blocks[i]
		keep := block.phis[:0]
		for _, id := range block.phis {
			if alias[id] == id {
				keep = append(keep, id)
			}
		}
		block.phis = keep
		if block.term != "" && block.value >= 0 {
			block.value = root(block.value)
		}
	}
}
