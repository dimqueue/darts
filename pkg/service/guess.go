package service

import "github.com/dimqueue/darts/pkg/repository"

type GuessService struct {
	repo repository.Guess
}

func NewGuessService(repo repository.Guess) *GuessService {
	return &GuessService{repo: repo}
}

func (s *GuessService) CreateGuess(i int) error {
	return nil
}

func (s *GuessService) GetAllGuessByGame(i int) error {
	return nil
}

func (s *GuessService) GetGuessById(i int) error {
	return nil
}
