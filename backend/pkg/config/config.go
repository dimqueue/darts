package config

import (
	"time"

	"github.com/spf13/viper"
)

var (
	GameTTL            = 2 * time.Hour
	GameComputeTimeout = 30 * time.Second
	ScorePerWin        = 100
	GameMaxGuesses     = 0 // 0 = unlimited

)

var (
	LeaderboardMinGamesGlobal = 3
	LeaderboardMinGamesPeriod = 1
	LeaderboardDefaultLimit   = 50
	LeaderboardMaxLimit       = 100
)

func LoadFromViper() {
	if v := viper.GetInt("game.ttl"); v > 0 {
		GameTTL = time.Duration(v) * time.Second
	}
	if v := viper.GetInt("game.computeTimeout"); v > 0 {
		GameComputeTimeout = time.Duration(v) * time.Second
	}
	if v := viper.GetInt("game.maxGuesses"); v > 0 {
		GameMaxGuesses = v
	}
	if v := viper.GetInt("leaderboard.minGamesGlobal"); v > 0 {
		LeaderboardMinGamesGlobal = v
	}
	if v := viper.GetInt("leaderboard.minGamesPeriod"); v > 0 {
		LeaderboardMinGamesPeriod = v
	}
	if v := viper.GetInt("leaderboard.defaultLimit"); v > 0 {
		LeaderboardDefaultLimit = v
	}
	if v := viper.GetInt("leaderboard.maxLimit"); v > 0 {
		LeaderboardMaxLimit = v
	}
}
