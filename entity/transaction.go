package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID        uuid.UUID
	Type          string
	Amount        float64
	Remarks       string
	BalanceBefore float64
	BalanceAfter  float64
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_date"`
	DeletedAt     gorm.DeletedAt
}
