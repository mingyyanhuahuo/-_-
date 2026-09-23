package handler

import (
	"lostfound/pkg/errcode"
	"lostfound/pkg/response"
	"lostfound/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListCategories(c *gin.Context) {
	var query struct {
		IncludeDisabled bool `form:"includeDisabled"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		BindError(c, err)
		return
	}
	list, err := service.ListCategories(c.GetString("role"), query.IncludeDisabled)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, list)
}

func CreateCategory(c *gin.Context) {
	var req service.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BindError(c, err)
		return
	}
	category, err := service.CreateCategory(&req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, category)
}

func UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("categoryId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrCategoryInvalid)
		return
	}
	var req service.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BindError(c, err)
		return
	}
	category, err := service.UpdateCategory(uint(categoryID), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, category)
}

func DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("categoryId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrCategoryInvalid)
		return
	}
	if err := service.DeleteCategory(uint(categoryID)); err != nil {
		c.Error(err)
		return
	}
	response.OK(c, nil)
}
