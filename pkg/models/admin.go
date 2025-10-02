package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRole string

const (
	RoleSuperAdmin AdminRole = "super_admin"
	RoleAdmin      AdminRole = "admin"
)

type Admin struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	Username     string         `json:"username" gorm:"uniqueIndex;not null"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	Role         AdminRole      `json:"role" gorm:"default:'admin'"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	LastLogin    *time.Time     `json:"last_login"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate hook to generate UUID
func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (Admin) TableName() string {
	return "admins"
}

// LoginRequest represents the login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token     string    `json:"token"`
	Admin     Admin     `json:"admin"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CreateAdminRequest represents the request body for creating an admin
type CreateAdminRequest struct {
	Username string    `json:"username" binding:"required,min=3"`
	Email    string    `json:"email" binding:"required,email"`
	Password string    `json:"password" binding:"required,min=6"`
	Role     AdminRole `json:"role" binding:"required,oneof=admin super_admin"`
}
