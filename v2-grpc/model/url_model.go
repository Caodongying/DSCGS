package model

// 使用gorm定义数据库model

import (
	"time"
)

type URLMapping struct {
	ID uint64 `gorm:"primaryKey"` // 使用雪花算法计算的64位ID
	OriginalUrl string `gorm:"type:varchar(2048);not null"` // 不包含域名
	ShortUrl string `gorm:"type:varchar(9);uniqueIndex;not null"` // 不包含域名
	ExpireAt time.Time `gorm:"type:datetime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	AccessCount uint64 `gorm:"default:0"`
}

// check: https://www.tizi365.com/archives/6.html