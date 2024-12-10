package dao

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"strconv"
	"yuqueppbackend/service-base/models"
	"yuqueppbackend/service-base/util"
)

type CommentDAO struct {
	db *gorm.DB
}

// NewCommentDAO 创建一个新的 CommentDAO 实例
func NewCommentDAO(db *gorm.DB) *CommentDAO {
	return &CommentDAO{db: db}
}

// CreateComment 创建新评论
func (dao *CommentDAO) CreateComment(comment *models.DocumentComment) error {
	if err := dao.db.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

// GetCommentByID 根据评论 ID 获取评论
func (dao *CommentDAO) GetCommentByID(commentID int64) (*models.DocumentComment, error) {
	var comment models.DocumentComment
	if err := dao.db.Preload("User").
		Where("id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

// GetRootCommentsByDocumentID 获取某文档下的顶级评论（支持分页）
func (dao *CommentDAO) GetRootCommentsByDocumentID(documentID int64, page, pageSize int) ([]models.DocumentComment, int64, error) {
	var comments []models.DocumentComment
	var total int64

	offset := (page - 1) * pageSize
	if err := dao.db.Model(&models.DocumentComment{}).Where("document_id = ? and parent_id = null", documentID).
		Where("parent_id IS NULL").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dao.db.Preload("User").
		Where("document_id = ?", documentID).
		Where("parent_id IS NULL").
		Limit(pageSize).Offset(offset).
		Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// GetRootCommentsByDocumentID 获取某个顶级评论下的子评论（支持分页）
func (dao *CommentDAO) GetChildrenCommentsByRootIdFromRedis(rootCommentId int64, page, pageSize int64) ([]map[string]interface{}, error) {
	key := "rootComment:" + strconv.FormatInt(rootCommentId, 10)
	var childrenComments []map[string]interface{}
	values, err := util.GetRedisClient().ZRevRangeWithScores(context.Background(), key, page*pageSize, pageSize).Result()
	if err != nil {
		return nil, err
	}
	for _, value := range values {
		var childrenComment map[string]interface{}
		if err := json.Unmarshal([]byte(value.Member.(string)), &childrenComment); err != nil {
			// 如果解析失败，记录错误但继续处理其他记录
			log.Printf("Failed to parse recent document entry: %v", err)
			continue
		}
		log.Println(childrenComment)
		childrenComments = append(childrenComments, childrenComment)
	}
	return childrenComments, nil
}

// UpdateComment 更新评论
func (dao *CommentDAO) UpdateComment(comment *models.DocumentComment) error {
	if err := dao.db.Save(comment).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCommentByID 根据 ID 逻辑删除评论
func (dao *CommentDAO) DeleteCommentByID(commentID int64) error {
	if err := dao.db.Model(&models.DocumentComment{}).Where("id = ?", commentID).
		Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCommentsByDocumentID 根据文档 ID 逻辑删除所有评论
func (dao *CommentDAO) DeleteCommentsByDocumentID(documentID int64) error {
	if err := dao.db.Model(&models.DocumentComment{}).Where("document_id = ?", documentID).
		Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}

// GetRepliesByCommentID 获取某条评论的所有回复
func (dao *CommentDAO) GetRepliesByCommentID(commentID int64) ([]models.DocumentComment, error) {
	var replies []models.DocumentComment
	if err := dao.db.Preload("User").Where("parent_id = ?", commentID).
		Order("created_at ASC").Find(&replies).Error; err != nil {
		return nil, err
	}
	return replies, nil
}

// HasRepliesByCommentID 判断某条评论是否有子评论
func (dao *CommentDAO) HasRepliesByCommentID(commentID int64) (bool, error) {
	var count int64
	// 查询子评论的数量，避免加载所有子评论
	if err := dao.db.Model(&models.DocumentComment{}).
		Where("parent_id = ?", commentID).
		Count(&count).Error; err != nil {
		return false, err
	}

	// 如果子评论数量大于 0，则存在子评论
	return count > 0, nil
}

func (dao *CommentDAO) InsertCommentToRedis(dc models.DocumentComment) error {
	key := "comment:" + strconv.FormatInt(int64(dc.DocumentID), 10)

	member := map[string]interface{}{
		"comment_id":         dc.ID,
		"user_id":            dc.UserID,
		"nickname":           dc.User.Nickname,
		"doc_id":             dc.DocumentID,
		"comment_content":    dc.Content,
		"comment_created_at": dc.CreatedAt,
		"comment_updated_at": dc.UpdatedAt,
	}
	memberJSON, err := json.Marshal(member)
	if err != nil {
		return err
	}
	res := util.GetRedisClient().ZAdd(context.Background(), key, &redis.Z{
		Score:  float64(dc.CreatedAt.Unix()),
		Member: memberJSON,
	})
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (dao *CommentDAO) InsertReplyCommentToRedis(dc models.DocumentComment) error {
	key := "rootComment:" + strconv.FormatInt(*(dc.RootID), 10)
	parentComment, err := dao.GetCommentByID(*dc.ParentID)
	curComment, err := dao.GetCommentByID(dc.ID)
	if err != nil {
		return err
	}
	member := map[string]interface{}{
		"comment_id":                   dc.ID,
		"user_id":                      dc.UserID,
		"nickname":                     curComment.User.Nickname,
		"doc_id":                       dc.DocumentID,
		"root_comment_id":              dc.RootID,
		"parent_comment_id":            parentComment.ID,
		"parent_comment_user_id":       parentComment.UserID,
		"parent_comment_user_nickname": parentComment.User.Nickname,
		"comment_content":              dc.Content,
		"comment_created_at":           dc.CreatedAt,
		"comment_updated_at":           dc.UpdatedAt,
		"comment_like_count":           dc.LikeCount,
		"comment_dislike_count":        dc.DislikeCount,
	}
	memberJSON, err := json.Marshal(member)
	if err != nil {
		return err
	}
	res := util.GetRedisClient().ZAdd(context.Background(), key, &redis.Z{
		Score:  float64(dc.CreatedAt.Unix()),
		Member: memberJSON,
	})
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}
