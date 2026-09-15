package router

import (
	"lostfound/handle"
	"lostfound/middleware"
	"lostfound/model"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	api := r.Group("/api/v1")

	api.GET("/announcements", handle.AnnouncementList)
	api.GET("/announcements/:announcementId", handle.GetAnnouncement)

	admin := api.Group("/announcements",
		middleware.RequireAuthMiddleware(),
		middleware.RequireRole(model.RoleSysAdmin))
	admin.POST("", handle.GenerateAnnouncement)
	admin.PUT("/:announcementId", handle.UpdateAnnouncement)
	admin.PATCH("/:announcementId/publish", handle.ChangeAnnouncementStatus)
	admin.DELETE("/:announcementId", handle.DeleteAnnouncement)
}
