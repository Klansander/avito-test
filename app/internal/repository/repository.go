package repository

import (
	"avito/app/pkg/postgresql"
	"avito/app/pkg/redis"
	"sync"
)

type Repository struct {
	Banner Banner
}

var once sync.Once
var repository *Repository

func NewRepository(db *postgresql.Postgres, dbR *rediscl.Redis) *Repository {

	once.Do(func() {
		repository = &Repository{
			Banner: NewBanner(db, dbR),
		}
	})

	return repository

}
