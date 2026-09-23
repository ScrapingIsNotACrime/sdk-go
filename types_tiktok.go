package sinac

import "encoding/json"

// TikTokProfile is GET /tiktok/profile/{username}.
type TikTokProfile struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Bio        string `json:"bio"`
	BioLink    string `json:"bio_link"`
	Avatar     string `json:"avatar"`
	SecUID     string `json:"sec_uid"`
	Followers  int64  `json:"followers"`
	Following  int64  `json:"following"`
	Hearts     int64  `json:"hearts"`
	Videos     int64  `json:"videos"`
	IsPrivate  bool   `json:"is_private"`
	IsVerified bool   `json:"is_verified"`
}

// TikTokVideo is GET /tiktok/video/{videoId}. No documented example yet; Extra holds every other key.
type TikTokVideo struct {
	ID    *string        `json:"id,omitempty"`
	Extra map[string]any `json:"-"`
}

func (v *TikTokVideo) UnmarshalJSON(data []byte) error {
	type known TikTokVideo
	if err := json.Unmarshal(data, (*known)(v)); err != nil {
		return err
	}
	var all map[string]any
	if err := json.Unmarshal(data, &all); err != nil {
		return err
	}
	delete(all, "id")
	v.Extra = all
	return nil
}
