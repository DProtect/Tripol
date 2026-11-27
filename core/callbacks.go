package core

import "strings"

const callbackPrefix = "func("

func callbackType(args []string, result string) string {
	return callbackPrefix + strings.Join(args, ",") + ")->" + result
}

func callbackParts(typ string) (args []string, result string, ok bool) {
	signature, ok := strings.CutPrefix(typ, callbackPrefix)
	if !ok {
		return nil, "", false
	}
	params, result, ok := strings.Cut(signature, ")->")
	if params != "" {
		args = strings.Split(params, ",")
	}
	return args, result, ok
}

func functionCallbackType(f fn) string {
	args := make([]string, len(f.args))
	for i, param := range f.args {
		args[i] = param.typ
		if strings.HasPrefix(param.typ, callbackPrefix) || callbackListElement(param.typ) != "" {
			return ""
		}
	}
	if f.record || strings.HasPrefix(f.ret, callbackPrefix) || callbackListElement(f.ret) != "" {
		return ""
	}
	return callbackType(args, f.ret)
}

func callbackTargets(fs []fn, main []stmt) map[string]bool {
	var targets map[string]bool
	var expression func(*expr)
	expression = func(e *expr) {
		if e == nil {
			return
		}
		if e.op == "function-ref" {
			if targets == nil {
				targets = make(map[string]bool)
			}
			targets[e.s] = true
		}
		for _, arg := range e.args {
			expression(arg)
		}
	}
	var block func([]stmt)
	block = func(stmts []stmt) {
		for _, s := range stmts {
			expression(s.e)
			expression(s.bound)
			block(s.yes)
			block(s.no)
			block(s.final)
		}
	}
	block(main)
	for _, f := range fs {
		block(f.body)
	}
	return targets
}

func (c *compiler) functionReference(name string) string {
	f, ok := c.fs[name]
	if !ok {
		c.err("unknown callback target " + name)
	}
	typ := functionCallbackType(f)
	if typ == "" {
		c.err("callback target must be a function with ordinary argument and result types: " + name)
	}
	c.a.rel("fn:"+name, 0x48, 0x8d, 0x05)
	return typ
}

func (c *compiler) callCallback(name string, v variable, args []*expr) (string, bool) {
	params, result, ok := callbackParts(v.typ)
	if !ok {
		return "", false
	}
	if len(args) != len(params) {
		c.err("wrong callback argument count: " + name)
	}
	a := c.a
	a.load(v.off)
	a.emit(0x50)
	for i, arg := range args {
		c.require(c.expression(arg), params[i])
		a.emit(0x50)
	}
	// Arguments are above the saved target on the ordinary TRIPOL stack.
	a.emit(0xff, 0x94, 0x24)
	a.u32(len(args) * 8)
	a.dropStack(len(args) + 1)
	return result, true
}
