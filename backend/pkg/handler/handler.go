package handler

import (
	"time"

	"github.com/dimqueue/darts/pkg/config"
	"github.com/dimqueue/darts/pkg/service"
	"github.com/dimqueue/darts/pkg/validation"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services  *service.Service
	validator *validation.Validator
}

func NewHandler(services *service.Service, validator *validation.Validator) *Handler {
	return &Handler{
		services:  services,
		validator: validator,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(RequestIDMiddleware())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", RequestIDHeader},
		ExposeHeaders:    []string{"Content-Length", RequestIDHeader},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	api := router.Group("/api", h.userIdentity)
	{
		games := api.Group("/games")
		{
			games.POST("", h.createGame)
			games.GET("", h.getAllGames)
			games.GET("/:id", h.getGameById)

			guess := games.Group(":id/guesses")
			{
				guess.POST("", h.createGuess)
				guess.GET("", h.getAllGuessByGame)
				guess.GET("/:guess_id", h.getGuessById)

			}
		}

		profile := api.Group("/profile")
		{
			profile.GET("", h.getMyProfile)
			profile.PUT("", h.updateMyProfile)
			profile.GET("/settings", h.getMySettings)
			profile.PUT("/settings", h.updateMySettings)
			profile.GET("/statistics", h.getMyStatistics)
			profile.GET("/statistics/languages", h.getMyLanguageStats)
		}

		leaderboard := api.Group("/leaderboard")
		{
			leaderboard.GET("", h.getLeaderboard)
			leaderboard.GET("/my-rank", h.getMyRank)
		}
	}

	public := router.Group("/public")
	{
		public.GET("/profile/:username", h.getProfileByUsername)
	}

	return router
}
