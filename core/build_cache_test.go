package core

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildCacheKey(t *testing.T) {
	f, s := parse("show 42")
	f2, s2 := parse("show 43")
	var digest [32]byte
	key, ok := buildCacheKey(digest, f, s, "linux", "turbo")
	if !ok || key == "" { t.Fatal("missing key") }
	same, _ := buildCacheKey(digest, f, s, "linux", "turbo")
	if key != same { t.Fatal("nondeterministic key") }
	for _, x := range []struct{ d [32]byte; f []fn; s []stmt; target, profile string }{
		{[32]byte{1},f,s,"linux","turbo"},
		{digest,f2,s2,"linux","turbo"},
		{digest,f,s,"windows","turbo"},
		{digest,f,s,"linux","nano"},
	} {
		k, ok := buildCacheKey(x.d,x.f,x.s,x.target,x.profile)
		if !ok || k == key { t.Fatal("cache key collision") }
	}
}

func TestBuildCacheIntegrity(t *testing.T) {
	p := filepath.Join(t.TempDir(), "cache", "entry.bin")
	image := []byte("\x7fELFtest")
	if readBuildCache(p,"linux") != nil { t.Fatal("unexpected hit") }
	writeBuildCache(p,image)
	if !bytes.Equal(readBuildCache(p,"linux"),image) { t.Fatal("cache round trip") }
	if readBuildCache(p,"windows") != nil { t.Fatal("wrong target accepted") }
	data, err := os.ReadFile(p)
	if err != nil { t.Fatal(err) }
	data[len(data)-1] ^= 1
	if err = os.WriteFile(p,data,0600); err != nil { t.Fatal(err) }
	if readBuildCache(p,"linux") != nil { t.Fatal("corrupt cache accepted") }
}

func TestBuildCacheNativeImages(t *testing.T) {
	f, s := parse("func double(x:int)->int {return x*2};show double(21)")
	for _, target := range []string{"linux", "windows"} {
		for _, profile := range []string{"nano", "balanced", "turbo"} {
			t.Run(target+"/"+profile, func(t *testing.T) {
				image := compileParsed(f,s,target,profile)
				p := filepath.Join(t.TempDir(),"image.bin")
				writeBuildCache(p,image)
				got := readBuildCache(p,target)
				if !bytes.Equal(got,image) { t.Fatal("cached native image differs") }
				if !bytes.Equal(got,compileParsed(f,s,target,profile)) { t.Fatal("non-deterministic native image") }
			})
		}
	}
}
