package main

import (
	"bytes"
	pecheck "debug/pe"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func testTCP(t *testing.T, request, response, program string, fail bool) {
	t.Helper()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skipf("loopback unavailable: %v", err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	result := make(chan error, 1)
	go func() {
		c, err := listener.Accept()
		if err != nil {
			result <- err
			return
		}
		defer c.Close()
		c.SetDeadline(time.Now().Add(4 * time.Second))
		in, err := io.ReadAll(c)
		if err != nil {
			result <- err
			return
		}
		if string(in) != request {
			result <- fmt.Errorf("request %q", in)
			return
		}
		for i := 0; i < len(response); i += 3000 {
			end := i + 3000
			if end > len(response) {
				end = len(response)
			}
			if _, err = c.Write([]byte(response[i:end])); err != nil {
				result <- err
				return
			}
		}
		result <- nil
	}()
	source := fmt.Sprintf(program, port)
	p := filepath.Join(t.TempDir(), "test")
	if runtime.GOOS == "windows" {
		p += ".exe"
	}
	if err = os.WriteFile(p, compile(source, runtime.GOOS), 0755); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command(p).CombinedOutput()
	if fail {
		if err == nil || !strings.Contains(string(out), "TCP connection") {
			t.Fatalf("expected TCP error: %v %q", err, out)
		}
	} else if err != nil || string(out) != response+"\n" {
		t.Fatalf("execution %v, output length %d, expected %d: %.100q", err, len(out), len(response)+1, out)
	}
	select {
	case e := <-result:
		if e != nil && !fail {
			t.Fatal(e)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server timeout")
	}
}
func TestConnectTCP(t *testing.T) {
	src := `import tripol.network
show network.exchange(network.ipv4(127,0,0,1),%d,"hello")`
	for _, n := range []int{0, 5, 12000, 65535, 65536} {
		t.Run(fmt.Sprint(n), func(t *testing.T) { testTCP(t, "hello", strings.Repeat("a", n), src, n == 65536) })
	}
	t.Run("nul", func(t *testing.T) { testTCP(t, "hello", "a\x00b", src, true) })
}
func TestConnectHTTP(t *testing.T) {
	testTCP(t, "GET /test HTTP/1.0\r\nHost: localhost\r\nConnection: close\r\n\r\n", "HTTP/1.0 200 OK\r\nContent-Length: 2\r\n\r\nOK", `import tripol.http
show http.get_raw(2130706433,%d,"localhost","/test")`, false)
}
func TestConnectInvalid(t *testing.T) {
	for _, src := range []string{`show tcp_exchange(-1,1,"")`, `show tcp_exchange(1,0,"")`, `show tcp_exchange(4294967296,1,"")`} {
		runCase(t, src, "TRIPOL runtime error: TCP connection, transfer or response limit failed\n", true)
	}
	mustReject(t, func() { compile(`show tcp_exchange("127.0.0.1",1,"")`, runtime.GOOS) })
}
func TestConnectPE(t *testing.T) {
	src := `import tripol.network
show network.exchange(network.loopback(),8080,"hello")`
	b := compile(src, "windows")
	f, e := pecheck.NewFile(bytes.NewReader(b))
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	symbols, e := f.ImportedSymbols()
	if e != nil {
		t.Fatal(e)
	}
	for _, s := range []string{"WSAStartup:WS2_32.dll", "socket:WS2_32.dll", "connect:WS2_32.dll", "setsockopt:WS2_32.dll", "send:WS2_32.dll", "shutdown:WS2_32.dll", "recv:WS2_32.dll", "closesocket:WS2_32.dll", "WSACleanup:WS2_32.dll", "VirtualAlloc:KERNEL32.dll"} {
		found := false
		for _, v := range symbols {
			if v == s {
				found = true
			}
		}
		if !found {
			t.Fatalf("missing %s: %v", s, symbols)
		}
	}
	for _, target := range []string{"linux", "windows"} {
		if !bytes.Equal(compile("show 42", target), compile("import tripol.network\nimport tripol.http\nshow 42", target)) {
			t.Fatal("unused network changes binary")
		}
	}
}
func TestConnectHTTPParsing(t *testing.T) {
	runCase(t, `import tripol.http
r="HTTP/1.1 201 Created\r\nX-Test: ok\r\n\r\nHello"
show http.status(r),http.body(r)`, "201 Hello\n", false)
	runCase(t, `import tripol.http
show http.status("HTTP/1.0 200 OK\r\n\r\n")
show http.body("HTTP/1.0 200 OK\r\n\r\n")`, "200\n\n", false)
}
func TestConnectRefused(t *testing.T) {
	l, e := net.Listen("tcp4", "127.0.0.1:0")
	if e != nil {
		t.Skip(e)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	runCase(t, fmt.Sprintf(`show tcp_exchange(2130706433,%d,"x")`, port), "TRIPOL runtime error: TCP connection, transfer or response limit failed\n", true)
}
func TestConnectFindBytes(t *testing.T) {
	runCase(t, `import tripol.text
show text.index_of("banana","ana"),text.index_of("banana","xyz"),text.index_of("",""),text.index_of("","a"),text.index_of("éa","a"),text.index_of("abc","abcd"),text.contains("abc","bc")`, "1 -1 0 -1 2 -1 true\n", false)
	runCase(t, `import tripol.text
n=0
repeat 1000000 {if text.contains("TRIPOL","IP") {n+=1}}
show n`, "1000000\n", false)
}
