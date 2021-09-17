package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func MyLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		host := context.Request.Host
		url := context.Request.URL
		method := context.Request.Method
		log.Printf("%s::%s \t %s \t %s ", time.Now().Format("2021-06-01 12:00:00"), host, url, method)
		// 记得执行next
		context.Next()
		// 考虑将时间测量移动到这一起输出  		context.Set() context.Get()
	}
}
