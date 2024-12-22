package models

import (
	"gorm.io/gorm"
	"time"
)

// KnowledgeBase 模型，作为主表
type KnowledgeBase struct {
	ID          int64     `json:"kb_id" gorm:"primaryKey"` // 使用 int64 存储雪花算法生成的 ID
	Name        string    `json:"kb_name" binding:"required"`
	Description string    `json:"kb_description"`
	IsPublic    bool      `json:"kb_is_public"`             // 是否公开
	OwnerID     int64     `json:"kb_owner_id" gorm:"index"` // 所有者
	CreatedAt   time.Time `json:"kb_created_at"`
	UpdatedAt   time.Time `json:"kb_updated_at"`

	User User `json:"user" gorm:"foreignKey:OwnerID;references:ID"`
	// 一对多关系
	Documents []Document `json:"directories" gorm:"foreignKey:KnowledgeBaseID;constraint:OnDelete:CASCADE;"`
}

// Document 模型，表示知识库中的文档
type Document struct {
	ID              int64     `json:"doc_id" gorm:"primaryKey"`     // 使用 int64 存储雪花算法生成的 ID
	Title           string    `json:"doc_title" binding:"required"` // 文档标题
	Content         string    `json:"doc_content"`                  // 文档内容
	OwnerId         int64     `json:"userid" gorm:"index"`
	KnowledgeBaseID int64     `json:"kb_id" gorm:"index"` // 外键，所属知识库
	ParentID        *int64    `json:"doc_parent_id"`      // 自引用外键，父文档 ID
	Status          string    `json:"doc_status"`         // 文档状态（如草稿、发布等）
	Tags            string    `json:"doc_tags"`           // 标签，逗号分隔
	ViewCount       uint      `json:"doc_view_count"`     // 浏览次数
	CommentCount    uint      `json:"doc_comment_count"`  // 评论数量
	CreatedAt       time.Time `json:"doc_created_at"`     // 创建时间
	UpdatedAt       time.Time `json:"doc_updated_at"`     // 更新时间
	Type            string    `json:"doc_type"`           // 文档类型（如文章、教程、参考等）

	// 关联的知识库
	KnowledgeBase KnowledgeBase `json:"knowledge_base" gorm:"foreignKey:KnowledgeBaseID;references:ID"`

	// 子文档，适用于文档的层级结构
	Children []Document `json:"children" gorm:"foreignKey:ParentID"`
}

// 使用 BeforeCreate 钩子自动生成雪花 ID
func (kb *KnowledgeBase) BeforeCreate(tx *gorm.DB) (err error) {
	kb.ID = node.Generate().Int64() // 使用雪花算法生成唯一 ID
	return
}

func (doc *Document) BeforeCreate(tx *gorm.DB) (err error) {
	doc.ID = node.Generate().Int64() // 使用雪花算法生成唯一 ID
	return
}
