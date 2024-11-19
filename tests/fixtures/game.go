package fixtures

import (
	"github.com/google/uuid"
	"knb/app/entities"
)

const (
	AlreadyExistGameUuid = "769b8215-0345-4c6c-aa08-ca18f4318fa0"
)

func createGamesTestFixtures() []entities.Game {
	game1 := entities.NewGame([]entities.Player{
		{
			ID: uuid.MustParse(AlreadyExistPlayer1Uuid),
		},
	})
	game1.ID = uuid.MustParse(AlreadyExistGameUuid)

	return []entities.Game{
		*game1,
	}
}
