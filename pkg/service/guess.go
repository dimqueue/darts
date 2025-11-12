package service

import (
	"errors"

	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/repository"
)

type GuessService struct {
	client    connections.Client
	guessRepo repository.Guess
	gameRepo  repository.Game
}

func NewGuessService(guessRepo repository.Guess, gameRepo repository.Game, client connections.Client) *GuessService {
	return &GuessService{
		guessRepo: guessRepo,
		gameRepo:  gameRepo,
		client:    client,
	}
}

func (s *GuessService) CreateGuess(userId, gameId int, guess string) (int, error) {
	game, err := s.gameRepo.GetGameById(gameId)
	if err != nil {
		return 0, err
	}

	if game.UserId != userId {
		return 0, errors.New("unauthorized: cannot add guess to another user's game")
	}

	if game.Status != "in_progress" {
		return 0, errors.New("game is not active")
	}

	// distance, err := s.client.CalculateDistance(game.Word, guess)
	// if err != nil {
	//     return 0, err
	// }
	distance := 5

	err = s.guessRepo.CreateGuess(gameId, guess, distance)
	if err != nil {
		return 0, err
	}

	// if distance == 0 {
	//     s.gameRepo.UpdateGame(gameId, "won")
	// }

	return distance, nil
}

func (s *GuessService) GetAllGuessByGame(i int) error {
	return nil
}

func (s *GuessService) GetGuessById(i int) error {
	return nil
}
