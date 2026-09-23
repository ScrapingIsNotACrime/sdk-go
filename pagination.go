package sinac

import (
	"context"
	"encoding/json"
	"iter"
	"strconv"
)

type pageKind int

const (
	cursorPages pageKind = iota
	numberedPages
)

type pageSpec struct {
	path     string
	query    []queryParam
	kind     pageKind
	itemsKey string
	cursor   string // cursorPages: "" is the first page
	page     int    // numberedPages: the page to request (always sent)
	// maxPage > 0 marks numbered endpoints without a has_more flag (App Store
	// reviews): more while the page is non-empty and below this cap.
	maxPage int
}

func (s pageSpec) route() route {
	query := append([]queryParam(nil), s.query...)
	if s.kind == cursorPages {
		query = append(query, qs("cursor", s.cursor)...)
	} else {
		query = append(query, queryParam{"page", strconv.Itoa(s.page)})
	}
	return route{path: s.path, query: query}
}

// Page is one page of results. Data is the full page object; Items, HasMore,
// NextCursor ("" = none) and NextPage (0 = none) are read from it.
type Page[T, R any] struct {
	Items      []T
	HasMore    bool
	NextCursor string
	NextPage   int
	Data       R

	http *httpCore
	next *pageSpec
}

func fetchPage[T, R any](ctx context.Context, h *httpCore, spec pageSpec, specErr error) (*Page[T, R], error) {
	if specErr != nil {
		return nil, specErr
	}
	raw, err := h.get(ctx, spec.route())
	if err != nil {
		return nil, err
	}
	page := &Page[T, R]{http: h}
	if err := decode(raw, &page.Data); err != nil {
		return nil, err
	}
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(raw, &fields) // nil when data is not an object
	if rawItems, ok := fields[spec.itemsKey]; ok {
		if err := decode(rawItems, &page.Items); err != nil {
			return nil, err
		}
	}
	var hasMore bool
	_ = json.Unmarshal(fields["has_more"], &hasMore)

	switch spec.kind {
	case cursorPages:
		_ = json.Unmarshal(fields["next_cursor"], &page.NextCursor)
		page.HasMore = hasMore && page.NextCursor != ""
	case numberedPages:
		if spec.maxPage > 0 {
			page.HasMore = len(page.Items) > 0 && spec.page < spec.maxPage
		} else {
			page.HasMore = hasMore
		}
		if page.HasMore {
			page.NextPage = spec.page + 1
		}
	}
	if page.HasMore && len(page.Items) > 0 {
		next := spec
		next.cursor, next.page = page.NextCursor, page.NextPage
		page.next = &next
	}
	return page, nil
}

// Next fetches the following page. It returns nil, nil when there is none.
func (p *Page[T, R]) Next(ctx context.Context) (*Page[T, R], error) {
	if p.next == nil {
		return nil, nil
	}
	return fetchPage[T, R](ctx, p.http, *p.next, nil)
}

// All yields every item from this page onward, fetching later pages lazily
// (each page is one billed request). Breaking out of the loop stops fetching.
// A fetch error is yielded once, with the zero item, and ends the iteration.
func (p *Page[T, R]) All(ctx context.Context) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for page := p; page != nil; {
			for _, item := range page.Items {
				if !yield(item, nil) {
					return
				}
			}
			next, err := page.Next(ctx)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			page = next
		}
	}
}
