package service

import (
	"context"

	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (int64, error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(token string) (int64, error) // ParseToken doesn't need ctx - it's CPU-only JWT parsing
}

type Word interface {
	SelectWord(ctx context.Context, language string) (*model.Word, error)
	GetWordById(ctx context.Context, wordId int64) (*model.Word, error)
}

type Game interface {
	CreateGame(ctx context.Context, userId int64, lang string) (int64, error)
	GetAllGames(ctx context.Context, userId int64) ([]model.Game, error)
	GetGameById(ctx context.Context, userId, gameId int64) (*model.Game, error)
	UpdateGameStatus(ctx context.Context, gameId int64, status string) error
	MakeGuess(ctx context.Context, userId, gameId int64, guess string) (int, error)
	GetAllGuessByGame(ctx context.Context, userId, gameId int64) ([]model.Guess, error)
}

type Stats interface {
	InitializeStats(ctx context.Context, q repository.Querier, userId int64) error
	UpdateGameEndStats(ctx context.Context, q repository.Querier, update model.StatisticsUpdate) error
	GetStatistics(ctx context.Context, userId int64) (*model.UserStatistics, error)
	GetLanguageStats(ctx context.Context, userId int64, language string) (*model.UserLanguageStats, error)
	GetAllLanguageStats(ctx context.Context, userId int64) ([]model.UserLanguageStats, error)
}

type Profile interface {
	GetProfile(ctx context.Context, userId int64) (*model.UserProfile, error)
	GetProfileSummary(ctx context.Context, userId int64) (*model.UserProfileSummary, error)
	GetProfileByUsername(ctx context.Context, username string) (*model.UserProfileSummary, error)
	UpdateProfile(ctx context.Context, userId int64, input model.UpdateProfileInput) error
	GetSettings(ctx context.Context, userId int64) (*model.UserSettings, error)
	UpdateSettings(ctx context.Context, userId int64, input model.UpdateSettingsInput) error
	GetStatistics(ctx context.Context, userId int64) (*model.UserStatistics, error)
	GetLanguageStats(ctx context.Context, userId int64, language string) (*model.UserLanguageStats, error)
	GetAllLanguageStats(ctx context.Context, userId int64) ([]model.UserLanguageStats, error)
}

type Leaderboard interface {
	GetLeaderboard(ctx context.Context, query model.LeaderboardQuery) (*model.LeaderboardResponse, error)
	GetLeaderboardWithUserRank(ctx context.Context, userId int64, query model.LeaderboardQuery) (*model.LeaderboardResponse, error)
	GetUserRank(ctx context.Context, userId int64, query model.LeaderboardQuery) (*int, error)
	GetAllUserRanks(ctx context.Context, userId int64) (*model.UserRanks, error)
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
