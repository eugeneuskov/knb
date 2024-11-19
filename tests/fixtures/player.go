package fixtures

import (
	"github.com/google/uuid"
	"knb/app/entities"
)

const (
	AlreadyExistPlayer1Uuid        = "6c425956-3a9c-47c9-910d-f43a0d7234e8"
	AlreadyExistPlayer1Email       = "exist-email@test.com"
	AlreadyExistPlayer1Password    = "strong-pass"
	AlreadyExistPlayer1DisplayName = "John Doe"

	AlreadyExistPlayer2Uuid        = "bcfccdd5-c9dc-4acd-973f-b487fb46b9fb"
	AlreadyExistPlayer2DisplayName = "Jane Doe"
)

func createPlayersTestFixtures() []entities.Player {
	player1 := entities.NewPlayer(AlreadyExistPlayer1Email, AlreadyExistPlayer1Password, AlreadyExistPlayer1DisplayName)
	player1.ID = uuid.MustParse(AlreadyExistPlayer1Uuid)

	player2 := entities.NewPlayer("", "", AlreadyExistPlayer2DisplayName)
	player2.ID = uuid.MustParse(AlreadyExistPlayer2Uuid)

	return []entities.Player{
		*player1,
		*player2,
	}
}
