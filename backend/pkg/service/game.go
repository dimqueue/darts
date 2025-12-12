package service

import (
	"context"
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
	log := logger.Op(ctx, "GameService.CreateGame").With(
		logger.FieldUserID, userId,
		logger.FieldLanguage, lang,
	)

	selectedWord, err := s.wordService.SelectWord(ctx, lang)
	if err != nil {
		log.Error("select word failed", logger.FieldError, err)
		return 0, Logged(err)
	}

	log.Debug("secret word selected", logger.FieldWord, selectedWord.Word)

	grpcCtx, cancel := context.WithTimeout(ctx, config.GameComputeTimeout)
	defer cancel()

	req := &connections.StartGameRequest{
		Language:   lang,
		SecretWord: selectedWord.Word,
	}

	resp, err := s.computeClient.StartGame(grpcCtx, req)
	if err != nil {
		log.Error("compute service StartGame failed", logger.FieldError, err)
		return 0, ErrComputeService
	}

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
		log.Error("db create game failed", logger.FieldError, err)
		return 0, Logged(err)
	}

	log.Info("game created",
		logger.FieldGameID, gameId,
		logger.FieldWord, selectedWord.Word,
		logger.FieldCalculationTime, resp.CalculationTime,
	)

	return gameId, nil
}

func (s *GameService) GetAllGames(ctx context.Context, userId int64) ([]model.Game, error) {
	return s.gameRepo.GetAllGames(ctx, userId)
}

func (s *GameService) GetActiveGame(ctx context.Context, userId int64) (*model.Game, error) {
	return s.gameRepo.GetActiveGame(ctx, userId)
}

func (s *GameService) GetGameById(ctx context.Context, userId, gameId int64) (*model.Game, error) {
	game, err := s.gameRepo.GetGameById(ctx, s.txManager.DB(), gameId, false)
	if err != nil {
		return nil, ErrGameNotFound
	}

	if game.UserId != userId {
		return nil, ErrForbidden
	}

	return game, nil
}

func (s *GameService) UpdateGameStatus(ctx context.Context, gameId int64, status string) error {
	return s.gameRepo.UpdateGameStatus(ctx, s.txManager.DB(), gameId, status)
}

func (s *GameService) MakeGuess(ctx context.Context, userId, gameId int64, guess string) (int, error) {
	log := logger.Op(ctx, "GameService.MakeGuess").With(
		logger.FieldUserID, userId,
		logger.FieldGameID, gameId,
		logger.FieldGuess, guess,
	)

	game, err := s.GetGameById(ctx, userId, gameId)
	if err != nil {
		return 0, err
	}

	if game.Status != "in_progress" {
		return 0, ErrGameNotActive
	}

	exists, err := s.gameRepo.GuessExists(ctx, s.txManager.DB(), gameId, guess)
	if err != nil {
		log.Error("db check guess exists failed", logger.FieldError, err)
		return 0, Logged(err)
	}
	if exists {
		return 0, ErrWordAlreadyUsed
	}

	word, err := s.wordService.GetWordById(ctx, game.WordId)
	if err != nil {
		log.Error("get target word failed", logger.FieldError, err)
		return 0, Logged(err)
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

		resp, err := s.computeClient.MakeGuess(grpcCtx, req)
		if err != nil {
			log.Error("compute service MakeGuess failed", logger.FieldError, err)
			return 0, ErrComputeService
		}
		distance = resp.Distance

		if distance == -1 {
			return -1, ErrWordNotFound
		}
		if distance == 0 {
			return 0, ErrWordTooFar
		}
	}

	err = s.txManager.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		lockedGame, err := s.gameRepo.GetGameById(ctx, tx, gameId, true)
		if err != nil {
			return err
		}

		if lockedGame.Status != "in_progress" {
			return ErrGameNotActive
		}

		existsTx, err := s.gameRepo.GuessExists(ctx, tx, gameId, guess)
		if err != nil {
			return err
		}
		if existsTx {
			return ErrWordAlreadyUsed
		}

		newGuess := model.Guess{
			GameId:    gameId,
			GuessWord: guess,
			Distance:  distance,
		}
		if err := s.gameRepo.CreateGuess(ctx, tx, &newGuess); err != nil {
			return err
		}

		if isWinningGuess {
			if err := s.gameRepo.UpdateGameStatus(ctx, tx, gameId, "won"); err != nil {
				return err
			}

			guessCount, err := s.gameRepo.CountGuessesByGame(ctx, tx, gameId)
			if err != nil {
				return err
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
				return err
			}

			log.Info("game won",
				"guess_count", guessCount,
				"time_seconds", elapsed,
			)
		}
		return nil
	})

	if err != nil {
		if !IsDomainError(err) {
			log.Error("guess transaction failed", logger.FieldError, err)
			return 0, Logged(err)
		}
		return 0, err
	}

	if !isWinningGuess {
		log.Debug("guess processed", logger.FieldDistance, distance)
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

func (s *GameService) AbandonGame(ctx context.Context, userId, gameId int64) error {
	log := logger.Op(ctx, "GameService.AbandonGame").With(
		logger.FieldUserID, userId,
		logger.FieldGameID, gameId,
	)

	game, err := s.GetGameById(ctx, userId, gameId)
	if err != nil {
		return err
	}

	if game.Status != "in_progress" {
		return ErrGameNotActive
	}

	err = s.txManager.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		lockedGame, err := s.gameRepo.GetGameById(ctx, tx, gameId, true)
		if err != nil {
			return err
		}

		if lockedGame.Status != "in_progress" {
			return nil
		}

		if err := s.gameRepo.UpdateGameStatus(ctx, tx, gameId, "abandoned"); err != nil {
			return err
		}

		guessCount, err := s.gameRepo.CountGuessesByGame(ctx, tx, gameId)
		if err != nil {
			return err
		}

		statsUpdate := model.StatisticsUpdate{
			UserId:      userId,
			Language:    lockedGame.Language,
			IsWin:       false,
			GuessCount:  guessCount,
			ScoreEarned: 0,
		}

		if err := s.statsService.UpdateGameEndStats(ctx, tx, statsUpdate); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Error("abandon transaction failed", logger.FieldError, err)
		return Logged(err)
	}

	log.Info("game abandoned")
	return nil
}
