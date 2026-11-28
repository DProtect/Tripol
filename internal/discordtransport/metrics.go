package discordtransport

import (
	"encoding/json"
	"maps"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

type transportMetrics struct {
	connected     atomic.Bool
	heartbeatMS   atomic.Int64
	dispatches    atomic.Uint64
	guildsAtReady atomic.Int64
	queueDepth    atomic.Int64
	connections   atomic.Uint64
	reconnects    atomic.Uint64
	resumes       atomic.Uint64
	gatewayErrors atomic.Uint64
	restAttempts  atomic.Uint64
	restErrors    atomic.Uint64
	rateLimits    atomic.Uint64
	breakdown     atomic.Pointer[dispatchMetrics]
}

type dispatchMetrics struct {
	mu           sync.Mutex
	events       map[string]uint64
	commands     map[string]uint64
	interactions map[string]uint64
}

type Metrics struct {
	UptimeMS      int64             `json:"uptime_ms"`
	Connected     bool              `json:"connected"`
	HeartbeatMS   int64             `json:"heartbeat_ms"`
	Dispatches    uint64            `json:"dispatches"`
	GuildsOnReady int64             `json:"guilds_on_ready"`
	QueueDepth    int64             `json:"event_queue_depth"`
	Connections   uint64            `json:"websocket_connections"`
	Reconnects    uint64            `json:"reconnects"`
	Resumes       uint64            `json:"resumes"`
	GatewayErrors uint64            `json:"gateway_errors"`
	RESTAttempts  uint64            `json:"rest_attempts"`
	RESTErrors    uint64            `json:"rest_errors"`
	RateLimits    uint64            `json:"rate_limits"`
	Events        map[string]uint64 `json:"events_by_type,omitempty"`
	Commands      map[string]uint64 `json:"commands_by_name,omitempty"`
	Interactions  map[string]uint64 `json:"interactions_by_type,omitempty"`
}

// EnableDispatchBreakdown is opt-in so ordinary bots avoid parsing command
// metrics and tracking names on the dispatch path.
func (c *Client) EnableDispatchBreakdown() {
	c.metrics.breakdown.CompareAndSwap(nil, &dispatchMetrics{
		events: make(map[string]uint64), commands: make(map[string]uint64), interactions: make(map[string]uint64),
	})
}

func (m *transportMetrics) recordDispatch(event string, data json.RawMessage) {
	detail := m.breakdown.Load()
	if detail == nil {
		return
	}
	interaction, command := "", ""
	if event == "INTERACTION_CREATE" {
		var decoded struct {
			Type int `json:"type"`
			Data struct {
				Type int    `json:"type"`
				Name string `json:"name"`
			} `json:"data"`
		}
		if json.Unmarshal(data, &decoded) == nil {
			switch decoded.Type {
			case 2:
				interaction = "command"
				name := decoded.Data.Name
				if name != "" && utf8.RuneCountInString(name) <= 32 {
					switch decoded.Data.Type {
					case 1:
						command = "slash/" + name
					case 2:
						command = "user/" + name
					case 3:
						command = "message/" + name
					}
				}
			case 3:
				interaction = "component"
			case 4:
				interaction = "autocomplete"
			case 5:
				interaction = "modal"
			}
		}
	}
	detail.mu.Lock()
	defer detail.mu.Unlock()
	if event != "" && len(event) <= 80 {
		countDispatch(detail.events, event)
	}
	if interaction != "" {
		detail.interactions[interaction]++
	}
	if command != "" {
		countDispatch(detail.commands, command)
	}
}

func countDispatch(counter map[string]uint64, key string) {
	if len(counter) >= 1024 && counter[key] == 0 {
		key = "other"
	}
	counter[key]++
}

// Metrics returns process-local counters with no event payloads or credentials.
func (c *Client) Metrics() Metrics {
	m := &c.metrics
	snapshot := Metrics{
		UptimeMS: time.Since(c.started).Milliseconds(), Connected: m.connected.Load(),
		HeartbeatMS: m.heartbeatMS.Load(), Dispatches: m.dispatches.Load(),
		GuildsOnReady: m.guildsAtReady.Load(), QueueDepth: m.queueDepth.Load(),
		Connections: m.connections.Load(), Reconnects: m.reconnects.Load(),
		Resumes: m.resumes.Load(), GatewayErrors: m.gatewayErrors.Load(),
		RESTAttempts: m.restAttempts.Load(), RESTErrors: m.restErrors.Load(), RateLimits: m.rateLimits.Load(),
	}
	if detail := m.breakdown.Load(); detail != nil {
		detail.mu.Lock()
		snapshot.Events = maps.Clone(detail.events)
		snapshot.Commands = maps.Clone(detail.commands)
		snapshot.Interactions = maps.Clone(detail.interactions)
		detail.mu.Unlock()
	}
	return snapshot
}
