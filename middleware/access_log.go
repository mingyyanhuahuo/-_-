package middleware

import (
	"errors"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		errMsg := ""
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var bizErr *errcode.BizError
			if errors.As(err, &bizErr) {
				errMsg = bizErr.Error()
			}
		}
		if errMsg != "" {
			logger.Logger.Warn("请求日志",
				zap.Uint("user_id", c.GetUint("id")),
				zap.Int("status", c.Writer.Status()),
				zap.String("ip", c.ClientIP()),
				zap.String("error", errMsg),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Duration("duration", time.Since(startTime)),
				zap.String("client_ip", c.ClientIP()),
			)
			return
		}
		logger.Logger.Info("请求日志",
			zap.Uint("user_id", c.GetUint("id")),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("duration", time.Since(startTime)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
