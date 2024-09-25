package responses

import "github.com/google/uuid"

type AuthRegistrationResponse struct {
	ID uuid.UUID `json:"id"`
}

type AuthLoginResponse struct {
	Token string `json:"token"`
}
