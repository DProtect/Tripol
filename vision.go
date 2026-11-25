package main

// Native Windows message boxes. Internal ABI: RAX title, RCX message, RDX flags.
func (c *compiler) uiBuiltin(e *expr) string {
	if len(e.args) != 3 {
		c.err("ui_dialog expects title, message and flags")
	}
	a := c.a
	c.require(c.expression(e.args[0]), "string")
	a.emit(0x50)
	c.require(c.expression(e.args[1]), "string")
	a.emit(0x50)
	c.require(c.expression(e.args[2]), "int")
	a.emit(0x48, 0x89, 0xc2, 0x59, 0x58)
	a.call("ui_dialog")
	return "int"
}
func (a *asm) visionRuntime() {
	if !a.needed("ui_dialog") {
		return
	}
	a.label("ui_dialog")
	if a.target == "windows" {
		// Save frame, align stack, reserve shadow space; MessageBoxA(NULL,msg,title,flags).
		a.emit(0x55, 0x48, 0x89, 0xe5, 0x48, 0x83, 0xe4, 0xf0, 0x48, 0x83, 0xec, 0x20, 0x49, 0x89, 0xc0, 0x41, 0x89, 0xd1, 0x48, 0x89, 0xca, 0x31, 0xc9)
		a.rel("iat:19", 0xff, 0x15)
		a.emit(0x85, 0xc0)
		a.j(0x84, "ui_error")
		a.emit(0x48, 0x98, 0xc9, 0xc3)
	}
	a.label("ui_error")
	msg := "TRIPOL runtime error: native dialogs require Windows or UI operation failed\n"
	a.str(msg)
	a.emit(0x48, 0x89, 0xc6, 0xba)
	a.u32(len(msg))
	a.call("writeerr")
	a.imm(17)
	a.jmp("exit")
}
