package responses

import (
	"github.com/google/uuid"
	"time"
)

type GameNewGameResponse struct {
	ID uuid.UUID `json:"id"`
}

type GamePlayerResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type GameJoinGameResponse struct {
	ID        uuid.UUID            `json:"id"`
	StartedAt time.Time            `json:"started_at"`
	Players   []GamePlayerResponse `json:"players"`
}
