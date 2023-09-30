package cache

import (
	"os"
	"strconv"
)

type ICache interface {
	SaveAccessRefreshTokens(userID uint, deviceID, accessToken, refreshToken string) error
}

type Cache struct {
	Connection ICache
}

func (cache *Cache) InitializeConnection() error {
	redisEnabled, _ := strconv.ParseBool(os.Getenv("REDIS_ENABLED"))
	if redisEnabled {
		redisClient, err := initializeRedisConnection()
		if err != nil {
			return err
		}
		redisCache := new(RedisCache)
		redisCache.client = redisClient
		cache.Connection = redisCache
	} else {
		cache.Connection = initializeNoCache()
	}
	return nil
}
