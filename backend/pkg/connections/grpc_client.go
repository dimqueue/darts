package connections

import (
	"context"
	"fmt"
	"time"

	computev1 "github.com/dimqueue/darts/pkg/proto/compute/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCComputeClient struct {
	conn    *grpc.ClientConn
	client  computev1.ComputeServiceClient
	timeout time.Duration
}

func NewGRPCComputeClient(cfg Config) (*GRPCComputeClient, error) {
	timeout := 30
	if cfg.Timeout > 0 {
		timeout = cfg.Timeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to gRPC server
	conn, err := grpc.DialContext(
		ctx,
		cfg.BaseURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &GRPCComputeClient{
		conn:    conn,
		client:  computev1.NewComputeServiceClient(conn),
		timeout: time.Duration(timeout) * time.Second,
	}, nil
}

func (c *GRPCComputeClient) StartGame(ctx context.Context, req *StartGameRequest) (*StartGameResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Convert DTO to Proto
	protoReq := &computev1.StartGameRequest{
		Language:   req.Language,
		SecretWord: req.SecretWord,
		TopN:       10000, // Default value, could be added to DTO if needed
	}

	protoResp, err := c.client.StartGame(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("grpc StartGame failed: %w", err)
	}

	// Convert Proto to DTO
	return &StartGameResponse{
		CalculationTime: protoResp.CalculationTime,
		HintWord:        protoResp.HintWord,
	}, nil
}

func (c *GRPCComputeClient) MakeGuess(ctx context.Context, req *GuessRequest) (*GuessResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Convert DTO to Proto
	protoReq := &computev1.MakeGuessRequest{
		SecretWord: req.SecretWord,
		Guess:      req.Guess,
		Language:   req.Language,
	}

	protoResp, err := c.client.MakeGuess(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("grpc MakeGuess failed: %w", err)
	}

	// Convert Proto to DTO
	return &GuessResponse{
		Distance: int(protoResp.Distance),
	}, nil
}

func (c *GRPCComputeClient) HealthCheck(ctx context.Context) (*HealthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	protoReq := &computev1.HealthCheckRequest{}

	protoResp, err := c.client.HealthCheck(ctx, protoReq)
	if err != nil {
		return nil, fmt.Errorf("grpc HealthCheck failed: %w", err)
	}

	// Convert Proto to DTO
	return &HealthResponse{
		Status:          protoResp.Status,
		LoadedLanguages: protoResp.LoadedLanguages,
	}, nil
}

func (c *GRPCComputeClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
