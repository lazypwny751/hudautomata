package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lazypwny751/hudautomata/pkg/database"
	"github.com/lazypwny751/hudautomata/pkg/models"
)

// Logger middleware logs all requests to system_logs
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Skip logging for health checks
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/api/v1/ping" {
			return
		}

		// Create log entry
		log := models.SystemLog{
			Action:    c.Request.Method + " " + c.Request.URL.Path,
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			CreatedAt: time.Now(),
		}

		// Get admin ID from context if authenticated
		if adminID, exists := c.Get("admin_id"); exists {
			if id, ok := adminID.(uuid.UUID); ok {
				log.AdminID = &id
			}
		}

		// Add response details
		duration := time.Since(start)
		log.Details = fmt.Sprintf(`{"status":%d,"duration":"%s"}`, c.Writer.Status(), duration.String())

		// Save log to database (non-blocking)
		go func() {
			database.DB.Create(&log)
		}()
	}
}
