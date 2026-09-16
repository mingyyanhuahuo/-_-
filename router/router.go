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

	api.GET("/announcements", handler.AnnouncementList)
	api.GET("/announcements/:announcementId", handler.GetAnnouncement)

	admin := api.Group("/announcements",
		middleware.RequireAuthMiddleware(),
		middleware.RequireRole(model.RoleSysAdmin))
	admin.POST("", handler.GenerateAnnouncement)
	admin.PUT("/:announcementId", handler.UpdateAnnouncement)
	admin.PATCH("/:announcementId/publish", handler.ChangeAnnouncementStatus)
	admin.DELETE("/:announcementId", handler.DeleteAnnouncement)
}
