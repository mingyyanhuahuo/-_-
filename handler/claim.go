package handler

import (
	"lostfound/dto"
	"lostfound/pkg/errcode"
	"lostfound/pkg/response"
	"lostfound/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GenerateClaim(c *gin.Context) {
	var req dto.ClaimCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.GenerateClaim(c.GetUint("id"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}
func ListMyClaims(c *gin.Context) {
	var req dto.ClaimListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.ListMyClaims(c.GetUint("id"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}

func GetClaimDetail(c *gin.Context) {
	claimID, err := strconv.ParseUint(c.Param("claimId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrClaimReqNotFound)
		return
	}
	result, err := service.GetClaimDetail(c.GetUint("id"), uint(claimID), c.GetString("role"))
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}
func CancelClaim(c *gin.Context) {
	claimID, err := strconv.ParseUint(c.Param("claimId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrClaimReqNotFound)
		return
	}
	if err := service.CancelClaim(c.GetUint("id"), uint(claimID)); err != nil {
		c.Error(err)
		return
	}
	response.OK(c, nil)
}
func ListItemClaims(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
	if err != nil {
		c.Error(errcode.ErrLostInfoNotFound)
		return
	}
	var req dto.ClaimListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		BindError(c, err)
		return
	}
	result, err := service.ListReceivedClaims(c.GetUint("id"), uint(itemID), c.GetString("role"), &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, result)
}
