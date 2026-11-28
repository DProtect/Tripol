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
	ID             int
	Name           string
	Path           string
	expires        time.Time
	delivered      bool
	delivery       <-chan bool
	cancelDelivery context.CancelFunc
}

type httpContinuations map[int]httpContinuation

func (pending httpContinuations) attach(service *httptransport.Service, raw []byte) bool {
	pending.prune(service)
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
	for _, existing := range pending {
		if existing.Path == next.Path {
			return false
		}
	}
	next.expires = time.Now().Add(14 * time.Minute)
	pending[next.ID] = next
	return true
}

func (pending httpContinuations) prune(service *httptransport.Service) {
	now := time.Now()
	for id, next := range pending {
		if now.After(next.expires) {
			service.Discard(id)
			if next.cancelDelivery != nil {
				next.cancelDelivery()
			}
			delete(pending, id)
		}
	}
}

func (pending httpContinuations) close(service *httptransport.Service) {
	for id, next := range pending {
		service.Discard(id)
		if next.cancelDelivery != nil {
			next.cancelDelivery()
		}
		delete(pending, id)
	}
}

func (pending httpContinuations) cancel(service *httptransport.Service, id int) bool {
	pending.prune(service)
	entry, ok := pending[id]
	return ok && !entry.delivered && service.Cancel(id)
}

func (pending httpContinuations) event(service *httptransport.Service) []byte {
	pending.prune(service)
	for id, next := range pending {
		if next.delivery != nil {
			select {
			case ok := <-next.delivery:
				delete(pending, id)
				data, _ := json.Marshal(map[string]any{"op": 0, "t": "TRIPOL_HTTP_DELIVERY:" + next.Name,
					"d": map[string]any{"receipt": id, "ok": ok}})
				return data
			default:
			}
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

func (pending httpContinuations) edit(ctx context.Context, client *Client, raw []byte, queued bool) bool {
	var request struct {
		ID   int
		Body json.RawMessage
	}
	if json.Unmarshal(raw, &request) != nil || len(request.Body) == 0 || request.Body[0] != '{' {
		return false
	}
	continuation, ok := pending[request.ID]
	if !ok || !continuation.delivered || continuation.delivery != nil || time.Now().After(continuation.expires) {
		return false
	}
	ctx, cancel := context.WithDeadline(ctx, continuation.expires)
	if queued {
		delivery := make(chan bool, 1)
		continuation.delivery, continuation.cancelDelivery = delivery, cancel
		pending[request.ID] = continuation
		go func() {
			defer cancel()
			_, err := client.EditInteraction(ctx, continuation.Path, request.Body)
			delivery <- err == nil
		}()
		return true
	}
	delete(pending, request.ID) // An ambiguous PATCH must not be repeated by a replay.
	defer cancel()
	_, err := client.EditInteraction(ctx, continuation.Path, request.Body)
	return err == nil
}

func (pending httpContinuations) complete(ctx context.Context, client *Client, raw []byte) bool {
	return pending.edit(ctx, client, raw, false)
}
