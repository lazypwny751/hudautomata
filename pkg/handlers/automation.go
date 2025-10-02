package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazypwny751/hudautomata/pkg/database"
	"github.com/lazypwny751/hudautomata/pkg/models"
	"gorm.io/gorm"
)

// AutomationScan handles RFID scan from automation device
func AutomationScan(c *gin.Context) {
	var req models.AutomationScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by RFID
	var user models.User
	if err := database.DB.Where("rfid_card_id = ? AND is_active = ?", req.RFIDCardID, true).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, models.AutomationScanResponse{
			Success: false,
			Message: "RFID kartı kayıtlı değil veya kullanıcı aktif değil",
		})
		return
	}

	// Check balance
	if user.Balance < req.ServiceCost {
		c.JSON(http.StatusOK, models.AutomationScanResponse{
			Success:        false,
			UserID:         user.ID,
			UserName:       user.Name,
			CurrentBalance: user.Balance,
			RequiredAmount: req.ServiceCost,
			Deficit:        req.ServiceCost - user.Balance,
			Message:        "Yetersiz bakiye. Lütfen yöneticiye başvurun.",
		})
		return
	}

	// Process transaction
	balanceBefore := user.Balance
	balanceAfter := balanceBefore - req.ServiceCost

	transaction := models.Transaction{
		UserID:        user.ID,
		Type:          models.TypeDebit,
		Amount:        req.ServiceCost,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   req.Description,
		Source:        models.SourceAutomation,
	}

	// Update balance and create transaction atomically
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transaction"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, models.AutomationScanResponse{
		Success:       true,
		UserID:        user.ID,
		UserName:      user.Name,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		TransactionID: transaction.ID,
		Message:       "Hizmet verildi",
	})
}

// CheckBalance checks user balance without deducting
func CheckBalance(c *gin.Context) {
	var req struct {
		RFIDCardID string `json:"rfid_card_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("rfid_card_id = ?", req.RFIDCardID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":   user.ID,
		"user_name": user.Name,
		"balance":   user.Balance,
		"is_active": user.IsActive,
	})
}

// GetAutomationHistory returns automation transaction history
func GetAutomationHistory(c *gin.Context) {
	var transactions []models.Transaction
	
	query := database.DB.Where("source = ?", models.SourceAutomation).
		Preload("User").
		Order("created_at DESC")

	// Filter by date range
	if from := c.Query("from"); from != "" {
		query = query.Where("created_at >= ?", from)
	}
	if to := c.Query("to"); to != "" {
		query = query.Where("created_at <= ?", to)
	}

	if err := query.Limit(200).Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
