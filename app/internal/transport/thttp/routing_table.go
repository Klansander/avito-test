package thttp

import (
	"avito/app/docs"
	"avito/app/internal/model"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"sync"
)

var Wg sync.WaitGroup

// SetRoutingTable функция с маршрутами
func (r *Router) SetRoutingTable() {
	docs.SwaggerInfo.BasePath = ""
	r.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Router.GET("/user_banner", r.authorize(model.UserToken), r.userBanner)
	banner := r.Router.Group("/banner")
	{
		banner.GET("", r.authorize(model.AdminToken), r.listBanner)
		banner.POST("", r.authorize(model.AdminToken), r.createBanner)
		banner.PATCH("/:id", r.authorize(model.AdminToken), r.updateBanner)
		banner.DELETE("/:id", r.authorize(model.AdminToken), r.deleteBanner)
		banner.DELETE("/", r.authorize(model.AdminToken), r.deleteBannerByTagIDOrFeatureID)
		banner.GET("/version", r.authorize(model.AdminToken), r.getVersionBanner)

	}

	r.Router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, "Маршрут не найден")
	})

	go func() {
		Wg.Wait()
	}()

}

// функция авторизации
func (r *Router) authorize(action string) gin.HandlerFunc {

	return func(c *gin.Context) {

		token := c.GetHeader("token")

		// Проверяем, присутствует ли токен
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
			return
		}

		// Проверяем наличие прав у пользователя
		if _isValidToken(c, token, action) {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Пользователь не имеет доступа"})
		}

	}

}

// Функция для проверки токена
func _isValidToken(c *gin.Context, token string, checkToken string) bool {

	switch token {
	case model.AdminToken:
		c.Set("adm", true)
		return true
	case checkToken:
		c.Set("adm", false)
		return true
	default:
		return false

	}

}
