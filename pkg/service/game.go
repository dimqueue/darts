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

func (s *GameService) CreateGame(userId int, lang string) (int, error) {
	var game model.Game
	//add validation language
	game.Language = lang
	game.Status = "in_progress"
	game.WordId = 1
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
