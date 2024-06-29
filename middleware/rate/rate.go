package rate

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nxadm/tail/ratelimiter"
)

func LimitDefault() func(c *gin.Context) {
	return Limit(1000, 1)
}

// Limit : leak * time.second with size
func Limit(size uint16, leak time.Duration) func(c *gin.Context) {
	bucket := ratelimiter.NewLeakyBucket(size, leak*time.Second)
	return func(c *gin.Context) {
		if !bucket.Pour(1) {
			c.String(http.StatusServiceUnavailable, "")
			c.Abort()
			return
		}
		c.Next()
	}
}
