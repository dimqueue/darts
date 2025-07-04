package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
}

type Game interface {
}

type Repository struct {
	Authorization
	Game
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
