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
	if err != nil || !valKey && !isAdmin {
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
		row sql.NullString
	)

	query := "select o_json,o_res from public.fn_banner_get($1, $2,$3)"

	// i_tag_id int,
	// i_feature_id int,
	// i_is_admin bool,
	// out o_json json,
	//	out o_mes text

	if err = r.db.DB.QueryRow(ctx, query, userBannerQuery.TagID, userBannerQuery.FeatureID, isAdmin).Scan(&row, &res); err != nil {
		err = errors.Wrap(err)

		return
	}

	if res != 0 || !row.Valid {
		err = errors.New(http.StatusNotFound, "Баннер не найден")
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

	// i_tag_id int,
	// i_feature_id int,
	// i_limit int,
	// i_offset int,
	// out o_json json

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
	//	i_content json
	//	id int
	//  o_res int
	// 	o_mes text
	//	)

	tx, err := r.db.DB.Begin(ctx)
	if err != nil {
		return 0, errors.Wrap(err)
	}

	defer func() {
		txTransaction(c, tx, err)
	}()

	var (
		res  int
		mes  string
		row  sql.NullString
		data model.Banner
	)

	if err = tx.QueryRow(ctx, query, headerBanner.IsActive, headerBanner.TagID, headerBanner.FeatureID, time.Now(), time.Now(), headerBanner.Content).Scan(&id, &res, &mes); err != nil {
		err = errors.Wrap(err)
		return id, err
	}

	if res != 0 {
		err = errors.New(http.StatusBadRequest, mes)
		return id, err
	}

	query = "select o_json,o_res,o_mes from public.fn_banner_get_by_id($1)"
	if err = tx.QueryRow(ctx, query, id).Scan(&row, &res, &mes); err != nil {
		err = errors.Wrap(err)
		return id, err
	}

	if res != 0 {
		err = errors.New(http.StatusBadRequest, mes)
		return id, err
	}

	if row.Valid {
		if err = json.Unmarshal([]byte(row.String), &data); err != nil {
			err = errors.Wrap(err)
			return id, err
		}
	}

	r.saveRedis(c, data)

	return id, err

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
	//	out o_mes text
	//	out o_json json
	//	)

	if err = r.db.DB.QueryRow(ctx, query, bannerID).Scan(&row, &res, &mes); err != nil {
		return errors.Wrap(err)
	}
	if res != 0 {
		return errors.New(http.StatusNotFound, mes)
	}

	if row.Valid {
		if err = json.Unmarshal([]byte(row.String), &data); err != nil {
			return errors.Wrap(err)
		}
	}

	r.updateRedis(c, data)

	tx, err := r.db.DB.Begin(ctx)
	if err != nil {
		return errors.Wrap(err)
	}

	defer func() {
		txTransaction(c, tx, err)
	}()

	if _, ok := headerBanner[model.FieldTagID]; !ok {
		for k, v := range headerBanner {
			if k == model.FieldFeatureID {

				if err = r.updateBanner(c, tx, bannerID, model.FieldBannerID, k, model.TableContents, v); err != nil {
					return errors.Wrap(err)
				}
				continue
			}

			if err = r.updateBanner(c, tx, bannerID, model.FieldID, k, model.TableBanners, v); err != nil {
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

	//	i_banner_id int,
	//	i_tag_id int,
	//	i_feature_id int,
	//	out o_res int,
	//	out o_mes text)

	if err = tx.QueryRow(ctx, query, bannerID, updateBanner.TagID, updateBanner.FeatureID).Scan(&res, &mes); err != nil {
		return errors.Wrap(err)
	}
	if res != 0 {
		return errors.New(http.StatusNotFound, mes)
	}

	for k, v := range headerBanner {

		if k == model.FieldContent || k == model.FieldIsActive {
			if err = r.updateBanner(c, tx, bannerID, model.FieldID, k, model.TableBanners, v); err != nil {
				return errors.Wrap(err)
			}
		}

	}

	return nil

}

func (r *BannerRepository) updateBanner(ctx context.Context, tx pgx.Tx, bannerID int, whereField string, field string, table string, value interface{}) error {
	query := ""

	if field == model.FieldFeatureID {

		query = fmt.Sprintf(
			"update public.%s "+
				"set %s = $1 "+
				"where %s = $2",
			table, field, whereField,
		)
		if _, err := tx.Exec(ctx, query, value, bannerID); err != nil {
			return errors.Wrap(err)
		}
		return nil
	}

	query = fmt.Sprintf(
		"update public.%s "+
			"set %s = $1, "+
			"updated_at = $2 "+
			"where %s = $3",
		table, field, whereField,
	)
	if _, err := tx.Exec(ctx, query, value, time.Now(), bannerID); err != nil {
		return errors.Wrap(err)
	}

	return nil

}

func _mergeData(dataOld map[string]interface{}, dataNew map[string]interface{}) (res model.Banner, err error) {

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

	if err = json.Unmarshal(resByte, &res); err != nil {
		return model.Banner{}, errors.Wrap(err)
	}

	return res, nil

}

func (r *BannerRepository) updateRedis(c context.Context, obj map[string]interface{}) {

	cfg := pc.GetConfig(c)

	ctx, cancel := context.WithTimeout(c, cfg.Redis.Timeout)
	defer cancel()

	jsonData, err := json.Marshal(obj)
	if err != nil {
		logrus.Errorln("Ошибка при сериализации данных:", err)
		return
	}

	err = r.dbR.DB.LPush(ctx, strconv.Itoa(int(obj[model.FieldBannerID].(float64))), jsonData).Err()
	if err != nil {
		logrus.Errorln("Ошибка при добавлении новой версии:", err)
		return
	}

	listLen, err := r.dbR.DB.LLen(ctx, strconv.Itoa(int(obj[model.FieldBannerID].(float64)))).Result()
	if err != nil {
		logrus.Errorln("Ошибка при получении длины списка:", err)
		return
	}

	if listLen > cfg.Redis.LenStack {
		err := r.dbR.DB.RPop(ctx, strconv.Itoa(int(obj[model.FieldBannerID].(float64)))).Err()
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
		txTransaction(c, tx, err)
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

}

func (r *BannerRepository) saveRedis(c context.Context, item model.Banner) {
	cfg := pc.GetConfig(c)

	for _, tagID := range item.TagID {

		key := fmt.Sprintf("new%d%d", tagID, item.FeatureID)
		keyActive := fmt.Sprintf("new%d%dIsActive", tagID, item.FeatureID)

		jsonData, err := json.Marshal(item.Content)
		if err != nil {
			logrus.Errorln("Ошибка добавления объекта в Redis:", err)
			return
		}

		err = r.dbR.DB.Set(c, key, jsonData, cfg.Redis.TTL).Err()
		if err != nil {
			logrus.Errorln("Ошибка добавления объекта в Redis:", err)
			return
		}
		err = r.dbR.DB.Set(c, keyActive, item.IsActive, cfg.Redis.TTL).Err()
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
			logrus.Errorln("Ошибка при Unmarshal:", err)
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
	//	out o_res int,
	//	out o_mes text)

	var (
		res int
		mes string
	)

	if err = r.db.DB.QueryRow(ctx, query, banner.TagID, banner.FeatureID).Scan(&res, &mes); err != nil {

		return errors.Wrap(err)
	}

	if res != 0 {

		return errors.New(http.StatusNotFound, mes)
	}

	wg.Add(1)

	go func(c context.Context) {
		defer wg.Done()
		ctx := context.WithoutCancel(c)
		tx, err := r.db.DB.Begin(ctx)
		if err != nil {
			logrus.Errorln("err", err)
			return
		}

		defer func() {
			txTransaction(ctx, tx, err)
		}()

		query = "select  from public.fn_banner_del_by_tag_or_feature_id($1,$2)"
		if _, err = tx.Exec(ctx, query, banner.TagID, banner.FeatureID); err != nil {
			logrus.Errorln("Ошибка при удалении по тегу или фиче", err)
			return
		}

	}(c)

	return nil
}
func txTransaction(c context.Context, tx pgx.Tx, err error) {

	if err != nil {
		errRollback := tx.Rollback(c)
		if errRollback != nil {
			logrus.Errorln("Ошибка отката тразакции", errRollback)
		}
	} else {
		errCommit := tx.Commit(c)
		if err != nil {
			logrus.Errorln("Ошибка коммита тразакции", errCommit)
		}
	}

}
