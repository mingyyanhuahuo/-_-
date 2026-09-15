package handler

import (
	"lostfound/dto"
	"lostfound/pkg/errcode"
	"lostfound/pkg/response"
	"lostfound/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GenerateAnnouncement(c *gin.Context) {
	var req dto.AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errcode.ErrBadRequest)
		return
	}
	announcement, err := service.GenerateAnnouncement(int(c.GetUint("id")), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, announcement)
}
func GetAnnouncement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("announcementId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrAnnouncementNotFound)
		return
	}
	announcement, err := service.GetAnnouncement(uint(id), c.GetString("role"))
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, announcement)
}

func UpdateAnnouncement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("announcementId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrAnnouncementNotFound)
		return
	}
	var req dto.AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errcode.ErrBadRequest)
		return
	}
	announcement, err := service.UpdateAnnouncement(uint(id), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, announcement)
}

func ChangeAnnouncementStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("announcementId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrAnnouncementNotFound)
		return
	}
	var req struct {
		Action string `json:"action" binding:"required,oneof=publish offline"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errcode.ErrBadRequest)
		return
	}
	err = service.ChangeAnnouncementStatus(uint(id), req.Action)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, nil)
}

func DeleteAnnouncement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("announcementId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrAnnouncementNotFound)
		return
	}
	err = service.DeleteAnnouncement(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, nil)
}

func AnnouncementList(c *gin.Context) {
	var query struct {
		Status string `form:"status" binding:"omitempty,oneof=published draft offline"`
		Page   int    `form:"page" binding:"omitempty,min=1"`
		Size   int    `form:"pageSize" binding:"omitempty,min=1,max=50"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(errcode.ErrBadRequest)
		return
	}
	result, err := service.ListAnnouncements(c.GetString("role"), query.Status, query.Page, query.Size)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}
