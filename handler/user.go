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
