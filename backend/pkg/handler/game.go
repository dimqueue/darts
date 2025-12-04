package handler

import (
	"net/http"
	"strconv"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/gin-gonic/gin"
)

type CreateGameInput struct {
	Language string `json:"language" binding:"required"`
}

// @Summary      CreateGame
// @Security     ApiKeyAuth
// @Description  create game
// @Tags         game
// @Accept       json
// @Produce      json
// @Param        input body CreateGameInput true "game info"
// @Success      200  {integer}   integer 1
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/games/ [post]
func (h *Handler) createGame(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	var input CreateGameInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.ValidateLanguage(input.Language); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	gameId, err := h.services.Game.CreateGame(userId, input.Language)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": gameId,
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
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	games, err := h.services.GetAllGames(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllGamesResponse{
		Data: games,
	})
}

type getGameByIdResponse struct {
	Data *model.Game `json:"data"`
}

func (h *Handler) getGameById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	gameId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid game id")
		return
	}

	game, err := h.services.GetGameById(userId, gameId)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, getGameByIdResponse{
		Data: game,
	})
}
