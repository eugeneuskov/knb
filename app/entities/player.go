package entities

import (
	"github.com/google/uuid"
	"time"
)

type Player struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey'"`
	Email       string    `gorm:"size:255;unique;not null"`
	Password    string    `gorm:"size:255;not null"`
	DisplayName string    `gorm:"size:255;null"`
	Points      uint      `gorm:"not null;default:0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func NewPlayer(email, password, displayName string) *Player {
	return &Player{
		ID:          uuid.New(),
		Email:       email,
		Password:    password,
		DisplayName: displayName,
	}
}
