package core

import (
	"bytes"
	pecheck "debug/pe"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestOmegaReduction(t *testing.T) {
	cases := []string{
		`a=ints(0);s=8;for i in range(0){s+=a[i]};show s`,
		`a=ints(5);for i in range(5){a[i]=i-2};s=8;for i in range(5){s+=a[i]};show s`,
		`a=ints(5);a[0]=2;s=8;for i in range(-1){s+=a[i]};show s`,
		`a=ints(5);s=8;for i in range(6){s+=a[i]};show s`,
		`a=ints(3);a[0]=1;s=9223372036854775807;for i in range(4){s+=a[i]};show s`,
		`a=ints(3);a[0]=9223372036854775807;a[1]=1;a[2]=-1;s=0;for i in range(3){s+=a[i]};show s`,
		`a=ints(3);a[0]=-9223372036854775807-1;a[1]=-1;s=0;for i in range(3){s+=a[i]};show s`,
		`a=ints(4);s=3;try{for i in range(5){s+=a[i]}}catch code{show code};show s`,
		`a=ints(5);s=0;func n(){show 7;return 3};for i in range(n()){s+=a[i]};show s`,
		`a=ints(4);s=4;for i in range(s){s+=a[i]};show s`,
		`a=ints(4);s=0;for i in range(1,4){s+=a[i]};show s`,
		`a=[1,2,3];s=0;for i in range(3){s+=a[i]};show s`,
		`func f(a:ints,n:int,s:int)->int {for i in range(n){s+=a[i]};return s};a=ints(3);a[0]=1;show f(a,3,7)`,
	}
	for i, src := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			w, b := hyperRun(t, src, "debug")
			for _, p := range []string{"nano", "balanced", "turbo"} {
				o, e := hyperRun(t, src, p)
				if o != w || e != b {
					t.Fatalf("%s: %q/%v != %q/%v", p, o, e, w, b)
				}
			}
		})
	}
}
func TestOmegaReductionRandom(t *testing.T) {
	r := rand.New(rand.NewSource(20))
	var s, w strings.Builder
	for i := 0; i < 35; i++ {
		n := r.Intn(21)
		k := 0
		if n > 0 {
			k = r.Intn(n + 1)
		}
		fmt.Fprintf(&s, "a%d=ints(%d);s%d=0;", i, n, i)
		var sum int
		for j := 0; j < n; j++ {
			v := r.Intn(200) - 100
			fmt.Fprintf(&s, "a%d[%d]=%d;", i, j, v)
			if j < k {
				sum += v
			}
		}
		fmt.Fprintf(&s, "for i in range(%d){s%d+=a%d[i]};show s%d\n", k, i, i, i)
		fmt.Fprintln(&w, sum)
	}
	o, e := hyperRun(t, s.String(), "turbo")
	if e || o != w.String() {
		t.Fatalf("mismatch %q %v", o, e)
	}
}
func TestOmegaSIMD(t *testing.T) {
	runCase(t, `import tripol.simd
 a=ints(0);b=ints(0);show simd.equal(a,b)
 c=ints(3);d=ints(3);show simd.equal(c,d),simd.equal(a,c)
 c[2]=99;show simd.equal(c,d);d[2]=99;show simd.equal(c,d)
 c[0]=4294967296;show simd.equal(c,d);d[0]=4294967296;show simd.equal(c,d)
 c[1]=-9223372036854775807-1;show simd.equal(c,d)`, "true\ntrue false\nfalse\ntrue\nfalse\ntrue\nfalse\n", false)
}
func TestOmegaDiscordParsing(t *testing.T) {
	runCase(t, `import tripol.discord
 m="{\"id\":\"123\",\"content\":\"!hola\",\"author\":{\"id\":\"456\"}}"
 show discord.message_id(m),discord.matches(m,"!hola"),discord.matches(m,"!no")
 show discord.matches("{\"content\":\"!hola\",\"author\":{\"bot\":true}}","!hola")
 show discord.matches("{\"content\":\"!hola\",\"author\":{},\"webhook_id\":\"42\"}","!hola")
 show discord.message_id("{}"),discord.matches("{}","!hola")
 show discord.channel_path("123")`, "123 true false\nfalse\nfalse\n0 false\n/api/v10/channels/123/messages\n", false)
}
func TestOmegaDiscordPE(t *testing.T) {
	p, e := pecheck.NewFile(bytes.NewReader(compile(`import tripol.discord
discord.listen("123","!hola","Hola")`, "windows")))
	if e != nil {
		t.Fatal(e)
	}
	defer p.Close()
	syms, e := p.ImportedSymbols()
	if e != nil {
		t.Fatal(e)
	}
	all := strings.Join(syms, " ")
	for _, sym := range []string{"GetEnvironmentVariableA:KERNEL32.dll", "WinHttpOpen:WINHTTP.dll", "WinHttpSendRequest:WINHTTP.dll", "WinHttpQueryHeaders:WINHTTP.dll", "Sleep:KERNEL32.dll"} {
		if !strings.Contains(all, sym) {
			t.Fatalf("missing %s: %s", sym, all)
		}
	}
}
