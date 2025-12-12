package mocks

import (
	"context"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type MockAuthRepository struct {
	CreateUserFn        func(ctx context.Context, q repository.Querier, user model.User) (int64, error)
	GetUserByUsernameFn func(ctx context.Context, username string) (model.User, error)
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, q repository.Querier, user model.User) (int64, error) {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, q, user)
	}
	return 0, nil
}

func (m *MockAuthRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	if m.GetUserByUsernameFn != nil {
		return m.GetUserByUsernameFn(ctx, username)
	}
	return model.User{}, nil
}
