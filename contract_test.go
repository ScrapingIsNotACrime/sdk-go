package sinac

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type contractCase struct {
	name     string // "GitHub.Followers"
	fixture  string // "" = placeholder
	request  string // expected path?query, as in the fixture
	itemsKey string // "" = not paginated; otherwise the fixture/placeholder key holding the items array
	call     func(ctx context.Context, c *Client) (any, error)
}

var placeholders = map[string]string{
	"Instagram.Media":     `{"id":"1","shortcode":"C8xQz1aP9Kv"}`,
	"TikTok.Video":        `{"id":"7300000000000000000"}`,
	"GitHub.Following":    `{"items":[],"total":null,"has_more":false}`,
	"HackerNews.Comments": `{"items":[],"total":0,"page":0,"has_more":false}`,
}

var contractCases = []contractCase{
	{"Instagram.Profile", "ig-profile", "/instagram/profile/instagram", "",
		func(ctx context.Context, c *Client) (any, error) { return c.Instagram.Profile(ctx, "instagram") }},
	{"Instagram.Contact", "ig-contact", "/instagram/profile/cafedaesquina/contact", "",
		func(ctx context.Context, c *Client) (any, error) { return c.Instagram.Contact(ctx, "cafedaesquina") }},
	{"Instagram.LatestPosts", "ig-timeline", "/instagram/profile/instagram/timeline/latest", "",
		func(ctx context.Context, c *Client) (any, error) { return c.Instagram.LatestPosts(ctx, "instagram") }},
	{"Instagram.Posts", "ig-timeline-paged", "/instagram/profile/nasa/timeline?count=12&cursor=3950671748375397992_528817151", "medias",
		func(ctx context.Context, c *Client) (any, error) {
			return c.Instagram.Posts(ctx, "nasa", &InstagramPostsParams{Count: 12, Cursor: "3950671748375397992_528817151"})
		}},
	{"Instagram.Highlights", "ig-highlights", "/instagram/profile/nasa/highlights", "",
		func(ctx context.Context, c *Client) (any, error) { return c.Instagram.Highlights(ctx, "nasa") }},
	{"Instagram.Highlight", "ig-highlight-content", "/instagram/highlights/highlight%3A18201653992314974", "",
		func(ctx context.Context, c *Client) (any, error) {
			return c.Instagram.Highlight(ctx, "highlight:18201653992314974")
		}},
	{"Instagram.MediaByID", "ig-media-by-id", "/instagram/profile/instagram/media/3123456789012345678", "",
		func(ctx context.Context, c *Client) (any, error) {
			return c.Instagram.MediaByID(ctx, "instagram", "3123456789012345678")
		}},
	{"Instagram.Media", "", "/instagram/media/C8xQz1aP9Kv", "",
		func(ctx context.Context, c *Client) (any, error) { return c.Instagram.Media(ctx, "C8xQz1aP9Kv") }},
	{"Instagram.Download", "ig-media-download", "/instagram/media/DbtErSrlB2J/download", "",
		func(ctx context.Context, c *Client) (any, error) { return c.Instagram.Download(ctx, "DbtErSrlB2J") }},
	{"Instagram.ShortcodeToID", "ig-shortcode-to-id", "/instagram/media/Dbn-XJhk0_-/id", "",
		func(ctx context.Context, c *Client) (any, error) {
			return c.Instagram.ShortcodeToID(ctx, "Dbn-XJhk0_-")
		}},
	{"Instagram.IDToShortcode", "ig-id-to-shortcode", "/instagram/media/id/3956405067326902270", "",
		func(ctx context.Context, c *Client) (any, error) {
			return c.Instagram.IDToShortcode(ctx, "3956405067326902270")
		}},
	{"Instagram.Reel", "ig-reels", "/instagram/reels/DyKlMnOpQrS", "",
		func(ctx context.Context, c *Client) (any, error) { return c.Instagram.Reel(ctx, "DyKlMnOpQrS") }},
	{"TikTok.Profile", "tt-profile", "/tiktok/profile/tiktok", "",
		func(ctx context.Context, c *Client) (any, error) { return c.TikTok.Profile(ctx, "tiktok") }},
	{"TikTok.Video", "", "/tiktok/video/7300000000000000000", "",
		func(ctx context.Context, c *Client) (any, error) { return c.TikTok.Video(ctx, "7300000000000000000") }},
	{"YouTube.Videos", "yt-channel-videos", "/youtube/channel/youtube/videos", "",
		func(ctx context.Context, c *Client) (any, error) { return c.YouTube.Videos(ctx, "youtube") }},
	{"AppStore.Search", "as-search", "/appstore/search?term=instagram&country=us&limit=1", "",
		func(ctx context.Context, c *Client) (any, error) {
			return c.AppStore.Search(ctx, "instagram", &AppStoreSearchParams{Country: "us", Limit: 1})
		}},
	{"AppStore.Reviews", "as-reviews", "/appstore/reviews?appId=389801252&country=us&page=1", "reviews",
		func(ctx context.Context, c *Client) (any, error) {
			return c.AppStore.Reviews(ctx, "389801252", &AppStoreReviewsParams{Country: "us", Page: 1})
		}},
	{"GitHub.Profile", "gh-profile", "/github/profiles/torvalds", "",
		func(ctx context.Context, c *Client) (any, error) { return c.GitHub.Profile(ctx, "torvalds") }},
	{"GitHub.Followers", "gh-followers", "/github/profiles/torvalds/followers?limit=30&page=1", "items",
		func(ctx context.Context, c *Client) (any, error) {
			return c.GitHub.Followers(ctx, "torvalds", &GitHubListParams{Limit: 30, Page: 1})
		}},
	{"GitHub.Following", "", "/github/profiles/torvalds/following?limit=5&page=1", "items",
		func(ctx context.Context, c *Client) (any, error) {
			return c.GitHub.Following(ctx, "torvalds", &GitHubListParams{Limit: 5})
		}},
	{"GitHub.Repositories", "gh-repos", "/github/profiles/torvalds/repositories?limit=30&page=1", "items",
		func(ctx context.Context, c *Client) (any, error) {
			return c.GitHub.Repositories(ctx, "torvalds", &GitHubListParams{Limit: 30})
		}},
	{"GitHub.SearchRepositories", "gh-search-repos", "/github/repositories?q=stars:%3E10000+language:php&limit=1&page=1", "items",
		func(ctx context.Context, c *Client) (any, error) {
			return c.GitHub.SearchRepositories(ctx, "stars:>10000 language:php", &GitHubListParams{Limit: 1})
		}},
	{"GitHub.Trending", "gh-trending", "/github/trending/repositories?since=weekly&language=php&limit=1", "",
		func(ctx context.Context, c *Client) (any, error) {
			return c.GitHub.Trending(ctx, &GitHubTrendingParams{Since: TrendingWeekly, Language: "php", Limit: 1})
		}},
	{"HackerNews.Feed", "hn-feed", "/hackernews/feeds/top?limit=20&page=0", "items",
		func(ctx context.Context, c *Client) (any, error) {
			return c.HackerNews.Feed(ctx, FeedTop, &HackerNewsListParams{Limit: 20, Page: 0})
		}},
	{"HackerNews.Item", "hn-item", "/hackernews/items/8863", "",
		func(ctx context.Context, c *Client) (any, error) { return c.HackerNews.Item(ctx, 8863) }},
	{"HackerNews.Search", "hn-search", "/hackernews/search?q=postgres&limit=20&page=0", "items",
		func(ctx context.Context, c *Client) (any, error) {
			return c.HackerNews.Search(ctx, "postgres", &HackerNewsListParams{Limit: 20, Page: 0})
		}},
	{"HackerNews.User", "hn-user", "/hackernews/users/pg", "",
		func(ctx context.Context, c *Client) (any, error) { return c.HackerNews.User(ctx, "pg") }},
	{"HackerNews.Submissions", "hn-user-submissions", "/hackernews/users/pg/submissions?limit=20&page=0", "items",
		func(ctx context.Context, c *Client) (any, error) {
			return c.HackerNews.Submissions(ctx, "pg", &HackerNewsListParams{Limit: 20, Page: 0})
		}},
	{"HackerNews.Comments", "", "/hackernews/users/pg/comments?limit=10&page=0", "items",
		func(ctx context.Context, c *Client) (any, error) {
			return c.HackerNews.Comments(ctx, "pg", &HackerNewsListParams{Limit: 10})
		}},
	{"Bluesky.Profile", "bs-profile", "/bluesky/profiles/bsky.app", "",
		func(ctx context.Context, c *Client) (any, error) { return c.Bluesky.Profile(ctx, "bsky.app") }},
	{"Bluesky.Posts", "bs-posts", "/bluesky/profiles/bsky.app/posts?limit=25", "posts",
		func(ctx context.Context, c *Client) (any, error) {
			return c.Bluesky.Posts(ctx, "bsky.app", &BlueskyPostsParams{Limit: 25})
		}},
	{"Twitch.Profile", "tw-profile", "/twitch/profiles/ninja", "",
		func(ctx context.Context, c *Client) (any, error) { return c.Twitch.Profile(ctx, "ninja") }},
	{"Twitch.Videos", "tw-videos", "/twitch/profiles/ninja/videos?limit=20", "",
		func(ctx context.Context, c *Client) (any, error) {
			return c.Twitch.Videos(ctx, "ninja", &TwitchVideosParams{Limit: 20})
		}},
	{"Linktree.Profile", "lt-profile", "/linktree/profiles/linktree", "",
		func(ctx context.Context, c *Client) (any, error) { return c.Linktree.Profile(ctx, "linktree") }},
}

// queryPairs decodes a raw query into ordered key/value pairs, so encoding
// differences (":" vs "%3A") don't matter but order and values do.
func queryPairs(raw string) [][2]string {
	var pairs [][2]string
	for _, part := range strings.Split(raw, "&") {
		if part == "" {
			continue
		}
		k, v, _ := strings.Cut(part, "=")
		k, _ = url.QueryUnescape(k)
		v, _ = url.QueryUnescape(v)
		pairs = append(pairs, [2]string{k, v})
	}
	return pairs
}

func TestContractCoversAllMethods(t *testing.T) {
	names := map[string]bool{}
	for _, c := range contractCases {
		names[c.name] = true
	}
	if len(names) != 34 || len(contractCases) != 34 {
		t.Fatalf("cases=%d unique=%d", len(contractCases), len(names))
	}
}

func TestContract(t *testing.T) {
	for _, tc := range contractCases {
		t.Run(tc.name, func(t *testing.T) {
			data := json.RawMessage(placeholders[tc.name])
			if tc.fixture != "" {
				data = loadFixture(t, tc.fixture).Response.Data
			}
			var got *http.Request
			client, err := NewClient(WithAPIKey("sinac_test"), WithMaxRetries(0), WithHTTPClient(&http.Client{
				Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
					got = r
					body, _ := json.Marshal(map[string]any{"message": "ok", "data": data})
					return response(200, string(body), nil), nil
				}),
			}))
			if err != nil {
				t.Fatal(err)
			}
			result, err := tc.call(context.Background(), client)
			if err != nil {
				t.Fatal(err)
			}
			wantPath, wantQuery, _ := strings.Cut(tc.request, "?")
			gotPath := strings.TrimPrefix(got.URL.EscapedPath(), "/v1")
			if gotPath != wantPath || !reflect.DeepEqual(queryPairs(got.URL.RawQuery), queryPairs(wantQuery)) {
				t.Fatalf("request = %s?%s, want %s", gotPath, got.URL.RawQuery, tc.request)
			}
			// The returned value re-encodes to the documented data (paginated: its Data field).
			value := reflect.ValueOf(result)
			elem := value.Elem()
			if strings.HasPrefix(elem.Type().Name(), "Page[") {
				if tc.itemsKey != "" {
					assertItemsLen(t, data, tc.itemsKey, elem.FieldByName("Items").Len())
				}
				value = elem.FieldByName("Data")
			}
			assertSameJSON(t, data, value.Interface())
		})
	}
}

// assertItemsLen checks that the Page's Items has as many elements as the
// fixture's (or placeholder's) items array, under itemsKey.
func assertItemsLen(t *testing.T, data json.RawMessage, itemsKey string, gotLen int) {
	t.Helper()
	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		t.Fatal(err)
	}
	var items []json.RawMessage
	_ = json.Unmarshal(object[itemsKey], &items)
	if gotLen != len(items) {
		t.Fatalf("len(Items) = %d, want %d (key %q)", gotLen, len(items), itemsKey)
	}
}

// assertSameJSON checks that every key/value in want is present with the same value in got's encoding.
func assertSameJSON(t *testing.T, want json.RawMessage, got any) {
	t.Helper()
	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	var w, g any
	_ = json.Unmarshal(want, &w)
	_ = json.Unmarshal(encoded, &g)
	if diff := subset(w, g, "data"); diff != "" {
		t.Fatal(diff)
	}
}

// isEmptyJSON reports null, [] or {} — the values omitempty leaves out.
func isEmptyJSON(v any) bool {
	switch x := v.(type) {
	case nil:
		return true
	case []any:
		return len(x) == 0
	case map[string]any:
		return len(x) == 0
	}
	return false
}

// subset reports the first place where want is not contained in got. Keys that
// are null in want may be absent in got (pointer fields use omitempty).
func subset(want, got any, at string) string {
	switch w := want.(type) {
	case map[string]any:
		g, ok := got.(map[string]any)
		if !ok {
			return at + ": not an object"
		}
		for k, wv := range w {
			gv, present := g[k]
			if !present {
				// omitempty drops null, empty arrays and empty objects on re-encoding
				if isEmptyJSON(wv) {
					continue
				}
				return at + "." + k + ": missing"
			}
			if d := subset(wv, gv, at+"."+k); d != "" {
				return d
			}
		}
	case []any:
		g, ok := got.([]any)
		if !ok || len(g) != len(w) {
			return at + ": array length differs"
		}
		for i := range w {
			if d := subset(w[i], g[i], at); d != "" {
				return d
			}
		}
	default:
		if !reflect.DeepEqual(want, got) {
			return at + ": value differs"
		}
	}
	return ""
}

// TestPageNextKeepsQueryAndAdvances checks, for one contract case per
// pagination strategy, that Next's request keeps the method's original query
// and advances the page/cursor: numbered-with-has_more (GitHub), numbered
// capped without has_more (App Store), cursor (Bluesky) and a 0-based
// numbered feed with no explicit params (Hacker News).
func TestPageNextKeepsQueryAndAdvances(t *testing.T) {
	cases := []struct {
		name  string
		body1 map[string]any
		call  func(ctx context.Context, c *Client) (any, error)
		want  string // second request, as path?query
	}{
		{
			name:  "GitHub.SearchRepositories",
			body1: map[string]any{"items": []any{map[string]any{"name": "a"}}, "has_more": true},
			call: func(ctx context.Context, c *Client) (any, error) {
				return c.GitHub.SearchRepositories(ctx, "go", &GitHubListParams{Limit: 1})
			},
			want: "/github/repositories?q=go&limit=1&page=2",
		},
		{
			name:  "AppStore.Reviews",
			body1: map[string]any{"reviews": []any{map[string]any{"id": "1"}}},
			call: func(ctx context.Context, c *Client) (any, error) {
				return c.AppStore.Reviews(ctx, "1", &AppStoreReviewsParams{Country: "br"})
			},
			want: "/appstore/reviews?appId=1&country=br&page=2",
		},
		{
			name:  "Bluesky.Posts",
			body1: map[string]any{"posts": []any{map[string]any{"uri": "x"}}, "has_more": true, "next_cursor": "c2"},
			call: func(ctx context.Context, c *Client) (any, error) {
				return c.Bluesky.Posts(ctx, "bsky.app", &BlueskyPostsParams{Limit: 5})
			},
			want: "/bluesky/profiles/bsky.app/posts?limit=5&cursor=c2",
		},
		{
			name:  "HackerNews.Search",
			body1: map[string]any{"items": []any{map[string]any{"id": 1}}, "has_more": true},
			call: func(ctx context.Context, c *Client) (any, error) {
				return c.HackerNews.Search(ctx, "go", nil)
			},
			want: "/hackernews/search?q=go&page=1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var reqs []*http.Request
			client, err := NewClient(WithAPIKey("sinac_test"), WithMaxRetries(0), WithHTTPClient(&http.Client{
				Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
					reqs = append(reqs, r)
					body, _ := json.Marshal(map[string]any{"message": "ok", "data": tc.body1})
					return response(200, string(body), nil), nil
				}),
			}))
			if err != nil {
				t.Fatal(err)
			}
			result, err := tc.call(context.Background(), client)
			if err != nil {
				t.Fatal(err)
			}
			out := reflect.ValueOf(result).MethodByName("Next").Call([]reflect.Value{reflect.ValueOf(context.Background())})
			if !out[1].IsNil() {
				t.Fatalf("Next err = %v", out[1].Interface())
			}
			if len(reqs) != 2 {
				t.Fatalf("requests = %d, want 2", len(reqs))
			}
			wantPath, wantQuery, _ := strings.Cut(tc.want, "?")
			gotPath := strings.TrimPrefix(reqs[1].URL.EscapedPath(), "/v1")
			if gotPath != wantPath || !reflect.DeepEqual(queryPairs(reqs[1].URL.RawQuery), queryPairs(wantQuery)) {
				t.Fatalf("second request = %s?%s, want %s", gotPath, reqs[1].URL.RawQuery, tc.want)
			}
		})
	}
}
