package database

import (
	"github.com/ZXstrike/api-gateway/internal/config"
	"github.com/redis/go-redis/v9"
)

func RedisConnect(redisConf *config.RedisConfig) (*redis.Client, error) {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConf.Host + ":" + redisConf.Port,
		Password: redisConf.Password, // no password set
		DB:       redisConf.Database, // use default DB
		Protocol: 2,
	})

	return redisClient, nil
}
