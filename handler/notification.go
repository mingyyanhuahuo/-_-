package handler

import (
	"lostfound/dto"
	"lostfound/pkg/response"
	"lostfound/service"

	"github.com/gin-gonic/gin"
)

func ListNotifications(c *gin.Context) {
	var req dto.NotificationListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.ListNotifications(c.GetUint("id"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}