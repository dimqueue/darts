package service

import (
	"errors"

	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/model"
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

func (s *GuessService) CreateGuess(userId, gameId int, word string) (int, error) {
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

	guess := model.Guess{
		GameId:    gameId,
		GuessWord: word,
		Distance:  distance,
	}
	err = s.guessRepo.CreateGuess(&guess)
	if err != nil {
		return 0, err
	}

	// if distance == 0 {
	//     s.gameRepo.UpdateGame(gameId, "won")
	// }

	return distance, nil
}

func (s *GuessService) GetAllGuessByGame(userId, gameId int) ([]model.Guess, error) {
	game, err := s.gameRepo.GetGameById(gameId)
	if err != nil {
		return nil, err
	}

	if game.UserId != userId {
		return nil, errors.New("unauthorized: cannot view guesses for another user's game")
	}

	guesses, err := s.guessRepo.GetAllGuessByGame(gameId)
	if err != nil {
		return nil, err
	}

	return guesses, nil
}

func (s *GuessService) GetGuessById(i int) error {
	return nil
}
