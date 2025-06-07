package redisdb

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewClient(addr, password string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := rdb.Ping(Ctx).Result(); err != nil {
		log.Fatal("Redis недоступен:", err)
	}

	log.Println("Redis подключён!")
	return rdb
}
