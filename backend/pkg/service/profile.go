package service

import (
	"context"
	"fmt"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type ProfileService struct {
	profileRepo  repository.Profile
	statsService *StatsService
}

func NewProfileService(profileRepo repository.Profile, statsService *StatsService) *ProfileService {
	return &ProfileService{
		profileRepo:  profileRepo,
		statsService: statsService,
	}
}

func (s *ProfileService) GetProfile(ctx context.Context, userId int64) (*model.UserProfile, error) {
	return s.profileRepo.GetProfile(ctx, userId)
}

func (s *ProfileService) GetProfileSummary(ctx context.Context, userId int64) (*model.UserProfileSummary, error) {
	return s.profileRepo.GetProfileSummary(ctx, userId)
}

func (s *ProfileService) GetProfileByUsername(ctx context.Context, username string) (*model.UserProfileSummary, error) {
	return s.profileRepo.GetProfileByUsername(ctx, username)
}

func (s *ProfileService) UpdateProfile(ctx context.Context, userId int64, input model.UpdateProfileInput) error {
	_, err := s.profileRepo.GetProfile(ctx, userId)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	return s.profileRepo.UpdateProfile(ctx, userId, input)
}

func (s *ProfileService) GetSettings(ctx context.Context, userId int64) (*model.UserSettings, error) {
	return s.profileRepo.GetSettings(ctx, userId)
}

func (s *ProfileService) UpdateSettings(ctx context.Context, userId int64, input model.UpdateSettingsInput) error {
	return s.profileRepo.UpdateSettings(ctx, userId, input)
}

func (s *ProfileService) GetStatistics(ctx context.Context, userId int64) (*model.UserStatistics, error) {
	return s.statsService.GetStatistics(ctx, userId)
}

func (s *ProfileService) GetLanguageStats(ctx context.Context, userId int64, language string) (*model.UserLanguageStats, error) {
	return s.statsService.GetLanguageStats(ctx, userId, language)
}

func (s *ProfileService) GetAllLanguageStats(ctx context.Context, userId int64) ([]model.UserLanguageStats, error) {
	return s.statsService.GetAllLanguageStats(ctx, userId)
}
