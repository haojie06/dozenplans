package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// 时间测量
func Metic() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		useTime := time.Since(startTime)
		log.Printf("耗时:%fs", useTime.Seconds())
	}
}
