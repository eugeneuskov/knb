package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func (h *Handler) userAccessIdentity(c *gin.Context) {
	headerToken, err := h.checkHeader(c, authorizationToken)
	if err != nil {
		h.response.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	playerId, err := h.service.Security.ParseAuthToken(headerToken)
	if err != nil {
		h.response.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if playerId == "" {
		h.response.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	c.Set(authorizationContext, playerId)
}

func (h *Handler) checkHeader(c *gin.Context, headerName string) (string, error) {
	header := c.GetHeader(headerName)
	if header == "" {
		return "", errors.New(fmt.Sprintf("empty '%s' header", headerName))
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) == 0 {
		return "", errors.New(fmt.Sprintf("invalid '%s' header", headerName))
	}

	headerValue := headerParts[0]
	if headerValue == "" {
		return "", errors.New(fmt.Sprintf("invalid '%s' header", headerName))
	}

	return headerValue, nil
}

func (h *Handler) getAccessContext(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(authorizationContext)
	if !ok {
		return uuid.Nil, errors.New("request is invalid")
	}

	playerIdValue, ok := id.(string)
	if !ok {
		return uuid.Nil, errors.New("player is invalid")
	}

	playerId, err := uuid.Parse(playerIdValue)
	if err != nil {
		return uuid.Nil, errors.New("player is invalid")
	}

	return playerId, nil
}
