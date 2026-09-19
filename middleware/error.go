package middleware

import (
	"errors"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"
	"lostfound/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Written() {
			return
		}
		errs := c.Errors
		if len(errs) == 0 {
			return
		}
		err := errs.Last().Err
		var bizErr *errcode.BizError
		if ok := errors.As(err, &bizErr); ok {
			response.Err(c, bizErr)
			return
		}
		logger.Logger.Error("中间件截取到的错误", zap.Error(err))
		response.Err(c, errcode.ErrInternalServer)
	}
}
