// 用于为长链生成短链
package app

import (
	"context"
	"crypto/md5"
	"dscgs/utils"
	"dscgs/myredis"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

)


const ShortUrlServer string = "localhost:1234"

type getShortUrlRequest struct {
	Url string  `json:"url" binding:"required"` // 如果Url不大写就无法被ShouldBind()访问到
	// 需要json标签是表明结构体中的Url是和JSON数据中的url绑定，注意两者大小写不一样，不写标签就无法完全匹配
}

func getShortUrl(longUrl string) (bool, string){
	// 检查是否已经为初始的长链生成的短链
	// 1. 查询布隆过滤器，如果不存在就直接返回空串
	if !IsInBloomFilter(longUrl) {
		return false, ""
	}

	// 2. 访问Redis，若命中，直接返回短链
	redisClient := myredis.GetRedisClient()
	keyLongUrl := "long:" + longUrl
	result, err := redisClient.Get(context.Background(), keyLongUrl).Result()
	if err == redis.Nil {
		// keyLongUrl不存在
		fmt.Println("keyLongUrl不存在于Redis中", err)
	} else if err != nil {
		// 其他错误
		fmt.Println("Redis访问出错", err)
	} else {
		// err为空，key存在
		return true, result
	}

	// 3. 访问MySQL，若命中，写如Redis，并返回

	// 4. 若MySQL不命中，就生成
	return false, ""
}

func getShortUrlHandler(context *gin.Context){
	// 1. 解析请求参数
	var request getShortUrlRequest
	if err := context.ShouldBind(&request); err != nil {
		panic("参数有误")
	}

	longUrl := request.Url
	shortUrl := ""

	// 2. 检查是否已经为此长链生成了短链
	if hasGenerated, shortUrl := getShortUrl(longUrl); hasGenerated {
		context.JSON(200, gin.H{
			"shorturl": shortUrl,
		})
	}

	// 3. 生成短链
	shortUrl = generateShortUrl(longUrl)


	// TODO: 检查生成的短链是否存在 - 不确定是否需要，因为短链的生成方案已经很大程度避免重复了

	// TODO: 将(长链，短链)对，写入Redis并进行持久化保存

	// TODO: 允许自行设置短链有效期
	// TODO: 设置短链服务器的域名，开发期间暂时使用localhost （应该使用配置文件进行指定）

	// 写入response
	context.JSON(200, gin.H{
		"shorturl": shortUrl,
	})
}

/*
核心模块：短链生成

生成方法：（获取Redis的自增ID，转化为62进制后）-（Hash(长链加盐，时间戳作为盐)，取前5位）
*/
func generateShortUrl(longUrl string) string{
	redisClient := myredis.GetRedisClient()
	context := context.Background()
	// 获取自增ID，并进行Base62编码
	id := int(redisClient.Incr(context, "id").Val())
	id62 := utils.DecimalTo62(id)

	// 加盐，避免其他短链被解码猜到
	suffixHash := md5.Sum([]byte(fmt.Sprintf("%s%d", longUrl,time.Now().Unix())))
	suffixString := fmt.Sprintf("%x", suffixHash) // 将[16]byte数组（16个二进制数）转为16进制
	suffix := suffixString[:5] // 取前5位，左开右闭

	// 将十进制的id转为62进制
	shortUrl := id62 + "-" + suffix
	return shortUrl
}