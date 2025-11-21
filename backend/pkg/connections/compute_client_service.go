package connections

import "fmt"

type ComputeClientService struct {
	client Client
}

func NewComputeClientService(client Client) *ComputeClientService {
	return &ComputeClientService{
		client: client,
	}
}

func (s *ComputeClientService) StartGame(secretWord, language string) (*StartGameResponse, error) {
	req := StartGameRequest{
		Language:   language,
		SecretWord: secretWord,
	}

	var resp StartGameResponse
	if err := s.client.Call("POST", "/start-game", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to start game: %w", err)
	}

	return &resp, nil
}

func (s *ComputeClientService) MakeGuess(secretWord, guess, language string) (*GuessResponse, error) {
	req := GuessRequest{
		SecretWord: secretWord,
		Guess:      guess,
		Language:   language,
	}

	var resp GuessResponse
	if err := s.client.Call("POST", "/guess", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to make guess: %w", err)
	}

	return &resp, nil
}

func (s *ComputeClientService) HealthCheck() (*HealthResponse, error) {
	var resp HealthResponse
	if err := s.client.Call("GET", "/health", nil, &resp); err != nil {
		return nil, fmt.Errorf("health check failed: %w", err)
	}

	return &resp, nil
}
