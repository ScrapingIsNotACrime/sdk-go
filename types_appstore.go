package sinac

// AppStoreApp is one app in an App Store search result.
type AppStoreApp struct {
	ID          int64    `json:"id"`
	BundleID    string   `json:"bundleId"`
	Name        string   `json:"name"`
	Developer   string   `json:"developer"`
	URL         string   `json:"url"`
	IconURL     string   `json:"iconUrl"`
	Price       float64  `json:"price"`
	Currency    string   `json:"currency"`
	Rating      float64  `json:"rating"`
	RatingCount int64    `json:"ratingCount"`
	Version     string   `json:"version"`
	Genres      []string `json:"genres"`
	Screenshots []string `json:"screenshots"`
}

// AppStoreSearch is GET /appstore/search.
type AppStoreSearch struct {
	Term        string        `json:"term"`
	Country     string        `json:"country"`
	ResultCount int64         `json:"resultCount"`
	Apps        []AppStoreApp `json:"apps"`
}

// AppStoreReview is one customer review of an app.
type AppStoreReview struct {
	ID        string  `json:"id"`
	Author    string  `json:"author"`
	Rating    float64 `json:"rating"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	Version   string  `json:"version"`
	UpdatedAt string  `json:"updatedAt"`
	VoteCount int64   `json:"voteCount"`
	VoteSum   int64   `json:"voteSum"`
}

// AppStoreReviewPage is GET /appstore/reviews.
type AppStoreReviewPage struct {
	AppID   string           `json:"appId"`
	Country string           `json:"country"`
	Page    int64            `json:"page"`
	Reviews []AppStoreReview `json:"reviews"`
}
