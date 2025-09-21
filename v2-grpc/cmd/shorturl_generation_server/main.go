package main

import (
	"google.golang.org/grpc"
	pb_gen "dscgs/v2-grpc/service/shorturl_generation_service"
	"context"
	"log"
	"net"
	redis "dscgs/v2-grpc/utils/redis"
)

type generationServer struct {
	pb_gen.UnimplementedShortURLGenerationServiceServer // 嵌入，不是字段！不可以是server pb....
}

func (gs *generationServer) GenerateShortURL(ctx context.Context, req *pb_gen.GenerateShortURLRequest) (*pb_gen.GenerateShortURLResponse, error) {
	originalUrl := req.OriginalUrl
	shortUrl := ""
	var err error = nil

	// ------ 检查Redis是否已经存在生成的短链 ------
	redisUtils := redis.RedisUtils{ServerAddr: "localhost:6379"}
	key := "long:" + originalUrl
	if result, exists := redisUtils.GetKey(key); exists {
		// ------- Redis里已经存有当前长链的信息 -------
		// 检查短链是否过期
		if isExpired := redisUtils.IsExpired(key); isExpired {
			// 短链已经过期
			shortUrl = createShortURL(originalUrl)
		} else {
			shortUrl = result.(string) // 类型断言，将any类型的result转换为string
		}
	} else {
		// ------- Redis里不存在当前长链的键 -------
		log.Println("Redis里不存在", key, "查询布隆过滤器......")
		// 检查布隆过滤器
		filterName := "GeneratedOriginalUrlBF"
		if exists = redisUtils.BFExists(filterName, originalUrl); exists {
			// 布隆过滤器里存在当前长链，访问MySQL

		} else {
			// 布隆过滤器里不存在当前长链，则数据库中也必然不存在长链信息
			// 生成短链并返回
			shortUrl = createShortURL(originalUrl)
		}

	}

	return &pb_gen.GenerateShortURLResponse{ShortUrl: shortUrl}, err
}

// 👇🏻 真正的开始生成短链的逻辑
func createShortURL(originalUrl string) string {
	// 获取分布式锁

	return ""
}

func main() {
	// 创建一个新的gRPC服务器实例
	grpcServer := grpc.NewServer() // 固定搭配
	// 创建短链注册服务器
	generationServer := &generationServer{}
	// 注册短链生成服务到gRPC服务器
	pb_gen.RegisterShortURLGenerationServiceServer(grpcServer, generationServer)

	// 监听端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 启动服务器
	log.Printf("Generation server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Generation server failed to serve: %v", err)
	}
}