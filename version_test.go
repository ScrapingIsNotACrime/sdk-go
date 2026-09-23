package sinac

import (
	"strings"
	"testing"
)

func TestUserAgent(t *testing.T) {
	ua := userAgent()
	if !strings.HasPrefix(ua, "scrapingisnotacrime-go/") {
		t.Fatalf("userAgent() = %q", ua)
	}
	if strings.Contains(ua, "/v") {
		t.Fatalf("version must not keep the v prefix: %q", ua)
	}
	if Version == "" {
		t.Fatal("Version is empty")
	}
}
