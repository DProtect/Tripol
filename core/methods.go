package core

import "strings"

// Keep method symbols distinct from module-qualified functions.
const methodPrefix = "@method:"

func methodName(receiver, method string) string {
	return methodPrefix + receiver + "." + method
}

func splitMethodName(name string) (receiver, member string, ok bool) {
	if !strings.HasPrefix(name, methodPrefix) {
		return "", "", false
	}
	name = strings.TrimPrefix(name, methodPrefix)
	at := strings.LastIndexByte(name, '.')
	return name[:at], name[at+1:], true
}

func (a *asm) stackCall(name string, args int) {
	a.call("fn:" + name)
	a.dropStack(args)
}

func (a *asm) dropStack(args int) {
	size := args * 8
	if size == 0 {
		return
	}
	if a.opt && size <= 127 {
		a.emit(0x48, 0x83, 0xc4, byte(size))
		return
	}
	a.emit(0x48, 0x81, 0xc4)
	a.u32(size)
}

// RAX already holds the checked receiver. Use the ordinary stack ABI so
// a method and its callers agree even across different optimization passes.
func (c *compiler) emitMethodCall(receiver, member string, explicit []*expr) string {
	if isList(receiver) {
		return c.emitListMethod(receiver, member, explicit, false)
	}
	f, ok := c.fs[methodName(receiver, member)]
	if !ok || !c.isRecord(receiver) {
		c.err("unknown method " + receiver + "." + member)
	}
	if strings.HasPrefix(member, "__") && isStandardFunction(receiver) && c.source != f.source {
		c.err("private method " + receiver + "." + member)
	}
	placeholder := make([]*expr, len(explicit)+1)
	copy(placeholder[1:], explicit)
	args := c.callArguments(f, &expr{s: receiver + "." + member, args: placeholder})
	a := c.a
	a.emit(0x50)
	for i, arg := range args[1:] {
		c.require(c.expression(arg), f.args[i+1].typ)
		a.emit(0x50)
	}
	a.stackCall(f.name, len(args))
	return f.ret
}

func methodCandidates(fs []fn) map[string][]string {
	byMember := make(map[string][]string)
	for _, f := range fs {
		if _, member, ok := splitMethodName(f.name); ok {
			byMember[member] = append(byMember[member], f.name)
		}
	}
	return byMember
}
