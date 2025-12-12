package mocks

import (
	"context"

	"github.com/dimqueue/darts/pkg/model"
)

type MockProfileService struct {
	GetProfileFn           func(ctx context.Context, userId int64) (*model.UserProfile, error)
	GetProfileSummaryFn    func(ctx context.Context, userId int64) (*model.UserProfileSummary, error)
	GetProfileByUsernameFn func(ctx context.Context, username string) (*model.UserProfileSummary, error)
	UpdateProfileFn        func(ctx context.Context, userId int64, input model.UpdateProfileInput) error
	GetSettingsFn          func(ctx context.Context, userId int64) (*model.UserSettings, error)
	UpdateSettingsFn       func(ctx context.Context, userId int64, input model.UpdateSettingsInput) error
	GetStatisticsFn        func(ctx context.Context, userId int64) (*model.UserStatistics, error)
	GetLanguageStatsFn     func(ctx context.Context, userId int64, language string) (*model.UserLanguageStats, error)
	GetAllLanguageStatsFn  func(ctx context.Context, userId int64) ([]model.UserLanguageStats, error)
}

func (m *MockProfileService) GetProfile(ctx context.Context, userId int64) (*model.UserProfile, error) {
	if m.GetProfileFn != nil {
		return m.GetProfileFn(ctx, userId)
	}
	return &model.UserProfile{UserId: userId}, nil
}

func (m *MockProfileService) GetProfileSummary(ctx context.Context, userId int64) (*model.UserProfileSummary, error) {
	if m.GetProfileSummaryFn != nil {
		return m.GetProfileSummaryFn(ctx, userId)
	}
	return &model.UserProfileSummary{Id: userId}, nil
}

func (m *MockProfileService) GetProfileByUsername(ctx context.Context, username string) (*model.UserProfileSummary, error) {
	if m.GetProfileByUsernameFn != nil {
		return m.GetProfileByUsernameFn(ctx, username)
	}
	return &model.UserProfileSummary{Username: username}, nil
}

func (m *MockProfileService) UpdateProfile(ctx context.Context, userId int64, input model.UpdateProfileInput) error {
	if m.UpdateProfileFn != nil {
		return m.UpdateProfileFn(ctx, userId, input)
	}
	return nil
}

func (m *MockProfileService) GetSettings(ctx context.Context, userId int64) (*model.UserSettings, error) {
	if m.GetSettingsFn != nil {
		return m.GetSettingsFn(ctx, userId)
	}
	return &model.UserSettings{UserId: userId}, nil
}

func (m *MockProfileService) UpdateSettings(ctx context.Context, userId int64, input model.UpdateSettingsInput) error {
	if m.UpdateSettingsFn != nil {
		return m.UpdateSettingsFn(ctx, userId, input)
	}
	return nil
}

func (m *MockProfileService) GetStatistics(ctx context.Context, userId int64) (*model.UserStatistics, error) {
	if m.GetStatisticsFn != nil {
		return m.GetStatisticsFn(ctx, userId)
	}
	return &model.UserStatistics{UserId: userId}, nil
}

func (m *MockProfileService) GetLanguageStats(ctx context.Context, userId int64, language string) (*model.UserLanguageStats, error) {
	if m.GetLanguageStatsFn != nil {
		return m.GetLanguageStatsFn(ctx, userId, language)
	}
	return &model.UserLanguageStats{UserId: userId, Language: language}, nil
}

func (m *MockProfileService) GetAllLanguageStats(ctx context.Context, userId int64) ([]model.UserLanguageStats, error) {
	if m.GetAllLanguageStatsFn != nil {
		return m.GetAllLanguageStatsFn(ctx, userId)
	}
	return []model.UserLanguageStats{}, nil
}
