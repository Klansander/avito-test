package service

import (
	"avito/app/internal/model"
	"avito/app/internal/repository"
	pc "avito/app/pkg/context"
	"context"
	json "github.com/json-iterator/go"
	"sync"
)

type Banner interface {
	UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet, isAdmin bool) (data map[string]interface{}, err error)
	ListBanner(c context.Context, userBannerQuery model.UserBannerQueryList) (data []model.Banner, err error)
	CreateBanner(c context.Context, headerBanner model.NewBanner) (id int, err error)
	UpdateBanner(c context.Context, bannerID int, headerBanner model.HeaderBanner) error
	DeleteBanner(c context.Context, bannerID int) (string, error)
	GetVersionBanner(c context.Context, headerBanner model.BannerVersion) (dataArr []model.Banner, err error)
	DeleteBannerByTagIdOrFeatureId(c context.Context, banner model.BannerTagOrFeatureID, wg *sync.WaitGroup) (err error)
}

type BannerService struct {
	r *repository.Repository
}

func NewBanner(r *repository.Repository) *BannerService {

	return &BannerService{r: r}

}

func (s *BannerService) UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet, isAdmin bool) (data map[string]interface{}, err error) {

	if !userBannerQuery.UseLastRevision {
		data, err = s.r.Banner.UserBannerRedis(c, userBannerQuery, isAdmin)
		if err != nil {
			return nil, err
		}
		return
	}
	data, err = s.r.Banner.UserBanner(c, userBannerQuery, isAdmin)
	if err != nil {
		return nil, err
	}

	return
}

func (s *BannerService) ListBanner(c context.Context, userBannerQuery model.UserBannerQueryList) (data []model.Banner, err error) {

	if userBannerQuery.Limit == nil || *userBannerQuery.Limit <= 0 {
		cfg := pc.GetConfig(c)
		userBannerQuery.Limit = &cfg.PSQL.LimitMax
	}

	if userBannerQuery.Offset == nil || *userBannerQuery.Offset < 0 {
		offset := 0
		userBannerQuery.Offset = &offset
	}

	data, err = s.r.Banner.ListBanner(c, userBannerQuery)
	if err != nil {
		return nil, err
	}

	return
}
func (s *BannerService) CreateBanner(c context.Context, headerBanner model.NewBanner) (id int, err error) {

	id, err = s.r.Banner.CreateBanner(c, headerBanner)
	if err != nil {
		return 0, err
	}

	return
}

func (s *BannerService) UpdateBanner(c context.Context, bannerID int, headerBanner model.HeaderBanner) error {

	bannerByte, err := json.Marshal(headerBanner)
	if err != nil {
		return err
	}

	data := make(map[string]interface{}, 4)

	err = json.Unmarshal(bannerByte, &data)
	if err != nil {
		return err
	}

	return s.r.Banner.UpdateBanner(c, bannerID, data)
}
func (s *BannerService) DeleteBanner(c context.Context, bannerID int) (string, error) {

	return s.r.Banner.DeleteBanner(c, bannerID)
}

func (s *BannerService) GetVersionBanner(c context.Context, headerBanner model.BannerVersion) (dataArr []model.Banner, err error) {
	return s.r.Banner.GetVersionBanner(c, headerBanner)
}
func (s *BannerService) DeleteBannerByTagIdOrFeatureId(c context.Context, banner model.BannerTagOrFeatureID, wg *sync.WaitGroup) (err error) {
	return s.r.Banner.DeleteBannerByTagOrFeature(c, banner, wg)
}
