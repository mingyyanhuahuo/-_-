package redisdb

import (
	"context"
	"lostfound/config"
	"lostfound/pkg/errcode"
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

func RequestsTimeLimit(key string, limit int64, window time.Duration) (*errcode.BizError, error) {
	ctx := context.Background()

	count, err := Rdb.Incr(ctx, key).Result()
	if err != nil {
		return errcode.ErrInternalServer, err
	}

	if count == 1 {
		if err := Rdb.Expire(ctx, key, window).Err(); err != nil {
			Rdb.Del(ctx, key)
			return errcode.ErrInternalServer, err
		}
	}

	if count > limit {
		return errcode.ErrTooManyRequests, nil
	}
	return nil, nil
}
