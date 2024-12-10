package dao

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"gorm.io/gorm"
	"log"
	"strconv"
	"yuqueppbackend/service-base/models"
)

type KBDAO struct {
	DB       *gorm.DB
	esClient *elasticsearch.Client
}

func NewKBDAO(db *gorm.DB, esClient *elasticsearch.Client) *KBDAO {
	return &KBDAO{DB: db, esClient: esClient}
}

func (dao *KBDAO) CreateKB(kb *models.KnowledgeBase) error {
	return dao.DB.Create(kb).Error
}

func (dao *KBDAO) GetKBListByOwnerId(ownerId int64) ([]models.KnowledgeBase, error) {
	var knowledgeBases []models.KnowledgeBase
	if err := dao.DB.Where("owner_id", ownerId).Find(&knowledgeBases).Error; err != nil {
		return nil, err
	}
	return knowledgeBases, nil
}

func (dao *KBDAO) DeleteKB(kb models.KnowledgeBase) error {
	return dao.DB.Delete(&kb).Error
}

// FindKB 查找知识库，使用 ownerId 和 id 作为查询条件
func (dao *KBDAO) FindKB(ownerId, id int64) (kb models.KnowledgeBase, err error) {
	// 查询知识库，使用 ownerId 和 id 作为条件
	err = dao.DB.Where("owner_id = ? AND id = ?", ownerId, id).First(&kb).Error
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

// FindKB 查找知识库，使用 ownerId 和 id 作为查询条件
func (dao *KBDAO) GetKnowledgeBaseById(kbId int64) (kb models.KnowledgeBase, err error) {
	// 查询知识库，使用 ownerId 和 id 作为条件
	err = dao.DB.Where("id = ?", kbId).First(&kb).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果没有找到记录，返回一个适当的错误
			return kb, fmt.Errorf("knowledge base not found for id: %d", kbId)
		}
		// 其他数据库查询错误
		return kb, fmt.Errorf("failed to find knowledge base: %v", err)
	}

	// 返回查询到的知识库对象
	return kb, nil
}

// UpdateKB 更新知识库信息，使用 ownerId 和 id 作为查询条件
func (dao *KBDAO) UpdateKB(ownerId, id int64, updatedKB models.KnowledgeBase) (models.KnowledgeBase, error) {
	var kb models.KnowledgeBase

	// 查找指定 ownerId 和 id 的知识库
	err := dao.DB.Where("owner_id = ? AND id = ?", ownerId, id).First(&kb).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果没有找到记录，返回一个适当的错误
			return kb, fmt.Errorf("knowledge base not found for ownerId: %d, id: %d", ownerId, id)
		}
		// 其他数据库查询错误
		return kb, fmt.Errorf("failed to find knowledge base: %v", err)
	}

	// 更新字段，只更新传入的字段
	err = dao.DB.Model(&kb).Updates(updatedKB).Error
	if err != nil {
		// 更新失败，返回错误
		return kb, fmt.Errorf("failed to update knowledge base: %v", err)
	}

	// 返回更新后的知识库对象
	return kb, nil
}

func (dao *KBDAO) InsertKBToEs(kb models.KnowledgeBase) error {
	strKbId := strconv.FormatInt(kb.ID, 10)
	kbInfo := map[string]interface{}{
		"kb_id":          strKbId,
		"kb_name":        kb.Name,
		"kb_description": kb.Description,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(kbInfo); err != nil {
		log.Fatalf("Error encoding document: %s", err)
	}
	res, err := dao.esClient.Index(
		"knowledgebase",
		&buf,
		dao.esClient.Index.WithDocumentID(strKbId),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (dao *KBDAO) UpdateKBToES(kb models.KnowledgeBase) error {
	// 创建需要更新的字段数据
	doc := map[string]interface{}{
		"kb": map[string]interface{}{
			"kb_name":        kb.Name,
			"kb_description": kb.Description,
		},
	}

	// 将更新数据转为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		log.Fatalf("Error encoding update data: %s", err)
	}

	// 执行更新操作
	res, err := dao.esClient.Update(
		"knowledgebase",              // 索引名称
		strconv.FormatInt(kb.ID, 10), // 文档 ID（字符串）
		&buf,
	)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	defer res.Body.Close()

	fmt.Println(res)
	return nil
}

func (dao *KBDAO) DeleteKBFromES(knowledgeBaseID int64) error {
	// 执行删除操作
	res, err := dao.esClient.Delete(
		"knowledgebase",                        // 索引名称
		strconv.FormatInt(knowledgeBaseID, 10), // 文档 ID
	)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	fmt.Println(res)
	return nil
}
