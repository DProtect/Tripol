package main

import "fmt"

// Horizon records are typed references to arena-resident sequences of 8-byte slots.
func (c *compiler) isRecord(t string) bool { f, ok := c.fs[t]; return ok && f.record }
func (c *compiler) validType(t string) bool {
	switch t {
	case "int", "float", "bool", "string", "ints":
		return true
	}
	return c.isRecord(t)
}
func (c *compiler) validateTypes(fs []fn) {
	for _, f := range fs {
		if !c.validType(f.ret) {
			die("unknown type %s in %s", f.ret, f.name)
		}
		if f.record && (len(f.args) == 0 || len(f.args) > 64) {
			die("struct %s requires 1 to 64 fields", f.name)
		}
		seen := map[string]bool{}
		for _, p := range f.args {
			if seen[p.name] {
				die("duplicate field or parameter %s in %s", p.name, f.name)
			}
			seen[p.name] = true
			if !c.validType(p.typ) {
				die("unknown type %s in %s", p.typ, f.name)
			}
			if f.record && c.isRecord(p.typ) {
				die("nested struct fields are not supported yet: %s.%s", f.name, p.name)
			}
		}
	}
}
func (c *compiler) recordField(t, n string) (param, int) {
	f, ok := c.fs[t]
	if !ok || !f.record {
		c.err("field access requires a struct, got " + t)
	}
	for i, p := range f.args {
		if p.name == n {
			return p, i * 8
		}
	}
	c.err("unknown field " + t + "." + n)
	return param{}, 0
}
func (c *compiler) recordConstruct(f fn, e *expr) string {
	if len(f.args) != len(e.args) {
		c.err("wrong argument count: " + f.name)
	}
	a := c.a
	for i, x := range e.args {
		c.require(c.expression(x), f.args[i].typ)
		a.emit(0x50)
	}
	a.imm(int64(len(f.args) * 8))
	a.call("alloc")
	for i := len(f.args) - 1; i >= 0; i-- {
		a.emit(0x59, 0x48, 0x89, 0x88)
		a.u32(i * 8)
	}
	return f.name
}
func (c *compiler) recordAssign(s stmt) {
	base := s.bound.args[0]
	if base.op != "var" {
		c.err("field assignment requires a named struct variable")
	}
	v, ok := c.lookup(base.s)
	if !ok {
		c.err("unknown variable " + base.s)
	}
	if !v.mut {
		c.err("cannot modify constant " + base.s)
	}
	f, off := c.recordField(v.typ, s.bound.s)
	c.require(c.expression(s.e), f.typ)
	c.a.loadRCX(v.off)
	c.a.emit(0x48, 0x89, 0x81)
	c.a.u32(off)
}
func learnSpanish() {
	fmt.Print(`TRIPOL 0.8 — Aprende en cinco minutos

1. Mostrar: show "Hola", 10 + 20
2. Preguntar: nombre = input("Tu nombre?")
   edad = input_int("Tu edad?")
3. Decidir: if edad >= 18 { show "Adulto" } else { show "Menor" }
4. Repetir: repeat 3 { show "Hola" }
   for i in range(5) { show i }
5. Datos: numeros = [10, 20, 30]
   show sum(numeros)
6. Guardar: write_text("nota.txt", "Hola")
   show read_text("nota.txt")
7. Estructuras:
   struct Jugador { nombre: string; puntos: int }
   j = Jugador("Ana", 10)
   j.puntos += 5
   show j.nombre, j.puntos
   assert j.puntos == 15, "Puntos incorrectos"

8. Bibliotecas: import tripol.math
   show math.square(7)
   tripol libs muestra los modulos disponibles.

Guarda el codigo en main.tri.
tripol run main.tri   Ejecuta el programa
tripol build main.tri Crea un ejecutable independiente
`)
}
