package service

import (
	"avito/internal/repository"
	"sync"
)

type Service struct {
}

var once sync.Once
var service *Service

func NewService(r *repository.Repository) *Service {

	once.Do(func() {
		service = &Service{}
	})

	return service

}
