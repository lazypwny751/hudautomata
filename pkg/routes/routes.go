package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lazypwny751/hudautomata/pkg/handlers"
	"github.com/lazypwny751/hudautomata/pkg/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Logger middleware
	r.Use(middleware.Logger())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Public routes
			v1.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "pong"})
			})

			// Auth routes
			auth := v1.Group("/auth")
			{
				auth.POST("/login", handlers.Login)
				auth.POST("/logout", middleware.AuthMiddleware(), handlers.Logout)
				auth.GET("/me", middleware.AuthMiddleware(), handlers.GetMe)
			}

			// Automation routes (no auth required for IoT devices)
			automation := v1.Group("/automation")
			{
				automation.POST("/scan", handlers.AutomationScan)
				automation.POST("/check-balance", handlers.CheckBalance)
				automation.GET("/history", middleware.AuthMiddleware(), handlers.GetAutomationHistory)
			}

			// Protected routes (require authentication)
			protected := v1.Group("")
			protected.Use(middleware.AuthMiddleware())
			{
				// Users
				users := protected.Group("/users")
				{
					users.GET("", handlers.ListUsers)
					users.POST("", handlers.CreateUser)
					users.GET("/:id", handlers.GetUser)
					users.PUT("/:id", handlers.UpdateUser)
					users.DELETE("/:id", handlers.DeleteUser)
					users.GET("/rfid/:cardId", handlers.GetUserByRFID)
					users.GET("/:id/balance", handlers.GetUserBalance)
					users.GET("/:id/transactions", handlers.GetUserTransactions)
				}

				// Transactions
				transactions := protected.Group("/transactions")
				{
					transactions.GET("", handlers.ListTransactions)
					transactions.POST("", handlers.CreateTransaction)
					transactions.GET("/:id", handlers.GetTransaction)
				}

				// Dashboard
				dashboard := protected.Group("/dashboard")
				{
					dashboard.GET("/stats", handlers.GetDashboardStats)
					dashboard.GET("/charts", handlers.GetChartData)
					dashboard.GET("/recent", handlers.GetRecentActivities)
				}

				// Logs
				protected.GET("/logs", handlers.ListLogs)

				// Admins (super admin only)
				admins := protected.Group("/admins")
				admins.Use(middleware.SuperAdminOnly())
				{
					admins.GET("", handlers.ListAdmins)
					admins.POST("", handlers.CreateAdmin)
					admins.GET("/:id", handlers.GetAdmin)
					admins.DELETE("/:id", handlers.DeleteAdmin)
				}
			}
		}
	}
}

