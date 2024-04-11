package thttp

import (
	"avito/app/internal/model"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"sync"
)

var Wg sync.WaitGroup

func (r *Router) SetRoutingTable() {

	r.Router.GET("/swagger/index.html", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Router.GET("/user_banner", r.authorize(model.UserToken), r.userBanner)
	banner := r.Router.Group("/banner")
	{
		banner.GET("", r.authorize(model.AdminToken), r.listBanner)
		banner.POST("", r.authorize(model.AdminToken), r.createBanner)
		banner.PATCH("/:id", r.authorize(model.AdminToken), r.updateBanner)
		banner.DELETE("/:id", r.authorize(model.AdminToken), r.deleteBanner)
		banner.DELETE("/", r.authorize(model.AdminToken), r.deleteBannerByTagIdOrFeatureId)
		banner.GET("/version", r.authorize(model.AdminToken), r.getVersionBanner)

	}

	r.Router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, "Маршрут не найден")
	})

	go func() {
		Wg.Wait()
	}()

}

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
			if token == model.AdminToken {
				c.Set("adm", true)
			}
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Пользователь не имеет доступа"})
		}

	}

}

// Функция для проверки токена
func _isValidToken(c *gin.Context, token string, checkToken string) bool {
	if token == model.AdminToken {
		c.Set("adm", true)
	} else if token == checkToken {
		c.Set("adm", false)
	} else {
		return false
	}

	return true
}
