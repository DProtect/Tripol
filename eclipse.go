package main

// Fixed private bot slots contain copies rather than references into scratch.
// Layout: token, channel, then 16 command/response pairs.
const botStateSize = 39 * 528

func (c *compiler) eclipseBuiltin(e *expr) (string, bool) {
	a := c.a
	var types []string
	var ret string
	switch e.s {
	case "bot_slot_set":
		types = []string{"int", "string"}
		ret = "bool"
	case "bot_slot_get":
		types = []string{"int"}
		ret = "string"
	case "array_copy":
		types = []string{"ints"}
		ret = "ints"
	case "array_copy_into":
		types = []string{"ints", "ints"}
		ret = "bool"
	case "array_find":
		types = []string{"ints", "int"}
		ret = "int"
	default:
		return "", false
	}
	if len(e.args) != len(types) {
		c.err(e.s + " has wrong argument count")
	}
	for i, x := range e.args {
		c.require(c.expression(x), types[i])
		a.emit(0x50)
	}
	if len(types) == 2 {
		a.emit(0x59)
	}
	a.emit(0x58)
	a.call(e.s)
	return ret, true
}
