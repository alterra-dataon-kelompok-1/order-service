package order

import "github.com/alterra-dataon-kelompok-1/order-service/internal/repository"

type Service interface {
}

type service struct {
	repository repository.Repository
}

func NewService(repository repository.Repository) Service {
	return &service{repository}
}
