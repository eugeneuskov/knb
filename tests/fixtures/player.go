package fixtures

import (
	"github.com/google/uuid"
	"knb/app/entities"
)

const (
	AlreadyExistUuid     = "6c425956-3a9c-47c9-910d-f43a0d7234e8"
	AlreadyExistEmail    = "exist-email@test.com"
	AlreadyExistPassword = "strong-pass"
)

func createPlayersTestFixtures() []entities.Player {
	player1 := entities.NewPlayer(AlreadyExistEmail, AlreadyExistPassword, "")
	player1.ID = uuid.MustParse(AlreadyExistUuid)

	return []entities.Player{
		*player1,
	}
}
