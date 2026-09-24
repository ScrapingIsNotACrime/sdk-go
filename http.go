package sinac

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type queryParam struct{ key, value string }

type route struct {
	path  string
	query []queryParam
}

// with returns a copy of r with the query parameters appended.
func (r route) with(query ...[]queryParam) route {
	r.query = append(append([]queryParam(nil), r.query...), concat(query)...)
	return r
}

func concat(query [][]queryParam) []queryParam {
	var all []queryParam
	for _, q := range query {
		all = append(all, q...)
	}
	return all
}

// qs is a string query parameter, omitted when empty.
func qs(key, value string) []queryParam {
	if value == "" {
		return nil
	}
	return []queryParam{{key, value}}
}

// qi is an integer query parameter, omitted when zero (the API's default applies).
func qi(key string, value int) []queryParam {
	if value == 0 {
		return nil
	}
	return []queryParam{{key, strconv.Itoa(value)}}
}

// segment encodes one path segment like JavaScript's encodeURIComponent and
// rejects values that would drop or climb a path level.
func segment(value string) (string, error) {
	if value == "" || value == "." || value == ".." {
		return "", &Error{Message: fmt.Sprintf("invalid path segment %q", value), kind: ErrBadRequest}
	}
	const hex = "0123456789ABCDEF"
	var b strings.Builder
	for i := 0; i < len(value); i++ {
		c := value[i]
		if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || strings.IndexByte("-_.!~*'()", c) >= 0 {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('%')
		b.WriteByte(hex[c>>4])
		b.WriteByte(hex[c&15])
	}
	return b.String(), nil
}

// pathOf fills each %s in format with an encoded segment.
func pathOf(format string, values ...string) (string, error) {
	encoded := make([]any, len(values))
	for i, value := range values {
		s, err := segment(value)
		if err != nil {
			return "", err
		}
		encoded[i] = s
	}
	return fmt.Sprintf(format, encoded...), nil
}

func buildURL(baseURL string, r route) string {
	if len(r.query) == 0 {
		return baseURL + r.path
	}
	parts := make([]string, len(r.query))
	for i, p := range r.query {
		parts[i] = url.QueryEscape(p.key) + "=" + url.QueryEscape(p.value)
	}
	return baseURL + r.path + "?" + strings.Join(parts, "&")
}

type httpCore struct {
	cfg    config
	client *http.Client
	sleep  func(context.Context, time.Duration) error
	random func() float64
	now    func() time.Time
}

func newHTTPCore(cfg config) *httpCore {
	base := cfg.httpClient
	if base == nil {
		base = &http.Client{}
	}
	// A shallow copy: the caller's client is never modified. Redirects are not
	// followed because the X-Api-Key header would be forwarded to the new host.
	client := *base
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &httpCore{cfg: cfg, client: &client, sleep: sleepContext, random: rand.Float64, now: time.Now}
}

func sleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

type failure struct {
	err        *Error
	retryAfter string
}

func (h *httpCore) get(ctx context.Context, r route) (json.RawMessage, error) {
	target := buildURL(h.cfg.baseURL, r)
	for attempt := 0; ; attempt++ {
		data, fail := h.attempt(ctx, target)
		if fail == nil {
			return data, nil
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, fmt.Errorf("sinac: %w", ctxErr)
		}
		if attempt >= h.cfg.maxRetries || !isRetryable(fail.err) {
			return nil, fail.err
		}
		if err := h.sleep(ctx, retryDelay(attempt, fail.retryAfter, h.random, h.now())); err != nil {
			return nil, fmt.Errorf("sinac: %w", err)
		}
	}
}

// attempt makes one request under its own timeout, which covers connect,
// headers and the whole body read.
func (h *httpCore) attempt(ctx context.Context, target string) (json.RawMessage, *failure) {
	ctx, cancel := context.WithTimeout(ctx, h.cfg.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, &failure{err: &Error{Message: "invalid request: " + err.Error(), kind: ErrAPI, cause: err}}
	}
	req.Header.Set("X-Api-Key", h.cfg.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", userAgent())
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, &failure{err: h.connectionError(ctx, err)}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &failure{err: h.connectionError(ctx, err)}
	}
	return interpret(resp.StatusCode, resp.Header, body)
}

// connectionError describes a transport failure. attemptCtx is the per-attempt
// context: when its deadline passed, the failure is our timeout, whatever error
// the transport chose to report.
func (h *httpCore) connectionError(attemptCtx context.Context, err error) *Error {
	message := "network error: " + err.Error()
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(attemptCtx.Err(), context.DeadlineExceeded) {
		message = fmt.Sprintf("request timed out after %s", h.cfg.timeout)
	}
	return &Error{Message: message, kind: ErrConnection, cause: err}
}

var whitespace = regexp.MustCompile(`\s+`)

func snippet(body []byte) string {
	flat := []rune(strings.TrimSpace(whitespace.ReplaceAllString(string(body), " ")))
	if len(flat) > 200 {
		return string(flat[:200]) + "…"
	}
	return string(flat)
}

// bodySnippet is snippet(body), with a readable fallback for an empty body.
// Error.Error() already appends "(HTTP n)", so callers never repeat the
// status themselves in the message they build from this.
func bodySnippet(body []byte) string {
	if s := snippet(body); s != "" {
		return s
	}
	return "empty response body"
}

func interpret(status int, header http.Header, body []byte) (json.RawMessage, *failure) {
	requestID := header.Get("X-Request-Id")
	var envelope map[string]json.RawMessage
	isEnvelope := json.Unmarshal(body, &envelope) == nil && envelope != nil

	if status >= 200 && status < 300 {
		if data, ok := envelope["data"]; isEnvelope && ok {
			return data, nil
		}
		message := "unexpected response body: " + bodySnippet(body)
		return nil, &failure{err: &Error{Status: status, Message: message, RequestID: requestID, kind: ErrAPI}}
	}

	retryAfter := header.Get("Retry-After")
	if status >= 300 && status < 400 {
		message := "redirect not followed"
		if location := header.Get("Location"); location != "" {
			message = fmt.Sprintf("redirect to %s not followed", location)
		}
		return nil, &failure{err: &Error{Status: status, Message: message, RequestID: requestID, kind: ErrAPI}, retryAfter: retryAfter}
	}

	var message string
	if raw, ok := envelope["message"]; !isEnvelope || !ok || json.Unmarshal(raw, &message) != nil {
		message = bodySnippet(body)
	}
	return nil, &failure{err: errorFromStatus(status, message, requestID), retryAfter: retryAfter}
}

// decode fills out from raw. A field whose JSON type does not match its
// destination is left at its zero value instead of failing the whole call.
// A mismatch at the root (out itself has the wrong shape, e.g. an array or a
// string where a struct was expected) still fails: json.UnmarshalTypeError
// only names a field for a mismatch found inside out, leaving Field empty at
// the root.
func decode(raw json.RawMessage, out any) error {
	err := json.Unmarshal(raw, out)
	var typeErr *json.UnmarshalTypeError
	if err == nil || (errors.As(err, &typeErr) && typeErr.Field != "") {
		return nil
	}
	return &Error{Message: "unexpected response data: " + err.Error(), kind: ErrAPI, cause: err}
}

func getJSON[T any](ctx context.Context, h *httpCore, r route, routeErr error) (*T, error) {
	if routeErr != nil {
		return nil, routeErr
	}
	raw, err := h.get(ctx, r)
	if err != nil {
		return nil, err
	}
	var out T
	if err := decode(raw, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
