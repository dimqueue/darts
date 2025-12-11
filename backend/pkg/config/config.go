package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Server
var (
	Port = "8080"
)

// Database
var (
	DBHost     = "localhost"
	DBPort     = "5432"
	DBUser     = "admin"
	DBName     = "dartsdb"
	DBSSLMode  = "disable"
	DBPassword string
)

// Auth
var (
	JWTSecret string
	TokenTTL  = 12 * time.Hour
)

// CORS
var (
	CORSOrigins = []string{"http://localhost:5173"}
)

// Word Service
var (
	WordServiceType    = "grpc"
	WordServiceURL     = "localhost:50051"
	WordServiceTimeout = 30
)

// Logger
var (
	LogLevel        = "debug"
	LogFormat       = "text"
	LogOutput       = "stdout"
	LogReportCaller = false
)

// Game
var (
	GameTTL            = 2 * time.Hour
	GameComputeTimeout = 30 * time.Second
	ScorePerWin        = 100
	GameMaxGuesses     = 0 // 0 = unlimited
)

// Leaderboard
var (
	LeaderboardMinGamesGlobal = 3
	LeaderboardMinGamesPeriod = 1
	LeaderboardDefaultLimit   = 50
	LeaderboardMaxLimit       = 100
)

func Load() {
	loadServer()
	loadDatabase()
	loadAuth()
	loadCORS()
	loadWordService()
	loadLogger()
	loadGame()
	loadLeaderboard()
}

func loadServer() {
	Port = getEnvOrViper("APP_PORT", "port", Port)
}

func loadDatabase() {
	DBHost = getEnvOrViper("POSTGRES_HOST", "db.host", DBHost)
	DBPort = getEnvOrViper("POSTGRES_PORT", "db.port", DBPort)
	DBUser = getEnvOrViper("POSTGRES_USER", "db.username", DBUser)
	DBName = getEnvOrViper("POSTGRES_DB", "db.dbname", DBName)
	DBSSLMode = getEnvOrViper("POSTGRES_SSLMODE", "db.sslmode", DBSSLMode)
	DBPassword = os.Getenv("POSTGRES_PASSWORD")
}

func loadAuth() {
	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		panic(fmt.Sprintf("JWT_SECRET environment variable is required"))
	}

	if v := getEnvOrViperInt("TOKEN_TTL_HOURS", "auth.tokenTTLHours", 12); v > 0 {
		TokenTTL = time.Duration(v) * time.Hour
	}
}

func loadCORS() {
	origins := getEnvOrViper("CORS_ORIGINS", "cors.origins", "")
	if origins != "" {
		CORSOrigins = parseCommaSeparated(origins)
	}
}

func loadWordService() {
	WordServiceType = getEnvOrViper("WORD_SERVICE_TYPE", "word-service.type", WordServiceType)
	WordServiceURL = getEnvOrViper("WORD_SERVICE_URL", "word-service.baseURL", WordServiceURL)
	WordServiceTimeout = getEnvOrViperInt("WORD_SERVICE_TIMEOUT", "word-service.timeout", WordServiceTimeout)
}

func loadLogger() {
	LogLevel = getEnvOrViper("LOG_LEVEL", "logger.level", LogLevel)
	LogFormat = getEnvOrViper("LOG_FORMAT", "logger.format", LogFormat)
	LogOutput = getEnvOrViper("LOG_OUTPUT", "logger.output", LogOutput)
	LogReportCaller = viper.GetBool("logger.reportCaller")
}

func loadGame() {
	if v := viper.GetInt("game.ttl"); v > 0 {
		GameTTL = time.Duration(v) * time.Second
	}
	if v := viper.GetInt("game.computeTimeout"); v > 0 {
		GameComputeTimeout = time.Duration(v) * time.Second
	}
	if v := viper.GetInt("game.maxGuesses"); v > 0 {
		GameMaxGuesses = v
	}
}

func loadLeaderboard() {
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

func getEnvOrViper(envKey, viperKey, defaultVal string) string {
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	if v := viper.GetString(viperKey); v != "" {
		return v
	}
	return defaultVal
}

func getEnvOrViperInt(envKey, viperKey string, defaultVal int) int {
	if v := os.Getenv(envKey); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	if v := viper.GetInt(viperKey); v != 0 {
		return v
	}
	return defaultVal
}

func parseCommaSeparated(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
