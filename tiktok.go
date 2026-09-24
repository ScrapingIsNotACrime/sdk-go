package sinac

import "context"

// TikTokService groups the TikTok endpoints.
type TikTokService struct{ http *httpCore }

// Profile is GET /tiktok/profile/{username}.
func (s *TikTokService) Profile(ctx context.Context, username string) (*TikTokProfile, error) {
	r, err := tiktokProfile(username)
	return getJSON[TikTokProfile](ctx, s.http, r, err)
}

// Video is GET /tiktok/video/{videoId}.
func (s *TikTokService) Video(ctx context.Context, videoID string) (*TikTokVideo, error) {
	r, err := tiktokVideo(videoID)
	return getJSON[TikTokVideo](ctx, s.http, r, err)
}
