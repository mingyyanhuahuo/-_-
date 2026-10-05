package handler

import (
	"lostfound/dto"
	"lostfound/pkg/errcode"
	"lostfound/pkg/response"
	"lostfound/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetItemMatches(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrLostInfoNotFound)
		return
	}
	var req dto.MatchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.GetItemMatches(uint(itemID), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}