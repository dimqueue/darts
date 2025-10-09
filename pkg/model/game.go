package model

import "time"

type Game struct {
	Id        int       `db:"id" json:"id"`
	UserId    int       `db:"user_id" json:"user_id"`
	WordId    int       `db:"word_id" json:"word_id"`
	Status    string    `db:"status" json:"status"`
	StartedAt time.Time `db:"started_at" json:"started-at"`
}

type Guess struct {
	Id        int       `db:"id" json:"id"`
	GameId    int       `db:"game_id" json:"game_id"`
	GuessWord string    `db:"guess_word" json:"guess_word"`
	Distance  int       `db:"distance" json:"distance"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
type ListGame struct {
	Id     int
	UserId int
	GameId int
}
