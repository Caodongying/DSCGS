// 用于为长链生成短链
package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strconv"
)

const RedisServer string = "localhost:6379"
const ShortUrlServer string = "localhost:1234"

type getShortUrlRequest struct {
	Url string  `json:"url" binding:"required"` // 如果Url不大写就无法被ShouldBind()访问到
	// 需要json标签是表明结构体中的Url是和JSON数据中的url绑定，注意两者大小写不一样，不写标签就无法完全匹配
}

func getShortUrlHandler(context *gin.Context){
	var request getShortUrlRequest
	if err := context.ShouldBind(&request); err != nil {
		panic("参数有误")
	}
	shortUrl := generateShortUrl(request.Url)

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

生成方法：获取Redis的自增ID，在末尾添加时间戳，然后进行Base62编码
*/
func generateShortUrl(longUrl string) string{
	// 创建Redis客户端
	context := context.Background()
	redisDb := redis.NewClient(&redis.Options{
		Addr: RedisServer,
		Password: "", //暂时还没有设置密码
		DB: 0, //使用默认DB
	})

	fmt.Println("已经连接到Redis")

	// 获取自增ID
	id := redisDb.Incr(context, "id").Val()

	// TODO: 转为Base62
	return strconv.FormatInt(id, 10)
}