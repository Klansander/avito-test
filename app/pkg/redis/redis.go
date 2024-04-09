package rediscl

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"

	pc "avito/pkg/context"
	"avito/pkg/errors"
)

type Redis struct {
	DB *redis.Client
}

var (
	pgInstance *Redis
	pgOnce     sync.Once
)

func New(ctx context.Context) (*Redis, error) {

	var err error

	pgOnce.Do(func() {

		cfg := pc.GetConfig(ctx).Redis

		addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		client := redis.NewClient(&redis.Options{
			Addr: addr,
		})

		pgInstance = &Redis{client}

	})

	if err != nil {
		return nil, err
	}

	return pgInstance, errors.Wrap(err)

}

func (pg *Redis) Close() {

	pgInstance.DB.Close()

}
