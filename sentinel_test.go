package main

import (
	"bytes"
	pecheck "debug/pe"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func TestJSONRead(t *testing.T) {
	for _, x := range []string{`{"a":[1,true,null,{"b":-1.2e-4}],"c":"hola"}`, `[]`, `{}`, ` "a\u00ff" `, `"é😃"`, `0`, `-0.2E+3`, `true`, `null`, `[1,]`, `{"x":1,}`, `01`, `1.`, `1e`, `{"a" 1}`, `"\q"`, `"\u0x00"`, ``, `true false`, `[[]]`} {
		t.Run(fmt.Sprintf("%q", x), func(t *testing.T) {
			q, _ := json.Marshal(x)
			runCase(t, "import tripol.json\nshow json.valid("+string(q)+")", fmt.Sprintln(json.Valid([]byte(x))), false)
		})
	}
	runCase(t, `import tripol.json
s="{\"a\": [1, 2], \"b\": {\"x\":true}}"
show json.get_raw(s,"a")
show json.get_raw(s,"b")
show json.get_raw(s,"missing")`, "[1, 2]\n{\"x\":true}\n\n", false)
}
func TestTraceFrames(t *testing.T) {
	diagnosticTrace = true
	defer func() { diagnosticTrace = false }()
	runCase(t, `func inner(n){return 7/n}
func outer(n){return inner(n)}
show outer(0)`, "TRIPOL runtime error: invalid integer arithmetic\n  at inner\n  at outer\n  at @main\n", true)
	runCase(t, `func inner(n){return 7/n}
try {show inner(0)} catch e {show e}
throw 100`, "10\nTRIPOL runtime error: unhandled exception\n  at @main\n", true)
}
func TestSentinelPE(t *testing.T) {
	src := `import tripol.ui
import tripol.discord
import tripol.http
import tripol.network
show ui.form("Title","Name","Save")
show discord.gateway()
show http.get_raw(network.loopback(),8765,"localhost","/")`
	p, e := pecheck.NewFile(bytes.NewReader(compile(src, "windows")))
	if e != nil {
		t.Fatal(e)
	}
	syms, e := p.ImportedSymbols()
	if e != nil {
		t.Fatal(e)
	}
	all := strings.Join(syms, " ")
	for _, sym := range []string{"DialogBoxIndirectParamW:USER32.dll", "GetDlgItemTextA:USER32.dll", "GetModuleHandleA:KERNEL32.dll", "WinHttpOpen:WINHTTP.dll", "WinHttpQueryHeaders:WINHTTP.dll", "WSAStartup:WS2_32.dll"} {
		if !strings.Contains(all, sym) {
			t.Fatalf("missing %s in %s", sym, all)
		}
	}
	if runtime.GOOS == "linux" {
		runCase(t, `import tripol.ui
show ui.form("T","L","B")`, "TRIPOL runtime error: native dialogs require Windows or UI operation failed\n", true)
		runCase(t, `import tripol.https
show https.get("example.com","/")`, "TRIPOL runtime error: HTTPS requires Windows, or request failed\n", true)
	}
}
func TestFixSuggestion(t *testing.T) {
	defer func() {
		d, ok := recover().(diagnostic)
		if !ok || d.Suggestion != "numero" {
			t.Fatalf("%#v", d)
		}
	}()
	compile("numero=1\nshow numro", "linux")
}
func TestVectorLength(t *testing.T) {
	for _, n := range []int{0, 1, 15, 16, 17, 31, 32, 33, 255, 4095, 4096, 65535} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			q, _ := json.Marshal(strings.Repeat("a", n))
			runCase(t, "s="+string(q)+"\nshow len(s)\nshow len(s+\"\")", fmt.Sprintf("%d\n%d\n", n, n), false)
		})
	}
}
func TestScratchMemory(t *testing.T) {
	runCase(t, `total=0
repeat 2000 {scratch {a=ints(65535);a[0]=2;total+=a[0]}}
show total
try {scratch {a=ints(65535);throw 105}} catch e {show e}
scratch {scratch {show "nested"}}
show "done"`, "4000\n105\nnested\ndone\n", false)
}
func TestScratchEscape(t *testing.T) {
	for _, s := range []string{`s="a";scratch {s="b"+"c"}`, `a=ints(2);scratch {a[0]=1}`, `struct R {s:string};r=R("a");scratch {r.s="b"}`, `func f()->string {scratch {return "a"};return "b"}`} {
		t.Run(s, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected rejection")
				}
			}()
			compile(s, "linux")
		})
	}
}
