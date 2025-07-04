package service

import "darts/pkg/repository"

type Authorization interface {
}

type Game interface {
}

type Service struct {
	Authorization
	Game
}

func NewService(repository *repository.Repository) *Service {
	return &Service{}
}
