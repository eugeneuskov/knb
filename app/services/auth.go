package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"knb/app/entities"
	customErrors "knb/app/errors"
	"knb/app/interfaces"
	"knb/app/repositories"
)

type authService struct {
	playerRepository interfaces.RepositoryPlayer
}

func newAuthService(authRepository interfaces.RepositoryPlayer) *authService {
	return &authService{authRepository}
}

func (a *authService) Registration(login, password string) (uuid.UUID, error) {
	player, err := a.playerRepository.Create(login, password)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == repositories.UniqueViolation {
				return uuid.Nil, customErrors.NewUniqueViolationError("This login already exists")
			}
		}

		return uuid.Nil, err
	}

	return player.ID, nil
}

func (a *authService) Login(login, password string) (*entities.Player, error) {
	player, err := a.playerRepository.FindByLoginAndPassword(login, password)
	if err != nil {
		var notFoundErr *customErrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			return nil, customErrors.NewWrongLoginError("Login or Password are incorrect")
		}

		return nil, err
	}

	return player, nil
}
