package handler

import (
	"github.com/dimqueue/darts/pkg/logger"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Code      string `json:"code,omitempty"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

type successResponse struct {
	Data interface{} `json:"data"`
}

type paginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	requestID := GetRequestID(c)

	logger.WithContext(c.Request.Context()).WithField("status_code", statusCode).Error(message)

	c.AbortWithStatusJSON(statusCode, errorResponse{
		Message:   message,
		RequestID: requestID,
	})
}

func newSuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, successResponse{Data: data})
}

func newPaginatedResponse(c *gin.Context, data interface{}, total, page, perPage int) {
	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}
	c.JSON(200, paginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}
