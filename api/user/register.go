package user

import (
	"github.com/gin-gonic/gin"
)

type RegisterBody struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Password  string `json:"password" binding:"required,min=8,max=20"`
	Nickname  string `json:"nickname" binding:"required,min=2,max=20"`
	StudentNo string `json:"studentNo" binding:"required,len=10"`
	Phone     string `json:"phone" binding:"required,len=11"`
	Email     string `json:"email" binding:"required,email"`
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var RB RegisterBody
		if err := c.ShouldBindJSON(&RB); err != nil {

		}
	}
}
