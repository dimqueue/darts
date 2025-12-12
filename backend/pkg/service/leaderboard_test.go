package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository/mocks"
)

func TestGetLeaderboard_Global_Success(t *testing.T) {
	mockRepo := &mocks.MockLeaderboardRepository{
		GetGlobalLeaderboardFn: func(ctx context.Context, limit, offset int) ([]model.LeaderboardUser, error) {
			return []model.LeaderboardUser{
				{Rank: 1, UserId: 1, Username: "user1", TotalScore: 1000},
				{Rank: 2, UserId: 2, Username: "user2", TotalScore: 900},
			}, nil
		},
		GetGlobalLeaderboardCountFn: func(ctx context.Context) (int, error) {
			return 100, nil
		},
	}

	service := NewLeaderboardService(mockRepo)
	query := model.LeaderboardQuery{Type: "global", Limit: 10, Offset: 0}

	response, err := service.GetLeaderboard(context.Background(), query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.LeaderboardType != "global" {
		t.Errorf("expected type 'global', got: %s", response.LeaderboardType)
	}
	if len(response.Users) != 2 {
		t.Errorf("expected 2 users, got: %d", len(response.Users))
	}
	if response.Total != 100 {
		t.Errorf("expected total 100, got: %d", response.Total)
	}
}

func TestGetLeaderboard_Global_ByLanguage(t *testing.T) {
	mockRepo := &mocks.MockLeaderboardRepository{
		GetGlobalLeaderboardByLanguageFn: func(ctx context.Context, language string, limit, offset int) ([]model.LeaderboardUser, error) {
			if language != "en" {
				t.Errorf("expected language 'en', got: %s", language)
			}
			return []model.LeaderboardUser{
				{Rank: 1, UserId: 1, Username: "user1", TotalScore: 500},
			}, nil
		},
		GetGlobalLeaderboardByLanguageCountFn: func(ctx context.Context, language string) (int, error) {
			return 50, nil
		},
	}

	service := NewLeaderboardService(mockRepo)
	lang := "en"
	query := model.LeaderboardQuery{Type: "global", Language: &lang, Limit: 10, Offset: 0}

	response, err := service.GetLeaderboard(context.Background(), query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if *response.Language != "en" {
		t.Errorf("expected language 'en', got: %v", response.Language)
	}
	if response.Total != 50 {
		t.Errorf("expected total 50, got: %d", response.Total)
	}
}

func TestGetLeaderboard_Period_Weekly(t *testing.T) {
	mockRepo := &mocks.MockLeaderboardRepository{
		GetPeriodLeaderboardFn: func(ctx context.Context, periodStart time.Time, language *string, limit, offset int) ([]model.LeaderboardUser, error) {
			return []model.LeaderboardUser{
				{Rank: 1, UserId: 1, Username: "user1", TotalScore: 200},
			}, nil
		},
		GetPeriodLeaderboardCountFn: func(ctx context.Context, periodStart time.Time, language *string) (int, error) {
			return 25, nil
		},
	}

	service := NewLeaderboardService(mockRepo)
	query := model.LeaderboardQuery{Type: "weekly", Limit: 10, Offset: 0}

	response, err := service.GetLeaderboard(context.Background(), query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.LeaderboardType != "weekly" {
		t.Errorf("expected type 'weekly', got: %s", response.LeaderboardType)
	}
	if len(response.Users) != 1 {
		t.Errorf("expected 1 user, got: %d", len(response.Users))
	}
}

func TestGetLeaderboard_Period_Monthly(t *testing.T) {
	mockRepo := &mocks.MockLeaderboardRepository{
		GetPeriodLeaderboardFn: func(ctx context.Context, periodStart time.Time, language *string, limit, offset int) ([]model.LeaderboardUser, error) {
			return []model.LeaderboardUser{}, nil
		},
		GetPeriodLeaderboardCountFn: func(ctx context.Context, periodStart time.Time, language *string) (int, error) {
			return 0, nil
		},
	}

	service := NewLeaderboardService(mockRepo)
	query := model.LeaderboardQuery{Type: "monthly", Limit: 10, Offset: 0}

	response, err := service.GetLeaderboard(context.Background(), query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.LeaderboardType != "monthly" {
		t.Errorf("expected type 'monthly', got: %s", response.LeaderboardType)
	}
}

func TestGetLeaderboard_Period_Daily(t *testing.T) {
	mockRepo := &mocks.MockLeaderboardRepository{
		GetPeriodLeaderboardFn: func(ctx context.Context, periodStart time.Time, language *string, limit, offset int) ([]model.LeaderboardUser, error) {
			return []model.LeaderboardUser{}, nil
		},
		GetPeriodLeaderboardCountFn: func(ctx context.Context, periodStart time.Time, language *string) (int, error) {
			return 0, nil
		},
	}

	service := NewLeaderboardService(mockRepo)
	query := model.LeaderboardQuery{Type: "daily", Limit: 10, Offset: 0}

	response, err := service.GetLeaderboard(context.Background(), query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.LeaderboardType != "daily" {
		t.Errorf("expected type 'daily', got: %s", response.LeaderboardType)
	}
}

func TestGetLeaderboard_DefaultLimit(t *testing.T) {
	mockRepo := &mocks.MockLeaderboardRepository{
		GetGlobalLeaderboardFn: func(ctx context.Context, limit, offset int) ([]model.LeaderboardUser, error) {
			if limit <= 0 {
				t.Error("limit should be set to default value")
			}
			return []model.LeaderboardUser{}, nil
		},
		GetGlobalLeaderboardCountFn: func(ctx context.Context) (int, error) {
			return 0, nil
		},
	}

	service := NewLeaderboardService(mockRepo)
	query := model.LeaderboardQuery{Type: "global", Limit: 0, Offset: 0}

	_, err := service.GetLeaderboard(context.Background(), query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetLeaderboard_RepoError(t *testing.T) {
	mockRepo := &mocks.MockLeaderboardRepository{
		GetGlobalLeaderboardFn: func(ctx context.Context, limit, offset int) ([]model.LeaderboardUser, error) {
			return nil, errors.New("database error")
		},
	}

	service := NewLeaderboardService(mockRepo)
	query := model.LeaderboardQuery{Type: "global", Limit: 10, Offset: 0}

	_, err := service.GetLeaderboard(context.Background(), query)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetUserRank_Global(t *testing.T) {
	rank := 5
	mockRepo := &mocks.MockLeaderboardRepository{
		GetGlobalUserRankFn: func(ctx context.Context, userId int64) (*int, error) {
			if userId != 42 {
				t.Errorf("expected userId 42, got: %d", userId)
			}
			return &rank, nil
		},
	}

	service := NewLeaderboardService(mockRepo)
	query := model.LeaderboardQuery{Type: "global"}

	result, err := service.GetUserRank(context.Background(), 42, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil || *result != 5 {
		t.Errorf("expected rank 5, got: %v", result)
	}
}

func TestGetUserRank_GlobalByLanguage(t *testing.T) {
	rank := 3
	mockRepo := &mocks.MockLeaderboardRepository{
		GetGlobalUserRankByLanguageFn: func(ctx context.Context, userId int64, language string) (*int, error) {
			if language != "en" {
				t.Errorf("expected language 'en', got: %s", language)
			}
			return &rank, nil
		},
	}

	service := NewLeaderboardService(mockRepo)
	lang := "en"
	query := model.LeaderboardQuery{Type: "global", Language: &lang}

	result, err := service.GetUserRank(context.Background(), 42, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil || *result != 3 {
		t.Errorf("expected rank 3, got: %v", result)
	}
}

func TestGetUserRank_Period(t *testing.T) {
	rank := 10
	mockRepo := &mocks.MockLeaderboardRepository{
		GetPeriodUserRankFn: func(ctx context.Context, userId int64, periodStart time.Time, language *string) (*int, error) {
			return &rank, nil
		},
	}

	service := NewLeaderboardService(mockRepo)
	query := model.LeaderboardQuery{Type: "weekly"}

	result, err := service.GetUserRank(context.Background(), 42, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil || *result != 10 {
		t.Errorf("expected rank 10, got: %v", result)
	}
}

func TestGetUserRank_NotRanked(t *testing.T) {
	mockRepo := &mocks.MockLeaderboardRepository{
		GetGlobalUserRankFn: func(ctx context.Context, userId int64) (*int, error) {
			return nil, nil
		},
	}

	service := NewLeaderboardService(mockRepo)
	query := model.LeaderboardQuery{Type: "global"}

	result, err := service.GetUserRank(context.Background(), 42, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != nil {
		t.Errorf("expected nil rank, got: %v", result)
	}
}

func TestGetLeaderboardWithUserRank_Success(t *testing.T) {
	rank := 5
	mockRepo := &mocks.MockLeaderboardRepository{
		GetGlobalLeaderboardFn: func(ctx context.Context, limit, offset int) ([]model.LeaderboardUser, error) {
			return []model.LeaderboardUser{
				{Rank: 1, UserId: 1, Username: "user1"},
			}, nil
		},
		GetGlobalLeaderboardCountFn: func(ctx context.Context) (int, error) {
			return 100, nil
		},
		GetGlobalUserRankFn: func(ctx context.Context, userId int64) (*int, error) {
			return &rank, nil
		},
	}

	service := NewLeaderboardService(mockRepo)
	query := model.LeaderboardQuery{Type: "global", Limit: 10}

	response, err := service.GetLeaderboardWithUserRank(context.Background(), 42, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.CurrentUserRank == nil || *response.CurrentUserRank != 5 {
		t.Errorf("expected current user rank 5, got: %v", response.CurrentUserRank)
	}
}

func TestGetAllUserRanks_Success(t *testing.T) {
	globalRank := 10
	dailyRank := 5
	weeklyRank := 8
	monthlyRank := 12

	mockRepo := &mocks.MockLeaderboardRepository{
		GetGlobalUserRankFn: func(ctx context.Context, userId int64) (*int, error) {
			return &globalRank, nil
		},
		GetPeriodUserRankFn: func(ctx context.Context, userId int64, periodStart time.Time, language *string) (*int, error) {
			now := time.Now()
			dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

			if periodStart.Equal(dayStart) {
				return &dailyRank, nil
			}
			if periodStart.Equal(monthStart) {
				return &monthlyRank, nil
			}
			return &weeklyRank, nil
		},
	}

	service := NewLeaderboardService(mockRepo)

	ranks, err := service.GetAllUserRanks(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ranks.GlobalRank == nil || *ranks.GlobalRank != 10 {
		t.Errorf("expected global rank 10, got: %v", ranks.GlobalRank)
	}
}

func TestGetAllUserRanks_PartialErrors(t *testing.T) {
	globalRank := 10

	mockRepo := &mocks.MockLeaderboardRepository{
		GetGlobalUserRankFn: func(ctx context.Context, userId int64) (*int, error) {
			return &globalRank, nil
		},
		GetPeriodUserRankFn: func(ctx context.Context, userId int64, periodStart time.Time, language *string) (*int, error) {
			return nil, errors.New("period rank error")
		},
	}

	service := NewLeaderboardService(mockRepo)

	ranks, err := service.GetAllUserRanks(context.Background(), 42)
	if ranks == nil {
		t.Fatal("expected partial ranks, got nil")
	}
	if ranks.GlobalRank == nil || *ranks.GlobalRank != 10 {
		t.Errorf("expected global rank 10, got: %v", ranks.GlobalRank)
	}
	if err == nil {
		t.Error("expected error for partial failures")
	}
}

func TestGetPeriodStart_Daily(t *testing.T) {
	start, err := GetPeriodStart("daily")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	now := time.Now()
	expected := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if !start.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, start)
	}
}

func TestGetPeriodStart_Weekly(t *testing.T) {
	start, err := GetPeriodStart("weekly")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if start.Weekday() != time.Monday {
		t.Errorf("expected Monday, got %v", start.Weekday())
	}
}

func TestGetPeriodStart_Monthly(t *testing.T) {
	start, err := GetPeriodStart("monthly")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if start.Day() != 1 {
		t.Errorf("expected day 1, got %d", start.Day())
	}
}

func TestGetPeriodStart_InvalidType(t *testing.T) {
	_, err := GetPeriodStart("invalid")
	if err == nil {
		t.Error("expected error for invalid period type")
	}
}
