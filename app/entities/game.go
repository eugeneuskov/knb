package entities

import (
	"github.com/google/uuid"
	"time"
)

type Game struct {
	ID         uuid.UUID    `gorm:"type:uuid;primaryKey"`
	StartedAt  time.Time    `gorm:"type:timestamp"`
	FinishedAt time.Time    `gorm:"type:timestamp"`
	Players    []Player     `gorm:"many2many:game_players"`
	Prizes     []GamePrize  `gorm:"foreignKey:GameID"`
	Result     []GameResult `gorm:"foreignKey:GameID"`
}

func NewGame(players []Player) *Game {
	return &Game{
		ID:      uuid.New(),
		Players: players,
	}
}

type GamePlayer struct {
	PlayerID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_game_player"`
	GameID   uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_game_player"`
}

type GamePrize struct {
	GameID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_game_prize"`
	Place  uint8     `gorm:"type:int;uniqueIndex:idx_game_prize"`
	Prize  uint      `gorm:"type:int"`
}

type GameResult struct {
	GameID   uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_game_result"`
	PlayerID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_game_result"`
	Place    uint8     `gorm:"type:int;uniqueIndex:idx_game_result"`
}
