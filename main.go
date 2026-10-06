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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func startFileClean() {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("文件清理服务异常退出", zap.Any("error", r))
		}
	}()
	tricker := time.NewTicker(1 * time.Hour)
	defer tricker.Stop()
	for range tricker.C {
		n, err := service.CleanOldFiles()
		if err != nil {
			logger.Logger.Error("清理过期文件失败", zap.Error(err))
			continue
		}
		if n > 0 {
			logger.Logger.Info("清理过期文件成功", zap.Int("count", n))
		}
	}
}
func startViewSync() {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("浏览量同步服务异常退出", zap.Any("error", r))
		}
	}()
	if redisdb.Refresh() {
		dao.SyncViewToDB()
	}
	tricker := time.NewTicker(5 * time.Second)
	defer tricker.Stop()
	for range tricker.C {
		wasAvailable := redisdb.Available.Load()
		nowAvailable := redisdb.Refresh()
		if wasAvailable != nowAvailable {
			if nowAvailable {
				logger.Logger.Info("Redis服务恢复可用")
			} else {
				logger.Logger.Warn("Redis服务不可用")
			}
		}
		if nowAvailable {
			dao.SyncViewToDB()
		}
	}
}

func initDB() *gorm.DB {
	dsn := config.GetConfig().Mysql.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库连接池失败: %v", err)
	}
	sqlDB.SetMaxOpenConns(100) // 设置最大连接数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接最大生命周期
	return db
}

func startPublishScheduler() {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("定时发布公告服务异常退出", zap.Any("error", r))
		}
	}()
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
	go startViewSync()
	go startFileClean()
	if err := jwtutil.Init(config.GetConfig().JWT.Secret); err != nil {
		log.Fatalf("JWT初始化失败: %v", err)
	}
	uploadDir := service.UpLoadDir()
	r := gin.Default()
	r.Group("/uploads", middleware.FrequentWare(middleware.LimitRule{
		Namespace: "static",
		Limit:     600,
		Window:    1 * time.Minute,
		ByIP:      true,
	})).Static("", uploadDir)
	r.Use(middleware.AccessLog())
	r.Use(middleware.ErrorMiddleware())
	router.InitRouter(r)

	// 前端 SPA 静态资源托管
	dist := "./dist"
	r.Static("/assets", filepath.Join(dist, "assets"))
	r.NoRoute(func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/api") || strings.HasPrefix(p, "/uploads") {
			c.JSON(404, gin.H{"code": 404, "msg": "not found", "timestamp": time.Now().Unix()})
			return
		}
		rel := strings.TrimPrefix(p, "/")
		if info, err := os.Stat(filepath.Join(dist, rel)); err == nil && !info.IsDir() {
			c.File(filepath.Join(dist, rel))
			return
		}
		c.File(filepath.Join(dist, "index.html"))
	})

	port := config.GetConfig().Server.Port
	certFile := config.GetConfig().SSL.CertFile
	keyFile := config.GetConfig().SSL.KeyFile
	log.Printf("服务启动，监听端口: %s (HTTPS)", port)
	if err := r.RunTLS(":"+port, certFile, keyFile); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
