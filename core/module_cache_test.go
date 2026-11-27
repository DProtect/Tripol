package core

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestModuleCacheSelectiveReuse(t *testing.T) {
	dir := t.TempDir()
	cache := &moduleCache{dir: filepath.Join(dir, "cache"), compiler: [32]byte{1}}
	write := func(name, source string) string {
		t.Helper()
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	const dep = "func twice(x:int, scale=2)->int {return x*scale}\n"
	const main = "import \"dep.tri\"\nimport tripol.math\nshow twice(6),math.square(3)"
	path := write("dep.tri", dep)
	entry := write("main.tri", main)
	fs, ss, _ := loadProjectCached(entry, cache)
	canonical, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatal(err)
	}
	entryPath, _ := cache.entry(canonical, []byte(dep))
	proof, err := os.ReadFile(entryPath)
	if err != nil {
		t.Fatal(err)
	}
	cachedFS, cachedSS, _ := loadProjectCached(entry, cache)
	if !reflect.DeepEqual(fs, cachedFS) || !reflect.DeepEqual(ss, cachedSS) {
		t.Fatal("warm cache changed the parsed program")
	}
	write("main.tri", main+"\nshow true")
	reloadedFS, reloadedSS, _ := loadProjectCached(entry, cache)
	if !reflect.DeepEqual(fs, reloadedFS) || reflect.DeepEqual(ss, reloadedSS) {
		t.Fatal("selective edit changed an unrelated module or missed the entry change")
	}
	stored, _ := os.ReadFile(entryPath)
	if !bytes.Equal(proof, stored) {
		t.Fatal("unmodified dependency was replaced")
	}
	for _, target := range []string{"linux", "windows"} {
		f, s, _ := loadProject(entry)
		if !bytes.Equal(compileParsed(f, s, target, "turbo"), cache.compile(reloadedFS, reloadedSS, target, "turbo")) {
			t.Fatalf("%s: the cache changed native output", target)
		}
	}

	write("dep.tri", "func twice(x:int)->int {return x+2}\n")
	changedFS, _, _ := loadProjectCached(entry, cache)
	if reflect.DeepEqual(fs, changedFS) {
		t.Fatal("dependency content change did not invalidate its entry")
	}
	writePath := func(data []byte) {
		t.Helper()
		if err := os.WriteFile(entryPath, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	writePath([]byte("corrupt"))
	repaired, _, _ := loadProjectCached(entry, cache)
	if !reflect.DeepEqual(changedFS, repaired) || readCachePayload(entryPath, maxCachedModule) == nil {
		t.Fatal("corrupt module entry was not rebuilt")
	}
}
