package cache

import (
	"os"
	"strconv"
)

type iCache interface {
	SaveAccessRefreshTokens(userID uint, deviceID, accessToken, refreshToken string) error
	MarkRefreshTokenAsUsed(refreshToken string) (int64, error)
	IsUsedRefreshToken(refreshToken string) (bool, error)
	MarkRefreshTokensAsCompromised(refreshTokens []string) error
}

type Cache struct {
	Connection iCache
	Enabled    bool
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
		cache.Enabled = true
	} else {
		cache.Connection = initializeNoCache()
		cache.Enabled = false
	}
	return nil
}
