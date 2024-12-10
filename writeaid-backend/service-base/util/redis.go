package util

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"yuqueppbackend/service-base/config"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

// GetRedisClient 获取 Redis 客户端实例
func GetRedisClient() *redis.Client {
	// 使用 sync.Once 确保只初始化一次 Redis 客户端
	once.Do(func() {
		redisConfig := config.GetRedisConfig()

		// 创建 Redis 客户端
		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisConfig["host"].(string) + ":" + redisConfig["port"].(string), // Redis 地址
			Password: redisConfig["password"].(string),                                  // Redis 密码
			DB:       redisConfig["db"].(int),                                           // 默认数据库
		})

		// 测试 Redis 连接
		_, err := redisClient.Ping(redisClient.Context()).Result()
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		}
		fmt.Println("Redis connected successfully!")
	})
	if redisClient == nil {
		panic("Redis client not initialized")
	}
	return redisClient
}
