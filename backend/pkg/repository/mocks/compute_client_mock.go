package mocks

import (
	"context"

	"github.com/dimqueue/darts/pkg/connections"
)

type MockComputeClient struct {
	StartGameFn   func(ctx context.Context, req *connections.StartGameRequest) (*connections.StartGameResponse, error)
	MakeGuessFn   func(ctx context.Context, req *connections.GuessRequest) (*connections.GuessResponse, error)
	HealthCheckFn func(ctx context.Context) (*connections.HealthResponse, error)
	CloseFn       func() error
}

func (m *MockComputeClient) StartGame(ctx context.Context, req *connections.StartGameRequest) (*connections.StartGameResponse, error) {
	if m.StartGameFn != nil {
		return m.StartGameFn(ctx, req)
	}
	return &connections.StartGameResponse{}, nil
}

func (m *MockComputeClient) MakeGuess(ctx context.Context, req *connections.GuessRequest) (*connections.GuessResponse, error) {
	if m.MakeGuessFn != nil {
		return m.MakeGuessFn(ctx, req)
	}
	return &connections.GuessResponse{}, nil
}

func (m *MockComputeClient) HealthCheck(ctx context.Context) (*connections.HealthResponse, error) {
	if m.HealthCheckFn != nil {
		return m.HealthCheckFn(ctx)
	}
	return &connections.HealthResponse{Status: "ok"}, nil
}

func (m *MockComputeClient) Close() error {
	if m.CloseFn != nil {
		return m.CloseFn()
	}
	return nil
}
