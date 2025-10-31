package service

import (
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Game interface {
	CreateGame(userId int) (int, error)
	GetAllGames(userId int) ([]model.Game, error)
	GetGameById(gameId int) (model.Game, error)
	UpdateGame(gameId int) (model.Game, error)
	DeleteGame(gameId int) (model.Game, error)
}

type Guess interface {
	CreateGuess(i int) error
	GetAllGuessByGame(i int) error
	GetGuessById(i int) error
}

type Service struct {
	Authorization
	Game
	Guess
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Game:          NewGameService(repos.Game),
		Guess:         NewGuessService(repos.Guess),
	}
}
