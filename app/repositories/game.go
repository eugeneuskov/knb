package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"knb/app/entities"
)

type gameRepository struct {
	db *gorm.DB
}

func newGameRepository(db *gorm.DB) *gameRepository {
	return &gameRepository{db}
}

func (g *gameRepository) CreateGame(players ...uuid.UUID) (*entities.Game, error) {
	var game *entities.Game

	if err := g.db.Transaction(func(tx *gorm.DB) error {
		gamePlayers := make([]entities.Player, 0, len(players))
		for _, playerId := range players {
			gamePlayers = append(gamePlayers, entities.Player{ID: playerId})
		}

		game = entities.NewGame(gamePlayers)

		return tx.Create(&game).Error
	}); err != nil {
		return nil, err
	}

	return game, nil
}

func (g *gameRepository) FindById(gameId uuid.UUID) (*entities.Game, error) {
	var game *entities.Game

	err := g.db.
		Preload("Players").
		First(&game, "id = ?", gameId).
		Error

	return game, err
}

func (g *gameRepository) AddPlayers(game *entities.Game, playerIds []uuid.UUID) error {
	players := make([]entities.Player, 0, len(playerIds))
	for _, playerId := range playerIds {
		players = append(players, entities.Player{ID: playerId})
	}

	return g.db.Model(&game).Association("Players").Append(&players)
}

func (g *gameRepository) StartGame(game *entities.Game) error {
	//TODO implement me
	panic("implement me")
}
