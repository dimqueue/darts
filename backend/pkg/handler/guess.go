package handler

import (
	"net/http"
	"strconv"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/gin-gonic/gin"
)

type createGuessInput struct {
	Guess string `json:"guess" binding:"required"`
}

func (h *Handler) createGuess(c *gin.Context) {
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

	var input createGuessInput
	if err = c.BindJSON(&input); err != nil {
		handleError(c, ErrBadRequest("invalid request body"))
		return
	}

	if err = h.validator.ValidateWord(input.Guess); err != nil {
		handleError(c, ErrValidation(err.Error()))
		return
	}

	result, err := h.services.Game.MakeGuess(c.Request.Context(), userId, gameId, input.Guess)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"rank":          result.Rank,
		"found":         result.Found,
		"in_vocabulary": result.InVocabulary,
	})
}

type getAllGuessByGameResponse struct {
	Data []model.Guess `json:"data"`
}

func (h *Handler) getAllGuessByGame(c *gin.Context) {
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

	guesses, err := h.services.Game.GetAllGuessByGame(c.Request.Context(), userId, gameId)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, getAllGuessByGameResponse{
		Data: guesses,
	})
}

func (h *Handler) getGuessById(c *gin.Context) {

}
