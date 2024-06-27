package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	FirstName    string
	LastName     string
	PhoneNumber  string `gorm:"unique"`
	Address      string
	PIN          string
	Balance      float64
	RefreshToken string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
