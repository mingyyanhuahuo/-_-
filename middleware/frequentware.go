package middleware

import (
	"lostfound/pkg/logger"
	"lostfound/pkg/redisdb"
	"lostfound/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LimitRule struct {
	Namespace string
	Limit     int64
	Window    time.Duration
	ByIP      bool
}

func FrequentWare(rule LimitRule) gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string
		if rule.ByIP {
			key = rule.Namespace + ":ip:" + c.ClientIP()
		} else if id := c.GetUint("id"); id != 0 {
			key = rule.Namespace + ":u:" + strconv.FormatUint(uint64(id), 10)
		} else {
			key = rule.Namespace + ":ip:" + c.ClientIP()
		}
		if redisdb.Available.Load() {
			result, err := redisdb.RequestsTimeLimit(key, rule.Limit, rule.Window)
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
		}
		c.Next()
	}
}
