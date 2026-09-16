package main

import (
	"log"
	"lostfound/config"
	"lostfound/dao"
	"lostfound/middleware"
	"lostfound/model"
	"lostfound/pkg/jwtutil"
	"lostfound/pkg/logger"
	"lostfound/pkg/redisdb"
	"lostfound/router"
	"lostfound/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dsn := config.GetConfig().Mysql.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	return db
}
func startPublishScheduler() {
	tricker := time.NewTicker(1 * time.Minute)
	defer tricker.Stop()
	for range tricker.C {
		n, err := dao.PublishSchedulerAnnouncement()
		if err != nil {
			logger.Logger.Error("定时发布公告失败", zap.Error(err))
			continue
		}
		if n > 0 {
			logger.Logger.Info("定时发布公告成功", zap.Int("count", int(n)))
		}
	}
}
func main() {
	logger.InitLogger()
	logger.Logger.Info("日志服务启动", zap.String("port", config.GetConfig().Server.Port))
	db := initDB()
	// 自动迁移数据库表
	err := db.AutoMigrate(&model.User{}, &model.Announcement{},
		&model.AuditLog{}, &model.Item{}, &model.Claim{},
		&model.File{}, &model.Follow{}, &model.RefreshToken{},
		&model.ItemType{}, &model.Comment{}, &model.Notification{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	dao.InitDB(db)
	go startPublishScheduler()
	redisdb.InitRedis()
	if err := jwtutil.Init(config.GetConfig().JWT.Secret); err != nil {
		log.Fatalf("JWT初始化失败: %v", err)
	}
	uploadDir := service.UpLoadDir()
	r := gin.Default()
	r.Static("/uploads", uploadDir)
	r.Use(middleware.AccessLog())
	r.Use(middleware.ErrorMiddleware())
	r.Use(middleware.JWTAuthMiddleware())
	router.InitRouter(r)
	port := config.GetConfig().Server.Port
	log.Printf("服务启动，监听端口: %s", port)
	r.Run(":" + port)
}
