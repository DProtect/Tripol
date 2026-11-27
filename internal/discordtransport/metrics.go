package discordtransport

import (
	"sync/atomic"
	"time"
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
}

type Metrics struct {
	UptimeMS      int64  `json:"uptime_ms"`
	Connected     bool   `json:"connected"`
	HeartbeatMS   int64  `json:"heartbeat_ms"`
	Dispatches    uint64 `json:"dispatches"`
	GuildsOnReady int64  `json:"guilds_on_ready"`
	QueueDepth    int64  `json:"event_queue_depth"`
	Connections   uint64 `json:"websocket_connections"`
	Reconnects    uint64 `json:"reconnects"`
	Resumes       uint64 `json:"resumes"`
	GatewayErrors uint64 `json:"gateway_errors"`
	RESTAttempts  uint64 `json:"rest_attempts"`
	RESTErrors    uint64 `json:"rest_errors"`
	RateLimits    uint64 `json:"rate_limits"`
}

// Metrics returns process-local counters with no event payloads or credentials.
func (c *Client) Metrics() Metrics {
	m := &c.metrics
	return Metrics{
		UptimeMS: time.Since(c.started).Milliseconds(), Connected: m.connected.Load(),
		HeartbeatMS: m.heartbeatMS.Load(), Dispatches: m.dispatches.Load(),
		GuildsOnReady: m.guildsAtReady.Load(), QueueDepth: m.queueDepth.Load(),
		Connections: m.connections.Load(), Reconnects: m.reconnects.Load(),
		Resumes: m.resumes.Load(), GatewayErrors: m.gatewayErrors.Load(),
		RESTAttempts: m.restAttempts.Load(), RESTErrors: m.restErrors.Load(), RateLimits: m.rateLimits.Load(),
	}
}
