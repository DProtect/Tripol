package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
	"reflect"
)

const maxBuildCacheImage = 64 << 20

// Hash the complete parsed project (including expanded imports), not merely
// the entry file or timestamps. Fail closed on a future unsupported AST kind.
func cacheAST(h hash.Hash, v reflect.Value) bool {
	h.Write([]byte{byte(v.Kind())})
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() { h.Write([]byte{0}); return true }
		h.Write([]byte{1})
		return cacheAST(h, v.Elem())
	case reflect.Struct:
		cacheNumber(h, uint64(v.NumField()))
		for i := 0; i < v.NumField(); i++ {
			if !cacheAST(h, v.Field(i)) { return false }
		}
	case reflect.Slice:
		if v.IsNil() { h.Write([]byte{0}); return true }
		h.Write([]byte{1})
		cacheNumber(h, uint64(v.Len()))
		for i := 0; i < v.Len(); i++ {
			if !cacheAST(h, v.Index(i)) { return false }
		}
	case reflect.String:
		cacheNumber(h, uint64(len(v.String())))
		io.WriteString(h, v.String())
	case reflect.Int, reflect.Int64:
		cacheNumber(h, uint64(v.Int()))
	case reflect.Bool:
		if v.Bool() { h.Write([]byte{1}) } else { h.Write([]byte{0}) }
	default:
		return false
	}
	return true
}

func cacheNumber(h hash.Hash, n uint64) {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], n)
	h.Write(b[:])
}

func buildCacheKey(compiler [32]byte, fs []fn, ss []stmt, target, profile string) (string, bool) {
	h := sha256.New()
	io.WriteString(h, "TRIPOL-CACHE-R1\x00")
	h.Write(compiler[:])
	for _, v := range []reflect.Value{reflect.ValueOf(target), reflect.ValueOf(profile), reflect.ValueOf(fs), reflect.ValueOf(ss)} {
		if !cacheAST(h, v) { return "", false }
	}
	return hex.EncodeToString(h.Sum(nil)), true
}

func cacheCompilerDigest() ([32]byte, bool) {
	var digest [32]byte
	p, err := os.Executable()
	if err != nil { return digest, false }
	f, err := os.Open(p)
	if err != nil { return digest, false }
	defer f.Close()
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil { return digest, false }
	copy(digest[:], h.Sum(nil))
	return digest, true
}

func readCachePayload(path string, limit int64) []byte {
	f, err := os.Open(path)
	if err != nil { return nil }
	defer f.Close()
	st, err := f.Stat()
	if err != nil || !st.Mode().IsRegular() || st.Size() <= sha256.Size || st.Size() > limit+sha256.Size { return nil }
	data, err := io.ReadAll(io.LimitReader(f, limit+sha256.Size+1))
	if err != nil || len(data) <= sha256.Size || int64(len(data)) > limit+sha256.Size { return nil }
	image := data[sha256.Size:]
	digest := sha256.Sum256(image)
	if !bytes.Equal(data[:sha256.Size], digest[:]) { return nil }
	return image
}

func readBuildCache(path string, target string) []byte {
	image := readCachePayload(path, maxBuildCacheImage)
	switch target {
	case "linux":
		if !bytes.HasPrefix(image, []byte("\x7fELF")) { return nil }
	case "windows":
		if !bytes.HasPrefix(image, []byte("MZ")) { return nil }
	default:
		return nil
	}
	return image
}

func writeBuildCache(path string, image []byte) {
	if len(image) == 0 || len(image) > maxBuildCacheImage { return }
	if os.MkdirAll(filepath.Dir(path), 0700) != nil { return }
	f, err := os.CreateTemp(filepath.Dir(path), ".tmp-")
	if err != nil { return }
	defer os.Remove(f.Name())
	digest := sha256.Sum256(image)
	_, e1 := f.Write(digest[:])
	_, e2 := f.Write(image)
	e3 := f.Close()
	if e1 == nil && e2 == nil && e3 == nil { _ = os.Rename(f.Name(), path) }
}

// Only called for explicit, non-debug build requests. Cache failures are
// ordinary misses. Never cache compilation errors or diagnostic-trace builds.
func newBuildCache(profile string) *moduleCache {
	if diagnosticTrace || profile == "debug" { return nil }
	dir, err := os.UserCacheDir()
	if err != nil { return nil }
	digest, ok := cacheCompilerDigest()
	if !ok { return nil }
	return &moduleCache{dir: filepath.Join(dir, "TRIPOL", "modules-r1"), compiler: digest}
}

func (c *moduleCache) compile(fs []fn, ss []stmt, target, profile string) []byte {
	if c == nil { return compileParsed(fs, ss, target, profile) }
	key, ok := buildCacheKey(c.compiler, fs, ss, target, profile)
	if !ok { return compileParsed(fs, ss, target, profile) }
	path := filepath.Join(filepath.Dir(c.dir), "compiler-r1", key+".bin")
	if image := readBuildCache(path, target); image != nil { return image }
	image := compileParsed(fs, ss, target, profile)
	writeBuildCache(path, image)
	return image
}
