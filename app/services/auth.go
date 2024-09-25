package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	customErrors "knb/app/errors"
	"knb/app/interfaces"
	"knb/app/repositories"
	"time"
)

type authService struct {
	playerRepository interfaces.RepositoryPlayer
	tokenSigningKey  string
}

func newAuthService(
	authRepository interfaces.RepositoryPlayer,
	tokenSigningKey string,
) *authService {
	return &authService{authRepository, tokenSigningKey}
}

type tokenClaims struct {
	jwt.StandardClaims
	PlayerId string `json:"player_id"`
}

func (a *authService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signed method")
			}
			return []byte(a.tokenSigningKey), nil
		},
	)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not type *tokenClaims")
	}

	return claims.PlayerId, err
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

func (a *authService) Login(login, password string) (string, error) {
	player, err := a.playerRepository.FindByLoginAndPassword(login, password)
	if err != nil {
		var notFoundErr *customErrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			return "", customErrors.NewWrongLoginError("Login or Password are incorrect")
		}

		return "", err
	}

	accessToken, err := a.generateAuthToken(player.ID)
	if err != nil {
		return "", errors.New("something went wrong")
	}

	return accessToken, nil
}

func (a *authService) generateAuthToken(playerId uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&tokenClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenTTL).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			PlayerId: playerId.String(),
		},
	)

	return token.SignedString([]byte(a.tokenSigningKey))
}
