package models

import (
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"time"
)

// 初始化雪花节点
var node *snowflake.Node

func init() {
	// 使用时间戳初始化节点
	var err error
	node, err = snowflake.NewNode(1) // 节点 ID，可以设置为机器的唯一标识
	if err != nil {
		panic(err)
	}
}

// 用户模型
type User struct {
	ID           int64     `json:"id" gorm:"primaryKey"`                         // 主键，自动增长
	Email        string    `json:"email" binding:"required" gorm:"unique;index"` // 唯一约束
	Nickname     string    `json:"nickname" binding:"required"`
	Password     string    `json:"password" binding:"required"`
	RegisteredAt time.Time `json:"registered_at"`
	LastLoginAt  time.Time `json:"last_login_at"`
	ExpiryAt     time.Time `json:"expiry_at"` // 会员到期时间
}

// 使用 BeforeCreate 钩子自动生成雪花 ID
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = node.Generate().Int64() // 使用雪花算法生成唯一 ID
	return
}
