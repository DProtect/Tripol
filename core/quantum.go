package core

// Mutable standard functions must not borrow fixed (potentially const) arrays.
func mutableArrayFunction(name string) bool {
	return name == "atomic.exchange" || name == "atomic.compare_exchange" || name == "atomic.add" || name == "stream.seek" || name == "stream.read_byte" || name == "stream.read_u16" || name == "heap.push" || name == "heap.pop" || name == "bitset.add" || name == "bitset.remove" || name == "stack.push" || name == "stack.pop" || name == "stack.clear" || name == "queue.push" || name == "queue.pop" || name == "queue.clear" || name == "dictionary.set" || name == "dictionary.remove" || name == "sprite.draw" || name == "sprite.draw_key" || name == "canvas.line" || name == "canvas.circle" || name == "canvas.blend" || name == "particles.step" || name == "sort.sift" || name == "buffer.push" || name == "buffer.set" || name == "buffer.pop" || name == "canvas.clear" || name == "canvas.pixel" || name == "canvas.rectangle"
}

func (c *compiler) quantumBuiltin(e *expr) string {
	signatures := map[string][]string{
		"arena_used": {}, "arena_committed": {}, "array_fill": {"ints", "int"}, "array_fill_range": {"ints", "int", "int", "int"},
		"game_open": {"string", "int", "int"}, "game_running": {}, "game_present": {"ints", "int", "int"}, "game_key": {"int"}, "game_close": {},
	}
	types := signatures[e.s]
	if len(e.args) != len(types) {
		c.err(e.s + " has wrong argument count")
	}
	a := c.a
	for i, x := range e.args {
		c.require(c.expression(x), types[i])
		a.emit(0x50)
	}
	if len(types) == 4 {
		a.emit(0x41, 0x58)
	}
	if len(types) >= 3 {
		a.emit(0x5a)
	}
	if len(types) >= 2 {
		a.emit(0x59)
	}
	if len(types) >= 1 {
		a.emit(0x58)
	}
	a.call(e.s)
	switch e.s {
	case "game_open", "game_running", "game_present", "game_key", "game_close":
		return "bool"
	}
	return "int"
}

const quantumGameExample = `import tripol.game
import tripol.canvas
import tripol.colors
import tripol.events
import tripol.physics

w=320
h=180
pixels=canvas.create(w,h)
x=150
y=80
show game.window("TRIPOL QUANTUM",960,540)
while game.running() {
    if events.pressed(events.escape()) {break}
    if events.pressed(events.left()) or game.key(65) {x-=3}
    if events.pressed(events.right()) or game.key(68) {x+=3}
    if events.pressed(events.up()) or game.key(87) {y-=3}
    if events.pressed(events.down()) or game.key(83) {y+=3}
    x=physics.clamp_position(x,16,w)
    y=physics.clamp_position(y,16,h)
    canvas.clear(pixels,colors.rgb(15,23,42))
    canvas.rectangle(pixels,w,h,0,150,w,30,colors.rgb(30,41,59))
    canvas.rectangle(pixels,w,h,48,30,16,16,colors.rgb(251,191,36))
    canvas.rectangle(pixels,w,h,x,y,16,16,colors.rgb(34,211,238))
    canvas.rectangle(pixels,w,h,x+4,y+4,3,3,16777215)
    canvas.rectangle(pixels,w,h,x+10,y+4,3,3,16777215)
    game.present(pixels,w,h)
    game.wait(16)
}
show game.close()
`

const novaGameExample = `import tripol.game
import tripol.canvas
import tripol.sprite
import tripol.physics
w=1280;h=720
p=canvas.create(w,h)
ship=sprite.solid(48,32,0)
canvas.rectangle(ship,48,32,0,12,48,8,65535)
canvas.rectangle(ship,48,32,14,0,20,32,65535)
canvas.rectangle(ship,48,32,20,4,8,10,16777215)
x=616;y=344
game.window("TRIPOL NOVA",1280,720)
while game.running() {
 if game.key(27) {break}
 if game.key(65) or game.key(37) {x-=8}
 if game.key(68) or game.key(39) {x+=8}
 if game.key(87) or game.key(38) {y-=8}
 if game.key(83) or game.key(40) {y+=8}
 x=physics.clamp_position(x,48,w);y=physics.clamp_position(y,32,h)
 canvas.clear(p,986895)
 canvas.circle(p,w,h,220,150,100,2516582)
 canvas.line(p,w,h,0,600,w-1,600,6710886)
 sprite.draw_key(p,w,h,ship,48,32,x,y,0)
 game.present(p,w,h)
 game.wait(16)
}
game.close()
`
