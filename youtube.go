package sinac

import "context"

// YouTubeService groups the YouTube endpoints.
type YouTubeService struct{ http *httpCore }

// Videos is GET /youtube/channel/{handle}/videos.
func (s *YouTubeService) Videos(ctx context.Context, handle string) (*YouTubeChannelVideos, error) {
	r, err := youtubeVideos(handle)
	return getJSON[YouTubeChannelVideos](ctx, s.http, r, err)
}
