package thttp

import (
	"avito/internal/model"
	"avito/pkg/errors"
	"fmt"
	"github.com/gin-gonic/gin"
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
// @Success			200 {object} model.Banner
// @Router			/user_banner [get]
func (r *Router) userBanner(c *gin.Context) {

	// Получаем параметры из запроса
	var queryParams model.UserBannerQueryGet
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		_ = c.Error(errors.Wrap(err))
		return
	}
	fmt.Println(c.MustGet("adm").(bool))
	data, err := r.service.Banner.UserBanner(c.Request.Context(), queryParams, c.MustGet("adm").(bool))
	if err != nil {
		_ = c.Error(errors.Wrap(err))
		return
	}
	c.JSON(http.StatusOK, data)
}

// listBanner
// @Summary			Получение всех баннеров c фильтрацией по фиче и/или тегу
// @BannerId		listBanner
// @Tags			Banner - Баннеры
// @Param 			tag_id query int false "Тэг пользователя"
// @Param			feature_id query int false "Идентификатор фичи"
// @Param limit query int false "Значение размера пачки, необходимо для реализации постраничного отображения списка"
// @Param offset query int false "Значение смещения, необходимо для реализации постраничного отображения списка"
// @Param			token header string true "Токен админа" example(admin_token)
// @Success			200 {object} []model.Banner
// @Router			/banner [get]
func (r *Router) listBanner(c *gin.Context) {

	// Получили тело запроса и сразу его закрыли, что бы
	// можно было получить отмену контекста.
	var userBannerQuery model.UserBannerQueryList
	if err := c.ShouldBindQuery(&userBannerQuery); err != nil {
		_ = c.Error(errors.New(http.StatusBadRequest, err.Error()))
		return
	}
	_ = c.Request.Body.Close()

	data, err := r.service.Banner.ListBanner(c.Request.Context(), userBannerQuery)
	if err != nil {
		_ = c.Error(errors.Wrap(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": data})

}

// createBanner
// @Summary			Создание нового баннера
// @BannerId		createBanner
// @Tags			Banner - Баннеры
// @Param 			data formData model.HeaderBanner true "Параметры для создания баннера"
// @Param			token header string true "Токен админа" example(admin_token)
// @Success			200 {object} model.Banner
// @Router			/banner [get]
func (r *Router) createBanner(c *gin.Context) {

	var err error

	// Получили тело запроса и сразу его закрыли, что бы
	// можно было получить отмену контекста.
	var headerBanner model.NewBanner
	//data := make(map[string]interface{}, 4)
	if err = c.ShouldBindJSON(&headerBanner); err != nil {
		_ = c.Error(errors.New(http.StatusBadRequest, err.Error()))
		return
	}
	_ = c.Request.Body.Close()

	id, err := r.service.Banner.CreateBanner(c.Request.Context(), headerBanner)
	if err != nil {
		_ = c.Error(errors.Wrap(err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"banner_id": id})

}

// updateBanner
// @Summary			Обновление содержимого баннера
// @BannerId		updateBanner
// @Tags			Banner - Баннеры
// @Param			id path int true "Идентификатор баннера"
// @Param 			data formData model.HeaderBanner true "Параметры для изменения баннера"
// @Param			token header string true "Токен админа" example(admin_token)
// @Success			200 {object} model.Banner
// @Router			/banner/{id} [patch]
func (r *Router) updateBanner(c *gin.Context) {

	bannerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(errors.New(http.StatusBadRequest, err.Error()))
		return
	}

	// Получили тело запроса и сразу его закрыли, что бы
	// можно было получить отмену контекста.
	var headerBanner model.HeaderBanner
	if err := c.ShouldBindJSON(&headerBanner); err != nil {
		_ = c.Error(errors.New(http.StatusBadRequest, err.Error()))
		return
	}
	_ = c.Request.Body.Close()

	err = r.service.Banner.UpdateBanner(c.Request.Context(), bannerID, headerBanner)
	if err != nil {
		_ = c.Error(errors.Wrap(err))
		return
	}
}

// updateBanner
// @Summary			Обновление содержимого баннера
// @BannerId		updateBanner
// @Tags			Banner - Баннеры
// @Param			id path int true "Идентификатор баннера"
// @Param 			data formData model.HeaderBanner true "Параметры для изменения баннера"
// @Param			token header string true "Токен админа" example(admin_token)
// @Success			200 {object} model.Banner
// @Router			/banner/{id} [patch]
func (r *Router) deleteBanner(c *gin.Context) {

	bannerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		_ = c.Error(errors.New(http.StatusBadRequest, err.Error()))
		return
	}

	// Получили тело запроса и сразу его закрыли, что бы
	// можно было получить отмену контекста.
	var newProject model.HeaderBanner
	if err := c.ShouldBindJSON(&newProject); err != nil {
		_ = c.Error(errors.New(http.StatusBadRequest, err.Error()))
		return
	}
	_ = c.Request.Body.Close()

	mes, err := r.service.Banner.DeleteBanner(c.Request.Context(), bannerID)
	if err != nil {
		_ = c.Error(errors.Wrap(err))
		return
	}

	c.JSON(http.StatusNoContent, mes)
}
