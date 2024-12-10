package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
	"yuqueppbackend/service-base/models"
	"yuqueppbackend/service-base/util"
)

type RecentType int

const (
	View RecentType = iota
	Edit
	Comment
)

// DocDao 处理与 Document 表相关的数据库操作
type DocDao struct {
	db       *gorm.DB
	esClient *elasticsearch.Client
}

// NewDocDao 创建一个新的 DocDao 实例
func NewDocDao(db *gorm.DB, esClient *elasticsearch.Client) *DocDao {
	return &DocDao{db: db, esClient: esClient}
}

// CreateDocument 创建文档
func (dao *DocDao) CreateDocument(doc *models.Document) error {
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()
	return dao.db.Create(doc).Error
}

// GetDocumentByID 根据 ID 获取文档
func (dao *DocDao) GetDocumentByID(id int64) (*models.Document, error) {
	var doc models.Document
	err := dao.db.Preload("KnowledgeBase").First(&doc, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// GetDocumentsByKnowledgeBaseID 获取某知识库下的所有文档
func (dao *DocDao) GetDocumentsByKnowledgeBaseID(kbID int64) ([]models.Document, error) {
	var docs []models.Document
	err := dao.db.Where("knowledge_base_id = ?", kbID).Find(&docs).Error
	return docs, err
}

// UpdateDocument 更新文档
func (dao *DocDao) UpdateDocument(doc *models.Document) error {
	doc.UpdatedAt = time.Now()
	return dao.db.Save(doc).Error
}

// DeleteDocumentByID 根据 ID 删除文档
func (dao *DocDao) DeleteDocumentByID(id int64) error {
	return dao.db.Delete(&models.Document{}, id).Error
}

// IncrementViewCount 增加文档的浏览次数
func (dao *DocDao) IncrementViewCount(id uint) error {
	return dao.db.Model(&models.Document{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1")).Error
}

// AddTag 为文档添加标签
func (dao *DocDao) AddTag(id int64, tag string) error {
	var doc models.Document
	err := dao.db.First(&doc, id).Error
	if err != nil {
		return err
	}

	// 追加标签
	if doc.Tags == "" {
		doc.Tags = tag
	} else {
		doc.Tags += "," + tag
	}

	doc.UpdatedAt = time.Now()
	return dao.db.Save(&doc).Error
}

// 向 redis中写入文档浏览记录
//func (dao *DocDao) UpdateRecentDocumentViewInRedis(userId, docId string) error {
//	// 写入 Redis 最近浏览记录
//	key := "user_recent_docs:" + userId
//	timestamp := float64(time.Now().Unix())
//	util.GetRedisClient().ZAdd(context.Background(), key, &redis.Z{
//		Score:  timestamp,
//		Member: docId,
//	})
//	// 限制最近浏览记录的条数
//	maxRecentCount := 50
//	util.GetRedisClient().ZRemRangeByRank(context.Background(), key, 0, int64(-maxRecentCount-1))
//	return nil
//}

// UpdateRecentDocumentViewInRedis 更新最近浏览记录到 Redis
func (dao *DocDao) UpdateRecentDocumentInRedis(recentType RecentType, document models.Document, kbName, userId string) error {
	var key string
	if recentType == View {
		key = "user_recent_view_docs:" + userId
	}
	if recentType == Edit {
		key = "user_recent_edit_docs:" + userId
	}
	if recentType == Comment {
		key = "user_recent_comment_docs:" + userId
	}
	timestamp := float64(time.Now().Unix())

	// 构造 JSON 数据作为 Redis ZSet 的 Member
	member := map[string]interface{}{
		"doc_id":    document.ID,
		"doc_title": document.Title,
		"kb_id":     document.KnowledgeBaseID,
		"kb_name":   kbName,
	}
	memberJSON, err := json.Marshal(member)
	if err != nil {
		return err
	}

	// 向 Redis ZSet 添加数据
	util.GetRedisClient().ZAdd(context.Background(), key, &redis.Z{
		Score:  timestamp,
		Member: memberJSON,
	})

	// 限制最多保留 50 条最近记录
	maxRecentCount := 50
	util.GetRedisClient().ZRemRangeByRank(context.Background(), key, 0, int64(-maxRecentCount-1))

	return nil
}

// GetRecentDocumentViewsWithScoresFromRedis 获取用户最近浏览的文档记录（带分数）
func (dao *DocDao) GetRecentDocumentWithScoresFromRedis(recentType RecentType, userId string, start, end int64) ([]map[string]interface{}, error) {
	var key string
	if recentType == View {
		key = "user_recent_view_docs:" + userId
	}
	if recentType == Edit {
		key = "user_recent_edit_docs:" + userId
	}
	if recentType == Comment {
		key = "user_recent_comment_docs:" + userId
	}
	// 解析数据
	var recentDocuments []map[string]interface{}
	var cur int64 = 0
	for cur <= end-start {
		// 获取带分数的数据
		values, err := util.GetRedisClient().ZRevRangeWithScores(context.Background(), key, start+cur, start+cur).Result()
		if len(values) == 0 || err != nil {
			break
		}
		for _, value := range values {
			var document map[string]interface{}
			if err := json.Unmarshal([]byte(value.Member.(string)), &document); err != nil {
				// 如果解析失败，记录错误但继续处理其他记录
				log.Printf("Failed to parse recent document entry: %v", err)
				continue
			}
			log.Println(document)
			// 判断当前文档是否已经被删除
			if docId, ok := document["doc_id"].(float64); ok {
				document["doc_id"] = int64(docId)
			}
			docId, ok := document["doc_id"].(int64)
			if ok == false {
				log.Printf("Failed to parse recent document entry: %v", err)
				continue
			}

			if res, err := dao.GetDocumentByID(docId); err != nil || res == nil {
				_, err = util.GetRedisClient().ZRem(context.Background(), key, value.Member).Result()
				log.Printf("the document %v maybe was removed", document["doc_title"])
				continue
			}

			// 附加分数（时间戳）
			document["timestamp"] = value.Score
			recentDocuments = append(recentDocuments, document)
			cur++
		}
	}
	return recentDocuments, nil
}

func (dao *DocDao) SetDocumentContentHash(documentId int64, hashValue string) error {
	res := util.GetRedisClient().Set(context.Background(), "documentContentHash:"+strconv.FormatInt(documentId, 10), hashValue, time.Hour*24*7)
	return res.Err()
}

func (dao *DocDao) GetDocumentContentHashByDocumentId(documentId int64) (string, error) {
	res := util.GetRedisClient().Get(context.Background(), "documentContentHash:"+strconv.FormatInt(documentId, 10))
	return res.Result()
}

// 将文档数据插入到ES

func (dao *DocDao) InsertDocToES(document models.Document, content string) error {
	docIdStr := strconv.FormatInt(int64(document.ID), 10)
	doc := map[string]interface{}{
		"id":      docIdStr,
		"title":   document.Title,
		"content": content,
	}

	// 将文档转为json
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		log.Fatalf("Error encoding document: %s", err)
	}

	// 索引文档
	res, err := dao.esClient.Index(
		"document",
		&buf,
		dao.esClient.Index.WithDocumentID(docIdStr), // 自定义 _id
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (dao *DocDao) UpdateDocToES(documentID int64, title string, content string) error {
	// 创建需要更新的字段数据
	doc := map[string]interface{}{
		"doc": map[string]interface{}{
			"title":   title,
			"content": content,
		},
	}

	// 将更新数据转为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		log.Fatalf("Error encoding update data: %s", err)
	}

	// 执行更新操作
	res, err := dao.esClient.Update(
		"document",                        // 索引名称
		strconv.FormatInt(documentID, 10), // 文档 ID（字符串）
		&buf,
	)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	defer res.Body.Close()

	fmt.Println(res)
	return nil
}

func (dao *DocDao) DeleteDocFromES(documentID int64) error {
	// 执行删除操作
	res, err := dao.esClient.Delete(
		"document",                        // 索引名称
		strconv.FormatInt(documentID, 10), // 文档 ID
	)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	fmt.Println(res)
	return nil
}

// FindKB 查找知识库，使用 ownerId 和 id 作为查询条件
func (dao *DocDao) FindKB(ownerId, id int64) (kb models.KnowledgeBase, err error) {
	// 查询知识库，使用 ownerId 和 id 作为条件
	err = dao.db.Where("owner_id = ? AND id = ?", ownerId, id).First(&kb).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果没有找到记录，返回一个适当的错误
			return kb, fmt.Errorf("knowledge base not found for ownerId: %d, id: %d", ownerId, id)
		}
		// 其他数据库查询错误
		return kb, fmt.Errorf("failed to find knowledge base: %v", err)
	}

	// 返回查询到的知识库对象
	return kb, nil
}
