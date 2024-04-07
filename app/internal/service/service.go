package service

import (
	"avito/internal/repository"
	"sync"
)

type Service struct {
	Banner Banner
}

var once sync.Once
var service *Service

func NewService(r *repository.Repository) *Service {

	once.Do(func() {
		service = &Service{
			Banner: NewBanner(r),
		}
	})

	return service

}
