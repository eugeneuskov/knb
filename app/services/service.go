package services

import (
	"knb/app/config"
	"knb/app/interfaces"
	"knb/app/repositories"
)

type Service struct {
	Security interfaces.ServiceSecurity
	Auth     interfaces.ServiceAuth
}

func NewService(repository *repositories.Repository, config *config.Config) *Service {
	return &Service{
		Security: newSecurityService(),
		Auth:     newAuthService(repository.Player, config.AuthConfig.TokenSigningKey),
	}
}
