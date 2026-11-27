package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
)

const maxCachedModule = 8 << 20

type moduleCache struct {
	dir      string
	compiler [32]byte
}

// The cache stores the parsed unit before resolving imports or qualifying
// function names. Source locations are attached using the current path.
type cachedExpr struct {
	Op           string
	Text         string `json:",omitempty"`
	Int          int64  `json:",omitempty"`
	Line, Column int    `json:",omitempty"`
	Args         []*cachedExpr
}

type cachedStmt struct {
	Op             string
	Name, Type     string      `json:",omitempty"`
	Expr, Bound    *cachedExpr `json:",omitempty"`
	Yes, No, Final []cachedStmt
	Line           int
	Mutable        bool `json:",omitempty"`
}

type cachedParam struct{ Name, Type string }

type cachedFunction struct {
	Name, Return string
	Args         []cachedParam
	Defaults     []*cachedExpr
	Body         []cachedStmt
	Implicit     []bool
	InferReturn  bool `json:",omitempty"`
	Record       bool `json:",omitempty"`
	Anonymous    bool `json:",omitempty"`
	Line, Column int
}

type cachedModule struct {
	Key        [32]byte
	Imports    []string
	Functions  []cachedFunction
	Statements []cachedStmt
}

func convertSlice[A, B any](in []A, convert func(A) B) []B {
	if in == nil {
		return nil
	}
	out := make([]B, len(in))
	for i, a := range in {
		out[i] = convert(a)
	}
	return out
}

func storeExpr(x *expr) *cachedExpr {
	if x == nil {
		return nil
	}
	return &cachedExpr{Op: x.op, Text: x.s, Int: x.n, Line: x.line, Column: x.col, Args: convertSlice(x.args, storeExpr)}
}

func restoreExpr(x *cachedExpr) *expr {
	if x == nil {
		return nil
	}
	return &expr{op: x.Op, s: x.Text, n: x.Int, line: x.Line, col: x.Column, args: convertSlice(x.Args, restoreExpr)}
}

func storeStmt(s stmt) cachedStmt {
	return cachedStmt{
		Op: s.op, Name: s.name, Type: s.typ, Expr: storeExpr(s.e), Bound: storeExpr(s.bound),
		Yes: convertSlice(s.yes, storeStmt), No: convertSlice(s.no, storeStmt), Final: convertSlice(s.final, storeStmt),
		Line: s.line, Mutable: s.mut,
	}
}

func restoreStmt(s cachedStmt) stmt {
	return stmt{
		op: s.Op, name: s.Name, typ: s.Type, e: restoreExpr(s.Expr), bound: restoreExpr(s.Bound),
		yes: convertSlice(s.Yes, restoreStmt), no: convertSlice(s.No, restoreStmt), final: convertSlice(s.Final, restoreStmt),
		line: s.Line, mut: s.Mutable,
	}
}

func storeFunction(f fn) cachedFunction {
	return cachedFunction{
		Name: f.name, Return: f.ret, Args: convertSlice(f.args, func(a param) cachedParam { return cachedParam{a.name, a.typ} }),
		Defaults: convertSlice(f.defaults, storeExpr), Body: convertSlice(f.body, storeStmt),
		Implicit: f.implicit, InferReturn: f.inferRet, Record: f.record, Anonymous: f.anonymous, Line: f.line, Column: f.col,
	}
}

func restoreFunction(f cachedFunction) fn {
	return fn{
		name: f.Name, ret: f.Return, args: convertSlice(f.Args, func(a cachedParam) param { return param{a.Name, a.Type} }),
		defaults: convertSlice(f.Defaults, restoreExpr), body: convertSlice(f.Body, restoreStmt),
		implicit: f.Implicit, inferRet: f.InferReturn, record: f.Record, anonymous: f.Anonymous, line: f.Line, col: f.Column,
	}
}

func (c *moduleCache) entry(name string, source []byte) (string, [32]byte) {
	pathDigest := sha256.Sum256([]byte(name))
	h := sha256.New()
	h.Write([]byte("TRIPOL-PARSED-MODULE-R2\x00"))
	h.Write(c.compiler[:])
	h.Write(source)
	var key [32]byte
	h.Sum(key[:0])
	path := filepath.Join(c.dir, hex.EncodeToString(pathDigest[:])+".bin")
	return path, key
}

func (c *moduleCache) parse(source []byte, name string) ([]string, []fn, []stmt) {
	if c == nil {
		return parseNamed(string(source), name)
	}
	path, key := c.entry(name, source)
	if data := readCachePayload(path, maxCachedModule); data != nil {
		var stored cachedModule
		if json.Unmarshal(data, &stored) == nil && stored.Key == key {
			fs, ss := convertSlice(stored.Functions, restoreFunction), convertSlice(stored.Statements, restoreStmt)
			for i := range fs {
				fs[i].source = name
				tagSource(fs[i].body, name)
			}
			tagSource(ss, name)
			return stored.Imports, fs, ss
		}
	}
	imports, fs, ss := parseNamed(string(source), name)
	stored := cachedModule{
		Key: key, Imports: imports,
		Functions: convertSlice(fs, storeFunction), Statements: convertSlice(ss, storeStmt),
	}
	if data, err := json.Marshal(stored); err == nil && len(data) <= maxCachedModule {
		writeBuildCache(path, data)
	}
	return imports, fs, ss
}
