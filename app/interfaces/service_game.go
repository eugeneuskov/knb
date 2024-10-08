package interfaces

import (
	"github.com/google/uuid"
	"knb/app/entities"
)

type ServiceGame interface {
	NewGameRequest(playerOwnerId uuid.UUID) (*entities.Game, error)
	FindGame(gameId uuid.UUID) (*entities.Game, error)
	JoinGame(playerId uuid.UUID, gameId uuid.UUID) error
	StartGame(gameId uuid.UUID) error
}
