package sinac

import (
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	maxRetryDelay = 10 * time.Second
	baseDelay     = 500 * time.Millisecond
)

var retryAfterSeconds = regexp.MustCompile(`^\d+(\.\d+)?$`)

func parseRetryAfter(value string, now time.Time) (time.Duration, bool) {
	text := strings.TrimSpace(value)
	if text == "" {
		return 0, false
	}
	if retryAfterSeconds.MatchString(text) {
		seconds, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return 0, false
		}
		if seconds*float64(time.Second) > float64(maxRetryDelay) {
			return maxRetryDelay, true
		}
		return time.Duration(seconds * float64(time.Second)), true
	}
	moment, err := http.ParseTime(text)
	if err != nil {
		return 0, false
	}
	return max(moment.Sub(now), 0), true
}

// retryDelay is the wait before retry number attempt (0 for the first retry).
func retryDelay(attempt int, retryAfter string, random func() float64, now time.Time) time.Duration {
	delay, ok := parseRetryAfter(retryAfter, now)
	if !ok {
		jitter := random() * float64(baseDelay) * math.Pow(2, float64(attempt))
		if jitter > float64(maxRetryDelay) {
			return maxRetryDelay
		}
		delay = time.Duration(jitter)
	}
	return min(delay, maxRetryDelay)
}
