package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const version = "0.6.0"

type token struct {
	s, k string
	line int
}

var lexRE = regexp.MustCompile(`^(?:([ \t\r]+)|(//[^\n]*|#[^\n]*)|(\n)|("(?:[^"\\\n]|\\.)*")|([0-9]+(?:\.[0-9]+)?)|([A-Za-z_][A-Za-z_0-9]*)|(\+=|-=|\*=|/=|%=|==|!=|<=|>=|&&|\|\||->|[+*/%<>=!(){}\[\],:;\-]))`)

func die(f string, a ...any) { panic(fmt.Errorf(f, a...)) }
func lex(s string) []token {
	var ts []token
	line := 1
	for len(s) > 0 {
		m := lexRE.FindStringSubmatch(s)
		if m == nil {
			die("line %d: unexpected character %q", line, s[0])
		}
		k := 0
		for i := 1; i < len(m); i++ {
			if m[i] != "" {
				k = i
				break
			}
		}
		if k > 2 {
			ts = append(ts, token{m[0], []string{"", "", "", "nl", "str", "int", "id", "op"}[k], line})
		}
		if k == 3 {
			line++
		}
		s = s[len(m[0]):]
	}
	return append(ts, token{"<eof>", "eof", line})
}

type expr struct {
	op, s string
	n     int64
	args  []*expr
}
type stmt struct {
	op, name, typ string
	e, bound      *expr
	yes, no       []stmt
	line          int
	mut           bool
}
type param struct{ name, typ string }
type fn struct {
	name, ret string
	args      []param
	body      []stmt
}
type parser struct {
	ts    []token
	i     int
	depth int
}

func (p *parser) t() token     { return p.ts[p.i] }
func (p *parser) err(s string) { die("line %d: %s (got %q)", p.t().line, s, p.t().s) }
func (p *parser) take() token {
	t := p.t()
	if t.k == "eof" {
		p.err("unexpected end of input")
	}
	p.i++
	return t
}
func (p *parser) eat(s string) bool {
	if p.t().s == s {
		p.take()
		return true
	}
	return false
}
func (p *parser) need(s string) {
	if !p.eat(s) {
		p.err("expected " + s)
	}
}
func (p *parser) skip() {
	for p.eat("\n") || p.eat(";") {
	}
}
func (p *parser) end() {
	switch p.t().s {
	case "\n", ";", "}", "<eof>":
	default:
		p.err("expected statement separator")
	}
	p.skip()
}
func (p *parser) name() string {
	t := p.t()
	if t.k != "id" || strings.Contains("|import|func|return|show|let|const|if|else|while|for|in|range|repeat|and|or|not|break|continue|true|false|int|bool|string|float|ints|void|", "|"+t.s+"|") {
		p.err("expected identifier")
	}
	return p.take().s
}
func (p *parser) typ() string {
	switch p.t().s {
	case "int", "bool", "string", "float", "ints":
		return p.take().s
	}
	p.err("expected type")
	return ""
}
func (p *parser) block() []stmt {
	p.skip()
	p.need("{")
	p.skip()
	var ss []stmt
	p.depth++
	if p.depth > 256 {
		p.err("nesting limit")
	}
	defer func() { p.depth-- }()
	for !p.eat("}") {
		ss = append(ss, p.statement())
		p.skip()
	}
	return ss
}
func (p *parser) statement() stmt {
	t := p.t()
	s := stmt{line: t.line}
	switch t.s {
	case "break", "continue":
		s.op = p.take().s
	case "repeat":
		p.take()
		s.op = "for"
		s.name = fmt.Sprintf("@repeat%d", p.i)
		s.e = &expr{op: "int"}
		s.bound = p.expression(1)
		s.yes = p.block()
		return s
	case "for":
		p.take()
		s.op = "for"
		s.name = p.name()
		p.need("in")
		if p.eat("range") {
			p.need("(")
			first := p.expression(1)
			if p.eat(",") {
				s.e = first
				s.bound = p.expression(1)
			} else {
				s.e = &expr{op: "int"}
				s.bound = first
			}
			p.need(")")
			s.yes = p.block()
		} else {
			array := p.name()
			item := s.name
			s.name = fmt.Sprintf("@each%d", p.i)
			s.e = &expr{op: "int"}
			s.bound = &expr{op: "call", s: "@array_len", args: []*expr{{op: "var", s: array}}}
			body := p.block()
			s.yes = append([]stmt{{op: "decl", name: item, line: t.line, e: &expr{op: "index", s: array, args: []*expr{{op: "var", s: s.name}}}}}, body...)
		}
		return s
	case "show", "return":
		s.op = p.take().s
		s.e = p.expression(1)
		if s.op == "show" && p.eat(",") {
			args := []*expr{s.e}
			for {
				args = append(args, p.expression(1))
				if !p.eat(",") {
					break
				}
			}
			s.op = "showmany"
			s.e = &expr{op: "list", args: args}
		}
	case "if", "while":
		s.op = p.take().s
		s.e = p.expression(1)
		s.yes = p.block()
		if s.op == "if" {
			p.skip()
			if p.eat("else") {
				p.skip()
				if p.t().s == "if" {
					p.depth++
					if p.depth > 256 {
						p.err("nesting limit")
					}
					s.no = []stmt{p.statement()}
					p.depth--
				} else {
					s.no = p.block()
				}
			}
		}
		return s
	case "let", "const":
		s.op = "decl"
		s.mut = p.take().s == "let"
		s.name = p.name()
		if p.eat(":") {
			s.typ = p.typ()
		}
		p.need("=")
		s.e = p.expression(1)
	default:
		if t.k == "id" && strings.Contains("|=|+=|-=|*=|/=|%=|", "|"+p.ts[p.i+1].s+"|") && p.ts[p.i+1].s != "|" {
			s.op = "assign"
			s.name = p.name()
			op := p.take().s
			s.e = p.expression(1)
			if op != "=" {
				s.e = &expr{op: "bin", s: op[:1], args: []*expr{{op: "var", s: s.name}, s.e}}
			}
		} else {
			s.op = "expr"
			s.e = p.expression(1)
			if s.e.op == "index" && p.eat("=") {
				s.op = "indexassign"
				s.name = s.e.s
				s.bound = s.e.args[0]
				s.e = p.expression(1)
			}
		}
	}
	p.end()
	return s
}

var prec = map[string]int{"||": 1, "or": 1, "&&": 2, "and": 2, "==": 3, "!=": 3, "<": 4, ">": 4, "<=": 4, ">=": 4, "+": 5, "-": 5, "*": 6, "/": 6, "%": 6}

func (p *parser) expression(min int) *expr {
	p.depth++
	if p.depth > 256 {
		p.err("nesting limit")
	}
	defer func() { p.depth-- }()
	t := p.take()
	e := &expr{s: t.s}
	switch {
	case t.s == "-" || t.s == "!" || t.s == "not":
		if t.s == "not" {
			e.s = "!"
		}
		e.op = "un"
		if t.s == "not" {
			e.args = []*expr{p.expression(3)}
		} else {
			e.args = []*expr{p.expression(7)}
		}
	case t.s == "[":
		e.op = "array"
		if !p.eat("]") {
			for {
				e.args = append(e.args, p.expression(1))
				if p.eat("]") {
					break
				}
				p.need(",")
			}
		}
	case t.s == "(":
		e = p.expression(1)
		p.need(")")
	case t.k == "int" && strings.Contains(t.s, "."):
		e.op = "float"
		f, err := strconv.ParseFloat(t.s, 64)
		if err != nil {
			p.err("invalid float literal")
		}
		e.n = int64(math.Float64bits(f))
	case t.k == "int":
		e.op = "int"
		n, err := strconv.ParseInt(t.s, 10, 64)
		if err != nil {
			die("line %d: integer outside signed 64-bit range", t.line)
		}
		e.n = n
	case t.k == "str":
		e.op = "string"
		if json.Unmarshal([]byte(t.s), &e.s) != nil || strings.ContainsRune(e.s, 0) {
			die("line %d: invalid string", t.line)
		}
	case t.s == "true" || t.s == "false":
		e.op = "bool"
		if t.s == "true" {
			e.n = 1
		}
	case t.k == "id":
		e.op = "var"
		if p.eat("(") {
			e.op = "call"
			if !p.eat(")") {
				for {
					e.args = append(e.args, p.expression(1))
					if p.eat(")") {
						break
					}
					p.need(",")
				}
			}
		}
	default:
		die("line %d: expected expression", t.line)
	}
	if e.op == "var" && p.eat("[") {
		e.op = "index"
		e.args = []*expr{p.expression(1)}
		p.need("]")
	}
	for prec[p.t().s] >= min {
		op := p.take().s
		priority := prec[op]
		if op == "and" {
			op = "&&"
		}
		if op == "or" {
			op = "||"
		}
		e = &expr{op: "bin", s: op, args: []*expr{e, p.expression(priority + 1)}}
	}
	return e
}
func parseUnit(s string) ([]string, []fn, []stmt) {
	p := parser{ts: lex(strings.TrimPrefix(s, "\ufeff"))}
	var fs []fn
	var ss []stmt
	var imports []string
	p.skip()
	for p.eat("import") {
		if p.t().k != "str" {
			p.err("expected quoted local module path")
		}
		var name string
		if json.Unmarshal([]byte(p.take().s), &name) != nil {
			p.err("invalid import path")
		}
		imports = append(imports, name)
		p.end()
	}
	for p.t().k != "eof" {
		if p.eat("func") {
			f := fn{name: p.name(), ret: "int"}
			p.need("(")
			if !p.eat(")") {
				for {
					a := param{p.name(), "int"}
					if p.eat(":") {
						a.typ = p.typ()
					}
					f.args = append(f.args, a)
					if p.eat(")") {
						break
					}
					p.need(",")
				}
			}
			if p.eat("->") {
				f.ret = p.typ()
			}
			f.body = p.block()
			fs = append(fs, f)
		} else {
			ss = append(ss, p.statement())
		}
		p.skip()
	}
	return imports, fs, ss
}
func parse(s string) ([]fn, []stmt) {
	imports, fs, ss := parseUnit(s)
	if len(imports) != 0 {
		die("imports require a source file")
	}
	return fs, ss
}

// Direct x86-64 assembler with relative relocations. No external toolchain.
type fix struct {
	at    int
	label string
}
type asm struct {
	b      []byte
	labels map[string]int
	fixes  []fix
	serial int
	data   map[string]string
	target string
	opt    bool
}

func (a *asm) emit(b ...byte) { a.b = append(a.b, b...) }
func (a *asm) u32(n int)      { a.b = binary.LittleEndian.AppendUint32(a.b, uint32(n)) }
func (a *asm) imm(n int64) {
	if a.opt && n >= 0 && n <= 0xffffffff {
		a.emit(0xb8)
		a.u32(int(n))
		return
	}
	if a.opt && n >= -2147483648 && n <= 2147483647 {
		a.emit(0x48, 0xc7, 0xc0)
		a.u32(int(n))
		return
	}
	a.emit(0x48, 0xb8)
	a.b = binary.LittleEndian.AppendUint64(a.b, uint64(n))
}
func (a *asm) label(s string) {
	if _, ok := a.labels[s]; ok {
		die("internal: duplicate label %s", s)
	}
	a.labels[s] = len(a.b)
}
func (a *asm) fresh() string { a.serial++; return fmt.Sprintf("@%d", a.serial) }
func (a *asm) rel(s string, b ...byte) {
	a.emit(b...)
	a.fixes = append(a.fixes, fix{len(a.b), s})
	a.u32(0)
}
func (a *asm) call(s string)      { a.rel(s, 0xe8) }
func (a *asm) jmp(s string)       { a.rel(s, 0xe9) }
func (a *asm) j(c byte, s string) { a.rel(s, 0x0f, c) }
func (a *asm) str(s string) {
	l, ok := a.data[s]
	if !ok {
		l = a.fresh()
		a.data[s] = l
	}
	a.rel(l, 0x48, 0x8d, 0x05)
}
func (a *asm) load(off int) {
	if a.opt && off >= -128 && off <= 127 {
		a.emit(0x48, 0x8b, 0x45, byte(off))
		return
	}
	a.emit(0x48, 0x8b, 0x85)
	a.u32(off)
}
func (a *asm) store(off int) {
	if a.opt && off >= -128 && off <= 127 {
		a.emit(0x48, 0x89, 0x45, byte(off))
		return
	}
	a.emit(0x48, 0x89, 0x85)
	a.u32(off)
}
func (a *asm) finish() {
	for _, f := range a.fixes {
		v, ok := a.labels[f.label]
		if !ok {
			die("internal missing label %s", f.label)
		}
		binary.LittleEndian.PutUint32(a.b[f.at:], uint32(v-f.at-4))
	}
}

type variable struct {
	typ string
	off int
	mut bool
}
type compiler struct {
	a           *asm
	fs          map[string]fn
	scopes      []map[string]variable
	slots, line int
	peak        int
	ret         string
	current     fn
	restart     string
	loops       [][2]string
}

func (c *compiler) err(s string) { die("line %d: %s", c.line, s) }
func (c *compiler) lookup(n string) (variable, bool) {
	for i := len(c.scopes) - 1; i >= 0; i-- {
		if v, ok := c.scopes[i][n]; ok {
			return v, true
		}
	}
	return variable{}, false
}
func (c *compiler) require(a, b string) {
	if a != b {
		c.err("expected " + b + ", got " + a)
	}
}
func (c *compiler) expression(e *expr) string {
	a := c.a
	if a.opt {
		if k, n, ok := fold(e); ok && e.op != "int" && e.op != "bool" {
			a.imm(n)
			return k
		}
	}
	switch e.op {
	case "int", "bool", "float":
		a.imm(e.n)
		return e.op
	case "string":
		a.str(e.s)
		return "string"
	case "var":
		v, ok := c.lookup(e.s)
		if !ok {
			c.err("unknown variable " + e.s)
		}
		if strings.HasPrefix(v.typ, "array:") {
			c.err("array requires an index")
		}
		a.load(v.off)
		return v.typ
	case "array":
		c.err("array literal only allowed in a local declaration")
	case "index":
		v := c.arrayVar(e.s)
		c.arrayIndex(v, e.args[0])
		if v.typ == "ints" {
			a.emit(0x48, 0x8b, 0x44, 0xc1, 0x08)
		} else {
			a.emit(0x48, 0x8b, 0x84, 0xc5)
			a.u32(v.off)
		}
		return "int"
	case "call":
		f, ok := c.fs[e.s]
		if !ok {
			if t, handled := c.builtin(e); handled {
				return t
			}
		}
		if !ok {
			c.err("unknown function " + e.s)
		}
		if len(f.args) != len(e.args) {
			c.err("wrong argument count: " + e.s)
		}
		for i, x := range e.args {
			c.require(c.expression(x), f.args[i].typ)
			a.emit(0x50)
		}
		a.call("fn:" + e.s)
		if len(e.args) > 0 {
			if a.opt && len(e.args) <= 15 {
				a.emit(0x48, 0x83, 0xc4, byte(8*len(e.args)))
			} else {
				a.emit(0x48, 0x81, 0xc4)
				a.u32(8 * len(e.args))
			}
		}
		return f.ret
	case "un":
		t := c.expression(e.args[0])
		if e.s == "!" {
			c.require(t, "bool")
			a.emit(0x48, 0x83, 0xf0, 1)
		} else {
			if t == "float" {
				a.emit(0x48, 0xba)
				a.b = binary.LittleEndian.AppendUint64(a.b, 1<<63)
				a.emit(0x48, 0x31, 0xd0)
				return t
			}
			c.require(t, "int")
			a.emit(0x48, 0xf7, 0xd8)
			a.j(0x80, "err")
		}
		return t
	}
	t := c.expression(e.args[0])
	op := e.s
	if op == "&&" || op == "||" {
		c.require(t, "bool")
		end := a.fresh()
		a.emit(0x48, 0x85, 0xc0)
		if op == "&&" {
			a.j(0x84, end)
		} else {
			a.j(0x85, end)
		}
		c.require(c.expression(e.args[1]), "bool")
		a.label(end)
		return "bool"
	}
	if a.opt && t == "int" && e.args[1].op == "int" && e.args[1].n <= 2147483647 {
		n := e.args[1].n
		switch op {
		case "+", "-", "*":
			if op == "*" {
				if n <= 127 {
					a.emit(0x48, 0x6b, 0xc0, byte(n))
				} else {
					a.emit(0x48, 0x69, 0xc0)
					a.u32(int(n))
				}
			} else {
				code := byte(0xc0)
				if op == "-" {
					code = 0xe8
				}
				if n <= 127 {
					a.emit(0x48, 0x83, code, byte(n))
				} else {
					a.emit(0x48, 0x81, code)
					a.u32(int(n))
				}
			}
			a.j(0x80, "err")
			return "int"
		case "==", "!=", "<", ">", "<=", ">=":
			a.cmpImm(n)
			a.emit(0x0f, compareCode[op], 0xc0, 0x0f, 0xb6, 0xc0)
			return "bool"
		}
	}
	if a.opt && e.args[1].op == "var" {
		v, ok := c.lookup(e.args[1].s)
		if !ok {
			c.err("unknown variable " + e.args[1].s)
		}
		c.require(v.typ, t)
		if strings.HasPrefix(v.typ, "array:") {
			c.err("array requires an index")
		}
		a.loadRCX(v.off)
	} else {
		a.emit(0x50)
		u := c.expression(e.args[1])
		c.require(u, t)
		a.emit(0x48, 0x89, 0xc1, 0x58)
	}
	if t == "float" {
		return c.floatOp(op)
	}
	if t == "string" && op == "+" {
		a.call("concat")
		return "string"
	}
	if t == "ints" {
		c.err("array references do not support operators")
	}
	if op == "==" || op == "!=" {
		if t == "string" {
			a.call("streq")
			if op == "!=" {
				a.emit(0x48, 0x83, 0xf0, 1)
			}
			return "bool"
		}
	} else {
		c.require(t, "int")
	}
	switch op {
	case "+":
		a.emit(0x48, 0x01, 0xc8)
		a.j(0x80, "err")
	case "-":
		a.emit(0x48, 0x29, 0xc8)
		a.j(0x80, "err")
	case "*":
		a.emit(0x48, 0x0f, 0xaf, 0xc1)
		a.j(0x80, "err")
	case "/", "%":
		a.emit(0x48, 0x85, 0xc9)
		a.j(0x84, "err")
		ok := a.fresh()
		a.emit(0x48, 0x83, 0xf9, 0xff)
		a.j(0x85, ok)
		a.emit(0x48, 0xba)
		a.b = binary.LittleEndian.AppendUint64(a.b, 1<<63)
		a.emit(0x48, 0x39, 0xd0)
		a.j(0x84, "err")
		a.label(ok)
		a.emit(0x48, 0x99, 0x48, 0xf7, 0xf9)
		if op == "%" {
			a.emit(0x48, 0x89, 0xd0)
		}
	default:
		cc, ok := map[string]byte{"==": 0x94, "!=": 0x95, "<": 0x9c, ">": 0x9f, "<=": 0x9e, ">=": 0x9d}[op]
		if !ok {
			c.err("unknown operator")
		}
		a.emit(0x48, 0x39, 0xc8, 0x0f, cc, 0xc0, 0x48, 0x0f, 0xb6, 0xc0)
		return "bool"
	}
	return "int"
}
func (c *compiler) block(ss []stmt, nested bool) {
	if nested {
		old := c.slots
		c.scopes = append(c.scopes, map[string]variable{})
		defer func() {
			c.scopes = c.scopes[:len(c.scopes)-1]
			if c.a.opt {
				c.slots = old
			}
		}()
	}
	a := c.a
	for _, s := range ss {
		c.line = s.line
		switch s.op {
		case "indexassign":
			v := c.arrayVar(s.name)
			if !v.mut {
				c.err("cannot modify constant " + s.name)
			}
			c.arrayIndex(v, s.bound)
			a.emit(0x50)
			c.require(c.expression(s.e), "int")
			a.emit(0x59)
			if v.typ == "ints" {
				// Preserve value while loading the array reference in rdx.
				a.emit(0x48, 0x8b, 0x95)
				a.u32(v.off)
				a.emit(0x48, 0x89, 0x44, 0xca, 0x08)
			} else {
				a.emit(0x48, 0x89, 0x84, 0xcd)
				a.u32(v.off)
			}
		case "decl", "assign":
			if s.e.op == "array" {
				c.arrayDeclare(s)
				continue
			}
			t := c.expression(s.e)
			v, ok := c.lookup(s.name)
			if s.op == "decl" && ok {
				c.err("duplicate or shadowed variable " + s.name)
			}
			if s.typ != "" {
				c.require(t, s.typ)
			}
			if ok {
				if !v.mut {
					c.err("cannot modify constant " + s.name)
				}
				c.require(t, v.typ)
			} else {
				c.slots++
				if c.slots > c.peak {
					c.peak = c.slots
				}
				if c.slots > 400 {
					c.err("too many local variables (maximum 400)")
				}
				v = variable{t, -8 * c.slots, s.mut || s.op == "assign"}
				c.scopes[len(c.scopes)-1][s.name] = v
			}
			a.store(v.off)
		case "showmany":
			for i, x := range s.e.args {
				if i > 0 {
					a.str(" ")
					a.call("raw_string")
				}
				t := c.expression(x)
				switch t {
				case "int", "string", "bool", "float":
					a.call("raw_" + t)
				default:
					c.err("show requires scalar or string")
				}
			}
			a.str("\n")
			a.call("raw_string")
		case "show":
			t := c.expression(s.e)
			if t == "ints" {
				c.err("show requires a scalar or string")
			}
			a.call("print:" + t)
		case "expr":
			c.expression(s.e)
		case "return":
			if c.ret == "" {
				c.err("return outside function")
			}
			if a.opt && s.e.op == "call" && s.e.s == c.current.name && len(s.e.args) == len(c.current.args) {
				for i, x := range s.e.args {
					c.require(c.expression(x), c.current.args[i].typ)
					a.emit(0x50)
				}
				for i := len(s.e.args) - 1; i >= 0; i-- {
					a.emit(0x58)
					a.store(16 + 8*(len(s.e.args)-1-i))
				}
				a.jmp(c.restart)
			} else {
				c.require(c.expression(s.e), c.ret)
				a.emit(0xc9, 0xc3)
			}
		case "break", "continue":
			if len(c.loops) == 0 {
				c.err(s.op + " outside loop")
			}
			l := c.loops[len(c.loops)-1]
			if s.op == "break" {
				a.jmp(l[0])
			} else {
				a.jmp(l[1])
			}
		case "for":
			c.rangeLoop(s)
		case "if", "while":
			start, end := a.fresh(), a.fresh()
			a.label(start)
			other := end
			if len(s.no) > 0 {
				other = a.fresh()
			}
			c.branchFalse(s.e, other)
			if s.op == "while" {
				c.loops = append(c.loops, [2]string{end, start})
			}
			c.block(s.yes, true)
			if s.op == "while" {
				c.loops = c.loops[:len(c.loops)-1]
			}
			if s.op == "while" {
				a.jmp(start)
			}
			if len(s.no) > 0 {
				a.jmp(end)
				a.label(other)
				c.block(s.no, true)
			}
			a.label(end)
		}
	}
}
func returns(ss []stmt) bool {
	for _, s := range ss {
		if s.op == "return" || (s.op == "if" && returns(s.yes) && returns(s.no)) {
			return true
		}
	}
	return false
}
func (c *compiler) function(f fn, main bool) {
	a := c.a
	c.slots = 0
	c.peak = 0
	c.current = f
	c.loops = nil
	c.ret = f.ret
	c.scopes = []map[string]variable{{}}
	for i, p := range f.args {
		if _, ok := c.scopes[0][p.name]; ok {
			die("duplicate parameter %s", p.name)
		}
		c.scopes[0][p.name] = variable{p.typ, 16 + 8*(len(f.args)-1-i), true}
	}
	a.label("fn:" + f.name)
	a.emit(0x55, 0x48, 0x89, 0xe5)
	at := -1
	if !a.opt || needsFrame(f.body) {
		a.emit(0x48, 0x81, 0xec)
		at = len(a.b)
		a.u32(0)
	}
	c.restart = a.fresh()
	a.label(c.restart)
	c.block(f.body, false)
	if main {
		a.imm(0)
		a.jmp("exit")
	} else if !returns(f.body) {
		die("function %s may finish without return", f.name)
	}
	if at >= 0 {
		binary.LittleEndian.PutUint32(a.b[at:], uint32(c.peak*8))
	} else if c.peak != 0 {
		die("internal frame analysis error")
	}
}

// Runtime helpers use rax for value; write uses rsi=buffer, rdx=length.
func (a *asm) rt() {
	a.libertyRuntime()
	if !a.opt || a.needed("streq") {
		a.label("streq")
		l, n := a.fresh(), a.fresh()
		a.label(l)
		a.emit(0x8a, 0x10, 0x3a, 0x11)
		a.j(0x85, n)
		a.emit(0x84, 0xd2)
		eq := a.fresh()
		a.j(0x84, eq)
		a.emit(0x48, 0xff, 0xc0, 0x48, 0xff, 0xc1)
		a.jmp(l)
		a.label(n)
		a.emit(0x31, 0xc0, 0xc3)
		a.label(eq)
		a.emit(0xb8, 1, 0, 0, 0, 0xc3)
	}
	if !a.opt || a.needed("print:bool") {
		a.label("print:bool")
		a.emit(0x48, 0x85, 0xc0)
		f := a.fresh()
		a.j(0x84, f)
		a.str("true")
		a.jmp("print:string")
		a.label(f)
		a.str("false")
		a.jmp("print:string")
	}
	if !a.opt || a.needed("print:string") {
		a.label("print:string")
		a.emit(0x48, 0x89, 0xc6, 0x31, 0xd2)
		loop, done := a.fresh(), a.fresh()
		a.label(loop)
		a.emit(0x80, 0x3c, 0x16, 0)
		a.j(0x84, done)
		a.emit(0x48, 0xff, 0xc2)
		a.jmp(loop)
		a.label(done)
		a.call("write")
		a.str("\n")
		a.emit(0x48, 0x89, 0xc6, 0xba, 1, 0, 0, 0)
		a.jmp("write")
	}
	if !a.opt || a.needed("print:int") {
		a.label("print:int")
		a.emit(0x55, 0x48, 0x89, 0xe5, 0x48, 0x83, 0xec, 0x30, 0x48, 0x8d, 0x75, 0xff, 0xc6, 0x06, 0x0a, 0x41, 0xb8, 1, 0, 0, 0, 0x45, 0x31, 0xc9, 0x48, 0x85, 0xc0)
		pos := a.fresh()
		a.j(0x89, pos)
		a.emit(0x41, 0xb1, 1, 0x48, 0xf7, 0xd8)
		a.label(pos)
		a.emit(0xb9, 10, 0, 0, 0)
		digits := a.fresh()
		a.label(digits)
		a.emit(0x31, 0xd2, 0x48, 0xf7, 0xf1, 0x80, 0xc2, 0x30, 0x48, 0xff, 0xce, 0x88, 0x16, 0x49, 0xff, 0xc0, 0x48, 0x85, 0xc0)
		a.j(0x85, digits)
		sign := a.fresh()
		a.emit(0x45, 0x84, 0xc9)
		a.j(0x84, sign)
		a.emit(0x48, 0xff, 0xce, 0xc6, 0x06, 0x2d, 0x49, 0xff, 0xc0)
		a.label(sign)
		a.emit(0x4c, 0x89, 0xc2)
		a.call("write")
		a.emit(0xc9, 0xc3)
	}
	if !a.opt || a.needed("bounds") {
		a.label("bounds")
		a.str("TRIPOL runtime error: array index out of bounds\n")
		a.emit(0x48, 0x89, 0xc6, 0xba)
		a.u32(len("TRIPOL runtime error: array index out of bounds\n"))
		a.call("writeerr")
		a.imm(1)
		a.jmp("exit")
	}
	if !a.opt || a.needed("err") {
		a.label("err")
		a.str("TRIPOL runtime error: invalid integer arithmetic\n")
		a.emit(0x48, 0x89, 0xc6, 0xba)
		a.u32(len("TRIPOL runtime error: invalid integer arithmetic\n"))
		a.call("writeerr")
		a.imm(1)
		a.jmp("exit")
	}
	for _, s := range []string{"write", "writeerr"} {
		if a.opt && !a.needed(s) {
			continue
		}
		a.label(s)
		fd := byte(1)
		if s == "writeerr" {
			fd = 2
		}
		if a.target == "linux" {
			a.emit(0xb8, 1, 0, 0, 0, 0xbf, fd, 0, 0, 0, 0x0f, 0x05, 0xc3)
		} else { // Save buffer/length across GetStdHandle, align stack, reserve Win64 shadow space.
			a.emit(0x55, 0x48, 0x89, 0xe5, 0x48, 0x83, 0xe4, 0xf0, 0x48, 0x83, 0xec, 0x50, 0x48, 0x89, 0x74, 0x24, 0x30, 0x48, 0x89, 0x54, 0x24, 0x38, 0xb9)
			a.u32(-10 - int(fd))
			a.rel("iat:0", 0xff, 0x15)
			a.emit(0x48, 0x89, 0xc1, 0x48, 0x8b, 0x54, 0x24, 0x30, 0x44, 0x8b, 0x44, 0x24, 0x38, 0x4c, 0x8d, 0x4c, 0x24, 0x40, 0x48, 0xc7, 0x44, 0x24, 0x20, 0, 0, 0, 0)
			a.rel("iat:1", 0xff, 0x15)
			a.emit(0xc9, 0xc3)
		}
	}
	a.label("exit")
	if a.target == "linux" {
		a.emit(0x48, 0x89, 0xc7, 0xb8, 60, 0, 0, 0, 0x0f, 0x05)
	} else {
		a.emit(0x89, 0xc1, 0x48, 0x83, 0xe4, 0xf0, 0x48, 0x83, 0xec, 0x20)
		a.rel("iat:2", 0xff, 0x15)
	}
	a.emit(0x0f, 0x0b)
}
func align(n, a int) int              { return (n + a - 1) &^ (a - 1) }
func put16(b []byte, o int, v int)    { binary.LittleEndian.PutUint16(b[o:], uint16(v)) }
func put32(b []byte, o int, v int)    { binary.LittleEndian.PutUint32(b[o:], uint32(v)) }
func put64(b []byte, o int, v uint64) { binary.LittleEndian.PutUint64(b[o:], v) }
func elf(a *asm) []byte {
	const base = 0x400000
	const off = 0x100
	arena := a.needed("arena_state")
	arenaOff := align(off+len(a.b), 0x1000)
	if arena {
		a.labels["arena_state"] = arenaOff - off
	}
	a.finish()
	b := make([]byte, off+len(a.b))
	copy(b, []byte{0x7f, 'E', 'L', 'F', 2, 1, 1})
	put16(b, 16, 2)
	put16(b, 18, 62)
	put32(b, 20, 1)
	put64(b, 24, base+off)
	put64(b, 32, 64)
	put16(b, 52, 64)
	put16(b, 54, 56)
	put16(b, 56, 2)
	if arena {
		put16(b, 56, 3)
		put32(b, 176, 1)
		put32(b, 180, 6)
		put64(b, 184, uint64(arenaOff))
		put64(b, 192, uint64(base+arenaOff))
		put64(b, 200, uint64(base+arenaOff))
		put64(b, 216, 24)
		put64(b, 224, 0x1000)
	}
	put32(b, 64, 1)
	put32(b, 68, 5)
	put64(b, 80, base)
	put64(b, 88, base)
	put64(b, 96, uint64(len(b)))
	put64(b, 104, uint64(len(b)))
	put64(b, 112, 0x1000)
	put32(b, 120, 0x6474e551)
	put32(b, 124, 6)
	put64(b, 168, 16)
	copy(b[off:], a.b)
	return b
}
func pe(a *asm) []byte {
	const textRVA = 0x1000
	const header = 0x200
	raw := align(len(a.b), 0x200)
	idRVA := align(textRVA+len(a.b), 0x1000)
	id := make([]byte, 0x200)
	imports := []string{"GetStdHandle", "WriteFile", "ExitProcess"}
	for i, name := range []string{"ReadFile", "VirtualAlloc", "GetLastError", "CreateFileA", "CloseHandle"} {
		if a.needed(fmt.Sprintf("iat:%d", i+3)) {
			imports = append(imports, name)
		} else {
			imports = append(imports, "")
		}
	}
	// Preserve logical IAT labels while packing only used imports.
	count := 0
	for _, name := range imports {
		if name != "" {
			count++
		}
	}
	const ilt = 40
	iat := ilt + 8*(count+1)
	dll := iat + 8*(count+1)
	copy(id[dll:], "KERNEL32.dll\x00")
	cursor := align(dll+len("KERNEL32.dll")+1, 2)
	j := 0
	for i, s := range imports {
		if s == "" {
			continue
		}
		put64(id, ilt+8*j, uint64(idRVA+cursor))
		put64(id, iat+8*j, uint64(idRVA+cursor))
		copy(id[cursor+2:], s+"\x00")
		cursor = align(cursor+2+len(s)+1, 2)
		a.labels[fmt.Sprintf("iat:%d", i)] = idRVA - textRVA + iat + 8*j
		j++
	}
	if a.needed("arena_state") {
		if cursor > 480 {
			die("import table overflow")
		}
		a.labels["arena_state"] = idRVA - textRVA + 480
	}
	put32(id, 0, idRVA+ilt)
	put32(id, 12, idRVA+dll)
	put32(id, 16, idRVA+iat)
	a.finish()
	b := make([]byte, header+raw+len(id))
	copy(b, "MZ")
	put32(b, 0x3c, 0x80)
	copy(b[0x80:], "PE\x00\x00")
	co := 0x84
	put16(b, co, 0x8664)
	put16(b, co+2, 2)
	put16(b, co+16, 240)
	put16(b, co+18, 0x23)
	o := co + 20
	put16(b, o, 0x20b)
	put32(b, o+4, raw)
	put32(b, o+8, len(id))
	put32(b, o+16, textRVA)
	put32(b, o+20, textRVA)
	put64(b, o+24, 0x140000000)
	put32(b, o+32, 0x1000)
	put32(b, o+36, 0x200)
	put16(b, o+40, 6)
	put16(b, o+48, 6)
	put32(b, o+56, align(idRVA+len(id), 0x1000))
	put32(b, o+60, header)
	put16(b, o+68, 3)
	put16(b, o+70, 0x100)
	put64(b, o+72, 1<<20)
	put64(b, o+80, 4096)
	put64(b, o+88, 1<<20)
	put64(b, o+96, 4096)
	put32(b, o+108, 16)
	put32(b, o+120, idRVA)
	put32(b, o+124, 40)
	put32(b, o+208, idRVA+iat)
	put32(b, o+212, 8*(count+1))
	s := o + 240
	copy(b[s:], ".text")
	put32(b, s+8, len(a.b))
	put32(b, s+12, textRVA)
	put32(b, s+16, raw)
	put32(b, s+20, header)
	put32(b, s+36, 0x60000020)
	s += 40
	copy(b[s:], ".idata")
	put32(b, s+8, len(id))
	put32(b, s+12, idRVA)
	put32(b, s+16, len(id))
	put32(b, s+20, header+raw)
	put32(b, s+36, 0xc0000040)
	copy(b[header:], a.b)
	copy(b[header+raw:], id)
	return b
}
func compile(src, target string) []byte { return compileProfile(src, target, "nano") }
func compileProfile(src, target, profile string) []byte {
	fs, ss := parse(src)
	return compileParsed(fs, ss, target, profile)
}
func compileParsed(fs []fn, ss []stmt, target, profile string) []byte {
	// Validate every function before eliminating any unreachable code.
	if profile != "debug" {
		generate(fs, ss, target, false, nil)
	}
	return generate(fs, ss, target, profile != "debug", reachable(fs, ss))
}
func generate(fs []fn, ss []stmt, target string, opt bool, live map[string]bool) []byte {
	a := &asm{labels: map[string]int{}, data: map[string]string{}, target: target, opt: opt}
	c := compiler{a: a, fs: map[string]fn{}}
	for _, f := range fs {
		if _, ok := c.fs[f.name]; ok {
			die("duplicate function %s", f.name)
		}
		c.fs[f.name] = f
	}
	c.function(fn{name: "@main", body: ss}, true)
	for _, f := range fs {
		if !opt || live[f.name] {
			c.function(f, false)
		}
	}
	a.rt() // deterministic string order by assigned label
	for i := 1; i <= a.serial; i++ {
		l := fmt.Sprintf("@%d", i)
		for s, k := range a.data {
			if k == l {
				a.label(l)
				a.b = append(a.b, []byte(s)...)
				a.emit(0)
			}
		}
	}
	if target == "linux" {
		return elf(a)
	}
	if target == "windows" {
		return pe(a)
	}
	die("unsupported target %s", target)
	return nil
}
func cli() {
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "--help" {
		fmt.Println("TRIPOL 0.6.0 — NEXUS\ntripol new NAME [--template basic|calculator|files]\ntripol learn\ntripol doctor\ntripol run [FILE.tri] [--profile nano|debug]\ntripol build [FILE.tri] [-o OUTPUT] [--target windows|linux] [--profile nano|debug]\ntripol check [FILE.tri]\ntripol bench [FILE.tri] [--runs 1..100] [--profile nano|debug]\ntripol --version")
		return
	}
	if args[0] == "--version" {
		fmt.Println("TRIPOL " + version)
		return
	}
	if args[0] == "learn" && len(args) == 1 {
		learn()
		return
	}
	if args[0] == "doctor" && len(args) == 1 {
		doctor()
		return
	}
	if args[0] == "run" || args[0] == "build" || args[0] == "check" || args[0] == "bench" {
		if len(args) == 1 || strings.HasPrefix(args[1], "-") {
			args = append([]string{args[0], "main.tri"}, args[1:]...)
		}
	}
	if len(args) < 2 {
		die("source file required; try tripol --help")
	}
	cmd, src := args[0], args[1]
	if cmd == "new" {
		template := "basic"
		if len(args) == 4 && args[2] == "--template" {
			template = args[3]
		} else if len(args) != 2 {
			die("usage: tripol new NAME [--template basic|calculator|files]")
		}
		newProjectTemplate(src, template)
		return
	}
	profile := "nano"
	target := runtime.GOOS
	out := ""
	runs := 5
	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "-o", "--output", "--target", "--profile", "--runs":
			if i+1 >= len(args) {
				die("missing option value")
			}
			v := args[i+1]
			if args[i] == "--runs" {
				var err error
				runs, err = strconv.Atoi(v)
				if cmd != "bench" || err != nil || runs < 1 || runs > 100 {
					die("--runs requires bench and a count from 1 to 100")
				}
			} else if args[i] == "--profile" {
				if v != "nano" && v != "debug" {
					die("unknown profile %s", v)
				}
				profile = v
			} else if args[i] == "--target" {
				target = v
			} else {
				out = v
			}
			i++
		default:
			die("unknown option %s", args[i])
		}
	}
	if cmd != "build" && cmd != "run" && cmd != "check" && cmd != "bench" {
		die("unknown command %s", cmd)
	}
	if (cmd == "run" || cmd == "bench") && (target != runtime.GOOS || runtime.GOARCH != "amd64") {
		die("run requires local x86-64 target")
	}
	fs, ss, inputs := loadProject(src)
	bin := compileParsed(fs, ss, target, profile)
	if cmd == "check" {
		fmt.Println("OK")
		return
	}
	ext := ""
	if target == "windows" {
		ext = ".exe"
	}
	if cmd == "run" || cmd == "bench" {
		d, e := os.MkdirTemp("", "tripol-")
		if e != nil {
			die("%v", e)
		}
		defer os.RemoveAll(d)
		out = filepath.Join(d, "program"+ext)
		write(out, bin)
		if cmd == "bench" {
			benchmark(out, runs)
		} else {
			runNative(out)
		}
		return
	}
	if out == "" {
		out = strings.TrimSuffix(src, filepath.Ext(src)) + ext
	}
	o, _ := filepath.Abs(out)
	for _, source := range inputs {
		if strings.EqualFold(source, o) {
			die("output must not overwrite source")
		}
		if si, e := os.Stat(source); e == nil {
			if oi, e := os.Stat(out); e == nil && os.SameFile(si, oi) {
				die("output aliases source")
			}
		}
	}
	write(out, bin)
	fmt.Printf("Built %s (%d bytes)\n", out, len(bin))
}
func write(p string, b []byte) {
	if e := os.WriteFile(p, b, 0755); e != nil {
		die("%v", e)
	}
	if runtime.GOOS != "windows" {
		if e := os.Chmod(p, 0755); e != nil {
			die("%v", e)
		}
	}
}
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "tripol:", r)
			os.Exit(1)
		}
	}()
	cli()
}
