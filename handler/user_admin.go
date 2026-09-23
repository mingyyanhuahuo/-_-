package handler

import (
	"strconv"

	"lostfound/pkg/errcode"
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

func GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrBadRequest)
		return
	}
	result, err := service.GetUser(uint(userID))
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}

func DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrBadRequest)
		return
	}
	if err := service.DeleteUser(c.GetUint("id"), uint(userID)); err != nil {
		c.Error(err)
		return
	}
	response.OK(c, nil)
}
