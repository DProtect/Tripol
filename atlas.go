package main

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed stdlib/*.tri
var standardSources embed.FS

var standardNames = []string{"animation", "arrays", "binary", "bits", "bitset", "buffer", "camera", "canvas", "checksum", "collections", "colors", "compression", "console", "convert", "datetime", "dictionary", "discord", "encoding", "env", "errors", "events", "files", "game", "geometry", "http", "https", "image", "json", "math", "matrix", "memory", "network", "particles", "path", "physics", "profiler", "queue", "random", "search", "sequences", "simd", "sort", "sprite", "stack", "stats", "system", "testing", "text", "time", "ui", "validate", "vector"}

type stdLoader struct {
	state     map[string]int
	functions []fn
}

func standardLoader() *stdLoader { return &stdLoader{state: map[string]int{}} }
func standardSource(name string) string {
	found := false
	for _, n := range standardNames {
		if name == n {
			found = true
			break
		}
	}
	if !found {
		die("unknown standard library %s; try tripol libs", name)
	}
	b, err := standardSources.ReadFile("stdlib/" + name + ".tri")
	if err != nil {
		die("standard library %s: %v", name, err)
	}
	return string(b)
}
func (l *stdLoader) add(name string) {
	if l.state[name] == 1 {
		die("cyclic standard import %s", name)
	}
	if l.state[name] == 2 {
		return
	}
	l.state[name] = 1
	imports, fs, ss := parseNamed(standardSource(name), "tripol."+name)
	if len(ss) != 0 {
		die("standard library cannot have executable top-level statements")
	}
	for _, dep := range imports {
		if !strings.HasPrefix(dep, "@std:") {
			die("standard library must import standard modules")
		}
		l.add(strings.TrimPrefix(dep, "@std:"))
	}
	l.qualify(name, fs)
	l.state[name] = 2
}
func (l *stdLoader) qualify(name string, fs []fn) {
	names := map[string]bool{}
	for _, f := range fs {
		names[f.name] = true
	}
	// Resolve library-local functions and intrinsic calls independently of user symbols.
	var exprWalk func(*expr)
	exprWalk = func(e *expr) {
		if e == nil {
			return
		}
		if e.op == "call" && !strings.Contains(e.s, ".") {
			if names[e.s] {
				e.s = name + "." + e.s
			} else {
				e.s = "@builtin:" + e.s
			}
		}
		for _, a := range e.args {
			exprWalk(a)
		}
	}
	var stmtWalk func([]stmt)
	stmtWalk = func(ss []stmt) {
		for _, s := range ss {
			exprWalk(s.e)
			exprWalk(s.bound)
			stmtWalk(s.yes)
			stmtWalk(s.no)
			stmtWalk(s.final)
		}
	}
	for _, f := range fs {
		f.name = name + "." + f.name
		stmtWalk(f.body)
		l.functions = append(l.functions, f)
	}
}
func isStandardFunction(name string) bool {
	prefix, _, ok := strings.Cut(name, ".")
	if !ok {
		return false
	}
	for _, n := range standardNames {
		if prefix == n {
			return true
		}
	}
	return false
}
func listLibraries() {
	fmt.Println("TRIPOL Core — built-in standard libraries")
	for _, name := range standardNames {
		_, fs, _ := parseUnit(standardSource(name))
		fmt.Printf("tripol.%s (%d functions)\n", name, len(fs))
	}
	fmt.Println("Example: import tripol.math\n         show math.square(7)\nDetails: tripol docs math")
}
func libraryDocs(name string) {
	name = strings.TrimPrefix(name, "tripol.")
	fmt.Printf("import tripol.%s\n\n", name)
	// Source is the exact offline implementation, including signatures.
	fmt.Print(standardSource(name))
}

// Fixed stack arrays have a reserved length slot, so read-only standard-library
// calls can borrow the stack data. Official routines never return or retain this
// input reference; returning array routines allocate a new dynamic array.
func (c *compiler) promoteFixed(e *expr) bool {
	if e.op != "var" {
		return false
	}
	v, ok := c.lookup(e.s)
	if !ok || !strings.HasPrefix(v.typ, "array:") {
		return false
	}
	a := c.a
	a.imm(int64(arrayLen(v)))
	a.store(v.off - 8)
	a.emit(0x48, 0x8d, 0x85)
	a.u32(v.off - 8)
	return true
}
