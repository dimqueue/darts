package model

type Word struct {
	Id       int    `json:"id"`
	Word     string `json:"word"`
	Language string `json:"language"`
}
