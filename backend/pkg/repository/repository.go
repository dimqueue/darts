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
	CreateGame(game *model.Game) (int, error)
	GetAllGames(userId int) ([]model.Game, error)
	GetGameById(gameId int) (*model.Game, error)
	UpdateGame(gameId int) (*model.Game, error)
	UpdateGameStatus(gameId int, status string) error
	DeleteGame(gameId int) (*model.Game, error)
}

type Guess interface {
	CreateGuess(guess *model.Guess) error
	GetAllGuessByGame(gameId int) ([]model.Guess, error)
	GetGuessById(guessId int) error
}

type Word interface {
	GetWordById(wordId int) (*model.Word, error)
	GetRandomWordByLanguage(language string) (*model.Word, error)
	GetWordCountByLanguage(language string) (int, error)
}

type Repository struct {
	Authorization
	Game
	Guess
	Word
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Game:          NewGamePostgres(db),
		Guess:         NewGuessPostgres(db),
		Word:          NewWordPostgres(db),
	}
}
