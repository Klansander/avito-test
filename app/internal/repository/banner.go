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

	query := "select v_id from public.fn_banner_ins($1, $2, $3, $4)"

	//(i_is_active boolean,
	//	i_tag_id int[],
	//	i_feature_id int,
	//	i_content json)
	//returns int

	if err = r.db.DB.QueryRow(ctx, query, headerBanner.IsActive, headerBanner.TagID, headerBanner.FeatureID).Scan(&id); err != nil {
		err = errors.Wrap(err)
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

	query := "select o_json,o_res, o_mes from public.fn_banner_get($1, $2)"

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
	defer tx.Rollback(ctx)

	if _, ok := headerBanner["tag_id"]; ok == false {
		for k, v := range data {
			if k == "feature_id" {
				r.updateBanner(c, tx, bannerID, k, "banners", v)
				continue
			}
			r.updateBanner(c, tx, bannerID, k, "contents", v)
		}
		return nil
	}

	mes, _ = r.DeleteBanner(c, bannerID)

	//updateBanner, _ := _mergeData(data, headerBanner)

	return tx.Commit(ctx)

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

	query := "select o_res, o_mes from public.fn_banner_get($1, $2)"

	//	(i_banner_id int,
	//	out o_res int,
	//	out o_mes text)

	if err = r.db.DB.QueryRow(ctx, query, bannerID).Scan(&res, &mes); err != nil {
		err = errors.Wrap(err)
		return
	}
	if res != 0 {
		err = errors.New(http.StatusNotFound, mes)
		return
	}

	return
}
func (s *BannerRepository) updateBanner(ctx context.Context, tx pgx.Tx, bannerID int, field string, table string, value interface{}) error {

	query := fmt.Sprintf(
		"update public.%s "+
			"set %s = $1, "+
			"updated_by = $2, "+
			"updatedat = $3 "+
			"where id = $4",
		table, field,
	)

	if _, err := tx.Exec(ctx, query, value); err != nil {
		return err
	}

	return nil

}

func (s *BannerRepository) updateStructBanner(structOld map[string]interface{}, structNew map[string]interface{}) error {

	// Проходим по ключам первой мапы
	for key, valueOld := range structOld {
		// Проверяем, есть ли ключ во второй мапе
		if valueNew, ok := structNew[key]; ok {
			// Сравниваем значения
			if !isEqual(valueOld, valueNew) {
				// Если значения отличаются, обновляем значение в первой мапе
				structOld[key] = valueNew
			}
		}
	}

	return nil

}

// Функция для сравнения значений интерфейсов
func isEqual(value1, value2 interface{}) bool {
	// Если типы значений не совпадают, считаем их разными
	if fmt.Sprintf("%T", value1) != fmt.Sprintf("%T", value2) {
		return false
	}

	// Сравниваем значения
	return value1 == value2
}
