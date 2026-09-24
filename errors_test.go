package sinac

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestErrorFromStatusKinds(t *testing.T) {
	cases := map[int]error{
		400: ErrBadRequest, 401: ErrAuthentication, 402: ErrQuotaExceeded, 404: ErrNotFound,
		429: ErrRateLimit, 502: ErrUpstream, 500: ErrAPI, 503: ErrAPI, 418: ErrAPI,
	}
	for status, kind := range cases {
		err := errorFromStatus(status, "boom", "req_1")
		if !errors.Is(err, kind) {
			t.Errorf("%d: errors.Is(%v) = false", status, kind)
		}
		if err.Status != status || err.RequestID != "req_1" {
			t.Errorf("%d: fields %+v", status, err)
		}
	}
}

func TestQuotaMessageHasPricingURL(t *testing.T) {
	err := errorFromStatus(402, "No credits left", "")
	if err.Message != "No credits left — see "+PricingURL {
		t.Fatalf("message = %q", err.Message)
	}
}

func TestErrorStringAndUnwrap(t *testing.T) {
	err := &Error{Status: 404, Message: "Profile not found", RequestID: "req_9", kind: ErrNotFound}
	if got := err.Error(); !strings.Contains(got, "Profile not found") || !strings.Contains(got, "404") ||
		!strings.Contains(got, "req_9") {
		t.Fatalf("Error() = %q", got)
	}
	wrapped := &Error{Message: "network error", kind: ErrConnection, cause: io.ErrUnexpectedEOF}
	if !errors.Is(wrapped, ErrConnection) || !errors.Is(wrapped, io.ErrUnexpectedEOF) {
		t.Fatal("connection error must match its kind and its cause")
	}
	if errors.Is(wrapped, ErrNotFound) {
		t.Fatal("must not match another kind")
	}
}

func TestIsRetryable(t *testing.T) {
	for _, kind := range []error{ErrRateLimit, ErrUpstream, ErrConnection} {
		if !isRetryable(&Error{kind: kind}) {
			t.Errorf("%v should be retryable", kind)
		}
	}
	for _, kind := range []error{ErrBadRequest, ErrAuthentication, ErrQuotaExceeded, ErrNotFound, ErrAPI} {
		if isRetryable(&Error{kind: kind}) {
			t.Errorf("%v should not be retryable", kind)
		}
	}
}
