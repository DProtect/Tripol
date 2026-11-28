package core

import "strings"

func startsAwaitOperand(t token) bool {
	return t.k == "id" && t.s != "and" && t.s != "or" || t.s == "("
}

// An async block constructs a startable operation on the caller thread. Only
// its .start() submits work. The backend owns scheduling and result lifetime.
func (r *simpleResolver) asyncExpression(e *expr, scope simpleScope) string {
	inner := e.args[0]
	start := r.asyncOperation(inner, scope)
	if e.op == "async" {
		if start == nil {
			start = &expr{op: "method", s: "start", args: []*expr{inner}, line: e.line, col: e.col}
		}
		typ := r.expression(start, "", scope)
		if r.final {
			wait := methodName(typ, "await")
			if r.functions[wait] == nil && len(r.overloads[wait]) == 0 {
				r.failure(scope, e, "async needs a startable operation that returns a task with await(). Try http.get(url).")
			}
			*e = *start
		}
		return typ
	}
	if start != nil {
		inner = start
	}
	await := &expr{op: "method", s: "await", args: []*expr{inner}, line: e.line, col: e.col}
	typ := r.expression(await, "", scope)
	if r.final {
		if typ == "" {
			r.failure(scope, e, "await needs a task with await() or an HTTP operation.")
		}
		*e = *await
	}
	return typ
}

func (r *simpleResolver) asyncOperation(value *expr, scope simpleScope) *expr {
	var receiver *expr
	switch value.op {
	case "call":
		root, member, qualified := strings.Cut(value.s, ".")
		if !qualified {
			return nil
		}
		if root == "http" || root == "https" {
			method := map[string]string{"get": "GET", "head": "HEAD", "delete": "DELETE", "options": "OPTIONS"}[member]
			declared := r.functions[value.s] != nil || len(r.overloads[value.s]) != 0
			if declared && method != "" && len(value.args) == 1 && r.functions[root+".request"] != nil {
				receiver = &expr{op: "call", s: root + ".request", args: []*expr{
					{op: "string", s: method}, argumentValue(value.args[0]),
				}, line: value.line, col: value.col}
			}
		} else if member == "send" && len(value.args) == 0 {
			if binding, found := scope.locals[root]; found && binding.typ == "http.Request" {
				receiver = &expr{op: "var", s: root, line: value.line, col: value.col}
			}
		}
	case "method":
		if value.s == "send" && len(value.args) == 1 && r.expression(value.args[0], "", scope) == "http.Request" {
			receiver = value.args[0]
		}
	}
	if receiver == nil {
		return nil
	}
	return &expr{op: "method", s: "start", args: []*expr{receiver}, line: value.line, col: value.col}
}
