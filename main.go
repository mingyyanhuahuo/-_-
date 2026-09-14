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

	redisdb.InitRedis()
	if err := jwtutil.Init(config.GetConfig().JWT.Secret); err != nil {
		log.Fatalf("JWT初始化失败: %v", err)
	}

	r := gin.Default()
	r.Use(middleware.AccessLog(), middleware.ErrorMiddleware())
	port := config.GetConfig().Server.Port
	log.Printf("服务启动，监听端口: %s", port)
	r.Run(":" + port)
}
