package fixtures

import (
	"github.com/google/uuid"
	"knb/app/dictionary"
	"knb/app/entities"
	"time"
)

const (
	GameFinishedUuid = "715053c0-517c-4f82-9809-28dedea83ec5"
	GameStartedUuid  = "1eb63a04-b947-4d64-8ec1-51006762fc25"
	GameWaitingUuid  = "769b8215-0345-4c6c-aa08-ca18f4318fa0"
	GamePlannedUuid  = "690bbd55-3d63-4a49-8be3-04f31e52d759"
)

func createGamesTestFixtures() []entities.Game {
	gameFinished := entities.NewGame([]entities.Player{
		{
			ID: uuid.MustParse(Player1Uuid),
		},
		{
			ID: uuid.MustParse(Player2Uuid),
		},
	})
	gameFinished.ID = uuid.MustParse(GameFinishedUuid)
	gameFinished.StartedAt = time.Now().Add(-7 * 24 * time.Hour)
	gameFinished.FinishedAt = time.Now().Add(-7 * 23 * time.Hour)
	gameFinished.Status = dictionary.GameStatusFinished

	gameStarted := entities.NewGame([]entities.Player{
		{
			ID: uuid.MustParse(Player1Uuid),
		},
		{
			ID: uuid.MustParse(Player2Uuid),
		},
	})
	gameStarted.ID = uuid.MustParse(GameStartedUuid)
	gameStarted.StartedAt = time.Now().Add(-1 * time.Minute)
	gameStarted.Status = dictionary.GameStatusStarted

	gameWaiting := entities.NewGame([]entities.Player{
		{
			ID: uuid.MustParse(Player1Uuid),
		},
	})
	gameWaiting.ID = uuid.MustParse(GameWaitingUuid)
	gameWaiting.StartedAt = time.Now().Add(7 * 24 * time.Hour)
	gameWaiting.Status = dictionary.GameStatusWaiting

	gamePlanned := entities.NewGame([]entities.Player{
		{
			ID: uuid.MustParse(Player1Uuid),
		},
	})
	gameWaiting.ID = uuid.MustParse(GamePlannedUuid)
	gameWaiting.StartedAt = time.Now().Add(7 * 24 * time.Hour)
	gameWaiting.Status = dictionary.GameStatusPlanned

	return []entities.Game{
		*gameFinished,
		*gameStarted,
		*gameWaiting,
		*gamePlanned,
	}
}
