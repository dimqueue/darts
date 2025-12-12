package mocks

import (
	"context"

	"github.com/dimqueue/darts/pkg/model"
)

type MockAuthService struct {
	CreateUserFn    func(ctx context.Context, user model.User) (int64, error)
	GenerateTokenFn func(ctx context.Context, username, password string) (string, error)
	ParseTokenFn    func(token string) (int64, error)
}

func (m *MockAuthService) CreateUser(ctx context.Context, user model.User) (int64, error) {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, user)
	}
	return 1, nil
}

func (m *MockAuthService) GenerateToken(ctx context.Context, username, password string) (string, error) {
	if m.GenerateTokenFn != nil {
		return m.GenerateTokenFn(ctx, username, password)
	}
	return "mock-token", nil
}

func (m *MockAuthService) ParseToken(token string) (int64, error) {
	if m.ParseTokenFn != nil {
		return m.ParseTokenFn(token)
	}
	return 1, nil
}
