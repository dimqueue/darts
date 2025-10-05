package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createUser(c *gin.Context) {
	_, ok := c.Get("userCtx")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

}

func (h *Handler) getUser(c *gin.Context) {

}

func (h *Handler) getUserById(c *gin.Context) {

}

func (h *Handler) updateUserById(c *gin.Context) {

}

func (h *Handler) deleteUserById(c *gin.Context) {

}

func (h *Handler) getAllUser(c *gin.Context) {

}
