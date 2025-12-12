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

func setupGameHandler(mockAuth *servicemocks.MockAuthService, mockGame *servicemocks.MockGameService) *Handler {
	return &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Game:          mockGame,
		},
		validator: validation.New(),
	}
}

func setupAuthenticatedRouter(handler *Handler, mockAuth *servicemocks.MockAuthService) *gin.Engine {
	router := gin.New()
	router.Use(RequestIDMiddleware())
	return router
}

func addAuthContext(c *gin.Context, userId int64) {
	c.Set(userCtx, userId)
}

func TestCreateGame_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		CreateGameFn: func(ctx context.Context, userId int64, lang string) (int64, error) {
			if userId != 42 || lang != "en" {
				t.Errorf("unexpected params: userId=%d, lang=%s", userId, lang)
			}
			return 100, nil
		},
	}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games", handler.userIdentity, handler.createGame)

	body := `{"language": "en"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d, body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["id"] != float64(100) {
		t.Errorf("expected id 100, got: %v", response["id"])
	}
}

func TestCreateGame_Unauthorized(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 0, service.ErrUnauthorized
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games", handler.userIdentity, handler.createGame)

	body := `{"language": "en"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got: %d", w.Code)
	}
}

func TestCreateGame_InvalidJSON(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games", handler.userIdentity, handler.createGame)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/api/games", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestCreateGame_UnsupportedLanguage(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games", handler.userIdentity, handler.createGame)

	body := `{"language": "xyz"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestCreateGame_ServiceError(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		CreateGameFn: func(ctx context.Context, userId int64, lang string) (int64, error) {
			return 0, errors.New("service error")
		},
	}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games", handler.userIdentity, handler.createGame)

	body := `{"language": "en"}`
	req := httptest.NewRequest(http.MethodPost, "/api/games", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got: %d", w.Code)
	}
}

func TestGetAllGames_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		GetAllGamesFn: func(ctx context.Context, userId int64) ([]model.Game, error) {
			return []model.Game{
				{Id: 1, UserId: userId, Status: "active", Language: "en", StartedAt: time.Now()},
				{Id: 2, UserId: userId, Status: "won", Language: "en", StartedAt: time.Now()},
			}, nil
		},
	}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games", handler.userIdentity, handler.getAllGames)

	req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response getAllGamesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response.Data) != 2 {
		t.Errorf("expected 2 games, got: %d", len(response.Data))
	}
}

func TestGetAllGames_Unauthorized(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 0, service.ErrUnauthorized
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games", handler.userIdentity, handler.getAllGames)

	req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got: %d", w.Code)
	}
}

func TestGetAllGames_EmptyList(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		GetAllGamesFn: func(ctx context.Context, userId int64) ([]model.Game, error) {
			return []model.Game{}, nil
		},
	}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games", handler.userIdentity, handler.getAllGames)

	req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response getAllGamesResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response.Data) != 0 {
		t.Errorf("expected 0 games, got: %d", len(response.Data))
	}
}

func TestGetActiveGame_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		GetActiveGameFn: func(ctx context.Context, userId int64) (*model.Game, error) {
			return &model.Game{Id: 1, UserId: userId, Status: "active", Language: "en", StartedAt: time.Now()}, nil
		},
	}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games/active", handler.userIdentity, handler.getActiveGame)

	req := httptest.NewRequest(http.MethodGet, "/api/games/active", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestGetActiveGame_NoActiveGame(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		GetActiveGameFn: func(ctx context.Context, userId int64) (*model.Game, error) {
			return nil, nil
		},
	}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games/active", handler.userIdentity, handler.getActiveGame)

	req := httptest.NewRequest(http.MethodGet, "/api/games/active", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response getGameByIdResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Data != nil {
		t.Errorf("expected nil data, got: %v", response.Data)
	}
}

func TestGetGameById_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		GetGameByIdFn: func(ctx context.Context, userId, gameId int64) (*model.Game, error) {
			if gameId != 123 {
				t.Errorf("expected gameId 123, got: %d", gameId)
			}
			return &model.Game{Id: gameId, UserId: userId, Status: "active", Language: "en", StartedAt: time.Now()}, nil
		},
	}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games/:id", handler.userIdentity, handler.getGameById)

	req := httptest.NewRequest(http.MethodGet, "/api/games/123", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestGetGameById_InvalidId(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games/:id", handler.userIdentity, handler.getGameById)

	req := httptest.NewRequest(http.MethodGet, "/api/games/invalid", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestGetGameById_NotFound(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		GetGameByIdFn: func(ctx context.Context, userId, gameId int64) (*model.Game, error) {
			return nil, service.ErrGameNotFound
		},
	}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/games/:id", handler.userIdentity, handler.getGameById)

	req := httptest.NewRequest(http.MethodGet, "/api/games/999", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got: %d", w.Code)
	}
}

func TestAbandonGame_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		AbandonGameFn: func(ctx context.Context, userId, gameId int64) error {
			if gameId != 123 || userId != 42 {
				t.Errorf("unexpected params: userId=%d, gameId=%d", userId, gameId)
			}
			return nil
		},
	}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/abandon", handler.userIdentity, handler.abandonGame)

	req := httptest.NewRequest(http.MethodPost, "/api/games/123/abandon", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response statusResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Status != "abandoned" {
		t.Errorf("expected status 'abandoned', got: %s", response.Status)
	}
}

func TestAbandonGame_InvalidId(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/abandon", handler.userIdentity, handler.abandonGame)

	req := httptest.NewRequest(http.MethodPost, "/api/games/invalid/abandon", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestAbandonGame_NotFound(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		AbandonGameFn: func(ctx context.Context, userId, gameId int64) error {
			return service.ErrGameNotFound
		},
	}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/abandon", handler.userIdentity, handler.abandonGame)

	req := httptest.NewRequest(http.MethodPost, "/api/games/999/abandon", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got: %d", w.Code)
	}
}

func TestAbandonGame_GameNotActive(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockGame := &servicemocks.MockGameService{
		AbandonGameFn: func(ctx context.Context, userId, gameId int64) error {
			return service.ErrGameNotActive
		},
	}

	handler := setupGameHandler(mockAuth, mockGame)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/api/games/:id/abandon", handler.userIdentity, handler.abandonGame)

	req := httptest.NewRequest(http.MethodPost, "/api/games/123/abandon", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}
