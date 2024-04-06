package thttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (r *Router) setRoutingTable() {

	r.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Router.GET("/user_banner")
	banner := r.Router.Group("/banner")
	{
		banner.POST("")
		banner.GET("")
		banner.PATCH("/:id")
		banner.DELETE("/:id")
	}

	r.Router.NoRoute(func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, "Маршрут не найден")
	})

}

//func middlewareGetUser() gin.HandlerFunc {
//
//	return func(c *gin.Context) {
//
//		claims := jwt.ExtractClaims(c)
//		userID, err := uuid.FromString(claims["user_id"].(string))
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Пользователь не опознан системой."})
//			return
//		}
//
//		c.Request = c.Request.WithContext(pc.AddUserID(c.Request.Context(), userID))
//		c.Next()
//	}
//
//}
