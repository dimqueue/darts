package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/siruspen/logrus"
)

type err struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, err{message})
}
