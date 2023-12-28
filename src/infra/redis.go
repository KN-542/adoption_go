package infra

import (
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	// "GO_ENV" が "dev" でない場合にエラーを出力
	if os.Getenv("GO_ENV") != "dev" {
		log.Fatalln("GO_ENV is not set to 'dev'. Please set it appropriately.")
	}

	// var ctx = context.Background()

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
