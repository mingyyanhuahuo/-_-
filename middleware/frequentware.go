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
		err, ErrCode := redisdb.RequestsTimeLimit(ip)
		if err != nil {
			logger.Logger.Error("未知错误", zap.Error(err))
			c.Abort()
			return
		}
		if ErrCode != nil {
			response.Err(c, ErrCode)
			c.Abort()
			return
		}
	}
}
