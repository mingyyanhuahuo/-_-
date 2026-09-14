package user

import (
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.Engine) {
	r.POST("auth/register", Register())
}
