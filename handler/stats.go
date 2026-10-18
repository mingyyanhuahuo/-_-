package handler

import (
	"lostfound/pkg/response"
	"lostfound/service"

	"github.com/gin-gonic/gin"
)

func Overview(c *gin.Context) {
	OR, err := service.Overview()
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, OR)
}

func Category(c *gin.Context) {
	var Body struct {
		ItemType     string `form:"type"`
		IncludeEmpty bool   `form:"includeEmpty"`
	}
	if err := c.ShouldBindQuery(&Body); err != nil {
		BindError(c, err)
		return
	}

	CR, err := service.Category(Body.ItemType, Body.IncludeEmpty)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, CR)
}

func Trend(c *gin.Context) {
	var Body struct {
		Days int `form:"days" binding:"omitempty,min=1,max=30"`
	}
	if err := c.ShouldBindQuery(&Body); err != nil {
		BindError(c, err)
		return
	}
	if Body.Days <= 0 {
		Body.Days = 7
	}

	TR, err := service.Trend(Body.Days)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, TR)
}

func ClaimRate(c *gin.Context) {
	var Body struct {
		StartTime string `form:"startTime" binding:"omitempty,datetime=2006-01-02"`
		EndTime   string `form:"endTime" binding:"omitempty,datetime=2006-01-02"`
	}
	if err := c.ShouldBindQuery(&Body); err != nil {
		BindError(c, err)
		return
	}

	CR, err := service.ClaimRate(Body.StartTime, Body.EndTime)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, CR)
}