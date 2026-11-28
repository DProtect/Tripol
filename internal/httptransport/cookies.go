package httptransport

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/publicsuffix"
	"tripol/internal/localfile"
)

const maxCookieFile = 1 << 20
const maxPersistentCookies = 512

type cookieSnapshot struct {
	Version int
	Origin  string
	Cookies []http.Cookie
}

// diskJar delegates matching and in-memory behavior to the standard jar.
// The snapshot is confined to the client's base origin and excludes cookies
// without an explicit lifetime.
type diskJar struct {
	jar    *cookiejar.Jar
	base   *url.URL
	path   string
	lock   *os.File
	mu     sync.Mutex
	data   cookieSnapshot
	failed bool
}

func openDiskJar(path, base string) (*diskJar, error) {
	u, err := parseURL(base)
	if err != nil || !absolute(u) || !asciiDomain(u.Hostname()) {
		return nil, errors.New("cookie_store_requires_base")
	}
	if path, err = filepath.Abs(path); err != nil {
		return nil, errors.New("cookie_store_path")
	}
	parent, err := os.Stat(filepath.Dir(path))
	if err != nil || !parent.IsDir() || runtime.GOOS != "windows" && parent.Mode().Perm()&0022 != 0 {
		return nil, errors.New("cookie_store_path")
	}
	lock, err := localfile.Lock(path + ".lock")
	if err != nil {
		return nil, errors.New("cookie_store_locked")
	}
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	// Use a credential-free base when rehydrating cookies.
	baseURL := &url.URL{Scheme: u.Scheme, Host: u.Host, Path: "/"}
	store := &diskJar{jar: jar, base: baseURL, path: path, lock: lock,
		data: cookieSnapshot{Version: 1, Origin: origin(baseURL)}}
	if err := store.load(); err != nil {
		lock.Close()
		return nil, err
	}
	return store, nil
}

func (s *diskJar) load() error {
	info, err := os.Lstat(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil || !info.Mode().IsRegular() || info.Size() > maxCookieFile ||
		runtime.GOOS != "windows" && info.Mode().Perm()&0077 != 0 {
		return errors.New("cookie_store_file")
	}
	file, err := os.Open(s.path)
	if err != nil {
		return errors.New("cookie_store_file")
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return errors.New("cookie_store_file")
	}
	data, err := io.ReadAll(io.LimitReader(file, maxCookieFile+1))
	var loaded cookieSnapshot
	if err != nil || len(data) > maxCookieFile || json.Unmarshal(data, &loaded) != nil ||
		loaded.Version != 1 || len(loaded.Cookies) > maxPersistentCookies {
		return errors.New("cookie_store_format")
	}
	if loaded.Origin != s.data.Origin {
		return errors.New("cookie_store_origin")
	}
	now := time.Now()
	for _, cookie := range loaded.Cookies {
		if _, valid := cookieIdentity(s.base, &cookie); !valid || cookie.Valid() != nil ||
			cookie.MaxAge != 0 || !strings.HasPrefix(cookie.Path, "/") || cookie.Expires.IsZero() ||
			cookie.Partitioned || cookie.Secure && s.base.Scheme != "https" {
			return errors.New("cookie_store_format")
		}
		if cookie.Expires.After(now) {
			s.data.Cookies = append(s.data.Cookies, cookie)
			s.jar.SetCookies(s.base, []*http.Cookie{&cookie})
		}
	}
	return nil
}

func asciiDomain(host string) bool {
	for i := 0; i < len(host); i++ {
		if host[i] < 33 || host[i] >= 127 {
			return false
		}
	}
	return true
}

type cookieKey struct{ domain, path, name string }

func cookieIdentity(u *url.URL, c *http.Cookie) (cookieKey, bool) {
	host := strings.ToLower(strings.TrimSuffix(u.Hostname(), "."))
	domain := strings.ToLower(strings.TrimPrefix(c.Domain, "."))
	switch {
	case c.Domain == "":
		domain = host
	case net.ParseIP(host) != nil:
		if c.Domain != host {
			return cookieKey{}, false
		}
	case domain == "" || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") || !asciiDomain(domain):
		return cookieKey{}, false
	case domain != host:
		if !strings.HasSuffix(host, "."+domain) {
			return cookieKey{}, false
		}
		if suffix, _ := publicsuffix.PublicSuffix(domain); suffix == domain {
			return cookieKey{}, false
		}
	}
	path := c.Path
	if !strings.HasPrefix(path, "/") {
		path = "/"
		if index := strings.LastIndex(u.Path, "/"); index > 0 {
			path = u.Path[:index]
		}
	}
	return cookieKey{domain, path, c.Name}, true
}

func (s *diskJar) Cookies(u *url.URL) []*http.Cookie {
	if origin(u) != s.data.Origin {
		return nil
	}
	return s.jar.Cookies(u)
}

func (s *diskJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if origin(u) != s.data.Origin || len(cookies) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.lock == nil {
		return
	}
	s.jar.SetCookies(u, cookies)
	now := time.Now()
	s.data.Cookies = slices.DeleteFunc(s.data.Cookies, func(c http.Cookie) bool {
		return !c.Expires.After(now)
	})
	changed := false
	rejected := false
	for _, cookie := range cookies {
		key, valid := cookieIdentity(u, cookie)
		if !valid || cookie.Valid() != nil || cookie.Partitioned || cookie.Secure && u.Scheme != "https" {
			continue
		}
		at := slices.IndexFunc(s.data.Cookies, func(saved http.Cookie) bool {
			id, _ := cookieIdentity(s.base, &saved)
			return id == key
		})
		expires := cookie.Expires
		if cookie.MaxAge > 0 {
			const maxAge = 365 * 100 * 24 * 60 * 60
			expires = now.Add(time.Duration(min(int64(cookie.MaxAge), maxAge)) * time.Second)
		}
		if cookie.MaxAge < 0 || !expires.After(now) {
			if at >= 0 {
				s.data.Cookies = slices.Delete(s.data.Cookies, at, at+1)
				changed = true
			}
			continue
		}
		saved := http.Cookie{Name: cookie.Name, Value: cookie.Value, Quoted: cookie.Quoted,
			Path: key.path, Domain: cookie.Domain, Expires: expires, Secure: cookie.Secure,
			HttpOnly: cookie.HttpOnly, SameSite: cookie.SameSite}
		if len(cookie.Name)+len(cookie.Value)+len(cookie.Domain)+len(key.path) > 4096 ||
			at < 0 && len(s.data.Cookies) >= maxPersistentCookies {
			rejected = true
			continue
		}
		if at < 0 {
			s.data.Cookies = append(s.data.Cookies, saved)
		} else {
			s.data.Cookies[at] = saved
		}
		changed = true
	}
	if changed || s.failed {
		s.failed = s.save() != nil
	}
	s.failed = s.failed || rejected
}

func (s *diskJar) save() error {
	data, err := json.Marshal(s.data)
	if err != nil || len(data) > maxCookieFile {
		return errors.New("cookie_store_full")
	}
	file, err := os.CreateTemp(filepath.Dir(s.path), ".tripol-cookies-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	_, err = file.Write(data)
	if err == nil {
		err = file.Sync()
	}
	closeErr := file.Close()
	if err != nil || closeErr != nil {
		return errors.Join(err, closeErr)
	}
	return os.Rename(file.Name(), s.path)
}

func (s *diskJar) failure() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.failed
}

func (s *diskJar) close() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.failed {
		s.failed = s.save() != nil
	}
	err := s.lock.Close()
	s.lock = nil
	return err == nil && !s.failed
}
