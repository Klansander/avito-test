package service

import (
	"avito/internal/model"
	"avito/internal/repository"
	"context"
)

type Banner interface {
	UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet) error
	ListBanner(c context.Context, userBannerQuery model.UserBannerQueryList) error
	CreateBanner(c context.Context, headerBanner model.HeaderBanner) error
	UpdateBanner(c context.Context, bannerID int, headerBanner model.HeaderBanner) error
	DeleteBanner(c context.Context, bannerID int) error
}

type BannerService struct {
	r *repository.Repository
}

func NewBanner(r *repository.Repository) *BannerService {

	return &BannerService{r: r}

}

func (s *BannerService) UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet) error {
	s.r.Banner.UserBanner(c, userBannerQuery)
	return nil
}

func (s *BannerService) ListBanner(c context.Context, userBannerQuery model.UserBannerQueryList) error {

	s.r.Banner.ListBanner(c, userBannerQuery)

	return nil
}
func (s *BannerService) CreateBanner(c context.Context, headerBanner model.HeaderBanner) error {
	s.r.Banner.CreateBanner(c, headerBanner)
	return nil
}

func (s *BannerService) UpdateBanner(c context.Context, bannerID int, headerBanner model.HeaderBanner) error {
	return s.r.Banner.UpdateBanner(c, bannerID, headerBanner)
}
func (s *BannerService) DeleteBanner(c context.Context, bannerID int) error {
	s.r.Banner.DeleteBanner(c, bannerID)

	return nil
}
