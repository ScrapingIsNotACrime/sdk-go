package sinac

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func env(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestResolveConfigDefaults(t *testing.T) {
	cfg, err := resolveConfig([]Option{WithAPIKey("  sinac_abc  ")}, env(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.apiKey != "sinac_abc" || cfg.baseURL != DefaultBaseURL || cfg.timeout != 30*time.Second || cfg.maxRetries != 2 {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestResolveConfigEnvKeyAndOverrides(t *testing.T) {
	client := &http.Client{}
	cfg, err := resolveConfig([]Option{
		WithBaseURL("http://localhost:8080/v1///"),
		WithTimeout(5 * time.Second),
		WithMaxRetries(0),
		WithHTTPClient(client),
	}, env(map[string]string{APIKeyEnv: "sinac_env"}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.apiKey != "sinac_env" || cfg.baseURL != "http://localhost:8080/v1" || cfg.timeout != 5*time.Second ||
		cfg.maxRetries != 0 || cfg.httpClient != client {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestResolveConfigRejects(t *testing.T) {
	cases := map[string][]Option{
		"missing key":      nil,
		"blank key":        {WithAPIKey("   ")},
		"zero timeout":     {WithAPIKey("k"), WithTimeout(0)},
		"negative retries": {WithAPIKey("k"), WithMaxRetries(-1)},
		"relative url":     {WithAPIKey("k"), WithBaseURL("api.example.com/v1")},
		"ftp url":          {WithAPIKey("k"), WithBaseURL("ftp://example.com")},
	}
	for name, opts := range cases {
		if _, err := resolveConfig(opts, env(nil)); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
	_, err := resolveConfig(nil, env(nil))
	if !strings.Contains(err.Error(), "WithAPIKey") || !strings.Contains(err.Error(), APIKeyEnv) {
		t.Fatalf("missing-key message must name WithAPIKey and %s: %v", APIKeyEnv, err)
	}
}
