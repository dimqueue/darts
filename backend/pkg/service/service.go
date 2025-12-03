package service

import (
	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user model.User) (int64, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int64, error)
}

type Word interface {
	SelectWord(language string) (*model.Word, error)
	GetWordById(wordId int64) (*model.Word, error)
}

type Game interface {
	CreateGame(userId int64, lang string) (int64, error)
	GetAllGames(userId int64) ([]model.Game, error)
	GetGameById(userId, gameId int64) (*model.Game, error)
	UpdateGame(gameId int64) (*model.Game, error)
	UpdateGameStatus(gameId int64, status string) error
	DeleteGame(gameId int64) (*model.Game, error)
	MakeGuess(userId, gameId int64, guess string) (int, error)
	GetAllGuessByGame(userId, gameId int64) ([]model.Guess, error)
}

type Stats interface {
	InitializeStats(tx *sqlx.Tx, userId int64) error
	UpdateGameEndStats(tx *sqlx.Tx, update model.StatisticsUpdate) error
	GetStatistics(userId int64) (*model.UserStatistics, error)
	GetLanguageStats(userId int64, language string) (*model.UserLanguageStats, error)
	GetAllLanguageStats(userId int64) ([]model.UserLanguageStats, error)
}

type Profile interface {
	GetProfile(userId int64) (*model.UserProfile, error)
	GetProfileSummary(userId int64) (*model.UserProfileSummary, error)
	GetProfileByUsername(username string) (*model.UserProfileSummary, error)
	UpdateProfile(userId int64, input model.UpdateProfileInput) error
	GetSettings(userId int64) (*model.UserSettings, error)
	UpdateSettings(userId int64, input model.UpdateSettingsInput) error
	GetStatistics(userId int64) (*model.UserStatistics, error)
	GetLanguageStats(userId int64, language string) (*model.UserLanguageStats, error)
	GetAllLanguageStats(userId int64) ([]model.UserLanguageStats, error)
}

type Leaderboard interface {
	GetLeaderboard(query model.LeaderboardQuery) (*model.LeaderboardResponse, error)
	GetLeaderboardWithUserRank(userId int64, query model.LeaderboardQuery) (*model.LeaderboardResponse, error)
	GetUserRank(userId int64, query model.LeaderboardQuery) (*int, error)
	GetAllUserRanks(userId int64) (*model.UserRanks, error)
}

type Service struct {
	Authorization
	Game
	Word
	Profile
	Leaderboard
	Stats
}

func NewService(repos *repository.Repository, computeClient connections.ComputeClient) *Service {
	wordService := NewWordService(repos.Word)
	statsService := NewStatsService(repos.Statistics, repos.TxManager)
	gameService := NewGameService(repos.Game, wordService, statsService, repos.TxManager, computeClient)
	profileService := NewProfileService(repos.Profile, statsService)
	leaderboardService := NewLeaderboardService(repos.Leaderboard)
	authService := NewAuthService(repos.Authorization, repos.Profile, statsService, repos.TxManager)

	return &Service{
		Authorization: authService,
		Game:          gameService,
		Word:          wordService,
		Profile:       profileService,
		Leaderboard:   leaderboardService,
		Stats:         statsService,
	}
}
