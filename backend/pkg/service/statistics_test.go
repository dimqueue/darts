package service

import (
	"testing"

	"github.com/dimqueue/darts/pkg/model"
)

func TestCalculateNewStreaks_FirstWin(t *testing.T) {
	service := &StatsService{}

	result := service.CalculateNewStreaks(nil, 42, true)

	if result.CurrentStreak != 1 {
		t.Errorf("expected current streak 1, got: %d", result.CurrentStreak)
	}
	if result.BestStreak != 1 {
		t.Errorf("expected best streak 1, got: %d", result.BestStreak)
	}
	if result.UserId != 42 {
		t.Errorf("expected userId 42, got: %d", result.UserId)
	}
}

func TestCalculateNewStreaks_ContinueWinStreak(t *testing.T) {
	service := &StatsService{}
	current := &model.UserGlobalStreaks{
		UserId:        42,
		CurrentStreak: 3,
		BestStreak:    5,
	}

	result := service.CalculateNewStreaks(current, 42, true)

	if result.CurrentStreak != 4 {
		t.Errorf("expected current streak 4, got: %d", result.CurrentStreak)
	}
	if result.BestStreak != 5 {
		t.Errorf("expected best streak 5, got: %d", result.BestStreak)
	}
}

func TestCalculateNewStreaks_NewBestStreak(t *testing.T) {
	service := &StatsService{}
	current := &model.UserGlobalStreaks{
		UserId:        42,
		CurrentStreak: 5,
		BestStreak:    5,
	}

	result := service.CalculateNewStreaks(current, 42, true)

	if result.CurrentStreak != 6 {
		t.Errorf("expected current streak 6, got: %d", result.CurrentStreak)
	}
	if result.BestStreak != 6 {
		t.Errorf("expected best streak 6, got: %d", result.BestStreak)
	}
}

func TestCalculateNewStreaks_LossResetsStreak(t *testing.T) {
	service := &StatsService{}
	current := &model.UserGlobalStreaks{
		UserId:        42,
		CurrentStreak: 5,
		BestStreak:    10,
	}

	result := service.CalculateNewStreaks(current, 42, false)

	if result.CurrentStreak != 0 {
		t.Errorf("expected current streak 0, got: %d", result.CurrentStreak)
	}
	if result.BestStreak != 10 {
		t.Errorf("expected best streak 10, got: %d", result.BestStreak)
	}
}

func TestCalculateNewStreaks_FirstLoss(t *testing.T) {
	service := &StatsService{}

	result := service.CalculateNewStreaks(nil, 42, false)

	if result.CurrentStreak != 0 {
		t.Errorf("expected current streak 0, got: %d", result.CurrentStreak)
	}
	if result.BestStreak != 0 {
		t.Errorf("expected best streak 0, got: %d", result.BestStreak)
	}
}

func TestCalculateNewLanguageStats_FirstGame_Win(t *testing.T) {
	service := &StatsService{}
	timeSeconds := 120

	result := service.CalculateNewLanguageStats(nil, 42, "en", true, 5, &timeSeconds, 100)

	if result.GamesPlayed != 1 {
		t.Errorf("expected games played 1, got: %d", result.GamesPlayed)
	}
	if result.GamesWon != 1 {
		t.Errorf("expected games won 1, got: %d", result.GamesWon)
	}
	if result.TotalGuesses != 5 {
		t.Errorf("expected total guesses 5, got: %d", result.TotalGuesses)
	}
	if result.CurrentStreak != 1 {
		t.Errorf("expected current streak 1, got: %d", result.CurrentStreak)
	}
	if result.BestStreak != 1 {
		t.Errorf("expected best streak 1, got: %d", result.BestStreak)
	}
	if result.TotalScore != 100 {
		t.Errorf("expected total score 100, got: %d", result.TotalScore)
	}
	if result.FastestWinSeconds == nil || *result.FastestWinSeconds != 120 {
		t.Errorf("expected fastest win 120s, got: %v", result.FastestWinSeconds)
	}
	if result.FewestGuessesWin == nil || *result.FewestGuessesWin != 5 {
		t.Errorf("expected fewest guesses 5, got: %v", result.FewestGuessesWin)
	}
}

func TestCalculateNewLanguageStats_FirstGame_Loss(t *testing.T) {
	service := &StatsService{}

	result := service.CalculateNewLanguageStats(nil, 42, "en", false, 10, nil, 0)

	if result.GamesPlayed != 1 {
		t.Errorf("expected games played 1, got: %d", result.GamesPlayed)
	}
	if result.GamesWon != 0 {
		t.Errorf("expected games won 0, got: %d", result.GamesWon)
	}
	if result.CurrentStreak != 0 {
		t.Errorf("expected current streak 0, got: %d", result.CurrentStreak)
	}
	if result.FastestWinSeconds != nil {
		t.Errorf("expected fastest win nil, got: %v", result.FastestWinSeconds)
	}
	if result.FewestGuessesWin != nil {
		t.Errorf("expected fewest guesses nil, got: %v", result.FewestGuessesWin)
	}
}

func TestCalculateNewLanguageStats_ContinueWinStreak(t *testing.T) {
	service := &StatsService{}
	current := &model.UserLanguageStats{
		UserId:        42,
		Language:      "en",
		GamesPlayed:   10,
		GamesWon:      8,
		TotalGuesses:  50,
		CurrentStreak: 3,
		BestStreak:    5,
		TotalScore:    800,
	}
	timeSeconds := 90

	result := service.CalculateNewLanguageStats(current, 42, "en", true, 4, &timeSeconds, 100)

	if result.GamesPlayed != 11 {
		t.Errorf("expected games played 11, got: %d", result.GamesPlayed)
	}
	if result.GamesWon != 9 {
		t.Errorf("expected games won 9, got: %d", result.GamesWon)
	}
	if result.TotalGuesses != 54 {
		t.Errorf("expected total guesses 54, got: %d", result.TotalGuesses)
	}
	if result.CurrentStreak != 4 {
		t.Errorf("expected current streak 4, got: %d", result.CurrentStreak)
	}
	if result.BestStreak != 5 {
		t.Errorf("expected best streak 5, got: %d", result.BestStreak)
	}
	if result.TotalScore != 900 {
		t.Errorf("expected total score 900, got: %d", result.TotalScore)
	}
}

func TestCalculateNewLanguageStats_NewBestStreak(t *testing.T) {
	service := &StatsService{}
	current := &model.UserLanguageStats{
		UserId:        42,
		Language:      "en",
		GamesPlayed:   5,
		GamesWon:      5,
		CurrentStreak: 5,
		BestStreak:    5,
	}

	result := service.CalculateNewLanguageStats(current, 42, "en", true, 3, nil, 100)

	if result.CurrentStreak != 6 {
		t.Errorf("expected current streak 6, got: %d", result.CurrentStreak)
	}
	if result.BestStreak != 6 {
		t.Errorf("expected best streak 6, got: %d", result.BestStreak)
	}
}

func TestCalculateNewLanguageStats_LossResetsStreak(t *testing.T) {
	service := &StatsService{}
	current := &model.UserLanguageStats{
		UserId:        42,
		Language:      "en",
		GamesPlayed:   10,
		GamesWon:      8,
		CurrentStreak: 5,
		BestStreak:    7,
	}

	result := service.CalculateNewLanguageStats(current, 42, "en", false, 10, nil, 0)

	if result.GamesPlayed != 11 {
		t.Errorf("expected games played 11, got: %d", result.GamesPlayed)
	}
	if result.GamesWon != 8 {
		t.Errorf("expected games won 8, got: %d", result.GamesWon)
	}
	if result.CurrentStreak != 0 {
		t.Errorf("expected current streak 0, got: %d", result.CurrentStreak)
	}
	if result.BestStreak != 7 {
		t.Errorf("expected best streak 7, got: %d", result.BestStreak)
	}
}

func TestCalculateNewLanguageStats_NewFastestWin(t *testing.T) {
	service := &StatsService{}
	oldFastestWin := 100
	current := &model.UserLanguageStats{
		UserId:            42,
		Language:          "en",
		GamesPlayed:       5,
		GamesWon:          4,
		FastestWinSeconds: &oldFastestWin,
	}
	newTimeSeconds := 80

	result := service.CalculateNewLanguageStats(current, 42, "en", true, 3, &newTimeSeconds, 100)

	if result.FastestWinSeconds == nil || *result.FastestWinSeconds != 80 {
		t.Errorf("expected fastest win 80s, got: %v", result.FastestWinSeconds)
	}
}

func TestCalculateNewLanguageStats_KeepOldFastestWin(t *testing.T) {
	service := &StatsService{}
	oldFastestWin := 60
	current := &model.UserLanguageStats{
		UserId:            42,
		Language:          "en",
		GamesPlayed:       5,
		GamesWon:          4,
		FastestWinSeconds: &oldFastestWin,
	}
	newTimeSeconds := 90

	result := service.CalculateNewLanguageStats(current, 42, "en", true, 3, &newTimeSeconds, 100)

	if result.FastestWinSeconds == nil || *result.FastestWinSeconds != 60 {
		t.Errorf("expected fastest win 60s (unchanged), got: %v", result.FastestWinSeconds)
	}
}

func TestCalculateNewLanguageStats_NewFewestGuesses(t *testing.T) {
	service := &StatsService{}
	oldFewest := 5
	current := &model.UserLanguageStats{
		UserId:           42,
		Language:         "en",
		GamesPlayed:      5,
		GamesWon:         4,
		FewestGuessesWin: &oldFewest,
	}

	result := service.CalculateNewLanguageStats(current, 42, "en", true, 3, nil, 100)

	if result.FewestGuessesWin == nil || *result.FewestGuessesWin != 3 {
		t.Errorf("expected fewest guesses 3, got: %v", result.FewestGuessesWin)
	}
}

func TestCalculateNewLanguageStats_KeepOldFewestGuesses(t *testing.T) {
	service := &StatsService{}
	oldFewest := 2
	current := &model.UserLanguageStats{
		UserId:           42,
		Language:         "en",
		GamesPlayed:      5,
		GamesWon:         4,
		FewestGuessesWin: &oldFewest,
	}

	result := service.CalculateNewLanguageStats(current, 42, "en", true, 4, nil, 100)

	if result.FewestGuessesWin == nil || *result.FewestGuessesWin != 2 {
		t.Errorf("expected fewest guesses 2 (unchanged), got: %v", result.FewestGuessesWin)
	}
}
