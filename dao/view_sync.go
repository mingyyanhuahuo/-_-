package dao

import (
	"context"
	"errors"
	"lostfound/model"
	"lostfound/pkg/logger"
	"lostfound/pkg/redisdb"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	viewKeyPrefix = "item:views:"
	viewKeyTTL    = 24 * time.Hour
)

func IncrItemView(itemID uint) (int64, bool) {
	if !redisdb.Available.Load() {
		return 0, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	key := viewKeyPrefix + strconv.FormatUint(uint64(itemID), 10)
	pipe := redisdb.Rdb.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, viewKeyTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		redisdb.Available.Store(false)
		logger.Logger.Error("Redis操作失败", zap.Error(err))
		return 0, false
	}
	return incr.Val(), true
}
func SyncViewToDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var cursor uint64
	for {
		keys, nextCursor, err := redisdb.Rdb.Scan(ctx, cursor, viewKeyPrefix+"*", 100).Result()
		if err != nil {
			logger.Logger.Error("Redis扫描 Key 失败", zap.Error(err))
			return
		}
		for _, key := range keys {
			if err := syncOneViewCount(ctx, key); err != nil {
				logger.Logger.Error("同步浏览量到数据库失败", zap.String("key", key), zap.Error(err))
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

}
func syncOneViewCount(ctx context.Context, key string) error {
	itemID, err := strconv.ParseUint(strings.TrimPrefix(key, viewKeyPrefix), 10, 64)
	if err != nil {
		redisdb.Rdb.Del(ctx, key)
		return err
	}
	delta, err := redisdb.Rdb.GetDel(ctx, key).Int64()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	if delta <= 0 {
		return nil
	}
	return db.Model(&model.Item{}).Where("id = ?", itemID).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", delta)).Error
}
