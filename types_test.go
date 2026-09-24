package sinac

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type fixture struct {
	Request  string `json:"request"`
	Response struct {
		Data json.RawMessage `json:"data"`
	} `json:"response"`
}

func loadFixture(t *testing.T, id string) fixture {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "fixtures", id+".json"))
	if err != nil {
		t.Fatal(err)
	}
	var f fixture
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatal(err)
	}
	return f
}

var fixtureTypes = map[string]any{
	"ig-profile": InstagramProfile{}, "ig-contact": InstagramContact{}, "ig-timeline": InstagramLatestPosts{},
	"ig-timeline-paged": InstagramTimelinePage{}, "ig-highlights": InstagramHighlights{},
	"ig-highlight-content": InstagramHighlight{}, "ig-media-by-id": InstagramMediaDetail{},
	"ig-media-download": InstagramDownload{}, "ig-shortcode-to-id": InstagramShortcodeID{},
	"ig-id-to-shortcode": InstagramShortcodeID{}, "ig-reels": InstagramReel{},
	"tt-profile": TikTokProfile{}, "yt-channel-videos": YouTubeChannelVideos{},
	"as-search": AppStoreSearch{}, "as-reviews": AppStoreReviewPage{},
	"gh-profile": GitHubProfile{}, "gh-followers": GitHubUserPage{}, "gh-repos": GitHubRepositoryPage{},
	"gh-search-repos": GitHubRepositorySearchPage{}, "gh-trending": GitHubTrending{},
	"hn-feed": HackerNewsStoryPage{}, "hn-item": HackerNewsItem{}, "hn-search": HackerNewsStoryPage{},
	"hn-user": HackerNewsUser{}, "hn-user-submissions": HackerNewsStoryPage{},
	"bs-profile": BlueskyProfile{}, "bs-posts": BlueskyPostPage{},
	"tw-profile": TwitchProfile{}, "tw-videos": TwitchVideos{}, "lt-profile": LinktreeProfile{},
}

func TestFixtureTableCoversEveryFixture(t *testing.T) {
	files, _ := filepath.Glob(filepath.Join("testdata", "fixtures", "*.json"))
	if len(files) != 30 || len(fixtureTypes) != 30 {
		t.Fatalf("files=%d table=%d", len(files), len(fixtureTypes))
	}
	for _, f := range files {
		if _, ok := fixtureTypes[strings.TrimSuffix(filepath.Base(f), ".json")]; !ok {
			t.Errorf("no type for %s", f)
		}
	}
}

// TestTypesMatchFixtures decodes each documented example strictly: an unknown
// key or a mismatched value type fails, and every non-pointer field must be present.
func TestTypesMatchFixtures(t *testing.T) {
	for id, sample := range fixtureTypes {
		data := loadFixture(t, id).Response.Data
		target := reflect.New(reflect.TypeOf(sample))
		dec := json.NewDecoder(bytes.NewReader(data))
		if reflect.TypeOf(sample) != reflect.TypeOf(TikTokVideo{}) {
			dec.DisallowUnknownFields()
		}
		if err := dec.Decode(target.Interface()); err != nil {
			t.Errorf("%s: %v", id, err)
			continue
		}
		var generic any
		_ = json.Unmarshal(data, &generic)
		for _, problem := range missingRequired(reflect.TypeOf(sample), generic, id) {
			t.Error(problem)
		}
	}
}

// missingRequired lists non-pointer, non-omitempty fields absent from the JSON, recursively.
func missingRequired(typ reflect.Type, value any, at string) []string {
	switch typ.Kind() {
	case reflect.Pointer:
		if value == nil {
			return nil
		}
		return missingRequired(typ.Elem(), value, at)
	case reflect.Slice:
		items, _ := value.([]any)
		var out []string
		for i, item := range items {
			out = append(out, missingRequired(typ.Elem(), item, at+"["+strconv.Itoa(i)+"]")...)
		}
		return out
	case reflect.Struct:
		object, ok := value.(map[string]any)
		if !ok {
			return nil
		}
		var out []string
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			name, opts, _ := strings.Cut(field.Tag.Get("json"), ",")
			if name == "-" || !field.IsExported() {
				continue
			}
			child, present := object[name]
			required := field.Type.Kind() != reflect.Pointer && !strings.Contains(opts, "omitempty")
			if required && !present {
				out = append(out, at+"."+name+": required field missing from the documented example")
			}
			out = append(out, missingRequired(field.Type, child, at+"."+name)...)
		}
		return out
	}
	return nil
}

func TestTikTokVideoKeepsExtras(t *testing.T) {
	var video TikTokVideo
	if err := json.Unmarshal([]byte(`{"id":"7","desc":"hi","stats":{"plays":3}}`), &video); err != nil {
		t.Fatal(err)
	}
	if video.ID == nil || *video.ID != "7" || video.Extra["desc"] != "hi" || video.Extra["id"] != nil {
		t.Fatalf("video = %+v", video)
	}
}

// TestTikTokVideoUnexpectedIDTypeKeepsExtra covers an id of the wrong JSON
// type (a number, here): ID stays nil and the raw id value is kept in Extra
// like any other key, instead of being dropped.
func TestTikTokVideoUnexpectedIDTypeKeepsExtra(t *testing.T) {
	var video TikTokVideo
	if err := json.Unmarshal([]byte(`{"id":123,"desc":"hello"}`), &video); err != nil {
		t.Fatal(err)
	}
	if video.ID != nil {
		t.Fatalf("ID = %v, want nil", *video.ID)
	}
	if video.Extra["desc"] != "hello" {
		t.Fatalf("Extra[desc] = %v", video.Extra["desc"])
	}
	if video.Extra["id"] != float64(123) {
		t.Fatalf("Extra[id] = %v, want float64(123)", video.Extra["id"])
	}
}

func TestTikTokVideoMarshalUnmarshalRoundTrip(t *testing.T) {
	id := "7"
	want := TikTokVideo{ID: &id, Extra: map[string]any{"desc": "hi", "stats": map[string]any{"plays": float64(3)}}}
	data, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got TikTokVideo
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID == nil || *got.ID != "7" || got.Extra["desc"] != "hi" || !reflect.DeepEqual(got.Extra["stats"], want.Extra["stats"]) {
		t.Fatalf("got = %+v, want %+v", got, want)
	}
}
