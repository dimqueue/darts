package handler

import (
	"github.com/dimqueue/darts/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createGame(c *gin.Context) {
	//like list
	//id, ok := c.Get(userCtx)
	//if !ok {
	//	newErrorResponse(c, http.StatusInternalServerError, "user not found")
	//}
	//var input

	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user not found")
	}

	var input model.Game

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Game.CreateGame(id.(int), input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllGames(c *gin.Context) {

}

func (h *Handler) getGameById(c *gin.Context) {

}

func (h *Handler) updateGame(c *gin.Context) {

}

func (h *Handler) deleteGame(c *gin.Context) {

}
