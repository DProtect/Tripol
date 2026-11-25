package core

// Native Gateway intrinsics. The worker never touches TRIPOL's arena.
func (c *compiler) nexusGatewayBuiltin(e *expr) (string, bool) {
	var args []string
	var ret string
	switch e.s {
	case "gw_open", "gw_close":
		ret = "bool"
	case "gw_abort":
		ret = "int"
	case "gw_send":
		args = []string{"string"}
		ret = "bool"
	case "gw_receive":
		args = []string{"int"}
		ret = "string"
	default:
		return "", false
	}
	if len(e.args) != len(args) {
		c.err(e.s + " has wrong argument count")
	}
	for i, x := range e.args {
		c.require(c.expression(x), args[i])
	}
	c.a.call(e.s)
	return ret, true
}
