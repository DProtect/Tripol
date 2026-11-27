package discordtransport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type Gateway struct {
	events  chan []byte
	done    chan struct{}
	cancel  context.CancelFunc
	err     error
	once    sync.Once
	metrics *transportMetrics
}

type gatewayOptions struct {
	allowURL      func(*url.URL) bool
	retryDelay    func(int) time.Duration
	identifyDelay time.Duration
	helloDeadline time.Duration
}

func discordGatewayURL(u *url.URL) bool {
	host := u.Hostname()
	return u.Scheme == "wss" && (u.Port() == "" || u.Port() == "443") &&
		(host == "gateway.discord.gg" || strings.HasPrefix(host, "gateway-") && strings.HasSuffix(host, ".discord.gg")) &&
		u.User == nil && (u.Path == "" || u.Path == "/") && u.Fragment == ""
}

func delayFor(attempt int) time.Duration {
	ceiling := time.Second << min(attempt, 5)
	return ceiling/2 + time.Duration(rand.Int64N(int64(ceiling/2)))
}

func (c *Client) StartGateway(parent context.Context, identify []byte) (*Gateway, error) {
	return c.startGateway(parent, identify, gatewayOptions{
		allowURL: discordGatewayURL, retryDelay: delayFor, identifyDelay: 5 * time.Second, helloDeadline: 15 * time.Second,
	})
}

func (c *Client) startGateway(parent context.Context, identify []byte, opt gatewayOptions) (*Gateway, error) {
	var hello struct {
		Op   int `json:"op"`
		Data struct {
			Token   string `json:"token"`
			Intents *int   `json:"intents"`
		} `json:"d"`
	}
	if len(identify) > 4096 || json.Unmarshal(identify, &hello) != nil ||
		hello.Op != 2 || !validToken(hello.Data.Token) || hello.Data.Intents == nil || *hello.Data.Intents < 0 {
		return nil, errors.New("invalid Discord Identify payload")
	}
	ctx, cancel := context.WithCancel(parent)
	g := &Gateway{events: make(chan []byte, 32), done: make(chan struct{}), cancel: cancel, metrics: &c.metrics}
	go func() {
		defer close(g.done)
		g.err = c.supervise(ctx, g.events, append([]byte(nil), identify...), hello.Data.Token, opt)
		if g.err != nil && ctx.Err() == nil {
			c.metrics.gatewayErrors.Add(1)
		}
	}()
	return g, nil
}

func (g *Gateway) Close() {
	g.once.Do(g.cancel)
	<-g.done
}

func (g *Gateway) Poll(ctx context.Context, timeout time.Duration) ([]byte, bool, error) {
	defer func() { g.metrics.queueDepth.Store(int64(len(g.events))) }()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case event := <-g.events:
		return event, false, nil
	case <-g.done:
		select {
		case event := <-g.events:
			return event, false, nil
		default:
			return nil, true, g.err
		}
	case <-timer.C:
		return nil, false, nil
	case <-ctx.Done():
		return nil, true, ctx.Err()
	}
}

type sessionState struct {
	id, resumeURL string
	sequence      *int64
}

func (s *sessionState) canResume() bool { return s.id != "" && s.resumeURL != "" && s.sequence != nil }
func (s *sessionState) clear()          { *s = sessionState{} }

type reconnect struct {
	terminal, invalidate, established bool
	err                               error
}

func (c *Client) supervise(ctx context.Context, output chan<- []byte, identify []byte, token string, opt gatewayOptions) error {
	state := sessionState{}
	var lastIdentify time.Time
	failures := 0
	for ctx.Err() == nil {
		if failures >= 8 {
			return errors.New("Discord Gateway reconnect budget exhausted")
		}
		var endpoint string
		var err error
		if state.canResume() {
			endpoint = state.resumeURL
		} else {
			if err := wait(ctx, time.Until(lastIdentify.Add(opt.identifyDelay))); err != nil {
				return err
			}
			endpoint, err = c.discover(ctx, token, opt)
			if err != nil {
				var apiErr HTTPError
				if errors.As(err, &apiErr) && (apiErr.Status == 401 || apiErr.Status == 403) {
					return err
				}
				if errors.Is(err, errIdentifyExhausted) || errors.Is(err, errGatewayURL) {
					return err
				}
				if err := wait(ctx, opt.retryDelay(failures)); err != nil {
					return err
				}
				failures++
				continue
			}
		}

		if failures > 0 {
			c.metrics.reconnects.Add(1)
		}
		connectedAt := time.Now()
		result := c.gatewayConnection(ctx, endpoint, identify, token, &state, &lastIdentify, output, opt)
		if result.terminal {
			return result.err
		}
		if result.established && time.Since(connectedAt) >= time.Minute {
			failures = 0
		}
		if result.invalidate {
			state.clear()
		}
		if err := wait(ctx, opt.retryDelay(failures)); err != nil {
			return err
		}
		failures++
	}
	return ctx.Err()
}

var (
	errGatewayURL        = errors.New("Discord returned an invalid Gateway URL")
	errIdentifyExhausted = errors.New("Discord Gateway session start quota is exhausted")
)

func gatewayURL(raw string, opt gatewayOptions) (string, error) {
	u, err := url.Parse(raw)
	if err != nil || !opt.allowURL(u) {
		return "", errGatewayURL
	}
	u.RawQuery = "v=10&encoding=json"
	return u.String(), nil
}

func (c *Client) discover(ctx context.Context, token string, opt gatewayOptions) (string, error) {
	payload, err := c.Request(ctx, http.MethodGet, "/api/v10/gateway/bot", token, nil)
	if err != nil {
		return "", err
	}
	var config struct {
		URL   string `json:"url"`
		Start struct {
			Remaining *int `json:"remaining"`
		} `json:"session_start_limit"`
	}
	if json.Unmarshal(payload, &config) != nil || config.Start.Remaining == nil {
		return "", errors.New("invalid Discord Gateway discovery response")
	}
	if *config.Start.Remaining <= 0 {
		return "", errIdentifyExhausted
	}
	return gatewayURL(config.URL, opt)
}

type gatewayMessage struct {
	Op       int             `json:"op"`
	Data     json.RawMessage `json:"d"`
	Sequence *int64          `json:"s"`
	Type     string          `json:"t"`
}
type gatewayFrame struct {
	message gatewayMessage
	raw     []byte
	err     error
}

func readFrame(ctx context.Context, ws *websocket.Conn) gatewayFrame {
	kind, raw, err := ws.Read(ctx)
	if err != nil {
		return gatewayFrame{err: err}
	}
	if kind != websocket.MessageText {
		return gatewayFrame{err: errors.New("Discord Gateway sent a non-text message")}
	}
	var message gatewayMessage
	message.Op = -1
	if json.Unmarshal(raw, &message) != nil || message.Op < 0 {
		return gatewayFrame{err: errors.New("Discord Gateway sent invalid JSON")}
	}
	return gatewayFrame{message: message, raw: raw}
}

func gatewayClose(err error) reconnect {
	code := int(websocket.CloseStatus(err))
	switch code {
	case 4004, 4010, 4011, 4012, 4013, 4014:
		return reconnect{terminal: true, err: fmt.Errorf("Discord Gateway rejected authentication or configuration (%d)", code)}
	case 4007, 4009:
		return reconnect{invalidate: true}
	default:
		return reconnect{}
	}
}

func sendGateway(ctx context.Context, ws *websocket.Conn, payload []byte) error {
	writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return ws.Write(writeCtx, websocket.MessageText, payload)
}

func (c *Client) gatewayConnection(ctx context.Context, endpoint string, identify []byte, token string, state *sessionState, lastIdentify *time.Time, output chan<- []byte, opt gatewayOptions) (result reconnect) {
	handshake, cancelOpen := context.WithTimeout(ctx, 12*time.Second)
	ws, res, err := websocket.Dial(handshake, endpoint, &websocket.DialOptions{HTTPClient: c.http, HTTPHeader: http.Header{"User-Agent": {userAgent}}})
	cancelOpen()
	if err != nil {
		if res != nil && (res.StatusCode == 401 || res.StatusCode == 403) {
			return reconnect{terminal: true, err: errors.New("Discord denied the Gateway connection")}
		}
		return reconnect{}
	}
	c.metrics.connections.Add(1)
	c.metrics.heartbeatMS.Store(-1)
	defer c.metrics.connected.Store(false)
	ws.SetReadLimit(maxMessage)
	connection, cancelConnection := context.WithCancel(ctx)
	var readerDone chan struct{}
	defer func() {
		cancelConnection()
		_ = ws.CloseNow()
		if readerDone != nil {
			<-readerDone
		}
	}()

	hello, cancelHello := context.WithTimeout(connection, opt.helloDeadline)
	first := readFrame(hello, ws)
	cancelHello()
	if first.err != nil {
		return gatewayClose(first.err)
	}
	var interval struct {
		Milliseconds int64 `json:"heartbeat_interval"`
	}
	if first.message.Op != 10 || json.Unmarshal(first.message.Data, &interval) != nil ||
		interval.Milliseconds < 1000 || interval.Milliseconds > 300000 {
		return reconnect{}
	}
	heartbeatEvery := time.Duration(interval.Milliseconds) * time.Millisecond

	initial := identify
	if state.canResume() {
		payload := struct {
			Op   int `json:"op"`
			Data struct {
				Token    string `json:"token"`
				ID       string `json:"session_id"`
				Sequence int64  `json:"seq"`
			} `json:"d"`
		}{Op: 6}
		payload.Data.Token, payload.Data.ID, payload.Data.Sequence = token, state.id, *state.sequence
		initial, _ = json.Marshal(payload)
	} else {
		*lastIdentify = time.Now()
	}
	if sendGateway(connection, ws, initial) != nil {
		return reconnect{}
	}

	frames := make(chan gatewayFrame, 1)
	readerDone = make(chan struct{})
	go func() {
		defer close(readerDone)
		for {
			record := readFrame(connection, ws)
			select {
			case frames <- record:
			case <-connection.Done():
				return
			}
			if record.err != nil {
				return
			}
		}
	}()

	heartbeat := time.NewTimer(time.Duration(rand.Int64N(int64(heartbeatEvery))))
	defer heartbeat.Stop()
	deadline := time.NewTimer(30 * time.Second)
	defer deadline.Stop()
	awaiting := false
	var heartbeatSent time.Time
	beat := func() error {
		body, _ := json.Marshal(struct {
			Op   int    `json:"op"`
			Data *int64 `json:"d"`
		}{Op: 1, Data: state.sequence})
		awaiting = true
		heartbeatSent = time.Now()
		return sendGateway(connection, ws, body)
	}
	for {
		select {
		case <-ctx.Done():
			return reconnect{terminal: true, err: ctx.Err()}
		case <-deadline.C:
			return result
		case <-heartbeat.C:
			if awaiting || beat() != nil {
				return result
			}
			heartbeat.Reset(heartbeatEvery)
		case record := <-frames:
			if record.err != nil {
				disposition := gatewayClose(record.err)
				disposition.established = result.established
				return disposition
			}
			message := record.message
			switch message.Op {
			case 0:
				if message.Sequence == nil || *message.Sequence < 0 {
					return result
				}
				if message.Type == "READY" {
					var ready struct {
						ID     string            `json:"session_id"`
						Resume string            `json:"resume_gateway_url"`
						Guilds []json.RawMessage `json:"guilds"`
					}
					if json.Unmarshal(message.Data, &ready) != nil || len(ready.ID) == 0 || len(ready.ID) > 512 {
						return result
					}
					resume, err := gatewayURL(ready.Resume, opt)
					if err != nil {
						return reconnect{terminal: true, err: err}
					}
					state.id, state.resumeURL = ready.ID, resume
					c.metrics.guildsAtReady.Store(int64(len(ready.Guilds)))
				}
				if message.Type == "READY" || message.Type == "RESUMED" {
					result.established = true
					c.metrics.connected.Store(true)
					if message.Type == "RESUMED" {
						c.metrics.resumes.Add(1)
					}
					deadline.Stop()
				}
				select {
				case output <- record.raw:
					// Only acknowledge an event after admission to the bounded queue.
					state.sequence = message.Sequence
					c.metrics.dispatches.Add(1)
					c.metrics.queueDepth.Store(int64(len(output)))
				default:
					return reconnect{terminal: true, err: errors.New("Discord Gateway event queue is full")}
				}
			case 1:
				if beat() != nil {
					return result
				}
			case 7:
				return result
			case 9:
				var resumable bool
				_ = json.Unmarshal(message.Data, &resumable)
				result.invalidate = !resumable
				return result
			case 11:
				if awaiting {
					c.metrics.heartbeatMS.Store(time.Since(heartbeatSent).Milliseconds())
				}
				awaiting = false
			}
		}
	}
}
