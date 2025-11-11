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
	GetGameById(userId, gameId int) (*model.Game, error)
	UpdateGame(gameId int) (*model.Game, error)
	DeleteGame(gameId int) (*model.Game, error)
}

type Guess interface {
	CreateGuess(i int) error
	GetAllGuessByGame(i int) error
	GetGuessById(i int) error
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
