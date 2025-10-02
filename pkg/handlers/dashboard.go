package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lazypwny751/hudautomata/pkg/database"
	"github.com/lazypwny751/hudautomata/pkg/models"
)

// GetDashboardStats returns dashboard statistics
func GetDashboardStats(c *gin.Context) {
	var stats struct {
		TotalUsers        int64   `json:"total_users"`
		ActiveUsers       int64   `json:"active_users"`
		TotalBalance      float64 `json:"total_balance"`
		TodayTransactions int64   `json:"today_transactions"`
		TodayRevenue      float64 `json:"today_revenue"`
	}

	// Total users
	database.DB.Model(&models.User{}).Count(&stats.TotalUsers)
	
	// Active users
	database.DB.Model(&models.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers)
	
	// Total balance
	database.DB.Model(&models.User{}).Select("COALESCE(SUM(balance), 0)").Scan(&stats.TotalBalance)
	
	// Today's transactions
	database.DB.Model(&models.Transaction{}).
		Where("DATE(created_at) = CURRENT_DATE").
		Count(&stats.TodayTransactions)
	
	// Today's revenue (debit transactions)
	database.DB.Model(&models.Transaction{}).
		Where("DATE(created_at) = CURRENT_DATE AND type = ?", models.TypeDebit).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.TodayRevenue)

	c.JSON(http.StatusOK, stats)
}

// GetRecentActivities returns recent transactions
func GetRecentActivities(c *gin.Context) {
	var transactions []models.Transaction
	
	if err := database.DB.
		Preload("User").
		Preload("Admin").
		Order("created_at DESC").
		Limit(20).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetChartData returns data for charts
func GetChartData(c *gin.Context) {
	period := c.DefaultQuery("period", "week") // day, week, month

	var chartData []struct {
		Date   string  `json:"date"`
		Amount float64 `json:"amount"`
		Count  int64   `json:"count"`
	}

	var dateFormat string
	var dateGroup string

	switch period {
	case "day":
		dateFormat = "YYYY-MM-DD HH24:00:00"
		dateGroup = "DATE_TRUNC('hour', created_at)"
	case "month":
		dateFormat = "YYYY-MM-DD"
		dateGroup = "DATE_TRUNC('day', created_at)"
	default: // week
		dateFormat = "YYYY-MM-DD"
		dateGroup = "DATE_TRUNC('day', created_at)"
	}

	// Note: This is PostgreSQL syntax, for SQLite you'll need to adjust
	database.DB.Model(&models.Transaction{}).
		Select("TO_CHAR("+dateGroup+", '"+dateFormat+"') as date, "+
			"COALESCE(SUM(amount), 0) as amount, "+
			"COUNT(*) as count").
		Where("created_at >= NOW() - INTERVAL '7 days'").
		Group(dateGroup).
		Order("date ASC").
		Scan(&chartData)

	c.JSON(http.StatusOK, chartData)
}

// ListLogs returns system logs
func ListLogs(c *gin.Context) {
	var logs []models.SystemLog
	
	query := database.DB.Model(&models.SystemLog{}).Preload("Admin")
	
	// Filter by action
	if action := c.Query("action"); action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}

	// Filter by admin
	if adminID := c.Query("admin_id"); adminID != "" {
		query = query.Where("admin_id = ?", adminID)
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
	
	if err := query.Order("created_at DESC").Limit(100).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
	})
}
