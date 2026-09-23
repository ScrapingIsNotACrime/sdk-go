package sinac

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// roundTrip lets a test answer requests without a network.
type roundTrip func(*http.Request) (*http.Response, error)

func (f roundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func response(status int, body string, header http.Header) *http.Response {
	if header == nil {
		header = http.Header{}
	}
	return &http.Response{StatusCode: status, Header: header, Body: io.NopCloser(strings.NewReader(body))}
}

func testCore(t *testing.T, rt http.RoundTripper, opts ...Option) *httpCore {
	t.Helper()
	all := append([]Option{WithAPIKey("sinac_test"), WithHTTPClient(&http.Client{Transport: rt})}, opts...)
	cfg, err := resolveConfig(all, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	core := newHTTPCore(cfg)
	core.sleep = func(ctx context.Context, _ time.Duration) error { return ctx.Err() }
	core.random = func() float64 { return 0 }
	return core
}

func TestSegmentEncodesAndRejects(t *testing.T) {
	cases := map[string]string{
		"nasa": "nasa", "a/b": "a%2Fb", "#tag": "%23tag", "é": "%C3%A9", "highlight:1": "highlight%3A1",
		"a b": "a%20b", "-_.!~*'()": "-_.!~*'()", "a?b&c=d": "a%3Fb%26c%3Dd",
	}
	for in, want := range cases {
		got, err := segment(in)
		if err != nil || got != want {
			t.Errorf("segment(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, bad := range []string{"", ".", ".."} {
		_, err := segment(bad)
		var apiErr *Error
		if !errors.As(err, &apiErr) || !errors.Is(err, ErrBadRequest) || apiErr.Status != 0 {
			t.Errorf("segment(%q) error = %v", bad, err)
		}
	}
}

func TestBuildURLKeepsOrderAndSkipsEmpty(t *testing.T) {
	r := route{path: "/github/repositories", query: append(append(qs("q", "stars:>1 go"), qi("limit", 0)...), qi("page", 2)...)}
	got := buildURL("https://api.example.com/v1", r)
	if got != "https://api.example.com/v1/github/repositories?q=stars%3A%3E1+go&page=2" {
		t.Fatalf("buildURL = %s", got)
	}
	if got := buildURL("https://x/v1", route{path: "/a"}); got != "https://x/v1/a" {
		t.Fatalf("no query = %s", got)
	}
}

func TestGetSendsHeadersAndReturnsData(t *testing.T) {
	var seen *http.Request
	core := testCore(t, roundTrip(func(r *http.Request) (*http.Response, error) {
		seen = r
		return response(200, `{"message":"ok","data":{"username":"nasa"}}`, nil), nil
	}))
	data, err := core.get(context.Background(), route{path: "/instagram/profile/nasa"})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"username":"nasa"}` {
		t.Fatalf("data = %s", data)
	}
	if seen.URL.String() != DefaultBaseURL+"/instagram/profile/nasa" || seen.Header.Get("X-Api-Key") != "sinac_test" ||
		seen.Header.Get("Accept") != "application/json" ||
		!strings.HasPrefix(seen.Header.Get("User-Agent"), "scrapingisnotacrime-go/") {
		t.Fatalf("request = %s %v", seen.URL, seen.Header)
	}
}

func TestErrorMapping(t *testing.T) {
	cases := []struct {
		status int
		body   string
		kind   error
		msg    string
	}{
		{404, `{"message":"Profile not found","data":null}`, ErrNotFound, "Profile not found"},
		{401, `{"message":"Invalid API key"}`, ErrAuthentication, "Invalid API key"},
		{402, `{"message":"No credits"}`, ErrQuotaExceeded, "No credits — see " + PricingURL},
		{500, `<html>oops</html>`, ErrAPI, "HTTP 500: <html>oops</html>"},
		{404, `not json`, ErrNotFound, "HTTP 404: not json"},
		{400, `{"message":42}`, ErrBadRequest, `HTTP 400: {"message":42}`},
	}
	for _, c := range cases {
		core := testCore(t, roundTrip(func(*http.Request) (*http.Response, error) {
			return response(c.status, c.body, http.Header{"X-Request-Id": {"req_7"}}), nil
		}), WithMaxRetries(0))
		_, err := core.get(context.Background(), route{path: "/x"})
		var apiErr *Error
		if !errors.As(err, &apiErr) || !errors.Is(err, c.kind) || apiErr.Status != c.status ||
			apiErr.Message != c.msg || apiErr.RequestID != "req_7" {
			t.Errorf("%d %s: got %#v", c.status, c.body, err)
		}
	}
}

func TestTwoxxWithoutEnvelopeIsAPIError(t *testing.T) {
	for _, body := range []string{`<html>proxy</html>`, `{"message":"ok"}`, `[1,2]`, ``} {
		core := testCore(t, roundTrip(func(*http.Request) (*http.Response, error) {
			return response(200, body, nil), nil
		}))
		_, err := core.get(context.Background(), route{path: "/x"})
		var apiErr *Error
		if !errors.As(err, &apiErr) || !errors.Is(err, ErrAPI) || apiErr.Status != 200 ||
			!strings.HasPrefix(apiErr.Message, "unexpected response body (HTTP 200)") {
			t.Errorf("body %q: got %v", body, err)
		}
	}
}

func TestRedirectIsNotFollowed(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		http.Redirect(w, r, "https://elsewhere.example/v1"+r.URL.Path, http.StatusMovedPermanently)
	}))
	defer server.Close()
	userClient := server.Client()
	core := testCore(t, nil, WithHTTPClient(userClient), WithBaseURL(server.URL+"/v1"))
	_, err := core.get(context.Background(), route{path: "/x"})
	var apiErr *Error
	if !errors.As(err, &apiErr) || !errors.Is(err, ErrAPI) || apiErr.Status != 301 ||
		apiErr.Message != "HTTP 301: redirect to https://elsewhere.example/v1/v1/x not followed" {
		t.Fatalf("got %v", err)
	}
	if calls.Load() != 1 || userClient.CheckRedirect != nil {
		t.Fatalf("calls=%d, user client modified=%v", calls.Load(), userClient.CheckRedirect != nil)
	}
}

func TestRetriesRetryableThenSucceeds(t *testing.T) {
	var calls atomic.Int32
	var waits []time.Duration
	core := testCore(t, roundTrip(func(*http.Request) (*http.Response, error) {
		switch calls.Add(1) {
		case 1:
			return response(429, `{"message":"slow down"}`, http.Header{"Retry-After": {"2"}}), nil
		case 2:
			return response(502, `{"message":"upstream"}`, nil), nil
		default:
			return response(200, `{"data":{"ok":true}}`, nil), nil
		}
	}))
	core.sleep = func(_ context.Context, d time.Duration) error { waits = append(waits, d); return nil }
	core.random = func() float64 { return 1 }
	data, err := core.get(context.Background(), route{path: "/x"})
	if err != nil || string(data) != `{"ok":true}` {
		t.Fatalf("data=%s err=%v", data, err)
	}
	if calls.Load() != 3 || len(waits) != 2 || waits[0] != 2*time.Second || waits[1] != time.Second {
		t.Fatalf("calls=%d waits=%v", calls.Load(), waits)
	}
}

func TestNoRetryForClientErrorsAndMaxRetries(t *testing.T) {
	var calls atomic.Int32
	core := testCore(t, roundTrip(func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return response(404, `{"message":"nope"}`, nil), nil
	}))
	if _, err := core.get(context.Background(), route{path: "/x"}); !errors.Is(err, ErrNotFound) || calls.Load() != 1 {
		t.Fatalf("404: err=%v calls=%d", err, calls.Load())
	}
	calls.Store(0)
	core = testCore(t, roundTrip(func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return response(502, `{"message":"down"}`, nil), nil
	}), WithMaxRetries(1))
	if _, err := core.get(context.Background(), route{path: "/x"}); !errors.Is(err, ErrUpstream) || calls.Load() != 2 {
		t.Fatalf("502: err=%v calls=%d", err, calls.Load())
	}
}

func TestNetworkErrorIsRetriedAndWrapped(t *testing.T) {
	var calls atomic.Int32
	core := testCore(t, roundTrip(func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return nil, io.ErrUnexpectedEOF
	}), WithMaxRetries(1))
	_, err := core.get(context.Background(), route{path: "/x"})
	var apiErr *Error
	if !errors.As(err, &apiErr) || !errors.Is(err, ErrConnection) || !errors.Is(err, io.ErrUnexpectedEOF) ||
		apiErr.Status != 0 || calls.Load() != 2 {
		t.Fatalf("err=%v calls=%d", err, calls.Load())
	}
}

func TestPerAttemptTimeoutIsRetriedAsConnectionError(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(200)
		w.(http.Flusher).Flush()
		select { // headers sent, body stalls: the per-attempt timeout must cover the body read
		case <-r.Context().Done():
		case <-time.After(2 * time.Second):
		}
	}))
	defer server.Close()
	core := testCore(t, nil, WithHTTPClient(server.Client()), WithBaseURL(server.URL), WithTimeout(100*time.Millisecond),
		WithMaxRetries(1))
	start := time.Now()
	_, err := core.get(context.Background(), route{path: "/x"})
	if !errors.Is(err, ErrConnection) || calls.Load() != 2 || time.Since(start) > 1500*time.Millisecond {
		t.Fatalf("err=%v calls=%d elapsed=%v", err, calls.Load(), time.Since(start))
	}
	if !strings.Contains(err.Error(), "timed out after 100ms") {
		t.Fatalf("message = %v", err)
	}
}

func TestCallerCancelDuringBackoff(t *testing.T) {
	var calls atomic.Int32
	ctx, cancel := context.WithCancel(context.Background())
	core := testCore(t, roundTrip(func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return response(429, `{"message":"slow"}`, nil), nil
	}))
	core.sleep = sleepContext
	core.random = func() float64 { return 1 }
	go func() { time.Sleep(50 * time.Millisecond); cancel() }()
	start := time.Now()
	_, err := core.get(ctx, route{path: "/x"})
	if !errors.Is(err, context.Canceled) || calls.Load() != 1 || time.Since(start) > 400*time.Millisecond {
		t.Fatalf("err=%v calls=%d elapsed=%v", err, calls.Load(), time.Since(start))
	}
}

func TestCallerDeadlineIsNotRetried(t *testing.T) {
	var calls atomic.Int32
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	core := testCore(t, roundTrip(func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		return nil, r.Context().Err()
	}))
	_, err := core.get(ctx, route{path: "/x"})
	if !errors.Is(err, context.Canceled) || calls.Load() > 1 {
		t.Fatalf("err=%v calls=%d", err, calls.Load())
	}
}

func TestTypeMismatchIsTolerated(t *testing.T) {
	type profile struct {
		Username  string `json:"username"`
		Followers int64  `json:"followers"`
	}
	core := testCore(t, roundTrip(func(*http.Request) (*http.Response, error) {
		return response(200, `{"data":{"username":"nasa","followers":"many"}}`, nil), nil
	}))
	got, err := getJSON[profile](context.Background(), core, route{path: "/x"}, nil)
	if err != nil || got.Username != "nasa" || got.Followers != 0 {
		t.Fatalf("got=%+v err=%v", got, err)
	}
}

func TestGetJSONReturnsRouteErrorWithoutRequest(t *testing.T) {
	core := testCore(t, roundTrip(func(*http.Request) (*http.Response, error) {
		t.Fatal("no request expected")
		return nil, nil
	}))
	_, routeErr := segment("..")
	if _, err := getJSON[json.RawMessage](context.Background(), core, route{}, routeErr); !errors.Is(err, ErrBadRequest) {
		t.Fatalf("err = %v", err)
	}
}
