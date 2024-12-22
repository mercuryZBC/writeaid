package models

import (
	"gorm.io/gorm"
	"time"
)

type NotifyType int

const (
	Comment NotifyType = iota
	DocumentOper
	KnowledgeBaseOper
	System
)

// PushNotification model represents a notification to be pushed
type NotificationHistory struct {
	ID            int64      `gorm:"primary_key;not null" json:"id"`
	ReceiveUserID int64      `gorm:"not null" json:"receive_userId"`
	NotifyType    NotifyType `json:"notify_type"`
	Content       string     `json:"message"`
	Link          string     `json:"link"`
	LinkDescribe  string     `json:"link_describe"`
	IsCheck       bool       `json:"is_check"`
	AvatarLink    string     `json:"avatar_link"`
	CreatedAt     time.Time  `json:"created_at"`
}

// 使用 BeforeCreate 钩子自动生成雪花 ID
func (push *NotificationHistory) BeforeCreate(tx *gorm.DB) (err error) {
	push.ID = node.Generate().Int64() // 使用雪花算法生成唯一 ID
	return
}
