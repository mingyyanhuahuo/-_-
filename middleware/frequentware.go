package middleware

import (
	"lostfound/pkg/logger"
	"lostfound/pkg/redisdb"
	"lostfound/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FrequentWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		result, err := redisdb.RequestsTimeLimit(ip)
		if err != nil {
			logger.Logger.Error("频繁请求内部错误", zap.Error(err))
			response.Err(c, result)
			c.Abort()
			return
		}
		if result != nil {
			response.Err(c, result)
			c.Abort()
			return
		}

		c.Next()
	}
}
