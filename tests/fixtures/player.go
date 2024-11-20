package fixtures

import (
	"github.com/google/uuid"
	"knb/app/entities"
)

const (
	Player1Uuid        = "6c425956-3a9c-47c9-910d-f43a0d7234e8"
	Player1Email       = "exist-email@test.com"
	Player1Password    = "strong-pass"
	Player1DisplayName = "John Doe"

	Player2Uuid        = "bcfccdd5-c9dc-4acd-973f-b487fb46b9fb"
	Player2DisplayName = "Jane Doe"

	Player3Uuid        = "8e5a5f40-7015-4892-95c0-ffe407767f6c"
	Player3DisplayName = "Arnold Doe"
)

func createPlayersTestFixtures() []entities.Player {
	player1 := entities.NewPlayer(Player1Email, Player1Password, Player1DisplayName)
	player1.ID = uuid.MustParse(Player1Uuid)

	player2 := entities.NewPlayer("", "", Player2DisplayName)
	player2.ID = uuid.MustParse(Player2Uuid)

	player3 := entities.NewPlayer("", "", Player3DisplayName)
	player3.ID = uuid.MustParse(Player3Uuid)

	return []entities.Player{
		*player1,
		*player2,
		*player3,
	}
}
