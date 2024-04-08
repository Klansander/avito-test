package repository

import (
	"avito/internal/model"
	pc "avito/pkg/context"
	"avito/pkg/errors"
	"avito/pkg/postgresql"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
)

type Banner interface {
	UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet) (data map[string]interface{}, err error)
	ListBanner(c context.Context, userBannerQuery model.UserBannerQueryList) (data []model.Banner, err error)
	CreateBanner(c context.Context, headerBanner model.HeaderBanner) (id int, err error)
	UpdateBanner(c context.Context, bannerID int, headerBanner model.HeaderBanner) error
	DeleteBanner(c context.Context, bannerID int) (mes string, err error)
}

type BannerRepository struct {
	db *postgresql.Postgres
}

func NewBanner(db *postgresql.Postgres) *BannerRepository {

	return &BannerRepository{db: db}

}

func (r *BannerRepository) UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet) (data map[string]interface{}, err error) {

	cfg := pc.GetConfig(c)

	ctx, cancel := context.WithTimeout(c, cfg.PSQL.Timeout)
	defer cancel()

	var (
		res int
		mes string
		row sql.NullString
	)

	query := "select o_json,o_res,o_mes from public.fn_banner_get($1, $2)"

	// i_tag_id int,
	//i_feature_id int,
	//out o_json json,
	//	out o_res int,
	//	out o_mes text

	if err = r.db.DB.QueryRow(ctx, query, userBannerQuery.TagID, userBannerQuery.FeatureID).Scan(&row, &res, &mes); err != nil {
		err = errors.Wrap(err)

		return
	}
	if res != 0 {
		err = errors.New(http.StatusNotFound, mes)
		return
	}

	if row.Valid {
		if err = json.Unmarshal([]byte(row.String), &data); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return

}

func (r *BannerRepository) ListBanner(c context.Context, userBannerQuery model.UserBannerQueryList) (data []model.Banner, err error) {

	cfg := pc.GetConfig(c)

	ctx, cancel := context.WithTimeout(c, cfg.PSQL.Timeout)
	defer cancel()

	var row sql.NullString

	query := "select o_json from public.fn_banner_list($1, $2, $3, $4)"

	//i_tag_id int,
	//i_feature_id int,
	//i_limit int,
	//i_offset int,
	//	out o_json json

	if err = r.db.DB.QueryRow(ctx, query, userBannerQuery.TagID, userBannerQuery.FeatureID, userBannerQuery.Limit, userBannerQuery.Offset).Scan(&row); err != nil {
		err = errors.Wrap(err)
		return
	}

	if row.Valid {
		if err = json.Unmarshal([]byte(row.String), &data); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		data = make([]model.Banner, 0)
	}

	return

}

func (r *BannerRepository) CreateBanner(c context.Context, headerBanner model.HeaderBanner) (id int, err error) {

	cfg := pc.GetConfig(c)

	ctx, cancel := context.WithTimeout(c, cfg.PSQL.Timeout)
	defer cancel()

	var row sql.NullString

	query := "select v_id from public.fn_banner_ins($1, $2, $3, $4)"

	//(i_is_active boolean,
	//	i_tag_id int[],
	//	i_feature_id int,
	//	i_content json)
	//returns int

	if err = r.db.DB.QueryRow(ctx, query, headerBanner.IsActive, headerBanner.TagID, headerBanner.FeatureID).Scan(&row, &id); err != nil {
		err = errors.Wrap(err)
		return
	}

	return

}
func (r *BannerRepository) UpdateBanner(c context.Context, bannerID int, headerBanner model.HeaderBanner) error {
	return nil
}
func (r *BannerRepository) DeleteBanner(c context.Context, bannerID int) (mes string, err error) {
	cfg := pc.GetConfig(c)

	ctx, cancel := context.WithTimeout(c, cfg.PSQL.Timeout)
	defer cancel()

	var (
		res int
		row sql.NullString
	)

	query := "select o_res, o_mes from public.fn_banner_get($1, $2)"

	//	(i_banner_id int,
	//	out o_res int,
	//	out o_mes text)

	if err = r.db.DB.QueryRow(ctx, query, bannerID).Scan(&row, &res, &mes); err != nil {
		err = errors.Wrap(err)
		return
	}
	if res != 0 {
		err = errors.New(http.StatusNotFound, mes)
		return
	}

	return
}
