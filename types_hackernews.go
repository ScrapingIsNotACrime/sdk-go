package sinac

// HackerNewsStory is a story, as returned by feeds, search, and a user's submissions.
type HackerNewsStory struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Points      int64  `json:"points"`
	NumComments int64  `json:"num_comments"`
	// External link; null for self-posts (Ask HN, etc.).
	URL *string `json:"url,omitempty"`
	// Self-post body as HTML; null for link posts.
	Text      *string `json:"text,omitempty"`
	CreatedAt string  `json:"created_at"`
	HNURL     string  `json:"hn_url"`
}

// HackerNewsStoryPage is GET /hackernews/feeds/{feed}, /hackernews/search, and
// /hackernews/users/{username}/submissions.
type HackerNewsStoryPage struct {
	Items   []HackerNewsStory `json:"items"`
	Total   int64             `json:"total"`
	Page    int64             `json:"page"`
	HasMore bool              `json:"has_more"`
}

// HackerNewsComment is a comment inside an item's comment tree; replies nest recursively.
type HackerNewsComment struct {
	ID        int64               `json:"id"`
	Author    string              `json:"author"`
	Text      string              `json:"text"`
	CreatedAt string              `json:"created_at"`
	Replies   []HackerNewsComment `json:"replies"`
}

// HackerNewsItem is GET /hackernews/items/{id}.
type HackerNewsItem struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
	// Null for comments and other untitled item types.
	Title  *string `json:"title,omitempty"`
	Author string  `json:"author"`
	Points *int64  `json:"points,omitempty"`
	// External link; null for self-posts.
	URL *string `json:"url,omitempty"`
	// Self-post body as HTML; null for link posts.
	Text      *string             `json:"text,omitempty"`
	CreatedAt string              `json:"created_at"`
	HNURL     string              `json:"hn_url"`
	Comments  []HackerNewsComment `json:"comments"`
}

// HackerNewsUser is GET /hackernews/users/{username}.
type HackerNewsUser struct {
	Username string `json:"username"`
	Karma    int64  `json:"karma"`
	// Profile bio as HTML; null when the user has not written one.
	About           *string `json:"about,omitempty"`
	CreatedAt       string  `json:"created_at"`
	SubmissionCount int64   `json:"submission_count"`
	HNURL           string  `json:"hn_url"`
}

// HackerNewsUserComment is a comment in a user's comment listing; unlike item-tree nodes,
// Replies may be absent.
type HackerNewsUserComment struct {
	ID        int64               `json:"id"`
	Author    string              `json:"author"`
	Text      string              `json:"text"`
	CreatedAt string              `json:"created_at"`
	Replies   []HackerNewsComment `json:"replies,omitempty"`
}

// HackerNewsUserCommentPage is GET /hackernews/users/{username}/comments — same page envelope
// as the other listings.
type HackerNewsUserCommentPage struct {
	Items   []HackerNewsUserComment `json:"items"`
	Total   int64                   `json:"total"`
	Page    int64                   `json:"page"`
	HasMore bool                    `json:"has_more"`
}
