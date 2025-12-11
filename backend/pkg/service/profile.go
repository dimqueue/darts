package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dimqueue/darts/pkg/logger"
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
	log := logger.Op(ctx, "ProfileService.GetProfile").With(logger.FieldUserID, userId)

	profile, err := s.profileRepo.GetProfile(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		log.Error("failed to get profile", logger.FieldError, err)
		return nil, Logged(err)
	}
	return profile, nil
}

func (s *ProfileService) GetProfileSummary(ctx context.Context, userId int64) (*model.UserProfileSummary, error) {
	log := logger.Op(ctx, "ProfileService.GetProfileSummary").With(logger.FieldUserID, userId)

	summary, err := s.profileRepo.GetProfileSummary(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		log.Error("failed to get profile summary", logger.FieldError, err)
		return nil, Logged(err)
	}
	return summary, nil
}

func (s *ProfileService) GetProfileByUsername(ctx context.Context, username string) (*model.UserProfileSummary, error) {
	log := logger.Op(ctx, "ProfileService.GetProfileByUsername").With(logger.FieldUsername, username)

	profile, err := s.profileRepo.GetProfileByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		log.Error("failed to get profile by username", logger.FieldError, err)
		return nil, Logged(err)
	}
	return profile, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, userId int64, input model.UpdateProfileInput) error {
	log := logger.Op(ctx, "ProfileService.UpdateProfile").With(logger.FieldUserID, userId)

	_, err := s.profileRepo.GetProfile(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProfileNotFound
		}
		log.Error("failed to get profile", logger.FieldError, err)
		return Logged(err)
	}

	if err := s.profileRepo.UpdateProfile(ctx, userId, input); err != nil {
		log.Error("failed to update profile", logger.FieldError, err)
		return Logged(err)
	}
	return nil
}

func (s *ProfileService) GetSettings(ctx context.Context, userId int64) (*model.UserSettings, error) {
	log := logger.Op(ctx, "ProfileService.GetSettings").With(logger.FieldUserID, userId)

	settings, err := s.profileRepo.GetSettings(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		log.Error("failed to get settings", logger.FieldError, err)
		return nil, Logged(err)
	}
	return settings, nil
}

func (s *ProfileService) UpdateSettings(ctx context.Context, userId int64, input model.UpdateSettingsInput) error {
	log := logger.Op(ctx, "ProfileService.UpdateSettings").With(logger.FieldUserID, userId)

	if err := s.profileRepo.UpdateSettings(ctx, userId, input); err != nil {
		log.Error("failed to update settings", logger.FieldError, err)
		return Logged(err)
	}
	return nil
}

func (s *ProfileService) GetStatistics(ctx context.Context, userId int64) (*model.UserStatistics, error) {
	log := logger.Op(ctx, "ProfileService.GetStatistics").With(logger.FieldUserID, userId)

	stats, err := s.statsService.GetStatistics(ctx, userId)
	if err != nil {
		log.Error("failed to get statistics", logger.FieldError, err)
		return nil, Logged(err)
	}
	return stats, nil
}

func (s *ProfileService) GetLanguageStats(ctx context.Context, userId int64, language string) (*model.UserLanguageStats, error) {
	log := logger.Op(ctx, "ProfileService.GetLanguageStats").With(logger.FieldUserID, userId, logger.FieldLanguage, language)

	stats, err := s.statsService.GetLanguageStats(ctx, userId, language)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		log.Error("failed to get language stats", logger.FieldError, err)
		return nil, Logged(err)
	}
	return stats, nil
}

func (s *ProfileService) GetAllLanguageStats(ctx context.Context, userId int64) ([]model.UserLanguageStats, error) {
	log := logger.Op(ctx, "ProfileService.GetAllLanguageStats").With(logger.FieldUserID, userId)

	stats, err := s.statsService.GetAllLanguageStats(ctx, userId)
	if err != nil {
		log.Error("failed to get all language stats", logger.FieldError, err)
		return nil, Logged(err)
	}
	return stats, nil
}
