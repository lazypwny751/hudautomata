package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SystemLog struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	AdminID    *uuid.UUID `json:"admin_id" gorm:"type:uuid;index"`
	Admin      *Admin    `json:"admin,omitempty" gorm:"foreignKey:AdminID"`
	Action     string    `json:"action" gorm:"not null;index"`
	Resource   string    `json:"resource"`
	ResourceID string    `json:"resource_id"`
	Details    string    `json:"details" gorm:"type:text"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at" gorm:"index"`
}

// BeforeCreate hook to generate UUID
func (l *SystemLog) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (SystemLog) TableName() string {
	return "system_logs"
}
