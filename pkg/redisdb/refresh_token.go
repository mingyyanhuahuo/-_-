package redisdb

import (
	"context"
	"time"
)

func AddRevokedToken(tokenHash string, ttl time.Duration) error {
	if !Available.Load() {
		return nil
	}
	if ttl <= 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	return Rdb.Set(ctx, "revoked:"+tokenHash, "1", ttl).Err()
}

func IsRevokedToken(tokenHash string) (bool, error) {
	if !Available.Load() {
		return false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	n, err := Rdb.Exists(ctx, "revoked:"+tokenHash).Result()
	if err != nil {
		return false, nil
	}
	return n > 0, nil
}
