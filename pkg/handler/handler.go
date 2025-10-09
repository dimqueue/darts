package handler

import (
	"github.com/dimqueue/darts/pkg/service"
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
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	api := router.Group("/api", h.userIdentity)
	{
		games := api.Group("/games")
		{
			games.POST("/", h.createGame)
			games.GET("/", h.getAllGames)
			games.GET("/:id", h.getGameById)
			games.PUT("/:id", h.updateGame)
			games.DELETE("/:id", h.deleteGame)

			guess := games.Group(":id/guess")
			{
				guess.POST("/", h.createGuess)
				guess.GET("/", h.getAllGuessByGame)
				guess.GET("/:guess_id", h.getGuessById)

			}
		}
	}
	return router
}
