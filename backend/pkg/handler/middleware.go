package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dimqueue/darts/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	RequestIDHeader     = "X-Request-Id"
	RequestIDKey        = "requestId"
)

func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" || !isValidUUID(requestID) {
			requestID = uuid.New().String()
		}

		c.Set(RequestIDKey, requestID)

		ctx := context.WithValue(c.Request.Context(), logger.RequestIDCtxKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(RequestIDKey); exists {
		if requestID, ok := id.(string); ok {
			return requestID
		}
	}
	return ""
}

func GetRequestIDFromContext(ctx context.Context) string {
	if id := ctx.Value(logger.RequestIDCtxKey); id != nil {
		if requestID, ok := id.(string); ok {
			return requestID
		}
	}
	return ""
}

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)

	c.Next()
}

func getUserId(c *gin.Context) (int64, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt64, ok := id.(int64)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt64, nil
}
