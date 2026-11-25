package core

import (
	"encoding/json"
	"math"
	"regexp"
	"strings"
)

// Default arguments are scalar constants. No implicit scope or side effects.
func (p *parser) defaultValue() *expr {
	e := p.expression(1)
	if e.op == "string" || e.op == "float" {
		return e
	}
	if e.op == "un" && e.s == "-" && e.args[0].op == "float" {
		return &expr{op: "float", n: int64(math.Float64bits(-math.Float64frombits(uint64(e.args[0].n))))}
	}
	if t, n, ok := fold(e); ok {
		return &expr{op: t, n: n}
	}
	p.err("default parameter requires a scalar constant")
	return nil
}
func (c *compiler) callArguments(f fn, e *expr) []*expr {
	if len(e.args) > len(f.args) {
		c.err("wrong argument count: " + e.s)
	}
	if len(e.args) == len(f.args) {
		return e.args
	}
	args := append([]*expr(nil), e.args...)
	for i := len(args); i < len(f.args); i++ {
		if i >= len(f.defaults) || f.defaults[i] == nil {
			c.err("wrong argument count: " + e.s)
		}
		args = append(args, f.defaults[i])
	}
	return args
}

var interpolationName = regexp.MustCompile(`^[A-Za-z_][A-Za-z_0-9]*$`)

// f"...{name}..." is opt-in, so existing JSON and literal braces remain intact.
func (p *parser) interpolation() *expr {
	t := p.take()
	var s string
	if json.Unmarshal([]byte(t.s), &s) != nil || strings.ContainsRune(s, 0) {
		p.err("invalid interpolated string")
	}
	e := &expr{op: "format"}
	var buf strings.Builder
	flush := func() {
		if buf.Len() > 0 {
			e.args = append(e.args, &expr{op: "string", s: buf.String()})
			buf.Reset()
		}
	}
	for i := 0; i < len(s); {
		if i+1 < len(s) && (s[i:i+2] == "{{" || s[i:i+2] == "}}") {
			buf.WriteByte(s[i])
			i += 2
			continue
		}
		if s[i] == '{' {
			flush()
			j := strings.IndexByte(s[i+1:], '}')
			if j < 0 {
				p.err("unclosed interpolation")
			}
			j += i + 1
			name := s[i+1 : j]
			if !interpolationName.MatchString(name) {
				p.err("interpolation requires a variable name")
			}
			e.args = append(e.args, &expr{op: "var", s: name})
			i = j + 1
		} else if s[i] == '}' {
			p.err("unmatched interpolation brace")
		} else {
			buf.WriteByte(s[i])
			i++
		}
	}
	flush()
	if len(e.args) == 0 {
		return &expr{op: "string", s: ""}
	}
	return e
}
func (c *compiler) formatExpression(e *expr) string {
	a := c.a
	for i, x := range e.args {
		if i > 0 {
			a.emit(0x50)
		}
		switch t := c.expression(x); t {
		case "string":
		case "int":
			a.call("int_text")
		case "bool":
			no, done := a.fresh(), a.fresh()
			a.emit(0x48, 0x85, 0xc0)
			a.j(0x84, no)
			a.str("true")
			a.jmp(done)
			a.label(no)
			a.str("false")
			a.label(done)
		default:
			c.err("interpolation requires string, int or bool, got " + t)
		}
		if i > 0 {
			a.emit(0x48, 0x89, 0xc1, 0x58)
			a.call("concat")
		}
	}
	return "string"
}
