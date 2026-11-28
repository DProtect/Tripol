package discordtransport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
)

func TestUploadMessageStreamsMultipartAndRetriesExplicit429(t *testing.T) {
	const route = "/api/v10/channels/123/messages"
	names := []string{"report.txt", "binary.dat"}
	data := [][]byte{[]byte("informe\n"), {0, 255, '\n', 0}}
	dir := t.TempDir()
	paths := make([]string, 0, len(names))
	for i, name := range names {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, data[i], 0600); err != nil {
			t.Fatal(err)
		}
		paths = append(paths, path)
	}
	var requests atomic.Int32
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := requests.Add(1)
		if r.URL.Path != route || r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bot test-token" || r.ContentLength < 1 {
			t.Error("missing upload route, authorization or length")
		}
		multipart, err := r.MultipartReader()
		if err != nil {
			t.Error(err)
			w.WriteHeader(400)
			return
		}
		part, err := multipart.NextPart()
		if err != nil || part.FormName() != "payload_json" {
			t.Error("missing payload_json")
			w.WriteHeader(400)
			return
		}
		var message struct {
			Attachments []uploadAttachment                    `json:"attachments"`
			Flags       int                                   `json:"flags"`
			Components  []struct{ File struct{ URL string } } `json:"components"`
		}
		if json.NewDecoder(part).Decode(&message) != nil || len(message.Attachments) != 2 ||
			message.Flags != 32768 || len(message.Components) != 1 || message.Components[0].File.URL != "attachment://report.txt" {
			t.Errorf("invalid message metadata: %+v", message)
		}
		for i, name := range names {
			part, err = multipart.NextPart()
			if err != nil {
				t.Error(err)
				break
			}
			body, err := io.ReadAll(part)
			if err != nil || part.FormName() != fmt.Sprintf("files[%d]", i) || part.FileName() != name || !bytes.Equal(body, data[i]) ||
				len(message.Attachments) <= i || message.Attachments[i] != (uploadAttachment{ID: i, Filename: name}) {
				t.Errorf("invalid file part %d", i)
			}
		}
		if _, err := multipart.NextPart(); err != io.EOF {
			t.Errorf("unexpected trailer: %v", err)
		}
		if n == 1 {
			w.Header().Set("Retry-After", "0.01")
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprint(w, `{"retry_after":0.01}`)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id":"642"}`)
	}))
	defer server.Close()
	client := newClient(server.Client(), server.URL)
	defer client.Close()
	create := func(files []string) []byte {
		value := struct {
			Payload json.RawMessage `json:"payload"`
			Files   []string        `json:"files"`
		}{json.RawMessage(`{"flags":32768,"components":[{"type":13,"file":{"url":"attachment://report.txt"}}]}`), files}
		raw, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		return raw
	}
	result, err := client.UploadMessage(context.Background(), route, "test-token", create(paths))
	if err != nil || string(result) != `{"id":"642"}` || requests.Load() != 2 {
		t.Fatalf("%s, %v, %d", result, err, requests.Load())
	}
	for _, files := range [][]string{{dir}, {paths[0], paths[0]}, {filepath.Join(dir, "missing.txt")}} {
		if _, err := client.UploadMessage(context.Background(), route, "test-token", create(files)); err == nil {
			t.Fatal("accepted invalid file selection")
		}
	}
	if _, err := client.UploadMessage(context.Background(), route+"?wait=true", "test-token", create(paths)); err == nil {
		t.Fatal("accepted invalid upload route")
	}
	large, err := os.Create(filepath.Join(dir, "large.bin"))
	if err != nil {
		t.Fatal(err)
	}
	if err := large.Truncate(maxUploadRequest + 1); err != nil {
		large.Close()
		t.Fatal(err)
	}
	large.Close()
	if _, err := client.UploadMessage(context.Background(), route, "test-token", create([]string{large.Name()})); err == nil {
		t.Fatal("accepted oversized file")
	}
	if requests.Load() != 2 {
		t.Fatalf("invalid request reached the server: %d", requests.Load())
	}
}
