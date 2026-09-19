package handler

import (
	"lostfound/model"
	"lostfound/pkg/response"
	"lostfound/service"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var body model.RegisterBody
	if err := c.ShouldBindJSON(&body); err != nil {
		BindError(c, err)
		return
	}

	userId, err := service.Register(&body)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, gin.H{"userId": userId})

}

func Login(c *gin.Context) {
	var Body struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&Body); err != nil {
		BindError(c, err)
		return
	}

	LR, err := service.Login(Body.Username, Body.Password)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, LR)
}

func RefreshToken(c *gin.Context) {
	var Body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.ShouldBindJSON(&Body); err != nil {
		BindError(c, err)
		return
	}

	RR, err := service.RefreshToken(Body.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, RR)
}

func Logout(c *gin.Context) {
	var Body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.ShouldBindJSON(&Body); err != nil {
		BindError(c, err)
		return
	}

	if err := service.Logout(Body.RefreshToken); err != nil {
		c.Error(err)
		return
	}

	response.OK(c, "null")
}

func GetMe(c *gin.Context) {
	userid := c.GetUint("id")

	UF, err := service.GetMe(userid)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, UF)
}

func UpdatePassword(c *gin.Context) {
	var Body struct {
		OldPwd       string `json:"oldPassword"`
		NewPwd       string `json:"newPassword"`
		RefreshToken string `json:"rfreshToken"`
	}
	if err := c.ShouldBindJSON(&Body); err != nil {
		BindError(c, err)
		return
	}

	if err := service.UpdatePassaard(Body.OldPwd, Body.NewPwd, c.GetUint("id"), Body.RefreshToken); err != nil {
		c.Error(err)
		return
	}
	response.OK(c, "null")
}
