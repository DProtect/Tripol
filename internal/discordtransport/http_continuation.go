package discordtransport

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"tripol/internal/httptransport"
)

type httpContinuation struct {
	ID        int
	Name      string
	Path      string
	expires   time.Time
	delivered bool
}

type httpContinuations map[int]httpContinuation

func (pending httpContinuations) attach(service *httptransport.Service, raw []byte) bool {
	var next httpContinuation
	if json.Unmarshal(raw, &next) != nil || next.ID <= 0 || len(next.Name) == 0 || len(next.Name) > 32 || len(pending) >= 16 ||
		!strings.HasSuffix(next.Path, "/messages/@original") || !service.Pending(next.ID) {
		return false
	}
	for _, b := range next.Name {
		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '_' || b == '-') {
			return false
		}
	}
	if _, err := validateInteractionRoute(http.MethodPatch, next.Path); err != nil {
		return false
	}
	if _, exists := pending[next.ID]; exists {
		return false
	}
	next.expires = time.Now().Add(14 * time.Minute)
	pending[next.ID] = next
	return true
}

func (pending httpContinuations) event(service *httptransport.Service) []byte {
	for id, next := range pending {
		if time.Now().After(next.expires) {
			delete(pending, id)
			continue
		}
		if next.delivered {
			continue
		}
		result := service.TakeReady(id)
		if len(result) == 0 {
			continue
		}
		var frame struct {
			Op   int    `json:"op"`
			Type string `json:"t"`
			Data struct {
				Receipt int             `json:"receipt"`
				Result  json.RawMessage `json:"result"`
			} `json:"d"`
		}
		frame.Type = "TRIPOL_HTTP:" + next.Name
		frame.Data.Receipt = id
		frame.Data.Result = result
		data, _ := json.Marshal(frame)
		if len(data) > maxMessage {
			frame.Data.Result = json.RawMessage(`{"status":0,"error":"response_too_large"}`)
			data, _ = json.Marshal(frame)
		}
		next.delivered = true
		pending[id] = next
		return data
	}
	return nil
}

func (pending httpContinuations) complete(ctx context.Context, client *Client, raw []byte) bool {
	var request struct {
		ID   int
		Body json.RawMessage
	}
	if json.Unmarshal(raw, &request) != nil || len(request.Body) == 0 {
		return false
	}
	continuation, ok := pending[request.ID]
	if !ok || !continuation.delivered || time.Now().After(continuation.expires) {
		return false
	}
	delete(pending, request.ID) // An ambiguous PATCH must not be repeated by a replay.
	_, err := client.EditInteraction(ctx, continuation.Path, request.Body)
	return err == nil
}
