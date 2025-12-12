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

func setupProfileHandler(mockAuth *servicemocks.MockAuthService, mockProfile *servicemocks.MockProfileService) *Handler {
	return &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Profile:       mockProfile,
		},
		validator: validation.New(),
	}
}

func TestGetMyProfile_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{
		GetProfileSummaryFn: func(ctx context.Context, userId int64) (*model.UserProfileSummary, error) {
			return &model.UserProfileSummary{
				Id:          userId,
				Username:    "testuser",
				Name:        "Test User",
				TotalGames:  100,
				TotalWins:   75,
				MemberSince: time.Now(),
			}, nil
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/profile", handler.userIdentity, handler.getMyProfile)

	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response model.UserProfileSummary
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Username != "testuser" {
		t.Errorf("expected username 'testuser', got: %s", response.Username)
	}
}

func TestGetMyProfile_Unauthorized(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 0, service.ErrUnauthorized
		},
	}
	mockProfile := &servicemocks.MockProfileService{}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/profile", handler.userIdentity, handler.getMyProfile)

	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got: %d", w.Code)
	}
}

func TestGetMyProfile_NotFound(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{
		GetProfileSummaryFn: func(ctx context.Context, userId int64) (*model.UserProfileSummary, error) {
			return nil, service.ErrProfileNotFound
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/profile", handler.userIdentity, handler.getMyProfile)

	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got: %d", w.Code)
	}
}

func TestGetProfileByUsername_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{}
	mockProfile := &servicemocks.MockProfileService{
		GetProfileByUsernameFn: func(ctx context.Context, username string) (*model.UserProfileSummary, error) {
			return &model.UserProfileSummary{
				Id:                42,
				Username:          username,
				Name:              "Public User",
				ShowProfilePublic: true,
			}, nil
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/public/profile/:username", handler.getProfileByUsername)

	req := httptest.NewRequest(http.MethodGet, "/public/profile/publicuser", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestGetProfileByUsername_NotFound(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{}
	mockProfile := &servicemocks.MockProfileService{
		GetProfileByUsernameFn: func(ctx context.Context, username string) (*model.UserProfileSummary, error) {
			return nil, service.ErrProfileNotFound
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/public/profile/:username", handler.getProfileByUsername)

	req := httptest.NewRequest(http.MethodGet, "/public/profile/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got: %d", w.Code)
	}
}

func TestGetProfileByUsername_PrivateProfile(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{}
	mockProfile := &servicemocks.MockProfileService{
		GetProfileByUsernameFn: func(ctx context.Context, username string) (*model.UserProfileSummary, error) {
			return &model.UserProfileSummary{
				Id:                42,
				Username:          username,
				ShowProfilePublic: false,
			}, nil
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/public/profile/:username", handler.getProfileByUsername)

	req := httptest.NewRequest(http.MethodGet, "/public/profile/privateuser", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got: %d", w.Code)
	}
}

func TestUpdateMyProfile_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{
		UpdateProfileFn: func(ctx context.Context, userId int64, input model.UpdateProfileInput) error {
			return nil
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.PUT("/api/profile", handler.userIdentity, handler.updateMyProfile)

	body := `{"bio": "New bio", "timezone": "UTC"}`
	req := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestUpdateMyProfile_InvalidJSON(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.PUT("/api/profile", handler.userIdentity, handler.updateMyProfile)

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestUpdateMyProfile_ServiceError(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{
		UpdateProfileFn: func(ctx context.Context, userId int64, input model.UpdateProfileInput) error {
			return errors.New("update failed")
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.PUT("/api/profile", handler.userIdentity, handler.updateMyProfile)

	body := `{"bio": "New bio"}`
	req := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got: %d", w.Code)
	}
}

func TestGetMySettings_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{
		GetSettingsFn: func(ctx context.Context, userId int64) (*model.UserSettings, error) {
			return &model.UserSettings{
				UserId:            userId,
				PreferredLanguage: "en",
				Theme:             "dark",
				SoundEnabled:      true,
			}, nil
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/profile/settings", handler.userIdentity, handler.getMySettings)

	req := httptest.NewRequest(http.MethodGet, "/api/profile/settings", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response model.UserSettings
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Theme != "dark" {
		t.Errorf("expected theme 'dark', got: %s", response.Theme)
	}
}

func TestGetMySettings_Error(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{
		GetSettingsFn: func(ctx context.Context, userId int64) (*model.UserSettings, error) {
			return nil, service.ErrNotFound
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/profile/settings", handler.userIdentity, handler.getMySettings)

	req := httptest.NewRequest(http.MethodGet, "/api/profile/settings", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got: %d", w.Code)
	}
}

func TestUpdateMySettings_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{
		UpdateSettingsFn: func(ctx context.Context, userId int64, input model.UpdateSettingsInput) error {
			return nil
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.PUT("/api/profile/settings", handler.userIdentity, handler.updateMySettings)

	body := `{"theme": "light", "sound_enabled": false}`
	req := httptest.NewRequest(http.MethodPut, "/api/profile/settings", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestUpdateMySettings_InvalidJSON(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.PUT("/api/profile/settings", handler.userIdentity, handler.updateMySettings)

	body := `{invalid}`
	req := httptest.NewRequest(http.MethodPut, "/api/profile/settings", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestGetMyStatistics_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{
		GetStatisticsFn: func(ctx context.Context, userId int64) (*model.UserStatistics, error) {
			return &model.UserStatistics{
				UserId:           userId,
				TotalGames:       100,
				TotalWins:        80,
				CurrentWinStreak: 5,
				BestWinStreak:    15,
				TotalScore:       5000,
			}, nil
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/profile/statistics", handler.userIdentity, handler.getMyStatistics)

	req := httptest.NewRequest(http.MethodGet, "/api/profile/statistics", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response model.UserStatistics
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.TotalGames != 100 {
		t.Errorf("expected 100 total games, got: %d", response.TotalGames)
	}
}

func TestGetMyLanguageStats_AllLanguages(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{
		GetAllLanguageStatsFn: func(ctx context.Context, userId int64) ([]model.UserLanguageStats, error) {
			return []model.UserLanguageStats{
				{UserId: userId, Language: "en", GamesPlayed: 50},
				{UserId: userId, Language: "ua", GamesPlayed: 30},
			}, nil
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/profile/statistics/languages", handler.userIdentity, handler.getMyLanguageStats)

	req := httptest.NewRequest(http.MethodGet, "/api/profile/statistics/languages", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestGetMyLanguageStats_SpecificLanguage(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{
		GetLanguageStatsFn: func(ctx context.Context, userId int64, language string) (*model.UserLanguageStats, error) {
			if language != "en" {
				t.Errorf("expected language 'en', got: %s", language)
			}
			return &model.UserLanguageStats{
				UserId:      userId,
				Language:    language,
				GamesPlayed: 50,
				GamesWon:    40,
			}, nil
		},
	}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/profile/statistics/languages", handler.userIdentity, handler.getMyLanguageStats)

	req := httptest.NewRequest(http.MethodGet, "/api/profile/statistics/languages?language=en", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestGetMyLanguageStats_InvalidLanguage(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockProfile := &servicemocks.MockProfileService{}

	handler := setupProfileHandler(mockAuth, mockProfile)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/profile/statistics/languages", handler.userIdentity, handler.getMyLanguageStats)

	req := httptest.NewRequest(http.MethodGet, "/api/profile/statistics/languages?language=xyz", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}
