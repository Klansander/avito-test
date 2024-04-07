package repository

import (
	"avito/internal/model"
	"avito/pkg/postgresql"
	"context"
)

type Banner interface {
	UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet) error
	ListBanner(c context.Context, userBannerQuery model.UserBannerQueryList) error
	CreateBanner(c context.Context, headerBanner model.HeaderBanner) error
	UpdateBanner(c context.Context, bannerID int, headerBanner model.HeaderBanner) error
	DeleteBanner(c context.Context, bannerID int) error
}

type BannerRepository struct {
	db *postgresql.Postgres
}

func NewBanner(db *postgresql.Postgres) *BannerRepository {

	return &BannerRepository{db: db}

}

func (r *BannerRepository) UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet) error {
	return nil
}

func (r *BannerRepository) ListBanner(c context.Context, userBannerQuery model.UserBannerQueryList) error {
	return nil
}

func (r *BannerRepository) CreateBanner(c context.Context, headerBanner model.HeaderBanner) error {
	return nil
}
func (r *BannerRepository) UpdateBanner(c context.Context, bannerID int, headerBanner model.HeaderBanner) error {
	return nil
}
func (r *BannerRepository) DeleteBanner(c context.Context, bannerID int) error {
	return nil
}
