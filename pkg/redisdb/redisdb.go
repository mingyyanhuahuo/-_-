package redisdb

import (
	"context"
	"lostfound/config"

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
