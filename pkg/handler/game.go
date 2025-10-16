package handler

import (
	"github.com/dimqueue/darts/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

//fix input

// @Summary      CreateGame
// @Description  create game
// @Tags         game
// @Accept       json
// @Produce      json
// @Param        input body model.Game true "game info"
// @Success      200  {integer}   integer 1
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/games/ [post]
func (h *Handler) createGame(c *gin.Context) {
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
