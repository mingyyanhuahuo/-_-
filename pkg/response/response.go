package response

import (
	"lostfound/pkg/errcode"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
)

type resp struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, resp{
		Code:      200,
		Msg:       "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}
func Err(c *gin.Context, bizErr *errcode.BizError) {
	c.JSON(bizErr.HttpStatus, resp{
		Code:      bizErr.Code,
		Msg:       bizErr.Message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}
