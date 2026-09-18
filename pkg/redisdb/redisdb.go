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

func RequestsTimeLimit(ip string) (*errcode.BizError, error) {
	ctx := context.Background()

	_, err := Rdb.Get(ctx, ip).Result()
	if err != nil {
		if err == redis.Nil {
			if _, err := Rdb.Set(ctx, ip, "0", 60).Result(); err != nil {
				return errcode.ErrInternalServer, err
			}
			return nil, nil
		}
		return errcode.ErrInternalServer, err
	}

	count, err := Rdb.Incr(ctx, ip).Result()
	if err != nil {
		return errcode.ErrInternalServer, err
	} else if count > 11 {
		return errcode.ErrTooManyRequests, nil
	}
	return nil, nil
}
