package core

import (
	"math/bits"
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
			use(v.y)
			state.def.add(id)
		}
		if block.term != "" {
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

// Use live-in and live-out sets to keep values across forward and backward
// edges. Covering holes conservatively permits a single location per value.
func allocateFlow(p flowSSA) ([]ssaLocation, int) {
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
	for i, block := range p.blocks {
		position++
		for _, id := range block.phis {
			at(id, false)
		}
		live[i].in.each(func(id int) { at(id, true) })
		for _, id := range block.code {
			position++
			v := p.values[id]
			at(v.x, true)
			at(v.y, true)
			at(id, false)
		}
		position++
		if block.term != "" {
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
	locations, frame := allocateFlow(p)
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
	edge := func(from, to, next int) {
		// Stack temporaries break cycles when phi destinations overlap
		// source locations on an edge. Skip identity and unused moves.
		var copies [][2]int
		p.edgeInputs(from, to, func(dst, src int) {
			if locations[dst].present && locations[dst] != locations[src] {
				copies = append(copies, [2]int{dst, src})
			}
		})
		if len(copies) == 1 {
			load(copies[0][1], 0)
			save(copies[0][0])
		} else {
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
	for i, block := range p.blocks {
		if block.term == "" {
			continue
		}
		next := i + 1
		for next < len(p.blocks) && p.blocks[next].term == "" {
			next++
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
			default:
				load(v.x, 0)
				switch v.op {
				case "+", "-", "*":
					right := p.values[v.y]
					if right.op != "literal" || !a.checkedArithmeticImmediate(v.op, right.number) {
						load(v.y, 1)
						a.regArith(v.op, 0, 1)
					}
				case "/", "%":
					load(v.y, 1)
					a.checkedSSADivision(v.op == "%")
				default:
					load(v.y, 1)
					a.emit(0x48, 0x39, 0xc8, 0x0f, compareCode[v.op], 0xc0, 0x0f, 0xb6, 0xc0)
				}
			}
			save(id)
		}
		switch block.term {
		case "return":
			load(block.value, 0)
			a.emit(0xc9, 0xc3)
		case "jump":
			edge(i, block.yes, next)
		case "branch":
			load(block.value, 0)
			first, second, condition := block.yes, block.no, byte(0x84)
			if block.yes == next {
				first, second, condition = block.no, block.yes, 0x85
			}
			other := a.fresh()
			a.emit(0x48, 0x85, 0xc0)
			a.j(condition, other)
			edge(i, first, -1)
			a.label(other)
			edge(i, second, next)
		}
	}
	a.flowSSAFunctions++
	return true
}
