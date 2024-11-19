package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"knb/app/handlers/responses"
	"knb/app/services"
)

const (
	authorizationToken   = "Access-Token"
	authorizationContext = "authorizationCtx"
)

type Handler struct {
	service  *services.Service
	response *responses.Response
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		service:  service,
		response: responses.NewResponse(),
	}
}

func (h *Handler) InitRoutes(mode string) *gin.Engine {
	gin.SetMode(mode)
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/registration", h.authRegistration)
		auth.POST("/login", h.authLogin)
	}

	game := router.Group("/game", h.userAccessIdentity)
	{
		game.POST("/new", h.gameNewGame)
		game.POST("/join/:id", h.gameJoinGame)
		game.POST("/start/:id", h.gameStart)
	}

	return router
}

func (h *Handler) checkGetParam(c *gin.Context, paramName string) (string, error) {
	paramValue := c.Param(paramName)
	if paramValue == "" {
		return "", fmt.Errorf("param %s is required", paramName)
	}

	return paramValue, nil
}
