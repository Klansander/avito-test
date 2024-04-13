package thttp

import (
	"avito/app/internal/model"
	"avito/app/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// userBanner
// @Summary			Получение баннера для пользователя
// @BannerId		userBanner
// @Tags			Banner - Баннеры
// @Param 			tag_id query int true "Тэг пользователя"
// @Param			feature_id query int true "Идентификатор фичи"
// @Param			use_last_revision query boolean false "Получать актуальную информацию" default false
// @Param			token header string false "Токен пользователя" example(user_token)
// @Success			200 {object} swagger.Banner "Баннер пользователя"
// @Failure 		400 {object} swagger.Error "Некорректные данные"
// @Failure 		401  "Пользователь не авторизован"
// @Failure 		403 "Пользователь не имеет доступа"
// @Failure 		404  "Баннер не найден"
// @Failure 		500 {object} swagger.Error "Внутренняя ошибка сервера"
// @Router			/user_banner [get]
func (r *Router) userBanner(c *gin.Context) {

	// Получаем параметры из запроса
	var queryParams model.UserBannerQueryGet
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		err = c.Error(errors.New(http.StatusBadRequest, "Проверьте валидность данных"))
		logrus.Errorln("Ошибка при получении баннера ", err)
		return
	}

	data, err := r.service.Banner.UserBanner(c.Request.Context(), queryParams, c.MustGet("adm").(bool))
	if err != nil {
		err = c.Error(errors.Wrap(err))
		logrus.Errorln("Ошибка при получении баннера ", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": data})
}

// listBanner
// @Summary			Получение всех баннеров c фильтрацией по фиче и/или тегу
// @BannerId		listBanner
// @Tags			Banner - Баннеры
// @Param 			tag_id query int false "Тэг пользователя"
// @Param			feature_id query int false "Идентификатор фичи"
// @Param limit query int false "Значение размера пачки"
// @Param offset query int false "Значение смещения"
// @Param			token header string true "Токен админа" example(admin_token)
// @Success			200 {object} []swagger.Banner "OK"
// @Failure 		400 {object} swagger.Error "Некорректные данные"
// @Failure 		401  "Пользователь не авторизован"
// @Failure 		403 "Пользователь не имеет доступа"
// @Failure 		500 {object} swagger.Error "Внутренняя ошибка сервера"
// @Router			/banner [get]
func (r *Router) listBanner(c *gin.Context) {

	var userBannerQuery model.UserBannerQueryList
	if err := c.ShouldBindQuery(&userBannerQuery); err != nil {
		err = c.Error(errors.New(http.StatusBadRequest, "Проверьте валидность данных"))
		logrus.Errorln("Ошибка при получении списка баннеров ", err)
		return
	}

	data, err := r.service.Banner.ListBanner(c.Request.Context(), userBannerQuery)
	if err != nil {
		err = c.Error(errors.Wrap(err))
		logrus.Errorln("Ошибка при получении списка баннеров ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": data})

}

// createBanner
// @Summary			Создание нового баннера
// @BannerId		createBanner
// @Tags			Banner - Баннеры
// @Param 			data formData swagger.HeaderBanner true "Параметры для создания баннера"
// @Param			token header string true "Токен админа" example(admin_token)
// @Success 		201 {object} swagger.CreateBanner "Идентификатор созданного баннера"
// @Failure 		400 {object} swagger.Error "Некорректные данные"
// @Failure 		401  "Пользователь не авторизован"
// @Failure 		403 "Пользователь не имеет доступа"
// @Failure 		500 {object} swagger.Error "Внутренняя ошибка сервера"
// @Router			/banner [post]
func (r *Router) createBanner(c *gin.Context) {

	var err error

	var headerBanner model.NewBanner
	if err = c.ShouldBindJSON(&headerBanner); err != nil {
		err = c.Error(errors.New(http.StatusBadRequest, "Проверьте валидность данных"))
		logrus.Errorln("Ошибка при создании баннера ", err)
		return
	}

	id, err := r.service.Banner.CreateBanner(c.Request.Context(), headerBanner)
	if err != nil {
		err = c.Error(errors.Wrap(err))
		logrus.Errorln("Ошибка при создании баннера ", err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"banner_id": id})

}

// updateBanner
// @Summary			Обновление содержимого баннера
// @BannerId		updateBanner
// @Tags			Banner - Баннеры
// @Param			id path int true "Идентификатор баннера"
// @Param 			data formData swagger.HeaderBanner true "Параметры для изменения баннера"
// @Param			token header string true "Токен админа" example(admin_token)
// @Success			200 "OK"
// @Failure 		400 {object} swagger.Error "Некорректные данные"
// @Failure 		401  "Пользователь не авторизован"
// @Failure 		403 "Пользователь не имеет доступа"
// @Failure 		404  "Баннер не найден"
// @Failure 		500 {object} swagger.Error "Внутренняя ошибка сервера"
// @Router			/banner/{id} [patch]
func (r *Router) updateBanner(c *gin.Context) {

	bannerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err = c.Error(errors.New(http.StatusBadRequest, "Проверьте валидность данных"))
		logrus.Errorln("Ошибка при обновлении баннера ", err)
		return
	}

	var headerBanner *model.HeaderBanner
	if err = c.ShouldBindJSON(&headerBanner); err != nil || headerBanner == nil {
		err = c.Error(errors.New(http.StatusBadRequest, "Проверьте валидность данных"))
		logrus.Errorln("Ошибка при обновлении баннера ", err)
		return
	}

	err = r.service.Banner.UpdateBanner(c.Request.Context(), bannerID, *headerBanner)
	if err != nil {
		err = c.Error(errors.Wrap(err))
		logrus.Errorln("Ошибка при обновлении баннера ", err)
		return
	}
}

// deleteBanner
// @Summary			Удаление баннера по идентификатору
// @BannerId		deleteBanner
// @Tags			Banner - Баннеры
// @Param			id path int true "Идентификатор баннера"
// @Param			token header string true "Токен админа" example(admin_token)
// @Success			204 "OK"
// @Failure 		400 {object} swagger.Error "Некорректные данные"
// @Failure 		401  "Пользователь не авторизован"
// @Failure 		403 "Пользователь не имеет доступа"
// @Failure 		404  "Баннер не найден"
// @Failure 		500 {object} swagger.Error "Внутренняя ошибка сервера"
// @Router			/banner/{id} [delete]
func (r *Router) deleteBanner(c *gin.Context) {

	bannerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err = c.Error(errors.New(http.StatusBadRequest, "Проверьте валидность данных"))
		logrus.Errorln("Ошибка при удалении баннера ", err)
		return
	}

	mes, err := r.service.Banner.DeleteBanner(c.Request.Context(), bannerID)
	if err != nil {
		err = c.Error(errors.Wrap(err))
		logrus.Errorln("Ошибка при удалении баннера ", err)
		return
	}

	c.JSON(http.StatusNoContent, mes)
}

// deleteBannerByTagIDOrFeatureID
// @Summary			Удаление баннера по тегу или фиче
// @BannerId		deleteBannerByTagIDOrFeatureID
// @Tags			Banner - Баннеры
// @Param 			tag_id query int false "Тэг пользователя"
// @Param			feature_id query int false "Идентификатор фичи"
// @Param			token header string true "Токен админа" example(admin_token)
// @Success			204 "OK"
// @Failure 		400 {object} swagger.Error "Некорректные данные"
// @Failure 		401  "Пользователь не авторизован"
// @Failure 		403 "Пользователь не имеет доступа"
// @Failure 		404  "Баннер не найден"
// @Failure 		500 {object} swagger.Error "Внутренняя ошибка сервера"
// @Router			/banner [delete]
func (r *Router) deleteBannerByTagIDOrFeatureID(c *gin.Context) {

	var banner model.BannerTagOrFeatureID
	if err := c.ShouldBindQuery(&banner); err != nil {
		err = c.Error(errors.New(http.StatusBadRequest, "Проверьте валидность данных"))
		logrus.Errorln("Ошибка при удалении баннера по тегу или фиче ", err)
		return
	}

	if banner.TagID == nil && banner.FeatureID == nil {
		err := c.Error(errors.New(http.StatusBadRequest, "Проверьте валидность данных"))
		logrus.Errorln("Ошибка при удалении баннера по тегу или фиче ", err)
		return
	}

	err := r.service.Banner.DeleteBannerByTagIDOrFeatureID(c.Request.Context(), banner, &Wg)
	if err != nil {
		err = c.Error(errors.Wrap(err))
		logrus.Errorln("Ошибка при удалении баннера по тегу или фиче ", err)
		return
	}

	c.JSON(http.StatusNoContent, "")
}

// getVersionBanner
// @Summary			Получение версий баннера
// @BannerId		getVersionBanner
// @Tags			Banner - Баннеры
// @Param 			banner_id query int true "Тэг пользователя"
// @Param			version query int false "Версия баннера" example(1-3)
// @Param			token header string true "Токен админа" example(admin_token)
// @Success			200 {object} swagger.Banner
// @Failure 		400 {object} swagger.Error "Некорректные данные"
// @Failure 		401  "Пользователь не авторизован"
// @Failure 		403 "Пользователь не имеет доступа"
// @Failure 		404  "Баннер не найден"
// @Failure 		500 {object} swagger.Error "Внутренняя ошибка сервера"
// @Router			/banner/version [get]
func (r *Router) getVersionBanner(c *gin.Context) {

	var version model.BannerVersion
	if err := c.ShouldBindQuery(&version); err != nil {
		err = c.Error(errors.New(http.StatusBadRequest, "Проверьте валидность данных"))
		logrus.Errorln("Ошибка при получение версий баннера ", err)
		return
	}

	data, err := r.service.Banner.GetVersionBanner(c.Request.Context(), version)
	if err != nil {
		err = c.Error(errors.Wrap(err))
		logrus.Errorln("Ошибка при получение версий баннера ", err)
		return
	}

	c.JSON(http.StatusOK, data)
}
