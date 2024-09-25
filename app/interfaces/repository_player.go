package interfaces

import (
	"github.com/google/uuid"
	"knb/app/entities"
)

type RepositoryPlayer interface {
	Create(login, password string) (*entities.Player, error)
	FindById(id uuid.UUID) (*entities.Player, error)
	FindByLoginAndPassword(login, password string) (*entities.Player, error)
}
