package service

import (
	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Game interface {
	CreateGame(userId int, lang string) (int, error)
	GetAllGames(userId int) ([]model.Game, error)
	GetGameById(userId, gameId int) (*model.Game, error)
	UpdateGame(gameId int) (*model.Game, error)
	DeleteGame(gameId int) (*model.Game, error)
}

type Guess interface {
	CreateGuess(userId, gameId int, guess string) (int, error)
	GetAllGuessByGame(i int) error
	GetGuessById(i int) error
}

type Service struct {
	Authorization
	Game
	Guess
}

func NewService(repos *repository.Repository, computeClient connections.Client) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Game:          NewGameService(repos.Game),
		Guess:         NewGuessService(repos.Guess, repos.Game, computeClient),
	}
}
