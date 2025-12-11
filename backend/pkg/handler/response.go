package handler

import (
	"errors"
	"net/http"

	"github.com/dimqueue/darts/pkg/logger"
	"github.com/dimqueue/darts/pkg/service"
	"github.com/gin-gonic/gin"
)

const (
	CodeInternalError   = "INTERNAL_ERROR"
	CodeBadRequest      = "BAD_REQUEST"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
	CodeNotFound        = "NOT_FOUND"
	CodeValidation      = "VALIDATION_ERROR"
	CodeGameNotFound    = "GAME_NOT_FOUND"
	CodeGameNotActive   = "GAME_NOT_ACTIVE"
	CodeGameExpired     = "GAME_EXPIRED"
	CodeWordNotFound    = "WORD_NOT_FOUND"
	CodeWordAlreadyUsed = "WORD_ALREADY_USED"
	CodeWordTooFar      = "WORD_TOO_FAR"
	CodeComputeError    = "COMPUTE_SERVICE_ERROR"
	CodeProfilePrivate  = "PROFILE_PRIVATE"
	CodeProfileNotFound = "PROFILE_NOT_FOUND"
	CodeUserExists      = "USER_EXISTS"
)

type httpError struct {
	Status  int
	Code    string
	Message string
}

var errorMap = map[error]httpError{
	// Game errors
	service.ErrGameNotFound:  {http.StatusNotFound, CodeGameNotFound, "Game not found"},
	service.ErrGameNotActive: {http.StatusBadRequest, CodeGameNotActive, "Game is not active"},
	service.ErrGameExpired:   {http.StatusBadRequest, CodeGameExpired, "Game has expired"},

	// Guess errors
	service.ErrWordNotFound:    {http.StatusOK, CodeWordNotFound, "Word not found in vocabulary"},
	service.ErrWordAlreadyUsed: {http.StatusBadRequest, CodeWordAlreadyUsed, "This word has already been guessed"},
	service.ErrWordTooFar:      {http.StatusOK, CodeWordTooFar, "Word is too far from target"},

	// Auth errors
	service.ErrUnauthorized: {http.StatusUnauthorized, CodeUnauthorized, "Invalid credentials"},
	service.ErrUserExists:   {http.StatusConflict, CodeUserExists, "User already exists"},

	// Access errors
	service.ErrForbidden:       {http.StatusForbidden, CodeForbidden, "Access denied"},
	service.ErrProfilePrivate:  {http.StatusForbidden, CodeProfilePrivate, "This profile is private"},
	service.ErrProfileNotFound: {http.StatusNotFound, CodeProfileNotFound, "Profile not found"},

	// Generic errors
	service.ErrNotFound:     {http.StatusNotFound, CodeNotFound, "Resource not found"},
	service.ErrInvalidInput: {http.StatusBadRequest, CodeValidation, "Invalid input"},

	// External service errors
	service.ErrComputeService: {http.StatusServiceUnavailable, CodeComputeError, "Service temporarily unavailable"},
}

type errorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func handleError(c *gin.Context, err error) {
	requestID := GetRequestID(c)

	// Check if it's a handler-level error (input validation, parsing)
	if he, ok := err.(*handlerError); ok {
		c.AbortWithStatusJSON(he.status, errorResponse{
			Code:      he.code,
			Message:   he.message,
			RequestID: requestID,
		})
		return
	}

	for domainErr, httpErr := range errorMap {
		if errors.Is(err, domainErr) {
			c.AbortWithStatusJSON(httpErr.Status, errorResponse{
				Code:      httpErr.Code,
				Message:   httpErr.Message,
				RequestID: requestID,
			})
			return
		}
	}

	if !service.IsLogged(err) {
		log := logger.FromContext(c.Request.Context())
		if userId, exists := c.Get(userCtx); exists {
			log = log.With(logger.FieldUserID, userId)
		}
		log.Error("unhandled error", logger.FieldErrorCode, CodeInternalError, logger.FieldError, err)
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{
		Code:      CodeInternalError,
		Message:   "An unexpected error occurred",
		RequestID: requestID,
	})
}

func handleValidationError(c *gin.Context, message string) {
	requestID := GetRequestID(c)
	log := logger.FromContext(c.Request.Context())

	if userId, exists := c.Get(userCtx); exists {
		log = log.With(logger.FieldUserID, userId)
	}

	log.Warn(message, logger.FieldErrorCode, CodeValidation)

	c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{
		Code:      CodeValidation,
		Message:   message,
		RequestID: requestID,
	})
}

func ErrBadRequest(message string) error {
	return &handlerError{status: http.StatusBadRequest, code: CodeBadRequest, message: message}
}

func ErrUnauthorized(message string) error {
	return &handlerError{status: http.StatusUnauthorized, code: CodeUnauthorized, message: message}
}

func ErrValidation(message string) error {
	return &handlerError{status: http.StatusBadRequest, code: CodeValidation, message: message}
}

func ErrNotFound(message string) error {
	return &handlerError{status: http.StatusNotFound, code: CodeNotFound, message: message}
}

func ErrForbidden(message string) error {
	return &handlerError{status: http.StatusForbidden, code: CodeForbidden, message: message}
}

type handlerError struct {
	status  int
	code    string
	message string
}

func (e *handlerError) Error() string {
	return e.message
}
