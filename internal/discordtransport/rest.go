// Package discordtransport implements the Discord TLS transport used by TRIPOL's Linux launcher.
package discordtransport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	maxMessage = 65535
	apiHost    = "https://discord.com"
	userAgent  = "DiscordBot (https://github.com/Tripollas/Tripol, 0.28.0)"
	maxWait    = 15 * time.Second
)

type Client struct {
	http    *http.Client
	base    string
	started time.Time
	metrics transportMetrics

	mu      sync.Mutex
	global  time.Time
	routes  map[string]string
	blocked map[string]time.Time
}

func New() *Client {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			ForceAttemptHTTP2: true,
			MaxIdleConns:      8, MaxIdleConnsPerHost: 4, IdleConnTimeout: 60 * time.Second,
			TLSHandshakeTimeout: 8 * time.Second, ResponseHeaderTimeout: 10 * time.Second,
		},
		Timeout:       15 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	return newClient(httpClient, apiHost)
}

func newClient(httpClient *http.Client, base string) *Client {
	client := &Client{http: httpClient, base: base, started: time.Now(), routes: make(map[string]string), blocked: make(map[string]time.Time)}
	client.metrics.heartbeatMS.Store(-1)
	return client
}

func (c *Client) Close() { c.http.CloseIdleConnections() }

func validToken(token string) bool {
	if len(token) == 0 || len(token) > 512 {
		return false
	}
	for _, r := range token {
		if !(r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.') {
			return false
		}
	}
	return true
}

type HTTPError struct{ Status int }

func (e HTTPError) Error() string { return fmt.Sprintf("Discord HTTP status %d", e.Status) }

// Request sends a bounded authenticated API request. A POST is retried only
// after an explicit 429 rejection, never after an ambiguous I/O error.
func (c *Client) Request(ctx context.Context, method, path, token string, body []byte) (payload []byte, err error) {
	defer func() {
		if err != nil {
			c.metrics.restErrors.Add(1)
		}
	}()
	target, err := url.ParseRequestURI(path)
	if err != nil || target.IsAbs() || target.Host != "" || !strings.HasPrefix(path, "/api/v10/") ||
		strings.ContainsAny(path, "\\#") || len(path) > 1024 || !validToken(token) || len(body) > maxMessage {
		return nil, errors.New("invalid Discord API request")
	}
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
	default:
		return nil, errors.New("unsupported Discord API method")
	}
	route, major := routeKey(method, target.Path)
	for attempt := 0; attempt < 3; attempt++ {
		if err := c.waitLimit(ctx, route); err != nil {
			return nil, err
		}
		req, err := http.NewRequestWithContext(ctx, method, c.base+path, bytes.NewReader(body))
		if err != nil {
			return nil, errors.New("invalid Discord API request")
		}
		req.Header.Set("Authorization", "Bot "+token)
		req.Header.Set("User-Agent", userAgent)
		if len(body) != 0 {
			req.Header.Set("Content-Type", "application/json")
		}

		c.metrics.restAttempts.Add(1)
		res, err := c.http.Do(req)
		if err != nil {
			return nil, errors.New("Discord HTTPS request failed")
		}
		payload, readErr := io.ReadAll(io.LimitReader(res.Body, maxMessage+1))
		res.Body.Close()
		if readErr != nil || len(payload) > maxMessage {
			return nil, errors.New("invalid or oversized Discord HTTPS response")
		}
		status := res.StatusCode
		if status == http.StatusTooManyRequests {
			c.metrics.rateLimits.Add(1)
		}
		c.updateLimit(route, major, res, payload)
		if status == http.StatusTooManyRequests && attempt < 2 {
			continue
		}
		if status < 200 || status >= 300 {
			return nil, HTTPError{status}
		}
		if len(payload) == 0 {
			return []byte("{}"), nil
		}
		if !json.Valid(payload) {
			return nil, errors.New("invalid Discord JSON response")
		}
		return payload, nil
	}
	return nil, HTTPError{http.StatusTooManyRequests}
}

func routeKey(method, path string) (route, major string) {
	// Discord's major parameters keep guilds, channels and webhooks independent.
	segs := strings.Split(strings.TrimPrefix(path, "/api/v10/"), "/")
	if len(segs) >= 2 {
		switch segs[0] {
		case "channels", "guilds", "webhooks":
			major = segs[0] + "/" + segs[1]
			if segs[0] == "webhooks" && len(segs) >= 3 {
				segs[2] = ":token"
			}
		case "interactions":
			if len(segs) >= 3 {
				segs[2] = ":token"
			}
		}
	}
	for i, s := range segs {
		if s == "" {
			continue
		}
		if strings.IndexFunc(s, func(r rune) bool { return r < '0' || r > '9' }) == -1 {
			segs[i] = ":id"
		}
	}
	return method + " " + strings.Join(segs, "/") + " " + major, major
}

func seconds(text string) time.Duration {
	value, err := strconv.ParseFloat(text, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) || value <= 0 {
		return 0
	}
	if value > 86400 {
		value = 86400
	}
	return time.Duration(value * float64(time.Second))
}

func wait(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Client) waitLimit(ctx context.Context, route string) error {
	for {
		c.mu.Lock()
		key := c.routes[route]
		if key == "" {
			key = route
		}
		until := c.blocked[key]
		if c.global.After(until) {
			until = c.global
		}
		c.mu.Unlock()
		delay := time.Until(until)
		if delay <= 0 {
			return ctx.Err()
		}
		if delay > maxWait {
			return errors.New("Discord retry interval exceeds this request's wait budget")
		}
		if err := wait(ctx, delay); err != nil {
			return err
		}
	}
}

func (c *Client) updateLimit(route, major string, res *http.Response, body []byte) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	key := c.routes[route]
	if key == "" {
		key = route
	}
	if bucket := res.Header.Get("X-RateLimit-Bucket"); bucket != "" && len(bucket) <= 128 {
		if len(c.routes) >= 4096 {
			c.collapseLimits()
			clear(c.routes)
		}
		key = bucket + " " + major
		c.routes[route] = key
	}
	delay := seconds(res.Header.Get("X-RateLimit-Reset-After"))
	global := res.Header.Get("X-RateLimit-Global") == "true"
	if res.StatusCode == http.StatusTooManyRequests {
		var reply struct {
			RetryAfter json.Number `json:"retry_after"`
			Global     bool        `json:"global"`
		}
		_ = json.Unmarshal(body, &reply)
		delay = max(delay, seconds(reply.RetryAfter.String()), seconds(res.Header.Get("Retry-After")))
		if delay == 0 {
			delay = time.Second
		}
		global = global || reply.Global || res.Header.Get("X-RateLimit-Scope") == "global"
	} else if res.Header.Get("X-RateLimit-Remaining") != "0" {
		return
	}
	until := now.Add(delay)
	if global {
		if until.After(c.global) {
			c.global = until
		}
		return
	}
	if len(c.blocked) >= 4096 {
		for k, expiry := range c.blocked {
			if !expiry.After(now) {
				delete(c.blocked, k)
			}
		}
		if len(c.blocked) >= 4096 {
			// Conservatively honor all existing blocks when the route cache is full.
			c.collapseLimits()
		}
	}
	if until.After(c.blocked[key]) {
		c.blocked[key] = until
	}
}

func (c *Client) collapseLimits() {
	for _, expiry := range c.blocked {
		if expiry.After(c.global) {
			c.global = expiry
		}
	}
	clear(c.blocked)
}
