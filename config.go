package sinac

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the production API.
	DefaultBaseURL = "https://api.scrapingisnotacrime.com/v1"
	// APIKeyEnv is read when WithAPIKey is not given.
	APIKeyEnv = "SCRAPINGISNOTACRIME_API_KEY"
	// DefaultTimeout bounds each attempt: connect, headers and body.
	DefaultTimeout = 30 * time.Second
	// DefaultMaxRetries is the number of extra attempts for 429, 502 and network errors.
	DefaultMaxRetries = 2
)

// Option configures a Client.
type Option func(*config)

// WithAPIKey sets the API key (sinac_…). Defaults to the SCRAPINGISNOTACRIME_API_KEY environment variable.
func WithAPIKey(key string) Option { return func(c *config) { c.apiKey = key } }

// WithBaseURL overrides the API base URL. It must be the final HTTPS URL: redirects are not followed.
func WithBaseURL(baseURL string) Option { return func(c *config) { c.baseURL = baseURL } }

// WithTimeout sets the per-attempt timeout (connect, headers and body). Defaults to 30 s.
func WithTimeout(timeout time.Duration) Option { return func(c *config) { c.timeout = timeout } }

// WithMaxRetries sets the extra attempts for 429, 502 and network errors. Defaults to 2; 0 disables retries.
func WithMaxRetries(n int) Option { return func(c *config) { c.maxRetries = n } }

// WithHTTPClient sets the *http.Client used for requests (proxies, tests). The client is never modified.
func WithHTTPClient(client *http.Client) Option { return func(c *config) { c.httpClient = client } }

type config struct {
	apiKey     string
	baseURL    string
	timeout    time.Duration
	maxRetries int
	httpClient *http.Client
}

func resolveConfig(opts []Option, getenv func(string) string) (config, error) {
	cfg := config{baseURL: DefaultBaseURL, timeout: DefaultTimeout, maxRetries: DefaultMaxRetries}
	for _, opt := range opts {
		opt(&cfg)
	}
	if strings.TrimSpace(cfg.apiKey) == "" {
		cfg.apiKey = getenv(APIKeyEnv)
	}
	cfg.apiKey = strings.TrimSpace(cfg.apiKey)
	if cfg.apiKey == "" {
		return cfg, fmt.Errorf("sinac: missing API key: pass sinac.WithAPIKey or set the %s environment variable", APIKeyEnv)
	}
	if cfg.timeout <= 0 {
		return cfg, errors.New("sinac: timeout must be positive")
	}
	if cfg.maxRetries < 0 {
		return cfg, errors.New("sinac: max retries must not be negative")
	}
	cfg.baseURL = strings.TrimRight(cfg.baseURL, "/")
	parsed, err := url.Parse(cfg.baseURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return cfg, fmt.Errorf("sinac: invalid base URL %q", cfg.baseURL)
	}
	return cfg, nil
}
