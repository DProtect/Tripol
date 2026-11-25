package main

import (
	"os"
	"path/filepath"
	"strings"
)

// Local modules share the top-level function namespace; no global statements in libraries.
// Paths are resolved relative to the importer, and confined to the root project directory.
func loadProject(source string) ([]fn, []stmt, []string) {
	abs, e := filepath.Abs(source)
	if e != nil {
		die("%v", e)
	}
	abs, e = filepath.EvalSymlinks(abs)
	if e != nil {
		die("%v", e)
	}
	root := filepath.Dir(abs)
	states := map[string]int{}
	libs := standardLoader()
	var fs []fn
	var main []stmt
	var inputs []string
	total := 0
	var visit func(string, bool, int)
	visit = func(path string, entry bool, depth int) {
		if depth > 64 {
			die("module depth limit")
		}
		path, e := filepath.EvalSymlinks(path)
		if e != nil {
			die("module %s: %v", path, e)
		}
		rel, e := filepath.Rel(root, path)
		if e != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			die("module outside project: %s", path)
		}
		if states[path] == 1 {
			die("cyclic import: %s", path)
		}
		if states[path] == 2 {
			return
		}
		states[path] = 1
		if len(states) > 256 {
			die("too many modules")
		}
		b, e := os.ReadFile(path)
		if e != nil {
			die("%v", e)
		}
		total += len(b)
		if total > 4<<20 {
			die("project source too large")
		}
		imports, local, ss := parseUnit(string(b))
		if !entry && len(ss) != 0 {
			die("module %s may contain functions only", path)
		}
		for _, name := range imports {
			if strings.HasPrefix(name, "@std:") {
				libs.add(strings.TrimPrefix(name, "@std:"))
				continue
			}
			if name == "" || filepath.IsAbs(name) || filepath.VolumeName(name) != "" || filepath.Ext(name) != ".tri" || strings.ContainsRune(name, 0) {
				die("invalid local import: %s", name)
			}
			visit(filepath.Join(filepath.Dir(path), name), false, depth+1)
		}
		fs = append(fs, local...)
		if entry {
			main = ss
		}
		inputs = append(inputs, path)
		states[path] = 2
	}
	visit(abs, true, 0)
	return append(libs.functions, fs...), main, inputs
}
