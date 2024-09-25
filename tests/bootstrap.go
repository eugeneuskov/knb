package tests

import (
	"fmt"
	"gorm.io/gorm"
	"knb/app/config"
	"knb/db"
	"log"
)

type BootstrapTest struct {
	db     *db.DB
	config *config.Config
}

func NewBootstrapTest(envFilePath string) *BootstrapTest {
	conf, err := new(config.Config).Init(envFilePath)
	if err != nil {
		log.Fatalf("Failed initializing test config: %s\n", err.Error())
	}

	return &BootstrapTest{
		db:     new(db.DB),
		config: conf,
	}
}

func (bt *BootstrapTest) SetupTestDB() error {
	if err := bt.db.NewPostgresDb(&bt.config.DbConfig); err != nil {
		return fmt.Errorf("Failed connect to test DB: %s\n", err.Error())
	}
	if err := bt.db.Migrate(); err != nil {
		return fmt.Errorf("Failed migration test DB: %s\n", err.Error())
	}

	return nil
}

func (bt *BootstrapTest) TeardownTestDB() error {
	if err := bt.db.DropMigrate(); err != nil {
		return fmt.Errorf("Failed drop migration test DB: %s\n", err.Error())
	}

	return nil
}

func (bt *BootstrapTest) DB() *gorm.DB {
	return bt.db.DB()
}

func (bt *BootstrapTest) Config() *config.Config {
	return bt.config
}
