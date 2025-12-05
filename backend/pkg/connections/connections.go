package connections

import "context"

type ComputeClient interface {
	StartGame(ctx context.Context, req *StartGameRequest) (*StartGameResponse, error)
	MakeGuess(ctx context.Context, req *GuessRequest) (*GuessResponse, error)
	HealthCheck(ctx context.Context) (*HealthResponse, error)
	Close() error
}

type Config struct {
	Type    string
	BaseURL string
	Timeout int
}

func NewComputeClient(cfg Config) (ComputeClient, error) {
	switch cfg.Type {
	case "grpc":
		return NewGRPCComputeClient(cfg)
	case "http":
		return NewHTTPComputeClient(cfg)
	default:
		return NewHTTPComputeClient(cfg) // default to HTTP
	}
}
