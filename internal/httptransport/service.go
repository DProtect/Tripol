package httptransport

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"sync/atomic"
)

const maxSessions = 32
const maxTasks = 16

type pending struct {
	session  int
	cancel   context.CancelFunc
	done     chan result
	progress transfer
}

// Service is scoped to one generated program. Handle runs on its bridge
// reader, while requests execute independently of the Gateway and each other.
type Service struct {
	anonymous *session
	sessions  map[int]*session
	tasks     map[int]*pending
	next      int
	workers   sync.WaitGroup
	active    atomic.Int64
}

func (s *Service) Close() {
	for _, task := range s.tasks {
		task.cancel()
	}
	s.workers.Wait()
	for _, client := range s.sessions {
		client.close()
	}
	if s.anonymous != nil {
		s.anonymous.client.CloseIdleConnections()
	}
}

func (s *Service) Handle(ctx context.Context, action string, data []byte) []byte {
	reply := encode
	switch action {
	case "send2", "await2", "poll2":
		action = action[:len(action)-1]
		reply = encodeTextFrame
	}
	switch action {
	case "open":
		if len(s.sessions) >= maxSessions {
			return encode(result{Error: "client_limit"})
		}
		var settings struct {
			Cookies bool
			Path    string
			Base    string
		}
		if json.Unmarshal(data, &settings) != nil && json.Unmarshal(data, &settings.Cookies) != nil {
			return encode(result{Error: "invalid_client"})
		}
		client := newSession(settings.Cookies && settings.Path == "")
		client.stats = &clientMetrics{}
		if settings.Path != "" {
			store, err := openDiskJar(settings.Path, settings.Base)
			if err != nil {
				return encode(result{Error: err.Error()})
			}
			client.store, client.client.Jar = store, store
		}
		if s.sessions == nil {
			s.sessions = make(map[int]*session)
		}
		s.next++
		s.sessions[s.next] = client
		return strconv.AppendInt(nil, int64(s.next), 10)
	case "metrics":
		id, err := strconv.Atoi(string(data))
		client := s.sessions[id]
		if err != nil || client == nil {
			return encode(result{Error: "client_closed"})
		}
		return client.stats.snapshot()
	case "abandon":
		id, _ := strconv.Atoi(string(data))
		s.Discard(id)
		return []byte("1")
	case "close":
		id, err := strconv.Atoi(string(data))
		client := s.sessions[id]
		if err != nil || client == nil {
			return []byte("0")
		}
		for _, task := range s.tasks {
			if task.session == id {
				task.cancel()
			}
		}
		delete(s.sessions, id)
		if client.close() {
			return []byte("1")
		}
		return []byte("0")
	case "send", "start":
		var request spec
		if json.Unmarshal(data, &request) != nil {
			return encode(result{Error: "invalid_request"})
		}
		client := s.sessions[request.Session]
		if request.Session == 0 {
			if s.anonymous == nil {
				s.anonymous = newSession(false)
			}
			client = s.anonymous
		}
		if client == nil {
			return encode(result{Error: "client_closed"})
		}
		if action == "send" {
			return reply(client.send(ctx, request, &transfer{}))
		}
		if len(s.tasks) >= maxTasks || s.active.Load() >= maxTasks {
			return encode(result{Error: "task_limit"})
		}
		if s.tasks == nil {
			s.tasks = make(map[int]*pending)
		}
		child, cancel := context.WithCancel(ctx)
		task := &pending{session: request.Session, cancel: cancel, done: make(chan result, 1)}
		task.progress.total.Store(-1)
		s.next++
		s.tasks[s.next] = task
		s.workers.Add(1)
		s.active.Add(1)
		go func() {
			defer s.workers.Done()
			defer cancel()
			response := client.send(child, request, &task.progress)
			s.active.Add(-1)
			task.done <- response
		}()
		return strconv.AppendInt(nil, int64(s.next), 10)
	case "poll", "await", "progress", "cancel":
		id, err := strconv.Atoi(string(data))
		task := s.tasks[id]
		if err != nil || task == nil {
			return encode(result{Error: "unknown_task"})
		}
		switch action {
		case "cancel":
			task.cancel()
			return []byte("1")
		case "progress":
			data, _ := json.Marshal(struct {
				Bytes       int64 `json:"bytes"`
				Total       int64 `json:"total"`
				Uploaded    int64 `json:"uploaded"`
				UploadTotal int64 `json:"upload_total"`
			}{task.progress.bytes.Load(), task.progress.total.Load(), task.progress.uploaded.Load(), task.progress.uploadTotal.Load()})
			return data
		case "poll":
			select {
			case response := <-task.done:
				delete(s.tasks, id)
				return reply(response)
			default:
				return []byte("0")
			}
		default:
			select {
			case response := <-task.done:
				delete(s.tasks, id)
				return reply(response)
			case <-ctx.Done():
				return encode(result{Error: "canceled"})
			}
		}
	default:
		return encode(result{Error: "unknown_operation"})
	}
}

// Pending reports whether a job is owned by this launcher and still collectible.
func (s *Service) Pending(id int) bool { return s.tasks[id] != nil }

// Cancel requests cancellation while leaving the result available for collection.
func (s *Service) Cancel(id int) bool {
	task := s.tasks[id]
	if task == nil {
		return false
	}
	task.cancel()
	return true
}

// Discard cancels a job and relinquishes its slot after its owner closes.
func (s *Service) Discard(id int) {
	if s.Cancel(id) {
		delete(s.tasks, id)
	}
}

// TakeReady collects a completed job without blocking the bridge.
func (s *Service) TakeReady(id int) []byte {
	task := s.tasks[id]
	if task == nil {
		return nil
	}
	select {
	case response := <-task.done:
		delete(s.tasks, id)
		return encode(response)
	default:
		return nil
	}
}
