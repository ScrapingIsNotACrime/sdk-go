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

// UnmarshalJSON decodes id into ID only when it is a JSON string; any other
// shape (number, null, ...) leaves ID nil and keeps the raw value in Extra,
// alongside every other key, so no data is ever dropped.
func (v *TikTokVideo) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	v.ID = nil
	extra := make(map[string]any, len(fields))
	for key, raw := range fields {
		var value any
		if err := json.Unmarshal(raw, &value); err != nil {
			return err
		}
		if key == "id" {
			if id, ok := value.(string); ok {
				v.ID = &id
				continue
			}
		}
		extra[key] = value
	}
	v.Extra = extra
	return nil
}

// MarshalJSON re-encodes Extra plus id (when ID is set), so round-tripping
// through the SDK never loses data.
func (v TikTokVideo) MarshalJSON() ([]byte, error) {
	out := make(map[string]any, len(v.Extra)+1)
	for key, value := range v.Extra {
		out[key] = value
	}
	if v.ID != nil {
		out["id"] = *v.ID
	}
	return json.Marshal(out)
}
