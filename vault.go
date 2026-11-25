package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type vaultEntry struct {
	Version string `json:"version"`
	SHA256  string `json:"sha256"`
}
type vaultLock struct {
	Format   int                   `json:"format"`
	Packages map[string]vaultEntry `json:"packages"`
}

var vaultName = regexp.MustCompile(`^[a-z][a-z0-9_]{0,47}$`)
var vaultVersion = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)
var vaultHash = regexp.MustCompile(`^[a-f0-9]{64}$`)

func validPackage(name string) bool {
	return vaultName.MatchString(name) && !isStandardFunction(name+".x")
}
func hashSource(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }
func readLock(root string) vaultLock {
	l := vaultLock{1, map[string]vaultEntry{}}
	p := filepath.Join(root, "tripol.lock")
	b, e := os.ReadFile(p)
	if os.IsNotExist(e) {
		return l
	}
	if e != nil {
		die("Vault lock: %v", e)
	}
	if len(b) > 1<<20 || json.Unmarshal(b, &l) != nil || l.Format != 1 || l.Packages == nil || len(l.Packages) > 128 {
		die("invalid Vault lock")
	}
	for n, x := range l.Packages {
		if !validPackage(n) || !vaultVersion.MatchString(x.Version) || len(x.Version) > 64 || !vaultHash.MatchString(x.SHA256) {
			die("invalid Vault package entry %s", n)
		}
	}
	return l
}
func writeLock(root string, l vaultLock) {
	b, e := json.MarshalIndent(l, "", "  ")
	if e != nil {
		die("%v", e)
	}
	f, e := os.CreateTemp(root, ".tripol-lock-")
	if e != nil {
		die("%v", e)
	}
	p := f.Name()
	defer os.Remove(p)
	_, e = f.Write(append(b, '\n'))
	ce := f.Close()
	if e != nil {
		die("%v", e)
	}
	if ce != nil {
		die("%v", ce)
	}
	if e = os.Rename(p, filepath.Join(root, "tripol.lock")); e != nil {
		die("Vault lock update: %v", e)
	}
}
func vaultPath(root, hash string) string {
	return filepath.Join(root, ".tripol", "packages", hash+".tri")
}
func vaultSource(root, name string, x vaultEntry) []byte {
	// Resolve both paths: Windows temporary directories can contain 8.3 aliases.
	root, e := filepath.Abs(root)
	if e != nil {
		die("Vault root: %v", e)
	}
	root, e = filepath.EvalSymlinks(root)
	if e != nil {
		die("Vault root: %v", e)
	}
	p := vaultPath(root, x.SHA256)
	resolved, e := filepath.EvalSymlinks(p)
	if e != nil {
		die("Vault %s: %v", name, e)
	}
	r, e := filepath.Rel(root, resolved)
	if e != nil || r == ".." || strings.HasPrefix(r, ".."+string(filepath.Separator)) {
		die("Vault package outside project")
	}
	b, e := os.ReadFile(resolved)
	if e != nil || len(b) > 256<<10 || hashSource(b) != x.SHA256 {
		die("Vault integrity check failed: %s", name)
	}
	return b
}
func checkPackage(b []byte) ([]string, []fn) {
	if len(b) > 256<<10 {
		die("Vault source exceeds 256 KiB")
	}
	im, fs, ss := parseUnit(string(b))
	if len(ss) != 0 || len(fs) == 0 {
		die("Vault packages require functions only, no top-level execution")
	}
	for _, f := range fs {
		if f.record {
			die("Vault R1 packages cannot export structs yet")
		}
		for _, p := range f.args {
			if !vaultScalar(p.typ) {
				die("Vault R1 invalid parameter type %s", p.typ)
			}
		}
		if !vaultScalar(f.ret) {
			die("Vault R1 invalid return type %s", f.ret)
		}
	}
	for _, n := range im {
		if !strings.HasPrefix(n, "@std:") && !strings.HasPrefix(n, "@pkg:") {
			die("Vault packages require named imports")
		}
	}
	return im, fs
}
func vaultScalar(t string) bool {
	return t == "int" || t == "bool" || t == "float" || t == "string" || t == "ints"
}

type vaultLoader struct {
	root   string
	lock   *vaultLock
	libs   *stdLoader
	state  map[string]int
	total  int
	inputs []string
}

func newVaultLoader(root string, l *stdLoader) *vaultLoader {
	return &vaultLoader{root: root, libs: l, state: map[string]int{}}
}
func (v *vaultLoader) add(name string) {
	if !validPackage(name) {
		die("invalid Vault name %s", name)
	}
	if v.state[name] == 1 {
		die("cyclic Vault import %s", name)
	}
	if v.state[name] == 2 {
		return
	}
	if v.lock == nil {
		l := readLock(v.root)
		v.lock = &l
		v.inputs = append(v.inputs, filepath.Join(v.root, "tripol.lock"))
	}
	x, ok := v.lock.Packages[name]
	if !ok {
		die("Vault package %s not installed; use tripol add", name)
	}
	v.state[name] = 1
	b := vaultSource(v.root, name, x)
	v.inputs = append(v.inputs, vaultPath(v.root, x.SHA256))
	v.total += len(b)
	if v.total > 4<<20 {
		die("Vault source limit exceeded")
	}
	var im []string
	var fs []fn
	file := vaultPath(v.root, x.SHA256)
	attachSource(file, func() { im, fs = checkPackage(b) })
	for i := range fs {
		fs[i].source = file
		tagSource(fs[i].body, file)
	}
	allowed := map[string]bool{name: true}
	for _, dep := range im {
		_, prefix, _ := strings.Cut(dep, ":")
		allowed[prefix] = true
	}
	var checkExpr func(*expr)
	checkExpr = func(e *expr) {
		if e == nil {
			return
		}
		if e.op == "call" {
			if prefix, _, qualified := strings.Cut(e.s, "."); qualified && !allowed[prefix] {
				die("Vault package %s must declare import for %s", name, prefix)
			}
		}
		for _, a := range e.args {
			checkExpr(a)
		}
	}
	var checkStmts func([]stmt)
	checkStmts = func(ss []stmt) {
		for _, s := range ss {
			checkExpr(s.e)
			checkExpr(s.bound)
			checkStmts(s.yes)
			checkStmts(s.no)
			checkStmts(s.final)
		}
	}
	for _, f := range fs {
		checkStmts(f.body)
	}
	for _, dep := range im {
		if strings.HasPrefix(dep, "@std:") {
			v.libs.add(strings.TrimPrefix(dep, "@std:"))
		} else {
			v.add(strings.TrimPrefix(dep, "@pkg:"))
		}
	}
	v.libs.qualify(name, fs)
	v.state[name] = 2
}
func addVault(root, name, source, version string) {
	if !validPackage(name) || !vaultVersion.MatchString(version) || len(version) > 64 {
		die("invalid package name or version (expected x.y.z)")
	}
	b, e := os.ReadFile(source)
	if e != nil {
		die("%v", e)
	}
	checkPackage(b)
	l := readLock(root)
	x := vaultEntry{version, hashSource(b)}
	if old, ok := l.Packages[name]; ok && old.Version == version && old.SHA256 != x.SHA256 {
		die("version already locked to different content; increment version")
	}
	if len(l.Packages) >= 128 {
		if _, ok := l.Packages[name]; !ok {
			die("Vault package limit")
		}
	}
	p := vaultPath(root, x.SHA256)
	if e = os.MkdirAll(filepath.Dir(p), 0755); e != nil {
		die("%v", e)
	}
	if existing, e := os.ReadFile(p); e == nil {
		if !bytes.Equal(existing, b) {
			die("Vault cache integrity failure")
		}
	} else if os.IsNotExist(e) {
		f, e := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if e != nil {
			die("%v", e)
		}
		_, e = f.Write(b)
		ce := f.Close()
		if e != nil || ce != nil {
			os.Remove(p)
			die("Vault cache write failed")
		}
	} else {
		die("%v", e)
	}
	l.Packages[name] = x
	// Verify complete dependency graph and typecheck before changing the lock.
	loader := newVaultLoader(root, standardLoader())
	loader.lock = &l
	for _, n := range sortedVault(l) {
		loader.add(n)
	}
	compileParsed(loader.libs.functions, nil, "linux", "nano")
	writeLock(root, l)
}
func sortedVault(l vaultLock) []string {
	var names []string
	for n := range l.Packages {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
func vaultCLI(args []string) bool {
	if len(args) == 0 {
		return false
	}
	switch args[0] {
	case "add", "remove", "packages", "verify", "bench-build":
	default:
		return false
	}
	if args[0] == "bench-build" {
		benchBuild(args[1:])
		return true
	}
	root, e := os.Getwd()
	if e != nil {
		die("%v", e)
	}
	root, e = filepath.EvalSymlinks(root)
	if e != nil {
		die("%v", e)
	}
	switch args[0] {
	case "add":
		if len(args) != 5 || args[3] != "--version" {
			die("usage: tripol add NAME FILE.tri --version X.Y.Z")
		}
		addVault(root, args[1], args[2], args[4])
		fmt.Printf("Locked %s %s\nImport: import vault.%s\n", args[1], args[4], args[1])
	case "remove":
		if len(args) != 2 {
			die("usage: tripol remove NAME")
		}
		l := readLock(root)
		if _, ok := l.Packages[args[1]]; !ok {
			die("package not installed")
		}
		delete(l.Packages, args[1])
		loader := newVaultLoader(root, standardLoader())
		loader.lock = &l
		for _, n := range sortedVault(l) {
			loader.add(n)
		}
		compileParsed(loader.libs.functions, nil, "linux", "nano")
		writeLock(root, l)
		fmt.Println("Removed", args[1], "(cached source preserved)")
	case "packages", "verify":
		if len(args) != 1 {
			die("unexpected argument")
		}
		l := readLock(root)
		if args[0] == "verify" {
			loader := newVaultLoader(root, standardLoader())
			loader.lock = &l
			for _, n := range sortedVault(l) {
				loader.add(n)
			}
			compileParsed(loader.libs.functions, nil, "linux", "nano")
			fmt.Printf("Verified %d locked packages\n", len(l.Packages))
		} else {
			for _, n := range sortedVault(l) {
				x := l.Packages[n]
				fmt.Printf("%s %s sha256:%s\n", n, x.Version, x.SHA256)
			}
			fmt.Printf("%d local packages\n", len(l.Packages))
		}
	}
	return true
}
func benchBuild(args []string) {
	src := "main.tri"
	runs := 10
	profile := "nano"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		src = args[0]
		args = args[1:]
	}
	for len(args) > 0 {
		if len(args) < 2 {
			die("missing benchmark option value")
		}
		switch args[0] {
		case "--runs":
			n, e := strconv.Atoi(args[1])
			if e != nil || n < 1 || n > 100 {
				die("runs must be 1..100")
			}
			runs = n
		case "--profile":
			profile = args[1]
			if profile != "nano" && profile != "balanced" && profile != "turbo" && profile != "debug" {
				die("unknown profile %s", profile)
			}
		default:
			die("unknown benchmark option %s", args[0])
		}
		args = args[2:]
	}
	build := func() []byte { fs, ss, _ := loadProject(src); return compileParsed(fs, ss, "windows", profile) }
	ref := build()
	sample := func(n int) float64 {
		start := time.Now()
		for i := 0; i < n; i++ {
			if !bytes.Equal(ref, build()) {
				die("non-deterministic compilation")
			}
		}
		return float64(time.Since(start)) / float64(time.Millisecond)
	}
	batch := 1
	for batch < 1024 && sample(batch) < 15 {
		batch *= 2
	}
	times := make([]float64, runs)
	for i := range times {
		times[i] = sample(batch) / float64(batch)
	}
	sort.Float64s(times)
	median := (times[(runs-1)/2] + times[runs/2]) / 2
	fmt.Printf("Build median: %.6f ms | Min: %.6f ms | Max: %.6f ms | Output: %d bytes\n", median, times[0], times[runs-1], len(ref))
	fmt.Printf("Windows x64 %s; %d builds per sample; includes source loading, compilation and deterministic equality check, excludes executable write and compiler startup. Warm-up and calibration excluded.\n", profile, batch)
}
