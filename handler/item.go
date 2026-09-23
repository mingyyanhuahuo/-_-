package handler

import (
	"lostfound/dto"
	"lostfound/pkg/errcode"
	"lostfound/pkg/response"
	"lostfound/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GenerateItem(c *gin.Context) {
	var req dto.ItemCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.GenerateItem(c.GetUint("id"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}
func GetItemDetail(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrLostInfoNotFound)
		return
	}
	result, err := service.GetItemDetail(c.GetUint("id"), uint(itemID), c.GetString("role"))
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}
func UpdateItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrLostInfoNotFound)
		return
	}
	var req dto.ItemUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.UpdateItem(c.GetUint("id"), uint(itemID), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}
func DeleteItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrLostInfoNotFound)
		return
	}
	err = service.DeleteItem(c.GetUint("id"), uint(itemID), c.GetString("role"))
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, nil)
}

func ListItems(c *gin.Context) {
	var req dto.ItemListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.ListItems(c.GetString("role"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}
func ListMyItems(c *gin.Context) {
	var req dto.ItemListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.ListMyItems(c.GetUint("id"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}
func UpdateItemStatus(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrLostInfoNotFound)
		return
	}
	var req dto.UpdateItemStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BindError(c, err)
		return
	}
	if err := service.UpdateItemStatus(c.GetUint("id"), uint(itemID),
		c.GetString("role"), req.Status, req.Remark); err != nil {
		c.Error(err)
		return
	}
	response.OK(c, nil)
}
