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
	CreateGuess(userId, gameId int, word string) (int, error)
	GetAllGuessByGame(userId, gameId int) ([]model.Guess, error)
	GetGuessById(i int) error
}

type Word interface {
	SelectWord(language string) (*model.Word, error)
	GetWordById(wordId int) (*model.Word, error)
}

type Service struct {
	Authorization
	Game
	Guess
	Word
}

func NewService(repos *repository.Repository, computeClient *connections.ComputeClientService) *Service {
	wordService := NewWordService(repos.Word)
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Game:          NewGameService(repos.Game, wordService, computeClient),
		Guess:         NewGuessService(repos.Guess, repos.Game, repos.Word, computeClient),
		Word:          wordService,
	}
}
