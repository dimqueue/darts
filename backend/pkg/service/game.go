package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dimqueue/darts/pkg/config"
	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/logger"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/jmoiron/sqlx"
)

type GameService struct {
	computeClient connections.ComputeClient
	wordService   Word
	gameRepo      repository.Game
	statsService  *StatsService
	txManager     *repository.TransactionManager
}

func NewGameService(gameRepo repository.Game, wordService Word, statsService *StatsService, txManager *repository.TransactionManager, computeClient connections.ComputeClient) *GameService {
	return &GameService{
		computeClient: computeClient,
		wordService:   wordService,
		gameRepo:      gameRepo,
		statsService:  statsService,
		txManager:     txManager,
	}
}

func (s *GameService) CreateGame(ctx context.Context, userId int64, lang string) (int64, error) {
	log := logger.WithContext(ctx).WithFields(map[string]interface{}{
		"user_id":  userId,
		"language": lang,
	})

	selectedWord, err := s.wordService.SelectWord(ctx, lang)
	if err != nil {
		log.WithError(err).Error("failed to select word")
		return 0, fmt.Errorf("failed to select word: %w", err)
	}

	grpcCtx, cancel := context.WithTimeout(ctx, config.GameComputeTimeout)
	defer cancel()

	req := &connections.StartGameRequest{
		Language:   lang,
		SecretWord: selectedWord.Word,
	}

	log.WithField("word", selectedWord.Word).Debug("calling compute service to start game")
	resp, err := s.computeClient.StartGame(grpcCtx, req)
	if err != nil {
		log.WithError(err).Error("grpc StartGame failed")
		return 0, fmt.Errorf("failed to start game on compute client: %w", err)
	}
	log.WithFields(map[string]interface{}{
		"word":      selectedWord.Word,
		"hint_word": resp.HintWord,
	}).Debug("compute service responded")

	expiresAt := time.Now().Add(config.GameTTL)

	game := model.Game{
		Language:  lang,
		Status:    "in_progress",
		UserId:    userId,
		WordId:    selectedWord.Id,
		ExpiresAt: &expiresAt,
	}

	gameId, err := s.gameRepo.CreateGame(ctx, &game)
	if err != nil {
		log.WithError(err).Error("failed to save game to database")
		return 0, err
	}

	log.WithFields(map[string]interface{}{
		"game_id":          gameId,
		"calculation_time": resp.CalculationTime,
	}).Info("game created")

	return gameId, nil
}

func (s *GameService) GetAllGames(ctx context.Context, userId int64) ([]model.Game, error) {
	return s.gameRepo.GetAllGames(ctx, userId)
}

func (s *GameService) GetGameById(ctx context.Context, userId, gameId int64) (*model.Game, error) {
	game, err := s.gameRepo.GetGameById(ctx, s.txManager.DB(), gameId, false)
	if err != nil {
		return nil, err
	}

	if game.UserId != userId {
		return nil, fmt.Errorf("unauthorized: access denied")
	}

	return game, nil
}

func (s *GameService) UpdateGameStatus(ctx context.Context, gameId int64, status string) error {
	return s.gameRepo.UpdateGameStatus(ctx, s.txManager.DB(), gameId, status)
}

func (s *GameService) MakeGuess(ctx context.Context, userId, gameId int64, guess string) (int, error) {
	log := logger.WithContext(ctx).WithFields(map[string]interface{}{
		"user_id": userId,
		"game_id": gameId,
	})

	game, err := s.GetGameById(ctx, userId, gameId)
	if err != nil {
		return 0, err
	}

	if game.Status != "in_progress" {
		log.Warn("attempt to guess on inactive game")
		return 0, errors.New("game is not active")
	}

	exists, err := s.gameRepo.GuessExists(ctx, s.txManager.DB(), gameId, guess)
	if err != nil {
		log.WithError(err).Error("failed to check if guess exists")
		return 0, fmt.Errorf("failed to check guess: %w", err)
	}
	if exists {
		return 0, errors.New("word already guessed")
	}

	word, err := s.wordService.GetWordById(ctx, game.WordId)
	if err != nil {
		log.WithError(err).Error("failed to get target word")
		return 0, err
	}

	var distance int
	isWinningGuess := guess == word.Word

	if isWinningGuess {
		distance = 1
	} else {
		grpcCtx, cancel := context.WithTimeout(ctx, config.GameComputeTimeout)
		defer cancel()

		req := &connections.GuessRequest{
			SecretWord: word.Word,
			Guess:      guess,
			Language:   game.Language,
		}

		log.WithFields(map[string]interface{}{
			"target_word": word.Word,
			"guess":       guess,
		}).Debug("calling compute service for distance calculation")
		resp, err := s.computeClient.MakeGuess(grpcCtx, req)
		if err != nil {
			log.WithError(err).Error("grpc MakeGuess failed")
			return 0, err
		}
		distance = resp.Distance
		log.WithField("distance", distance).Debug("distance calculated")
	}

	err = s.txManager.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		lockedGame, err := s.gameRepo.GetGameById(ctx, tx, gameId, true)
		if err != nil {
			return fmt.Errorf("failed to lock game: %w", err)
		}

		if lockedGame.Status != "in_progress" {
			return errors.New("game is no longer active")
		}

		existsTx, err := s.gameRepo.GuessExists(ctx, tx, gameId, guess)
		if err != nil {
			return fmt.Errorf("failed to check guess in transaction: %w", err)
		}
		if existsTx {
			return errors.New("word already guessed")
		}

		newGuess := model.Guess{
			GameId:    gameId,
			GuessWord: guess,
			Distance:  distance,
		}
		if err := s.gameRepo.CreateGuess(ctx, tx, &newGuess); err != nil {
			return fmt.Errorf("failed to create guess: %w", err)
		}

		if isWinningGuess {
			if err := s.gameRepo.UpdateGameStatus(ctx, tx, gameId, "won"); err != nil {
				return fmt.Errorf("failed to update game status: %w", err)
			}

			guessCount, err := s.gameRepo.CountGuessesByGame(ctx, tx, gameId)
			if err != nil {
				return fmt.Errorf("failed to count guesses: %w", err)
			}

			elapsed := int(time.Since(lockedGame.StartedAt).Seconds())

			statsUpdate := model.StatisticsUpdate{
				UserId:      userId,
				Language:    lockedGame.Language,
				IsWin:       true,
				GuessCount:  guessCount,
				TimeSeconds: &elapsed,
				ScoreEarned: config.ScorePerWin,
			}

			if err := s.statsService.UpdateGameEndStats(ctx, tx, statsUpdate); err != nil {
				return fmt.Errorf("failed to update statistics: %w", err)
			}

			log.WithFields(map[string]interface{}{
				"guess_count":  guessCount,
				"time_seconds": elapsed,
			}).Info("game won")
		}
		return nil
	})

	if err != nil {
		log.WithError(err).Error("failed to process guess transaction")
		return 0, err
	}

	return distance, nil
}

func (s *GameService) GetAllGuessByGame(ctx context.Context, userId, gameId int64) ([]model.Guess, error) {
	_, err := s.GetGameById(ctx, userId, gameId)
	if err != nil {
		return nil, err
	}

	return s.gameRepo.GetAllGuessByGame(ctx, gameId)
}
