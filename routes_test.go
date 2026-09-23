package sinac

import (
	"errors"
	"testing"
)

func TestRoutesRejectInvalidSegments(t *testing.T) {
	if _, err := instagramProfile(""); !errors.Is(err, ErrBadRequest) {
		t.Fatalf("instagramProfile: %v", err)
	}
	if _, err := instagramMediaByID("nasa", ".."); !errors.Is(err, ErrBadRequest) {
		t.Fatalf("instagramMediaByID: %v", err)
	}
	if _, err := githubFollowers(".", nil); !errors.Is(err, ErrBadRequest) {
		t.Fatalf("githubFollowers: %v", err)
	}
	if _, err := hackernewsFeed("", nil); !errors.Is(err, ErrBadRequest) {
		t.Fatalf("hackernewsFeed: %v", err)
	}
}

func TestPageDefaults(t *testing.T) {
	gh, _ := githubFollowers("torvalds", nil)
	hn, _ := hackernewsSearch("go", nil)
	as, _ := appstoreReviews("1", nil)
	if gh.page != 1 || hn.page != 0 || as.page != 1 || as.maxPage != 10 {
		t.Fatalf("gh=%d hn=%d as=%d/%d", gh.page, hn.page, as.page, as.maxPage)
	}
}
