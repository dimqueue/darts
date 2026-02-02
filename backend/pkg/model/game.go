package model

import "time"

type Game struct {
	Id        int64      `db:"id" json:"id"`
	UserId    int64      `db:"user_id" json:"user_id"`
	WordId    int64      `db:"word_id" json:"word_id"`
	Status    string     `db:"status" json:"status"`
	Language  string     `db:"language" json:"language"`
	StartedAt time.Time  `db:"started_at" json:"started_at"`
	EndedAt   *time.Time `db:"ended_at" json:"ended_at,omitempty"`
	ExpiresAt *time.Time `db:"expires_at" json:"expires_at,omitempty"`
}

type Guess struct {
	Id        int64     `db:"id" json:"id"`
	GameId    int64     `db:"game_id" json:"game_id"`
	GuessWord string    `db:"guess_word" json:"guess_word"`
	Distance  int       `db:"distance" json:"distance"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type GuessResult struct {
	Rank         int  `json:"rank"`
	Found        bool `json:"found"`
	InVocabulary bool `json:"in_vocabulary"`
}
