package middlewares

import "github.com/gin-gonic/gin"

var requestAmount int = 0

func GetMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestAmount += 1
		c.Next()
	}
}

func GetRequestAmount() int {
	return requestAmount
}
