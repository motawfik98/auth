package cache

import (
	"github.com/redis/go-redis/v9"
	"os"
)

type Cache struct {
	connection *redis.Client
}

func (cache *Cache) SetCache(connection *redis.Client) {
	cache.connection = connection
}

func InitializeConnection() (*redis.Client, error) {
	opts, err := redis.ParseURL(os.ExpandEnv("redis://${REDIS_USERNAME}:${REDIS_PASSWORD}@${REDIS_HOST}:${REDIS_PORT}/#{REDIS_DB}"))
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opts), nil
}
