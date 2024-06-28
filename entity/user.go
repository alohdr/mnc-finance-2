package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uuid.UUID      `gorm:"column:id;type:uuid;default:uuid_generate_v4()" json:"id"`
	FirstName    string         `gorm:"column:first_name;type:varchar"`
	LastName     string         `gorm:"column:last_name;type:varchar"`
	PhoneNumber  string         `gorm:"unique"`
	Address      string         `gorm:"column:address;type:varchar"`
	PIN          string         `gorm:"column:pin;type:varchar"`
	Balance      float64        `gorm:"column:balance;type:decimal"`
	RefreshToken string         `gorm:"column:refresh_token;type:varchar"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamp"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamp"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp"`
}
