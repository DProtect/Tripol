package httptransport

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

func TestSessionMethodsPoolingCookiesAndRedirectPolicy(t *testing.T) {
	var reached atomic.Int32
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached.Add(1)
		w.WriteHeader(204)
	}))
	defer other.Close()
	var connections atomic.Int32
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cross":
			http.Redirect(w, r, other.URL, 302)
		case "/same":
			http.Redirect(w, r, "/session", 302)
		default:
			if r.URL.Query().Get("q") != "a & b" && r.URL.Path == "/query" {
				t.Error("query was not encoded")
			}
			if r.Header.Get("Authorization") != "Bearer fixture-token" || r.Header.Get("X-Api-Key") != "sample" {
				t.Error("request missed credentials")
			}
			w.Header().Set("X-Sample", r.Method)
			observed := ""
			if cookie, err := r.Cookie("session"); err == nil {
				observed = cookie.Value
			}
			http.SetCookie(w, &http.Cookie{Name: "session", Value: "ready", Path: "/"})
			fmt.Fprint(w, observed)
		}
	}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			connections.Add(1)
		}
	}
	server.Start()
	defer server.Close()
	s := newSession(true)
	defer s.client.CloseIdleConnections()
	req := spec{Session: 1, Base: server.URL + "/", URL: "query", Method: "GET", Bearer: "fixture-token", KeyName: "X-Api-Key", KeyValue: "sample", Query: []pair{{"q", "a & b"}}}
	for _, method := range []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"} {
		req.Method = method
		resp := s.send(context.Background(), req, &transfer{})
		if resp.Error != "" || resp.Status != 200 || resp.Headers["x-sample"] != method {
			t.Fatalf("%s: %+v", method, resp)
		}
		if method != "GET" && method != "HEAD" && resp.Text != "ready" {
			t.Fatalf("cookie did not persist: %+v", resp)
		}
	}
	if n := connections.Load(); n != 1 {
		t.Fatalf("expected connection reuse, opened %d", n)
	}
	req.Method, req.URL, req.Redirect = "GET", "same", "same-origin"
	if out := s.send(context.Background(), req, &transfer{}); out.Status != 200 {
		t.Fatal(out)
	}
	req.URL = "cross"
	if out := s.send(context.Background(), req, &transfer{}); out.Status != 302 || reached.Load() != 0 {
		t.Fatalf("followed cross-origin: %+v", out)
	}
	req.URL = other.URL
	if out := s.send(context.Background(), req, &transfer{}); out.Error != "cross_origin" {
		t.Fatal(out)
	}
}

func TestRetriesRejectCertificateAndDoNotReplayWrites(t *testing.T) {
	var retries atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := retries.Add(1)
		w.Header().Set("Retry-After", "0")
		if n == 1 || r.Method == "POST" {
			w.WriteHeader(503)
			return
		}
		w.WriteHeader(204)
	}))
	defer server.Close()
	s := newSession(false)
	defer s.client.CloseIdleConnections()
	q := spec{URL: server.URL, Method: "GET", Retries: 2}
	if out := s.send(context.Background(), q, &transfer{}); out.Status != 204 || out.Attempts != 2 {
		t.Fatal(out)
	}
	q.Method = "POST"
	if out := s.send(context.Background(), q, &transfer{}); out.Status != 503 || out.Attempts != 1 {
		t.Fatal(out)
	}
	q.Headers = []pair{{"X-New-Line", "unsafe\r\nvalue"}}
	if out := s.send(context.Background(), q, &transfer{}); out.Error != "invalid_header" || retries.Load() != 3 {
		t.Fatal(out, retries.Load())
	}

	tlsServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) }))
	tlsServer.EnableHTTP2 = true
	tlsServer.StartTLS()
	defer tlsServer.Close()
	q = spec{URL: tlsServer.URL, Method: "GET", Secure: true, Retries: 3}
	if out := s.send(context.Background(), q, &transfer{}); out.Error != "tls_certificate" || out.Attempts != 1 {
		t.Fatal(out)
	}
	fixture := session{client: tlsServer.Client()}
	if out := fixture.send(context.Background(), q, &transfer{}); out.Error != "" || out.Protocol != "HTTP/2.0" {
		t.Fatal(out)
	}
}

func TestServiceStartsConcurrentTasksAndCancelsOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	started := make(chan string, 2)
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started <- r.URL.Path
		select {
		case <-release:
			fmt.Fprint(w, "done")
		case <-r.Context().Done():
		}
	}))
	defer server.Close()
	var service Service
	defer service.Close()
	start := func(path string) []byte {
		raw, _ := json.Marshal(spec{Method: "GET", URL: server.URL + path})
		reply := service.Handle(ctx, "start", raw)
		if _, err := strconv.Atoi(string(reply)); err != nil {
			t.Fatalf("start %s: %s", path, reply)
		}
		return reply
	}
	first, second := start("/one"), start("/two")
	for i := 0; i < 2; i++ {
		select {
		case <-started:
		case <-ctx.Done():
			t.Fatal("requests did not overlap")
		}
	}
	if out := service.Handle(ctx, "poll", second); string(out) != "0" {
		t.Fatalf("premature result: %s", out)
	}
	service.Handle(ctx, "cancel", first)
	var got result
	if raw := service.Handle(ctx, "await", first); json.Unmarshal(raw, &got) != nil || got.Error != "canceled" {
		t.Fatalf("cancel result %s", raw)
	}
	close(release)
	if raw := service.Handle(ctx, "await", second); json.Unmarshal(raw, &got) != nil || got.Text != "done" || got.Status != 200 {
		t.Fatalf("second result %s", raw)
	}
}

func TestFileTransfersStreamAndKeepDestinationAtomic(t *testing.T) {
	raw := bytes.Repeat([]byte{0, 1, 2, 255}, 512<<10)
	dir := t.TempDir()
	source := filepath.Join(dir, "upload.dat")
	dest := filepath.Join(dir, "download.dat")
	partial := filepath.Join(dir, "partial.dat")
	if err := os.WriteFile(source, raw, 0600); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/short":
			w.Header().Set("Content-Length", "40")
			_, _ = w.Write([]byte("incomplete"))
		case "/download":
			w.Header().Set("Content-Length", strconv.Itoa(len(raw)))
			_, _ = w.Write(raw)
		case "/upload":
			mime, err := r.MultipartReader()
			if err != nil {
				t.Error(err)
				w.WriteHeader(400)
				return
			}
			field, err := mime.NextPart()
			if err != nil {
				t.Error(err)
				w.WriteHeader(400)
				return
			}
			text, err := io.ReadAll(field)
			if err != nil || string(text) != "note & content" || field.FormName() != "note" {
				t.Error("missing form field")
			}
			part, err := mime.NextPart()
			if err != nil {
				t.Error(err)
				w.WriteHeader(400)
				return
			}
			body, err := io.ReadAll(part)
			if err != nil || !bytes.Equal(body, raw) || part.FileName() != "upload.dat" || part.FormName() != "file" {
				t.Error("wrong multipart file")
			}
			w.WriteHeader(204)
		}
	}))
	defer server.Close()
	s := newSession(false)
	defer s.client.CloseIdleConnections()
	ctx := context.Background()
	progress := &transfer{}
	q := spec{URL: server.URL + "/download", Method: "GET", Destination: dest}
	if out := s.send(ctx, q, progress); out.Error != "" || out.Bytes != int64(len(raw)) || progress.total.Load() != out.Bytes || progress.bytes.Load() != out.Bytes {
		t.Fatal(out)
	}
	if written, err := os.ReadFile(dest); err != nil || !bytes.Equal(written, raw) {
		t.Fatal("download was not written intact", err)
	}
	if out := s.send(ctx, q, &transfer{}); out.Error != "file_exists" {
		t.Fatal(out)
	}
	q.URL, q.Destination = server.URL+"/short", partial
	if out := s.send(ctx, q, &transfer{}); out.Error != "transfer_failed" {
		t.Fatal(out)
	}
	if _, err := os.Stat(partial); !os.IsNotExist(err) {
		t.Fatal("partial destination exists", err)
	}
	q = spec{URL: server.URL + "/upload", Method: "POST", Form: []pair{{"note", "note & content"}}, Files: []upload{{"file", source}}}
	if out := s.send(ctx, q, &transfer{}); out.Error != "" || out.Status != 204 || out.Uploaded != int64(len(raw)) || out.UploadTotal != int64(len(raw)) {
		t.Fatal(out)
	}
}

func TestNativeResponseErrorsDoNotRevealURLsAndEnforceLimits(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/timeout" {
			<-r.Context().Done()
			return
		}
		w.WriteHeader(404)
		fmt.Fprint(w, strings.Repeat("a", maxInline+1))
	}))
	defer server.Close()
	s := newSession(false)
	defer s.client.CloseIdleConnections()
	q := spec{URL: server.URL + "/timeout?key=secret", Method: "GET", Timeout: 20}
	if out := s.send(context.Background(), q, &transfer{}); out.Error != "timeout" {
		t.Fatal(out)
	}
	q.URL, q.Timeout = server.URL, 2000
	if out := s.send(context.Background(), q, &transfer{}); out.Error != "body_too_large" || out.Status != 404 {
		t.Fatal(out)
	}
}

func TestNegotiatedCompression(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoding := strings.TrimPrefix(r.URL.Path, "/")
		if !strings.Contains(r.Header.Get("Accept-Encoding"), encoding) {
			t.Error("encoding was not advertised")
		}
		w.Header().Set("Content-Encoding", encoding)
		var writer io.WriteCloser
		switch encoding {
		case "gzip":
			writer = gzip.NewWriter(w)
		case "br":
			writer = brotli.NewWriter(w)
		case "zstd":
			encoded, err := zstd.NewWriter(w)
			if err != nil {
				t.Error(err)
				w.WriteHeader(500)
				return
			}
			writer = encoded
		default:
			return
		}
		fmt.Fprint(writer, "compressed payload")
		writer.Close()
	}))
	defer server.Close()
	client := newSession(false)
	defer client.client.CloseIdleConnections()
	for _, encoding := range []string{"gzip", "br", "zstd"} {
		progress := &transfer{}
		out := client.send(context.Background(), spec{Method: "GET", URL: server.URL + "/" + encoding}, progress)
		if out.Error != "" || out.Text != "compressed payload" || progress.total.Load() != -1 {
			t.Fatalf("%s: %+v", encoding, out)
		}
	}
}
