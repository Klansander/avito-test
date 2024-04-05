package context

import (
	"context"
	"fmt"

	"avito/pkg/config"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type ctxConfig struct{}
type ctxUser struct{}

// AddConfig adds config to context
func AddConfig(ctx context.Context, c *config.Config) context.Context {

	return context.WithValue(ctx, ctxConfig{}, c)

}

// GetConfig returns config from context
func GetConfig(ctx context.Context) *config.Config {

	if c, ok := ctx.Value(ctxConfig{}).(*config.Config); ok {
		return c
	}

	logrus.Error("Отсутствует инициализация системы конфигурации")

	return nil

}

func GetUserID(ctx context.Context) uuid.UUID {

	if id, ok := ctx.Value(ctxUser{}).(uuid.UUID); ok {
		fmt.Println(id)
		return id
	}
	fmt.Println("!!")

	return uuid.Nil

}
