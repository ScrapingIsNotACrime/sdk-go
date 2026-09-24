package sinac

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestNewClientMissingKey(t *testing.T) {
	t.Setenv(APIKeyEnv, "")
	if _, err := NewClient(); err == nil {
		t.Fatal("expected an error")
	}
}

func TestNewClientReadsEnv(t *testing.T) {
	t.Setenv(APIKeyEnv, "sinac_env")
	var key string
	client, err := NewClient(WithHTTPClient(&http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
		key = r.Header.Get("X-Api-Key")
		return response(200, `{"data":{"username":"x"}}`, nil), nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Linktree.Profile(context.Background(), "x"); err != nil || key != "sinac_env" {
		t.Fatalf("key=%q err=%v", key, err)
	}
}

func TestInvalidArgumentMakesNoRequest(t *testing.T) {
	client, _ := NewClient(WithAPIKey("k"), WithHTTPClient(&http.Client{Transport: roundTrip(func(*http.Request) (*http.Response, error) {
		t.Fatal("no request expected")
		return nil, nil
	})}))
	if _, err := client.Instagram.Profile(context.Background(), ""); !errors.Is(err, ErrBadRequest) {
		t.Fatalf("Profile: %v", err)
	}
	if _, err := client.GitHub.Followers(context.Background(), "..", nil); !errors.Is(err, ErrBadRequest) {
		t.Fatalf("Followers: %v", err)
	}
}

func TestUserClientIsNotModified(t *testing.T) {
	user := &http.Client{}
	if _, err := NewClient(WithAPIKey("k"), WithHTTPClient(user)); err != nil {
		t.Fatal(err)
	}
	if user.CheckRedirect != nil {
		t.Fatal("NewClient modified the caller's http.Client")
	}
}
