// 用于为长链生成短链
package app

import(
	"github.com/gin-gonic/gin"
)

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
	// 写入response
	context.JSON(200, gin.H{
		"shorturl": shortUrl,
	})
}

func generateShortUrl(longUrl string) string{
	return longUrl
}