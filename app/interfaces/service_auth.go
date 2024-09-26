package interfaces

import (
	"github.com/google/uuid"
	"knb/app/entities"
)

type ServiceAuth interface {
	Registration(login, password string) (uuid.UUID, error)
	Login(login, password string) (*entities.Player, error)
}
