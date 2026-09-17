package redisdb

import (
	"context"
	"lostfound/config"
	"lostfound/pkg/logger"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Available atomic.Bool

func Refresh() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	ok := Rdb.Ping(ctx).Err() == nil
	Available.Store(ok)
	return ok
}
func InitRedis() {
	redisConfig := config.GetConfig().Redis
	Rdb = redis.NewClient(&redis.Options{
		Addr:                  redisConfig.Host + ":" + redisConfig.Port,
		Password:              redisConfig.Password,
		DB:                    0,
		DialTimeout:           500 * time.Millisecond,
		ReadTimeout:           300 * time.Millisecond,
		WriteTimeout:          300 * time.Millisecond,
		MaxRetries:            -1,
		DialerRetries:         1,
		ContextTimeoutEnabled: true,
	})
	if !Refresh() {
		logger.Logger.Error("Redis无法运行，请检查配置和网络连接")
	}
}
