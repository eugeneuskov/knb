package fixtures

import (
	"gorm.io/gorm"
	"knb/app/services"
)

type Fixture struct {
	db      *gorm.DB
	service *services.Service
}

func NewFixtures(db *gorm.DB, service *services.Service) *Fixture {
	return &Fixture{db, service}
}

func (f *Fixture) LoadPlayersFixture() error {
	for _, player := range createPlayersTestFixtures() {
		player.Password = f.service.Security.GeneratePasswordHash(player.Password)
		if err := f.db.Create(&player).Error; err != nil {
			return err
		}
	}

	return nil
}

func (f *Fixture) LoadGamesFixture() error {
	for _, game := range createGamesTestFixtures() {
		if err := f.db.Create(&game).Error; err != nil {
			return err
		}
	}

	return nil
}
