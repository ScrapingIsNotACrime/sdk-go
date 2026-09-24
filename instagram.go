package sinac

import "context"

// InstagramService groups the Instagram endpoints.
type InstagramService struct{ http *httpCore }

// Profile is GET /instagram/profile/{username} — username without @.
func (s *InstagramService) Profile(ctx context.Context, username string) (*InstagramProfile, error) {
	r, err := instagramProfile(username)
	return getJSON[InstagramProfile](ctx, s.http, r, err)
}

// Contact is GET /instagram/profile/{username}/contact — public business contact (email, phone, address).
func (s *InstagramService) Contact(ctx context.Context, username string) (*InstagramContact, error) {
	r, err := instagramContact(username)
	return getJSON[InstagramContact](ctx, s.http, r, err)
}

// LatestPosts is GET /instagram/profile/{username}/timeline/latest — the first page of posts.
func (s *InstagramService) LatestPosts(ctx context.Context, username string) (*InstagramLatestPosts, error) {
	r, err := instagramLatestPosts(username)
	return getJSON[InstagramLatestPosts](ctx, s.http, r, err)
}

// Posts is GET /instagram/profile/{username}/timeline — full history, cursor-paginated; Count 1-50 (default 12).
func (s *InstagramService) Posts(ctx context.Context, username string, p *InstagramPostsParams) (*Page[InstagramMedia, InstagramTimelinePage], error) {
	spec, err := instagramPosts(username, p)
	return fetchPage[InstagramMedia, InstagramTimelinePage](ctx, s.http, spec, err)
}

// Highlights is GET /instagram/profile/{username}/highlights.
func (s *InstagramService) Highlights(ctx context.Context, username string) (*InstagramHighlights, error) {
	r, err := instagramHighlights(username)
	return getJSON[InstagramHighlights](ctx, s.http, r, err)
}

// Highlight is GET /instagram/highlights/{highlightId}.
func (s *InstagramService) Highlight(ctx context.Context, highlightID string) (*InstagramHighlight, error) {
	r, err := instagramHighlight(highlightID)
	return getJSON[InstagramHighlight](ctx, s.http, r, err)
}

// MediaByID is GET /instagram/profile/{username}/media/{mediaId}.
func (s *InstagramService) MediaByID(ctx context.Context, username, mediaID string) (*InstagramMediaDetail, error) {
	r, err := instagramMediaByID(username, mediaID)
	return getJSON[InstagramMediaDetail](ctx, s.http, r, err)
}

// Media is GET /instagram/media/{shortcode} — shortcode from instagram.com/p/{shortcode}/.
func (s *InstagramService) Media(ctx context.Context, shortcode string) (*InstagramMediaDetail, error) {
	r, err := instagramMedia(shortcode)
	return getJSON[InstagramMediaDetail](ctx, s.http, r, err)
}

// Download is GET /instagram/media/{shortcode}/download — assets[0] is the best primary asset.
func (s *InstagramService) Download(ctx context.Context, shortcode string) (*InstagramDownload, error) {
	r, err := instagramDownload(shortcode)
	return getJSON[InstagramDownload](ctx, s.http, r, err)
}

// ShortcodeToID is GET /instagram/media/{shortcode}/id.
func (s *InstagramService) ShortcodeToID(ctx context.Context, shortcode string) (*InstagramShortcodeID, error) {
	r, err := instagramShortcodeToID(shortcode)
	return getJSON[InstagramShortcodeID](ctx, s.http, r, err)
}

// IDToShortcode is GET /instagram/media/id/{mediaId}.
func (s *InstagramService) IDToShortcode(ctx context.Context, mediaID string) (*InstagramShortcodeID, error) {
	r, err := instagramIDToShortcode(mediaID)
	return getJSON[InstagramShortcodeID](ctx, s.http, r, err)
}

// Reel is GET /instagram/reels/{shortcode}.
func (s *InstagramService) Reel(ctx context.Context, shortcode string) (*InstagramReel, error) {
	r, err := instagramReel(shortcode)
	return getJSON[InstagramReel](ctx, s.http, r, err)
}
