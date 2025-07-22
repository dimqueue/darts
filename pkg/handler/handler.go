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
		users := api.Group("/users")
		{
			users.POST("/", h.createUser)
			users.GET("/", h.getAllUser)
			users.GET("/:id", h.getUserById)
			users.PUT("/:id", h.updateUserById)
			users.DELETE("/:id", h.deleteUserById)

			games := users.Group(":id/games")
			{
				games.POST("/", h.createGame)
				games.GET("/", h.getAllGame)
				games.GET("/:game_id", h.getGameById)

			}
		}
	}
	return router
}
