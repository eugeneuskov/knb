package repositories

import (
	"gorm.io/gorm"
	"knb/app/interfaces"
)

const (
	UniqueViolation     = "23505" // Код ошибки для уникального ограничения
	foreignKeyViolation = "23503" // Код ошибки для внешнего ключа
	notNullViolation    = "23502" // Код ошибки для NOT NULL ограничения
	checkViolation      = "23514" // Код ошибки для CHECK ограничения
	exclusionViolation  = "23504" // Код ошибки для EXCLUSION ограничения
)

type Repository struct {
	Player interfaces.RepositoryPlayer
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Player: newPlayerRepository(db),
	}
}
