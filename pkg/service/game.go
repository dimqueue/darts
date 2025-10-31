package service

import (
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type GameService struct {
	repo repository.Game
}

func NewGameService(repo repository.Game) *GameService {
	return &GameService{repo: repo}
}

func (s *GameService) CreateGame(userId int) (int, error) {
	var game model.Game
	return s.repo.CreateGame(userId, game)
}

func (s *GameService) GetAllGames(userId int) ([]model.Game, error) {
	return s.repo.GetAllGames(userId)
}

func (s *GameService) GetGameById(gameId int) (model.Game, error) {
	return model.Game{}, nil
}

// ///
func (s *GameService) UpdateGame(gameId int) (model.Game, error) {
	return model.Game{}, nil
}
func (s *GameService) DeleteGame(gameId int) (model.Game, error) {
	return model.Game{}, nil
}
