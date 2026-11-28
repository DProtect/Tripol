package discordtransport

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"tripol/internal/localfile"
)

const (
	maxCooldownEntries = 200_000
	maxCooldownFile    = 32 << 20
	maxCooldownSeconds = 365 * 24 * 60 * 60
)

type cooldowns struct {
	mu       sync.Mutex
	entries  map[string]time.Time
	path     string
	fileLock *os.File
	prunedAt time.Time
}

type cooldownRequest struct {
	Command string
	Scope   string
	Subject string
	Seconds int
}

func (s *cooldowns) close() {
	if s.fileLock != nil {
		s.fileLock.Close()
	}
}

func (s *cooldowns) load(path string, now time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if path == "" {
		return errors.New("invalid Discord cooldown store configuration")
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return errors.New("invalid Discord cooldown storage path")
	}
	if s.fileLock != nil && s.path == absolute {
		return nil
	}
	if s.fileLock != nil || len(s.entries) != 0 {
		return errors.New("cooldown storage cannot be changed after use")
	}
	lock, err := localfile.Lock(absolute + ".lock")
	if err != nil {
		return errors.New("Discord cooldown store is already in use or cannot be locked")
	}
	if err := s.readFile(absolute, now); err != nil {
		lock.Close()
		return err
	}
	s.path, s.fileLock = absolute, lock
	return nil
}

func (s *cooldowns) readFile(path string, now time.Time) error {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil || !info.Mode().IsRegular() || info.Size() > maxCooldownFile {
		return errors.New("invalid Discord cooldown store")
	}
	file, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return errors.New("could not read Discord cooldown store")
	}
	defer file.Close()
	actual, err := file.Stat()
	if err != nil || !os.SameFile(info, actual) {
		return errors.New("Discord cooldown file changed before it was opened")
	}
	data, err := io.ReadAll(io.LimitReader(file, maxCooldownFile+1))
	var state struct {
		Version int
		Entries map[string]int64
	}
	if err != nil || len(data) > maxCooldownFile || json.Unmarshal(data, &state) != nil ||
		state.Version != 1 || state.Entries == nil || len(state.Entries) > maxCooldownEntries {
		return errors.New("invalid Discord cooldown store")
	}
	entries := make(map[string]time.Time, len(state.Entries))
	for key, unix := range state.Entries {
		if len(key) > 100 || key == "" {
			return errors.New("invalid Discord cooldown key in store")
		}
		expiry := time.UnixMilli(unix)
		if expiry.After(now) {
			entries[key] = expiry
		}
	}
	s.entries = entries
	return nil
}

func (s *cooldowns) save() error {
	if s.path == "" {
		return nil
	}
	file, err := os.CreateTemp(filepath.Dir(s.path), ".discord-cooldowns-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	state := struct {
		Version int
		Entries map[string]int64
	}{1, make(map[string]int64, len(s.entries))}
	for key, expiry := range s.entries {
		state.Entries[key] = expiry.UnixMilli()
	}
	err = json.NewEncoder(file).Encode(state)
	if err == nil {
		err = file.Sync()
	}
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(file.Name(), s.path)
}

func (s *cooldowns) prune(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pruneLocked(now)
}

func (s *cooldowns) pruneLocked(now time.Time) {
	if len(s.entries) == 0 || now.Sub(s.prunedAt) < time.Minute {
		return
	}
	for key, expiry := range s.entries {
		if !expiry.After(now) {
			delete(s.entries, key)
		}
	}
	s.prunedAt = now
}

func (q cooldownRequest) key() (string, bool) {
	name, guild, ok := strings.Cut(q.Command, "/")
	if !ok || len(name) < 2 || len(name) > 33 || name[0] != 's' || (guild != "" && !snowflake(guild)) {
		return "", false
	}
	for _, c := range name[1:] {
		if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_' || c == '-') {
			return "", false
		}
	}
	switch q.Scope {
	case "user", "channel", "guild":
		if !snowflake(q.Subject) || q.Subject == "0" {
			return "", false
		}
	case "global":
		if q.Subject != "" {
			return "", false
		}
	default:
		return "", false
	}
	return q.Command + ":" + q.Scope + ":" + q.Subject, true
}

func (s *cooldowns) operation(action string, body []byte, now time.Time) (int64, error) {
	if action == "load" {
		var path string
		if json.Unmarshal(body, &path) != nil {
			return 0, errors.New("invalid Discord cooldown storage path")
		}
		return 0, s.load(path, now)
	}
	var request cooldownRequest
	if json.Unmarshal(body, &request) != nil {
		return 0, errors.New("invalid Discord cooldown request")
	}
	key, ok := request.key()
	if !ok {
		return 0, errors.New("invalid Discord cooldown bucket")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pruneLocked(now)
	expiry, exists := s.entries[key]
	left := expiry.Sub(now)
	remaining := int64(0)
	if left > 0 {
		remaining = int64(1 + (left-1)/time.Second)
	}
	switch action {
	case "remaining":
		return remaining, nil
	case "reset":
		if !exists {
			return 0, nil
		}
		delete(s.entries, key)
	case "use":
		if request.Seconds < 1 || request.Seconds > maxCooldownSeconds {
			return 0, errors.New("invalid Discord cooldown duration")
		}
		if remaining > 0 {
			return remaining, nil
		}
		if len(s.entries) >= maxCooldownEntries && !exists {
			return 0, errors.New("Discord cooldown store is full")
		}
		if s.entries == nil {
			s.entries = make(map[string]time.Time)
		}
		s.entries[key] = now.Add(time.Duration(request.Seconds) * time.Second)
	default:
		return 0, errors.New("unsupported Discord cooldown operation")
	}
	if err := s.save(); err != nil {
		if exists {
			s.entries[key] = expiry
		} else {
			delete(s.entries, key)
		}
		return 0, errors.New("could not persist Discord cooldown")
	}
	return 0, nil
}

func (s *cooldowns) respond(action string, body []byte) ([]byte, error) {
	seconds, err := s.operation(action, body, time.Now())
	return strconv.AppendInt(nil, seconds, 10), err
}
