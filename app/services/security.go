package services

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

const (
	salt     = "8yU9nA`tx$=k-XPB/4@(ctf[n$Me;#Mhk^T]jBZC/K"
	tokenTTL = 12 * time.Hour
)

type securityService struct {
	tokenSigningKey string
}

func newSecurityService(tokenSigningKey string) *securityService {
	return &securityService{tokenSigningKey}
}

type tokenClaims struct {
	jwt.StandardClaims
	PlayerId string `json:"player_id"`
}

func (s *securityService) GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *securityService) GenerateAuthToken(playerId uuid.UUID) (string, error) {
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

	return token.SignedString([]byte(s.tokenSigningKey))
}

func (s *securityService) ParseAuthToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signed method")
			}
			return []byte(s.tokenSigningKey), nil
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
