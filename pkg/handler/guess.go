package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createGuessInput struct {
	GameId int    `json:"game_id" binding:"required"`
	Guess  string `json:"guess" binding:"required"`
}

func (h *Handler) createGuess(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user not found")
	}
	var input createGuessInput

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = validateWord(input.Guess); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	distance, err := h.services.Guess.CreateGuess(userId, input.GameId, input.Guess)

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"distance": distance,
	})

}

func validateWord(word string) error {
	return nil
}

func (h *Handler) getAllGuessByGame(c *gin.Context) {

}

func (h *Handler) getGuessById(c *gin.Context) {

}
