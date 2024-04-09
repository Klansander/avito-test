package repository

import (
	"avito/internal/model"
	pc "avito/pkg/context"
	"avito/pkg/errors"
	"avito/pkg/postgresql"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/jackc/pgx/v4"
	"net/http"
	"time"
)

type Banner interface {
	UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet, is_admin bool) (data map[string]interface{}, err error)
	ListBanner(c context.Context, userBannerQuery model.UserBannerQueryList) (data []model.Banner, err error)
	CreateBanner(c context.Context, headerBanner model.NewBanner) (id int, err error)
	UpdateBanner(c context.Context, bannerID int, headerBanner map[string]interface{}) error
	DeleteBanner(c context.Context, bannerID int) (mes string, err error)
}

type BannerRepository struct {
	db *postgresql.Postgres
}

func NewBanner(db *postgresql.Postgres) *BannerRepository {

	return &BannerRepository{db: db}

}

func (r *BannerRepository) UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet, is_admin bool) (data map[string]interface{}, err error) {

	cfg := pc.GetConfig(c)

	ctx, cancel := context.WithTimeout(c, cfg.PSQL.Timeout)
	defer cancel()

	var (
		res int
		mes string
		row sql.NullString
	)

	query := "select o_json,o_res,o_mes from public.fn_banner_get($1, $2,$3)"

	// i_tag_id int,
	//i_feature_id int,
	//i_is_admin bool,
	//out o_json json,
	//	out o_res int,
	//	out o_mes text

	if err = r.db.DB.QueryRow(ctx, query, userBannerQuery.TagID, userBannerQuery.FeatureID, is_admin).Scan(&row, &res, &mes); err != nil {
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

func (r *BannerRepository) CreateBanner(c context.Context, headerBanner model.NewBanner) (id int, err error) {

	cfg := pc.GetConfig(c)

	ctx, cancel := context.WithTimeout(c, cfg.PSQL.Timeout)
	defer cancel()

	query := "select v_id,o_res,o_mes from public.fn_banner_ins($1, $2, $3, $4,$5,$6)"

	//	(
	//	i_is_active boolean,
	//	i_tag_id int[],
	//	i_feature_id int,
	//	i_created_at timestamp,
	//	i_updated_at timestamp
	//	i_content json)
	//	returns int

	tx, err := r.db.DB.Begin(ctx)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	var (
		res int
		mes string
	)

	if err = tx.QueryRow(ctx, query, headerBanner.IsActive, headerBanner.TagID, headerBanner.FeatureID, time.Now(), time.Now(), headerBanner.Content).Scan(&id, &res, &mes); err != nil {
		err = errors.Wrap(err)
		return
	}

	if res != 0 {
		err = errors.New(http.StatusNotFound, mes)
		return
	}

	return

}
func (r *BannerRepository) UpdateBanner(c context.Context, bannerID int, headerBanner map[string]interface{}) (err error) {
	cfg := pc.GetConfig(c)

	ctx, cancel := context.WithTimeout(c, cfg.PSQL.Timeout)
	defer cancel()

	var (
		res  int
		row  sql.NullString
		mes  string
		data map[string]interface{}
	)

	query := "select o_json,o_res, o_mes from public.fn_banner_get_by_id($1)"

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

	if row.Valid {
		if err = json.Unmarshal([]byte(row.String), &data); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	tx, err := r.db.DB.Begin(ctx)
	if err != nil {
		return errors.Wrap(err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	fmt.Println(headerBanner)
	if _, ok := headerBanner["tag_id"]; ok == false {
		for k, v := range headerBanner {
			if k == "feature_id" {

				if err = r.updateBanner(c, tx, bannerID, k, "banner_id", "contents", v); err != nil {
					return errors.Wrap(err)
				}
				continue
			}

			if err = r.updateBanner(c, tx, bannerID, k, "id", "banners", v); err != nil {
				return errors.Wrap(err)
			}

		}
		return nil
	}

	updateBanner, err := _mergeData(data, headerBanner)
	if err != nil {
		return errors.Wrap(err)
	}

	query = "select * from public.fn_banner_upd($1,$2,$3)"
	if _, err = tx.Exec(ctx, query, bannerID, updateBanner.TagID, updateBanner.FeatureID); err != nil {
		err = errors.Wrap(err)
		return
	}

	for k, v := range headerBanner {
		if k == "content" || k == "is_active" {
			if err = r.updateBanner(c, tx, bannerID, "id", k, "banners", v); err != nil {
				return errors.Wrap(err)
			}
		}

	}

	return nil

}

func (s *BannerRepository) updateBanner(ctx context.Context, tx pgx.Tx, bannerID int, whereField string, field string, table string, value interface{}) error {

	query := fmt.Sprintf(
		"update public.%s "+
			"set %s = $1, "+
			"updated_at = $2 "+
			"where %s = $3",
		table, field, whereField,
	)
	fmt.Println(query, value)

	if _, err := tx.Exec(ctx, query, value, time.Now(), bannerID); err != nil {
		return err
	}

	return nil

}

func _mergeData(dataOld map[string]interface{}, dataNew map[string]interface{}) (model.Banner, error) {
	getByte, err := json.Marshal(dataOld)
	if err != nil {
		return model.Banner{}, errors.Wrap(err)
	}

	newByte, err := json.Marshal(dataNew)
	if err != nil {
		return model.Banner{}, errors.Wrap(err)
	}

	resByte, err := jsonpatch.MergePatch(getByte, newByte)
	if err != nil {
		return model.Banner{}, errors.Wrap(err)
	}
	var res model.Banner
	//for key, value1 := range dataOld {
	//	// Получаем значение по ключу из второй мапы
	//	if value2, ok := dataNew[key]; ok {
	//		// Если значения отличаются, добавляем их в новую мапу
	//		if value1 != value2 {
	//			dataOld[key] = value2
	//		}
	//	}
	//}
	fmt.Println(string(resByte))
	if err := json.Unmarshal(resByte, &res); err != nil {
		return model.Banner{}, errors.Wrap(err)
	}
	return res, nil

}

func (r *BannerRepository) DeleteBanner(c context.Context, bannerID int) (mes string, err error) {
	cfg := pc.GetConfig(c)

	ctx, cancel := context.WithTimeout(c, cfg.PSQL.Timeout)
	defer cancel()

	res := 0

	query := "select o_res, o_mes from public.fn_banner_del($1)"

	//	(i_banner_id int,
	//	out o_res int,
	//	out o_mes text)

	tx, err := r.db.DB.Begin(ctx)
	if err != nil {
		return "", errors.Wrap(err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	if err = tx.QueryRow(ctx, query, bannerID).Scan(&res, &mes); err != nil {
		err = errors.Wrap(err)
		return
	}
	if res != 0 {
		err = errors.New(http.StatusNotFound, mes)
		return
	}

	return
}
