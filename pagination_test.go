package sinac

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
)

type item struct {
	ID int `json:"id"`
}

type itemPage struct {
	Items      []item `json:"items"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
}

// pages serves bodies in order and records every request URL.
func pages(t *testing.T, bodies ...string) (*httpCore, *[]string) {
	t.Helper()
	var urls []string
	var n atomic.Int32
	core := testCore(t, roundTrip(func(r *http.Request) (*http.Response, error) {
		urls = append(urls, r.URL.RequestURI())
		i := int(n.Add(1)) - 1
		if i >= len(bodies) {
			t.Fatalf("unexpected request %d: %s", i, r.URL)
		}
		return response(200, bodies[i], nil), nil
	}))
	return core, &urls
}

func TestCursorPagination(t *testing.T) {
	core, urls := pages(t,
		`{"data":{"items":[{"id":1},{"id":2}],"has_more":true,"next_cursor":"c2"}}`,
		`{"data":{"items":[{"id":3}],"has_more":false,"next_cursor":""}}`,
	)
	spec := pageSpec{path: "/p", query: qi("limit", 2), kind: cursorPages, itemsKey: "items"}
	page, err := fetchPage[item, itemPage](context.Background(), core, spec, nil)
	if err != nil || len(page.Items) != 2 || !page.HasMore || page.NextCursor != "c2" || page.Data.NextCursor != "c2" {
		t.Fatalf("page=%+v err=%v", page, err)
	}
	var ids []int
	for it, err := range page.All(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, it.ID)
	}
	if len(ids) != 3 || ids[2] != 3 {
		t.Fatalf("ids = %v", ids)
	}
	if (*urls)[0] != "/v1/p?limit=2" || (*urls)[1] != "/v1/p?limit=2&cursor=c2" {
		t.Fatalf("urls = %v", *urls)
	}
}

func TestEmptyCursorEndsSequence(t *testing.T) {
	core, _ := pages(t, `{"data":{"items":[{"id":1}],"has_more":true,"next_cursor":""}}`)
	spec := pageSpec{path: "/p", kind: cursorPages, itemsKey: "items"}
	page, err := fetchPage[item, itemPage](context.Background(), core, spec, nil)
	if err != nil || page.HasMore {
		t.Fatalf("page=%+v err=%v", page, err)
	}
	next, err := page.Next(context.Background())
	if next != nil || err != nil {
		t.Fatalf("next=%v err=%v", next, err)
	}
}

func TestNumberedPaginationUsesHasMore(t *testing.T) {
	core, urls := pages(t,
		`{"data":{"items":[{"id":1}],"has_more":true}}`,
		`{"data":{"items":[{"id":2}],"has_more":false}}`,
	)
	spec := pageSpec{path: "/n", query: qi("limit", 1), kind: numberedPages, itemsKey: "items", page: 0}
	page, err := fetchPage[item, itemPage](context.Background(), core, spec, nil)
	if err != nil || page.NextPage != 1 || !page.HasMore {
		t.Fatalf("page=%+v err=%v", page, err)
	}
	next, err := page.Next(context.Background())
	if err != nil || next.HasMore || next.NextPage != 0 || next.Items[0].ID != 2 {
		t.Fatalf("next=%+v err=%v", next, err)
	}
	if (*urls)[0] != "/v1/n?limit=1&page=0" || (*urls)[1] != "/v1/n?limit=1&page=1" {
		t.Fatalf("urls = %v", *urls)
	}
}

func TestNumberedMaxPageWithoutHasMore(t *testing.T) {
	core, _ := pages(t, `{"data":{"reviews":[{"id":1}]}}`)
	spec := pageSpec{path: "/r", kind: numberedPages, itemsKey: "reviews", page: 10, maxPage: 10}
	page, err := fetchPage[item, map[string]any](context.Background(), core, spec, nil)
	if err != nil || page.HasMore || len(page.Items) != 1 {
		t.Fatalf("page=%+v err=%v", page, err)
	}
	core, _ = pages(t, `{"data":{"reviews":[{"id":1}]}}`)
	spec.page = 3
	page, _ = fetchPage[item, map[string]any](context.Background(), core, spec, nil)
	if !page.HasMore || page.NextPage != 4 {
		t.Fatalf("page 3 of 10 = %+v", page)
	}
}

func TestEmptyItemsStopsEvenWithHasMore(t *testing.T) {
	core, _ := pages(t, `{"data":{"items":[],"has_more":true}}`)
	spec := pageSpec{path: "/n", kind: numberedPages, itemsKey: "items", page: 1}
	page, _ := fetchPage[item, itemPage](context.Background(), core, spec, nil)
	if next, err := page.Next(context.Background()); next != nil || err != nil {
		t.Fatalf("next=%v err=%v", next, err)
	}
}

func TestAllBreakStopsFetching(t *testing.T) {
	core, urls := pages(t, `{"data":{"items":[{"id":1},{"id":2}],"has_more":true,"next_cursor":"c2"}}`)
	spec := pageSpec{path: "/p", kind: cursorPages, itemsKey: "items"}
	page, _ := fetchPage[item, itemPage](context.Background(), core, spec, nil)
	for range page.All(context.Background()) {
		break
	}
	if len(*urls) != 1 {
		t.Fatalf("requests after break: %v", *urls)
	}
}

func TestAllYieldsFetchErrorOnce(t *testing.T) {
	var n atomic.Int32
	core := testCore(t, roundTrip(func(*http.Request) (*http.Response, error) {
		if n.Add(1) == 1 {
			return response(200, `{"data":{"items":[{"id":1}],"has_more":true,"next_cursor":"c2"}}`, nil), nil
		}
		return response(404, `{"message":"gone"}`, nil), nil
	}))
	spec := pageSpec{path: "/p", kind: cursorPages, itemsKey: "items"}
	page, _ := fetchPage[item, itemPage](context.Background(), core, spec, nil)
	var got []error
	for _, err := range page.All(context.Background()) {
		got = append(got, err)
	}
	if len(got) != 2 || got[0] != nil || !errors.Is(got[1], ErrNotFound) {
		t.Fatalf("yields = %v", got)
	}
}
