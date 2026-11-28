package discordtransport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"tripol/internal/httptransport"
)

func TestHTTPContinuationDoesNotBlockAndCanCompleteOnce(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var edits atomic.Int32
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v10/webhooks/123/callback-fixture/messages/@original" || r.Header.Get("Authorization") != "" {
			t.Error("misdirected continuation")
		}
		edits.Add(1)
		w.WriteHeader(204)
	}))
	defer server.Close()
	external := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		select {
		case <-release:
			fmt.Fprint(w, "local profile")
		case <-r.Context().Done():
		}
	}))
	defer external.Close()
	var service httptransport.Service
	defer service.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	job := service.Handle(ctx, "start", []byte(fmt.Sprintf(`{"method":"GET","url":%q}`, external.URL)))
	id, err := strconv.Atoi(string(job))
	if err != nil {
		t.Fatal(string(job), err)
	}
	pending := make(httpContinuations)
	raw := []byte(fmt.Sprintf(`{"id":%d,"name":"profile","path":"/api/v10/webhooks/123/callback-fixture/messages/@original"}`, id))
	if !pending.attach(&service, raw) {
		t.Fatal("attachment failed")
	}
	select {
	case <-started:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	if event := pending.event(&service); event != nil {
		t.Fatalf("premature event %s", event)
	}
	close(release)
	var event []byte
	for len(event) == 0 {
		event = pending.event(&service)
		if ctx.Err() != nil {
			t.Fatal(ctx.Err())
		}
		if len(event) == 0 {
			time.Sleep(time.Millisecond)
		}
	}
	if !json.Valid(event) || !strings.Contains(string(event), "local profile") || !strings.Contains(string(event), "TRIPOL_HTTP:profile") ||
		strings.Contains(string(event), "callback-fixture") {
		t.Fatalf("invalid continuation event %s", event)
	}
	if duplicate := pending.event(&service); duplicate != nil {
		t.Fatal("duplicate completion dispatched")
	}
	client := newClient(server.Client(), server.URL)
	defer client.Close()
	finish := []byte(fmt.Sprintf(`{"id":%d,"body":{"content":"Resolved"}}`, id))
	if !pending.complete(ctx, client, finish) || pending.complete(ctx, client, finish) || edits.Load() != 1 {
		t.Fatal("a deferred interaction must be edited only once")
	}
}
