package redisdb

import (
	"context"
	"lostfound/config"
	"lostfound/pkg/errcode"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func InitRedis() {
	redisConfig := config.GetConfig().Redis
	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + redisConfig.Port,
		Password: redisConfig.Password,
		DB:       0,
	})
	if err := Rdb.Ping(context.Background()).Err(); err != nil {
		panic("Redis连接失败: " + err.Error())
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
