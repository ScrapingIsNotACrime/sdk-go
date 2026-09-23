package sinac

import "context"

// AppStoreService groups the App Store endpoints.
type AppStoreService struct{ http *httpCore }

// Search is GET /appstore/search — Country defaults to "us", Limit 1-200 (default 10).
func (s *AppStoreService) Search(ctx context.Context, term string, p *AppStoreSearchParams) (*AppStoreSearch, error) {
	r, err := appstoreSearch(term, p)
	return getJSON[AppStoreSearch](ctx, s.http, r, err)
}

// Reviews is GET /appstore/reviews — pages 1-10 (Apple's cap); the API returns 400 past page 10.
func (s *AppStoreService) Reviews(ctx context.Context, appID string, p *AppStoreReviewsParams) (*Page[AppStoreReview, AppStoreReviewPage], error) {
	spec, err := appstoreReviews(appID, p)
	return fetchPage[AppStoreReview, AppStoreReviewPage](ctx, s.http, spec, err)
}
