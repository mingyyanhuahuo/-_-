package handler

import (
	"lostfound/pkg/response"
	"lostfound/service"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		BindError(c, err)
		return
	}
	file, err := service.UploadFile(c.GetUint("id"), fileHeader)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, file)
}

// func DeleteFile(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("fileId"), 10, 64)
// 	if err != nil {
// 		c.Error(errcode.ErrFileNotFound)
// 		return
// 	}
// 	if err := service.DeleteFile(c.GetUint("id"), c.GetString("role"), uint(id)); err != nil {
// 		c.Error(err)
// 		return
// 	}
// 	response.OK(c, nil)
// }
