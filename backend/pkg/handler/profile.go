package handler

import (
	"net/http"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/service"
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
		handleError(c, ErrUnauthorized("user not found"))
		return
	}

	profile, err := h.services.Profile.GetProfileSummary(c.Request.Context(), userId)
	if err != nil {
		handleError(c, err)
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

	profile, err := h.services.Profile.GetProfileByUsername(c.Request.Context(), username)
	if err != nil {
		handleError(c, service.ErrProfileNotFound)
		return
	}

	if !profile.ShowProfilePublic {
		handleError(c, service.ErrProfilePrivate)
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
		handleError(c, ErrUnauthorized("user not found"))
		return
	}

	var input model.UpdateProfileInput
	if err := c.BindJSON(&input); err != nil {
		handleError(c, ErrBadRequest("invalid request body"))
		return
	}

	if err := h.services.Profile.UpdateProfile(c.Request.Context(), userId, input); err != nil {
		handleError(c, err)
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
		handleError(c, ErrUnauthorized("user not found"))
		return
	}

	settings, err := h.services.Profile.GetSettings(c.Request.Context(), userId)
	if err != nil {
		handleError(c, err)
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
		handleError(c, ErrUnauthorized("user not found"))
		return
	}

	var input model.UpdateSettingsInput
	if err := c.BindJSON(&input); err != nil {
		handleError(c, ErrBadRequest("invalid request body"))
		return
	}

	if err := h.services.Profile.UpdateSettings(c.Request.Context(), userId, input); err != nil {
		handleError(c, err)
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
		handleError(c, ErrUnauthorized("user not found"))
		return
	}

	stats, err := h.services.Profile.GetStatistics(c.Request.Context(), userId)
	if err != nil {
		handleError(c, err)
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
		handleError(c, ErrUnauthorized("user not found"))
		return
	}

	language := c.Query("language")

	if language != "" {
		if err := h.validator.ValidateLanguage(language); err != nil {
			handleError(c, ErrValidation(err.Error()))
			return
		}
		stats, err := h.services.Profile.GetLanguageStats(c.Request.Context(), userId, language)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, stats)
		return
	}

	stats, err := h.services.Profile.GetAllLanguageStats(c.Request.Context(), userId)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}
