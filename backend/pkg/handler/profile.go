package handler

import (
	"net/http"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/gin-gonic/gin"
)

// @Summary Get current user's profile
// @Tags Profile
// @Security ApiKeyAuth
// @Success 200 {object} model.UserProfileSummary
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/profile [get]
func (h *Handler) getMyProfile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	profile, err := h.services.Profile.GetProfileSummary(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, profile)
}

// @Summary Get user profile by username
// @Tags Profile
// @Param username path string true "Username"
// @Success 200 {object} model.UserProfileSummary
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/profile/{username} [get]
func (h *Handler) getProfileByUsername(c *gin.Context) {
	username := c.Param("username")

	profile, err := h.services.Profile.GetProfileByUsername(username)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	if !profile.ShowProfilePublic {
		newErrorResponse(c, http.StatusForbidden, "profile is private")
		return
	}

	c.JSON(http.StatusOK, profile)
}

// @Summary Update current user's profile
// @Tags Profile
// @Security ApiKeyAuth
// @Param input body model.UpdateProfileInput true "Profile update data"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/profile [put]
func (h *Handler) updateMyProfile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var input model.UpdateProfileInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := h.services.Profile.UpdateProfile(userId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

// @Summary Get current user's settings
// @Tags Profile
// @Security ApiKeyAuth
// @Success 200 {object} model.UserSettings
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/profile/settings [get]
func (h *Handler) getMySettings(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	settings, err := h.services.Profile.GetSettings(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, settings)
}

// @Summary Update current user's settings
// @Tags Profile
// @Security ApiKeyAuth
// @Param input body model.UpdateSettingsInput true "Settings update data"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/profile/settings [put]
func (h *Handler) updateMySettings(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var input model.UpdateSettingsInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := h.services.Profile.UpdateSettings(userId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

// @Summary Get current user's statistics
// @Tags Profile
// @Security ApiKeyAuth
// @Success 200 {object} model.UserStatistics
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/profile/statistics [get]
func (h *Handler) getMyStatistics(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	stats, err := h.services.Profile.GetStatistics(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary Get current user's language statistics
// @Tags Profile
// @Security ApiKeyAuth
// @Param language query string false "Language code (if empty, returns all)"
// @Success 200 {array} model.UserLanguageStats
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/profile/statistics/languages [get]
func (h *Handler) getMyLanguageStats(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	language := c.Query("language")

	if language != "" {
		stats, err := h.services.Profile.GetLanguageStats(userId, language)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, stats)
		return
	}

	stats, err := h.services.Profile.GetAllLanguageStats(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, stats)
}
