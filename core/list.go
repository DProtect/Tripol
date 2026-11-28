package core

import "strings"

const callbackListPrefix = "list["

func typedListElement(t string) string {
	if strings.HasPrefix(t, callbackListPrefix) && strings.HasSuffix(t, "]") {
		return t[len(callbackListPrefix) : len(t)-1]
	}
	return ""
}
func callbackListElement(t string) string {
	if elem := typedListElement(t); strings.HasPrefix(elem, callbackPrefix) {
		return elem
	}
	return ""
}
func listElement(t string) string {
	if element := typedListElement(t); element != "" {
		return element
	}
	switch t {
	case "list":
		return "int"
	case "string_list":
		return "string"
	case "bool_list":
		return "bool"
	case "float_list":
		return "float"
	}
	return ""
}
func listType(t string) string {
	switch t {
	case "int":
		return "list"
	case "string", "bool", "float":
		return t + "_list"
	case "", "void":
		return ""
	}
	if strings.HasPrefix(t, "@") || strings.HasPrefix(t, "array:") {
		return ""
	}
	return "list[" + t + "]"
}
func isList(t string) bool { return listElement(t) != "" }

func (c *compiler) listCreate(e *expr) string {
	a := c.a
	typ := listElement(e.s)
	if typ == "" || !c.validType(typ) {
		c.err("unknown list element type " + typ)
	}
	if len(e.args) > 1 {
		c.err("list expects an optional capacity or an array of elements")
	}
	if len(e.args) == 1 && e.args[0].op == "array" {
		items := e.args[0].args
		if len(items) > 400 {
			c.err("list initializer exceeds 400 elements")
		}
		a.imm(int64(len(items)))
		a.call("list_new")
		for i, item := range items {
			a.emit(0x50)
			t := c.expression(item)
			if i == 0 && e.s == "list" {
				typ = t
				if listType(typ) == "" {
					c.err("list requires matching concrete element types")
				}
			}
			c.require(t, typ)
			a.emit(0x48, 0x89, 0xc1, 0x58, 0x50)
			a.call("list_add")
			a.emit(0x58)
		}
		return listType(typ)
	}
	if len(e.args) == 1 {
		c.require(c.expression(e.args[0]), "int")
	} else {
		a.imm(4)
	}
	a.call("list_new")
	return listType(typ)
}
func (c *compiler) listMethod(e *expr) (string, bool) {
	name, method, ok := strings.Cut(e.s, ".")
	if !ok {
		return "", false
	}
	// These methods expose only values. Strings are immutable and function
	// references point to code, so reading them cannot leak a scratch allocation.
	v, ok := c.lookupBinding(name, true)
	if !ok || !isList(v.typ) {
		return "", false
	}
	if len(c.scratchScopes) > 0 && c.outerScratchBinding(name) &&
		!(method == "size" || method == "capacity" || method == "get" && nonAggregateList(v.typ)) {
		c.err("scratch cannot access outer aggregate " + name + " for this list operation")
	}
	c.a.load(v.off)
	return c.emitListMethod(v.typ, method, e.args, v.mut), true
}

func nonAggregateList(t string) bool {
	return typedListElement(t) == "" || callbackListElement(t) != ""
}

func (c *compiler) emitListMethod(typ, method string, args []*expr, mutable bool) string {
	elem := listElement(typ)
	a := c.a
	count, known := map[string]int{"add": 1, "remove_at": 1, "get": 1, "set": 2, "size": 0, "capacity": 0, "clear": 0, "contains": 1, "remove": 1, "sort": 0}[method]
	if !known {
		c.err("unknown list method " + method)
	}
	if len(args) != count {
		c.err("wrong list argument count")
	}
	if (method == "add" || method == "remove_at" || method == "set" || method == "clear" || method == "remove" || method == "sort") && !mutable {
		c.err("cannot modify a constant or computed list; bind it to a mutable variable first")
	}
	if (method == "contains" || method == "remove" || method == "sort") && (elem == "float" || !nonAggregateList(typ)) {
		c.err("this method requires comparable list elements")
	}
	if method == "sort" {
		if strings.HasPrefix(elem, callbackPrefix) {
			c.err("cannot sort callbacks")
		}
		c.sortList(elem)
		return "int"
	}
	if count == 0 {
		if method == "capacity" {
			a.emit(0x48, 0x8b, 0x00)
			return "int"
		}
		a.emit(0x48, 0x8b, 0x40, 0x08)
		if method == "clear" {
			a.emit(0x48, 0xc7, 0x00, 0, 0, 0, 0)
			a.imm(0)
		} else {
			a.emit(0x48, 0x8b, 0x00)
		}
		return "int"
	}
	a.emit(0x50)
	argType := "int"
	if method == "add" || method == "contains" || method == "remove" {
		argType = elem
	}
	c.require(c.expression(args[0]), argType)
	if count == 2 {
		a.emit(0x50)
		c.require(c.expression(args[1]), elem)
		// RAX=value, RCX=index, RDX=stable handle; reload backing after argument calls.
		a.emit(0x59, 0x5a, 0x48, 0x8b, 0x52, 0x08, 0x48, 0x3b, 0x0a)
		a.j(0x83, "bounds")
		a.emit(0x48, 0x89, 0x44, 0xca, 0x08)
		return elem
	}
	a.emit(0x48, 0x89, 0xc1, 0x58)
	if method == "contains" || method == "remove" {
		c.findList(elem, method == "remove")
		return "bool"
	}
	switch method {
	case "add":
		a.call("list_add")
	case "remove_at":
		a.call("list_remove")
	case "get":
		a.emit(0x48, 0x8b, 0x50, 0x08, 0x48, 0x3b, 0x0a)
		a.j(0x83, "bounds")
		a.emit(0x48, 0x8b, 0x44, 0xca, 0x08)
	}
	if method == "get" || method == "remove_at" {
		return elem
	}
	return "int"
}
