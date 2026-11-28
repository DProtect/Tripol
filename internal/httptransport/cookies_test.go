package httptransport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

func TestPersistentCookiesAndClientMetrics(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "windows" {
		t.Skip("file locks are implemented for Linux and Windows")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/profile/set":
			http.SetCookie(w, &http.Cookie{Name: "saved", Value: "persistent", MaxAge: 3600})
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "memory", Path: "/profile"})
			http.SetCookie(w, &http.Cookie{Name: "wrong", Value: "foreign", MaxAge: 3600, Domain: "example.com"})
		case "/profile/forget":
			http.SetCookie(w, &http.Cookie{Name: "saved", Value: "session-replacement", Path: "/profile"})
		case "/fail":
			w.WriteHeader(404)
			return
		}
		fmt.Fprint(w, r.Header.Get("Cookie"))
	}))
	defer server.Close()
	var service Service
	defer service.Close()
	ctx := context.Background()
	store := filepath.Join(t.TempDir(), "cookies.json")
	config := func(path, base string) []byte {
		data, _ := json.Marshal(map[string]any{"cookies": true, "path": path, "base": base})
		return data
	}
	open := func(path string) int {
		t.Helper()
		reply := service.Handle(ctx, "open", config(path, server.URL))
		id, err := strconv.Atoi(string(reply))
		if err != nil || id <= 0 {
			t.Fatalf("open: %s", reply)
		}
		return id
	}
	send := func(id int, route string) result {
		t.Helper()
		req, _ := json.Marshal(spec{Session: id, Method: "GET", URL: server.URL + route})
		var out result
		if err := json.Unmarshal(service.Handle(ctx, "send", req), &out); err != nil || out.Error != "" {
			t.Fatalf("send: %v %+v", err, out)
		}
		return out
	}
	closeClient := func(id int) {
		t.Helper()
		if reply := service.Handle(ctx, "close", []byte(strconv.Itoa(id))); string(reply) != "1" {
			t.Fatalf("close: %s", reply)
		}
	}
	first := open(store)
	if reply := string(service.Handle(ctx, "open", config(store, server.URL))); reply != `{"status":0,"error":"cookie_store_locked"}` {
		t.Fatal(reply)
	}
	send(first, "/profile/set")
	same := send(first, "/profile/get")
	if same.Text != "saved=persistent; session=memory" {
		t.Fatal(same)
	}
	send(first, "/fail")
	if send(open(filepath.Join(t.TempDir(), "other.json")), "/profile/get").Text != "" {
		t.Fatal("cookies crossed client stores")
	}
	var metrics map[string]int
	reply := service.Handle(ctx, "metrics", []byte(strconv.Itoa(first)))
	if err := json.Unmarshal(reply, &metrics); err != nil || metrics["completed"] != 3 ||
		metrics["attempts"] != 3 || metrics["failed"] != 1 || metrics["connections"] != 1 ||
		metrics["reused"] != 2 || metrics["downloaded"] != len(same.Text) || metrics["in_flight"] != 0 {
		t.Fatal(string(reply))
	}
	closeClient(first)

	if reply := string(service.Handle(ctx, "open", config(store, "https://example.com"))); reply != `{"status":0,"error":"cookie_store_origin"}` {
		t.Fatal(reply)
	}
	second := open(store)
	if out := send(second, "/other"); out.Text != "" {
		t.Fatal("persisted cookie lost its default path", out)
	}
	if out := send(second, "/profile/get"); out.Text != "saved=persistent" {
		t.Fatal(out)
	}
	send(second, "/profile/forget")
	closeClient(second)
	if out := send(open(store), "/profile/get"); out.Text != "" {
		t.Fatal("the in-memory replacement resurrected a saved credential", out)
	}
}
