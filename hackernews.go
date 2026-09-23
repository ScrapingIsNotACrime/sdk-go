package sinac

import "context"

// HackerNewsService groups the Hacker News endpoints.
type HackerNewsService struct{ http *httpCore }

// Feed is GET /hackernews/feeds/{feed} — limit 1-50 (default 20), 0-based pages.
func (s *HackerNewsService) Feed(ctx context.Context, feed HackerNewsFeed, p *HackerNewsListParams) (*Page[HackerNewsStory, HackerNewsStoryPage], error) {
	spec, err := hackernewsFeed(feed, p)
	return fetchPage[HackerNewsStory, HackerNewsStoryPage](ctx, s.http, spec, err)
}

// Item is GET /hackernews/items/{id} — the item with its full comment tree.
func (s *HackerNewsService) Item(ctx context.Context, id int64) (*HackerNewsItem, error) {
	r, err := hackernewsItem(id)
	return getJSON[HackerNewsItem](ctx, s.http, r, err)
}

// Search is GET /hackernews/search — limit 1-50 (default 20), 0-based pages.
func (s *HackerNewsService) Search(ctx context.Context, q string, p *HackerNewsListParams) (*Page[HackerNewsStory, HackerNewsStoryPage], error) {
	spec, err := hackernewsSearch(q, p)
	return fetchPage[HackerNewsStory, HackerNewsStoryPage](ctx, s.http, spec, err)
}

// User is GET /hackernews/users/{username}.
func (s *HackerNewsService) User(ctx context.Context, username string) (*HackerNewsUser, error) {
	r, err := hackernewsUser(username)
	return getJSON[HackerNewsUser](ctx, s.http, r, err)
}

// Submissions is GET /hackernews/users/{username}/submissions — limit 1-50 (default 20), 0-based pages.
func (s *HackerNewsService) Submissions(ctx context.Context, username string, p *HackerNewsListParams) (*Page[HackerNewsStory, HackerNewsStoryPage], error) {
	spec, err := hackernewsSubmissions(username, p)
	return fetchPage[HackerNewsStory, HackerNewsStoryPage](ctx, s.http, spec, err)
}

// Comments is GET /hackernews/users/{username}/comments — limit 1-50 (default 20), 0-based pages.
func (s *HackerNewsService) Comments(ctx context.Context, username string, p *HackerNewsListParams) (*Page[HackerNewsUserComment, HackerNewsUserCommentPage], error) {
	spec, err := hackernewsComments(username, p)
	return fetchPage[HackerNewsUserComment, HackerNewsUserCommentPage](ctx, s.http, spec, err)
}
