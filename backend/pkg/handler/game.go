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
		handleError(c, ErrUnauthorized("user not found"))
		return
	}

	var input CreateGameInput
	if err := c.BindJSON(&input); err != nil {
		handleError(c, ErrBadRequest("invalid request body"))
		return
	}

	if err := h.validator.ValidateLanguage(input.Language); err != nil {
		handleError(c, ErrValidation(err.Error()))
		return
	}

	gameId, err := h.services.Game.CreateGame(c.Request.Context(), userId, input.Language)
	if err != nil {
		handleError(c, err)
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
		handleError(c, ErrUnauthorized("user not found"))
		return
	}

	games, err := h.services.GetAllGames(c.Request.Context(), userId)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, getAllGamesResponse{
		Data: games,
	})
}

type getGameByIdResponse struct {
	Data *model.Game `json:"data"`
}

func (h *Handler) getActiveGame(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		handleError(c, ErrUnauthorized("user not found"))
		return
	}

	game, err := h.services.Game.GetActiveGame(c.Request.Context(), userId)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, getGameByIdResponse{
		Data: game,
	})
}

func (h *Handler) getGameById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		handleError(c, ErrUnauthorized("user not found"))
		return
	}

	gameId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		handleError(c, ErrBadRequest("invalid game id"))
		return
	}

	game, err := h.services.GetGameById(c.Request.Context(), userId, gameId)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, getGameByIdResponse{
		Data: game,
	})
}

func (h *Handler) abandonGame(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		handleError(c, ErrUnauthorized("user not found"))
		return
	}

	gameId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		handleError(c, ErrBadRequest("invalid game id"))
		return
	}

	if err := h.services.Game.AbandonGame(c.Request.Context(), userId, gameId); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "abandoned"})
}
