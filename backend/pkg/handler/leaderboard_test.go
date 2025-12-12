package handler

import (
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

func setupLeaderboardHandler(mockAuth *servicemocks.MockAuthService, mockLeaderboard *servicemocks.MockLeaderboardService) *Handler {
	return &Handler{
		services: &service.Service{
			Authorization: mockAuth,
			Leaderboard:   mockLeaderboard,
		},
		validator: validation.New(),
	}
}

func TestGetLeaderboard_Global_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockLeaderboard := &servicemocks.MockLeaderboardService{
		GetLeaderboardFn: func(ctx context.Context, query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
			if query.Type != "global" {
				t.Errorf("expected type 'global', got: %s", query.Type)
			}
			return &model.LeaderboardResponse{
				LeaderboardType: query.Type,
				Users: []model.LeaderboardUser{
					{Rank: 1, UserId: 1, Username: "top1", TotalScore: 1000},
					{Rank: 2, UserId: 2, Username: "top2", TotalScore: 900},
				},
				Total: 100,
			}, nil
		},
		GetUserRankFn: func(ctx context.Context, userId int64, query model.LeaderboardQuery) (*int, error) {
			rank := 15
			return &rank, nil
		},
	}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard", handler.userIdentity, handler.getLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=global", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response model.LeaderboardResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.LeaderboardType != "global" {
		t.Errorf("expected type 'global', got: %s", response.LeaderboardType)
	}
	if len(response.Users) != 2 {
		t.Errorf("expected 2 users, got: %d", len(response.Users))
	}
	if response.CurrentUserRank == nil || *response.CurrentUserRank != 15 {
		t.Errorf("expected current user rank 15, got: %v", response.CurrentUserRank)
	}
}

func TestGetLeaderboard_Weekly_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockLeaderboard := &servicemocks.MockLeaderboardService{
		GetLeaderboardFn: func(ctx context.Context, query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
			if query.Type != "weekly" {
				t.Errorf("expected type 'weekly', got: %s", query.Type)
			}
			return &model.LeaderboardResponse{
				LeaderboardType: query.Type,
				Users:           []model.LeaderboardUser{},
				Total:           0,
			}, nil
		},
		GetUserRankFn: func(ctx context.Context, userId int64, query model.LeaderboardQuery) (*int, error) {
			return nil, nil
		},
	}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard", handler.userIdentity, handler.getLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=weekly", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestGetLeaderboard_Monthly_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockLeaderboard := &servicemocks.MockLeaderboardService{
		GetLeaderboardFn: func(ctx context.Context, query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
			return &model.LeaderboardResponse{
				LeaderboardType: "monthly",
				Users:           []model.LeaderboardUser{},
				Total:           0,
			}, nil
		},
		GetUserRankFn: func(ctx context.Context, userId int64, query model.LeaderboardQuery) (*int, error) {
			return nil, nil
		},
	}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard", handler.userIdentity, handler.getLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=monthly", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestGetLeaderboard_Daily_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockLeaderboard := &servicemocks.MockLeaderboardService{
		GetLeaderboardFn: func(ctx context.Context, query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
			return &model.LeaderboardResponse{
				LeaderboardType: "daily",
				Users:           []model.LeaderboardUser{},
				Total:           0,
			}, nil
		},
		GetUserRankFn: func(ctx context.Context, userId int64, query model.LeaderboardQuery) (*int, error) {
			return nil, nil
		},
	}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard", handler.userIdentity, handler.getLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=daily", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestGetLeaderboard_MissingType(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockLeaderboard := &servicemocks.MockLeaderboardService{}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard", handler.userIdentity, handler.getLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestGetLeaderboard_InvalidType(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockLeaderboard := &servicemocks.MockLeaderboardService{}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard", handler.userIdentity, handler.getLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=invalid", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}
}

func TestGetLeaderboard_WithLanguageFilter(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockLeaderboard := &servicemocks.MockLeaderboardService{
		GetLeaderboardFn: func(ctx context.Context, query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
			if query.Language == nil || *query.Language != "en" {
				t.Errorf("expected language 'en', got: %v", query.Language)
			}
			lang := "en"
			return &model.LeaderboardResponse{
				LeaderboardType: query.Type,
				Language:        &lang,
				Users:           []model.LeaderboardUser{},
				Total:           0,
			}, nil
		},
		GetUserRankFn: func(ctx context.Context, userId int64, query model.LeaderboardQuery) (*int, error) {
			return nil, nil
		},
	}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard", handler.userIdentity, handler.getLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=global&language=en", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestGetLeaderboard_WithPagination(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockLeaderboard := &servicemocks.MockLeaderboardService{
		GetLeaderboardFn: func(ctx context.Context, query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
			if query.Limit != 20 {
				t.Errorf("expected limit 20, got: %d", query.Limit)
			}
			if query.Offset != 10 {
				t.Errorf("expected offset 10, got: %d", query.Offset)
			}
			return &model.LeaderboardResponse{
				LeaderboardType: query.Type,
				Users:           []model.LeaderboardUser{},
				Total:           100,
			}, nil
		},
		GetUserRankFn: func(ctx context.Context, userId int64, query model.LeaderboardQuery) (*int, error) {
			return nil, nil
		},
	}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard", handler.userIdentity, handler.getLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=global&limit=20&offset=10", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
}

func TestGetLeaderboard_ServiceError(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockLeaderboard := &servicemocks.MockLeaderboardService{
		GetLeaderboardFn: func(ctx context.Context, query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
			return nil, errors.New("database error")
		},
	}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard", handler.userIdentity, handler.getLeaderboard)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard?type=global", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got: %d", w.Code)
	}
}

// --- getMyRank Tests ---

func TestGetMyRank_Success(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	globalRank := 10
	dailyRank := 5
	weeklyRank := 8
	monthlyRank := 12
	mockLeaderboard := &servicemocks.MockLeaderboardService{
		GetAllUserRanksFn: func(ctx context.Context, userId int64) (*model.UserRanks, error) {
			if userId != 42 {
				t.Errorf("expected userId 42, got: %d", userId)
			}
			return &model.UserRanks{
				GlobalRank:  &globalRank,
				DailyRank:   &dailyRank,
				WeeklyRank:  &weeklyRank,
				MonthlyRank: &monthlyRank,
			}, nil
		},
	}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard/my-rank", handler.userIdentity, handler.getMyRank)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard/my-rank", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response model.UserRanks
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.GlobalRank == nil || *response.GlobalRank != 10 {
		t.Errorf("expected global rank 10, got: %v", response.GlobalRank)
	}
	if response.DailyRank == nil || *response.DailyRank != 5 {
		t.Errorf("expected daily rank 5, got: %v", response.DailyRank)
	}
}

func TestGetMyRank_Unauthorized(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 0, service.ErrUnauthorized
		},
	}
	mockLeaderboard := &servicemocks.MockLeaderboardService{}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard/my-rank", handler.userIdentity, handler.getMyRank)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard/my-rank", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got: %d", w.Code)
	}
}

func TestGetMyRank_NoRanks(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockLeaderboard := &servicemocks.MockLeaderboardService{
		GetAllUserRanksFn: func(ctx context.Context, userId int64) (*model.UserRanks, error) {
			return &model.UserRanks{}, nil
		},
	}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard/my-rank", handler.userIdentity, handler.getMyRank)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard/my-rank", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}

	var response model.UserRanks
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.GlobalRank != nil {
		t.Errorf("expected nil global rank, got: %v", response.GlobalRank)
	}
}

func TestGetMyRank_ServiceError(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 42, nil
		},
	}
	mockLeaderboard := &servicemocks.MockLeaderboardService{
		GetAllUserRanksFn: func(ctx context.Context, userId int64) (*model.UserRanks, error) {
			return nil, errors.New("database error")
		},
	}

	handler := setupLeaderboardHandler(mockAuth, mockLeaderboard)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/api/leaderboard/my-rank", handler.userIdentity, handler.getMyRank)

	req := httptest.NewRequest(http.MethodGet, "/api/leaderboard/my-rank", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got: %d", w.Code)
	}
}

func TestParseIntQueryWithBounds(t *testing.T) {
	tests := []struct {
		name         string
		queryValue   string
		defaultValue int
		min          int
		max          int
		expected     int
	}{
		{"default when empty", "", 50, 1, 100, 50},
		{"normal value", "25", 50, 1, 100, 25},
		{"below min", "0", 50, 1, 100, 1},
		{"above max", "200", 50, 1, 100, 100},
		{"invalid string", "abc", 50, 1, 100, 50},
		{"negative", "-5", 50, 1, 100, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			if tc.queryValue != "" {
				c.Request = httptest.NewRequest(http.MethodGet, "/?value="+tc.queryValue, nil)
			} else {
				c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
			}

			result := parseIntQueryWithBounds(c, "value", tc.defaultValue, tc.min, tc.max)

			if result != tc.expected {
				t.Errorf("expected %d, got: %d", tc.expected, result)
			}
		})
	}
}
