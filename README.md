# ScrapingIsNotACrime Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/ScrapingIsNotACrime/sdk-go.svg)](https://pkg.go.dev/github.com/ScrapingIsNotACrime/sdk-go)
[![license](https://img.shields.io/github/license/ScrapingIsNotACrime/sdk-go.svg)](./LICENSE)

Official Go SDK for the [ScrapingIsNotACrime](https://scrapingisnotacrime.com) public data API — typed access to Instagram, TikTok, YouTube, App Store, GitHub, Hacker News, Bluesky, Twitch and Linktree.

## Install

```bash
go get github.com/ScrapingIsNotACrime/sdk-go
```

Requires Go 1.23+. The module has no third-party dependencies: only the standard library.

## Quick start

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	sinac "github.com/ScrapingIsNotACrime/sdk-go"
)

func main() {
	client, err := sinac.NewClient(sinac.WithAPIKey("sinac_..."))
	if err != nil {
		log.Fatal(err)
	}

	profile, err := client.Instagram.Profile(context.Background(), "nasa")
	if errors.Is(err, sinac.ErrNotFound) {
		fmt.Println("no such profile")
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(profile.Username, profile.Followers)
}
```

`NewClient` also reads the API key from the `SCRAPINGISNOTACRIME_API_KEY` environment variable, so `sinac.NewClient()` with no options works when that variable is set. Get a key at [scrapingisnotacrime.com/dashboard/api-keys](https://scrapingisnotacrime.com/dashboard/api-keys). Keys start with `sinac_`.

More runnable examples are in [`example_test.go`](./example_test.go) (`Example` and `ExamplePage_All`), and on [pkg.go.dev](https://pkg.go.dev/github.com/ScrapingIsNotACrime/sdk-go).

## Configuration

```go
client, err := sinac.NewClient(
	sinac.WithAPIKey("sinac_..."),
	sinac.WithBaseURL("https://api.scrapingisnotacrime.com/v1"),
	sinac.WithTimeout(30*time.Second),
	sinac.WithMaxRetries(2),
	sinac.WithHTTPClient(nil),
)
```

| Option | Default | Description |
|---|---|---|
| `WithAPIKey` | `SCRAPINGISNOTACRIME_API_KEY` env var | Your API key (`sinac_…`). `NewClient` returns an error if none is set. |
| `WithBaseURL` | `https://api.scrapingisnotacrime.com/v1` | API base URL. Must be the final HTTPS URL: redirects are not followed, so a base URL that itself redirects fails (see [Errors](#errors)) — following one would forward the `X-Api-Key` header to whatever host it points to. |
| `WithTimeout` | `30 * time.Second` | Per attempt, covering connect, headers and the whole body read. |
| `WithMaxRetries` | `2` | Extra attempts for 429, 502 and network errors. `0` disables retries. |
| `WithHTTPClient` | a new `*http.Client` | Bring your own `*http.Client` (for tests, proxies or connection pooling). It is never modified: the SDK works with a shallow copy so it can set its own `CheckRedirect`. |

## Methods

`Client` groups its methods under a service field per platform: `client.Instagram`, `client.TikTok`, `client.YouTube`, `client.AppStore`, `client.GitHub`, `client.HackerNews`, `client.Bluesky`, `client.Twitch`, `client.Linktree`. Every method takes a `context.Context` as its first argument (omitted below for brevity) and returns the response envelope's `data`, decoded into a typed struct.

Methods marked `Page` return `*Page[T, R]` (see [Pagination](#pagination)) instead of a plain struct.

| Namespace | Method | Route | Page |
|---|---|---|---|
| Instagram | `Profile(username)` | `/instagram/profile/{username}` | |
| Instagram | `Contact(username)` | `/instagram/profile/{username}/contact` | |
| Instagram | `LatestPosts(username)` | `/instagram/profile/{username}/timeline/latest` | |
| Instagram | `Posts(username, p *InstagramPostsParams)` — `Count` 1-50 (default 12), `Cursor` | `/instagram/profile/{username}/timeline` | Page |
| Instagram | `Highlights(username)` | `/instagram/profile/{username}/highlights` | |
| Instagram | `Highlight(highlightID)` | `/instagram/highlights/{highlightId}` | |
| Instagram | `MediaByID(username, mediaID)` | `/instagram/profile/{username}/media/{mediaId}` | |
| Instagram | `Media(shortcode)` | `/instagram/media/{shortcode}` | |
| Instagram | `Download(shortcode)` | `/instagram/media/{shortcode}/download` | |
| Instagram | `ShortcodeToID(shortcode)` | `/instagram/media/{shortcode}/id` | |
| Instagram | `IDToShortcode(mediaID)` | `/instagram/media/id/{mediaId}` | |
| Instagram | `Reel(shortcode)` | `/instagram/reels/{shortcode}` | |
| TikTok | `Profile(username)` | `/tiktok/profile/{username}` | |
| TikTok | `Video(videoID)` | `/tiktok/video/{videoId}` | |
| YouTube | `Videos(handle)` | `/youtube/channel/{handle}/videos` | |
| App Store | `Search(term, p *AppStoreSearchParams)` — `Country` (default `"us"`), `Limit` 1-200 (default 10) | `/appstore/search` | |
| App Store | `Reviews(appID, p *AppStoreReviewsParams)` — `Country` (default `"us"`), `Page` 1-10 (default 1) | `/appstore/reviews` | Page |
| GitHub | `Profile(handle)` | `/github/profiles/{handle}` | |
| GitHub | `Followers(handle, p *GitHubListParams)` — `Limit` 1-100 (default 30), `Page` 1-based (default 1) | `/github/profiles/{handle}/followers` | Page |
| GitHub | `Following(handle, p *GitHubListParams)` — same params as `Followers` | `/github/profiles/{handle}/following` | Page |
| GitHub | `Repositories(handle, p *GitHubListParams)` — same params as `Followers` | `/github/profiles/{handle}/repositories` | Page |
| GitHub | `SearchRepositories(q, p *GitHubListParams)` — `q` in GitHub search syntax; same params as `Followers` | `/github/repositories` | Page |
| GitHub | `Trending(p *GitHubTrendingParams)` — `Since` (default `daily`), `Language`, `Limit` 1-100 (default 30) | `/github/trending/repositories` | |
| Hacker News | `Feed(feed, p *HackerNewsListParams)` — `Limit` 1-50 (default 20), `Page` 0-based (default 0) | `/hackernews/feeds/{feed}` | Page |
| Hacker News | `Item(id)` | `/hackernews/items/{id}` | |
| Hacker News | `Search(q, p *HackerNewsListParams)` — same params as `Feed` | `/hackernews/search` | Page |
| Hacker News | `User(username)` | `/hackernews/users/{username}` | |
| Hacker News | `Submissions(username, p *HackerNewsListParams)` — same params as `Feed` | `/hackernews/users/{username}/submissions` | Page |
| Hacker News | `Comments(username, p *HackerNewsListParams)` — same params as `Feed` | `/hackernews/users/{username}/comments` | Page |
| Bluesky | `Profile(handle)` | `/bluesky/profiles/{handle}` | |
| Bluesky | `Posts(handle, p *BlueskyPostsParams)` — `Limit` 1-100 (default 25), `Cursor` | `/bluesky/profiles/{handle}/posts` | Page |
| Twitch | `Profile(handle)` | `/twitch/profiles/{handle}` | |
| Twitch | `Videos(handle, p *TwitchVideosParams)` — `Limit` 1-100 (default 20) | `/twitch/profiles/{handle}/videos` | |
| Linktree | `Profile(handle)` | `/linktree/profiles/{handle}` | |

`GitHubTrendingSince` (`Trending`'s `Since`) and `HackerNewsFeed` (`Feed`'s `feed`) are typed string constants — `sinac.TrendingDaily`/`TrendingWeekly`/`TrendingMonthly` and `sinac.FeedTop`/`FeedNew`/`FeedBest`/`FeedAsk`/`FeedShow`/`FeedJob`. Pass `nil` for any `p *...Params` to use every default. Path arguments (`username`, `handle`, `shortcode`, …) are validated before any request: an empty string, `"."` or `".."` returns `*sinac.Error` with `Status: 0` and `errors.Is(err, sinac.ErrBadRequest)`, and the value is percent-encoded the way JavaScript's `encodeURIComponent` would encode it.

## Pagination

Methods marked `Page` return `*sinac.Page[T, R]`, where `T` is the item type and `R` is the full page response:

```go
type Page[T, R any] struct {
	Items      []T
	HasMore    bool
	NextCursor string // cursor endpoints; "" when there is none
	NextPage   int    // page-number endpoints; 0 when there is none
	Data       R      // the untouched response of this page

	// Next fetches the following page. It returns (nil, nil) when there is none.
	// All yields every item from this page onward, fetching later pages lazily.
}
```

Iterate every item with `All`, using Go 1.23's range-over-func. **Iterating fetches every remaining page; each page is one billed request**, so bound the loop:

```go
page, err := client.GitHub.Followers(ctx, "torvalds", &sinac.GitHubListParams{Limit: 100})
if err != nil {
	log.Fatal(err)
}

n := 0
for user, err := range page.All(ctx) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.Username)
	if n++; n >= 250 {
		break // stop early; no further pages are fetched
	}
}
```

Or walk pages one at a time with `Next`:

```go
page, err := client.Bluesky.Posts(ctx, "bsky.app", &sinac.BlueskyPostsParams{Limit: 25})
for page != nil && err == nil {
	for _, post := range page.Items {
		fmt.Println(post.Text)
	}
	page, err = page.Next(ctx)
}
if err != nil {
	log.Fatal(err)
}
```

`page.Data` gives you the untouched response of the current page, so fields such as `Total` stay reachable. Breaking out of the `All` loop, or simply not calling `Next` again, stops fetching — no further pages are requested.

## Errors

Every failure is returned as `*sinac.Error`, carrying `Status` (the HTTP status, or `0` for network errors and invalid arguments), `Message` and `RequestID` (the `X-Request-Id` header, when the API sends one). Match it by kind with `errors.Is` against the sentinels below, or unwrap it with `errors.As` to read `Status` and `RequestID`:

```go
var apiErr *sinac.Error
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.Status, apiErr.RequestID)
}
```

| Sentinel | Status | Retried |
|---|---|---|
| `ErrBadRequest` | 400, or an invalid path argument | no |
| `ErrAuthentication` | 401 | no |
| `ErrQuotaExceeded` | 402 | no |
| `ErrNotFound` | 404 | no |
| `ErrRateLimit` | 429 | yes |
| `ErrUpstream` | 502 | yes |
| `ErrConnection` | network failure or per-attempt timeout | yes |
| `ErrAPI` | any other status, a 2xx without the JSON envelope, or a redirect (not followed) | no |

```go
profile, err := client.TikTok.Profile(ctx, "this-user-does-not-exist-123")
switch {
case errors.Is(err, sinac.ErrNotFound):
	fmt.Println("no such profile")
case errors.Is(err, sinac.ErrQuotaExceeded):
	log.Fatal(err) // includes the pricing link in Message
case errors.Is(err, sinac.ErrRateLimit):
	fmt.Println("TikTok is rate limiting; already retried, try again later")
case err != nil:
	log.Fatal(err)
default:
	fmt.Println(profile.Username)
}
```

## Context and cancellation

Every method takes the caller's `context.Context` and honors it: a `ctx` that is already cancelled, or that reaches its deadline — including while the client is sleeping between retries — makes the call return immediately with an error satisfying `errors.Is(err, ctx.Err())`. Cancellation is never retried, even when the underlying failure would otherwise be.

## Retries

`ErrRateLimit` (429), `ErrUpstream` (502) and `ErrConnection` (network failure or per-attempt timeout) are retried automatically, up to `WithMaxRetries` additional attempts (default 2). These failures don't consume credits, so retrying them costs you nothing.

The delay before each retry uses the `Retry-After` header when the API sends one (seconds or an HTTP date); otherwise it's exponential backoff with full jitter, starting around 500 ms and doubling per attempt. Every wait is capped at 10 seconds. Set `WithMaxRetries(0)` to disable retries entirely.

## Releases and changelog

Every merge to `main` is released automatically: the version comes from the commit messages since the last release, following [Conventional Commits](https://www.conventionalcommits.org/).

| Commits since the last release | New version (while in 0.x) |
|---|---|
| only `docs:`, `chore:`, `test:`, `ci:`, `build:`, `refactor:` | none |
| at least one `fix:` | patch (`0.1.0` → `0.1.1`) |
| at least one `feat:` | minor (`0.1.1` → `0.2.0`) |
| `feat!:` or a `BREAKING CHANGE:` footer | minor while in 0.x |

The pipeline tags `vX.Y.Z` and publishes the GitHub Release with the notes; in Go the tag is the release — pkg.go.dev picks it up, no separate publish step. The changelog is the [Releases page](https://github.com/ScrapingIsNotACrime/sdk-go/releases); the version is resolved from git tags at build time (`runtime/debug.BuildInfo`), so no version number is written into the repository.

Merge PRs with a merge commit or rebase so each Conventional Commit is analysed; if you squash, the PR title must be a Conventional Commit (e.g. `feat: ...`).

## Links

- Docs: https://scrapingisnotacrime.com/docs
- Pricing: https://scrapingisnotacrime.com/#pricing
- Releases: https://github.com/ScrapingIsNotACrime/sdk-go/releases
- Python SDK: https://github.com/ScrapingIsNotACrime/sdk-python
- Node.js SDK: https://github.com/ScrapingIsNotACrime/sdk-nodejs
- MCP server: https://github.com/ScrapingIsNotACrime/mcp
- License: [MIT](./LICENSE)
