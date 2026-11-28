package discordtransport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCooldownIsAtomicAndScoped(t *testing.T) {
	var store cooldowns
	now := time.Now()
	const user = `{"command":"swork/321","scope":"user","subject":"456","seconds":60}`
	var admitted atomic.Int64
	var group sync.WaitGroup
	for i := 0; i < 64; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			left, err := store.operation("use", []byte(user), now)
			if err != nil || left != 60 && left != 0 {
				t.Errorf("decision %d %v", left, err)
			}
			if left == 0 {
				admitted.Add(1)
			}
		}()
	}
	group.Wait()
	if admitted.Load() != 1 {
		t.Fatalf("admitted %d", admitted.Load())
	}
	for _, scope := range []string{"channel", "guild", "global"} {
		subject := "456"
		if scope == "global" {
			subject = ""
		}
		payload := []byte(fmt.Sprintf(`{"command":"swork/321","scope":%q,"subject":%q,"seconds":60}`, scope, subject))
		if left, err := store.operation("use", payload, now); left != 0 || err != nil {
			t.Fatalf("scope %s unexpectedly blocked: %d %v", scope, left, err)
		}
	}
	if left, err := store.operation("remaining", []byte(user), now.Add(59*time.Second+time.Millisecond)); left != 1 || err != nil {
		t.Fatalf("remaining %d %v", left, err)
	}
	if _, err := store.operation("reset", []byte(user), now); err != nil {
		t.Fatal(err)
	}
	if left, err := store.operation("remaining", []byte(user), now); left != 0 || err != nil {
		t.Fatal(left, err)
	}
	store.prune(now.Add(61 * time.Second))
	if len(store.entries) != 0 {
		t.Fatal("expired entries remain")
	}
}

func TestCooldownPersistenceRecoversAndLocks(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "windows" {
		t.Skip("native cooldown persistence")
	}
	const entry = `{"command":"sdaily/","scope":"global","subject":"","seconds":90}`
	payload := []byte(entry)
	path := filepath.Join(t.TempDir(), "data.json")
	now := time.Now()
	first := &cooldowns{}
	defer first.close()
	if err := first.load(path, now); err != nil {
		t.Fatal(err)
	}
	if left, err := first.operation("use", payload, now); left != 0 || err != nil {
		t.Fatal(left, err)
	}
	second := &cooldowns{}
	defer second.close()
	if err := second.load(path, now); err == nil {
		t.Fatal("two processes can write the same store")
	}
	first.close()
	first.fileLock = nil
	if err := second.load(path, now); err != nil {
		t.Fatal(err)
	}
	if left, err := second.operation("remaining", payload, now.Add(20*time.Second)); left != 70 || err != nil {
		t.Fatal(left, err)
	}
	if _, err := second.operation("reset", payload, now); err != nil {
		t.Fatal(err)
	}
	second.close()
	second.fileLock = nil
	third := &cooldowns{}
	defer third.close()
	if err := third.load(path, now); err != nil {
		t.Fatal(err)
	}
	if left, err := third.operation("remaining", payload, now); left != 0 || err != nil {
		t.Fatal(left, err)
	}
	if info, err := os.Stat(path); err != nil || !info.Mode().IsRegular() {
		t.Fatalf("incorrect store type, error %v", err)
	}
}

func TestModerationRejectsOwnerAndHigherRoles(t *testing.T) {
	members := map[string]string{
		"101": `{"roles":["30"]}`,
		"102": `{"roles":["20"]}`,
		"103": `{"roles":["10"]}`,
		"104": `{"roles":[]}`,
		"105": `{"roles":["25"]}`,
		"106": `{"roles":["45"]}`,
		"107": `{"roles":["8"]}`,
	}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.Header.Get("Authorization") != "Bot test-token" {
			t.Error("unexpected moderation preflight method or authorization")
			w.WriteHeader(400)
			return
		}
		switch r.URL.Path {
		case "/api/v10/users/@me":
			fmt.Fprint(w, `{"id":"102"}`)
		case "/api/v10/guilds/100":
			fmt.Fprint(w, `{"owner_id":"104","roles":[{"id":"100","position":0,"permissions":"0"},{"id":"30","position":2,"permissions":"1099511627782"},{"id":"20","position":4,"permissions":"1099511627782"},{"id":"10","position":2,"permissions":"0"},{"id":"25","position":1,"permissions":"8"},{"id":"45","position":1,"permissions":"0"},{"id":"8","position":4,"permissions":"0"}]}`)
		default:
			value, ok := members[filepath.Base(r.URL.Path)]
			if !ok {
				w.WriteHeader(404)
				return
			}
			fmt.Fprint(w, value)
		}
	}))
	defer server.Close()
	client := newClient(server.Client(), server.URL)
	defer client.Close()
	for _, tc := range []struct {
		actor, target string
		permission    int
		want          string
	}{
		{"101", "104", 2, "owner"},
		{"101", "103", 2, "actor_role"},
		{"101", "106", 2, "1"},
		{"104", "106", 1, "1"},
		{"104", "107", 1, "bot_role"},
		{"101", "105", 40, "target_admin"},
		{"106", "103", 1, "actor_permission"},
		{"101", "101", 1, "self"},
	} {
		body, _ := json.Marshal(moderationRequest{"100", tc.actor, tc.target, tc.permission})
		got, err := client.authorizeModeration(context.Background(), "test-token", body)
		if err != nil || got != tc.want {
			t.Errorf("%s -> %s: got %s %v want %s", tc.actor, tc.target, got, err, tc.want)
		}
	}
}

func TestChannelLockPreservesAndRestoresOverwrite(t *testing.T) {
	const id, guild = "888", "999"
	original := overwrite{guild, 0, "3072", "8"}
	var current atomic.Pointer[overwrite]
	current.Store(&original)
	var edits atomic.Int64
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v10/channels/"+id && r.Method == http.MethodGet:
			entries := []overwrite{{"555", 1, "2048", "0"}}
			if item := current.Load(); item != nil {
				entries = append(entries, *item)
			}
			raw, _ := json.Marshal(entries)
			fmt.Fprintf(w, `{"guild_id":"%s","type":0,"permission_overwrites":%s}`, guild, raw)
		case r.URL.Path == "/api/v10/channels/"+id+"/permissions/"+guild && r.Method == http.MethodPut:
			item := overwrite{ID: guild}
			if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
				t.Error(err)
				w.WriteHeader(400)
				return
			}
			current.Store(&item)
			edits.Add(1)
			w.WriteHeader(204)
		case r.URL.Path == "/api/v10/channels/"+id+"/permissions/"+guild && r.Method == http.MethodDelete:
			current.Store(nil)
			edits.Add(1)
			w.WriteHeader(204)
		default:
			t.Error("unexpected lock request")
			w.WriteHeader(400)
		}
	}))
	defer server.Close()
	client := newClient(server.Client(), server.URL)
	defer client.Close()
	var locks channelLocks
	change := func(action string) error { return locks.change(context.Background(), client, id, "test-token", action) }
	if err := change("lock"); err != nil {
		t.Fatal(err)
	}
	frozen := current.Load()
	if frozen.Allow != "1024" || frozen.Deny != "2056" {
		t.Fatalf("wrong locked bits %+v", frozen)
	}
	if err := change("lock"); err != nil || edits.Load() != 1 {
		t.Fatal(edits.Load(), err)
	}
	current.Store(&overwrite{guild, 0, "4096", "2056"})
	if err := change("unlock"); err == nil || edits.Load() != 1 {
		t.Fatal("external edit was overwritten", edits.Load(), err)
	}
	current.Store(frozen)
	if err := change("unlock"); err != nil || *current.Load() != original {
		t.Fatal(current.Load(), err)
	}
	current.Store(nil)
	if err := change("lock"); err != nil {
		t.Fatal(err)
	}
	if err := change("unlock"); err != nil || current.Load() != nil {
		t.Fatal(current.Load(), err)
	}
}

func TestWebhookExecutesMultipartWithoutBotAuthorization(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.txt")
	if err := os.WriteFile(path, []byte("bytes\x00\xff"), 0600); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v10/webhooks/901/test-webhook" ||
			r.URL.Query().Get("wait") != "true" || r.URL.Query().Get("with_components") != "true" ||
			r.Header.Get("Authorization") != "" {
			t.Error("invalid webhook route, method or authorization")
			w.WriteHeader(400)
			return
		}
		if err := r.ParseMultipartForm(2048); err != nil {
			t.Error(err)
			w.WriteHeader(400)
			return
		}
		defer r.MultipartForm.RemoveAll()
		var body struct {
			Flags       int
			Attachments []uploadAttachment
		}
		if json.Unmarshal([]byte(r.FormValue("payload_json")), &body) != nil || body.Flags != 32768 ||
			len(body.Attachments) != 1 || body.Attachments[0].Filename != "report.txt" {
			t.Errorf("invalid attachment metadata: %+v", body)
		}
		f, header, err := r.FormFile("files[0]")
		if err != nil {
			t.Error(err)
			w.WriteHeader(400)
			return
		}
		f.Close()
		if header.Size != 7 || header.Filename != "report.txt" {
			t.Errorf("unexpected attachment %+v", header)
		}
		fmt.Fprint(w, `{"id":"960"}`)
	}))
	defer server.Close()
	client := newClient(server.Client(), server.URL)
	defer client.Close()
	data, _ := json.Marshal(map[string]any{"payload": map[string]int{"flags": 32768}, "files": []string{path}})
	result, err := client.executeWebhook(context.Background(), "901", "test-webhook", data, true)
	if err != nil || string(result) != `{"id":"960"}` {
		t.Fatalf("result %s %v", result, err)
	}
}
