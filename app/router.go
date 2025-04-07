/*
使用gin框架，定义路由
*/
package app

import (
	"github.com/gin-gonic/gin"
)


func StartServer(){
	// 1. 创建路由
	router := gin.Default()

	// 2. 绑定路由规则
	router.POST("/get-short-url", getShortUrlHandler) // gin.HandlerFunc函数不可以有返回值

	// 3. 启动路由
	router.Run(":1234")
}
