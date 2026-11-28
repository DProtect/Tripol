//go:build linux || windows

package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestHTTPNativeClientAndSimpleR3(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("requires x86-64")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/users":
			if r.Header.Get("Authorization") != "Bearer fixture-token" || r.URL.Query().Get("q") != "x y" {
				t.Error("incorrect request")
			}
			w.Header().Set("X-Test", "normal")
			fmt.Fprint(w, `{"ok":true}`)
		case "/json":
			data, err := io.ReadAll(r.Body)
			if err != nil || string(data) != `{"name":"Alex","active":true,"other":{"n":42}}` {
				t.Errorf("bad request JSON: %s %v", data, err)
			}
			w.WriteHeader(201)
		case "/download":
			fmt.Fprint(w, "text copied")
		case "/cookie":
			if cookie, err := r.Cookie("test"); err == nil {
				fmt.Fprint(w, cookie.Value)
			} else {
				http.SetCookie(w, &http.Cookie{Name: "test", Value: "saved", Path: "/", MaxAge: 60})
				fmt.Fprint(w, "new")
			}
		default:
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	path := filepath.Join(t.TempDir(), "native.exe")
	dest := filepath.Join(t.TempDir(), "file.txt")
	cookieFile := filepath.Join(t.TempDir(), "cookies.json")
	source := fmt.Sprintf(`import tripol.http
import tripol.https
import tripol.json
url=%q
client=http.client()
client.base_url(url)
client.timeout(5000)
client.bearer("fixture-token")
r=client.request("GET","/users")
r.query("q","x y")
task=r.start()
a=task.await()
show a.status,a.header("X-Test"),a.json()
b=http.post(url+"/json", {"name":"Alex","active":true,"other":{"n":42}})
show b.status
bad=https.get(url+"/users")
show bad.error
r2 = http.download(url+"/download", %q, func(progress) {
    if progress.percent==100 {show f"Succeeded: {progress.bytes}"}
})
show r2.ok()
show client.close()
first=http.client()
first.base_url(url)
first.cookie_file(%q)
from_first=first.get("/cookie")
metrics=first.metrics()
show from_first.text,metrics.completed,metrics.connections
show first.close()
second=http.client()
second.base_url(url)
second.cookie_file(%q)
from_second=second.get("/cookie")
show from_second.text
show second.close()
`, server.URL, dest, cookieFile, cookieFile)
	if err := os.WriteFile(path, compileProfile(source, runtime.GOOS, "turbo"), 0700); err != nil {
		t.Fatal(err)
	}
	if !needsDiscordLauncher(path) {
		t.Fatal("network launcher marker was not linked")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	if err := runNativeCommand(cmd, nil); err != nil {
		t.Fatalf("native program %v: %s", err, out.String())
	}
	want := "200 normal {\"ok\":true}\n201\ninvalid_url\nSucceeded: 11\ntrue\ntrue\nnew 1 1\ntrue\nsaved\ntrue\n"
	if output := strings.ReplaceAll(out.String(), "\r\n", "\n"); output != want {
		t.Fatalf("got %q want %q", output, want)
	}
	if data, err := os.ReadFile(dest); err != nil || string(data) != "text copied" {
		t.Fatal("invalid file", err)
	}
}
