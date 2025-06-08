package redisdb

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/skrolikov/vira-logger"
)

// Config описывает параметры подключения к Redis.
type Config struct {
	Addr     string
	Password string
	DB       int
}

// Redis инкапсулирует клиента Redis.
type Redis struct {
	client *redis.Client
	logger *log.Logger
}

// New создаёт новый экземпляр Redis и проверяет соединение.
func New(ctx context.Context, cfg Config, logger *log.Logger) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Пинг с таймаутом
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(pingCtx).Result(); err != nil {
		logger.Error("❌ Ошибка подключения к Redis: %v", err)
		return nil, err
	}

	logger.Info("✅ Подключение к Redis успешно: %s", cfg.Addr)

	return &Redis{client: rdb, logger: logger}, nil
}

// Client возвращает внутренний *redis.Client.
func (r *Redis) Client() *redis.Client {
	return r.client
}

// Close закрывает соединение с Redis.
func (r *Redis) Close() error {
	err := r.client.Close()
	if err != nil {
		r.logger.Error("❌ Ошибка закрытия Redis: %v", err)
	} else {
		r.logger.Info("🔌 Соединение с Redis закрыто")
	}
	return err
}
