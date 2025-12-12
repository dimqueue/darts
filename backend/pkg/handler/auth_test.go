package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/service"
	servicemocks "github.com/dimqueue/darts/pkg/service/mocks"
	"github.com/dimqueue/darts/pkg/validation"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAuthHandler(mockAuth *servicemocks.MockAuthService) *Handler {
	return &Handler{
		services: &service.Service{
			Authorization: mockAuth,
		},
		validator: validation.New(),
	}
}

func TestSignUp_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		CreateUserFn: func(ctx context.Context, user model.User) (int64, error) {
			return 42, nil
		},
	}

	handler := setupAuthHandler(mockAuth)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/auth/sign-up", handler.signUp)

	body := `{"name": "Test User", "username": "testuser", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["id"] != float64(42) {
		t.Errorf("expected id 42, got: %v", response["id"])
	}
}

func TestSignUp_InvalidJSON(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{}
	handler := setupAuthHandler(mockAuth)

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/auth/sign-up", handler.signUp)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestSignUp_MissingRequiredFields(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{}
	handler := setupAuthHandler(mockAuth)

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/auth/sign-up", handler.signUp)

	testCases := []struct {
		name string
		body string
	}{
		{"missing username", `{"name": "Test", "password": "pass123"}`},
		{"missing password", `{"name": "Test", "username": "testuser"}`},
		{"missing name", `{"username": "testuser", "password": "pass123"}`},
		{"empty body", `{}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got: %d", w.Code)
			}
		})
	}
}

func TestSignUp_UserAlreadyExists(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		CreateUserFn: func(ctx context.Context, user model.User) (int64, error) {
			return 0, service.ErrUserExists
		},
	}

	handler := setupAuthHandler(mockAuth)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/auth/sign-up", handler.signUp)

	body := `{"name": "Test User", "username": "existinguser", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status 409, got: %d", w.Code)
	}
}

func TestSignUp_InternalError(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		CreateUserFn: func(ctx context.Context, user model.User) (int64, error) {
			return 0, errors.New("database connection failed")
		},
	}

	handler := setupAuthHandler(mockAuth)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/auth/sign-up", handler.signUp)

	body := `{"name": "Test User", "username": "testuser", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got: %d", w.Code)
	}
}

func TestSignIn_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		GenerateTokenFn: func(ctx context.Context, username, password string) (string, error) {
			if username == "testuser" && password == "correctpassword" {
				return "valid-jwt-token", nil
			}
			return "", service.ErrUnauthorized
		},
	}

	handler := setupAuthHandler(mockAuth)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/auth/sign-in", handler.signIn)

	body := `{"username": "testuser", "password": "correctpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["token"] != "valid-jwt-token" {
		t.Errorf("expected token 'valid-jwt-token', got: %v", response["token"])
	}
}

func TestSignIn_InvalidJSON(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{}
	handler := setupAuthHandler(mockAuth)

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/auth/sign-in", handler.signIn)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestSignIn_MissingRequiredFields(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{}
	handler := setupAuthHandler(mockAuth)

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/auth/sign-in", handler.signIn)

	testCases := []struct {
		name string
		body string
	}{
		{"missing username", `{"password": "pass123"}`},
		{"missing password", `{"username": "testuser"}`},
		{"empty body", `{}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got: %d", w.Code)
			}
		})
	}
}

func TestSignIn_WrongPassword(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		GenerateTokenFn: func(ctx context.Context, username, password string) (string, error) {
			return "", service.ErrUnauthorized
		},
	}

	handler := setupAuthHandler(mockAuth)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/auth/sign-in", handler.signIn)

	body := `{"username": "testuser", "password": "wrongpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got: %d", w.Code)
	}
}

func TestSignIn_UserNotFound(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		GenerateTokenFn: func(ctx context.Context, username, password string) (string, error) {
			return "", service.ErrUnauthorized
		},
	}

	handler := setupAuthHandler(mockAuth)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/auth/sign-in", handler.signIn)

	body := `{"username": "nonexistent", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got: %d", w.Code)
	}
}

func TestSignIn_InternalError(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		GenerateTokenFn: func(ctx context.Context, username, password string) (string, error) {
			return "", errors.New("database connection failed")
		},
	}

	handler := setupAuthHandler(mockAuth)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.POST("/auth/sign-in", handler.signIn)

	body := `{"username": "testuser", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got: %d", w.Code)
	}
}
