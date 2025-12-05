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
		newErrorResponse(c, http.StatusInternalServerError, "user not found")
		return
	}
	gameId, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input createGuessInput

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.validator.ValidateWord(input.Guess); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	distance, err := h.services.Game.MakeGuess(userId, gameId, input.Guess)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"distance": distance,
	})
}

type getAllGuessByGameResponse struct {
	Data []model.Guess `json:"data"`
}

func (h *Handler) getAllGuessByGame(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	gameId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	guesses, err := h.services.Game.GetAllGuessByGame(userId, gameId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllGuessByGameResponse{
		Data: guesses,
	})
}

func (h *Handler) getGuessById(c *gin.Context) {

}
