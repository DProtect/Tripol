package core

import (
	"fmt"
	"maps"
	"strings"
)

// The shared frontend instantiates monomorphic signatures before validation,
// reachability analysis and either backend. It never adds a runtime value.
type simpleBinding struct {
	typ       string
	parameter int // one-based, or zero for a local
	literal   *expr
}

type simpleScope struct {
	function *fn
	locals   map[string]simpleBinding
	captures map[string]simpleBinding
	file     string
	line     int
}

type simpleResolver struct {
	functions       map[string]*fn
	overloads       map[string][]*fn
	overloadMembers map[string]bool
	captures        map[string]map[string]simpleBinding
	final, changed  bool
}

func resolveSimpleProgram(fs []fn, main []stmt) {
	if len(fs) == 0 {
		return
	}
	r := simpleResolver{
		functions: make(map[string]*fn, len(fs)),
		overloads: make(map[string][]*fn),
	}
	for i := range fs {
		f := &fs[i]
		if previous := r.functions[f.name]; previous != nil {
			if !previous.record && !f.record {
				if len(r.overloads[f.name]) == 0 {
					r.overloads[f.name] = []*fn{previous}
				}
				r.overloads[f.name] = append(r.overloads[f.name], f)
			}
		} else {
			r.functions[f.name] = f
		}
		if f.record {
			continue
		}
		for j, inferred := range f.implicit {
			if inferred {
				f.args[j].typ = ""
			}
		}
		if f.inferRet {
			f.ret = ""
			implicitTailReturn(f.body)
		}
		if f.anonymous && r.captures == nil {
			r.captures = make(map[string]map[string]simpleBinding)
		}
	}
	for key, group := range r.overloads {
		if r.overloadMembers == nil {
			r.overloadMembers = make(map[string]bool)
		}
		_, member, method := splitMethodName(key)
		offset := 0
		if method {
			offset = 1
			r.overloadMembers[member] = true
		}
		delete(r.functions, key)
		signatures := map[string]bool{}
		for i, f := range group {
			types := make([]string, 0, len(f.args)-offset)
			for _, arg := range f.args[offset:] {
				if arg.typ == "" {
					r.failure(simpleScope{function: f}, nil, "Overloads need explicit parameter types.")
				}
				types = append(types, arg.typ)
			}
			signature := strings.Join(types, ",")
			if signatures[signature] {
				label := key
				if receiver, member, method := splitMethodName(key); method {
					label = receiver + "." + member
				}
				r.failure(simpleScope{function: f}, nil, "Duplicate overload for "+label+"("+signature+").")
			}
			signatures[signature] = true
			f.name = fmt.Sprintf("%s$%d", key, i+1)
			r.functions[f.name] = f
		}
	}
	standardBodies := make([]bool, len(fs))
	for i := range fs {
		f := &fs[i]
		standardBodies[i] = !strings.HasPrefix(f.source, "tripol.") || f.inferRet || f.anonymous || r.needsResolution(f.body)
	}
	// Each update fixes an unknown parameter or result. The bound also covers
	// workspaces with a long chain of forward references.
	limit := 1
	for i := range fs {
		f := &fs[i]
		if f.inferRet {
			limit++
		}
		for _, p := range f.args {
			if p.typ == "" {
				limit++
			}
		}
	}
	for pass := 0; (limit > 1 || r.captures != nil) && pass < limit; pass++ {
		r.changed = false
		r.walkProgram(fs, main, standardBodies)
		if !r.changed {
			break
		}
	}
	for i := range fs {
		f := &fs[i]
		if f.record {
			continue
		}
		for _, arg := range f.args {
			if arg.typ == "" {
				r.failure(simpleScope{function: f}, nil, fmt.Sprintf("Cannot infer the type of parameter %q in %s. Add a type annotation or call it with a typed argument.", arg.name, readableFunction(f)))
			}
		}
		if f.ret == "" {
			f.ret = "int"
		}
	}
	r.final = true
	r.walkProgram(fs, main, standardBodies)
}

func (r *simpleResolver) needsResolution(body []stmt) bool {
	for _, s := range body {
		if r.contextualExpression(s.e) || r.contextualExpression(s.bound) ||
			r.needsResolution(s.yes) || r.needsResolution(s.no) || r.needsResolution(s.final) {
			return true
		}
	}
	return false
}

func (r *simpleResolver) contextualExpression(e *expr) bool {
	if e == nil {
		return false
	}
	switch e.op {
	case "var", "function-ref":
		if f := r.functions[e.s]; f != nil && !f.record && (e.op == "var" || f.anonymous || f.inferRet) {
			return true
		}
	case "field":
		if e.args[0].op == "var" && r.functions[e.args[0].s+"."+e.s] != nil {
			return true
		}
	case "call":
		if len(r.overloads[e.s]) > 0 {
			return true
		}
		if _, member, ok := strings.Cut(e.s, "."); ok && r.overloadMembers[member] {
			return true
		}
		if f := r.functions[e.s]; f != nil && (f.inferRet || hasInferredParameter(*f)) {
			return true
		}
	case "method":
		if r.overloadMembers[e.s] {
			return true
		}
	}
	for _, arg := range e.args {
		if r.contextualExpression(arg) {
			return true
		}
	}
	return false
}

func hasInferredParameter(f fn) bool {
	for _, inferred := range f.implicit {
		if inferred {
			return true
		}
	}
	return false
}

func implicitTailReturn(body []stmt) {
	if len(body) == 0 || returns(body) {
		return
	}
	last := &body[len(body)-1]
	if last.op == "expr" {
		last.op = "return"
	}
	if last.op == "if" && len(last.no) > 0 {
		implicitTailReturn(last.yes)
		implicitTailReturn(last.no)
	}
}

func readableFunction(f *fn) string {
	if f.anonymous {
		return "anonymous callback"
	}
	if receiver, member, ok := splitMethodName(f.name); ok {
		return receiver + "." + strings.SplitN(member, "$", 2)[0]
	}
	return f.name
}

func (r *simpleResolver) failure(scope simpleScope, e *expr, message string) {
	line, col, source := scope.line, 0, scope.file
	if scope.function != nil {
		if line == 0 {
			line, col = scope.function.line, scope.function.col
		}
		source = scope.function.source
	}
	if e != nil && e.line > 0 {
		line, col = e.line, e.col
	}
	d := diagnose("compiler", message, source, line, col)
	d.Code = "TP1010"
	d.Hint = "Add a type annotation or use a single type for each parameter."
	if strings.Contains(message, "overload") {
		d.Code = "TP1011"
		d.Hint = "Compare the argument types with the available method signatures."
	} else if strings.Contains(message, "capture") {
		d.Code = "TP1012"
		d.Hint = "Capture only a const scalar literal or pass the value as an argument."
	} else if strings.HasPrefix(message, "Callback ") {
		d.Code = "TP1013"
		d.Hint = "The callback must have the number of parameters in the expected function signature."
	}
	panic(d)
}

func (r *simpleResolver) walkProgram(fs []fn, main []stmt, bodies []bool) {
	r.block(main, simpleScope{locals: map[string]simpleBinding{}})
	for i := range fs {
		f := &fs[i]
		if f.record {
			continue
		}
		if !bodies[i] {
			continue
		}
		scope := simpleScope{function: f, locals: make(map[string]simpleBinding, len(f.args)), captures: r.captures[f.name]}
		for j, arg := range f.args {
			scope.locals[arg.name] = simpleBinding{typ: arg.typ, parameter: j + 1}
		}
		r.block(f.body, scope)
	}
}

func (r *simpleResolver) child(body []stmt, scope simpleScope) {
	if len(body) == 0 {
		return
	}
	scope.locals = maps.Clone(scope.locals)
	r.block(body, scope)
}

func (r *simpleResolver) block(body []stmt, scope simpleScope) {
	for _, s := range body {
		scope.line = s.line
		if s.source != "" {
			scope.file = s.source
		}
		want := ""
		switch s.op {
		case "return":
			if scope.function != nil {
				want = scope.function.ret
			}
		case "decl", "assign":
			want = s.typ
			if want == "" && s.op == "assign" {
				if _, local := scope.locals[s.name]; !local && r.final {
					if _, captured := scope.captures[s.name]; captured {
						r.captureError(scope, s.e, s.name)
					}
				}
				want = scope.locals[s.name].typ
			}
		case "if", "while", "assert":
			want = "bool"
		case "for", "throw":
			want = "int"
		}
		typ := r.expression(s.e, want, scope)
		r.expression(s.bound, "", scope)
		switch s.op {
		case "return":
			if f := scope.function; f != nil && f.inferRet && typ != "" {
				r.setType(&f.ret, typ, scope, s.e, "return type of "+readableFunction(f))
			}
		case "decl", "assign":
			if current, ok := scope.locals[s.name]; s.op == "assign" && ok && current.parameter != 0 {
				r.parameter(scope.function, current.parameter-1, typ, scope, s.e)
			} else {
				if s.typ != "" {
					typ = s.typ
				}
				var literal *expr
				if s.op == "decl" && !s.mut && s.e != nil {
					switch s.e.op {
					case "int", "float", "bool", "string":
						literal = s.e
					}
				}
				scope.locals[s.name] = simpleBinding{typ: typ, literal: literal}
			}
		}
		if s.op == "for" {
			loop := scope
			loop.locals = maps.Clone(scope.locals)
			loop.locals[s.name] = simpleBinding{typ: "int"}
			r.block(s.yes, loop)
		} else {
			r.child(s.yes, scope)
		}
		if s.op == "try" {
			caught := scope
			caught.locals = maps.Clone(scope.locals)
			caught.locals[s.name] = simpleBinding{typ: "int"}
			r.block(s.no, caught)
		} else {
			r.child(s.no, scope)
		}
		r.child(s.final, scope)
	}
}

func (r *simpleResolver) setType(dst *string, typ string, scope simpleScope, e *expr, what string) {
	if typ == "" || strings.HasPrefix(typ, "array:") {
		return
	}
	if *dst == "" {
		*dst = typ
		r.changed = true
	} else if *dst != typ {
		r.failure(scope, e, fmt.Sprintf("Conflicting types for %s: %s and %s. Use one type or write separate functions.", what, *dst, typ))
	}
}

func (r *simpleResolver) parameter(f *fn, i int, typ string, scope simpleScope, e *expr) {
	if i < len(f.implicit) && f.implicit[i] {
		r.setType(&f.args[i].typ, typ, scope, e, fmt.Sprintf("parameter %q of %s", f.args[i].name, readableFunction(f)))
	}
}

func (r *simpleResolver) captureError(scope simpleScope, e *expr, name string) {
	r.failure(scope, e, fmt.Sprintf("Anonymous callback cannot capture %q. Use a const scalar literal or pass the value through the callback context.", name))
}

func (r *simpleResolver) binding(e *expr, name, want string, scope simpleScope) (string, bool) {
	if b, ok := scope.locals[name]; ok {
		if b.parameter > 0 {
			i := b.parameter - 1
			r.parameter(scope.function, i, want, scope, e)
			return scope.function.args[i].typ, true
		}
		return b.typ, true
	}
	if b, ok := scope.captures[name]; ok {
		if b.literal == nil {
			r.captureError(scope, e, name)
		}
		if r.final && e.op == "var" {
			*e = expr{op: b.literal.op, s: b.literal.s, n: b.literal.n, line: e.line, col: e.col}
		}
		return b.typ, true
	}
	return "", false
}

func (r *simpleResolver) reference(e *expr, f *fn, want string, scope simpleScope) string {
	if f.anonymous && !r.final {
		captures := maps.Clone(scope.locals)
		for name, b := range scope.captures {
			if _, found := captures[name]; !found {
				captures[name] = b
			}
		}
		if previous := r.captures[f.name]; !maps.Equal(previous, captures) {
			r.captures[f.name] = captures
			r.changed = true
		}
	}
	if params, result, ok := callbackParts(want); ok {
		if len(params) != len(f.args) {
			r.failure(scope, e, fmt.Sprintf("Callback %s needs %d parameters here; its declaration has %d.", readableFunction(f), len(params), len(f.args)))
		}
		for i, typ := range params {
			r.parameter(f, i, typ, scope, e)
		}
		if f.inferRet {
			r.setType(&f.ret, result, scope, e, "return type of "+readableFunction(f))
		}
	}
	if r.final {
		e.op, e.s, e.args = "function-ref", f.name, nil
	}
	if f.ret == "" {
		return ""
	}
	for _, arg := range f.args {
		if arg.typ == "" {
			return ""
		}
	}
	return functionCallbackType(*f)
}

func (r *simpleResolver) expression(e *expr, want string, scope simpleScope) string {
	if e == nil {
		return ""
	}
	switch e.op {
	case "int", "string", "bool", "float":
		return e.op
	case "var":
		if typ, ok := r.binding(e, e.s, want, scope); ok {
			return typ
		}
		if f := r.functions[e.s]; f != nil && !f.record {
			return r.reference(e, f, want, scope)
		}
	case "function-ref":
		if f := r.functions[e.s]; f != nil && !f.record {
			return r.reference(e, f, want, scope)
		}
	case "field":
		root := e.args[0]
		if root.op == "var" {
			if _, local := r.binding(root, root.s, "", scope); !local {
				if f := r.functions[root.s+"."+e.s]; f != nil && !f.record {
					return r.reference(e, f, want, scope)
				}
			}
		}
		if record := r.functions[r.expression(root, "", scope)]; record != nil && record.record {
			for _, field := range record.args {
				if field.name == e.s {
					return field.typ
				}
			}
		}
	case "index":
		r.expression(e.args[0], "int", scope)
		typ, _ := r.binding(e, e.s, "", scope)
		if isList(typ) {
			return listElement(typ)
		}
		if typ == "ints" || strings.HasPrefix(typ, "array:") {
			return "int"
		}
	case "un":
		if e.s == "!" {
			r.expression(e.args[0], "bool", scope)
			return "bool"
		}
		return r.expression(e.args[0], want, scope)
	case "bin":
		switch e.s {
		case "&&", "||":
			for _, arg := range e.args {
				r.expression(arg, "bool", scope)
			}
			return "bool"
		}
		left, right := r.expression(e.args[0], "", scope), r.expression(e.args[1], "", scope)
		if left == "" {
			left = r.expression(e.args[0], right, scope)
		}
		if right == "" {
			right = r.expression(e.args[1], left, scope)
		}
		switch e.s {
		case "==", "!=", "<", "<=", ">", ">=":
			return "bool"
		}
		if left != "" {
			return left
		}
		if right != "" {
			return right
		}
		if want == "int" || want == "float" || (want == "string" && e.s == "+") {
			r.expression(e.args[0], want, scope)
			r.expression(e.args[1], want, scope)
			return want
		}
	case "call":
		return r.call(e, want, scope)
	case "method":
		return r.method(e, r.expression(e.args[0], "", scope), e.s, 1, scope)
	case "format":
		for _, arg := range e.args {
			r.expression(arg, "", scope)
		}
		return "string"
	case "dictionary":
		if want == "json.Object" {
			return r.jsonObject(e, scope)
		}
		for _, arg := range e.args {
			r.expression(arg, "", scope)
		}
		return ""
	case "array", "list":
		for _, arg := range e.args {
			r.expression(arg, "", scope)
		}
		if e.op == "array" {
			return "array:literal"
		}
	}
	return ""
}

func (r *simpleResolver) call(e *expr, want string, scope simpleScope) string {
	if f := r.functions[e.s]; f != nil {
		return r.arguments(f, e.args, 0, scope)
	}
	if group := r.overloads[e.s]; len(group) != 0 {
		selected := r.selectOverload(e, e.s, e.args, group, 0, scope)
		if selected == nil {
			return ""
		}
		if r.final {
			e.s = selected.name
		}
		return r.arguments(selected, e.args, 0, scope)
	}
	name := strings.TrimPrefix(e.s, "@builtin:")
	if receiver, member, ok := strings.Cut(name, "."); ok {
		if typ, found := r.binding(e, receiver, "", scope); found {
			if isList(typ) {
				for _, arg := range e.args {
					r.expression(arg, "", scope)
				}
				switch member {
				case "get", "set", "remove_at":
					return listElement(typ)
				case "contains", "remove":
					return "bool"
				default:
					return "int"
				}
			}
			return r.method(e, typ, member, 0, scope)
		}
	}
	if typ, found := r.binding(e, name, "", scope); found {
		if params, result, ok := callbackParts(typ); ok {
			for i, arg := range e.args {
				if i < len(params) {
					r.expression(arg, params[i], scope)
				}
			}
			return result
		}
	}
	if name == "list" && len(e.args) == 1 && e.args[0].op == "array" && len(e.args[0].args) > 0 {
		element := ""
		if isList(want) {
			element = listElement(want)
		}
		result := ""
		for _, arg := range e.args[0].args {
			typ := r.expression(arg, element, scope)
			if result == "" {
				result = listType(typ)
			}
		}
		return result
	}
	for _, arg := range e.args {
		r.expression(arg, "", scope)
	}
	if isList(name) {
		return name
	}
	switch name {
	case "input", "text", "slice", "utf8_char", "json_quote", "platform_name", "compiler_version", "env_get":
		return "string"
	case "len", "abs", "min", "max", "input_int", "parse_int", "sum", "@array_len", "byte_at", "monotonic_ms", "monotonic_ns":
		return "int"
	case "has_prefix", "has_suffix":
		return "bool"
	case "float", "sqrt":
		return "float"
	case "ints":
		return "ints"
	}
	return ""
}

func (r *simpleResolver) arguments(f *fn, args []*expr, offset int, scope simpleScope) string {
	for i, arg := range args {
		j := i + offset
		if j >= len(f.args) {
			r.expression(arg, "", scope)
			continue
		}
		typ := r.expression(arg, f.args[j].typ, scope)
		r.parameter(f, j, typ, scope, arg)
	}
	return f.ret
}

func (r *simpleResolver) method(e *expr, receiver, member string, offset int, scope simpleScope) string {
	key := methodName(receiver, member)
	if f := r.functions[key]; f != nil {
		return r.arguments(f, e.args[offset:], 1, scope)
	}
	group := r.overloads[key]
	if len(group) == 0 {
		for _, arg := range e.args[offset:] {
			r.expression(arg, "", scope)
		}
		return ""
	}
	args := e.args[offset:]
	selected := r.selectOverload(e, receiver+"."+member, args, group, 1, scope)
	if selected == nil {
		return ""
	}
	if r.final {
		_, resolved, _ := splitMethodName(selected.name)
		if e.op == "method" {
			e.s = resolved
		} else {
			root, _, _ := strings.Cut(e.s, ".")
			e.s = root + "." + resolved
		}
	}
	return r.arguments(selected, args, 1, scope)
}

func (r *simpleResolver) selectOverload(e *expr, label string, args []*expr, group []*fn, paramOffset int, scope simpleScope) *fn {
	types := make([]string, len(args))
	for i, arg := range args {
		types[i] = r.expression(arg, "", scope)
	}
	var selected *fn
	count := 0
	var signatures []string
	for _, f := range group {
		params := f.args[paramOffset:]
		matches := len(args) <= len(params)
		for i, arg := range params {
			if i < len(args) {
				if args[i].op == "dictionary" && arg.typ != "json.Object" {
					matches = false
				} else if types[i] != "" && types[i] != arg.typ {
					matches = false
				}
			} else if i+paramOffset >= len(f.defaults) || f.defaults[i+paramOffset] == nil {
				matches = false
			}
		}
		if matches {
			selected = f
			count++
		}
		if r.final {
			items := make([]string, len(params))
			for i, p := range params {
				items[i] = p.typ
			}
			signatures = append(signatures, "("+strings.Join(items, ", ")+")")
		}
	}
	if count != 1 {
		if r.final {
			reason := "No overload matches"
			if count > 1 {
				reason = "Ambiguous overload for"
			}
			r.failure(scope, e, fmt.Sprintf("%s %s(%s). Available signatures: %s.", reason, label, strings.Join(types, ", "), strings.Join(signatures, " or ")))
		}
		return nil
	}
	return selected
}

func (r *simpleResolver) jsonObject(e *expr, scope simpleScope) string {
	if r.functions["json.object"] == nil {
		if r.final {
			r.failure(scope, e, "Import tripol.json to use this JSON object.")
		}
		return ""
	}
	value := &expr{op: "call", s: "json.object", line: e.line, col: e.col}
	for i := 0; i < len(e.args); i += 2 {
		item := e.args[i+1]
		value = &expr{op: "method", s: "with", args: []*expr{value, e.args[i], item}, line: e.line, col: e.col}
	}
	r.expression(value, "json.Object", scope)
	if r.final {
		*e = *value
	}
	return "json.Object"
}
