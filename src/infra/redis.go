package infra

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
		PoolSize: 1000,
	})

	return rdb
}
