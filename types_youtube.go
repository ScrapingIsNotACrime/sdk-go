package sinac

// YouTubeChannel is the channel block in a channel-videos response.
type YouTubeChannel struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ExternalID  string `json:"externalId"`
	Avatar      string `json:"avatar"`
}

// YouTubeVideo is one video in a channel's video list.
type YouTubeVideo struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	URL          string   `json:"url"`
	Thumbnail    string   `json:"thumbnail"`
	MetadataText []string `json:"metadataText"`
}

// YouTubeChannelVideos is GET /youtube/channel/{handle}/videos.
type YouTubeChannelVideos struct {
	Channel YouTubeChannel `json:"channel"`
	Videos  []YouTubeVideo `json:"videos"`
}
