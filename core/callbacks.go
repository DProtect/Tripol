package core

import "strings"

const callbackPrefix = "func("
const closureEnvironment = "@closure_env"

func callbackType(args []string, result string) string {
	return callbackPrefix + strings.Join(args, ",") + ")->" + result
}

func callbackParts(typ string) (args []string, result string, ok bool) {
	signature, ok := strings.CutPrefix(typ, callbackPrefix)
	if !ok {
		return nil, "", false
	}
	depth, start := 0, 0
	for i, ch := range signature {
		switch ch {
		case '(', '[':
			depth++
		case ']', ')':
			if depth == 0 && ch == ')' && strings.HasPrefix(signature[i:], ")->") {
				if i > start {
					args = append(args, signature[start:i])
				}
				return args, signature[i+3:], true
			}
			depth--
		case ',':
			if depth == 0 {
				args = append(args, signature[start:i])
				start = i + 1
			}
		}
	}
	return nil, "", false
}

func functionCallbackType(f fn) string {
	parameters := f.args
	if len(parameters) > 0 && parameters[0].name == closureEnvironment {
		parameters = parameters[1:]
	}
	args := make([]string, len(parameters))
	for i, param := range parameters {
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
		if e.op == "function-ref" || e.op == "closure" {
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
	plain := a.fresh()
	// Native code addresses are even. An odd reference points at an arena
	// record containing the code pointer followed by the immutable captures.
	a.emit(0x31, 0xc9, 0xa8, 1)
	a.j(0x84, plain)
	a.emit(0x48, 0x8d, 0x48, 7, 0x48, 0x8b, 0x40, 0xff)
	a.label(plain)
	a.emit(0x50, 0x51)
	for i, arg := range args {
		c.require(c.expression(arg), params[i])
		a.emit(0x50)
	}
	// The environment is pushed before user arguments. Capture-free
	// functions use their original offsets and ignore that extra slot.
	a.emit(0xff, 0x94, 0x24)
	a.u32((len(args) + 1) * 8)
	a.dropStack(len(args) + 2)
	return result, true
}
