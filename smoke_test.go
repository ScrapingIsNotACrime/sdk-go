//go:build smoke

package sinac

import (
	"context"
	"os"
	"testing"
	"time"
)

// Run with: SCRAPINGISNOTACRIME_API_KEY=sinac_... ./scripts/dev.sh go test -tags smoke -run Smoke -v ./...
// Each call is one billed request.
func TestSmokeOneCallPerPlatform(t *testing.T) {
	if os.Getenv(APIKeyEnv) == "" {
		t.Skip("set " + APIKeyEnv + " to run the smoke test")
	}
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	calls := map[string]func() error{
		"instagram": func() error { _, err := client.Instagram.Profile(ctx, "instagram"); return err },
		"tiktok":    func() error { _, err := client.TikTok.Profile(ctx, "tiktok"); return err },
		"youtube":   func() error { _, err := client.YouTube.Videos(ctx, "youtube"); return err },
		"appstore": func() error {
			_, err := client.AppStore.Search(ctx, "instagram", &AppStoreSearchParams{Limit: 1})
			return err
		},
		"github":     func() error { _, err := client.GitHub.Profile(ctx, "torvalds"); return err },
		"hackernews": func() error { _, err := client.HackerNews.Item(ctx, 8863); return err },
		"bluesky":    func() error { _, err := client.Bluesky.Profile(ctx, "bsky.app"); return err },
		"twitch":     func() error { _, err := client.Twitch.Profile(ctx, "ninja"); return err },
		"linktree":   func() error { _, err := client.Linktree.Profile(ctx, "linktree"); return err },
	}
	for name, call := range calls {
		if err := call(); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}
