package handle

import (
	"github.com/gin-gonic/gin"
)

func NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
}

func Health(c *gin.Context) {
	c.JSON(200, "ok")
}
