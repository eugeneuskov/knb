package services

import (
	"crypto/sha1"
	"fmt"
	"time"
)

const (
	salt     = "8yU9nA`tx$=k-XPB/4@(ctf[n$Me;#Mhk^T]jBZC/K"
	tokenTTL = 12 * time.Hour
)

type securityService struct {
}

func newSecurityService() *securityService {
	return &securityService{}
}

func (s *securityService) GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *securityService) GenerateAuthToken(login, password string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *securityService) ParseAuthToken(token string) (string, error) {
	//TODO implement me
	panic("implement me")
}
