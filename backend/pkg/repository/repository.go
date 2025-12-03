package repository

import (
	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUserTx(tx *sqlx.Tx, user model.User) (int64, error)
	GetUser(username, password string) (model.User, error)
}

type Game interface {
	CreateGame(game *model.Game) (int64, error)
	GetAllGames(userId int64) ([]model.Game, error)
	GetGameById(gameId int64) (*model.Game, error)
	GetGameByIdForUpdate(tx *sqlx.Tx, gameId int64) (*model.Game, error)
	UpdateGame(gameId int64) (*model.Game, error)
	UpdateGameStatus(gameId int64, status string) error
	UpdateGameStatusTx(tx *sqlx.Tx, gameId int64, status string) error
	DeleteGame(gameId int64) (*model.Game, error)
	ExpireGames() (int64, error)
	GetExpiredGames() ([]model.Game, error)
	CreateGuess(guess *model.Guess) error
	CreateGuessTx(tx *sqlx.Tx, guess *model.Guess) error
	GetAllGuessByGame(gameId int64) ([]model.Guess, error)
	GetGuessById(guessId int64) error
	CountGuessesByGameTx(tx *sqlx.Tx, gameId int64) (int, error)
	GuessExists(gameId int64, guessWord string) (bool, error)
	GuessExistsTx(tx *sqlx.Tx, gameId int64, guessWord string) (bool, error)
}

type Word interface {
	GetWordById(wordId int64) (*model.Word, error)
	GetRandomWordByLanguage(language string) (*model.Word, error)
	GetWordCountByLanguage(language string) (int, error)
}

type Profile interface {
	GetProfile(userId int64) (*model.UserProfile, error)
	CreateProfile(tx *sqlx.Tx, profile *model.UserProfile) error
	UpdateProfile(userId int64, input model.UpdateProfileInput) error
	GetSettings(userId int64) (*model.UserSettings, error)
	CreateSettings(tx *sqlx.Tx, userId int64) error
	UpdateSettings(userId int64, input model.UpdateSettingsInput) error
	GetProfileSummary(userId int64) (*model.UserProfileSummary, error)
	GetProfileByUsername(username string) (*model.UserProfileSummary, error)
}

type Statistics interface {
	GetStatistics(userId int64) (*model.UserStatistics, error)
	CreateGlobalStreaks(tx *sqlx.Tx, userId int64) error
	UpdateGlobalStreaksAfterGame(tx *sqlx.Tx, update model.StatisticsUpdate) error
	GetLanguageStats(userId int64, language string) (*model.UserLanguageStats, error)
	GetAllLanguageStats(userId int64) ([]model.UserLanguageStats, error)
	UpdateLanguageStats(tx *sqlx.Tx, update model.StatisticsUpdate) error
}

type Leaderboard interface {
	GetLeaderboard(query model.LeaderboardQuery) ([]model.LeaderboardUser, error)
	GetLeaderboardCount(query model.LeaderboardQuery) (int, error)
	GetUserRank(userId int64, query model.LeaderboardQuery) (*int, error)
	GetAllUserRanks(userId int64) (*model.UserRanks, error)
}

type Repository struct {
	Authorization
	Game
	Word
	Profile
	Statistics
	Leaderboard
	TxManager *TransactionManager
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Game:          NewGamePostgres(db),
		Word:          NewWordPostgres(db),
		Profile:       NewProfilePostgres(db),
		Statistics:    NewStatisticsPostgres(db),
		Leaderboard:   NewLeaderboardPostgres(db),
		TxManager:     NewTransactionManager(db),
	}
}
