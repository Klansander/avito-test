package swagger

import "time"

type Banner struct {
	Content string `json:"content" format:"json" example:"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}"`
}
type ListBanner struct {
	Content   string    `json:"content" format:"json" example:"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}"`
	IsActive  bool      `json:"is_active,omitempty"`
	BannerID  int       `json:"banner_id,omitempty"`
	TagID     []int     `json:"tag_id,omitempty"`
	FeatureID int       `json:"feature_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAT time.Time `json:"updated_at,omitempty"`
}

type CreateBanner struct {
	BannerID int `json:"banner_id,omitempty"`
}

type NewBanner struct {
	Content   string `json:"content" format:"json" example:"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}"`
	IsActive  *bool  `json:"is_active" binding:"required"`
	TagID     *[]int `json:"tag_id" binding:"required"`
	FeatureID *int   `json:"feature_id" binding:"required"`
}

type HeaderBanner struct {
	Content   string `json:"content" format:"json" example:"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}"`
	IsActive  *bool  `json:"is_active,omitempty"`
	TagID     *[]int `json:"tag_id,omitempty"`
	FeatureID *int   `json:"feature_id,omitempty"`
}

type Error struct {
	Error string `json:"error"`
}
