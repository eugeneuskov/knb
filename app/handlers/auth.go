package handlers

import (
	"github.com/gin-gonic/gin"
	"knb/app/handlers/requests"
	"knb/app/handlers/responses"
	"net/http"
)

func (h *Handler) authRegistration(c *gin.Context) {
	if c.Request.Body == http.NoBody {
		h.response.NewErrorResponse(c, http.StatusBadRequest, "Request is empty.")
		return
	}

	var request requests.AuthRegistrationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.Auth.Registration(
		request.Login,
		h.service.Security.GeneratePasswordHash(request.Password),
	)
	if err != nil {
		h.response.ParseError(c, err)
		return
	}

	h.response.NewOkResponse(
		c, http.StatusCreated,
		responses.AuthRegistrationResponse{ID: response},
	)
}

func (h *Handler) authLogin(c *gin.Context) {
	if c.Request.Body == http.NoBody {
		h.response.NewErrorResponse(c, http.StatusBadRequest, "Request is empty.")
		return
	}

	var request requests.AuthLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.service.Auth.Login(
		request.Login,
		h.service.Security.GeneratePasswordHash(request.Password),
	)
	if err != nil {
		h.response.ParseError(c, err)
		return
	}

	h.response.NewOkResponse(
		c, http.StatusOK,
		responses.AuthLoginResponse{
			Token: response,
		},
	)
}
