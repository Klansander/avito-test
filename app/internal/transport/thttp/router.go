package thttp

import (
	"context"
	"net/http"
	"os"
	"sync"

	"avito/app/internal/service"
	pc "avito/app/pkg/context"
	"avito/app/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Router struct {
	service *service.Service
	Router  *gin.Engine
}

var once sync.Once
var r *Router

// NewRouter Инициализация роутера
func NewRouter(ctx context.Context, service *service.Service) (*Router, error) {

	var err error
	once.Do(func() {
		r = &Router{service: service}
		gin.SetMode(gin.ReleaseMode)

		r.Router = gin.New()
		r.Router.Routes()

		r.Router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
			Output: os.Stdout,
			//SkipPaths: []string{"/skip"}, // Можете добавить пути, которые вы хотите пропустить при логгировании
		}))

		r.Router.Use(gin.Recovery())

		// Мержим два контекста
		r.Router.Use(transferParentContext(ctx))

		// Единая точка обработки ошибок
		r.Router.Use(errorHandler())

	})

	if err != nil {
		return nil, err
	}

	r.SetRoutingTable()

	return r, nil

}

func transferParentContext(ctx context.Context) gin.HandlerFunc {

	return func(c *gin.Context) {

		ctxMegge := pc.Link(c.Request.Context(), ctx)
		c.Request = c.Request.WithContext(ctxMegge)
		c.Next()

	}

}

func parseError(err error) (int, string) {

	code := http.StatusInternalServerError
	message := err.Error()

	var target *errors.ErrorApp
	if errors.As(err, &target) {
		code = target.Code
		message = errors.Cause(err).Error()
	}

	return code, message

}

func errorHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			code, message := parseError(err.Err)

			if code == http.StatusInternalServerError {

				c.AbortWithStatusJSON(code, gin.H{"error": "Запрос выполнен с ошибкой. Свяжитесь со службой технической поддержки."})
			} else {
				c.AbortWithStatusJSON(code, gin.H{"error": message})
			}

		}
	}

}
