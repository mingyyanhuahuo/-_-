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
	auth.POST("/register", middleware.AuthLimit, handler.Register)
	auth.POST("/login", middleware.AuthLimit, handler.Login)

	auth.Use(middleware.JWTAuthMiddleware(), middleware.QueryLimit)
	auth.POST("/refresh", handler.RefreshToken)
	auth.POST("/logout", handler.Logout)
	auth.GET("/me", handler.GetMe)
	auth.PUT("/password", middleware.AuthLimit, handler.UpdatePassword)

	apiJWT := api.Group("", middleware.JWTAuthMiddleware(), middleware.QueryLimit)
	apiJWT.GET("/announcements", handler.AnnouncementList)
	apiJWT.GET("/announcements/:announcementId", handler.GetAnnouncement)
	apiJWT.GET("/items/:itemId", handler.GetItemDetail)
	apiJWT.GET("/items", handler.ListItems)

	login := api.Group("",
		middleware.JWTAuthMiddleware(),
		middleware.RequireAuthMiddleware(), middleware.QueryLimit)
	{
		login.POST("/files/upload", middleware.UploadLimit, handler.UploadFile)
		// login.DELETE("/files/:fileId", handler.DeleteFile)

		login.POST("/items", middleware.PublishLimit, handler.GenerateItem)
		login.GET("/items/mine", handler.ListMyItems)
		login.PUT("/items/:itemId", handler.UpdateItem)
		login.DELETE("/items/:itemId", handler.DeleteItem)
		login.PATCH("/items/:itemId/status", handler.UpdateItemStatus)

		login.POST("/claims", middleware.ClaimLimit, handler.GenerateClaim)
		login.GET("/claims/mine", handler.ListMyClaims)
		login.GET("/items/:itemId/claims", handler.ListItemClaims)
		login.GET("/claims/:claimId", handler.GetClaimDetail)
		login.DELETE("/claims/:claimId", handler.CancelClaim)

	}
	root := api.Group("",
		middleware.JWTAuthMiddleware(),
		middleware.RequireAuthMiddleware(),
		middleware.RequireRoleMiddleware(model.RoleSysAdmin),
		middleware.QueryLimit)
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
