package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	{
		api := r.Group("/api")
		{
			v1 := api.Group("/v1")
			{
				v1.GET("/ping", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "pong"})
				})
			}
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
