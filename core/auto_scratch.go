package core

import "strings"

// Conservative region proof: scalar values, immutable strings and locally
// owned integer arrays/lists. No calls to arbitrary user functions, reference
// stores into outer bindings, control transfer, IO or opaque operations.
// Conditional scopes are analyzed independently, including both branches.
func (c *compiler) localLoopRegion(ss []stmt) bool {
	if len(c.scratchScopes) > 0 || len(ss) == 0 {
		return false
	}
	local := map[string]string{}
	allocates := false
	get := func(name string) (string, bool) {
		if t, ok := local[name]; ok {
			return t, true
		}
		if v, ok := c.lookup(name); ok && (v.typ == "int" || v.typ == "bool" || v.typ == "string" || v.typ == "float") {
			return v.typ, true
		}
		return "", false
	}
	var exprType func(*expr) string
	exprType = func(e *expr) string {
		if e == nil {
			return ""
		}
		switch e.op {
		case "int", "bool", "string", "float":
			return e.op
		case "var":
			t, _ := get(e.s)
			return t
		case "index":
			if exprType(e.args[0]) == "int" {
				if local[e.s] == "ints" {
					return "int"
				}
				return listElement(local[e.s])
			}
		case "format":
			for _, arg := range e.args {
				t := exprType(arg)
				if t != "int" && t != "bool" && t != "string" {
					return ""
				}
			}
			allocates = true
			return "string"
		case "un":
			t := exprType(e.args[0])
			if (e.s == "-" && t == "int") || (e.s == "!" && t == "bool") {
				return t
			}
		case "bin":
			a, z := exprType(e.args[0]), exprType(e.args[1])
			if a == "" || a != z {
				return ""
			}
			switch e.s {
			case "+":
				if a == "string" {
					allocates = true
					return "string"
				}
				if a == "int" {
					return "int"
				}
			case "-", "*", "/", "%":
				if a == "int" {
					return "int"
				}
			case "==", "!=":
				if a == "int" || a == "bool" || a == "string" {
					return "bool"
				}
			case "<", "<=", ">", ">=":
				if a == "int" {
					return "bool"
				}
			case "&&", "||":
				if a == "bool" {
					return "bool"
				}
			}
		case "call":
			if _, exists := c.fs[e.s]; exists {
				return ""
			}
			if name, method, ok := strings.Cut(e.s, "."); ok && isList(local[name]) {
				element := listElement(local[name])
				count, known := map[string]int{"add": 1, "remove_at": 1, "get": 1, "set": 2, "size": 0, "capacity": 0, "clear": 0}[method]
				if !known || len(e.args) != count {
					return ""
				}
				for i, arg := range e.args {
					want := "int"
					if method == "add" || i == 1 {
						want = element
					}
					if exprType(arg) != want {
						return ""
					}
				}
				if method == "add" {
					allocates = true
				}
				if method == "get" || method == "set" || method == "remove_at" {
					return element
				}
				return "int"
			}
			switch e.s {
			case "ints":
				if len(e.args) == 1 && exprType(e.args[0]) == "int" {
					allocates = true
					return "ints"
				}
			case "list", "string_list", "bool_list", "float_list":
				if len(e.args) == 0 || (len(e.args) == 1 && exprType(e.args[0]) == "int") {
					allocates = true
					return e.s
				}
			case "len":
				if len(e.args) == 1 {
					t := exprType(e.args[0])
					if t == "ints" || isList(t) || t == "string" {
						return "int"
					}
				}
			case "text":
				if len(e.args) == 1 && exprType(e.args[0]) == "int" {
					allocates = true
					return "string"
				}
			case "slice":
				if len(e.args) == 3 && exprType(e.args[0]) == "string" && exprType(e.args[1]) == "int" && exprType(e.args[2]) == "int" {
					allocates = true
					return "string"
				}
			case "byte_at":
				if len(e.args) == 2 && exprType(e.args[0]) == "string" && exprType(e.args[1]) == "int" {
					return "int"
				}
			case "abs", "min", "max":
				n := 2
				if e.s == "abs" {
					n = 1
				}
				if len(e.args) != n {
					return ""
				}
				for _, arg := range e.args {
					if exprType(arg) != "int" {
						return ""
					}
				}
				return "int"
			}
		}
		return ""
	}
	var prove func([]stmt) bool
	child := func(body []stmt, counter string) bool {
		parent := local
		local = make(map[string]string, len(parent)+1)
		for name, typ := range parent {
			local[name] = typ
		}
		if counter != "" {
			local[counter] = "int"
		}
		ok := prove(body)
		local = parent
		return ok
	}
	prove = func(body []stmt) bool {
		for _, s := range body {
			switch s.op {
			case "if":
				if exprType(s.e) != "bool" {
					return false
				}
				// Existing bindings remain visible. New branch-local names do not
				// become available in the sibling branch or after the join.
				for _, branch := range [][]stmt{s.yes, s.no} {
					if !child(branch, "") {
						return false
					}
				}
			case "while", "for":
				// Outer allocations can be reclaimed around nested iteration only when
				// the nested loop itself allocates no temporaries on each backedge.
				previous := allocates
				allocates = false
				conditionOK := false
				if s.op == "for" {
					conditionOK = exprType(s.e) == "int" && exprType(s.bound) == "int"
				} else {
					conditionOK = exprType(s.e) == "bool"
				}
				if !conditionOK || !child(s.yes, s.name) || allocates {
					return false
				}
				allocates = previous
			case "decl", "assign":
				t := exprType(s.e)
				if t == "" {
					return false
				}
				if old, ok := get(s.name); ok {
					if old != t || s.op == "decl" {
						return false
					}
					// Outer strings are readable but cannot be reassigned within a region.
					if _, inside := local[s.name]; !inside && t == "string" {
						return false
					}
				} else {
					if _, exists := c.lookup(s.name); exists {
						return false
					}
					local[s.name] = t
				}
			case "indexassign":
				typ := listElement(local[s.name])
				if local[s.name] == "ints" {
					typ = "int"
				}
				if typ == "" || exprType(s.bound) != "int" || exprType(s.e) != typ {
					return false
				}
			case "assert":
				if exprType(s.e) != "bool" {
					return false
				}
			case "show", "expr":
				if exprType(s.e) == "" {
					return false
				}
			case "showmany":
				for _, e := range s.e.args {
					if exprType(e) == "" {
						return false
					}
				}
			default:
				return false
			}
		}
		return true
	}
	return prove(ss) && allocates
}
func (c *compiler) loopBody(ss []stmt) {
	if c.localLoopRegion(ss) {
		c.a.guard = true
		c.scratchBlock(stmt{op: "scratch", yes: ss})
		return
	}
	c.block(ss, true)
}
