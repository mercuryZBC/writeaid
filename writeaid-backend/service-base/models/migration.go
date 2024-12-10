package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

// Migration 记录迁移信息
type Migration struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Version   string    `json:"version" gorm:"unique;not null"` // 版本号或变更标识
	AppliedAt time.Time `json:"applied_at"`                     // 应用时间
}

// 自动迁移数据库表结构并记录迁移信息
func MigrateDB(db *gorm.DB) error {
	// 自动迁移 User 和 Migration 模型
	if err := db.AutoMigrate(
		&User{},
		&Migration{},
		&KnowledgeBase{},
		&Document{},
		&DocumentComment{},
	); err != nil {
		return err
	}

	// 获取当前迁移的版本号，可以使用时间戳或其他标识
	version := fmt.Sprintf("v1.0-%s", time.Now().Format("20060102150405"))

	// 记录当前迁移版本
	migration := &Migration{
		Version:   version,
		AppliedAt: time.Now(),
	}

	// 插入迁移记录
	if err := db.Create(migration).Error; err != nil {
		return err
	}

	log.Printf("迁移版本 %s 已应用\n", version)
	return nil
}
