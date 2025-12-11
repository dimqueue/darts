package handler

import (
	"net/http"
	"strconv"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/gin-gonic/gin"
)

// @Summary Get leaderboard
// @Tags Leaderboard
// @Param type query string true "Leaderboard type (global, weekly, monthly)"
// @Param language query string false "Language filter (en, ua, etc)"
// @Param limit query int false "Number of entries" default(50)
// @Param offset query int false "Pagination offset" default(0)
// @Success 200 {object} model.LeaderboardResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/leaderboard [get]
func (h *Handler) getLeaderboard(c *gin.Context) {
	leaderboardType := c.Query("type")
	if leaderboardType == "" {
		handleError(c, ErrValidation("type parameter is required (global, weekly, monthly)"))
		return
	}

	if leaderboardType != "global" && leaderboardType != "daily" && leaderboardType != "weekly" && leaderboardType != "monthly" {
		handleError(c, ErrValidation("invalid type: must be global, daily, weekly, or monthly"))
		return
	}

	query := model.LeaderboardQuery{
		Type:   leaderboardType,
		Limit:  parseIntQueryWithBounds(c, "limit", 50, 1, 100),
		Offset: parseIntQueryWithBounds(c, "offset", 0, 0, 10000),
	}

	language := c.Query("language")
	if language != "" {
		query.Language = &language
	}

	response, err := h.services.Leaderboard.GetLeaderboard(c.Request.Context(), query)
	if err != nil {
		handleError(c, err)
		return
	}

	userId, err := getUserId(c)
	if err == nil {
		rank, _ := h.services.Leaderboard.GetUserRank(c.Request.Context(), userId, query)
		response.CurrentUserRank = rank
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get current user's rank
// @Tags Leaderboard
// @Security ApiKeyAuth
// @Success 200 {object} model.UserRanks
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/leaderboard/my-rank [get]
func (h *Handler) getMyRank(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		handleError(c, ErrUnauthorized("user not found"))
		return
	}

	ranks, err := h.services.Leaderboard.GetAllUserRanks(c.Request.Context(), userId)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, ranks)
}

func parseIntQuery(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func parseIntQueryWithBounds(c *gin.Context, key string, defaultValue, min, max int) int {
	value := parseIntQuery(c, key, defaultValue)
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
