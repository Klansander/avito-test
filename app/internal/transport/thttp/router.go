package thttp

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"sync"

	"avito/internal/service"
	pc "avito/pkg/context"
	"avito/pkg/errors"
	//"github.com/go-playground/validator/v10"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type Router struct {
	service  *service.Service
	enforcer *casbin.Enforcer
	jwt      *jwt.GinJWTMiddleware
	Router   *gin.Engine
}

var once sync.Once
var r *Router

func NewRouter(ctx context.Context, service *service.Service) (*Router, error) {

	var err error
	once.Do(func() {
		r = &Router{service: service}
		r.Router = gin.New()
		r.Router.Routes()

		gin.SetMode(gin.ReleaseMode)

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

	r.setRoutingTable()

	return r, nil

}

func transferParentContext(ctx context.Context) gin.HandlerFunc {

	return func(c *gin.Context) {

		// parent.ContextWithParentContext(ctx, ctxParent)
		// ctx.Next()
		// WithContext возвращает поверхностную копию контекста.
		// Возможно нужно использовать NewRequestWithContext или
		// чтобы сделать глубокую копию запроса с новым контекстом, Request.Clone.
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

			logrus.Error(err.Err)
		}
	}

}
