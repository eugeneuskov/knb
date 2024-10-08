package services

import (
	"knb/app/config"
	"knb/app/interfaces"
	"knb/app/repositories"
)

type Service struct {
	Security interfaces.ServiceSecurity
	Auth     interfaces.ServiceAuth
	Game     interfaces.ServiceGame
}

func NewService(repository *repositories.Repository, config *config.Config) *Service {
	return &Service{
		Security: newSecurityService(config.AuthConfig.TokenSigningKey),
		Auth:     newAuthService(repository.Player),
		Game:     newGameService(repository.Game, repository.Player),
	}
}
