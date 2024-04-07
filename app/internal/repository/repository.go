package repository

import (
	"avito/pkg/postgresql"
	"sync"
)

type Repository struct {
	Banner Banner
}

var once sync.Once
var repository *Repository

func NewRepository(db *postgresql.Postgres) *Repository {

	once.Do(func() {
		repository = &Repository{
			Banner: NewBanner(db),
		}
	})

	return repository

}
