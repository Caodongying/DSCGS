package main

// gateway server相当于shorturl_conversion_server和shorturl_generation_server的客户端
// 使用gin框架

import (
	"context"
	"log"
	pb_gen "dscgs/v2-grpc/service/shorturl_generation_service"
	"github.com/gin-gonic/gin"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

// -----------------------------------------------------------------------------
// 定义gateway server结构体
// -----------------------------------------------------------------------------
type gatewayServer struct {
	generationClient pb_gen.ShortURLGenerationServiceClient
}

func (gw *gatewayServer) doGenerateShortUrl(c *gin.Context) {
	// 从请求体获取数据
	var requestBody struct {
		OriginalUrl string `json:"original_url"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	originalUrl := requestBody.OriginalUrl

	// 调用grpc微服务
	request := &pb_gen.GenerateShortURLRequest{OriginalUrl: originalUrl}
	response, err := gw.generationClient.GenerateShortURL(context.Background(), request)

	// 处理调用结果
	if err != nil {
		log.Fatalf("Cannot generate short url: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{ // H 是 map[string]interface{} 的类型别名
			"detail": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"short_url": response.ShortUrl,
	})
}


func main() {
	// -------------------------------------------------------------------------
	// 创建两个客户端并连接到grpc微服务
	// -------------------------------------------------------------------------
	genClientConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法创建与localhost:50051的连接：%v", err)
	}
	defer genClientConn.Close()

	gateway := gatewayServer{
		generationClient: pb_gen.NewShortURLGenerationServiceClient(genClientConn),
	}

	// -------------------------------------------------------------------------
	// 定义路由并调用相应的grpc微服务
	// -------------------------------------------------------------------------
	router := gin.Default()
	router.POST("/generate-shorturl", gateway.doGenerateShortUrl)


	// -------------------------------------------------------------------------
	// 启动gateway服务器
	// -------------------------------------------------------------------------
	log.Println("Gateway server starting on: 8080")
	if err:= router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start gateway: %v", err)
	}
}


