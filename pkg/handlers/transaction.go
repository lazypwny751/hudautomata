package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lazypwny751/hudautomata/pkg/database"
	"github.com/lazypwny751/hudautomata/pkg/models"
	"gorm.io/gorm"
)

// CreateTransaction creates a new transaction (admin initiated)
func CreateTransaction(c *gin.Context) {
	var req models.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, _ := c.Get("admin_id")
	adminUUID := adminID.(uuid.UUID)

	// Get user
	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Validate transaction
	if req.Type == models.TypeDebit && user.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Create transaction
	balanceBefore := user.Balance
	var balanceAfter float64

	if req.Type == models.TypeCredit || req.Type == models.TypeRefund {
		balanceAfter = balanceBefore + req.Amount
	} else {
		balanceAfter = balanceBefore - req.Amount
	}

	transaction := models.Transaction{
		UserID:        req.UserID,
		AdminID:       &adminUUID,
		Type:          req.Type,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   req.Description,
		Source:        models.SourceAdmin,
	}

	// Update user balance and create transaction in a transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		user.Balance = balanceAfter
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// ListTransactions returns all transactions with pagination
func ListTransactions(c *gin.Context) {
	var transactions []models.Transaction
	
	query := database.DB.Model(&models.Transaction{}).Preload("User").Preload("Admin")
	
	// Filter by user
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Filter by type
	if txType := c.Query("type"); txType != "" {
		query = query.Where("type = ?", txType)
	}

	// Filter by source
	if source := c.Query("source"); source != "" {
		query = query.Where("source = ?", source)
	}

	// Date range
	if from := c.Query("from"); from != "" {
		query = query.Where("created_at >= ?", from)
	}
	if to := c.Query("to"); to != "" {
		query = query.Where("created_at <= ?", to)
	}

	var total int64
	query.Count(&total)
	
	if err := query.Order("created_at DESC").Limit(100).Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  transactions,
		"total": total,
	})
}

// GetTransaction returns a single transaction
func GetTransaction(c *gin.Context) {
	id := c.Param("id")
	txID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var transaction models.Transaction
	if err := database.DB.Preload("User").Preload("Admin").First(&transaction, txID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// GetUserTransactions returns all transactions for a user
func GetUserTransactions(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var transactions []models.Transaction
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Admin").
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
