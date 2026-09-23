package sinac

// InstagramProfile is GET /instagram/profile/{username}.
type InstagramProfile struct {
	ID                 string   `json:"id"`
	FBID               string   `json:"fbid"`
	Username           string   `json:"username"`
	FullName           string   `json:"full_name"`
	Bio                string   `json:"bio"`
	BioLinks           []string `json:"bio_links"`
	Followers          int64    `json:"followers"`
	Following          int64    `json:"following"`
	Medias             int64    `json:"medias"`
	HighlightReelCount int64    `json:"highlight_reel_count"`
	ProfilePic         string   `json:"profile_pic"`
	HasAREffects       bool     `json:"has_ar_effects"`
	HasClips           bool     `json:"has_clips"`
	HasGuides          bool     `json:"has_guides"`
	HasChannel         bool     `json:"has_channel"`
	HasBlockedViewer   bool     `json:"has_blocked_viewer"`
	IsBusinessAccount  bool     `json:"is_business_account"`
	// Null in every observed example; no evidence of its populated shape.
	BusinessAddressJSON   any     `json:"business_address_json,omitempty"`
	BusinessContactMethod string  `json:"business_contact_method"`
	BusinessEmail         *string `json:"business_email,omitempty"`
	BusinessPhoneNumber   *string `json:"business_phone_number,omitempty"`
	BusinessCategoryName  string  `json:"business_category_name"`
	IsProfessionalAccount bool    `json:"is_professional_account"`
	CategoryName          string  `json:"category_name"`
	IsPrivate             bool    `json:"is_private"`
	IsVerified            bool    `json:"is_verified"`
}

// InstagramContactAddress is a profile's structured business address, from the contact endpoint.
type InstagramContactAddress struct {
	StreetAddress string `json:"street_address"`
	ZipCode       string `json:"zip_code"`
	CityName      string `json:"city_name"`
	RegionName    string `json:"region_name"`
	CountryCode   string `json:"country_code"`
}

// InstagramFoundContact is an email address or phone number found written into a profile's bio.
type InstagramFoundContact struct {
	Value  string `json:"value"`
	Source string `json:"source"`
}

// InstagramContact is GET /instagram/profile/{username}/contact.
type InstagramContact struct {
	Username    string  `json:"username"`
	FullName    string  `json:"full_name"`
	Biography   string  `json:"biography"`
	IsVerified  bool    `json:"is_verified"`
	IsBusiness  bool    `json:"is_business"`
	Category    string  `json:"category"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	ExternalURL string  `json:"external_url"`
	// Null for a profile with no contact information (per the endpoint's docs, all fields in that case are null).
	Address     *InstagramContactAddress `json:"address,omitempty"`
	EmailsFound []InstagramFoundContact  `json:"emails_found"`
	PhonesFound []InstagramFoundContact  `json:"phones_found"`
}

// InstagramClipsMusicAttribution is music attribution for a video/reel; null for original audio.
type InstagramClipsMusicAttribution struct {
	ArtistName        string `json:"artist_name"`
	SongName          string `json:"song_name"`
	UsesOriginalAudio bool   `json:"uses_original_audio"`
}

// InstagramLatestPostMedia is a post in a profile's latest-posts timeline — an image post or
// a video post; Type tells them apart. Video-only fields (VideoViews, VideoURL, HasAudio,
// ClipsMusicAttributionInfo) are nil for image posts.
type InstagramLatestPostMedia struct {
	ID        string `json:"id"`
	Shortcode string `json:"shortcode"`
	Type      string `json:"type"`
	// Video-only field.
	VideoViews *int64 `json:"video_views,omitempty"`
	Comments   int64  `json:"comments"`
	Likes      int64  `json:"likes"`
	// Null when the post has no caption.
	Caption            *string `json:"caption,omitempty"`
	Location           any     `json:"location,omitempty"`
	ThumbnailResources any     `json:"thumbnail_resources,omitempty"`
	DisplayURL         string  `json:"display_url"`
	// Video-only field.
	VideoURL *string `json:"video_url,omitempty"`
	// Video-only field.
	HasAudio *bool `json:"has_audio,omitempty"`
	// Video-only field.
	ClipsMusicAttributionInfo *InstagramClipsMusicAttribution `json:"clips_music_attribution_info,omitempty"`
	TakenAtTimestamp          string                          `json:"taken_at_timestamp"`
}

// InstagramLatestPosts is GET /instagram/profile/{username}/timeline/latest.
type InstagramLatestPosts struct {
	Count       int64                      `json:"count"`
	LatestCount int64                      `json:"latest_count"`
	Medias      []InstagramLatestPostMedia `json:"medias"`
}

// InstagramMedia is a post, as returned by the paged timeline and the highlight-content endpoints.
type InstagramMedia struct {
	ID        string `json:"id"`
	Shortcode string `json:"shortcode"`
	Type      string `json:"type"`
	// Null when the post has no caption.
	Caption          *string `json:"caption,omitempty"`
	Likes            int64   `json:"likes"`
	Comments         int64   `json:"comments"`
	PreviewComments  []any   `json:"preview_comments"`
	Location         any     `json:"location,omitempty"`
	DisplayURL       string  `json:"display_url"`
	TakenAtTimestamp string  `json:"taken_at_timestamp"`
}

// InstagramTimelinePage is GET /instagram/profile/{username}/timeline.
type InstagramTimelinePage struct {
	Medias     []InstagramMedia `json:"medias"`
	HasMore    bool             `json:"has_more"`
	NextCursor *string          `json:"next_cursor,omitempty"`
}

// InstagramHighlightSummary is a highlight reel's summary, as listed on a profile.
type InstagramHighlightSummary struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Cover string `json:"cover"`
}

// InstagramHighlights is GET /instagram/profile/{username}/highlights.
type InstagramHighlights struct {
	Username   string                      `json:"username"`
	UserID     string                      `json:"user_id"`
	Highlights []InstagramHighlightSummary `json:"highlights"`
}

// InstagramHighlight is GET /instagram/highlights/{highlightId}.
type InstagramHighlight struct {
	ID    string           `json:"id"`
	Title string           `json:"title"`
	Items []InstagramMedia `json:"items"`
}

// InstagramMediaDetail is GET /instagram/profile/{username}/media/{mediaId} — the same media
// object the timeline endpoints return, including the video-only fields when the media is a video.
type InstagramMediaDetail struct {
	ID        string `json:"id"`
	Shortcode string `json:"shortcode"`
	Type      string `json:"type"`
	Comments  int64  `json:"comments"`
	Likes     int64  `json:"likes"`
	// Null when the post has no caption.
	Caption                   *string                         `json:"caption,omitempty"`
	Location                  any                             `json:"location,omitempty"`
	ThumbnailResources        any                             `json:"thumbnail_resources,omitempty"`
	DisplayURL                string                          `json:"display_url"`
	TakenAtTimestamp          string                          `json:"taken_at_timestamp"`
	VideoViews                *int64                          `json:"video_views,omitempty"`
	VideoURL                  *string                         `json:"video_url,omitempty"`
	HasAudio                  *bool                           `json:"has_audio,omitempty"`
	ClipsMusicAttributionInfo *InstagramClipsMusicAttribution `json:"clips_music_attribution_info,omitempty"`
}

// InstagramDownloadAsset is one downloadable asset behind a post, reel, or carousel.
type InstagramDownloadAsset struct {
	Kind   string `json:"kind"`
	Index  int64  `json:"index"`
	URL    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
	// Null for non-video assets (e.g. thumbnails).
	Quality   *string `json:"quality,omitempty"`
	ExpiresAt string  `json:"expires_at"`
}

// InstagramDownload is GET /instagram/media/{shortcode}/download — assets[0] is always the
// best primary asset.
type InstagramDownload struct {
	Shortcode string                   `json:"shortcode"`
	Type      string                   `json:"type"`
	ExpiresAt string                   `json:"expires_at"`
	Assets    []InstagramDownloadAsset `json:"assets"`
}

// InstagramShortcodeID is GET /instagram/media/{shortcode}/id and /instagram/media/id/{mediaId}.
type InstagramShortcodeID struct {
	Shortcode string `json:"shortcode"`
	MediaID   string `json:"media_id"`
}

// InstagramReel is GET /instagram/reels/{shortcode} — clips_music_attribution_info is null for
// original audio.
type InstagramReel struct {
	ID         string `json:"id"`
	Shortcode  string `json:"shortcode"`
	Type       string `json:"type"`
	VideoViews int64  `json:"video_views"`
	Comments   int64  `json:"comments"`
	Likes      int64  `json:"likes"`
	// Null when the post has no caption.
	Caption                   *string                         `json:"caption,omitempty"`
	Location                  any                             `json:"location,omitempty"`
	ThumbnailResources        any                             `json:"thumbnail_resources,omitempty"`
	DisplayURL                string                          `json:"display_url"`
	VideoURL                  string                          `json:"video_url"`
	HasAudio                  bool                            `json:"has_audio"`
	ClipsMusicAttributionInfo *InstagramClipsMusicAttribution `json:"clips_music_attribution_info,omitempty"`
	TakenAtTimestamp          string                          `json:"taken_at_timestamp"`
}
