package model

type Word struct {
	Id         int64  `json:"id" db:"id"`
	Word       string `json:"word" db:"word"`
	Language   string `json:"language" db:"language"`
	Difficulty int    `json:"difficulty" db:"difficulty"`
	IsActive   bool   `json:"is_active" db:"is_active"`
}
