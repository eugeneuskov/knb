package interfaces

import "github.com/google/uuid"

type ServiceSecurity interface {
	GeneratePasswordHash(password string) string
	GenerateAuthToken(playerId uuid.UUID) (string, error)
	ParseAuthToken(accessToken string) (string, error)
}
