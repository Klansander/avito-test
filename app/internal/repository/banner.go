package repository

import (
	"avito/app/internal/model"
	pc "avito/app/pkg/context"
	"avito/app/pkg/errors"
	"avito/app/pkg/postgresql"
	"avito/app/pkg/redis"
	"context"
	"database/sql"
	json "github.com/json-iterator/go"
	"sync"

	"fmt"
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type Banner interface {
	UserBannerRedis(c context.Context, userBannerQuery model.UserBannerQueryGet, isAdmin bool) (data map[string]interface{}, err error)
	UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet, isAdmin bool) (data map[string]interface{}, err error)
	ListBanner(c context.Context, userBannerQuery model.UserBannerQueryList) (data []model.Banner, err error)
	CreateBanner(c context.Context, headerBanner model.NewBanner) (id int, err error)
	UpdateBanner(c context.Context, bannerID int, headerBanner map[string]interface{}) error
	DeleteBanner(c context.Context, bannerID int) (mes string, err error)
	SaveVersionBanner(c context.Context)
	GetVersionBanner(c context.Context, headerBanner model.BannerVersion) (dataArr []model.Banner, err error)
	DeleteBannerByTagOrFeature(c context.Context, banner model.BannerTagOrFeatureID, wg *sync.WaitGroup) (err error)
}

type BannerRepository struct {
	db  *postgresql.Postgres
	dbR *rediscl.Redis
}

func NewBanner(db *postgresql.Postgres, dbR *rediscl.Redis) *BannerRepository {

	return &BannerRepository{db: db, dbR: dbR}

}

func (r *BannerRepository) UserBannerRedis(c context.Context, userBannerQuery model.UserBannerQueryGet, isAdmin bool) (data map[string]interface{}, err error) {

	valKey, err := r.dbR.DB.Get(c, fmt.Sprintf("new%d%dIsActive", userBannerQuery.TagID, userBannerQuery.FeatureID)).Bool()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(http.StatusNotFound, "Баннер не найден")
	}

	if !valKey && !isAdmin {
		return nil, errors.New(http.StatusNotFound, "Баннер не найден")
	}

	val, err := r.dbR.DB.Get(c, fmt.Sprintf("new%d%d", userBannerQuery.TagID, userBannerQuery.FeatureID)).Result()
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, "")
	}

	if err = json.Unmarshal([]byte(val), &data); err != nil {
		return nil, errors.New(http.StatusInternalServerError, "")
	}

	return
}

func (r *BannerRepository) UserBanner(c context.Context, userBannerQuery model.UserBannerQueryGet, isAdmin bool) (data map[string]interface{}, err error) {

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

	if err = r.db.DB.QueryRow(ctx, query, userBannerQuery.TagID, userBannerQuery.FeatureID, isAdmin).Scan(&row, &res, &mes); err != nil {
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
		res  int
		mes  string
		row  sql.NullString
		data model.Banner
	)

	if err = tx.QueryRow(ctx, query, headerBanner.IsActive, headerBanner.TagID, headerBanner.FeatureID, time.Now(), time.Now(), headerBanner.Content).Scan(&id, &res, &mes); err != nil {
		err = errors.Wrap(err)
		return
	}

	if res != 0 {
		err = errors.New(http.StatusBadRequest, mes)
		return
	}

	query = "select o_json,o_res,o_mes from public.fn_banner_get_by_id($1)"
	if err = tx.QueryRow(ctx, query, id).Scan(&row, &res, &mes); err != nil {
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

	r.saveRedis(c, data)

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
	r.updateRedis(data)
	if err != nil {
		return errors.Wrap(err)
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

	query = "select o_res,o_mes from public.fn_banner_upd($1,$2,$3)"
	if err = tx.QueryRow(ctx, query, bannerID, updateBanner.TagID, updateBanner.FeatureID).Scan(&res, &mes); err != nil {
		err = errors.Wrap(err)
		return
	}
	if res != 0 {
		err = errors.New(http.StatusNotFound, mes)
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

	if err := json.Unmarshal(resByte, &res); err != nil {
		return model.Banner{}, errors.Wrap(err)
	}
	return res, nil

}

func (r *BannerRepository) updateRedis(obj map[string]interface{}) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonData, err := json.Marshal(obj)
	if err != nil {
		logrus.Errorln("Ошибка при сериализации данных:", err)
		return
	}

	err = r.dbR.DB.LPush(ctx, strconv.Itoa(int(obj["banner_id"].(float64))), jsonData).Err()
	if err != nil {
		logrus.Errorln("Ошибка при добавлении новой версии:", err)
		return
	}

	listLen, err := r.dbR.DB.LLen(ctx, strconv.Itoa(int(obj["banner_id"].(float64)))).Result()
	if err != nil {
		logrus.Errorln("Ошибка при получении длины списка:", err)
		return
	}

	if listLen > 3 {
		err := r.dbR.DB.RPop(ctx, strconv.Itoa(int(obj["banner_id"].(float64)))).Err()
		if err != nil {
			logrus.Errorln("Ошибка при удалении старой версии:", err)
			return
		}
	}

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

func (r *BannerRepository) SaveVersionBanner(c context.Context) {
	data, err := r.ListBanner(c, model.UserBannerQueryList{})
	if err != nil {
		logrus.Errorln("Ошибка добавления объекта в Redis:", err)
		return
	}
	for _, item := range data {
		r.saveRedis(c, item)

	}
	fmt.Println(time.Now())

}

func (r *BannerRepository) saveRedis(c context.Context, item model.Banner) {

	for _, tagId := range item.TagID {

		key := fmt.Sprintf("new%d%d", tagId, item.FeatureID)
		keyActive := fmt.Sprintf("new%d%dIsActive", tagId, item.FeatureID)

		jsonData, err := json.Marshal(item.Content)
		if err != nil {
			logrus.Errorln("Ошибка добавления объекта в Redis:", err)
			return
		}

		err = r.dbR.DB.Set(c, key, jsonData, 5*time.Minute).Err()
		if err != nil {
			logrus.Errorln("Ошибка добавления объекта в Redis:", err)
			return
		}
		err = r.dbR.DB.Set(c, keyActive, item.IsActive, 5*time.Minute).Err()
		if err != nil {
			logrus.Errorln("Ошибка добавления объекта в Redis:", err)
			return
		}
	}

}

func (r *BannerRepository) GetVersionBanner(c context.Context, headerBanner model.BannerVersion) (dataArr []model.Banner, err error) {

	dataArr = make([]model.Banner, 0, 3)

	key := fmt.Sprintf("%d", headerBanner.BannerID)

	startIndex := 0
	stopIndex := -1
	fmt.Println(headerBanner)
	if headerBanner.Version != nil {
		startIndex = *headerBanner.Version - 1
		stopIndex = *headerBanner.Version - 1
	}

	result, err := r.dbR.DB.LRange(c, key, int64(startIndex), int64(stopIndex)).Result()
	if err != nil {
		logrus.Errorln("Ошибка при получении элементов из списка:", err)
		return
	}
	if len(result) == 0 {
		err = errors.New(http.StatusNotFound, "Версии баннера не найдены")
		return
	}

	// Вывести полученные элементы
	for _, value := range result {

		var data model.Banner
		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			logrus.Errorln("Ошибка при получении элементов из списка:", err)
			return
		}
		dataArr = append(dataArr, data)

	}

	return
}

func (r *BannerRepository) DeleteBannerByTagOrFeature(c context.Context, banner model.BannerTagOrFeatureID, wg *sync.WaitGroup) (err error) {

	cfg := pc.GetConfig(c)

	ctx, cancel := context.WithTimeout(c, cfg.PSQL.Timeout)
	defer cancel()

	query := "select o_res,o_mes from public.fn_banner_get_by_tag_or_feature_id($1, $2)"

	//	i_tag_id int,
	//	i_feature_id int,

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

	var (
		res int
		mes string
	)

	if err = r.db.DB.QueryRow(ctx, query, banner.TagID, banner.FeatureID).Scan(&res, &mes); err != nil {
		err = errors.Wrap(err)
		return
	}

	if res != 0 {
		err = errors.New(http.StatusNotFound, mes)
		return
	}

	wg.Add(1)

	go func(c context.Context) {
		defer wg.Done()

		tx, err := r.db.DB.Begin(context.WithoutCancel(c))
		if err != nil {
			logrus.Errorln("err", err)
			return
		}

		defer func() {
			if err != nil {
				_ = tx.Rollback(context.WithoutCancel(c))
			} else {
				_ = tx.Commit(context.WithoutCancel(c))
			}
		}()

		query = "select  from public.fn_banner_del_by_tag_or_feature_id($1,$2)"
		if _, err = tx.Exec(context.WithoutCancel(c), query, banner.TagID, banner.FeatureID); err != nil {
			logrus.Errorln("err", err)
			return
		}

	}(c)

	return nil
}
