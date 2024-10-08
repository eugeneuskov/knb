package interfaces

import (
	"github.com/google/uuid"
	"knb/app/entities"
)

type GameRepository interface {
	CreateGame(players ...uuid.UUID) (*entities.Game, error)
	FindById(gameId uuid.UUID) (*entities.Game, error)
	AddPlayers(game *entities.Game, players []uuid.UUID) error
	StartGame(game *entities.Game) error
}
