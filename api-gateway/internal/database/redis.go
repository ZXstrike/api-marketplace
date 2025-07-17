package database

import (
	"context"
	"fmt"

	"github.com/ZXstrike/api-gateway/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func RedisConnect(redisConf *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConf.Host, redisConf.Port),
		Password: redisConf.Password,
		DB:       redisConf.Database,

		// --- Add Connection Pool Settings Here ---
		// Set the maximum number of connections. For a small VPS,
		// a value like 10-15 is a good starting point for Redis.
		PoolSize: 15,

		// Set the minimum number of idle connections.
		MinIdleConns: 5,
		// --- End of Connection Pool Settings ---
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	RedisClient = client

	return client, nil
}
