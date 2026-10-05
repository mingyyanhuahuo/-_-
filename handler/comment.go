package handler

import (
	"lostfound/dto"
	"lostfound/pkg/errcode"
	"lostfound/pkg/response"
	"lostfound/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListComments(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrLostInfoNotFound)
		return
	}
	var query struct {
		Page     int `form:"page" binding:"omitempty,min=1"`
		PageSize int `form:"pageSize" binding:"omitempty,min=1,max=50"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.ListComments(uint(itemID), query.Page, query.PageSize)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}

func CreateComment(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrLostInfoNotFound)
		return
	}
	var req dto.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.CreateComment(c.GetUint("id"), uint(itemID), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}

func DeleteComment(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrLostInfoNotFound)
		return
	}
	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrResourceNotFound)
		return
	}
	if err := service.DeleteComment(c.GetUint("id"), uint(itemID), uint(commentID), c.GetString("role")); err != nil {
		c.Error(err)
		return
	}
	response.OK(c, nil)
}