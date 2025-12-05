package service

import (
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

func (s *ProfileService) GetProfile(userId int64) (*model.UserProfile, error) {
	return s.profileRepo.GetProfile(userId)
}

func (s *ProfileService) GetProfileSummary(userId int64) (*model.UserProfileSummary, error) {
	return s.profileRepo.GetProfileSummary(userId)
}

func (s *ProfileService) GetProfileByUsername(username string) (*model.UserProfileSummary, error) {
	return s.profileRepo.GetProfileByUsername(username)
}

func (s *ProfileService) UpdateProfile(userId int64, input model.UpdateProfileInput) error {
	_, err := s.profileRepo.GetProfile(userId)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	return s.profileRepo.UpdateProfile(userId, input)
}

func (s *ProfileService) GetSettings(userId int64) (*model.UserSettings, error) {
	return s.profileRepo.GetSettings(userId)
}

func (s *ProfileService) UpdateSettings(userId int64, input model.UpdateSettingsInput) error {
	return s.profileRepo.UpdateSettings(userId, input)
}

func (s *ProfileService) GetStatistics(userId int64) (*model.UserStatistics, error) {
	return s.statsService.GetStatistics(userId)
}

func (s *ProfileService) GetLanguageStats(userId int64, language string) (*model.UserLanguageStats, error) {
	return s.statsService.GetLanguageStats(userId, language)
}

func (s *ProfileService) GetAllLanguageStats(userId int64) ([]model.UserLanguageStats, error) {
	return s.statsService.GetAllLanguageStats(userId)
}
