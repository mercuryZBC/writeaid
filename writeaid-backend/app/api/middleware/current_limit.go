package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

func CurrentLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := rate.NewLimiter(rate.Every(time.Second), 100)
		if limiter.Allow() {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "请求过于频繁，请稍后重试"})
		}
	}
}
