package sinac

// GitHubProfile is GET /github/profiles/{handle}.
type GitHubProfile struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
	// Null when the user has not set a display name.
	Name *string `json:"name,omitempty"`
	// Null when the user has not set a bio.
	Bio      *string `json:"bio,omitempty"`
	Company  *string `json:"company,omitempty"`
	Location *string `json:"location,omitempty"`
	// Website URL; empty or null when the user has not set one.
	Blog        *string `json:"blog,omitempty"`
	PublicRepos int64   `json:"public_repos"`
	Followers   int64   `json:"followers"`
	Following   int64   `json:"following"`
	Avatar      string  `json:"avatar"`
	CreatedAt   string  `json:"created_at"`
	URL         string  `json:"url"`
}

// GitHubUser is a user in a followers/following page.
type GitHubUser struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
	Avatar   string `json:"avatar"`
	URL      string `json:"url"`
}

// GitHubUserPage is GET /github/profiles/{handle}/followers and /following.
type GitHubUserPage struct {
	Items []GitHubUser `json:"items"`
	// Always null: GitHub's REST API does not report a count for this collection.
	Total   *int64 `json:"total,omitempty"`
	HasMore bool   `json:"has_more"`
}

// GitHubRepository is a repository, as returned by the profile repositories list, search, and
// trending endpoints.
type GitHubRepository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	// Null when the repository has no description.
	Description *string `json:"description,omitempty"`
	Stars       int64   `json:"stars"`
	Forks       int64   `json:"forks"`
	// Primary language; null when GitHub has not detected one.
	Language   *string  `json:"language,omitempty"`
	Topics     []string `json:"topics"`
	IsFork     bool     `json:"is_fork"`
	IsArchived bool     `json:"is_archived"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	URL        string   `json:"url"`
}

// GitHubRepositoryPage is GET /github/profiles/{handle}/repositories.
type GitHubRepositoryPage struct {
	Items []GitHubRepository `json:"items"`
	// Always null: GitHub's REST API does not report a count for this collection.
	Total   *int64 `json:"total,omitempty"`
	HasMore bool   `json:"has_more"`
}

// GitHubRepositorySearchPage is GET /github/repositories.
type GitHubRepositorySearchPage struct {
	Items   []GitHubRepository `json:"items"`
	Total   int64              `json:"total"`
	Page    int64              `json:"page"`
	HasMore bool               `json:"has_more"`
}

// GitHubTrending is GET /github/trending/repositories.
type GitHubTrending struct {
	Items   []GitHubRepository `json:"items"`
	Total   int64              `json:"total"`
	Page    int64              `json:"page"`
	HasMore bool               `json:"has_more"`
}
