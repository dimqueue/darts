package repository

import (
	"database/sql"
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

func (r *GamePostgres) CreateGame(userId int, game *model.Game) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (user_id,word_id,status,language,started_at) VALUES ($1,$2,$3,$4) RETURNING id", gamesTable)
	row := r.db.QueryRow(query, userId, game.WordId, game.Status, game.Language, game.StartedAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *GamePostgres) GetAllGames(userId int) ([]model.Game, error) {
	var games []model.Game

	query := fmt.Sprintf("SELECT * FROM %s WHERE g.user_id=$1",
		gamesTable)
	err := r.db.Select(&games, query, userId)
	return games, err
}

func (r *GamePostgres) GetGameById(gameId int) (*model.Game, error) {
	var game model.Game

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", gamesTable)
	err := r.db.Get(&game, query, gameId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game not found")
		}
		return nil, err
	}

	return &game, nil
}

func (r *GamePostgres) UpdateGame(gameId int) (*model.Game, error) {
	return nil, nil
}
func (r *GamePostgres) DeleteGame(gameId int) (*model.Game, error) {
	return nil, nil
}
