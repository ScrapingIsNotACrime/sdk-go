// Package sinac is the Go client for the ScrapingIsNotACrime API
// (https://scrapingisnotacrime.com): public data from Instagram, TikTok,
// YouTube, the App Store, GitHub, Hacker News, Bluesky, Twitch and Linktree.
//
//	import sinac "github.com/ScrapingIsNotACrime/sdk-go"
//
//	client, err := sinac.NewClient(sinac.WithAPIKey("sinac_..."))
//	profile, err := client.Instagram.Profile(ctx, "nasa")
//
// Errors are *Error values matched with errors.Is against the Err* sentinels.
// Paginated methods return a *Page; Page.All iterates every item lazily, and
// each page fetched is one billed request.
package sinac
