package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func initializeRedisConnection() (*redis.Client, error) {
	opts, err := redis.ParseURL(os.ExpandEnv("redis://${REDIS_USERNAME}:${REDIS_PASSWORD}@${REDIS_HOST}:${REDIS_PORT}/#{REDIS_DB}"))
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opts), nil
}

func (cache *RedisCache) SaveAccessRefreshTokens(userID uint, deviceID, accessToken, refreshToken string) error {
	ctx := context.Background()
	accessTokenKey := userDeviceAccessTokenKey(userID, deviceID)
	refreshTokenKey := userDeviceRefreshTokenKey(userID, deviceID)

	pipe := cache.client.Pipeline()
	pipe.Set(ctx, accessTokenKey, accessToken, time.Hour*24)
	pipe.Set(ctx, refreshTokenKey, refreshToken, time.Hour*24*90)
	_, err := pipe.Exec(ctx)
	return err
}

func (cache *RedisCache) MarkRefreshTokenAsUsed(refreshToken string) (int64, error) {
	ctx := context.Background()
	usedRefreshTokensKey := usedRefreshTokensKey()
	return cache.client.SAdd(ctx, usedRefreshTokensKey, refreshToken).Result()
}

func (cache *RedisCache) IsUsedRefreshToken(refreshToken string) (bool, error) {
	ctx := context.Background()
	usedRefreshTokensKey := usedRefreshTokensKey()
	return cache.client.SIsMember(ctx, usedRefreshTokensKey, refreshToken).Result()
}
