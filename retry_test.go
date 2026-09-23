package sinac

import (
	"testing"
	"time"
)

func TestRetryDelayJitter(t *testing.T) {
	now := time.Date(2026, 9, 23, 12, 0, 0, 0, time.UTC)
	half := func() float64 { return 0.5 }
	if got := retryDelay(0, "", half, now); got != 250*time.Millisecond {
		t.Fatalf("attempt 0 = %v", got)
	}
	if got := retryDelay(2, "", half, now); got != time.Second {
		t.Fatalf("attempt 2 = %v", got)
	}
	one := func() float64 { return 1 }
	if got := retryDelay(10, "", one, now); got != maxRetryDelay {
		t.Fatalf("cap = %v", got)
	}
}

func TestRetryDelayRetryAfter(t *testing.T) {
	now := time.Date(2026, 9, 23, 12, 0, 0, 0, time.UTC)
	zero := func() float64 { return 0 }
	if got := retryDelay(0, "3", zero, now); got != 3*time.Second {
		t.Fatalf("seconds = %v", got)
	}
	if got := retryDelay(0, "1.5", zero, now); got != 1500*time.Millisecond {
		t.Fatalf("fractional = %v", got)
	}
	if got := retryDelay(0, "99999999999999", zero, now); got != maxRetryDelay {
		t.Fatalf("huge = %v", got)
	}
	date := now.Add(4 * time.Second).Format(time.RFC1123)
	date = date[:len(date)-3] + "GMT"
	if got := retryDelay(0, date, zero, now); got != 4*time.Second {
		t.Fatalf("http date = %v", got)
	}
	past := now.Add(-time.Minute).Format(time.RFC1123)
	past = past[:len(past)-3] + "GMT"
	if got := retryDelay(0, past, zero, now); got != 0 {
		t.Fatalf("past date = %v", got)
	}
	if got := retryDelay(0, "soon", func() float64 { return 0.5 }, now); got != 250*time.Millisecond {
		t.Fatalf("garbage falls back to jitter = %v", got)
	}
}
