package repository

import (
	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
}

type Game interface {
	CreateGame(userId int, game *model.Game) (int, error)
	GetAllGames(userId int) ([]model.Game, error)
	GetGameById(gameId int) (*model.Game, error)
	UpdateGame(gameId int) (*model.Game, error)
	DeleteGame(gameId int) (*model.Game, error)
}

type Guess interface {
	CreateGuess(gameId int, guess string, distance int) error
	GetAllGuessByGame(gameId int) error
	GetGuessById(guessId int) error
}

type Repository struct {
	Authorization
	Game
	Guess
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Game:          NewGamePostgres(db),
		Guess:         NewGuessPostgres(db),
	}
}
