package handlers

import (
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

	return router
}
