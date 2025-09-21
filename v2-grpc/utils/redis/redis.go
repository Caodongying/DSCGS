package redis

import (
	"context"
	"log"

	redis "github.com/redis/go-redis/v9"
)

// TODO: 重构，client连接池？

type RedisUtils struct {
	ServerAddr string
}

func (re *RedisUtils) GetRedisClient() *redis.Client{
	redisClient := redis.NewClient(&redis.Options{
		Addr: re.serverAddr,
		Password: "", //暂时还没有设置密码
		DB: 0, //使用默认DB
	})

	// 检查是否连接成功
	if pong, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("无法连接到Redis", err)
		panic(err)
	} else {
		log.Println("已经连接到Redis", pong)
	}

	return redisClient
}

// 👇🏻 获取key对应的值
func (re *RedisUtils) GetKey(key string) (value any, exists bool) {
	client := re.GetRedisClient()
	result, err := client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		log.Printf("Redis里不存在键", key)
		return nil, false
	}

	if err != nil {
		log.Fatalf("Redis访问出错: %v", err)
		panic(err)
	}

	return result, true
}


// 👇🏻 判断某个键是否已经过期
func (re *RedisUtils) IsExpired(key string) bool {
	client := re.GetRedisClient()
	ttl, err := client.TTL(context.Background(), key).Result()
	if err != nil {
		log.Fatalf("无法判断键%v是否已经过期", err)
		panic(err)
	}
	return ttl == -2 // -2代表键不存在或者已经被删除, -1代表永久有效，大于0代表剩下的生存时间
}

// 👇🏻 检查某个值是否存在于指定布隆过滤器
func (re *RedisUtils) BFExists(filterName string, item string) bool {
	client := re.GetRedisClient()
	exists, err := client.BFExists(
		context.Background(),
		filterName,
		item).Result()
	if err != nil {
		log.Fatalf("无法检查%v是否存在于布隆过滤器%v中", item, filterName)
	}
	return exists
}

// 👇🏻 将某个值加入指定的布隆过滤器
func (re *RedisUtils) BFAdd(filterName string, item string) bool {
	client := re.GetRedisClient()
	result, err := client.BFAdd(context.Background(), filterName, item).Result()
	if err != nil {
		log.Fatalf("无法创建向布隆过滤器%v中添加%v", filterName, item)
	}
	return result
}


// 👇🏻 创建布隆过滤器
func (re *RedisUtils) BFReserve(filterName string, errorRate float64, capacity int64) error {
	client := re.GetRedisClient()
	_, err := client.BFReserve(
		context.Background(),
		filterName,
		errorRate,
		capacity).Result()
	if err != nil && err.Error() != "ERR item exists" {
		return err
	}
	return nil
}
