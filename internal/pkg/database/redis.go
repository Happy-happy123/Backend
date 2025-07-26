package database

import (
  "context"
  "time"
  
  "github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, password string, db int) *redis.Client {
  client := redis.NewClient(&redis.Options{
    Addr:     addr,
    Password: password,
    DB:       db,
    PoolSize: 50, // 连接池大小
  })

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  if err := client.Ping(ctx).Err(); err != nil {
    panic("Redis连接失败: " + err.Error())
  }

  return client
}