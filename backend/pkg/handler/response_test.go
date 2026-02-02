package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dimqueue/darts/pkg/service"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHandleError_DomainErrors(t *testing.T) {
	testCases := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   string
	}{
		{"ErrGameNotFound", service.ErrGameNotFound, http.StatusNotFound, CodeGameNotFound},
		{"ErrGameNotActive", service.ErrGameNotActive, http.StatusBadRequest, CodeGameNotActive},
		{"ErrGameExpired", service.ErrGameExpired, http.StatusBadRequest, CodeGameExpired},
		{"ErrWordNotFound", service.ErrWordNotFound, http.StatusOK, CodeWordNotFound},
		{"ErrWordAlreadyUsed", service.ErrWordAlreadyUsed, http.StatusBadRequest, CodeWordAlreadyUsed},
		{"ErrUnauthorized", service.ErrUnauthorized, http.StatusUnauthorized, CodeUnauthorized},
		{"ErrUserExists", service.ErrUserExists, http.StatusConflict, CodeUserExists},
		{"ErrForbidden", service.ErrForbidden, http.StatusForbidden, CodeForbidden},
		{"ErrProfilePrivate", service.ErrProfilePrivate, http.StatusForbidden, CodeProfilePrivate},
		{"ErrProfileNotFound", service.ErrProfileNotFound, http.StatusNotFound, CodeProfileNotFound},
		{"ErrNotFound", service.ErrNotFound, http.StatusNotFound, CodeNotFound},
		{"ErrInvalidInput", service.ErrInvalidInput, http.StatusBadRequest, CodeValidation},
		{"ErrComputeService", service.ErrComputeService, http.StatusServiceUnavailable, CodeComputeError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			handleError(c, tc.err)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got: %d", tc.expectedStatus, w.Code)
			}

			var resp errorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			if resp.Code != tc.expectedCode {
				t.Errorf("expected code %s, got: %s", tc.expectedCode, resp.Code)
			}
		})
	}
}

func TestHandleError_HandlerErrors(t *testing.T) {
	testCases := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   string
	}{
		{"ErrBadRequest", ErrBadRequest("bad request"), http.StatusBadRequest, CodeBadRequest},
		{"ErrUnauthorized", ErrUnauthorized("unauthorized"), http.StatusUnauthorized, CodeUnauthorized},
		{"ErrValidation", ErrValidation("validation error"), http.StatusBadRequest, CodeValidation},
		{"ErrNotFound", ErrNotFound("not found"), http.StatusNotFound, CodeNotFound},
		{"ErrForbidden", ErrForbidden("forbidden"), http.StatusForbidden, CodeForbidden},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			handleError(c, tc.err)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got: %d", tc.expectedStatus, w.Code)
			}

			var resp errorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			if resp.Code != tc.expectedCode {
				t.Errorf("expected code %s, got: %s", tc.expectedCode, resp.Code)
			}
		})
	}
}

func TestHandleError_UnknownError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	unknownErr := errors.New("unknown error")
	handleError(c, unknownErr)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got: %d", w.Code)
	}

	var resp errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Code != CodeInternalError {
		t.Errorf("expected code %s, got: %s", CodeInternalError, resp.Code)
	}
}

func TestHandleError_IncludesRequestID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	expectedRequestID := "test-request-id-123"
	c.Set(RequestIDKey, expectedRequestID)

	handleError(c, service.ErrGameNotFound)

	var resp errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.RequestID != expectedRequestID {
		t.Errorf("expected requestId %s, got: %s", expectedRequestID, resp.RequestID)
	}
}

func TestHandleValidationError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	expectedMessage := "field 'email' is required"
	handleValidationError(c, expectedMessage)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got: %d", w.Code)
	}

	var resp errorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Code != CodeValidation {
		t.Errorf("expected code %s, got: %s", CodeValidation, resp.Code)
	}

	if resp.Message != expectedMessage {
		t.Errorf("expected message %s, got: %s", expectedMessage, resp.Message)
	}
}

func TestErrBadRequest(t *testing.T) {
	err := ErrBadRequest("test message")
	he, ok := err.(*handlerError)
	if !ok {
		t.Fatal("expected handlerError type")
	}

	if he.status != http.StatusBadRequest {
		t.Errorf("expected status %d, got: %d", http.StatusBadRequest, he.status)
	}
	if he.code != CodeBadRequest {
		t.Errorf("expected code %s, got: %s", CodeBadRequest, he.code)
	}
	if he.message != "test message" {
		t.Errorf("expected message 'test message', got: %s", he.message)
	}
}

func TestErrUnauthorized_Handler(t *testing.T) {
	err := ErrUnauthorized("test message")
	he, ok := err.(*handlerError)
	if !ok {
		t.Fatal("expected handlerError type")
	}

	if he.status != http.StatusUnauthorized {
		t.Errorf("expected status %d, got: %d", http.StatusUnauthorized, he.status)
	}
	if he.code != CodeUnauthorized {
		t.Errorf("expected code %s, got: %s", CodeUnauthorized, he.code)
	}
}

func TestErrValidation(t *testing.T) {
	err := ErrValidation("test message")
	he, ok := err.(*handlerError)
	if !ok {
		t.Fatal("expected handlerError type")
	}

	if he.status != http.StatusBadRequest {
		t.Errorf("expected status %d, got: %d", http.StatusBadRequest, he.status)
	}
	if he.code != CodeValidation {
		t.Errorf("expected code %s, got: %s", CodeValidation, he.code)
	}
}

func TestErrNotFound_Handler(t *testing.T) {
	err := ErrNotFound("test message")
	he, ok := err.(*handlerError)
	if !ok {
		t.Fatal("expected handlerError type")
	}

	if he.status != http.StatusNotFound {
		t.Errorf("expected status %d, got: %d", http.StatusNotFound, he.status)
	}
	if he.code != CodeNotFound {
		t.Errorf("expected code %s, got: %s", CodeNotFound, he.code)
	}
}

func TestErrForbidden_Handler(t *testing.T) {
	err := ErrForbidden("test message")
	he, ok := err.(*handlerError)
	if !ok {
		t.Fatal("expected handlerError type")
	}

	if he.status != http.StatusForbidden {
		t.Errorf("expected status %d, got: %d", http.StatusForbidden, he.status)
	}
	if he.code != CodeForbidden {
		t.Errorf("expected code %s, got: %s", CodeForbidden, he.code)
	}
}

func TestHandlerError_ErrorMethod(t *testing.T) {
	err := ErrBadRequest("custom message")
	if err.Error() != "custom message" {
		t.Errorf("expected 'custom message', got: %s", err.Error())
	}
}
