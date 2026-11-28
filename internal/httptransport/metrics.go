package httptransport

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

type clientMetrics struct {
	inflight, completed, failed, canceled, attempts    atomic.Int64
	connections, reused, downloaded, uploaded, elapsed atomic.Int64
}

func (m *clientMetrics) record(out result, elapsed time.Duration) {
	m.attempts.Add(int64(out.Attempts))
	m.downloaded.Add(out.Bytes)
	m.uploaded.Add(out.Uploaded)
	m.elapsed.Add(elapsed.Microseconds())
	if out.Error == "canceled" {
		m.canceled.Add(1)
	}
	if out.Error != "" || out.Status >= 400 {
		m.failed.Add(1)
	}
	m.completed.Add(1)
	m.inflight.Add(-1)
}

func (m *clientMetrics) snapshot() []byte {
	data, _ := json.Marshal(struct {
		InFlight    int64 `json:"in_flight"`
		Completed   int64 `json:"completed"`
		Failed      int64 `json:"failed"`
		Canceled    int64 `json:"canceled"`
		Attempts    int64 `json:"attempts"`
		Connections int64 `json:"connections"`
		Reused      int64 `json:"reused"`
		Downloaded  int64 `json:"downloaded"`
		Uploaded    int64 `json:"uploaded"`
		ElapsedUS   int64 `json:"elapsed_us"`
	}{
		InFlight: m.inflight.Load(), Completed: m.completed.Load(), Failed: m.failed.Load(),
		Canceled: m.canceled.Load(), Attempts: m.attempts.Load(), Connections: m.connections.Load(),
		Reused: m.reused.Load(), Downloaded: m.downloaded.Load(), Uploaded: m.uploaded.Load(), ElapsedUS: m.elapsed.Load(),
	})
	return data
}
