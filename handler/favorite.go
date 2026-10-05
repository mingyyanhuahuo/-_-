package handler

import (
	"lostfound/pkg/errcode"
	"lostfound/pkg/response"
	"lostfound/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddFavorite(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrLostInfoNotFound)
		return
	}
	if err := service.AddFavorite(c.GetUint("id"), uint(itemID)); err != nil {
		c.Error(err)
		return
	}
	response.OK(c, nil)
}

func RemoveFavorite(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrLostInfoNotFound)
		return
	}
	if err := service.RemoveFavorite(c.GetUint("id"), uint(itemID)); err != nil {
		c.Error(err)
		return
	}
	response.OK(c, nil)
}

func ListFavorites(c *gin.Context) {
	var query struct {
		Page     int `form:"page" binding:"omitempty,min=1"`
		PageSize int `form:"pageSize" binding:"omitempty,min=1,max=50"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.ListFavorites(c.GetUint("id"), query.Page, query.PageSize)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}