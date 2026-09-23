package handler

import (
	"lostfound/dto"
	"lostfound/pkg/errcode"
	"lostfound/pkg/response"
	"lostfound/service"

	"strconv"

	"github.com/gin-gonic/gin"
)

func ListAuditItems(c *gin.Context) {
	var req dto.ItemListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.ListAuditItems(c.GetString("role"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}

func AuditItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrLostInfoNotFound)
		return
	}
	var req dto.AuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BindError(c, err)
		return
	}
	err = service.AuditItem(c.GetString("role"), req.Remark, req.Action, c.GetUint("id"), uint(itemID))
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, nil)
}
func ListAuditClaims(c *gin.Context) {
	var req dto.ClaimListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.ListAuditClaims(c.GetString("role"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}

func AuditClaim(c *gin.Context) {
	claimID, err := strconv.ParseUint(c.Param("claimId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrClaimReqNotFound)
		return
	}
	var req dto.AuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BindError(c, err)
		return
	}
	err = service.AuditClaim(c.GetString("role"), req.Remark, req.Action, c.GetUint("id"), uint(claimID))
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, nil)
}
