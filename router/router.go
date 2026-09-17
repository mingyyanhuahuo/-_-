package router

import (
	"lostfound/handler"
	"lostfound/middleware"
	"lostfound/model"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.POST("/auth/register", handler.Register)
	api.GET("/announcements", handler.AnnouncementList)
	api.GET("/announcements/:announcementId", handler.GetAnnouncement)
	api.GET("/items/:itemId", handler.GetItemDetail)
	api.GET("/items", handler.ListItems)

	login := api.Group("", middleware.RequireAuthMiddleware())
	{
		login.POST("/files/upload", handler.UploadFile)
		login.DELETE("/files/:fileId", handler.DeleteFile)

		login.POST("/items", handler.GenerateItem)
		login.GET("/items/mine", handler.ListMyItems)
		login.PUT("/items/:itemId", handler.UpdateItem)
		login.DELETE("/items/:itemId", handler.DeleteItem)
		login.PATCH("/items/:itemId/status", handler.UpdateItemStatus)
	}
	root := api.Group("", middleware.RequireRoleMiddleware(model.RoleSysAdmin),
		middleware.RequireAuthMiddleware())
	{
		root.POST("/announcements", handler.GenerateAnnouncement)
		root.PUT("/announcements/:announcementId", handler.UpdateAnnouncement)
		root.PATCH("/announcements/:announcementId/publish", handler.ChangeAnnouncementStatus)
		root.DELETE("/announcements/:announcementId", handler.DeleteAnnouncement)
	}

	/*
		admin := api.Group("", middleware.RequireRoleMiddleware(model.RoleLfAdmin,model.RoleSysAdmin),
		middleware.RequireAuthMiddleware())
		{
			admin.POST("/users", handler.CreateUser)
		}
	*/
}
