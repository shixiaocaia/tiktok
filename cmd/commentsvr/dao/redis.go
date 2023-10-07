package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/shixiaocaia/tiktok/cmd/commentsvr/config"
	"github.com/shixiaocaia/tiktok/cmd/commentsvr/log"

	"sync"
)

var (
	redisConn *redis.Client
	redSync   *redsync.Redsync
	redisOnce sync.Once
)

func InitRedis() {
	redisConfig := config.GetGlobalConfig().RedisConfig
	addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	redisConn = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisConfig.PassWord,
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSize,
	})
	if redisConn == nil {
		panic("failed to call redis.NewClient")
	}
	_, err := redisConn.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to ping redis, err:%s")
	}

	// 初始化分布式锁
	redisInstance := goredis.NewPool(redisConn)
	redSync = redsync.New(redisInstance)

	log.Infof("redis init success")
}

func CloseRedis() {
	if redisConn != nil {
		redisConn.Close()
	}
}

// GetRedisCli 获取数据库连接
func GetRedisCli() *redis.Client {
	redisOnce.Do(InitRedis)

	return redisConn
}

// GetRedSync 获取分布式锁
func GetRedSync() *redsync.Redsync {
	redisOnce.Do(InitRedis)

	return redSync
}
