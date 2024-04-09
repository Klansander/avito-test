package postgresql

import (
	"context"
	"fmt"
	"sync"

	pc "avito/pkg/context"
	"avito/pkg/errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func New(ctx context.Context) (*Postgres, error) {

	var err error

	pgOnce.Do(func() {

		var (
			pgxCfg *pgxpool.Config
			db     *pgxpool.Pool
		)

		cfg := pc.GetConfig(ctx).PSQL
		dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

		if pgxCfg, err = pgxpool.ParseConfig(dsn); err != nil {
			err = errors.Wrap(err)
			return
		}

		if db, err = pgxpool.ConnectConfig(ctx, pgxCfg); err != nil {
			err = errors.Wrap(err)
			return
		}

		pgInstance = &Postgres{db}

	})

	if err != nil {
		return nil, err
	}

	err = pgInstance.DB.Ping(ctx)

	return pgInstance, errors.Wrap(err)

}

func (pg *Postgres) Close() {

	pgInstance.DB.Close()

}
