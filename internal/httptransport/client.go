// Package httptransport serves native HTTP requests through the TRIPOL launcher.
package httptransport

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"golang.org/x/net/publicsuffix"
)

const maxFrame = 65535
const maxInline = 48 << 10
const defaultTimeout = 15 * time.Second

type pair struct {
	Name  string
	Value string
}

type upload struct {
	Name string
	Path string
}

type spec struct {
	Session     int
	URL         string
	Base        string
	Method      string
	Headers     []pair
	Query       []pair
	Form        []pair
	Files       []upload
	Body        string
	Encoding    string
	Bearer      string
	User        string
	Password    string
	Basic       bool
	KeyName     string
	KeyValue    string
	Timeout     int
	Redirect    string
	Retries     int
	Secure      bool
	Destination string
}

type result struct {
	Status      int               `json:"status"`
	Text        string            `json:"text,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Protocol    string            `json:"protocol,omitempty"`
	Error       string            `json:"error,omitempty"`
	Bytes       int64             `json:"bytes,omitempty"`
	Elapsed     int64             `json:"elapsed_ms,omitempty"`
	Attempts    int               `json:"attempts,omitempty"`
	File        string            `json:"file,omitempty"`
	Uploaded    int64             `json:"uploaded,omitempty"`
	UploadTotal int64             `json:"upload_total,omitempty"`
}

type transfer struct {
	bytes       atomic.Int64
	total       atomic.Int64
	uploaded    atomic.Int64
	uploadTotal atomic.Int64
}

type session struct {
	client *http.Client
}

func newSession(cookies bool) *session {
	transport := &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		DialContext:       (&net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		ForceAttemptHTTP2: true,
		MaxIdleConns:      16, MaxIdleConnsPerHost: 4, MaxConnsPerHost: 8,
		IdleConnTimeout: 60 * time.Second, TLSHandshakeTimeout: 8 * time.Second,
		MaxResponseHeaderBytes: 64 << 10, ExpectContinueTimeout: time.Second,
		TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}
	c := &http.Client{Transport: transport}
	if cookies {
		jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		c.Jar = jar
	}
	return &session{client: c}
}

func origin(u *url.URL) string {
	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	return u.Scheme + "://" + net.JoinHostPort(strings.ToLower(u.Hostname()), port)
}

func parseURL(text string) (*url.URL, error) {
	u, err := url.Parse(text)
	if err != nil || u.User != nil || u.Fragment != "" || u.Opaque != "" || strings.ContainsAny(text, "\\\r\n\t") {
		return nil, errors.New("invalid_url")
	}
	return u, nil
}

func absolute(u *url.URL) bool {
	if u.Scheme != "http" && u.Scheme != "https" || u.Hostname() == "" {
		return false
	}
	port := u.Port()
	if port != "" {
		number, err := strconv.ParseUint(port, 10, 16)
		return err == nil && number != 0
	}
	return true
}

func (q spec) target() (*url.URL, error) {
	u, err := parseURL(q.URL)
	if err != nil {
		return nil, err
	}
	if q.Session > 0 && q.Base == "" && (q.Bearer != "" || q.Basic || q.KeyName != "") {
		return nil, errors.New("authentication_requires_base")
	}
	if q.Base != "" {
		base, err := parseURL(q.Base)
		if err != nil || !absolute(base) {
			return nil, errors.New("invalid_base_url")
		}
		if !u.IsAbs() && u.Host != "" {
			return nil, errors.New("cross_origin")
		}
		u = base.ResolveReference(u)
		if origin(base) != origin(u) {
			return nil, errors.New("cross_origin")
		}
	}
	if !absolute(u) || q.Secure && u.Scheme != "https" {
		return nil, errors.New("invalid_url")
	}
	values, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, errors.New("invalid_query")
	}
	if len(q.Query) > 0 {
		for _, item := range q.Query {
			values.Add(item.Name, item.Value)
		}
		u.RawQuery = values.Encode()
	}
	return u, nil
}

func safeHeader(name, value string) bool {
	if name == "" || len(name) > 128 || len(value) > 8192 {
		return false
	}
	for _, ch := range name {
		if !(ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || strings.ContainsRune("!#$%&'*+-.^`_|~", ch)) {
			return false
		}
	}
	for _, ch := range []byte(value) {
		if ch != '\t' && (ch < 32 || ch == 127) {
			return false
		}
	}
	switch strings.ToLower(name) {
	case "host", "content-length", "transfer-encoding", "connection", "upgrade", "proxy-authorization", "proxy-connection", "trailer":
		return false
	}
	return true
}

func (q spec) validate() error {
	switch q.Method {
	case "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS":
	default:
		return errors.New("invalid_method")
	}
	if q.Timeout < 0 || q.Timeout > 3600000 || q.Retries < 0 || q.Retries > 3 {
		return errors.New("invalid_limits")
	}
	if q.Redirect != "" && q.Redirect != "none" && q.Redirect != "same-origin" {
		return errors.New("invalid_redirect")
	}
	for _, h := range q.Headers {
		if !safeHeader(h.Name, h.Value) {
			return errors.New("invalid_header")
		}
	}
	if q.Bearer != "" && (!safeHeader("Authorization", q.Bearer) || q.Basic) {
		return errors.New("invalid_authentication")
	}
	if q.KeyName != "" && !safeHeader(q.KeyName, q.KeyValue) {
		return errors.New("invalid_header")
	}
	if q.Basic && strings.Contains(q.User, ":") {
		return errors.New("invalid_authentication")
	}
	if q.Encoding == "json" && !json.Valid([]byte(q.Body)) {
		return errors.New("invalid_json")
	}
	if (len(q.Files) > 0 || len(q.Form) > 0) && (q.Body != "" || q.Encoding != "") {
		return errors.New("mixed_request_body")
	}
	if q.Destination != "" && q.Method != "GET" {
		return errors.New("invalid_download")
	}
	return nil
}

func (s *session) send(parent context.Context, q spec, progress *transfer) (out result) {
	start := time.Now()
	defer func() {
		out.Elapsed = time.Since(start).Milliseconds()
		out.Uploaded = progress.uploaded.Load()
		out.UploadTotal = progress.uploadTotal.Load()
	}()
	if err := q.validate(); err != nil {
		out.Error = err.Error()
		return
	}
	target, err := q.target()
	if err != nil {
		out.Error = err.Error()
		return
	}
	timeout := defaultTimeout
	if q.Timeout > 0 {
		timeout = time.Duration(q.Timeout) * time.Millisecond
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	client := *s.client
	if q.Redirect == "same-origin" {
		client.CheckRedirect = func(next *http.Request, via []*http.Request) error {
			if len(via) >= 10 || origin(next.URL) != origin(via[0].URL) ||
				(via[0].Method != "GET" && via[0].Method != "HEAD") {
				return http.ErrUseLastResponse
			}
			return nil
		}
	} else {
		client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	}

	body, contentType, bodyErr := q.requestBody(ctx, progress)
	if bodyErr != nil {
		out.Error = bodyErr.Error()
		return
	}
	if body != nil {
		defer body.Close()
	}
	retryableMethod := body == nil && (q.Method == "GET" || q.Method == "HEAD" || q.Method == "OPTIONS")
	for i := 0; i <= q.Retries; i++ {
		out.Attempts = i + 1
		req, err := http.NewRequestWithContext(ctx, q.Method, target.String(), body)
		if err != nil {
			out.Error = "invalid_request"
			return
		}
		for _, header := range q.Headers {
			req.Header.Add(header.Name, header.Value)
		}
		if req.Header.Get("Accept-Encoding") == "" && req.Header.Get("Range") == "" && q.Method != "HEAD" {
			req.Header.Set("Accept-Encoding", "gzip, br, zstd")
		}
		if q.Bearer != "" {
			req.Header.Set("Authorization", "Bearer "+q.Bearer)
		}
		if q.Basic {
			req.SetBasicAuth(q.User, q.Password)
		}
		if q.KeyName != "" {
			req.Header.Set(q.KeyName, q.KeyValue)
		}
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}
		if q.Encoding == "json" || (q.Encoding == "text" && contentType == "") {
			req.ContentLength = int64(len(q.Body))
		}
		response, err := client.Do(req)
		if err != nil {
			out.Error = networkError(ctx, err)
			if retryableMethod && out.Error == "connection" && i < q.Retries && delay(ctx, nil, i) {
				continue
			}
			return
		}
		if retryableMethod && i < q.Retries && retryStatus(response.StatusCode) {
			_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, maxInline))
			response.Body.Close()
			if delay(ctx, response.Header, i) {
				continue
			}
			out.Error = networkError(ctx, ctx.Err())
			out.Status = response.StatusCode
			return
		}
		out.Status = response.StatusCode
		out.Protocol = response.Proto
		out.Headers = make(map[string]string, len(response.Header))
		for k, v := range response.Header {
			if len(v) > 0 {
				out.Headers[strings.ToLower(k)] = v[0]
			}
		}
		if q.Method == "HEAD" || response.StatusCode == 204 || response.StatusCode == 304 {
			response.Body.Close()
			return
		}
		bodyReader, closeDecoder, err := decodedBody(response.Header.Get("Content-Encoding"), response.Body)
		if err != nil {
			response.Body.Close()
			out.Error = "content_encoding"
			return
		}
		defer closeDecoder()
		total := response.ContentLength
		if response.Header.Get("Content-Encoding") != "" && !strings.EqualFold(response.Header.Get("Content-Encoding"), "identity") {
			total = -1
		}
		progress.total.Store(total)
		if q.Destination != "" && response.StatusCode >= 200 && response.StatusCode < 300 {
			out.Bytes, out.Error = download(q.Destination, bodyReader, progress)
			if out.Error == "" {
				out.File = q.Destination
			} else if ctx.Err() != nil {
				out.Error = networkError(ctx, ctx.Err())
			}
		} else {
			raw, err := io.ReadAll(io.LimitReader(bodyReader, maxInline+1))
			out.Bytes = int64(len(raw))
			progress.bytes.Store(out.Bytes)
			switch {
			case err != nil:
				out.Error = networkError(ctx, err)
			case len(raw) > maxInline:
				out.Error = "body_too_large"
			case !utf8.Valid(raw) || strings.IndexByte(string(raw), 0) >= 0:
				out.Error = "binary_body"
			default:
				out.Text = string(raw)
			}
		}
		response.Body.Close()
		return
	}
	return
}

func networkError(ctx context.Context, err error) string {
	switch ctx.Err() {
	case context.DeadlineExceeded:
		return "timeout"
	case context.Canceled:
		return "canceled"
	}
	var verification *tls.CertificateVerificationError
	var unknown x509.UnknownAuthorityError
	if errors.As(err, &verification) || errors.As(err, &unknown) {
		return "tls_certificate"
	}
	var timeout net.Error
	if errors.As(err, &timeout) && timeout.Timeout() {
		return "timeout"
	}
	return "connection"
}

func retryStatus(status int) bool {
	return status == 429 || status == 502 || status == 503 || status == 504
}

func delay(ctx context.Context, headers http.Header, attempt int) bool {
	wait := time.Duration(100<<attempt) * time.Millisecond
	retry := headers.Get("Retry-After")
	if seconds, err := strconv.Atoi(retry); err == nil && seconds >= 0 {
		wait = time.Duration(min(seconds, 3600)) * time.Second
	} else if instant, err := http.ParseTime(retry); err == nil {
		wait = max(wait, time.Until(instant))
	}
	t := time.NewTimer(wait)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

func encode(out result) []byte {
	data, _ := json.Marshal(out)
	if len(data) > maxFrame {
		data, _ = json.Marshal(result{Status: out.Status, Protocol: out.Protocol, Bytes: out.Bytes, Elapsed: out.Elapsed, Attempts: out.Attempts, Error: "response_too_large"})
	}
	return data
}
