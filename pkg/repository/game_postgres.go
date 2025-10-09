package repository

import (
	"fmt"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
)

type GamePostgres struct {
	db *sqlx.DB
}

func NewGamePostgres(db *sqlx.DB) *GamePostgres {
	return &GamePostgres{db: db}
}

func (r *GamePostgres) CreateGame(userId int, game model.Game) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id,word_id,status,started_at) VALUES ($1,$2,$3) RETURNING id", gamesTable)
	row := r.db.QueryRow(query, userId, game.WordId, game.Status, game.StartedAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
