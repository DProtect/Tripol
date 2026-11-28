package core

import (
	"fmt"
	"strings"
)

func argumentValue(e *expr) *expr {
	if e.op == "named" {
		return e.args[0]
	}
	return e
}

// bindArgumentSlots uses the declared parameter names for both overload selection
// and lowering. Its indices do not include a method's implicit receiver.
func bindArgumentSlots(f *fn, args []*expr, offset int) ([]int, string) {
	params := f.args[offset:]
	if len(args) > len(params) {
		return nil, "too many arguments"
	}
	slots := make([]int, len(args))
	used := make([]bool, len(params))
	named := false
	for i, arg := range args {
		slot := i
		if arg.op == "named" {
			named = true
			slot = -1
			for j, p := range params {
				if p.name == arg.s {
					slot = j
					break
				}
			}
			if slot == -1 {
				return nil, fmt.Sprintf("unknown named argument %q", arg.s)
			}
		} else if named {
			return nil, "place positional arguments before named arguments"
		}
		if used[slot] {
			return nil, fmt.Sprintf("duplicate argument %q", params[slot].name)
		}
		used[slot] = true
		slots[i] = slot
	}
	for i, p := range params {
		if !used[i] && (i+offset >= len(f.defaults) || f.defaults[i+offset] == nil) {
			return nil, fmt.Sprintf("missing required argument %q", p.name)
		}
	}
	return slots, ""
}

func (r *simpleResolver) arguments(e *expr, f *fn, skip, offset int, scope simpleScope) string {
	args := e.args[skip:]
	named := false
	for _, arg := range args {
		named = named || arg.op == "named"
	}
	if !named {
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

	slots, issue := bindArgumentSlots(f, args, offset)
	if issue != "" {
		if r.final {
			names := make([]string, 0, len(f.args)-offset)
			for _, p := range f.args[offset:] {
				names = append(names, p.name)
			}
			r.failure(scope, e, fmt.Sprintf("%s in %s. Parameters: %s.", issue, readableFunction(f), strings.Join(names, ", ")))
		}
		return f.ret
	}
	for i, arg := range args {
		j := slots[i] + offset
		typ := r.expression(argumentValue(arg), f.args[j].typ, scope)
		r.parameter(f, j, typ, scope, arg)
	}
	if !r.final {
		return f.ret
	}
	if f.record && strings.HasPrefix(f.source, "tripol.") && scope.file != f.source {
		for _, slot := range slots {
			if strings.HasPrefix(f.args[slot].name, "__") {
				r.failure(scope, e, "private field "+f.name+"."+f.args[slot].name)
			}
		}
	}
	lowerNamedCall(e, f, args, slots, offset)
	return f.ret
}

func lowerNamedCall(e *expr, f *fn, args []*expr, slots []int, offset int) {
	inner := *e
	inner.args = make([]*expr, len(f.args))
	bindings := make([]*expr, 0, len(args)+offset+1)
	store := func(value *expr, slot int) {
		name := fmt.Sprintf("@arg%d", slot)
		bindings = append(bindings, &expr{op: "arg", s: name, args: []*expr{value}})
		inner.args[slot] = &expr{op: "var", s: name}
	}
	if offset == 1 {
		var receiver *expr
		if e.op == "method" {
			receiver = e.args[0]
		} else {
			root, _, _ := strings.Cut(e.s, ".")
			receiver = &expr{op: "var", s: root, line: e.line, col: e.col}
		}
		inner.op = "method"
		_, inner.s, _ = splitMethodName(f.name)
		store(receiver, 0)
	}
	for i, slot := range slots {
		store(argumentValue(args[i]), slot+offset)
	}
	for i := offset; i < len(inner.args); i++ {
		if inner.args[i] == nil {
			inner.args[i] = f.defaults[i]
		}
	}
	e.op, e.s, e.args = "bound-args", "", append(bindings, &inner)
}
