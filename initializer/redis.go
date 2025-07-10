package initializer

import (
	"context"
	"store/config"

	"github.com/redis/go-redis/v9"
)

func RedisInit() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisHost + ":" + config.RedisPort,
		DB: config.RedisDb,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return rdb, nil
}
