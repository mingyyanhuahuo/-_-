package handler

import (
	"lostfound/pkg/logger"
	"lostfound/pkg/response"
	"lostfound/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ListUsers(c *gin.Context) {
	var query struct {
		Role     string `form:"role" binding:"omitempty,oneof=student lf_admin sys_admin"`
		Keyword  string `form:"keyword" binding:"max=50"`
		Page     int    `form:"page" binding:"omitempty,min=1"`
		PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=50"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		BindError(c, err)
		return
	}

	result, err := service.ListUsers(query.Role, query.Keyword, query.Page, query.PageSize)
	if err != nil {
		c.Error(err)
		return
	}

	logger.Logger.Info("用户列表操作审计",
		zap.Uint("operatorId", c.GetUint("id")),
		zap.String("roleFilter", query.Role),
		zap.String("keyword", query.Keyword),
	)

	response.OK(c, result)
}
