# Vira Redis

Пакет `vira-redisdb` предоставляет обёртку над клиентом Redis с интегрированным логированием и проверкой соединения.

## Особенности

- Простая инициализация с проверкой соединения
- Интеграция с `vira-logger` для структурированного логирования
- Автоматическое логирование статуса подключения и ошибок
- Возможность получения нативного клиента Redis
- Грейсфул shutdown с закрытием соединения

## Установка

```bash
go get github.com/skrolikov/vira-redisdb
```

## Использование

### Инициализация

```go
import (
    "context"
    "github.com/skrolikov/vira-redisdb"
    "github.com/skrolikov/vira-logger"
)

func main() {
    logger := logger.New(logger.Config{...})
    
    cfg := redisdb.Config{
        Addr:     "localhost:6379",
        Password: "", // нет пароля
        DB:       0,  // база по умолчанию
    }

    ctx := context.Background()
    rdb, err := redisdb.New(ctx, cfg, logger)
    if err != nil {
        // Обработка ошибки подключения
    }
    defer rdb.Close()
    
    // Использование Redis...
}
```

### Доступ к нативному клиенту

```go
client := rdb.Client()

// Пример использования
val, err := client.Get(ctx, "key").Result()
if err != nil {
    logger.Error("Ошибка получения ключа: %v", err)
}
```

## Конфигурация

Параметры `Config`:

- `Addr` - адрес Redis сервера (host:port)
- `Password` - пароль (если требуется)
- `DB` - номер базы данных

## Логирование

Пакет автоматически логирует:
- Успешное подключение (`✅ Подключение к Redis успешно`)
- Ошибки подключения (`❌ Ошибка подключения к Redis`)
- Закрытие соединения (`🔌 Соединение с Redis закрыто`)

Пример вывода:
```
[INFO] 2023-10-01T15:04:05Z redis.go:25 ✅ Подключение к Redis успешно: localhost:6379
[INFO] 2023-10-01T15:05:05Z redis.go:36 🔌 Соединение с Redis закрыто
```

## Лучшие практики

1. Всегда проверяйте ошибку при инициализации
2. Используйте `defer rdb.Close()` для гарантированного закрытия соединения
3. Для сложных операций получайте нативный клиент через `Client()`
4. Передавайте контекст с таймаутами для операций с Redis
5. Используйте один экземпляр Redis для всего приложения

## Пример с middleware

```go
func RedisMiddleware(rdb *redisdb.Redis) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), "redis", rdb.Client())
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```