package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"yuqueppbackend/service-base/config"
	"yuqueppbackend/service-base/dao"
	"yuqueppbackend/service-base/models"
	"yuqueppbackend/service-base/util"
)

type DocumentController struct {
	docDao *dao.DocDao
}

func getDocumentStoragePath(docId string) string {
	return config.GetDocumentStoragePath() + "/" + docId + ".txt"
}

func getDocumentContentById(docId string) (string, error) {
	content, err := os.ReadFile(config.GetDocumentStoragePath() + "/" + docId + ".txt")
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(content), nil
}
func deleteDocumentFile(docId string) error {
	err := os.Remove(config.GetDocumentStoragePath() + "/" + docId + ".txt")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// NewDocumentController 创建新的 DocumentController
func NewDocumentController(docDao *dao.DocDao) *DocumentController {
	return &DocumentController{docDao: docDao}
}

// CreateDocumentHandler 创建文档
func (dc *DocumentController) CreateDocumentHandler(c *gin.Context) {

	var contextData struct {
		UserId int64  `json:"userid"`
		KbId   string `json:"kb_id" binding:"required"`
		Title  string `json:"doc_title"`
	}
	if id, exists := c.Get("userid"); exists {
		contextData.UserId = id.(int64)
	}

	if err := c.ShouldBindJSON(&contextData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	if contextData.Title == "" {
		contextData.Title = "无标题"
	}
	kbId64, err := strconv.ParseInt(contextData.KbId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	var doc models.Document = models.Document{
		KnowledgeBaseID: kbId64,
		Title:           contextData.Title,
		OwnerId:         contextData.UserId,
	}

	if err := dc.docDao.CreateDocument(&doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document"})
		return
	}
	str_doc_id := strconv.FormatInt(doc.ID, 10)
	doc_file, err := os.Create(config.GetDocumentStoragePath() + "/" + str_doc_id + ".txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	// 确保文件内容为空
	err = doc_file.Truncate(0)
	doc_content := "# " + doc.Title
	if _, err := doc_file.Write([]byte(doc_content)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}

	_ = dc.docDao.InsertDocToES(doc, doc_content)

	hashValue, err := util.HashDocumentContent(config.GetDocumentStoragePath() + "/" + str_doc_id + ".txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	err = dc.docDao.SetDocumentContentHash(doc.ID, hashValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	// 关闭文件
	defer doc_file.Close()

	c.JSON(http.StatusOK, gin.H{
		"doc_id":      str_doc_id,
		"kb_id":       doc.KnowledgeBaseID,
		"doc_title":   doc.Title,
		"doc_content": doc.Content,
	})
}

// GetDocumentByIDHandler 获取文档详情
func (dc *DocumentController) GetDocumentByIDHandler(c *gin.Context) {
	strDocId := c.Param("doc_id")
	docId, err := strconv.Atoi(strDocId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统错误请稍后重试"})
		return
	}
	// 判断用户请求文档是否存在
	doc, err := dc.docDao.GetDocumentByID(int64(docId))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document"})
		return
	}
	if doc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "当前文档不见了，快去新建吧"})
		return
	}

	userId, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "系统错误，请稍后再试"})
	}
	strUserId := strconv.FormatInt(userId.(int64), 10)
	strKbId := strconv.FormatInt(doc.KnowledgeBaseID, 10)

	// 获取文档内容
	docContent, err := getDocumentContentById(strDocId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}

	// 获取知识库名称
	kb, err := dc.docDao.FindKB(doc.OwnerId, doc.KnowledgeBaseID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve knowledge base"})
		return
	}
	kbName := kb.Name

	//将最近浏览记录写入到redis中
	err = dc.docDao.UpdateRecentDocumentInRedis(dao.View, *doc, kbName, strUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	// 在redis中写入文档内容哈希值
	hashValue, err := util.HashDocumentContent(config.GetDocumentStoragePath() + "/" + strDocId + ".txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	err = dc.docDao.SetDocumentContentHash(doc.ID, hashValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"doc_id":      strDocId,
		"kb_id":       strKbId,
		"doc_title":   doc.Title,
		"doc_content": docContent,
	})
	return
}

// GetDocumentsByKnowledgeBaseIDHandler 获取某知识库的所有文档
func (dc *DocumentController) GetDocumentsByKnowledgeBaseIDHandler(c *gin.Context) {
	kbIDParam := c.Param("kb_id")
	kbID, err := strconv.Atoi(kbIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}

	docs, err := dc.docDao.GetDocumentsByKnowledgeBaseID(int64(kbID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
		return
	}
	var docList []map[string]interface{}
	for _, doc := range docs {
		tmpMap := make(map[string]interface{})
		tmpMap["kb_id"] = strconv.FormatInt(doc.KnowledgeBaseID, 10)
		tmpMap["doc_id"] = strconv.FormatInt(doc.ID, 10)
		tmpMap["doc_title"] = doc.Title
		docList = append(docList, tmpMap)
	}
	c.JSON(http.StatusOK, gin.H{"doc_list": docList})
}

// UpdateDocumentHandler 更新文档
func (dc *DocumentController) UpdateDocumentHandler(c *gin.Context) {
	// 获取 doc_id 参数
	docIdStr := c.Param("doc_id")
	userId, _ := c.Get("userid")
	docId, err := strconv.ParseInt(docIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "文档不存在或已经删除，请重试"})
		return
	}
	doc, err := dc.docDao.GetDocumentByID(docId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "系统错误，请稍候再试"})
		return
	}
	if doc == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "系统错误，请稍候再试"})
		return
	}
	if doc.OwnerId != userId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无修改权限，修改失败"})
		return
	}

	// 获取知识库名称
	kb, err := dc.docDao.FindKB(doc.OwnerId, doc.KnowledgeBaseID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve knowledge base"})
		return
	}
	kbName := kb.Name
	// 获取上传的文件
	docFile, err := c.FormFile("file") // 注意字段名称是 'file'
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "系统错误，文件保存失败，请稍后再试"})
		return
	}
	// 保存文件到服务器
	err = c.SaveUploadedFile(docFile, getDocumentStoragePath(docIdStr))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "系统错误，文件保存失败，请稍后再试"})
		return
	}
	content, err := os.ReadFile(getDocumentStoragePath(docIdStr))
	if err != nil {
		return
	}
	strContent := string(content)
	err = dc.docDao.UpdateDocToES(docId, doc.Title, strContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "系统错误，文件保存失败，请稍后再试"})
		return
	}
	// 在redis中写入文档内容哈希值
	hashValue, err := util.HashDocumentContent(config.GetDocumentStoragePath() + "/" + docIdStr + ".txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	err = dc.docDao.SetDocumentContentHash(doc.ID, hashValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	err = dc.docDao.UpdateRecentDocumentInRedis(dao.Edit, *doc, kbName, strconv.FormatInt(userId.(int64), 10))
	if err != nil {
		log.Println("插入最近编辑记录到redis中失败")
		return
	}
	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "文件更新成功"})
}

// DeleteDocumentByIDHandler 删除文档
func (dc *DocumentController) DeleteDocumentByIDHandler(c *gin.Context) {
	// 从请求参数获取文档ID
	docIdStr := c.Param("doc_id")
	docId, err := strconv.ParseInt(docIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "错误的文档ID"})
		return
	}

	userID, exists := c.Get("userid") // 从上下文中获取当前用户的ID
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户没有权限"})
		return
	}

	// 查询文档的所有者
	document, err := dc.docDao.GetDocumentByID(int64(docId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document"})
		return
	}

	// 判断当前用户是否为文档的所有者
	if document.OwnerId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this document"})
		return
	}

	// 执行删除操作
	if err := dc.docDao.DeleteDocumentByID(int64(docId)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}
	deleteDocumentFile(docIdStr)

	err = dc.docDao.DeleteDocFromES(docId)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// IncrementViewCountHandler 增加浏览次数
func (dc *DocumentController) IncrementViewCountHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	if err := dc.docDao.IncrementViewCount(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment view count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "View count incremented"})
}

// GetRecentDocumentsHandler 获取最近浏览记录
func (dc *DocumentController) GetRecentViewDocumentsHandler(c *gin.Context) {
	// 默认获取最近 10 条记录，可以通过查询参数调整
	startParam := c.DefaultQuery("start", "0")
	limitParam := c.DefaultQuery("limit", "10")
	start, err := strconv.Atoi(startParam)
	if err != nil || start <= 0 {
		start = 0
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10 // 设置默认值
	}

	// 获取用户 ID（从中间件中设置的上下文）
	userId, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未授权"})
		return
	}
	strUserId := strconv.FormatInt(userId.(int64), 10)

	// 从 Redis 中获取最近浏览记录
	recentDocs, err := dc.docDao.GetRecentDocumentWithScoresFromRedis(dao.View, strUserId, int64(start), int64(limit-1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，获取最近浏览文档失败"})
		return
	}
	log.Print(recentDocs)
	// 构造返回数据
	var docList []map[string]interface{}
	for _, doc := range recentDocs {
		// 构造文档信息
		docMap := map[string]interface{}{
			"doc_id":    fmt.Sprintf("%d", doc["doc_id"].(int64)),  // 文档 ID
			"doc_title": doc["doc_title"],                          // 文档标题
			"kb_id":     fmt.Sprintf("%d", doc["kb_id"].(float64)), // 知识库 ID
			"kb_name":   doc["kb_name"],                            // 知识库名称
			"timestamp": doc["timestamp"],                          // 浏览时间（Unix 时间戳）
		}
		docList = append(docList, docMap)
	}

	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{"recent_docs": docList})
}

func (dc *DocumentController) GetRecentEditDocumentsHandler(c *gin.Context) {
	// 默认获取最近 10 条记录，可以通过查询参数调整
	startParam := c.DefaultQuery("start", "0")
	limitParam := c.DefaultQuery("limit", "10")
	start, err := strconv.Atoi(startParam)
	if err != nil || start <= 0 {
		start = 0
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10 // 设置默认值
	}

	// 获取用户 ID（从中间件中设置的上下文）
	userId, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未授权"})
		return
	}
	strUserId := strconv.FormatInt(userId.(int64), 10)

	// 从 Redis 中获取最近浏览记录
	recentDocs, err := dc.docDao.GetRecentDocumentWithScoresFromRedis(dao.Edit, strUserId, int64(start), int64(limit-1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，获取最近编辑文档失败"})
		return
	}
	log.Print(recentDocs)
	// 构造返回数据
	var docList []map[string]interface{}
	for _, doc := range recentDocs {
		// 构造文档信息
		docMap := map[string]interface{}{
			"doc_id":    fmt.Sprintf("%d", doc["doc_id"].(int64)),  // 文档 ID
			"doc_title": doc["doc_title"],                          // 文档标题
			"kb_id":     fmt.Sprintf("%d", doc["kb_id"].(float64)), // 知识库 ID
			"kb_name":   doc["kb_name"],                            // 知识库名称
			"timestamp": doc["timestamp"],                          // 浏览时间（Unix 时间戳）
		}
		docList = append(docList, docMap)
	}

	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{"recent_docs": docList})
}

// GetRecentDocumentsHandler 获取最近浏览记录
func (dc *DocumentController) GetRecentCommentDocumentsHandler(c *gin.Context) {
	// 默认获取最近 10 条记录，可以通过查询参数调整
	startParam := c.DefaultQuery("start", "0")
	limitParam := c.DefaultQuery("limit", "10")
	start, err := strconv.Atoi(startParam)
	if err != nil || start <= 0 {
		start = 0
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10 // 设置默认值
	}

	// 获取用户 ID（从中间件中设置的上下文）
	userId, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未授权"})
		return
	}
	strUserId := strconv.FormatInt(userId.(int64), 10)

	// 从 Redis 中获取最近浏览记录
	recentDocs, err := dc.docDao.GetRecentDocumentWithScoresFromRedis(dao.Comment, strUserId, int64(start), int64(limit-1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，获取最近评论文档失败"})
		return
	}
	log.Print(recentDocs)
	// 构造返回数据
	var docList []map[string]interface{}
	for _, doc := range recentDocs {
		// 构造文档信息
		docMap := map[string]interface{}{
			"doc_id":    fmt.Sprintf("%d", doc["doc_id"].(int64)),  // 文档 ID
			"doc_title": doc["doc_title"],                          // 文档标题
			"kb_id":     fmt.Sprintf("%d", doc["kb_id"].(float64)), // 知识库 ID
			"kb_name":   doc["kb_name"],                            // 知识库名称
			"timestamp": doc["timestamp"],                          // 浏览时间（Unix 时间戳）
		}
		docList = append(docList, docMap)
	}

	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{"recent_docs": docList})
}

func (dc *DocumentController) GetDocumenHashByIdHandler(c *gin.Context) {
	// 从请求参数获取文档ID
	docIdStr := c.Param("doc_id")
	docId, err := strconv.ParseInt(docIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "错误的文档ID"})
		return
	}

	hash, err := dc.docDao.GetDocumentContentHashByDocumentId(docId)
	if err != nil {
		hash = ""
	}
	c.JSON(http.StatusOK, gin.H{"doc_content_hash": hash})
}
