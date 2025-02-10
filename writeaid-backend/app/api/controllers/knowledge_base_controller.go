package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strconv"
	"time"
	pb "writeaid-backend/app/knowledgeBaseService"
	"writeaid-backend/config"
	"writeaid-backend/dao"
	"writeaid-backend/util"
)

type KnowledgeBaseController struct {
	kbDao          *dao.KBDAO
	kbSeviceClient pb.KnowledgeServiceClient
}

var docDao *dao.DocDao = dao.NewDocDao(util.GetDB(), util.GetElasticSearchClient())

func NewKnowledgeBaseController(dao *dao.KBDAO) *KnowledgeBaseController {
	// 连接 gRPC 服务器
	conn, err := grpc.NewClient(config.GetKnowledgeBaseServiceAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect knowledgeBase service: %v", err)
	}
	client := pb.NewKnowledgeServiceClient(conn)
	return &KnowledgeBaseController{kbDao: dao, kbSeviceClient: client}
}

// CreateKnowledgeBase 创建知识库
func (kc *KnowledgeBaseController) CreateKnowledgeBase(c *gin.Context) {
	var contextData struct {
		Id          int64  `json:"userid"`
		Name        string `json:"kb_name" binding:"required"`
		Description string `json:"kb_description"`
		IsPublic    bool   `json:"kb_is_public"`
	}
	if id, exists := c.Get("userid"); exists {
		contextData.Id = id.(int64)
	}

	if err := c.ShouldBindJSON(&contextData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createReq := pb.CreateKnowledgeBaseReq{
		KnowledgeBase: &pb.KnowledgeBase{
			KbName:        contextData.Name,
			KbDescription: contextData.Description,
			KbIsPublic:    contextData.IsPublic,
			KbCreateAt:    time.Now().Unix(),
			KbUpdateAt:    time.Now().Unix(),
			KbOwnerId:     contextData.Id,
		},
	}
	createResp, err := kc.kbSeviceClient.CreateKnowledgeBase(context.Background(), &createReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"kb_id":          strconv.FormatInt(createResp.KnowledgeBase.KbId, 10),
		"kb_name":        createResp.KnowledgeBase.KbName,
		"kb_description": createResp.KnowledgeBase.KbDescription,
		"kb_is_public":   createResp.KnowledgeBase.KbIsPublic,
		"kb_created_at":  time.Unix(createResp.KnowledgeBase.KbCreateAt, 0),
		"kb_updated_at":  time.Unix(createResp.KnowledgeBase.KbUpdateAt, 0),
	})
}

// 获取用户创建的所有知识库
func (kc *KnowledgeBaseController) GetKnowledgeBaseList(c *gin.Context) {
	userId, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}
	getListResp, err := kc.kbSeviceClient.GetKnowledgeBaseList(context.Background(), &pb.GetKnowledgeBaseListReq{
		UserId: userId.(int64),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var kbListData []map[string]interface{}
	for _, kb := range getListResp.KnowledgeBases {
		tmpMap := make(map[string]interface{})
		tmpMap["kb_id"] = strconv.FormatInt(kb.KbId, 10)
		tmpMap["kb_name"] = kb.KbName
		tmpMap["kb_description"] = kb.KbDescription
		tmpMap["kb_is_public"] = kb.KbIsPublic
		tmpMap["kb_created_at"] = kb.KbCreateAt
		tmpMap["kb_updated_at"] = kb.KbUpdateAt
		kbListData = append(kbListData, tmpMap)
	}
	c.JSON(http.StatusOK, gin.H{"knowledge_bases": kbListData})
	return
}

// GetKnowledgeBaseDetail 根据用户ID和知识库ID获取知识库详情
func (kc *KnowledgeBaseController) GetKnowledgeBaseDetail(c *gin.Context) {
	kbId := c.Param("kb_id")
	id, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "系统错误请稍后再试"})
		return
	}
	kbId64, _ := strconv.ParseInt(kbId, 10, 64)
	userId64, _ := id.(int64)

	getDetailResp, err := kc.kbSeviceClient.GetKnowledgeBaseDetail(context.Background(), &pb.GetKnowledgeBaseDetailReq{
		UserId: userId64,
		KbId:   kbId64,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"kb_id":          strconv.FormatInt(getDetailResp.KnowledgeBase.KbId, 10),
		"kb_owner_id":    strconv.FormatInt(getDetailResp.KnowledgeBase.KbOwnerId, 10),
		"kb_name":        getDetailResp.KnowledgeBase.KbName,
		"kb_description": getDetailResp.KnowledgeBase.KbDescription,
		"kb_is_public":   getDetailResp.KnowledgeBase.KbIsPublic,
		"kb_created_at":  time.Unix(getDetailResp.KnowledgeBase.KbCreateAt, 0),
		"kb_updated_at":  time.Unix(getDetailResp.KnowledgeBase.KbUpdateAt, 0),
	})
}

// UpdateKnowledgeBase 更新知识库,可以更新的字段：Name,Description,IsPublic
func (kc *KnowledgeBaseController) UpdateKnowledgeBase(c *gin.Context) {
	var contextData struct {
		Id          int64     `json:"userid"`
		KBId        int64     `json:"kb_id" binding:"required"`
		Name        string    `json:"kb_name" binding:"required"`
		Description string    `json:"kb_description" binding:"required"`
		IsPublic    bool      `json:"kb_is_public" binding:"required"`
		updated_at  time.Time `json:"kb_updated_at"`
	}
	if id, exists := c.Get("userid"); exists {
		contextData.Id = id.(int64)
	}
	if err := c.ShouldBindJSON(&contextData); err != nil {
		log.Println("结构绑定失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kb := &pb.KnowledgeBase{
		KbName:        contextData.Name,
		KbDescription: contextData.Description,
		KbIsPublic:    contextData.IsPublic,
		KbCreateAt:    time.Now().Unix(),
		KbUpdateAt:    time.Now().Unix(),
		KbOwnerId:     contextData.Id,
	}
	updateResp, err := kc.kbSeviceClient.UpdateKnowledgeBase(context.Background(), &pb.UpdateKnowledgeBaseReq{
		KnowledgeBase: kb,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"kb_id":          strconv.FormatInt(updateResp.KnowledgeBase.KbId, 10),
		"kb_name":        updateResp.KnowledgeBase.KbName,
		"kb_description": updateResp.KnowledgeBase.KbDescription,
		"kb_is_public":   updateResp.KnowledgeBase.KbIsPublic,
		"kb_created_at":  time.Unix(updateResp.KnowledgeBase.KbCreateAt, 0),
		"kb_updated_at":  time.Unix(updateResp.KnowledgeBase.KbUpdateAt, 0),
	})
}

// DeleteKnowledgeBase 删除知识库
func (kc *KnowledgeBaseController) DeleteKnowledgeBase(c *gin.Context) {
	var contextData struct {
		Id   int64  `json:"userid"`
		KBId string `json:"kb_id" binding:"required"`
	}
	if id, exists := c.Get("userid"); exists {
		contextData.Id = id.(int64)
	}
	if err := c.ShouldBindJSON(&contextData); err != nil {
		log.Println("结构绑定失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用 strconv.ParseInt 将字符串转换为 int64
	kbId64, err := strconv.ParseInt(contextData.KBId, 10, 64) // 10 是十进制，64 表示返回 int64 类型
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统错误"})
		return
	}
	docs, err := docDao.GetDocumentsByKnowledgeBaseID(kbId64)
	for _, doc := range docs {
		err := docDao.DeleteDocumentByID(doc.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "系统错误"})
			return
		}
		err = docDao.DeleteDocFromES(doc.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "系统错误"})
			return
		}
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统错误"})
		return
	}
	_, err = kc.kbSeviceClient.DeleteKnowledgeBase(context.Background(), &pb.DeleteKnowledgeBaseReq{
		UserId: contextData.Id,
		KbId:   kbId64,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
