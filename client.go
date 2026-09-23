package sinac

import "os"

// Client is the ScrapingIsNotACrime API client. It is safe for concurrent use.
type Client struct {
	Instagram  *InstagramService
	TikTok     *TikTokService
	YouTube    *YouTubeService
	AppStore   *AppStoreService
	GitHub     *GitHubService
	HackerNews *HackerNewsService
	Bluesky    *BlueskyService
	Twitch     *TwitchService
	Linktree   *LinktreeService
}

// NewClient builds a client. It returns an error for a missing API key (pass
// WithAPIKey or set SCRAPINGISNOTACRIME_API_KEY) or an invalid option.
func NewClient(opts ...Option) (*Client, error) {
	cfg, err := resolveConfig(opts, os.Getenv)
	if err != nil {
		return nil, err
	}
	h := newHTTPCore(cfg)
	return &Client{
		Instagram: &InstagramService{h}, TikTok: &TikTokService{h}, YouTube: &YouTubeService{h},
		AppStore: &AppStoreService{h}, GitHub: &GitHubService{h}, HackerNews: &HackerNewsService{h},
		Bluesky: &BlueskyService{h}, Twitch: &TwitchService{h}, Linktree: &LinktreeService{h},
	}, nil
}
