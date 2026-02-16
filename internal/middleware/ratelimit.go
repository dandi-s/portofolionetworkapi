package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.RWMutex
)

// RateLimit limits requests per IP
// limit: max requests per minute
func RateLimit(limit int) gin.HandlerFunc {
	// Cleanup old visitors every 5 minutes
	go cleanupVisitors()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		defer mu.Unlock()

		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &visitor{
				lastSeen: time.Now(),
				count:    1,
			}
			c.Next()
			return
		}

		// Reset count if minute passed
		if time.Since(v.lastSeen) > time.Minute {
			v.count = 1
			v.lastSeen = time.Now()
			c.Next()
			return
		}

		v.count++

		// Check limit
		if v.count > limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
				"retry_after": 60 - int(time.Since(v.lastSeen).Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 10*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}
