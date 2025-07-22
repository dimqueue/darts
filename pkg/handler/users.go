package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createUser(c *gin.Context) {
	id, _ := c.Get("userCtx")
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
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
