package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string
type TransactionSource string

const (
	TypeCredit TransactionType = "credit"
	TypeDebit  TransactionType = "debit"
	TypeRefund TransactionType = "refund"

	SourceAdmin      TransactionSource = "admin"
	SourceAutomation TransactionSource = "automation"
	SourceSystem     TransactionSource = "system"
)

type Transaction struct {
	ID            uuid.UUID         `json:"id" gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID         `json:"user_id" gorm:"type:uuid;not null;index"`
	User          User              `json:"user,omitempty" gorm:"foreignKey:UserID"`
	AdminID       *uuid.UUID        `json:"admin_id" gorm:"type:uuid;index"`
	Admin         *Admin            `json:"admin,omitempty" gorm:"foreignKey:AdminID"`
	Type          TransactionType   `json:"type" gorm:"not null"`
	Amount        float64           `json:"amount" gorm:"type:decimal(10,2);not null"`
	BalanceBefore float64           `json:"balance_before" gorm:"type:decimal(10,2);not null"`
	BalanceAfter  float64           `json:"balance_after" gorm:"type:decimal(10,2);not null"`
	Description   string            `json:"description"`
	Source        TransactionSource `json:"source" gorm:"default:'admin'"`
	CreatedAt     time.Time         `json:"created_at" gorm:"index"`
}

// BeforeCreate hook to generate UUID
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name
func (Transaction) TableName() string {
	return "transactions"
}

// CreateTransactionRequest represents the request body for creating a transaction
type CreateTransactionRequest struct {
	UserID      uuid.UUID       `json:"user_id" binding:"required"`
	Type        TransactionType `json:"type" binding:"required,oneof=credit debit refund"`
	Amount      float64         `json:"amount" binding:"required,gt=0"`
	Description string          `json:"description"`
}

// AutomationScanRequest represents RFID scan request from automation device
type AutomationScanRequest struct {
	RFIDCardID  string  `json:"rfid_card_id" binding:"required"`
	ServiceCost float64 `json:"service_cost" binding:"required,gt=0"`
	Description string  `json:"description"`
}

// AutomationScanResponse represents the response for automation scan
type AutomationScanResponse struct {
	Success       bool      `json:"success"`
	UserID        uuid.UUID `json:"user_id,omitempty"`
	UserName      string    `json:"user_name,omitempty"`
	BalanceBefore float64   `json:"balance_before,omitempty"`
	BalanceAfter  float64   `json:"balance_after,omitempty"`
	TransactionID uuid.UUID `json:"transaction_id,omitempty"`
	CurrentBalance float64  `json:"current_balance,omitempty"`
	RequiredAmount float64  `json:"required_amount,omitempty"`
	Deficit       float64   `json:"deficit,omitempty"`
	Message       string    `json:"message"`
}
