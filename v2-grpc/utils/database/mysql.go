package database

import (
	"log"
	"fmt"
	"os"
	"time"
	logger "gorm.io/gorm/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"dscgs/v2-grpc/model"
)

type DBConfig struct {
	Username string
	Password string
	Host string
	Port int
	DBName string
	LogMode bool
}

var DB *gorm.DB

func NewDBConnection(cfg DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	// 设置GORM日志级别
	var logLevel logger.LogLevel
	if cfg.LogMode {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             2 * time.Second, // 慢查询阈值
				LogLevel:                  logLevel,               // 日志级别
				Colorful:                  true,                  // 使用彩色打印
			},
		),
	})

	if err != nil {
		return nil, fmt.Errorf("无法连接到MySQL: %w", err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.URLMapping{})
}

