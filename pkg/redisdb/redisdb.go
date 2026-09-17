package redisdb

import (
	"context"
	"lostfound/config"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"
	"strconv"
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

func RequestsTimeLimit(ip string) (error, *errcode.BizError) {
	ctx := context.Background()

	val_s, err1 := Rdb.Get(ctx, ip).Result()
	count, err := strconv.Atoi(val_s)
	if err != nil {
		return err, nil
	}
	if err1 == redis.Nil {
		if _, err := Rdb.Set(ctx, ip, 0, 60).Result(); err != nil {
			return err, nil
		}
	} else if count <= 10 {
		if _, err := Rdb.Incr(ctx, ip).Result(); err != nil {
			return err, nil
		}
	} else {
		return nil, errcode.ErrTooManyRequests
	}
	return nil, nil
}
