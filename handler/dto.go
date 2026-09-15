package handler

import (
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func BindError(c *gin.Context, err error) {
	logger.Logger.Warn("参数校验失败", zap.Error(err))
	c.Error(errcode.New(400, 10000, "请求参数错误"))
}
