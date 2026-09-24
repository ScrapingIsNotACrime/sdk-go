package sinac

import "strconv"

// GitHubTrendingSince is the period of GitHub's trending list.
type GitHubTrendingSince string

const (
	TrendingDaily   GitHubTrendingSince = "daily"
	TrendingWeekly  GitHubTrendingSince = "weekly"
	TrendingMonthly GitHubTrendingSince = "monthly"
)

// HackerNewsFeed is one of Hacker News' story lists.
type HackerNewsFeed string

const (
	FeedTop  HackerNewsFeed = "top"
	FeedNew  HackerNewsFeed = "new"
	FeedBest HackerNewsFeed = "best"
	FeedAsk  HackerNewsFeed = "ask"
	FeedShow HackerNewsFeed = "show"
	FeedJob  HackerNewsFeed = "job"
)

// InstagramPostsParams: Count 1-50 (default 12); Cursor from Page.NextCursor.
type InstagramPostsParams struct {
	Count  int
	Cursor string
}

// AppStoreSearchParams: Country is a two-letter store code (default "us"); Limit 1-200 (default 10).
type AppStoreSearchParams struct {
	Country string
	Limit   int
}

// AppStoreReviewsParams: Country (default "us"); Page 1-10 (default 1).
type AppStoreReviewsParams struct {
	Country string
	Page    int
}

// GitHubListParams: Limit 1-100 (default 30); Page is 1-based (default 1).
type GitHubListParams struct {
	Limit int
	Page  int
}

// GitHubTrendingParams: Since (default daily), Language (e.g. "go"), Limit 1-100 (default 30).
type GitHubTrendingParams struct {
	Since    GitHubTrendingSince
	Language string
	Limit    int
}

// HackerNewsListParams: Limit 1-50 (default 20); Page is 0-based (default 0).
type HackerNewsListParams struct {
	Limit int
	Page  int
}

// BlueskyPostsParams: Limit 1-100 (default 25); Cursor from Page.NextCursor.
type BlueskyPostsParams struct {
	Limit  int
	Cursor string
}

// TwitchVideosParams: Limit 1-100 (default 20).
type TwitchVideosParams struct {
	Limit int
}

func simple(format string, values ...string) (route, error) {
	path, err := pathOf(format, values...)
	return route{path: path}, err
}

func numbered(path string, err error, itemsKey string, page, firstPage, maxPage int, query ...[]queryParam) (pageSpec, error) {
	if page == 0 {
		page = firstPage
	}
	spec := pageSpec{path: path, kind: numberedPages, itemsKey: itemsKey, page: page, maxPage: maxPage}
	for _, q := range query {
		spec.query = append(spec.query, q...)
	}
	return spec, err
}

// Instagram
func instagramProfile(username string) (route, error) {
	return simple("/instagram/profile/%s", username)
}
func instagramContact(username string) (route, error) {
	return simple("/instagram/profile/%s/contact", username)
}
func instagramLatestPosts(username string) (route, error) {
	return simple("/instagram/profile/%s/timeline/latest", username)
}
func instagramPosts(username string, p *InstagramPostsParams) (pageSpec, error) {
	if p == nil {
		p = &InstagramPostsParams{}
	}
	path, err := pathOf("/instagram/profile/%s/timeline", username)
	return pageSpec{path: path, query: qi("count", p.Count), kind: cursorPages, itemsKey: "medias", cursor: p.Cursor}, err
}
func instagramHighlights(username string) (route, error) {
	return simple("/instagram/profile/%s/highlights", username)
}
func instagramHighlight(highlightID string) (route, error) {
	return simple("/instagram/highlights/%s", highlightID)
}
func instagramMediaByID(username, mediaID string) (route, error) {
	return simple("/instagram/profile/%s/media/%s", username, mediaID)
}
func instagramMedia(shortcode string) (route, error) { return simple("/instagram/media/%s", shortcode) }
func instagramDownload(shortcode string) (route, error) {
	return simple("/instagram/media/%s/download", shortcode)
}
func instagramShortcodeToID(shortcode string) (route, error) {
	return simple("/instagram/media/%s/id", shortcode)
}
func instagramIDToShortcode(mediaID string) (route, error) {
	return simple("/instagram/media/id/%s", mediaID)
}
func instagramReel(shortcode string) (route, error) { return simple("/instagram/reels/%s", shortcode) }

// TikTok
func tiktokProfile(username string) (route, error) { return simple("/tiktok/profile/%s", username) }
func tiktokVideo(videoID string) (route, error)    { return simple("/tiktok/video/%s", videoID) }

// YouTube
func youtubeVideos(handle string) (route, error) { return simple("/youtube/channel/%s/videos", handle) }

// App Store
func appstoreSearch(term string, p *AppStoreSearchParams) (route, error) {
	if p == nil {
		p = &AppStoreSearchParams{}
	}
	return route{path: "/appstore/search", query: []queryParam{{"term", term}}}.with(qs("country", p.Country), qi("limit", p.Limit)), nil
}
func appstoreReviews(appID string, p *AppStoreReviewsParams) (pageSpec, error) {
	if p == nil {
		p = &AppStoreReviewsParams{}
	}
	// Apple's feed caps at 10 pages and the payload has no has_more flag.
	return numbered("/appstore/reviews", nil, "reviews", p.Page, 1, 10, []queryParam{{"appId", appID}}, qs("country", p.Country))
}

// GitHub (1-based pages)
func githubProfile(handle string) (route, error) { return simple("/github/profiles/%s", handle) }
func githubList(format, handle string, p *GitHubListParams) (pageSpec, error) {
	if p == nil {
		p = &GitHubListParams{}
	}
	path, err := pathOf(format, handle)
	return numbered(path, err, "items", p.Page, 1, 0, qi("limit", p.Limit))
}
func githubFollowers(handle string, p *GitHubListParams) (pageSpec, error) {
	return githubList("/github/profiles/%s/followers", handle, p)
}
func githubFollowing(handle string, p *GitHubListParams) (pageSpec, error) {
	return githubList("/github/profiles/%s/following", handle, p)
}
func githubRepositories(handle string, p *GitHubListParams) (pageSpec, error) {
	return githubList("/github/profiles/%s/repositories", handle, p)
}
func githubSearchRepositories(q string, p *GitHubListParams) (pageSpec, error) {
	if p == nil {
		p = &GitHubListParams{}
	}
	return numbered("/github/repositories", nil, "items", p.Page, 1, 0, []queryParam{{"q", q}}, qi("limit", p.Limit))
}
func githubTrending(p *GitHubTrendingParams) (route, error) {
	if p == nil {
		p = &GitHubTrendingParams{}
	}
	return route{path: "/github/trending/repositories"}.with(qs("since", string(p.Since)), qs("language", p.Language), qi("limit", p.Limit)), nil
}

// Hacker News (0-based pages)
func hackernewsList(path string, err error, p *HackerNewsListParams, query ...[]queryParam) (pageSpec, error) {
	if p == nil {
		p = &HackerNewsListParams{}
	}
	return numbered(path, err, "items", p.Page, 0, 0, append(query, qi("limit", p.Limit))...)
}
func hackernewsFeed(feed HackerNewsFeed, p *HackerNewsListParams) (pageSpec, error) {
	path, err := pathOf("/hackernews/feeds/%s", string(feed))
	return hackernewsList(path, err, p)
}
func hackernewsItem(id int64) (route, error) {
	return simple("/hackernews/items/%s", strconv.FormatInt(id, 10))
}
func hackernewsSearch(q string, p *HackerNewsListParams) (pageSpec, error) {
	return hackernewsList("/hackernews/search", nil, p, []queryParam{{"q", q}})
}
func hackernewsUser(username string) (route, error) { return simple("/hackernews/users/%s", username) }
func hackernewsSubmissions(username string, p *HackerNewsListParams) (pageSpec, error) {
	path, err := pathOf("/hackernews/users/%s/submissions", username)
	return hackernewsList(path, err, p)
}
func hackernewsComments(username string, p *HackerNewsListParams) (pageSpec, error) {
	path, err := pathOf("/hackernews/users/%s/comments", username)
	return hackernewsList(path, err, p)
}

// Bluesky
func blueskyProfile(handle string) (route, error) { return simple("/bluesky/profiles/%s", handle) }
func blueskyPosts(handle string, p *BlueskyPostsParams) (pageSpec, error) {
	if p == nil {
		p = &BlueskyPostsParams{}
	}
	path, err := pathOf("/bluesky/profiles/%s/posts", handle)
	return pageSpec{path: path, query: qi("limit", p.Limit), kind: cursorPages, itemsKey: "posts", cursor: p.Cursor}, err
}

// Twitch
func twitchProfile(handle string) (route, error) { return simple("/twitch/profiles/%s", handle) }
func twitchVideos(handle string, p *TwitchVideosParams) (route, error) {
	if p == nil {
		p = &TwitchVideosParams{}
	}
	r, err := simple("/twitch/profiles/%s/videos", handle)
	return r.with(qi("limit", p.Limit)), err
}

// Linktree
func linktreeProfile(handle string) (route, error) { return simple("/linktree/profiles/%s", handle) }
