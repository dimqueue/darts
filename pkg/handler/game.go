package handler

import (
	"net/http"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/gin-gonic/gin"
)

// add details
type CreateGameInput struct {
	Language string `json:"language" binding:"required"`
}

// @Summary      CreateGame
// @Security     ApiKeyAuth
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
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
	}

	var input CreateGameInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err = h.services.Game.CreateGame(userId, input.Language)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": userId,
	})
}

type getAllGamesResponse struct {
	Data []model.Game `json:"data"`
}

// @Summary      GetAllGames
// @Security     ApiKeyAuth
// @Description  getAllGames
// @Tags         game
// @Accept       json
// @Produce      json
// @Success      200  {integer}   integer 1
// @Failure      500  {object}  errorResponse
// @Router       /api/games/ [get]
func (h *Handler) getAllGames(c *gin.Context) {
	idUser, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	games, err := h.services.GetAllGames(idUser)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, getAllGamesResponse{
		Data: games,
	})

}

func (h *Handler) getGameById(c *gin.Context) {

}

func (h *Handler) updateGame(c *gin.Context) {

}

func (h *Handler) deleteGame(c *gin.Context) {

}
