package core

type jsonDecoder struct {
	records   map[string]fn
	helpers   map[string]bool
	generated []fn
	hasJSON   bool
}

func lowerJSONDecoders(functions []fn, main []stmt) []fn {
	d := jsonDecoder{records: make(map[string]fn), helpers: make(map[string]bool)}
	for _, f := range functions {
		if f.record {
			d.records[f.name] = f
		}
		d.hasJSON = d.hasJSON || f.name == "json.__decode_object"
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
		if e.op != "call" || e.s != "json.decode" {
			return
		}
		fail := func(message string) {
			panic(diagnose("compiler", message, file, e.line, e.col))
		}
		if !d.hasJSON {
			fail("Import tripol.json before calling json.decode.")
		}
		if len(e.args) != 2 || e.args[0].op != "var" {
			fail("Use json.decode(YourStruct, jsonText). The first argument names a user-defined struct.")
		}
		typ := e.args[0].s
		if _, known := d.records[typ]; !known || isStandardFunction(typ) {
			fail("json.decode requires a user-defined struct, got " + typ + ".")
		}
		callee, unsupported := d.declare(typ, file, e.line)
		if unsupported != "" {
			fail("Cannot decode " + typ + " from JSON. Unsupported field type: " + unsupported + ".")
		}
		e.s, e.args = callee, e.args[1:]
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
	return append(functions, d.generated...)
}

func (d *jsonDecoder) declare(typ, file string, line int) (name, unsupported string) {
	switch typ {
	case "int", "bool", "string":
		return "json.__decode_" + typ, ""
	}
	name = "@json:" + typ
	if d.helpers[typ] {
		return name, ""
	}
	d.helpers[typ] = true
	call := func(name string, args ...*expr) *expr {
		return &expr{op: "call", s: name, args: args, line: line}
	}
	variable := func(name string) *expr { return &expr{op: "var", s: name, line: line} }
	create := func(name string, value *expr) stmt {
		return stmt{op: "decl", name: name, mut: true, e: value, line: line}
	}
	generated := fn{name: name, ret: typ, args: []param{{"raw", "string"}}, source: file, line: line}
	switch {
	case isList(typ):
		decode, bad := d.declare(listElement(typ), file, line)
		if bad != "" {
			return "", bad
		}
		generated.body = []stmt{
			create("items", call("json.__decode_array", variable("raw"))),
			create("decoded", call(typ, call("len", variable("items")))),
			{op: "for", name: "i", e: &expr{op: "int"}, bound: call("len", variable("items")), line: line,
				yes: []stmt{{op: "expr", line: line, e: call("decoded.add", call(decode, call("items.get", variable("i"))))}}},
			{op: "return", e: variable("decoded"), line: line},
		}
	default:
		record, known := d.records[typ]
		if !known || isStandardFunction(typ) {
			return "", typ
		}
		generated.body = append(generated.body, create("object", call("json.__decode_object", variable("raw"))))
		args := make([]*expr, 0, len(record.args))
		for _, field := range record.args {
			decode, bad := d.declare(field.typ, file, line)
			if bad != "" {
				return "", field.name + " (" + bad + ")"
			}
			args = append(args, call(decode, call("json.__decode_field", variable("object"),
				&expr{op: "string", s: field.name, line: line})))
		}
		generated.body = append(generated.body, stmt{op: "return", e: call(typ, args...), line: line})
	}
	d.generated = append(d.generated, generated)
	return name, ""
}
