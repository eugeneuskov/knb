package interfaces

import "github.com/google/uuid"

type RepositoryAuth interface {
	Insert(login, password string) (*uuid.UUID, error)
}
