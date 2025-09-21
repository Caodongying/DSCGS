package main

import (
	"google.golang.org/grpc"
	pb_gen "dscgs/v2-grpc/service/shorturl_generation_service"
	"context"
	"log"
	"net"
)

type generationServer struct {
	pb_gen.UnimplementedShortURLGenerationServiceServer // 嵌入，不是字段！不可以是server pb....
}

func (gs *generationServer) GenerateShortURL(ctx context.Context, req *pb_gen.GenerateShortURLRequest) (*pb_gen.GenerateShortURLResponse, error) {
	shortUrl := "placeholder"
	var err error = nil

	return &pb_gen.GenerateShortURLResponse{ShortUrl: shortUrl}, err
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