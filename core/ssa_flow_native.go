package core

import (
	"math/bits"
	"slices"
	"sort"
)

type flowSet [maxFlowValues / 64]uint64

func (s *flowSet) add(id int) {
	s[id/64] |= uint64(1) << (id % 64)
}

func (s flowSet) has(id int) bool {
	return s[id/64]&(uint64(1)<<(id%64)) != 0
}

func (s *flowSet) union(other flowSet) {
	for i, word := range other {
		s[i] |= word
	}
}

func (s flowSet) each(visit func(int)) {
	for i, word := range s {
		for word != 0 {
			visit(i*64 + bits.TrailingZeros64(word))
			word &= word - 1
		}
	}
}

func (p flowSSA) edgeInputs(from, to int, visit func(dst, src int)) {
	target := p.blocks[to]
	for i, pred := range target.predecessors {
		if pred == from {
			for _, phi := range target.phis {
				visit(phi, p.values[phi].inputs[i])
			}
			return
		}
	}
}

type flowLiveness struct {
	use, def, in, out flowSet
}

func flowShiftDivisor(n int64) bool { return n > 0 && n <= 1<<30 && n&(n-1) == 0 }

// A right operand encoded by the operation needs no register or stack slot.
func (p flowSSA) immediate(v flowValue) (int64, bool) {
	if v.y < 0 || p.values[v.y].op != "literal" {
		return 0, false
	}
	n := p.values[v.y].number
	if v.op == "/" || v.op == "%" {
		return n, flowShiftDivisor(n)
	}
	return n, fitsSignedImmediate32(n)
}

func (p flowSSA) liveness() []flowLiveness {
	live := make([]flowLiveness, len(p.blocks))
	for i, block := range p.blocks {
		state := &live[i]
		use := func(id int) {
			if id >= 0 && !state.def.has(id) {
				state.use.add(id)
			}
		}
		for _, id := range block.phis {
			state.def.add(id)
		}
		for _, id := range block.code {
			v := p.values[id]
			use(v.x)
			if _, immediate := p.immediate(v); !immediate {
				use(v.y)
			}
			state.def.add(id)
		}
		if block.term != "" && block.term != "branch-flags" {
			use(block.value)
		}
	}
	for changed := true; changed; {
		changed = false
		for i := len(p.blocks) - 1; i >= 0; i-- {
			block := p.blocks[i]
			if block.term == "" {
				continue
			}
			var out flowSet
			for _, next := range []int{block.yes, block.no} {
				if next >= 0 {
					out.union(live[next].in)
					p.edgeInputs(i, next, func(_, src int) { out.add(src) })
				}
			}
			state := &live[i]
			in := state.use
			for j, word := range out {
				in[j] |= word &^ state.def[j]
			}
			changed = changed || in != state.in
			state.in, state.out = in, out
		}
	}
	return live
}

type flowInterval struct {
	id, start, end int
	used           bool
}

// Prefer the true edge for layout where possible. Numbering intervals in the
// same order makes forward uses and loop backedges cheaper to allocate.
func (p flowSSA) order() []int {
	seen := make([]bool, len(p.blocks))
	order := make([]int, 0, len(p.blocks))
	var visit func(int)
	visit = func(id int) {
		if id < 0 || seen[id] {
			return
		}
		seen[id] = true
		block := p.blocks[id]
		visit(block.no)
		visit(block.yes)
		order = append(order, id)
	}
	visit(0)
	slices.Reverse(order)
	return order
}

// Use live-in and live-out sets to keep values across forward and backward
// edges. Covering holes conservatively permits a single location per value.
func allocateFlow(p flowSSA, order []int) ([]ssaLocation, int) {
	live := p.liveness()
	intervals := make([]flowInterval, len(p.values))
	for id := range intervals {
		intervals[id] = flowInterval{id: id, start: int(^uint(0) >> 1)}
	}
	position := 0
	at := func(id int, use bool) {
		if id >= 0 {
			r := &intervals[id]
			r.start = min(r.start, position)
			r.end = max(r.end, position)
			r.used = r.used || use
		}
	}
	for _, i := range order {
		block := p.blocks[i]
		position++
		for _, id := range block.phis {
			at(id, false)
		}
		live[i].in.each(func(id int) { at(id, true) })
		for _, id := range block.code {
			position++
			v := p.values[id]
			at(v.x, true)
			if _, immediate := p.immediate(v); !immediate {
				at(v.y, true)
			}
			at(id, false)
		}
		position++
		if block.term != "" && block.term != "branch-flags" {
			at(block.value, true)
		}
		live[i].out.each(func(id int) { at(id, true) })
	}
	sort.Slice(intervals, func(i, j int) bool {
		a, b := intervals[i], intervals[j]
		return a.start < b.start || a.start == b.start && a.id < b.id
	})
	locations := make([]ssaLocation, len(p.values))
	var active []flowInterval
	var freeSpills []int
	frame := 0
	for _, r := range intervals {
		if !r.used {
			continue
		}
		keep := active[:0]
		var registers [4]bool
		for _, a := range active {
			loc := locations[a.id]
			if a.end < r.start {
				if loc.spill != 0 {
					freeSpills = append(freeSpills, loc.spill)
				}
			} else {
				if loc.spill == 0 {
					registers[loc.reg-8] = true
				}
				keep = append(keep, a)
			}
		}
		active = append(keep, r)
		loc := ssaLocation{present: true}
		for i, used := range registers {
			if !used {
				loc.reg = byte(i + 8)
				break
			}
		}
		if loc.reg == 0 {
			if n := len(freeSpills); n > 0 {
				loc.spill = freeSpills[n-1]
				freeSpills = freeSpills[:n-1]
			} else {
				frame += 8
				loc.spill = -frame
			}
		}
		locations[r.id] = loc
	}
	return locations, frame
}

func (c *compiler) emitNativeFlowSSA(f fn) bool {
	a := c.a
	if !a.turbo || a.trace {
		return false
	}
	p, ok := buildFlowSSA(f)
	if !ok {
		return false
	}
	order := p.order()
	locations, frame := allocateFlow(p, order)
	labels := make([]string, len(p.blocks))
	for i := range labels {
		labels[i] = a.fresh()
	}
	load := func(id int, register byte) {
		loc := locations[id]
		if !loc.present {
			die("internal flow SSA operand has no location")
		}
		if loc.spill == 0 {
			a.regMove(register, loc.reg)
		} else if register == 1 {
			a.loadRCX(loc.spill)
		} else {
			a.load(loc.spill)
		}
	}
	save := func(id int) {
		loc := locations[id]
		if loc.present {
			if loc.spill == 0 {
				a.regMove(loc.reg, 0)
			} else {
				a.store(loc.spill)
			}
		}
	}
	edgeCopies := func(from, to int) [][2]int {
		var copies [][2]int
		p.edgeInputs(from, to, func(dst, src int) {
			if locations[dst].present && locations[dst] != locations[src] {
				copies = append(copies, [2]int{dst, src})
			}
		})
		return copies
	}
	edge := func(copies [][2]int, to, next int) {
		// Give a move priority when it cannot overwrite another source.
		// Use the stack for the remaining cycles only.
		for len(copies) > 0 {
			safe := -1
			for i, pair := range copies {
				needed := false
				for _, other := range copies {
					if locations[pair[0]] == locations[other[1]] {
						needed = true
						break
					}
				}
				if !needed {
					safe = i
					break
				}
			}
			if safe < 0 {
				break
			}
			pair := copies[safe]
			load(pair[1], 0)
			save(pair[0])
			copies = append(copies[:safe], copies[safe+1:]...)
		}
		if len(copies) > 0 {
			for _, pair := range copies {
				load(pair[1], 0)
				a.emit(0x50)
			}
			for i := len(copies) - 1; i >= 0; i-- {
				a.emit(0x58)
				save(copies[i][0])
			}
		}
		if to != next {
			a.jmp(labels[to])
		}
	}
	a.emit(0x55, 0x48, 0x89, 0xe5)
	if frame > 0 {
		a.emit(0x48, 0x81, 0xec)
		a.u32(frame)
	}
	for layout, i := range order {
		block := p.blocks[i]
		condition := byte(0x84)
		next := -1
		if layout+1 < len(order) {
			next = order[layout+1]
		}
		a.label(labels[i])
		for _, id := range block.code {
			v := p.values[id]
			switch v.op {
			case "literal":
				if !locations[id].present {
					continue
				}
				a.imm(v.number)
			case "param":
				if !locations[id].present {
					continue
				}
				a.load(16 + 8*(len(f.args)-1-int(v.number)))
			case "un-":
				load(v.x, 0)
				a.emit(0x48, 0xf7, 0xd8)
				a.j(0x80, "err")
			case "un!":
				load(v.x, 0)
				a.emit(0x48, 0x83, 0xf0, 1)
			case "mask==", "mask!=":
				load(v.x, 0)
				a.emit(0xa9)
				a.u32(int(v.number))
				condition = (compareCode[v.op[4:]] - 0x10) ^ 1
				continue
			default:
				load(v.x, 0)
				n, immediate := p.immediate(v)
				switch v.op {
				case "+", "-", "*":
					if immediate {
						a.checkedArithmeticImmediate(v.op, n)
					} else {
						load(v.y, 1)
						a.regArith(v.op, 0, 1)
					}
				case "/", "%":
					if immediate {
						a.divPowerOfTwo(n, v.op == "%")
					} else {
						load(v.y, 1)
						a.checkedSSADivision(v.op == "%")
					}
				default:
					if immediate {
						a.cmpImm(n)
					} else {
						load(v.y, 1)
						a.emit(0x48, 0x39, 0xc8)
					}
					cc := compareCode[v.op]
					if block.term == "branch-flags" && block.value == id {
						condition = (cc - 0x10) ^ 1
						continue
					}
					a.emit(0x0f, cc, 0xc0, 0x0f, 0xb6, 0xc0)
				}
			}
			save(id)
		}
		switch block.term {
		case "return":
			load(block.value, 0)
			a.emit(0xc9, 0xc3)
		case "jump":
			edge(edgeCopies(i, block.yes), block.yes, next)
		case "branch", "branch-flags":
			if block.term == "branch" {
				load(block.value, 0)
				a.emit(0x48, 0x85, 0xc0)
			}
			yes, no := edgeCopies(i, block.yes), edgeCopies(i, block.no)
			if len(no) == 0 && (block.yes == next || len(yes) != 0 || block.no != next) {
				a.j(condition, labels[block.no])
				edge(yes, block.yes, next)
				continue
			}
			if len(yes) == 0 {
				a.j(condition^1, labels[block.yes])
				edge(no, block.no, next)
				continue
			}
			other := a.fresh()
			if block.yes == next {
				a.j(condition^1, other)
				edge(no, block.no, -1)
				a.label(other)
				edge(yes, block.yes, next)
			} else {
				a.j(condition, other)
				edge(yes, block.yes, -1)
				a.label(other)
				edge(no, block.no, next)
			}
		}
	}
	a.flowSSAFunctions++
	return true
}
