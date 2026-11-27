package core

import (
	"tripol/stdlib"
	"fmt"
	"strings"
)

var standardSources = stdlib.Sources

var standardNames = []string{"animation", "arrays", "atomic", "audio", "binary", "bits", "bitset", "buffer", "camera", "canvas", "checksum", "cli", "collections", "colors", "compression", "console", "convert", "crypto", "csv", "datetime", "dictionary", "discord", "encoding", "env", "errors", "events", "files", "game", "geometry", "heap", "http", "https", "image", "io", "json", "math", "matrix", "memory", "network", "particles", "path", "physics", "process", "profiler", "queue", "random", "regex", "search", "sequences", "simd", "sort", "sprite", "stack", "stats", "stream", "system", "testing", "text", "thread", "time", "ui", "utf8", "validate", "vector"}

type stdLoader struct {
	state     map[string]int
	functions []fn
	cache     *moduleCache
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
	b, err := standardSources.ReadFile(name + ".tri")
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
	imports, fs, ss := l.cache.parse([]byte(standardSource(name)), "tripol."+name)
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
	names := map[string]string{}
	records := map[string]string{}
	for _, f := range fs {
		if f.record {
			records[f.name] = name + "." + f.name
		}
	}
	qualifyType := func(typ string) string {
		if args, ret, ok := callbackParts(typ); ok {
			for i, arg := range args {
				if qualified, found := records[arg]; found {
					args[i] = qualified
				}
			}
			if qualified, found := records[ret]; found {
				ret = qualified
			}
			return callbackType(args, ret)
		}
		if qualified, ok := records[typ]; ok {
			return qualified
		}
		return typ
	}
	for _, f := range fs {
		if receiver, method, ok := splitMethodName(f.name); ok {
			names[f.name] = methodName(qualifyType(receiver), method)
		} else {
			names[f.name] = name + "." + f.name
		}
	}
	// Resolve library-local functions and intrinsic calls independently of user symbols.
	var exprWalk func(*expr)
	exprWalk = func(e *expr) {
		if e == nil {
			return
		}
		if e.op == "call" && !strings.Contains(e.s, ".") {
			if qualified, ok := names[e.s]; ok {
				e.s = qualified
			} else {
				e.s = "@builtin:" + e.s
			}
		}
		if e.op == "function-ref" {
			if qualified, ok := names[e.s]; ok {
				e.s = qualified
			}
		}
		for _, a := range e.args {
			exprWalk(a)
		}
	}
	var stmtWalk func([]stmt)
	stmtWalk = func(ss []stmt) {
		for i := range ss {
			s := &ss[i]
			s.typ = qualifyType(s.typ)
			exprWalk(s.e)
			exprWalk(s.bound)
			stmtWalk(s.yes)
			stmtWalk(s.no)
			stmtWalk(s.final)
		}
	}
	for _, f := range fs {
		f.name = names[f.name]
		f.ret = qualifyType(f.ret)
		for i := range f.args {
			f.args[i].typ = qualifyType(f.args[i].typ)
		}
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
		count := 0
		for _, f := range fs {
			if !f.record && !strings.HasPrefix(f.name, methodPrefix) {
				count++
			}
		}
		fmt.Printf("tripol.%s (%d functions)\n", name, count)
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
