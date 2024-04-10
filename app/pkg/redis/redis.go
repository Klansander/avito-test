package rediscl

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"

	pc "avito/app/pkg/context"
	"avito/app/pkg/errors"
)

type Redis struct {
	DB *redis.Client
}

var (
	reInstance *Redis
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

		reInstance = &Redis{client}

	})

	if err != nil {
		return nil, err
	}
	_, err = reInstance.DB.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return reInstance, errors.Wrap(err)

}

func (pg *Redis) Close() {

	reInstance.DB.Close()

}
