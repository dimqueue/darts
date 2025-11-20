package handler

import (
	"time"

	"github.com/dimqueue/darts/pkg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
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
			games.PUT("/:id", h.updateGame)
			games.DELETE("/:id", h.deleteGame)

			guess := games.Group(":id/guesses")
			{
				guess.POST("", h.createGuess)
				guess.GET("", h.getAllGuessByGame)
				guess.GET("/:guess_id", h.getGuessById)

			}
		}
	}
	return router
}
