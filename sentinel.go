package main

import "encoding/binary"

func (c *compiler) sentinelBuiltin(e *expr) string {
	n := 2
	if e.s == "ui_form" {
		n = 3
	}
	if len(e.args) != n {
		c.err(e.s + " has wrong argument count")
	}
	a := c.a
	for _, x := range e.args {
		c.require(c.expression(x), "string")
		a.emit(0x50)
	}
	if n == 3 {
		a.emit(0x5a)
	}
	a.emit(0x59, 0x58)
	a.call(e.s)
	return "string"
}

// Standard in-memory Win32 dialog template, DWORD-aligned items. No resource compiler.
func (a *asm) dialogTemplate() {
	if !a.needed("ui_form_template") {
		return
	}
	for len(a.b)%4 != 0 {
		a.emit(0)
	}
	a.label("ui_form_template")
	word := func(n int) { a.b = binary.LittleEndian.AppendUint16(a.b, uint16(n)) }
	dword := func(n int) { a.u32(n) }
	dword(0x80c80080)
	dword(0)
	word(4)
	word(0)
	word(0)
	word(240)
	word(86)
	word(0)
	word(0)
	word(0)
	item := func(style, x, y, w, h, id, class int, title string) {
		for len(a.b)%4 != 0 {
			a.emit(0)
		}
		dword(style)
		dword(0)
		word(x)
		word(y)
		word(w)
		word(h)
		word(id)
		word(0xffff)
		word(class)
		for _, r := range title {
			word(int(r))
		}
		word(0)
		word(0)
	}
	item(0x50000000, 8, 8, 224, 16, 100, 0x82, "")
	item(0x50810080, 8, 26, 224, 16, 101, 0x81, "")
	item(0x50010001, 100, 57, 64, 18, 1, 0x80, "OK")
	item(0x50010000, 168, 57, 64, 18, 2, 0x80, "Cancel")
}
