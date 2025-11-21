package service

import (
	"fmt"

	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type GameService struct {
	computeClient *connections.ComputeClientService
	wordService   Word
	gameRepo      repository.Game
}

func NewGameService(gameRepo repository.Game, wordService Word, computeClient *connections.ComputeClientService) *GameService {
	return &GameService{
		computeClient: computeClient,
		wordService:   wordService,
		gameRepo:      gameRepo,
	}
}

func (s *GameService) CreateGame(userId int, lang string) (int, error) {
	selectedWord, err := s.wordService.SelectWord(lang)
	if err != nil {
		return 0, fmt.Errorf("failed to select word: %w", err)
	}

	// Start game on compute client with the selected word
	resp, err := s.computeClient.StartGame(selectedWord.Word, lang)
	if err != nil {
		return 0, fmt.Errorf("failed to start game on compute client: %w", err)
	}

	// Create game record in database
	var game model.Game
	game.Language = lang
	game.Status = "in_progress"
	game.UserId = userId
	game.WordId = selectedWord.Id

	fmt.Printf("Game started - Word: %s, Calculation time: %.3fs, Hint: %s\n",
		selectedWord.Word, resp.CalculationTime, resp.HintWord)

	return s.gameRepo.CreateGame(&game)
}

func (s *GameService) GetAllGames(userId int) ([]model.Game, error) {
	return s.gameRepo.GetAllGames(userId)
}

func (s *GameService) GetGameById(userId, gameId int) (*model.Game, error) {
	game, err := s.gameRepo.GetGameById(gameId)
	if err != nil {
		return nil, err
	}

	if game.UserId != userId {
		return nil, fmt.Errorf("unauthorized: access denied")
	}

	return game, nil
}

// ///
func (s *GameService) UpdateGame(gameId int) (*model.Game, error) {
	return nil, nil
}
func (s *GameService) DeleteGame(gameId int) (*model.Game, error) {
	return nil, nil
}
