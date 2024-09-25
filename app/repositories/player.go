package repositories

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"knb/app/entities"
	customErrors "knb/app/errors"
	"log"
)

type playerRepository struct {
	db *gorm.DB
}

func newPlayerRepository(db *gorm.DB) *playerRepository {
	return &playerRepository{db}
}

func (p *playerRepository) Create(login, password string) (*entities.Player, error) {
	player := entities.NewPlayer(login, password, "")

	if err := p.db.Create(player).Error; err != nil {
		return nil, err
	}

	return player, nil
}

func (p *playerRepository) FindById(id uuid.UUID) (*entities.Player, error) {
	var players []entities.Player
	result := p.db.Limit(1).Find(&players, "id = ?", id.String())
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	if len(players) == 0 {
		return nil, customErrors.NewNotFoundError(fmt.Sprintf("player with id %s not found", id))
	}

	return &players[0], nil
}

func (p *playerRepository) FindByLoginAndPassword(login, password string) (*entities.Player, error) {
	var players []entities.Player
	result := p.db.Limit(1).Find(
		&players,
		"email = ? AND password = ?",
		login, password,
	)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	if len(players) == 0 {
		return nil, customErrors.NewNotFoundError("player was not found")
	}

	return &players[0], nil
}
