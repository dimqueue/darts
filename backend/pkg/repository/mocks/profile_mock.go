package mocks

import (
	"context"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type MockProfileRepository struct {
	GetProfileFn           func(ctx context.Context, userId int64) (*model.UserProfile, error)
	CreateProfileFn        func(ctx context.Context, q repository.Querier, profile *model.UserProfile) error
	UpdateProfileFn        func(ctx context.Context, userId int64, input model.UpdateProfileInput) error
	GetSettingsFn          func(ctx context.Context, userId int64) (*model.UserSettings, error)
	CreateSettingsFn       func(ctx context.Context, q repository.Querier, userId int64) error
	UpdateSettingsFn       func(ctx context.Context, userId int64, input model.UpdateSettingsInput) error
	GetProfileSummaryFn    func(ctx context.Context, userId int64) (*model.UserProfileSummary, error)
	GetProfileByUsernameFn func(ctx context.Context, username string) (*model.UserProfileSummary, error)
}

func (m *MockProfileRepository) GetProfile(ctx context.Context, userId int64) (*model.UserProfile, error) {
	if m.GetProfileFn != nil {
		return m.GetProfileFn(ctx, userId)
	}
	return nil, nil
}

func (m *MockProfileRepository) CreateProfile(ctx context.Context, q repository.Querier, profile *model.UserProfile) error {
	if m.CreateProfileFn != nil {
		return m.CreateProfileFn(ctx, q, profile)
	}
	return nil
}

func (m *MockProfileRepository) UpdateProfile(ctx context.Context, userId int64, input model.UpdateProfileInput) error {
	if m.UpdateProfileFn != nil {
		return m.UpdateProfileFn(ctx, userId, input)
	}
	return nil
}

func (m *MockProfileRepository) GetSettings(ctx context.Context, userId int64) (*model.UserSettings, error) {
	if m.GetSettingsFn != nil {
		return m.GetSettingsFn(ctx, userId)
	}
	return nil, nil
}

func (m *MockProfileRepository) CreateSettings(ctx context.Context, q repository.Querier, userId int64) error {
	if m.CreateSettingsFn != nil {
		return m.CreateSettingsFn(ctx, q, userId)
	}
	return nil
}

func (m *MockProfileRepository) UpdateSettings(ctx context.Context, userId int64, input model.UpdateSettingsInput) error {
	if m.UpdateSettingsFn != nil {
		return m.UpdateSettingsFn(ctx, userId, input)
	}
	return nil
}

func (m *MockProfileRepository) GetProfileSummary(ctx context.Context, userId int64) (*model.UserProfileSummary, error) {
	if m.GetProfileSummaryFn != nil {
		return m.GetProfileSummaryFn(ctx, userId)
	}
	return nil, nil
}

func (m *MockProfileRepository) GetProfileByUsername(ctx context.Context, username string) (*model.UserProfileSummary, error) {
	if m.GetProfileByUsernameFn != nil {
		return m.GetProfileByUsernameFn(ctx, username)
	}
	return nil, nil
}
