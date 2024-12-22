package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"time"
	"writeaid-backend/models"
)

// NotificationHistoryDAO 负责对 NotificationHistory 数据表的操作
type NotificationHistoryDAO struct {
	db    *gorm.DB
	redis *redis.Client
	ctx   context.Context
}

// NewNotificationHistoryDAO 创建一个新的 NotificationHistoryDAO 实例
func NewNotificationHistoryDAO(db *gorm.DB, redis *redis.Client) *NotificationHistoryDAO {
	return &NotificationHistoryDAO{
		db:    db,
		redis: redis,
		ctx:   context.Background(),
	}
}

// CreateNotification 创建一条新的通知记录
func (dao *NotificationHistoryDAO) CreateNotification(notification *models.NotificationHistory) error {
	// 先插入数据库
	if err := dao.db.Create(notification).Error; err != nil {
		log.Printf("创建通知失败: %v", err)
		return err
	}
	// 将新创建的通知记录缓存到 Redis（根据通知 ID 存储）
	if err := dao.cacheNotification(notification); err != nil {
		log.Printf("缓存通知失败: %v", err)
	}
	err := dao.cacheUnreadNotification(notification.ReceiveUserID)
	if err != nil {
		log.Printf("缓存通知标记+1失败: %v", err)
		return err
	}

	return nil
}

// GetRecentNotifications 获取用户最近的 n 条通知
func (dao *NotificationHistoryDAO) GetRecentNotifications(userID int64, count int) ([]models.NotificationHistory, error) {
	var notifications []models.NotificationHistory

	// 获取用户的通知 ZSet 的键
	cacheKey := fmt.Sprintf("notifications:user:%d", userID)

	// 获取 ZSet 中最近的 count 条通知 ID（按时间倒序排列）
	notificationIDs, err := dao.redis.ZRevRange(dao.ctx, cacheKey, 0, int64(count-1)).Result()
	if err != nil {
		log.Printf("获取用户最近通知失败: %v", err)
		return nil, err
	}

	// 根据通知 ID 获取通知的详细信息
	for _, id := range notificationIDs {
		var notification models.NotificationHistory
		cacheKey := fmt.Sprintf("notification:%s", id)

		// 从 Redis 获取缓存的通知详细信息
		if err := dao.redis.Get(dao.ctx, cacheKey).Scan(&notification); err != nil {
			log.Printf("获取缓存的通知失败: %v", err)
			continue
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// cacheNotification 将通知缓存到 Redis
func (dao *NotificationHistoryDAO) cacheNotification(notification *models.NotificationHistory) error {
	cacheKey := fmt.Sprintf("notifications:user:%d", notification.ReceiveUserID)
	score := float64(notification.CreatedAt.Unix()) // 使用通知的创建时间作为 score

	// 将通知的 ID 和时间戳加入 ZSet
	if err := dao.redis.ZAdd(dao.ctx, cacheKey, &redis.Z{
		Score:  score,
		Member: notification.ID,
	}).Err(); err != nil {
		log.Printf("缓存通知失败: %v", err)
		return err
	}
	// 如果需要缓存通知内容，可以使用 Set 或 Hash 存储详细信息
	// 设置缓存的过期时间
	notificationJson, err := json.Marshal(notification)
	if err != nil {
		log.Printf("通知序列化为 JSON 失败: %v", err)
		return err
	}
	if err := dao.redis.Set(dao.ctx, fmt.Sprintf("notification:%d", notification.ID), notificationJson, 10*time.Minute).Err(); err != nil {
		log.Printf("缓存通知详细内容失败: %v", err)
		return err
	}

	return nil
}

// GetNotificationByID 根据通知 ID 获取通知记录
func (dao *NotificationHistoryDAO) GetNotificationByID(id int64) (*models.NotificationHistory, error) {
	var notification models.NotificationHistory
	cacheKey := fmt.Sprintf("notification:%d", id)

	// 首先尝试从 Redis 获取缓存数据
	cachedNotification, err := dao.redis.Get(dao.ctx, cacheKey).Result()
	if err == nil && cachedNotification != "" {
		// 如果缓存存在，直接返回缓存中的数据
		log.Printf("从缓存中获取通知：%s", cachedNotification)
		// 这里根据你的需求可以直接返回缓存的数据
		return &notification, nil
	}

	// 如果缓存不存在，则从数据库中获取
	if err := dao.db.First(&notification, id).Error; err != nil {
		log.Printf("根据 ID 获取通知失败: %v", err)
		return nil, err
	}

	// 将从数据库获取的数据缓存到 Redis
	if err := dao.cacheNotification(&notification); err != nil {
		log.Printf("将通知缓存失败: %v", err)
	}

	return &notification, nil
}

// GetNotificationsByUserID 根据用户 ID 获取用户的所有通知
func (dao *NotificationHistoryDAO) GetNotificationsByUserID(userID int64, limit, offset int) ([]models.NotificationHistory, error) {
	var notifications []models.NotificationHistory

	// 尝试从缓存获取数据
	cacheKey := fmt.Sprintf("user:%d:notifications", userID)
	cachedNotifications, err := dao.redis.Get(dao.ctx, cacheKey).Result()
	if err == nil && cachedNotifications != "" {
		log.Printf("从缓存中获取用户通知：%s", cachedNotifications)
		// 如果缓存存在，直接返回缓存中的数据
		// 此处需要解析缓存中的数据并返回
		return notifications, nil
	}

	// 如果缓存不存在，则从数据库获取
	if err := dao.db.Where("receive_user_id = ?", userID).Limit(limit).Offset(offset).Order("created_at desc").Find(&notifications).Error; err != nil {
		log.Printf("根据用户 ID 获取通知失败: %v", err)
		return nil, err
	}

	// 将数据库查询到的数据缓存到 Redis
	if err := dao.cacheUserNotifications(userID, notifications); err != nil {
		log.Printf("缓存用户通知失败: %v", err)
	}

	return notifications, nil
}

// cacheUserNotifications 缓存用户的通知列表到 Redis
func (dao *NotificationHistoryDAO) cacheUserNotifications(userID int64, notifications []models.NotificationHistory) error {
	cacheKey := fmt.Sprintf("user:%d:notifications", userID)
	// 在这里我们将整个通知列表缓存到 Redis
	// 一般来说，如果列表较大，可能会考虑分片或分页存储
	for _, notification := range notifications {
		// 将每个通知分别缓存
		if err := dao.cacheNotification(&notification); err != nil {
			log.Printf("缓存用户通知失败: %v", err)
			return err
		}
	}

	return dao.redis.Set(dao.ctx, cacheKey, notifications, 10*time.Minute).Err()
}

// DeleteNotification 删除通知记录，并清理缓存
func (dao *NotificationHistoryDAO) DeleteNotification(id int64) error {
	// 删除数据库中的通知
	if err := dao.db.Delete(&models.NotificationHistory{}, id).Error; err != nil {
		log.Printf("删除通知失败: %v", err)
		return err
	}

	// 删除 Redis 中的缓存
	cacheKey := fmt.Sprintf("notification:%d", id)
	if err := dao.redis.Del(dao.ctx, cacheKey).Err(); err != nil {
		log.Printf("删除缓存失败: %v", err)
	}

	return nil
}

// cacheUnreadNotification 增加未读通知数量
func (dao *NotificationHistoryDAO) cacheUnreadNotification(userID int64) error {
	cacheKey := fmt.Sprintf("unread_notifications:user:%d", userID)

	// 增加未读通知的数量
	if err := dao.redis.Incr(dao.ctx, cacheKey).Err(); err != nil {
		log.Printf("增加未读通知数量失败: %v", err)
		return err
	}
	return nil
}

// markNotificationAsRead 将通知标记为已读并减少未读通知数量
func (dao *NotificationHistoryDAO) markNotificationAsRead(userID int64, notificationID int64) error {
	cacheKey := fmt.Sprintf("unread_notifications:user:%d", userID)

	// 减少未读通知的数量
	if err := dao.redis.Decr(dao.ctx, cacheKey).Err(); err != nil {
		log.Printf("减少未读通知数量失败: %v", err)
		return err
	}

	// 标记通知为已读
	if err := dao.markAsChecked(notificationID); err != nil {
		log.Printf("标记通知为已读失败: %v", err)
		return err
	}
	return nil
}

// MarkAsChecked 将通知标记为已读
func (dao *NotificationHistoryDAO) markAsChecked(notificationID int64) error {
	// 更新数据库中的通知记录，将 is_check 字段设置为 true
	if err := dao.db.Model(&models.NotificationHistory{}).
		Where("id = ?", notificationID).
		Update("is_check", true).Error; err != nil {
		log.Printf("标记通知为已读失败: %v", err)
		return err
	}
	return nil
}

// GetUnreadNotificationCount 获取未读通知的数量
func (dao *NotificationHistoryDAO) GetUnreadNotificationCount(userID int64) (int64, error) {
	cacheKey := fmt.Sprintf("unread_notifications:user:%d", userID)

	// 获取未读通知的数量
	count, err := dao.redis.Get(dao.ctx, cacheKey).Int64()
	if err == redis.Nil {
		// 如果 Redis 中没有未读通知数量，说明用户没有未读通知
		count = 0
	} else if err != nil {
		log.Printf("获取未读通知数量失败: %v", err)
		return 0, err
	}

	return count, nil
}
