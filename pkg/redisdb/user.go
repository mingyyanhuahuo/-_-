package redisdb

import (
	"context"
	"time"
)

func RemveToken(token string, expiration time.Duration) error {
	ctx := context.Background()

	_, err := Rdb.Set(ctx, token, 0, expiration).Result()
	if err != nil {
		return err
	}
	return nil
}
