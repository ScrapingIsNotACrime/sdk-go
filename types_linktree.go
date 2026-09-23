package sinac

// LinktreeLink is one link in a Linktree profile.
type LinktreeLink struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"`
}

// LinktreeProfile is GET /linktree/profiles/{handle}.
type LinktreeProfile struct {
	Username    string         `json:"username"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Avatar      string         `json:"avatar"`
	IsVerified  bool           `json:"is_verified"`
	URL         string         `json:"url"`
	Links       []LinktreeLink `json:"links"`
}
