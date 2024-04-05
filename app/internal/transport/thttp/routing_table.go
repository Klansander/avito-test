package thttp

import (
	"net/http"

	"avito/docs"
	pc "avito/pkg/context"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (r *Router) setRoutingTable() {

	docs.SwaggerInfo.BasePath = "/api/v2"
	r.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Router.NoRoute(func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, "Маршрут не найден")
	})

}

func middlewareGetUser() gin.HandlerFunc {

	return func(c *gin.Context) {

		claims := jwt.ExtractClaims(c)
		userID, err := uuid.FromString(claims["user_id"].(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Пользователь не опознан системой."})
			return
		}

		c.Request = c.Request.WithContext(pc.AddUserID(c.Request.Context(), userID))
		c.Next()
	}

}

// func checkID() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Получаем значение параметра "id" из URI
// 		idParam := c.Param("id")

// 		// Преобразуем значение параметра "id" в число
// 		id, err := strconv.Atoi(idParam)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter in URI"})
// 			c.Abort()
// 			return
// 		}

// 		// Дополнительные проверки ID могут быть добавлены здесь

// 		// Устанавливаем значение ID в контекст для последующих обработчиков
// 		c.Set("id", id)

// 		// Продолжаем выполнение цепочки обработчиков
// 		c.Next()
// 	}
// }

func (r *Router) saveClient(ctx *gin.Context) {
	// Вытаскиваем user_id из контекста
	// Вызываем получение структуры пользователя r.getUser
	// Сохраняем пользователя в мапе
}
