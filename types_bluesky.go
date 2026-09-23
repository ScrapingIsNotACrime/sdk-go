package sinac

// BlueskyProfile is GET /bluesky/profiles/{handle}.
type BlueskyProfile struct {
	DID         string `json:"did"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	// Null when the profile has no bio.
	Description *string `json:"description,omitempty"`
	// Null when the profile has no avatar set.
	Avatar *string `json:"avatar,omitempty"`
	// Banner image URL; null when the profile has none set.
	Banner    *string `json:"banner,omitempty"`
	Followers int64   `json:"followers"`
	Following int64   `json:"following"`
	Posts     int64   `json:"posts"`
	CreatedAt string  `json:"created_at"`
	URL       string  `json:"url"`
}

// BlueskyPost is one post in a profile's posts feed.
type BlueskyPost struct {
	URI       string `json:"uri"`
	CID       string `json:"cid"`
	Text      string `json:"text"`
	Author    string `json:"author"`
	Likes     int64  `json:"likes"`
	Reposts   int64  `json:"reposts"`
	Replies   int64  `json:"replies"`
	Quotes    int64  `json:"quotes"`
	CreatedAt string `json:"created_at"`
	IndexedAt string `json:"indexed_at"`
	URL       string `json:"url"`
}

// BlueskyPostPage is GET /bluesky/profiles/{handle}/posts.
type BlueskyPostPage struct {
	Posts      []BlueskyPost `json:"posts"`
	NextCursor *string       `json:"next_cursor,omitempty"`
	HasMore    bool          `json:"has_more"`
}
