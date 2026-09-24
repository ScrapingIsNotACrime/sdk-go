package sinac

import "context"

// BlueskyService groups the Bluesky endpoints.
type BlueskyService struct{ http *httpCore }

// Profile is GET /bluesky/profiles/{handle} — full handle including the domain.
func (s *BlueskyService) Profile(ctx context.Context, handle string) (*BlueskyProfile, error) {
	r, err := blueskyProfile(handle)
	return getJSON[BlueskyProfile](ctx, s.http, r, err)
}

// Posts is GET /bluesky/profiles/{handle}/posts — limit 1-100 (default 25), cursor-paginated.
func (s *BlueskyService) Posts(ctx context.Context, handle string, p *BlueskyPostsParams) (*Page[BlueskyPost, BlueskyPostPage], error) {
	spec, err := blueskyPosts(handle, p)
	return fetchPage[BlueskyPost, BlueskyPostPage](ctx, s.http, spec, err)
}
