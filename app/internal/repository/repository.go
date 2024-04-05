package repository

import (
	"avito/pkg/postgresql"
	"sync"
)

type Repository struct {
}

var once sync.Once
var repository *Repository

func NewRepository(db *postgresql.Postgres) *Repository {

	once.Do(func() {
		repository = &Repository{}
	})

	return repository

}
