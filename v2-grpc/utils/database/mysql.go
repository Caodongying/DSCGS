package database

import (
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"dscgs/v2-grpc/model"
)

var DB *gorm.DB

func Init() {
	dsn := "username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local" // data source name
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) // 可以设置慢查询日志等
	if err != nil {
		log.Fatalf("无法连接到MySQL: %v", err)
	}

	// 自动建表
	// 生产环境取消，并用数据库迁移工具，比如goose（小项目）或者golang-migrate/migrate
	if err := DB.AutoMigrate(&model.URLMapping{}); err != nil {
		log.Fatalf("无法创建URLMapping表: %v", err)
	}
}
