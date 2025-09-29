package db

import (
	"context"
	"fmt"
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client
var redisEnabled bool

// InitFromConfig 根据 conf.Redis 初始化 Redis 客户端
func InitFromConfig(cfg conf.Redis) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	timeout := time.Duration(cfg.Timeout) * time.Second
	redisEnabled = true
	rdb = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.Database,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		MinIdleConns: cfg.Pool.MinIdle,
		PoolSize:     cfg.Pool.MaxActive,
	})
}

// RedisGet 获取 Redis 中的字符串值
func RedisGet(key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		// key 不存在
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", errors.WithStack(err)
	}
	return val, nil
}

// RedisSet 写入 Redis 字符串值
func RedisSet(key string, value string, expire time.Duration) error {
	return errors.WithStack(rdb.Set(ctx, key, value, expire).Err())
}

// RedisUpdate 更新 Redis 值（存在则更新，不存在则创建）
func RedisUpdate(key string, value string, expire time.Duration) error {
	// Redis 本身 Set 就是覆盖，所以直接调用 Set
	return RedisSet(key, value, expire)
}

// RedisDelete 删除 Redis key
func RedisDelete(key string) error {
	n, err := rdb.Del(ctx, key).Result()
	if err != nil {
		return errors.WithStack(err)
	}
	if n == 0 {
		return errors.New("key not found")
	}
	return nil
}

// RedisGetFunc 返回一个 func 获取 Redis 的值
func RedisGetFunc(key string, enabled bool) func() (string, error) {
	if !enabled {
		return nil
	}
	return func() (string, error) {
		return RedisGet(key)
	}
}

// RedisSetFunc 返回一个 func 写入 Redis 的值
func RedisSetFunc(key string, enabled bool, expire time.Duration) func(string) error {
	if !enabled {
		return nil
	}
	return func(value string) error {
		return RedisSet(key, value, expire)
	}
}

// IsRedisEnabled 判断 Redis 是否启用
func IsRedisEnabled() bool {
	return redisEnabled
}
