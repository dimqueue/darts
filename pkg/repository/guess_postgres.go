package repository

import (
	"github.com/jmoiron/sqlx"
)

type GuessPostgres struct {
	db *sqlx.DB
}

func NewGuessPostgres(db *sqlx.DB) *GuessPostgres {
	return &GuessPostgres{db: db}
}

func (r *GuessPostgres) CreateGuess(gameId int, guess string, distance int) error {
	return nil
}

func (r *GuessPostgres) GetAllGuessByGame(i int) error {
	return nil
}

func (r *GuessPostgres) GetGuessById(i int) error {
	return nil
}
