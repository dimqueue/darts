package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dimqueue/darts/pkg/service"
	servicemocks "github.com/dimqueue/darts/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRequestIDMiddleware_GeneratesNewID(t *testing.T) {
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	requestID := w.Header().Get(RequestIDHeader)
	if requestID == "" {
		t.Error("expected X-Request-Id header to be set")
	}

	if _, err := uuid.Parse(requestID); err != nil {
		t.Errorf("expected valid UUID, got: %s", requestID)
	}
}

func TestRequestIDMiddleware_UsesProvidedID(t *testing.T) {
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	providedID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(RequestIDHeader, providedID)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	requestID := w.Header().Get(RequestIDHeader)
	if requestID != providedID {
		t.Errorf("expected request ID %s, got: %s", providedID, requestID)
	}
}

func TestRequestIDMiddleware_IgnoresInvalidUUID(t *testing.T) {
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(RequestIDHeader, "not-a-valid-uuid")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	requestID := w.Header().Get(RequestIDHeader)
	if requestID == "not-a-valid-uuid" {
		t.Error("should have generated new UUID instead of using invalid one")
	}

	if _, err := uuid.Parse(requestID); err != nil {
		t.Errorf("expected valid UUID, got: %s", requestID)
	}
}

func TestRequestIDMiddleware_SetsContextValue(t *testing.T) {
	router := gin.New()
	router.Use(RequestIDMiddleware())

	var capturedID string
	router.GET("/test", func(c *gin.Context) {
		capturedID = GetRequestID(c)
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if capturedID == "" {
		t.Error("GetRequestID should return the request ID from context")
	}

	responseID := w.Header().Get(RequestIDHeader)
	if capturedID != responseID {
		t.Errorf("context ID %s doesn't match response ID %s", capturedID, responseID)
	}
}

func TestGetRequestIDFromContext(t *testing.T) {
	router := gin.New()
	router.Use(RequestIDMiddleware())

	var capturedID string
	router.GET("/test", func(c *gin.Context) {
		capturedID = GetRequestIDFromContext(c.Request.Context())
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if capturedID == "" {
		t.Error("GetRequestIDFromContext should return the request ID")
	}
}

func TestGetRequestIDFromContext_Empty(t *testing.T) {
	ctx := context.Background()
	id := GetRequestIDFromContext(ctx)
	if id != "" {
		t.Errorf("expected empty string for context without request ID, got: %s", id)
	}
}

func TestUserIdentity_ValidToken(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			if token == "valid-token" {
				return 42, nil
			}
			return 0, errors.New("invalid token")
		},
	}

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
		},
	}

	router := gin.New()
	router.Use(RequestIDMiddleware())

	var capturedUserId int64
	router.GET("/test", handler.userIdentity, func(c *gin.Context) {
		userId, err := getUserId(c)
		if err != nil {
			t.Errorf("getUserId failed: %v", err)
		}
		capturedUserId = userId
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got: %d", w.Code)
	}
	if capturedUserId != 42 {
		t.Errorf("expected userId 42, got: %d", capturedUserId)
	}
}

func TestUserIdentity_MissingHeader(t *testing.T) {
	handler := &Handler{
		services: &service.Service{
			Authorization: &servicemocks.MockAuthService{},
		},
	}

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", handler.userIdentity, func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No Authorization header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got: %d", w.Code)
	}
}

func TestUserIdentity_InvalidHeaderFormat(t *testing.T) {
	handler := &Handler{
		services: &service.Service{
			Authorization: &servicemocks.MockAuthService{},
		},
	}

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", handler.userIdentity, func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	testCases := []struct {
		name   string
		header string
	}{
		{"no space", "Bearertoken"},
		{"too many parts", "Bearer token extra"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", tc.header)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected status 401, got: %d", w.Code)
			}
		})
	}
}

func TestUserIdentity_InvalidToken(t *testing.T) {
	mockAuth := &servicemocks.MockAuthService{
		ParseTokenFn: func(token string) (int64, error) {
			return 0, service.ErrUnauthorized
		},
	}

	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuth,
		},
	}

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/test", handler.userIdentity, func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got: %d", w.Code)
	}
}

func TestGetUserId_Success(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(userCtx, int64(123))

	userId, err := getUserId(c)
	if err != nil {
		t.Fatalf("getUserId failed: %v", err)
	}
	if userId != 123 {
		t.Errorf("expected userId 123, got: %d", userId)
	}
}

func TestGetUserId_NotSet(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	// userId not set in context

	_, err := getUserId(c)
	if err == nil {
		t.Error("expected error when userId not in context")
	}
}

func TestGetUserId_WrongType(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(userCtx, "not-an-int64") // Wrong type

	_, err := getUserId(c)
	if err == nil {
		t.Error("expected error when userId is wrong type")
	}
}

func TestIsValidUUID_Valid(t *testing.T) {
	validUUIDs := []string{
		uuid.New().String(),
		"550e8400-e29b-41d4-a716-446655440000",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	}

	for _, id := range validUUIDs {
		if !isValidUUID(id) {
			t.Errorf("expected %s to be valid UUID", id)
		}
	}
}

func TestIsValidUUID_Invalid(t *testing.T) {
	invalidUUIDs := []string{
		"",
		"not-a-uuid",
		"550e8400-e29b-41d4-a716", // Too short
		"550e8400-e29b-41d4-a716-446655440000-extra", // Too long
		"gggggggg-gggg-gggg-gggg-gggggggggggg",       // Invalid chars
	}

	for _, id := range invalidUUIDs {
		if isValidUUID(id) {
			t.Errorf("expected %s to be invalid UUID", id)
		}
	}
}

func TestGetRequestID_Success(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	expectedID := uuid.New().String()
	c.Set(RequestIDKey, expectedID)

	id := GetRequestID(c)
	if id != expectedID {
		t.Errorf("expected %s, got: %s", expectedID, id)
	}
}

func TestGetRequestID_NotSet(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	id := GetRequestID(c)
	if id != "" {
		t.Errorf("expected empty string, got: %s", id)
	}
}

func TestGetRequestID_WrongType(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(RequestIDKey, 12345) // Wrong type

	id := GetRequestID(c)
	if id != "" {
		t.Errorf("expected empty string for wrong type, got: %s", id)
	}
}
