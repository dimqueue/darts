package model

type Game struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
}

type Guess struct {
	Id        int `json:"id"`
	GameId    int `json:"game_id"`
	UserId    int `json:"user_id"`
	WordId    int `json:"word_id"`
	Distance  int `json:"distance"`
	CreatedAt int `json:"created_at"`
}
type ListGame struct {
	Id     int
	UserId int
	GameId int
}
