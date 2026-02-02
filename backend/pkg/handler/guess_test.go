package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/service"
	servicemocks "github.com/dimqueue/darts/pkg/service/mocks"
	"github.com/dimqueue/darts/pkg/validation"
	"github.com/gin-gonic/gin"
)

func setupGuessHandler(mockAuth *servicemocks.MockAuthService, mockGame *servicemocks.MockGameService) *Handler {
	return &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Game:          mockGame,
		},
		validator: validation.New(),
	}
}

func TestCreateGuess_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		MakeGuessFn: func(ctx context.Context, userId, gameId int64, guess string) (*model.GuessResult, error) {
			if userId != 42 {
				t.Errorf("expected userId 42, got: %d", userId)
			}
			if gameId != 123 {
				t.Errorf("expected gameId 123, got: %d", gameId)
			}
			if guess != "hello" {
				t.Errorf("expected guess 'hello', got: %s", guess)
			}
			return &model.GuessResult{Rank: 150, Found: false, InVocabulary: true}, nil
		},
	}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/guesses", handler.userIdentity, handler.createGuess)

	body := `{"guess": "hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games/123/guesses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["rank"] != float64(150) {
		t.Errorf("expected rank 150, got: %v", response["rank"])
	}
	if response["found"] != false {
		t.Errorf("expected found false, got: %v", response["found"])
	}
	if response["in_vocabulary"] != true {
		t.Errorf("expected in_vocabulary true, got: %v", response["in_vocabulary"])
	}
}

func TestCreateGuess_CorrectWord(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		MakeGuessFn: func(ctx context.Context, userId, gameId int64, guess string) (*model.GuessResult, error) {
			return &model.GuessResult{Rank: 1, Found: true, InVocabulary: true}, nil
		},
	}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/guesses", handler.userIdentity, handler.createGuess)

	body := `{"guess": "correct"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games/123/guesses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["rank"] != float64(1) {
		t.Errorf("expected rank 1, got: %v", response["rank"])
	}
	if response["found"] != true {
		t.Errorf("expected found true, got: %v", response["found"])
	}
}

func TestCreateGuess_Unauthorized(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 0, service.ErrUnauthorized
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/guesses", handler.userIdentity, handler.createGuess)

	body := `{"guess": "hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games/123/guesses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got: %d", w.Code)
	}
}

func TestCreateGuess_InvalidGameId(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/guesses", handler.userIdentity, handler.createGuess)

	body := `{"guess": "hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games/invalid/guesses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestCreateGuess_InvalidJSON(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/guesses", handler.userIdentity, handler.createGuess)

	body := `{invalid}`
	req := httptest.NewRequest(http.MethodPost, "/api/games/123/guesses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestCreateGuess_EmptyGuess(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/guesses", handler.userIdentity, handler.createGuess)

	body := `{"guess": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/games/123/guesses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestCreateGuess_TooShort(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/guesses", handler.userIdentity, handler.createGuess)

	body := `{"guess": "a"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games/123/guesses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestCreateGuess_GameNotFound(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		MakeGuessFn: func(ctx context.Context, userId, gameId int64, guess string) (*model.GuessResult, error) {
			return nil, service.ErrGameNotFound
		},
	}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/guesses", handler.userIdentity, handler.createGuess)

	body := `{"guess": "hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games/999/guesses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got: %d", w.Code)
	}
}

func TestCreateGuess_GameNotActive(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		MakeGuessFn: func(ctx context.Context, userId, gameId int64, guess string) (*model.GuessResult, error) {
			return nil, service.ErrGameNotActive
		},
	}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/guesses", handler.userIdentity, handler.createGuess)

	body := `{"guess": "hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games/123/guesses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestCreateGuess_WordAlreadyUsed(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		MakeGuessFn: func(ctx context.Context, userId, gameId int64, guess string) (*model.GuessResult, error) {
			return nil, service.ErrWordAlreadyUsed
		},
	}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/guesses", handler.userIdentity, handler.createGuess)

	body := `{"guess": "duplicate"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games/123/guesses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestCreateGuess_WordNotFound(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		MakeGuessFn: func(ctx context.Context, userId, gameId int64, guess string) (*model.GuessResult, error) {
			return &model.GuessResult{Rank: 0, Found: false, InVocabulary: false}, service.ErrWordNotFound
		},
	}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/guesses", handler.userIdentity, handler.createGuess)

	body := `{"guess": "xyzabc"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games/123/guesses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestCreateGuess_ServiceError(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		MakeGuessFn: func(ctx context.Context, userId, gameId int64, guess string) (*model.GuessResult, error) {
			return nil, errors.New("compute service error")
		},
	}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/guesses", handler.userIdentity, handler.createGuess)

	body := `{"guess": "hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games/123/guesses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got: %d", w.Code)
	}
}

func TestGetAllGuessByGame_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		GetAllGuessByGameFn: func(ctx context.Context, userId, gameId int64) ([]model.Guess, error) {
			return []model.Guess{
				{Id: 1, GameId: gameId, GuessWord: "hello", Distance: 500, CreatedAt: time.Now()},
				{Id: 2, GameId: gameId, GuessWord: "world", Distance: 200, CreatedAt: time.Now()},
			}, nil
		},
	}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games/:id/guesses", handler.userIdentity, handler.getAllGuessByGame)

	req := httptest.NewRequest(http.MethodGet, "/api/games/123/guesses", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response getAllGuessByGameResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response.Data) != 2 {
		t.Errorf("expected 2 guesses, got: %d", len(response.Data))
	}
}

func TestGetAllGuessByGame_Empty(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		GetAllGuessByGameFn: func(ctx context.Context, userId, gameId int64) ([]model.Guess, error) {
			return []model.Guess{}, nil
		},
	}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games/:id/guesses", handler.userIdentity, handler.getAllGuessByGame)

	req := httptest.NewRequest(http.MethodGet, "/api/games/123/guesses", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response getAllGuessByGameResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response.Data) != 0 {
		t.Errorf("expected 0 guesses, got: %d", len(response.Data))
	}
}

func TestGetAllGuessByGame_InvalidGameId(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games/:id/guesses", handler.userIdentity, handler.getAllGuessByGame)

	req := httptest.NewRequest(http.MethodGet, "/api/games/invalid/guesses", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestGetAllGuessByGame_GameNotFound(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		GetAllGuessByGameFn: func(ctx context.Context, userId, gameId int64) ([]model.Guess, error) {
			return nil, service.ErrGameNotFound
		},
	}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games/:id/guesses", handler.userIdentity, handler.getAllGuessByGame)

	req := httptest.NewRequest(http.MethodGet, "/api/games/999/guesses", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got: %d", w.Code)
	}
}

func TestGetAllGuessByGame_Unauthorized(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 0, service.ErrUnauthorized
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGuessHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games/:id/guesses", handler.userIdentity, handler.getAllGuessByGame)

	req := httptest.NewRequest(http.MethodGet, "/api/games/123/guesses", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got: %d", w.Code)
	}
}
