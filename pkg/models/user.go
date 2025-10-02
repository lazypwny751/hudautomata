package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	RFIDCardID string         `json:"rfid_card_id" gorm:"column:rfid_card_id;uniqueIndex;not null"`
	Name       string         `json:"name" gorm:"not null"`
	Email      string         `json:"email" gorm:"index"`
	Phone      string         `json:"phone"`
	Balance    float64        `json:"balance" gorm:"type:decimal(10,2);default:0"`
	IsActive   bool           `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	RFIDCardID string  `json:"rfid_card_id" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	Email      string  `json:"email" binding:"omitempty,email"`
	Phone      string  `json:"phone"`
	Balance    float64 `json:"balance"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name     string  `json:"name"`
	Email    string  `json:"email" binding:"omitempty,email"`
	Phone    string  `json:"phone"`
	IsActive *bool   `json:"is_active"`
}
