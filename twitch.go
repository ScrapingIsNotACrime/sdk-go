package sinac

import "context"

// TwitchService groups the Twitch endpoints.
type TwitchService struct{ http *httpCore }

// Profile is GET /twitch/profiles/{handle}.
func (s *TwitchService) Profile(ctx context.Context, handle string) (*TwitchProfile, error) {
	r, err := twitchProfile(handle)
	return getJSON[TwitchProfile](ctx, s.http, r, err)
}

// Videos is GET /twitch/profiles/{handle}/videos — limit 1-100, default 20.
func (s *TwitchService) Videos(ctx context.Context, handle string, p *TwitchVideosParams) (*TwitchVideos, error) {
	r, err := twitchVideos(handle, p)
	return getJSON[TwitchVideos](ctx, s.http, r, err)
}
