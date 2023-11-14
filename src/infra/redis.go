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

	// rdb.Set(ctx, "mykey1", "hoge", 0) // キー名 mykey1で文字列hogeをセット
	// ret, err := rdb.Get(ctx, "mykey1").Result() // キー名mykey1を取得
	// if err != nil {
	// 	log.Fatalln(err)
	// 	return rdb
	// }

	// println("Result: ", ret)

	return rdb
}
