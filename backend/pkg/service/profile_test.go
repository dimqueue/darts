package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository/mocks"
)

func setupProfileService(profileRepo *mocks.MockProfileRepository, statsRepo *mocks.MockStatisticsRepository) *ProfileService {
	statsService := &StatsService{statsRepo: statsRepo}
	return &ProfileService{
		profileRepo:  profileRepo,
		statsService: statsService,
	}
}

func TestGetProfile_Success(t *testing.T) {
	now := time.Now()
	mockProfileRepo := &mocks.MockProfileRepository{
		GetProfileFn: func(ctx context.Context, userId int64) (*model.UserProfile, error) {
			return &model.UserProfile{
				UserId:    userId,
				Timezone:  "UTC",
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	profile, err := service.GetProfile(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if profile.UserId != 42 {
		t.Errorf("expected userId 42, got: %d", profile.UserId)
	}
	if profile.Timezone != "UTC" {
		t.Errorf("expected timezone 'UTC', got: %s", profile.Timezone)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	mockProfileRepo := &mocks.MockProfileRepository{
		GetProfileFn: func(ctx context.Context, userId int64) (*model.UserProfile, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	_, err := service.GetProfile(context.Background(), 42)
	if !errors.Is(err, ErrProfileNotFound) {
		t.Errorf("expected ErrProfileNotFound, got: %v", err)
	}
}

func TestGetProfile_RepoError(t *testing.T) {
	mockProfileRepo := &mocks.MockProfileRepository{
		GetProfileFn: func(ctx context.Context, userId int64) (*model.UserProfile, error) {
			return nil, errors.New("database error")
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	_, err := service.GetProfile(context.Background(), 42)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetProfileSummary_Success(t *testing.T) {
	now := time.Now()
	mockProfileRepo := &mocks.MockProfileRepository{
		GetProfileSummaryFn: func(ctx context.Context, userId int64) (*model.UserProfileSummary, error) {
			return &model.UserProfileSummary{
				Id:                userId,
				Username:          "testuser",
				Name:              "Test User",
				MemberSince:       now,
				TotalGames:        50,
				TotalWins:         35,
				CurrentWinStreak:  5,
				BestWinStreak:     10,
				ShowProfilePublic: true,
			}, nil
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	summary, err := service.GetProfileSummary(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.Username != "testuser" {
		t.Errorf("expected username 'testuser', got: %s", summary.Username)
	}
	if summary.TotalGames != 50 {
		t.Errorf("expected 50 total games, got: %d", summary.TotalGames)
	}
}

func TestGetProfileSummary_NotFound(t *testing.T) {
	mockProfileRepo := &mocks.MockProfileRepository{
		GetProfileSummaryFn: func(ctx context.Context, userId int64) (*model.UserProfileSummary, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	_, err := service.GetProfileSummary(context.Background(), 42)
	if !errors.Is(err, ErrProfileNotFound) {
		t.Errorf("expected ErrProfileNotFound, got: %v", err)
	}
}

func TestGetProfileByUsername_Success(t *testing.T) {
	mockProfileRepo := &mocks.MockProfileRepository{
		GetProfileByUsernameFn: func(ctx context.Context, username string) (*model.UserProfileSummary, error) {
			if username != "testuser" {
				t.Errorf("expected username 'testuser', got: %s", username)
			}
			return &model.UserProfileSummary{
				Id:       42,
				Username: username,
				Name:     "Test User",
			}, nil
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	profile, err := service.GetProfileByUsername(context.Background(), "testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if profile.Id != 42 {
		t.Errorf("expected id 42, got: %d", profile.Id)
	}
}

func TestGetProfileByUsername_NotFound(t *testing.T) {
	mockProfileRepo := &mocks.MockProfileRepository{
		GetProfileByUsernameFn: func(ctx context.Context, username string) (*model.UserProfileSummary, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	_, err := service.GetProfileByUsername(context.Background(), "nonexistent")
	if !errors.Is(err, ErrProfileNotFound) {
		t.Errorf("expected ErrProfileNotFound, got: %v", err)
	}
}

func TestUpdateProfile_Success(t *testing.T) {
	mockProfileRepo := &mocks.MockProfileRepository{
		GetProfileFn: func(ctx context.Context, userId int64) (*model.UserProfile, error) {
			return &model.UserProfile{UserId: userId}, nil
		},
		UpdateProfileFn: func(ctx context.Context, userId int64, input model.UpdateProfileInput) error {
			return nil
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	bio := "New bio"
	input := model.UpdateProfileInput{Bio: &bio}

	err := service.UpdateProfile(context.Background(), 42, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateProfile_NotFound(t *testing.T) {
	mockProfileRepo := &mocks.MockProfileRepository{
		GetProfileFn: func(ctx context.Context, userId int64) (*model.UserProfile, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	bio := "New bio"
	input := model.UpdateProfileInput{Bio: &bio}

	err := service.UpdateProfile(context.Background(), 42, input)
	if !errors.Is(err, ErrProfileNotFound) {
		t.Errorf("expected ErrProfileNotFound, got: %v", err)
	}
}

func TestUpdateProfile_UpdateError(t *testing.T) {
	mockProfileRepo := &mocks.MockProfileRepository{
		GetProfileFn: func(ctx context.Context, userId int64) (*model.UserProfile, error) {
			return &model.UserProfile{UserId: userId}, nil
		},
		UpdateProfileFn: func(ctx context.Context, userId int64, input model.UpdateProfileInput) error {
			return errors.New("update failed")
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	bio := "New bio"
	input := model.UpdateProfileInput{Bio: &bio}

	err := service.UpdateProfile(context.Background(), 42, input)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetSettings_Success(t *testing.T) {
	mockProfileRepo := &mocks.MockProfileRepository{
		GetSettingsFn: func(ctx context.Context, userId int64) (*model.UserSettings, error) {
			return &model.UserSettings{
				UserId:            userId,
				PreferredLanguage: "en",
				Theme:             "dark",
				SoundEnabled:      true,
			}, nil
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	settings, err := service.GetSettings(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if settings.PreferredLanguage != "en" {
		t.Errorf("expected language 'en', got: %s", settings.PreferredLanguage)
	}
	if settings.Theme != "dark" {
		t.Errorf("expected theme 'dark', got: %s", settings.Theme)
	}
}

func TestGetSettings_NotFound(t *testing.T) {
	mockProfileRepo := &mocks.MockProfileRepository{
		GetSettingsFn: func(ctx context.Context, userId int64) (*model.UserSettings, error) {
			return nil, sql.ErrNoRows
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	_, err := service.GetSettings(context.Background(), 42)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestUpdateSettings_Success(t *testing.T) {
	mockProfileRepo := &mocks.MockProfileRepository{
		UpdateSettingsFn: func(ctx context.Context, userId int64, input model.UpdateSettingsInput) error {
			return nil
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	theme := "light"
	input := model.UpdateSettingsInput{Theme: &theme}

	err := service.UpdateSettings(context.Background(), 42, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateSettings_Error(t *testing.T) {
	mockProfileRepo := &mocks.MockProfileRepository{
		UpdateSettingsFn: func(ctx context.Context, userId int64, input model.UpdateSettingsInput) error {
			return errors.New("update failed")
		},
	}

	service := setupProfileService(mockProfileRepo, &mocks.MockStatisticsRepository{})

	theme := "light"
	input := model.UpdateSettingsInput{Theme: &theme}

	err := service.UpdateSettings(context.Background(), 42, input)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestProfileService_GetStatistics(t *testing.T) {
	mockStatsRepo := &mocks.MockStatisticsRepository{
		GetStatisticsFn: func(ctx context.Context, userId int64) (*model.UserStatistics, error) {
			return &model.UserStatistics{
				UserId:           userId,
				TotalGames:       100,
				TotalWins:        75,
				CurrentWinStreak: 5,
				BestWinStreak:    15,
				TotalScore:       5000,
			}, nil
		},
	}

	service := setupProfileService(&mocks.MockProfileRepository{}, mockStatsRepo)

	stats, err := service.GetStatistics(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.TotalGames != 100 {
		t.Errorf("expected 100 total games, got: %d", stats.TotalGames)
	}
	if stats.TotalWins != 75 {
		t.Errorf("expected 75 wins, got: %d", stats.TotalWins)
	}
}

func TestProfileService_GetLanguageStats(t *testing.T) {
	t.Skip("requires txManager setup for DB access")
}

func TestProfileService_GetAllLanguageStats(t *testing.T) {
	mockStatsRepo := &mocks.MockStatisticsRepository{
		GetAllLanguageStatsFn: func(ctx context.Context, userId int64) ([]model.UserLanguageStats, error) {
			return []model.UserLanguageStats{
				{UserId: userId, Language: "en", GamesPlayed: 30},
				{UserId: userId, Language: "ua", GamesPlayed: 20},
			}, nil
		},
	}

	service := setupProfileService(&mocks.MockProfileRepository{}, mockStatsRepo)

	stats, err := service.GetAllLanguageStats(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stats) != 2 {
		t.Errorf("expected 2 language stats, got: %d", len(stats))
	}
}
