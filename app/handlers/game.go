package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			ID: game.ID,
		},
	)
}

func (h *Handler) gameJoinGame(c *gin.Context) {
	playerId, err := h.getAccessContext(c)
	if err != nil {
		h.response.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	gameIdParam, err := h.checkGetParam(c, "id")
	if err != nil {
		h.response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	gameId, err := uuid.Parse(gameIdParam)
	if err != nil {
		h.response.NewErrorResponse(c, http.StatusBadRequest, "game id is invalid")
		return
	}

	game, err := h.service.Game.JoinGame(playerId, gameId)
	if err != nil {
		h.response.ParseError(c, err)
		return
	}

	gamePlayers := make([]responses.GamePlayerResponse, 0, len(game.Players))
	for _, player := range game.Players {
		gamePlayers = append(gamePlayers, responses.GamePlayerResponse{
			ID:   player.ID,
			Name: player.DisplayName,
		})
	}

	h.response.NewOkResponse(
		c, http.StatusOK,
		responses.GameJoinGameResponse{
			ID:        game.ID,
			StartedAt: game.StartedAt,
			Players:   gamePlayers,
		},
	)
}

func (h *Handler) gameStart(c *gin.Context) {
	h.response.NewErrorResponse(c, http.StatusMethodNotAllowed, "/game/start")
}
