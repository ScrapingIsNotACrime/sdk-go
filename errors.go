package sinac

import (
	"errors"
	"fmt"
)

// PricingURL is appended to quota-exceeded messages.
const PricingURL = "https://scrapingisnotacrime.com/#pricing"

// Sentinel kinds, matched with errors.Is on any error the client returns.
var (
	ErrBadRequest     = errors.New("sinac: bad request")           // 400, or an invalid path argument
	ErrAuthentication = errors.New("sinac: authentication failed") // 401
	ErrQuotaExceeded  = errors.New("sinac: quota exceeded")        // 402
	ErrNotFound       = errors.New("sinac: not found")             // 404
	ErrRateLimit      = errors.New("sinac: rate limited")          // 429, retried, not charged
	ErrUpstream       = errors.New("sinac: upstream failure")      // 502, retried, not charged
	ErrConnection     = errors.New("sinac: connection error")      // network failure or timeout, retried
	ErrAPI            = errors.New("sinac: API error")             // anything else
)

// Error is returned for every API, network and argument failure.
type Error struct {
	Status    int    // HTTP status; 0 when there was no response
	Message   string // the API's message, or a description of the failure
	RequestID string // the x-request-id header, "" when absent
	kind      error
	cause     error
}

func (e *Error) Error() string {
	text := "sinac: " + e.Message
	if e.Status != 0 {
		text += fmt.Sprintf(" (HTTP %d)", e.Status)
	}
	if e.RequestID != "" {
		text += " [request " + e.RequestID + "]"
	}
	return text
}

// Is reports whether target is this error's kind sentinel.
func (e *Error) Is(target error) bool { return target == e.kind }

// Unwrap returns the underlying cause (network errors), or nil.
func (e *Error) Unwrap() error { return e.cause }

var kindByStatus = map[int]error{
	400: ErrBadRequest,
	401: ErrAuthentication,
	402: ErrQuotaExceeded,
	404: ErrNotFound,
	429: ErrRateLimit,
	502: ErrUpstream,
}

func errorFromStatus(status int, message, requestID string) *Error {
	kind, ok := kindByStatus[status]
	if !ok {
		kind = ErrAPI
	}
	if kind == ErrQuotaExceeded {
		message += " — see " + PricingURL
	}
	return &Error{Status: status, Message: message, RequestID: requestID, kind: kind}
}

// 429 and 502 don't consume credits, so retrying them costs the customer nothing.
func isRetryable(err *Error) bool {
	return err.kind == ErrRateLimit || err.kind == ErrUpstream || err.kind == ErrConnection
}
