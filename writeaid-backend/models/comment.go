package models

import (
	"gorm.io/gorm"
	"time"
)

type DocumentComment struct {
	ID           int64     `json:"comment_id" gorm:"primaryKey"`       // 评论 ID，主键
	DocumentID   int64     `json:"document_id" gorm:"index"`           // 外键，关联文档，添加索引
	UserID       int64     `json:"user_id" gorm:"index"`               // 外键,发表评论的用户，添加索引
	Content      string    `json:"comment_content" binding:"required"` // 评论内容
	RootID       *int64    `json:"root_id" gorm:"index"`
	ParentID     *int64    `json:"parent_id" gorm:"index"`          // 自引用外键，父评论 ID（用于支持评论的回复），添加索引
	Status       string    `json:"status" gorm:"index"`             // 评论状态（如审核中、已发布、已删除等），添加索引
	LikeCount    uint      `json:"comment_like_count"`              // 点赞数量
	DislikeCount uint      `json:"comment_dislike_count"`           // 点踩数量
	CreatedAt    time.Time `json:"comment_created_at" gorm:"index"` // 评论创建时间，添加索引
	UpdatedAt    time.Time `json:"comment_updated_at" gorm:"index"` // 评论更新时间，添加索引
	IsDeleted    bool      `json:"comment_is_deleted"`              // 是否已删除，逻辑删除
	CreatedAtBy  string    `json:"comment_created_at_by"`           // 创建评论的IP地址或来源（用于审计）
	EditedAtBy   string    `json:"comment_edited_at_by"`            // 编辑评论的IP地址或来源（用于审计）

	// 关联的用户
	User User `json:"user" gorm:"foreignKey:UserID;references:ID"`

	// 关联的文档
	Document Document `json:"document" gorm:"foreignKey:DocumentID;references:ID"`

	// 评论的回复（子评论）
	Children []DocumentComment `json:"children" gorm:"foreignKey:ParentID"`

	// 是否为匿名评论
	IsAnonymous bool `json:"is_anonymous"`
}

// 使用 BeforeCreate 钩子自动生成雪花 ID
func (kb *DocumentComment) BeforeCreate(tx *gorm.DB) (err error) {
	kb.ID = node.Generate().Int64() // 使用雪花算法生成唯一 ID
	return
}
