package myredis

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

const RedisServer string = "localhost:6379"

// 全局变量
var (
	redisClient *redis.Client
	once sync.Once
)

func GetRedisClient() *redis.Client{
	// 运用单例模式，创建Redis客户端

	once.Do(func(){
		redisClient = redis.NewClient(&redis.Options{
			Addr: RedisServer,
			Password: "", //暂时还没有设置密码
			DB: 0, //使用默认DB
		})

		// 检查是否连接成功
		if pong, err := redisClient.Ping(context.Background()).Result(); err != nil {
			fmt.Println("无法连接到Redis", err)
			panic(err)
		} else {
			fmt.Println("已经连接到Redis", pong)
		}
	})

	return redisClient
}

