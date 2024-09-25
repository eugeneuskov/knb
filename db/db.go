package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"knb/app/config"
	"knb/app/entities"
)

type DB struct {
	db *gorm.DB
}

func (db *DB) NewPostgresDb(config *config.DbConfig) error {
	dbInstance, err := gorm.Open(
		postgres.New(postgres.Config{
			DriverName: "pgx",
			DSN: fmt.Sprintf(
				"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
				config.Host,
				config.Port,
				config.User,
				config.DbName,
				config.Password,
				config.SslMode,
			),
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{},
	)

	if err != nil {
		return err
	}

	db.db = dbInstance

	return nil
}

func (db *DB) Migrate() error {
	return db.db.AutoMigrate(
		&entities.Player{},
		&entities.Game{},
		&entities.GamePrize{},
		&entities.GameResult{},
	)
}

func (db *DB) DropMigrate() error {
	return db.db.Migrator().DropTable(
		&entities.GameResult{},
		&entities.GamePlayer{},
		&entities.GamePrize{},
		&entities.Game{},
		&entities.Player{},
	)
}

func (db *DB) Insert(item interface{}) error {
	return db.db.Create(item).Error
}

func (db *DB) DB() *gorm.DB {
	return db.db
}
