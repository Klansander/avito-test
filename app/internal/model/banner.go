package model

import "time"

type UserBannerQueryGet struct {
	TagID           int `json:"tag_id" binding:"required"`
	FeatureID       int `json:"feature_id" binding:"required"`
	UseLastRevision int `json:"use_last_revision" binding:"omitempty"`
}
type UserBannerQueryList struct {
	TagID     int `json:"tag_id" `
	FeatureID int `json:"feature_id" `
	Limit     int `json:"limit"`
	Offset    int `json:"offset"`
}

type HeaderBanner struct {
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active,omitempty"`
	BannerID  int                    `json:"banner_id,omitempty"`
	TagID     int                    `json:"tag_id,omitempty"`
	FeatureID int                    `json:"feature_id,omitempty"`
}

type Banner struct {
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active,omitempty"`
	BannerID  int                    `json:"banner_id,omitempty"`
	TagID     int                    `json:"tag_id,omitempty"`
	FeatureID int                    `json:"feature_id,omitempty"`
	CreatedAt time.Time              `json:"created_at,omitempty"`
	UpdatedAT time.Time              `json:"updated_at,omitempty"`
}
