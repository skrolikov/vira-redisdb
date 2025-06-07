package redisdb

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewClient создаёт и возвращает подключение к Redis.
// Принимает контекст, адрес, пароль и номер базы.
// Проверяет соединение, возвращает клиента и ошибку, если не удалось подключиться.
func NewClient(ctx context.Context, addr, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Ping с контекстом и таймаутом
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	log.Println("Redis подключён!")
	return rdb, nil
}
