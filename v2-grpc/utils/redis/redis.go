package redis

import (
	"context"
	"log"
	"time"
	redis "github.com/redis/go-redis/v9"
	"math/rand"
)

// TODO: 重构，client连接池？

type RedisUtils struct {
	ServerAddr string
	client *redis.Client
	// mu sync.Mutex
}

func (re *RedisUtils) GetRedisClient() *redis.Client{
	// re.mu.Lock()
	// defer re.mu.Unlock()

	// 使用单例模式进行初始化
	// TODO: 是否应该用sync.Once?
	if re.client != nil {
		return re.client
	}

	re.client = redis.NewClient(&redis.Options{
		Addr: re.ServerAddr,
		Password: "", //暂时还没有设置密码
		DB: 0, //使用默认DB
	})

	// 检查是否连接成功
	if pong, err := re.client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("无法连接到Redis: %v", err)
	} else {
		log.Println("已经连接到Redis: ", pong)
	}

	return re.client
}

// 👇🏻 获取key对应的值
func (re *RedisUtils) GetKey(ctx context.Context, key string) (value any, exists bool) {
	client := re.GetRedisClient()
	result, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Redis里不存在键: %v", key)
		return nil, false
	}

	if err != nil {
		log.Printf("Redis访问出错: %v", err)
		return nil, false
	}

	return result, true
}

// 👇🏻 将某个键值对加入Redis(值为string)，并设置过期时间
func (re *RedisUtils) AddKeyEx(ctx context.Context, key string, value string, duration float64) error {
	client := re.GetRedisClient()
	// 为了防止缓存雪崩，生成一个1-3之间的随机数
	var randExtraTime float64 = 0
	if duration != 0 {
		randExtraTime = rand.Float64()*2 + 1
	}
	result := client.Set(ctx, key, value, time.Duration(duration+randExtraTime)*time.Hour)
	if result.Err() != nil {
		log.Printf("无法向Redis中添加键值对: %v, %v", key, value)
		return result.Err()
	}
	return nil
}

// 👇🏻 将某个键值对加入Redis(值为string)，无过期时间
func (re *RedisUtils) AddKey(ctx context.Context,key string, value string) error {
	return re.AddKeyEx(ctx, key, value, 0)
}

// 👇🏻 删除某个键
func (re *RedisUtils) DeleteKey(ctx context.Context, key string) error {
	client := re.GetRedisClient()
	_, err := client.Del(ctx, key).Result()
	if err != nil {
		log.Printf("无法删除Redis中的键: %v", err)
		return err
	}
	return nil
}

// 👇🏻 判断某个键是否已经过期
func (re *RedisUtils) IsExpired(ctx context.Context, key string) bool {
	client := re.GetRedisClient()
	ttl, err := client.TTL(ctx, key).Result()
	if err != nil {
		log.Printf("无法判断键%v是否已经过期", err)
		return false
	}
	return ttl == -2 // -2代表键不存在或者已经被删除, -1代表永久有效，大于0代表剩下的生存时间
}

// 👇🏻 检查某个值是否存在于指定布隆过滤器
func (re *RedisUtils) BFExists(ctx context.Context, filterName string, item string) bool {
	client := re.GetRedisClient()
	exists, err := client.BFExists(
		ctx,
		filterName,
		item).Result()
	if err != nil {
		log.Printf("无法检查%v是否存在于布隆过滤器%v中", item, filterName)
		return false
	}
	return exists
}

// 👇🏻 将某个值加入指定的布隆过滤器
func (re *RedisUtils) BFAdd(ctx context.Context, filterName string, item string) bool {
	client := re.GetRedisClient()
	result, err := client.BFAdd(ctx, filterName, item).Result()
	if err != nil {
		log.Printf("无法创建向布隆过滤器%v中添加%v", filterName, item)
		return false
	}
	return result
}


// 👇🏻 创建布隆过滤器
func (re *RedisUtils) BFReserve(ctx context.Context, filterName string, errorRate float64, capacity int64) error {
	client := re.GetRedisClient()
	_, err := client.BFReserve(
		ctx,
		filterName,
		errorRate,
		capacity).Result()
	if err != nil && err.Error() != "ERR item exists" {
		return err
	}
	return nil
}
