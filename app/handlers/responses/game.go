package responses

import "github.com/google/uuid"

type GameNewGameResponse struct {
	Id uuid.UUID `json:"id"`
}
