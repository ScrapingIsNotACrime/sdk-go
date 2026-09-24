package sinac

// TwitchBroadcast is a channel's most recent broadcast.
type TwitchBroadcast struct {
	Title     string `json:"title"`
	StartedAt string `json:"started_at"`
}

// TwitchProfile is GET /twitch/profiles/{handle}.
type TwitchProfile struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	Followers   int64  `json:"followers"`
	IsPartner   bool   `json:"is_partner"`
	IsAffiliate bool   `json:"is_affiliate"`
	CreatedAt   string `json:"created_at"`
	IsLive      bool   `json:"is_live"`
	// Current viewer count while live; null whenever is_live is false.
	LiveViewers *int64 `json:"live_viewers,omitempty"`
	// Null when the channel has never broadcast (or the info is unavailable).
	LastBroadcast *TwitchBroadcast `json:"last_broadcast,omitempty"`
	URL           string           `json:"url"`
}

// TwitchVideo is one video in a channel's published videos list.
type TwitchVideo struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	DurationSeconds int64  `json:"duration_seconds"`
	Views           int64  `json:"views"`
	PublishedAt     string `json:"published_at"`
	Thumbnail       string `json:"thumbnail"`
	URL             string `json:"url"`
}

// TwitchVideos is GET /twitch/profiles/{handle}/videos.
type TwitchVideos struct {
	Videos []TwitchVideo `json:"videos"`
	Count  int64         `json:"count"`
}
