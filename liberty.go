package main

// Liberty R1: native scalar floats and immutable, process-lifetime UTF-8 text.
// The generated program uses only OS services; no C runtime is linked.
func (c *compiler) builtin(e *expr) (string, bool) {
	a := c.a
	if t, ok := c.nexusBuiltin(e); ok {
		return t, true
	}
	switch e.s {
	case "arena_used", "arena_committed", "array_fill", "game_open", "game_running", "game_present", "game_key", "game_close":
		return c.quantumBuiltin(e), true
	case "ui_form", "https_get":
		return c.sentinelBuiltin(e), true
	case "ui_dialog":
		return c.uiBuiltin(e), true
	case "find_bytes":
		if len(e.args) != 2 {
			c.err("find_bytes expects two strings")
		}
		c.require(c.expression(e.args[0]), "string")
		a.emit(0x50)
		c.require(c.expression(e.args[1]), "string")
		a.emit(0x48, 0x89, 0xc1, 0x58)
		a.call("find_bytes")
		return "int", true
	case "tcp_exchange":
		return c.networkBuiltin(e), true
	case "platform_name", "compiler_version", "monotonic_ms", "sleep_ms", "has_prefix", "has_suffix":
		return c.ecosystemBuiltin(e), true
	case "byte_at":
		return c.byteAt(e), true
	case "ints":
		if len(e.args) != 1 {
			c.err("ints expects one length")
		}
		c.require(c.expression(e.args[0]), "int")
		a.call("ints_alloc")
		return "ints", true
	case "text":
		if len(e.args) != 1 {
			c.err("text expects one integer")
		}
		c.require(c.expression(e.args[0]), "int")
		a.call("int_text")
		return "string", true
	case "slice":
		if len(e.args) != 3 {
			c.err("slice expects text, start, end")
		}
		c.require(c.expression(e.args[0]), "string")
		a.emit(0x50)
		c.require(c.expression(e.args[1]), "int")
		a.emit(0x50)
		c.require(c.expression(e.args[2]), "int")
		a.emit(0x48, 0x89, 0xc2, 0x59, 0x58)
		a.call("text_slice")
		return "string", true
	case "read_text", "write_text", "append_text":
		return c.fileBuiltin(e), true
	case "len":
		if len(e.args) != 1 {
			c.err("len expects one argument")
		}
		if a.opt && e.args[0].op == "string" {
			a.imm(int64(len(e.args[0].s)))
			return "int", true
		}
		if e.args[0].op == "var" {
			if v, ok := c.lookup(e.args[0].s); ok && len(v.typ) >= 6 && v.typ[:6] == "array:" {
				a.imm(int64(arrayLen(v)))
				return "int", true
			}
		}
		t := c.expression(e.args[0])
		if t == "ints" {
			a.emit(0x48, 0x8b, 0x00)
			return "int", true
		}
		c.require(t, "string")
		a.call("strlen")
		return "int", true
	case "input":
		c.inputPrompt(e)
		a.call("input")
		return "string", true
	case "parse_int":
		if len(e.args) != 1 {
			c.err("parse_int expects one argument")
		}
		c.require(c.expression(e.args[0]), "string")
		a.call("parseint")
		return "int", true
	case "float":
		if len(e.args) != 1 {
			c.err("float expects one argument")
		}
		c.require(c.expression(e.args[0]), "int")
		a.emit(0xf2, 0x48, 0x0f, 0x2a, 0xc0, 0x66, 0x48, 0x0f, 0x7e, 0xc0)
		return "float", true
	case "sqrt":
		if len(e.args) != 1 {
			c.err("sqrt expects one argument")
		}
		c.require(c.expression(e.args[0]), "float")
		a.emit(0x66, 0x48, 0x0f, 0x6e, 0xc0, 0xf2, 0x0f, 0x51, 0xc0, 0x66, 0x48, 0x0f, 0x7e, 0xc0)
		return "float", true
	}
	return "", false
}
func (c *compiler) floatOp(op string) string {
	a := c.a
	a.emit(0x66, 0x48, 0x0f, 0x6e, 0xc0, 0x66, 0x48, 0x0f, 0x6e, 0xc9) // movq xmm0,rax; xmm1,rcx
	if code, ok := map[string]byte{"+": 0x58, "-": 0x5c, "*": 0x59, "/": 0x5e}[op]; ok {
		a.emit(0xf2, 0x0f, code, 0xc1, 0x66, 0x48, 0x0f, 0x7e, 0xc0)
		return "float"
	}
	if code, ok := map[string]byte{"==": 0x94, "!=": 0x95, "<": 0x92, "<=": 0x96, ">": 0x97, ">=": 0x93}[op]; ok {
		a.emit(0x66, 0x0f, 0x2e, 0xc1, 0x0f, code, 0xc0) // ucomisd / setcc
		if op == "!=" {
			a.emit(0x0f, 0x9a, 0xc2, 0x08, 0xd0)
		} else {
			a.emit(0x0f, 0x9b, 0xc2, 0x20, 0xd0)
		}
		a.emit(0x0f, 0xb6, 0xc0)
		return "bool"
	}
	c.err("unsupported float operator " + op)
	return ""
}

// Runtime sections contain hand-written x86-64, with relocations resolved by Forge.
// Ordered dependency closure keeps unused helpers out of tiny binaries.
func (a *asm) libertyRuntime() {
	order := []string{"game_open", "game_running", "game_present", "game_key", "game_close", "game_error", "arena_used", "arena_committed", "array_fill", "ui_form", "ui_form_proc", "https_get", "https_error", "tcp_exchange", "net_connect", "net_send", "net_shutdown", "net_recv", "net_close", "net_error", "find_bytes", "clock_ms", "sleep_ms", "has_prefix", "has_suffix", "raw_bool", "raw_int", "raw_float", "raw_string", "print:float", "input", "concat", "parseint", "file_read", "file_write", "file_append", "file_write_common", "file_open", "file_chunk_read", "file_chunk_write", "file_close", "file_error", "int_text", "text_slice", "ints_alloc", "strlen", "alloc", "osalloc", "oscommit", "readbyte", "text_error", "float_error"}
	for _, name := range order {
		if name == "ui_form_proc" && a.target != "windows" {
			continue
		}
		if !a.needed(name) {
			continue
		}
		key := name
		if name == "oscommit" || name == "game_open" || name == "game_running" || name == "game_present" || name == "game_key" || name == "game_close" || name == "ui_form" || name == "https_get" || name == "net_connect" || name == "net_send" || name == "net_shutdown" || name == "net_recv" || name == "net_close" || name == "clock_ms" || name == "sleep_ms" || name == "osalloc" || name == "readbyte" || name == "file_open" || name == "file_chunk_read" || name == "file_chunk_write" || name == "file_close" {
			key += ":" + a.target
		}
		t, ok := libertyCode[key]
		if !ok {
			die("missing runtime %s", key)
		}
		a.label(name)
		base := len(a.b)
		a.b = append(a.b, t.code...)
		for _, r := range t.reloc {
			a.fixes = append(a.fixes, fix{base + r.at, r.label})
		}
	}
}

type nativeTemplate struct {
	code  []byte
	reloc []fix
}
