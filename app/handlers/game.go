package handlers

import (
	"github.com/gin-gonic/gin"
	"knb/app/handlers/responses"
	"net/http"
)

func (h *Handler) gameNewGame(c *gin.Context) {
	playerId, err := h.getAccessContext(c)
	if err != nil {
		h.response.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	game, err := h.service.Game.NewGameRequest(playerId)
	if err != nil {
		h.response.ParseError(c, err)
		return
	}

	h.response.NewOkResponse(
		c, http.StatusCreated,
		responses.GameNewGameResponse{
			Id: game.ID,
		},
	)
}

func (h *Handler) gameJoinGame(c *gin.Context) {
	h.response.NewErrorResponse(c, http.StatusMethodNotAllowed, "/game/join")
}

func (h *Handler) gameStart(c *gin.Context) {
	h.response.NewErrorResponse(c, http.StatusMethodNotAllowed, "/game/start")
}
