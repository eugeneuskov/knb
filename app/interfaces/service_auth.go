package interfaces

import (
	"github.com/google/uuid"
)

type ServiceAuth interface {
	ParseToken(accessToken string) (string, error)
	Registration(login, password string) (uuid.UUID, error)
	Login(login, password string) (string, error)
}
