package router

import (
	"lostfound/handler"
	"lostfound/middleware"
	"lostfound/model"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
	apiJWT := api.Group("", middleware.JWTAuthMiddleware())
	apiJWT.GET("/announcements", handler.AnnouncementList)
	apiJWT.GET("/announcements/:announcementId", handler.GetAnnouncement)
	apiJWT.GET("/items/:itemId", handler.GetItemDetail)
	apiJWT.GET("/items", handler.ListItems)

	login := api.Group("",
		middleware.JWTAuthMiddleware(),
		middleware.RequireAuthMiddleware())
	{
		login.POST("/files/upload", handler.UploadFile)
		// login.DELETE("/files/:fileId", handler.DeleteFile)

		login.POST("/items", handler.GenerateItem)
		login.GET("/items/mine", handler.ListMyItems)
		login.PUT("/items/:itemId", handler.UpdateItem)
		login.DELETE("/items/:itemId", handler.DeleteItem)
		login.PATCH("/items/:itemId/status", handler.UpdateItemStatus)
	}
	root := api.Group("",
		middleware.JWTAuthMiddleware(),
		middleware.RequireAuthMiddleware(),
		middleware.RequireRoleMiddleware(model.RoleSysAdmin))
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
