package discordtransport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
)

const testIdentify = `{"op":2,"d":{"token":"fixture-token","intents":1}}`

func TestREST429UsesSharedConnection(t *testing.T) {
	var attempts atomic.Int32
	var sent atomic.Int64
	var firstRemote atomic.Value
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost || req.Header.Get("Authorization") != "Bot fixture-token" {
			t.Error("wrong method or authorization")
		}
		got, err := io.ReadAll(req.Body)
		if err != nil || string(got) != `{"content":"hello"}` {
			t.Errorf("unexpected body: %q, %v", got, err)
		}
		attempt := attempts.Add(1)
		if attempt == 1 {
			firstRemote.Store(req.RemoteAddr)
			sent.Store(time.Now().UnixNano())
			w.Header().Set("X-RateLimit-Bucket", "messages")
			w.WriteHeader(429)
			fmt.Fprint(w, `{"retry_after":0.050,"global":true}`)
			return
		}
		if req.RemoteAddr != firstRemote.Load().(string) {
			t.Error("HTTPS connection was not reused")
		}
		if time.Since(time.Unix(0, sent.Load())) < 45*time.Millisecond {
			t.Error("retried before Retry-After")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := newClient(server.Client(), server.URL)
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	got, err := client.Request(ctx, http.MethodPost, "/api/v10/channels/123/messages", "fixture-token", []byte(`{"content":"hello"}`))
	if err != nil || string(got) != "{}" || attempts.Load() != 2 {
		t.Fatalf("attempts=%d response=%q error=%v", attempts.Load(), got, err)
	}
	stats := client.Metrics()
	if stats.RESTAttempts != 2 || stats.RateLimits != 1 || stats.RESTErrors != 0 {
		t.Fatalf("incorrect REST metrics: %+v", stats)
	}
}

func TestRESTRetryBudgetAndNoAmbiguousPOST(t *testing.T) {
	var seen atomic.Int32
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		seen.Add(1)
		switch req.URL.Path {
		case "/api/v10/slow":
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(429)
		case "/api/v10/drop":
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Error(err)
				return
			}
			_ = conn.Close()
		}
	}))
	defer server.Close()
	client := newClient(server.Client(), server.URL)
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for i, path := range []string{"/api/v10/slow", "/api/v10/drop"} {
		if _, err := client.Request(ctx, "POST", path, "fixture-token", []byte("{}")); err == nil || seen.Load() != int32(i+1) {
			t.Fatalf("path=%s attempts=%d error=%v", path, seen.Load(), err)
		}
	}
}

func TestGatewayReconnectResumeAndFatalClose(t *testing.T) {
	var connections, discoveries atomic.Int32
	var target string
	testCtx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	router := http.NewServeMux()
	router.HandleFunc("/api/v10/gateway/bot", func(w http.ResponseWriter, r *http.Request) {
		discoveries.Add(1)
		if r.Header.Get("Authorization") != "Bot fixture-token" {
			t.Error("missing discovery authorization")
		}
		fmt.Fprintf(w, `{"url":%q,"session_start_limit":{"remaining":3}}`, target)
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		number := connections.Add(1)
		ws, err := websocket.Accept(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer ws.CloseNow()
		send := func(payload string) bool {
			if err := ws.Write(testCtx, websocket.MessageText, []byte(payload)); err != nil {
				t.Error(err)
				return false
			}
			return true
		}
		if !send(`{"op":10,"d":{"heartbeat_interval":1000}}`) {
			return
		}
		first := readFrame(testCtx, ws)
		if first.err != nil {
			t.Error(first.err)
			return
		}
		if number == 1 {
			if first.message.Op != 2 {
				t.Errorf("expected Identify, got %d", first.message.Op)
				return
			}
			if !send(fmt.Sprintf(`{"op":0,"t":"READY","s":41,"d":{"session_id":"fixture-session","resume_gateway_url":%q}}`, target)) {
				return
			}
			send(`{"op":7,"d":null}`)
			return
		}
		var resume struct {
			Token   string `json:"token"`
			Session string `json:"session_id"`
			Seq     int    `json:"seq"`
		}
		if err := json.Unmarshal(first.message.Data, &resume); first.message.Op != 6 || err != nil ||
			resume.Token != "fixture-token" || resume.Session != "fixture-session" || resume.Seq != 41 {
			t.Error("invalid Resume")
			return
		}
		if !send(`{"op":0,"t":"RESUMED","s":42,"d":{}}`) {
			return
		}
		if !send(`{"op":1,"d":null}`) {
			return
		}
		beat := readFrame(testCtx, ws)
		// A scheduled beat can race the RESUMED dispatch. The requested
		// heartbeat after that dispatch must report its sequence.
		if beat.err == nil && beat.message.Op == 1 && string(beat.message.Data) == "41" {
			beat = readFrame(testCtx, ws)
		}
		if beat.err != nil || beat.message.Op != 1 || string(beat.message.Data) != "42" {
			t.Errorf("incorrect heartbeat: %+v", beat)
			return
		}
		if !send(`{"op":11,"d":null}`) {
			return
		}
		_ = ws.Close(websocket.StatusCode(4014), "fixture")
	})
	server := httptest.NewTLSServer(router)
	defer server.Close()
	target = "wss" + strings.TrimPrefix(server.URL, "https")
	client := newClient(server.Client(), server.URL)
	defer client.Close()
	g, err := client.startGateway(testCtx, []byte(testIdentify), gatewayOptions{
		allowURL:      func(u *url.URL) bool { return u.String() == target },
		retryDelay:    func(int) time.Duration { return time.Millisecond },
		helloDeadline: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	for _, name := range []string{"READY", "RESUMED"} {
		message, done, err := g.Poll(testCtx, 2*time.Second)
		var got gatewayMessage
		if err != nil || done || json.Unmarshal(message, &got) != nil || got.Type != name {
			t.Fatalf("want %s, got %q (%v, %v)", name, message, done, err)
		}
	}
	_, done, err := g.Poll(testCtx, 2*time.Second)
	if !done || err == nil || !strings.Contains(err.Error(), "4014") || discoveries.Load() != 1 || connections.Load() != 2 {
		t.Fatalf("done=%t err=%v discovery=%d sockets=%d", done, err, discoveries.Load(), connections.Load())
	}
	stats := client.Metrics()
	raw, marshalErr := json.Marshal(stats)
	if marshalErr != nil || stats.Connected || stats.Resumes != 1 || stats.Dispatches != 2 || stats.QueueDepth != 0 ||
		stats.HeartbeatMS < 0 || stats.GatewayErrors != 1 || strings.Contains(string(raw), "fixture-token") ||
		strings.Contains(string(raw), "fixture-session") {
		t.Fatalf("incorrect or sensitive Gateway metrics: %s", raw)
	}
}

func TestDiscordGatewayTrustAndIdentifyQuota(t *testing.T) {
	for _, endpoint := range []string{"ws://gateway.discord.gg", "wss://gateway.discord.gg.attacker.invalid", "wss://gateway.discord.gg@other.invalid", "wss://gateway.discord.gg:444", "wss://gateway.discord.gg/other"} {
		if _, err := gatewayURL(endpoint, gatewayOptions{allowURL: discordGatewayURL}); !errors.Is(err, errGatewayURL) {
			t.Fatalf("accepted %s", endpoint)
		}
	}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `{"url":"wss://gateway.discord.gg","session_start_limit":{"remaining":0}}`)
	}))
	defer server.Close()
	client := newClient(server.Client(), server.URL)
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	g, err := client.StartGateway(ctx, []byte(testIdentify))
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	_, done, err := g.Poll(ctx, time.Second)
	if !done || !errors.Is(err, errIdentifyExhausted) {
		t.Fatalf("done=%t error=%v", done, err)
	}
}
