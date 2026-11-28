package core

import "fmt"

type jsonLowering struct {
	records   map[string]fn
	helpers   map[string]bool
	generated []fn
	hasJSON   bool
}

func lowerJSONDecoders(functions []fn, main []stmt) []fn {
	j := jsonLowering{records: make(map[string]fn), helpers: make(map[string]bool)}
	for _, f := range functions {
		if f.record {
			j.records[f.name] = f
		}
		j.hasJSON = j.hasJSON || f.name == "json.__decode_object"
	}

	var walk func([]stmt, string)
	var expression func(*expr, string)
	expression = func(e *expr, file string) {
		if e == nil {
			return
		}
		for _, arg := range e.args {
			expression(arg, file)
		}
		if e.op != "call" || e.s != "json.decode" && e.s != "json.encode" {
			return
		}
		fail := func(message string) {
			panic(diagnose("compiler", message, file, e.line, e.col))
		}
		if !j.hasJSON {
			fail("Import tripol.json before calling " + e.s + ".")
		}
		decode := e.s == "json.decode"
		if len(e.args) != 2 || e.args[0].op != "var" {
			fail("Use " + e.s + "(YourStruct, value). The first argument names a user-defined struct.")
		}
		typ := e.args[0].s
		if _, known := j.records[typ]; !known || isStandardFunction(typ) {
			fail(e.s + " requires a user-defined struct, got " + typ + ".")
		}
		ast := jsonAST{file, e.line}
		var callee, unsupported string
		if decode {
			callee, unsupported = j.decode(typ, ast)
		} else {
			callee, unsupported = j.encode(typ, ast)
		}
		if unsupported != "" {
			fail("Cannot use " + e.s + " with " + typ + ". Unsupported field: " + unsupported + ".")
		}
		e.s, e.args = callee, e.args[1:]
		if !decode {
			e.args = append(e.args, &expr{op: "int", n: 32})
		}
	}
	walk = func(body []stmt, file string) {
		for _, s := range body {
			if s.source != "" {
				file = s.source
			}
			expression(s.e, file)
			expression(s.bound, file)
			walk(s.yes, file)
			walk(s.no, file)
			walk(s.final, file)
		}
	}
	for _, f := range functions {
		walk(f.body, f.source)
	}
	walk(main, "")
	return append(functions, j.generated...)
}

type jsonAST struct {
	file string
	line int
}

func (a jsonAST) call(name string, args ...*expr) *expr {
	return &expr{op: "call", s: name, args: args, line: a.line}
}

func (a jsonAST) variable(name string) *expr { return &expr{op: "var", s: name, line: a.line} }

func (a jsonAST) text(value string) *expr { return &expr{op: "string", s: value, line: a.line} }

func (a jsonAST) bin(op string, lhs, rhs *expr) *expr {
	return &expr{op: "bin", s: op, args: []*expr{lhs, rhs}, line: a.line}
}

func (a jsonAST) declare(name string, value *expr) stmt {
	return stmt{op: "decl", name: name, e: value, mut: true, line: a.line}
}

func (a jsonAST) set(name string, value *expr) stmt {
	return stmt{op: "assign", name: name, e: value, line: a.line}
}

func (a jsonAST) ret(value *expr) stmt { return stmt{op: "return", e: value, line: a.line} }

func (a jsonAST) loop(count *expr, body ...stmt) stmt {
	return stmt{op: "for", name: "i", e: &expr{op: "int"}, bound: count, yes: body, line: a.line}
}

func (a jsonAST) check(condition *expr, message string) stmt {
	return stmt{op: "assert", name: message, e: condition, line: a.line}
}

func (j *jsonLowering) decode(typ string, a jsonAST) (string, string) {
	switch typ {
	case "int", "bool", "string", "float":
		return "json.__decode_" + typ, ""
	}
	name := "@json-decode:" + typ
	if j.helpers[name] {
		return name, ""
	}
	j.helpers[name] = true
	body := []stmt{}
	raw := a.variable("raw")
	if element := optionalElement(typ); element != "" {
		decoder, bad := j.decode(element, a)
		if bad != "" {
			return "", bad
		}
		empty := a.bin("||", a.bin("==", raw, a.text("")), a.bin("==", raw, a.text("null")))
		body = append(body, stmt{op: "if", e: empty, yes: []stmt{a.ret(&expr{op: "none", s: typ})}, line: a.line},
			a.ret(&expr{op: "some", s: typ, args: []*expr{a.call(decoder, raw)}}))
	} else if element := listElement(typ); element != "" {
		decoder, bad := j.decode(element, a)
		if bad != "" {
			return "", bad
		}
		size := a.call("len", a.variable("items"))
		body = append(body, a.declare("items", a.call("json.__decode_array", raw)),
			a.declare("decoded", a.call(typ, size)),
			a.loop(size, stmt{op: "expr", line: a.line,
				e: a.call("decoded.add", a.call(decoder, a.call("items.get", a.variable("i"))))}),
			a.ret(a.variable("decoded")))
	} else {
		record, known := j.records[typ]
		if !known || isStandardFunction(typ) {
			return "", typ
		}
		body = append(body, a.declare("object", a.call("json.__decode_object", raw)))
		values := make([]*expr, 0, len(record.args))
		for i, field := range record.args {
			decoder, bad := j.decode(field.typ, a)
			if bad != "" {
				return "", field.name + " (" + bad + ")"
			}
			source := fmt.Sprintf("field%d", i)
			value := fmt.Sprintf("value%d", i)
			body = append(body, a.declare(source, a.call("json.__decode_field", a.variable("object"), a.text(field.name))))
			expected := field.typ
			if elem := optionalElement(expected); elem != "" {
				expected = elem
			}
			kind := expected
			if isList(expected) {
				kind = "list"
			} else if _, known := j.records[expected]; known {
				kind = "object"
			}
			valid := a.call("json.__decode_type", a.variable(source), a.text(kind),
				&expr{op: "bool", n: 0})
			if optionalElement(field.typ) != "" {
				valid.args[2].n = 1
			}
			check := a.check(valid, "JSON field "+typ+"."+field.name+" must be "+field.typ)
			decoded := a.call(decoder, a.variable(source))
			present := a.bin("!=", a.variable(source), a.text(""))
			if i < len(record.defaults) && record.defaults[i] != nil {
				body = append(body, a.declare(value, record.defaults[i]),
					stmt{op: "if", e: present, yes: []stmt{check, a.set(value, decoded)}, line: a.line})
			} else {
				if optionalElement(field.typ) == "" {
					body = append(body, a.check(present, "Missing required JSON field "+typ+"."+field.name))
				}
				body = append(body, check, a.declare(value, decoded))
			}
			values = append(values, a.variable(value))
		}
		body = append(body, a.ret(a.call(typ, values...)))
	}
	j.generated = append(j.generated, fn{name: name, ret: typ, args: []param{{"raw", "string"}}, body: body, source: a.file, line: a.line})
	return name, ""
}

func jsonScalarEncoder(typ string) string {
	switch typ {
	case "int":
		return "json.integer"
	case "bool":
		return "json.boolean"
	case "string":
		return "json.quote"
	case "float":
		return "json.__encode_float"
	default:
		return ""
	}
}

func (j *jsonLowering) encode(typ string, a jsonAST) (string, string) {
	name := "@json-encode:" + typ
	if j.helpers[name] {
		return name, ""
	}
	j.helpers[name] = true
	next := a.bin("-", a.variable("depth"), &expr{op: "int", n: 1})
	body := []stmt{a.check(a.bin(">", a.variable("depth"), &expr{op: "int"}), "JSON serialization exceeded depth 32; check for a reference cycle")}
	item := a.variable("item")
	encodeField := func(t string, value *expr) (*expr, string) {
		if scalar := jsonScalarEncoder(t); scalar != "" {
			return a.call(scalar, value), ""
		}
		child, bad := j.encode(t, a)
		return a.call(child, value, next), bad
	}
	if element := optionalElement(typ); element != "" {
		value, bad := encodeField(element, a.call("item.unwrap"))
		if bad != "" {
			return "", bad
		}
		body = append(body, stmt{op: "if", e: a.call("item.has"), yes: []stmt{a.ret(value)}, line: a.line}, a.ret(a.text("null")))
	} else if element := listElement(typ); element != "" {
		value, bad := encodeField(element, a.call("item.get", a.variable("i")))
		if bad != "" {
			return "", bad
		}
		body = append(body, a.declare("output", a.text("[")),
			a.loop(a.call("len", item),
				a.set("output", a.call("json.__append", a.variable("output"), value))),
			a.ret(a.bin("+", a.variable("output"), a.text("]"))))
	} else {
		record, known := j.records[typ]
		if !known || isStandardFunction(typ) {
			return "", typ
		}
		body = append(body, a.declare("output", a.text("{")))
		for _, field := range record.args {
			value, bad := encodeField(field.typ, &expr{op: "field", s: field.name, args: []*expr{item}, line: a.line})
			if bad != "" {
				return "", field.name + " (" + bad + ")"
			}
			member := a.bin("+", a.bin("+", a.call("json.quote", a.text(field.name)), a.text(":")), value)
			body = append(body, a.set("output", a.call("json.__append", a.variable("output"), member)))
		}
		body = append(body, a.ret(a.bin("+", a.variable("output"), a.text("}"))))
	}
	j.generated = append(j.generated, fn{name: name, ret: "string", args: []param{{"item", typ}, {"depth", "int"}}, body: body, source: a.file, line: a.line})
	return name, ""
}
