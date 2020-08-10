package router

import (
	"github.com/gin-gonic/gin"
)

func Use(engine *gin.Engine) {
	// middleware
	// interface
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
