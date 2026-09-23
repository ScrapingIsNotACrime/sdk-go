package sinac

import "context"

// GitHubService groups the GitHub endpoints.
type GitHubService struct{ http *httpCore }

// Profile is GET /github/profiles/{handle}.
func (s *GitHubService) Profile(ctx context.Context, handle string) (*GitHubProfile, error) {
	r, err := githubProfile(handle)
	return getJSON[GitHubProfile](ctx, s.http, r, err)
}

// Followers is GET /github/profiles/{handle}/followers — limit 1-100 (default 30), 1-based pages.
func (s *GitHubService) Followers(ctx context.Context, handle string, p *GitHubListParams) (*Page[GitHubUser, GitHubUserPage], error) {
	spec, err := githubFollowers(handle, p)
	return fetchPage[GitHubUser, GitHubUserPage](ctx, s.http, spec, err)
}

// Following is GET /github/profiles/{handle}/following — limit 1-100 (default 30), 1-based pages.
func (s *GitHubService) Following(ctx context.Context, handle string, p *GitHubListParams) (*Page[GitHubUser, GitHubUserPage], error) {
	spec, err := githubFollowing(handle, p)
	return fetchPage[GitHubUser, GitHubUserPage](ctx, s.http, spec, err)
}

// Repositories is GET /github/profiles/{handle}/repositories — limit 1-100 (default 30), 1-based pages.
func (s *GitHubService) Repositories(ctx context.Context, handle string, p *GitHubListParams) (*Page[GitHubRepository, GitHubRepositoryPage], error) {
	spec, err := githubRepositories(handle, p)
	return fetchPage[GitHubRepository, GitHubRepositoryPage](ctx, s.http, spec, err)
}

// SearchRepositories is GET /github/repositories — q in GitHub search syntax; limit 1-100 (default 30), 1-based pages.
func (s *GitHubService) SearchRepositories(ctx context.Context, q string, p *GitHubListParams) (*Page[GitHubRepository, GitHubRepositorySearchPage], error) {
	spec, err := githubSearchRepositories(q, p)
	return fetchPage[GitHubRepository, GitHubRepositorySearchPage](ctx, s.http, spec, err)
}

// Trending is GET /github/trending/repositories — Since defaults to "daily"; Limit 1-100 (default 30). Not paginated.
func (s *GitHubService) Trending(ctx context.Context, p *GitHubTrendingParams) (*GitHubTrending, error) {
	r, err := githubTrending(p)
	return getJSON[GitHubTrending](ctx, s.http, r, err)
}
