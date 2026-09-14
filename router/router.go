package router

import (
	"log"
	"lostfound/config"
	"lostfound/handler"
	"lostfound/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()
	r.Use(middleware.AccessLog(), middleware.ErrorMiddleware())
	port := config.GetConfig().Server.Port
	log.Printf("服务启动，监听端口: %s", port)

	auth := r.Group("/auth")
	auth.POST("/register", handler.Register)

	r.Run(":" + port)
}
