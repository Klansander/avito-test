package service

import (
	"avito/internal/model"
	"avito/internal/repository"
	"context"
)

type Banner interface {
	UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet) (data map[string]interface{}, err error)
	ListBanner(c context.Context, userBannerQuery model.UserBannerQueryList) (data []model.Banner, err error)
	CreateBanner(c context.Context, headerBanner model.HeaderBanner) (id int, err error)
	UpdateBanner(c context.Context, bannerID int, headerBanner model.HeaderBanner) error
	DeleteBanner(c context.Context, bannerID int) (string, error)
}

type BannerService struct {
	r *repository.Repository
}

func NewBanner(r *repository.Repository) *BannerService {

	return &BannerService{r: r}

}

func (s *BannerService) UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet) (data map[string]interface{}, err error) {

	data, err = s.r.Banner.UserBanner(c, userBannerQuery)
	if err != nil {
		return nil, err
	}

	return
}

func (s *BannerService) ListBanner(c context.Context, userBannerQuery model.UserBannerQueryList) (data []model.Banner, err error) {

	data, err = s.r.Banner.ListBanner(c, userBannerQuery)
	if err != nil {
		return nil, err
	}

	return
}
func (s *BannerService) CreateBanner(c context.Context, headerBanner model.HeaderBanner) (id int, err error) {

	id, err = s.r.Banner.CreateBanner(c, headerBanner)
	if err != nil {
		return 0, err
	}

	return
}

func (s *BannerService) UpdateBanner(c context.Context, bannerID int, headerBanner model.HeaderBanner) error {
	return s.r.Banner.UpdateBanner(c, bannerID, headerBanner)
}
func (s *BannerService) DeleteBanner(c context.Context, bannerID int) (string, error) {

	return s.r.Banner.DeleteBanner(c, bannerID)
}
