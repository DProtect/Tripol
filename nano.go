package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var compareCode = map[string]byte{"==": 0x94, "!=": 0x95, "<": 0x9c, ">": 0x9f, "<=": 0x9e, ">=": 0x9d}

func (a *asm) cmpImm(n int64) {
	if n >= -128 && n <= 127 {
		a.emit(0x48, 0x83, 0xf8, byte(n))
	} else {
		a.emit(0x48, 0x3d)
		a.u32(int(n))
	}
}
func (c *compiler) branchFalse(e *expr, label string) {
	a := c.a
	if a.opt && e.op == "bin" && e.args[1].op == "int" && e.args[1].n <= math.MaxInt32 {
		if cc, ok := compareCode[e.s]; ok {
			c.require(c.expression(e.args[0]), "int")
			a.cmpImm(e.args[1].n)
			a.j((cc-0x10)^1, label)
			return
		}
	}
	c.require(c.expression(e), "bool")
	a.emit(0x48, 0x85, 0xc0)
	a.j(0x84, label)
}

// Fold only pure, well-typed, non-trapping constants. Validation is separate.
func fold(e *expr) (string, int64, bool) {
	if e == nil {
		return "", 0, false
	}
	if e.op == "int" || e.op == "bool" {
		return e.op, e.n, true
	}
	if e.op != "un" && e.op != "bin" {
		return "", 0, false
	}
	t, x, ok := fold(e.args[0])
	if !ok {
		return "", 0, false
	}
	if e.op == "un" {
		if e.s == "!" && t == "bool" {
			return t, x ^ 1, true
		}
		if e.s == "-" && t == "int" && x != math.MinInt64 {
			return t, -x, true
		}
		return "", 0, false
	}
	u, y, ok := fold(e.args[1])
	if !ok || u != t {
		return "", 0, false
	}
	b := func(v bool) (string, int64, bool) {
		if v {
			return "bool", 1, true
		}
		return "bool", 0, true
	}
	switch e.s {
	case "==":
		return b(x == y)
	case "!=":
		return b(x != y)
	case "&&":
		if t == "bool" {
			return b(x != 0 && y != 0)
		}
	case "||":
		if t == "bool" {
			return b(x != 0 || y != 0)
		}
	}
	if t != "int" {
		return "", 0, false
	}
	switch e.s {
	case "<":
		return b(x < y)
	case ">":
		return b(x > y)
	case "<=":
		return b(x <= y)
	case ">=":
		return b(x >= y)
	case "+":
		if (y > 0 && x > math.MaxInt64-y) || (y < 0 && x < math.MinInt64-y) {
			break
		}
		return t, x + y, true
	case "-":
		if (y < 0 && x > math.MaxInt64+y) || (y > 0 && x < math.MinInt64+y) {
			break
		}
		return t, x - y, true
	case "*":
		if x == 0 || y == 0 {
			return t, 0, true
		}
		if (x == math.MinInt64 && y == -1) || (y == math.MinInt64 && x == -1) {
			break
		}
		z := x * y
		if z/y == x {
			return t, z, true
		}
	case "/", "%":
		if y != 0 && !(x == math.MinInt64 && y == -1) {
			if e.s == "/" {
				return t, x / y, true
			}
			return t, x % y, true
		}
	}
	return "", 0, false
}
func reachable(fs []fn, ss []stmt) map[string]bool {
	by := map[string]fn{}
	for _, f := range fs {
		by[f.name] = f
	}
	live := map[string]bool{}
	var visit func([]stmt)
	var ex func(*expr)
	ex = func(e *expr) {
		if e == nil {
			return
		}
		if e.op == "call" && !live[e.s] {
			live[e.s] = true
			visit(by[e.s].body)
		}
		for _, a := range e.args {
			ex(a)
		}
	}
	visit = func(ss []stmt) {
		for _, s := range ss {
			ex(s.e)
			ex(s.bound)
			visit(s.yes)
			visit(s.no)
		}
	}
	visit(ss)
	return live
}
func (c *compiler) alloc(typ string, mut bool) variable {
	c.slots++
	if c.slots > 400 {
		c.err("too many local variables (maximum 400)")
	}
	if c.slots > c.peak {
		c.peak = c.slots
	}
	return variable{typ, -8 * c.slots, mut}
}
func (c *compiler) rangeLoop(s stmt) {
	a := c.a
	if _, ok := c.lookup(s.name); ok {
		c.err("duplicate or shadowed variable " + s.name)
	}
	old := c.slots
	c.require(c.expression(s.e), "int")
	v := c.alloc("int", false)
	a.store(v.off)
	c.require(c.expression(s.bound), "int")
	bound := c.alloc("int", false)
	a.store(bound.off)
	c.scopes = append(c.scopes, map[string]variable{s.name: v})
	defer func() {
		c.scopes = c.scopes[:len(c.scopes)-1]
		if a.opt {
			c.slots = old
		}
	}()
	start, inc, end := a.fresh(), a.fresh(), a.fresh()
	a.label(start)
	a.load(bound.off)
	a.emit(0x50)
	a.load(v.off)
	a.emit(0x59, 0x48, 0x39, 0xc8)
	a.j(0x8d, end)
	c.loops = append(c.loops, [2]string{end, inc})
	c.block(s.yes, true)
	c.loops = c.loops[:len(c.loops)-1]
	a.label(inc)
	a.load(v.off)
	a.emit(0x48, 0x83, 0xc0, 1)
	a.j(0x80, "err")
	a.store(v.off)
	a.jmp(start)
	a.label(end)
}
func newProject(dir string) {
	if err := os.Mkdir(dir, 0755); err != nil {
		die("cannot create project: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.tri"), []byte("show \"Hello, TRIPOL!\"\n"), 0644); err != nil {
		die("%v", err)
	}
	fmt.Printf("Created %s\ncd %s\ntripol run main.tri\n", dir, dir)
}

// Fixed-size local int arrays: inline stack storage, no heap or GC.
func arrayLen(v variable) int { n, _ := strconv.Atoi(strings.TrimPrefix(v.typ, "array:")); return n }
func (c *compiler) arrayVar(name string) variable {
	v, ok := c.lookup(name)
	if !ok || !strings.HasPrefix(v.typ, "array:") {
		c.err("expected array: " + name)
	}
	return v
}
func (c *compiler) arrayIndex(v variable, e *expr) {
	c.require(c.expression(e), "int")
	a := c.a
	a.emit(0x48, 0x3d)
	a.u32(arrayLen(v))
	a.j(0x83, "bounds") /* unsigned >= also catches negatives */
}
func (c *compiler) arrayDeclare(s stmt) {
	if _, ok := c.lookup(s.name); ok {
		c.err("array declaration requires a new local: " + s.name)
	}
	if s.typ != "" {
		c.err("array type is inferred")
	}
	n := len(s.e.args)
	if n == 0 || n > 400 {
		c.err("array must have 1 to 400 elements")
	}
	for i := 0; i < n; i++ {
		c.alloc("int", true)
	}
	v := variable{fmt.Sprintf("array:%d", n), -8 * c.slots, s.mut || s.op == "assign"}
	for i, e := range s.e.args {
		c.require(c.expression(e), "int")
		c.a.store(v.off + 8*i)
	}
	c.scopes[len(c.scopes)-1][s.name] = v
}

func (a *asm) needed(label string) bool {
	for _, f := range a.fixes {
		if f.label == label {
			return true
		}
	}
	return false
}
