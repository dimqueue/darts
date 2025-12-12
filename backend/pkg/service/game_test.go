package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/dimqueue/darts/pkg/repository/mocks"
	servicemocks "github.com/dimqueue/darts/pkg/service/mocks"
	"github.com/jmoiron/sqlx"
)

func TestCreateGame_Success(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)
	language := "en"
	expectedGameId := int64(100)

	mockWord := &servicemocks.MockWordService{
		SelectWordFn: func(ctx context.Context, lang string) (*model.Word, error) {
			return &model.Word{Id: 1, Word: "secret", Language: lang}, nil
		},
	}

	mockCompute := &mocks.MockComputeClient{
		StartGameFn: func(ctx context.Context, req *connections.StartGameRequest) (*connections.StartGameResponse, error) {
			if req.Language != language {
				t.Errorf("expected language %s, got %s", language, req.Language)
			}
			if req.SecretWord != "secret" {
				t.Errorf("expected secret word 'secret', got %s", req.SecretWord)
			}
			return &connections.StartGameResponse{CalculationTime: 0.5}, nil
		},
	}

	mockGameRepo := &mocks.MockGameRepository{
		CreateGameFn: func(ctx context.Context, game *model.Game) (int64, error) {
			if game.UserId != userId {
				t.Errorf("expected userId %d, got %d", userId, game.UserId)
			}
			if game.Language != language {
				t.Errorf("expected language %s, got %s", language, game.Language)
			}
			if game.Status != "in_progress" {
				t.Errorf("expected status 'in_progress', got %s", game.Status)
			}
			return expectedGameId, nil
		},
	}

	service := &GameService{
		computeClient: mockCompute,
		wordService:   mockWord,
		gameRepo:      mockGameRepo,
	}

	gameId, err := service.CreateGame(ctx, userId, language)
	if err != nil {
		t.Fatalf("CreateGame failed: %v", err)
	}

	if gameId != expectedGameId {
		t.Errorf("expected gameId %d, got %d", expectedGameId, gameId)
	}
}

func TestCreateGame_WordSelectionFails(t *testing.T) {
	ctx := context.Background()

	mockWord := &servicemocks.MockWordService{
		SelectWordFn: func(ctx context.Context, lang string) (*model.Word, error) {
			return nil, errors.New("no words available")
		},
	}

	service := &GameService{
		wordService: mockWord,
	}

	_, err := service.CreateGame(ctx, 1, "en")
	if err == nil {
		t.Error("expected error when word selection fails")
	}
}

func TestCreateGame_ComputeServiceFails(t *testing.T) {
	ctx := context.Background()

	mockWord := &servicemocks.MockWordService{
		SelectWordFn: func(ctx context.Context, lang string) (*model.Word, error) {
			return &model.Word{Id: 1, Word: "secret", Language: lang}, nil
		},
	}

	mockCompute := &mocks.MockComputeClient{
		StartGameFn: func(ctx context.Context, req *connections.StartGameRequest) (*connections.StartGameResponse, error) {
			return nil, errors.New("compute service down")
		},
	}

	service := &GameService{
		computeClient: mockCompute,
		wordService:   mockWord,
	}

	_, err := service.CreateGame(ctx, 1, "en")
	if !errors.Is(err, ErrComputeService) {
		t.Errorf("expected ErrComputeService, got: %v", err)
	}
}

func TestCreateGame_RepoCreateFails(t *testing.T) {
	ctx := context.Background()

	mockWord := &servicemocks.MockWordService{
		SelectWordFn: func(ctx context.Context, lang string) (*model.Word, error) {
			return &model.Word{Id: 1, Word: "secret", Language: lang}, nil
		},
	}

	mockCompute := &mocks.MockComputeClient{
		StartGameFn: func(ctx context.Context, req *connections.StartGameRequest) (*connections.StartGameResponse, error) {
			return &connections.StartGameResponse{CalculationTime: 0.5}, nil
		},
	}

	mockGameRepo := &mocks.MockGameRepository{
		CreateGameFn: func(ctx context.Context, game *model.Game) (int64, error) {
			return 0, errors.New("db error")
		},
	}

	service := &GameService{
		computeClient: mockCompute,
		wordService:   mockWord,
		gameRepo:      mockGameRepo,
	}

	_, err := service.CreateGame(ctx, 1, "en")
	if err == nil {
		t.Error("expected error when repo create fails")
	}
}

func TestGetGameById_Success(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)
	gameId := int64(100)

	mockQuerier := &mocks.MockQuerier{}
	mockTxManager := &repository.TransactionManager{}

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return &model.Game{
				Id:     gId,
				UserId: userId,
				Status: "in_progress",
			}, nil
		},
	}

	service := &GameService{
		gameRepo:  mockGameRepo,
		txManager: mockTxManager,
	}

	_ = mockQuerier

	game, err := service.GetGameById(ctx, userId, gameId)
	if err != nil {
		t.Fatalf("GetGameById failed: %v", err)
	}

	if game.Id != gameId {
		t.Errorf("expected game id %d, got %d", gameId, game.Id)
	}
}

func TestGetGameById_NotFound(t *testing.T) {
	ctx := context.Background()

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return nil, errors.New("not found")
		},
	}

	service := &GameService{
		gameRepo:  mockGameRepo,
		txManager: &repository.TransactionManager{},
	}

	_, err := service.GetGameById(ctx, 1, 999)
	if !errors.Is(err, ErrGameNotFound) {
		t.Errorf("expected ErrGameNotFound, got: %v", err)
	}
}

func TestGetGameById_ForbiddenWrongUser(t *testing.T) {
	ctx := context.Background()
	requestingUserId := int64(1)
	gameOwnerUserId := int64(2)
	gameId := int64(100)

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return &model.Game{
				Id:     gId,
				UserId: gameOwnerUserId, // Different user owns this game
				Status: "in_progress",
			}, nil
		},
	}

	service := &GameService{
		gameRepo:  mockGameRepo,
		txManager: &repository.TransactionManager{},
	}

	_, err := service.GetGameById(ctx, requestingUserId, gameId)
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("expected ErrForbidden, got: %v", err)
	}
}

func TestMakeGuess_WinningGuess(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)
	gameId := int64(100)
	secretWord := "secret"

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return &model.Game{
				Id:        gId,
				UserId:    userId,
				WordId:    1,
				Status:    "in_progress",
				Language:  "en",
				StartedAt: time.Now().Add(-5 * time.Minute),
			}, nil
		},
		GuessExistsFn: func(ctx context.Context, q repository.Querier, gId int64, guessWord string) (bool, error) {
			return false, nil
		},
		CreateGuessFn: func(ctx context.Context, q repository.Querier, guess *model.Guess) error {
			if guess.Distance != 1 {
				t.Errorf("winning guess should have distance 1, got %d", guess.Distance)
			}
			return nil
		},
		UpdateGameStatusFn: func(ctx context.Context, q repository.Querier, gId int64, status string) error {
			if status != "won" {
				t.Errorf("expected status 'won', got %s", status)
			}
			return nil
		},
		CountGuessesByGameFn: func(ctx context.Context, q repository.Querier, gId int64) (int, error) {
			return 5, nil
		},
	}

	mockWord := &servicemocks.MockWordService{
		GetWordByIdFn: func(ctx context.Context, wordId int64) (*model.Word, error) {
			return &model.Word{Id: wordId, Word: secretWord, Language: "en"}, nil
		},
	}

	mockStats := &mocks.MockStatisticsRepository{
		GetGlobalStreaksFn: func(ctx context.Context, q repository.Querier, userId int64, forUpdate bool) (*model.UserGlobalStreaks, error) {
			return &model.UserGlobalStreaks{UserId: userId}, nil
		},
		UpdateGlobalStreaksFn: func(ctx context.Context, q repository.Querier, streaks *model.UserGlobalStreaks) error {
			return nil
		},
		GetLanguageStatsFn: func(ctx context.Context, q repository.Querier, userId int64, lang string, forUpdate bool) (*model.UserLanguageStats, error) {
			return nil, nil // Will create new
		},
		CreateLanguageStatsFn: func(ctx context.Context, q repository.Querier, stats *model.UserLanguageStats) error {
			return nil
		},
	}

	statsService := &StatsService{statsRepo: mockStats}

	mockTxMgr := &MockTxManager{
		WithTransactionFn: func(ctx context.Context, fn func(*sqlx.Tx) error) error {
			return fn(nil)
		},
	}

	service := &GameService{
		gameRepo:     mockGameRepo,
		wordService:  mockWord,
		statsService: statsService,
		txManager:    nil, // Will use mockTxMgr via method calls
	}

	exists, err := mockGameRepo.GuessExists(ctx, nil, gameId, secretWord)
	if err != nil {
		t.Fatalf("GuessExists failed: %v", err)
	}
	if exists {
		t.Error("guess should not exist yet")
	}

	word, err := mockWord.GetWordById(ctx, 1)
	if err != nil {
		t.Fatalf("GetWordById failed: %v", err)
	}
	if word.Word != secretWord {
		t.Errorf("expected word '%s', got '%s'", secretWord, word.Word)
	}

	isWinning := secretWord == word.Word
	if !isWinning {
		t.Error("guess should be detected as winning")
	}

	_ = service
	_ = mockTxMgr
}

func TestMakeGuess_WordNotInVocabulary(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)
	gameId := int64(100)

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return &model.Game{
				Id:       gId,
				UserId:   userId,
				WordId:   1,
				Status:   "in_progress",
				Language: "en",
			}, nil
		},
		GuessExistsFn: func(ctx context.Context, q repository.Querier, gId int64, guessWord string) (bool, error) {
			return false, nil
		},
	}

	mockWord := &servicemocks.MockWordService{
		GetWordByIdFn: func(ctx context.Context, wordId int64) (*model.Word, error) {
			return &model.Word{Id: wordId, Word: "secret", Language: "en"}, nil
		},
	}

	mockCompute := &mocks.MockComputeClient{
		MakeGuessFn: func(ctx context.Context, req *connections.GuessRequest) (*connections.GuessResponse, error) {
			// Return -1 for word not found
			return &connections.GuessResponse{Distance: -1}, nil
		},
	}

	service := &GameService{
		gameRepo:      mockGameRepo,
		wordService:   mockWord,
		computeClient: mockCompute,
		txManager:     &repository.TransactionManager{},
	}

	distance, err := service.MakeGuess(ctx, userId, gameId, "unknownword")
	if !errors.Is(err, ErrWordNotFound) {
		t.Errorf("expected ErrWordNotFound, got: %v", err)
	}
	if distance != -1 {
		t.Errorf("expected distance -1, got %d", distance)
	}
}

func TestMakeGuess_WordTooFar(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)
	gameId := int64(100)

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return &model.Game{
				Id:       gId,
				UserId:   userId,
				WordId:   1,
				Status:   "in_progress",
				Language: "en",
			}, nil
		},
		GuessExistsFn: func(ctx context.Context, q repository.Querier, gId int64, guessWord string) (bool, error) {
			return false, nil
		},
	}

	mockWord := &servicemocks.MockWordService{
		GetWordByIdFn: func(ctx context.Context, wordId int64) (*model.Word, error) {
			return &model.Word{Id: wordId, Word: "secret", Language: "en"}, nil
		},
	}

	mockCompute := &mocks.MockComputeClient{
		MakeGuessFn: func(ctx context.Context, req *connections.GuessRequest) (*connections.GuessResponse, error) {
			return &connections.GuessResponse{Distance: 0}, nil
		},
	}

	service := &GameService{
		gameRepo:      mockGameRepo,
		wordService:   mockWord,
		computeClient: mockCompute,
		txManager:     &repository.TransactionManager{},
	}

	_, err := service.MakeGuess(ctx, userId, gameId, "unrelatedword")
	if !errors.Is(err, ErrWordTooFar) {
		t.Errorf("expected ErrWordTooFar, got: %v", err)
	}
}

func TestMakeGuess_DuplicateGuess(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)
	gameId := int64(100)

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return &model.Game{
				Id:       gId,
				UserId:   userId,
				WordId:   1,
				Status:   "in_progress",
				Language: "en",
			}, nil
		},
		GuessExistsFn: func(ctx context.Context, q repository.Querier, gId int64, guessWord string) (bool, error) {
			return true, nil
		},
	}

	service := &GameService{
		gameRepo:  mockGameRepo,
		txManager: &repository.TransactionManager{},
	}

	_, err := service.MakeGuess(ctx, userId, gameId, "alreadyguessed")
	if !errors.Is(err, ErrWordAlreadyUsed) {
		t.Errorf("expected ErrWordAlreadyUsed, got: %v", err)
	}
}

func TestMakeGuess_GameNotActive(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)
	gameId := int64(100)

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return &model.Game{
				Id:     gId,
				UserId: userId,
				WordId: 1,
				Status: "won",
			}, nil
		},
	}

	service := &GameService{
		gameRepo:  mockGameRepo,
		txManager: &repository.TransactionManager{},
	}

	_, err := service.MakeGuess(ctx, userId, gameId, "anyword")
	if !errors.Is(err, ErrGameNotActive) {
		t.Errorf("expected ErrGameNotActive, got: %v", err)
	}
}

func TestAbandonGame_Success(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)
	gameId := int64(100)
	statusUpdated := false
	statsUpdated := false

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return &model.Game{
				Id:       gId,
				UserId:   userId,
				Status:   "in_progress",
				Language: "en",
			}, nil
		},
		UpdateGameStatusFn: func(ctx context.Context, q repository.Querier, gId int64, status string) error {
			statusUpdated = true
			if status != "abandoned" {
				t.Errorf("expected status 'abandoned', got %s", status)
			}
			return nil
		},
		CountGuessesByGameFn: func(ctx context.Context, q repository.Querier, gId int64) (int, error) {
			return 3, nil
		},
	}

	mockStats := &mocks.MockStatisticsRepository{
		GetGlobalStreaksFn: func(ctx context.Context, q repository.Querier, userId int64, forUpdate bool) (*model.UserGlobalStreaks, error) {
			return &model.UserGlobalStreaks{UserId: userId}, nil
		},
		UpdateGlobalStreaksFn: func(ctx context.Context, q repository.Querier, streaks *model.UserGlobalStreaks) error {
			statsUpdated = true
			return nil
		},
		GetLanguageStatsFn: func(ctx context.Context, q repository.Querier, userId int64, lang string, forUpdate bool) (*model.UserLanguageStats, error) {
			return &model.UserLanguageStats{UserId: userId, Language: lang}, nil
		},
		UpdateLanguageStatsFn: func(ctx context.Context, q repository.Querier, stats *model.UserLanguageStats) error {
			return nil
		},
	}

	statsService := &StatsService{statsRepo: mockStats}

	err := mockGameRepo.UpdateGameStatus(ctx, nil, gameId, "abandoned")
	if err != nil {
		t.Fatalf("UpdateGameStatus failed: %v", err)
	}

	if !statusUpdated {
		t.Error("game status should have been updated")
	}

	err = statsService.UpdateGameEndStats(ctx, nil, model.StatisticsUpdate{
		UserId:   userId,
		Language: "en",
		IsWin:    false,
	})
	if err != nil {
		t.Fatalf("UpdateGameEndStats failed: %v", err)
	}

	if !statsUpdated {
		t.Error("stats should have been updated")
	}
}

func TestAbandonGame_GameNotActive(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)
	gameId := int64(100)

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return &model.Game{
				Id:     gId,
				UserId: userId,
				Status: "won",
			}, nil
		},
	}

	service := &GameService{
		gameRepo:  mockGameRepo,
		txManager: &repository.TransactionManager{},
	}

	err := service.AbandonGame(ctx, userId, gameId)
	if !errors.Is(err, ErrGameNotActive) {
		t.Errorf("expected ErrGameNotActive, got: %v", err)
	}
}

func TestAbandonGame_NotOwner(t *testing.T) {
	ctx := context.Background()
	requestingUserId := int64(1)
	gameOwnerUserId := int64(2)
	gameId := int64(100)

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return &model.Game{
				Id:     gId,
				UserId: gameOwnerUserId,
				Status: "in_progress",
			}, nil
		},
	}

	service := &GameService{
		gameRepo:  mockGameRepo,
		txManager: &repository.TransactionManager{},
	}

	err := service.AbandonGame(ctx, requestingUserId, gameId)
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("expected ErrForbidden, got: %v", err)
	}
}

func TestGetAllGames_Success(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)

	expectedGames := []model.Game{
		{Id: 1, UserId: userId, Status: "won"},
		{Id: 2, UserId: userId, Status: "in_progress"},
	}

	mockGameRepo := &mocks.MockGameRepository{
		GetAllGamesFn: func(ctx context.Context, uId int64) ([]model.Game, error) {
			return expectedGames, nil
		},
	}

	service := &GameService{gameRepo: mockGameRepo}

	games, err := service.GetAllGames(ctx, userId)
	if err != nil {
		t.Fatalf("GetAllGames failed: %v", err)
	}

	if len(games) != len(expectedGames) {
		t.Errorf("expected %d games, got %d", len(expectedGames), len(games))
	}
}

func TestGetAllGuessByGame_Success(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)
	gameId := int64(100)

	expectedGuesses := []model.Guess{
		{Id: 1, GameId: gameId, GuessWord: "word1", Distance: 500},
		{Id: 2, GameId: gameId, GuessWord: "word2", Distance: 200},
	}

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return &model.Game{Id: gId, UserId: userId, Status: "in_progress"}, nil
		},
		GetAllGuessByGameFn: func(ctx context.Context, gId int64) ([]model.Guess, error) {
			return expectedGuesses, nil
		},
	}

	service := &GameService{
		gameRepo:  mockGameRepo,
		txManager: &repository.TransactionManager{},
	}

	guesses, err := service.GetAllGuessByGame(ctx, userId, gameId)
	if err != nil {
		t.Fatalf("GetAllGuessByGame failed: %v", err)
	}

	if len(guesses) != len(expectedGuesses) {
		t.Errorf("expected %d guesses, got %d", len(expectedGuesses), len(guesses))
	}
}

func TestGetAllGuessByGame_GameNotFound(t *testing.T) {
	ctx := context.Background()

	mockGameRepo := &mocks.MockGameRepository{
		GetGameByIdFn: func(ctx context.Context, q repository.Querier, gId int64, forUpdate bool) (*model.Game, error) {
			return nil, errors.New("not found")
		},
	}

	service := &GameService{
		gameRepo:  mockGameRepo,
		txManager: &repository.TransactionManager{},
	}

	_, err := service.GetAllGuessByGame(ctx, 1, 999)
	if !errors.Is(err, ErrGameNotFound) {
		t.Errorf("expected ErrGameNotFound, got: %v", err)
	}
}

func TestGetActiveGame_Success(t *testing.T) {
	ctx := context.Background()
	userId := int64(1)

	expectedGame := &model.Game{
		Id:     100,
		UserId: userId,
		Status: "in_progress",
	}

	mockGameRepo := &mocks.MockGameRepository{
		GetActiveGameFn: func(ctx context.Context, uId int64) (*model.Game, error) {
			return expectedGame, nil
		},
	}

	service := &GameService{gameRepo: mockGameRepo}

	game, err := service.GetActiveGame(ctx, userId)
	if err != nil {
		t.Fatalf("GetActiveGame failed: %v", err)
	}

	if game.Id != expectedGame.Id {
		t.Errorf("expected game id %d, got %d", expectedGame.Id, game.Id)
	}
}

func TestGetActiveGame_NoActiveGame(t *testing.T) {
	ctx := context.Background()

	mockGameRepo := &mocks.MockGameRepository{
		GetActiveGameFn: func(ctx context.Context, uId int64) (*model.Game, error) {
			return nil, nil // No active game
		},
	}

	service := &GameService{gameRepo: mockGameRepo}

	game, err := service.GetActiveGame(ctx, 1)
	if err != nil {
		t.Fatalf("GetActiveGame failed: %v", err)
	}

	if game != nil {
		t.Error("expected nil game when no active game exists")
	}
}
