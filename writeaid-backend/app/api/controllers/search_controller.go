package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"writeaid-backend/dao"
	"writeaid-backend/util"
)

type SearchController struct {
	dao *dao.SearchDao
}

func NewSearchController(dao *dao.SearchDao) *SearchController {
	return &SearchController{dao: dao}
}

func (sc *SearchController) PersonalSearchKnowledgeBaseHandler(c *gin.Context) {
	searchText := c.Param("search_text")
	var strKBIds []string
	strKBIds, highlightKbNames, err := sc.dao.GetKnowledgeBaseIDByName(searchText)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	log.Println(strKBIds)
	kbDao := dao.NewKBDAO(util.GetDB(), util.GetElasticSearchClient())

	var kbInfo []map[string]interface{}
	for idx, strKbId := range strKBIds {
		kbId, err := strconv.ParseInt(strKbId, 10, 64)
		if err != nil {
			continue
		}
		kb, err := kbDao.GetKnowledgeBaseById(kbId)
		if err != nil {
			c.JSON(500, gin.H{"error": "系统错误，请稍后再试"})
			return
		}
		kbInfo = append(kbInfo, map[string]interface{}{
			"kb_id":             kbId,
			"kb_name":           kb.Name,
			"highlight_kb_name": highlightKbNames[idx],
			"kb_description":    kb.Description,
			"kb_updated_at":     kb.UpdatedAt,
			"kb_created_at":     kb.CreatedAt,
		})
	}
	c.JSON(200, gin.H{"knowledgeBases": kbInfo})
}

func (sc *SearchController) PersonalSearchDocumentTitleHandler(c *gin.Context) {
	searchText := c.Param("search_text")
	var strDocIds []string
	strDocIds, hightDocTities, err := sc.dao.GetDocIDByTitleFuzzy(searchText)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	log.Println(strDocIds)

	docDao := dao.NewDocDao(util.GetDB(), util.GetElasticSearchClient())
	var docInfo []map[string]interface{}
	for idx, strDocId := range strDocIds {
		docId, err := strconv.ParseInt(strDocId, 10, 64)
		if err != nil {
			continue
		}
		doc, err := docDao.GetDocumentByID(docId)
		if err != nil {
			c.JSON(500, gin.H{"error": "系统错误，请稍后再试"})
			return
		}
		docInfo = append(docInfo, map[string]interface{}{
			"doc_id":          strDocId,
			"doc_title":       doc.Title,
			"highlight_title": hightDocTities[idx],
			"doc_created_at":  doc.CreatedAt,
			"doc_updated_at":  doc.UpdatedAt,
			"kb_name":         doc.KnowledgeBase.Name,
		})
	}
	c.JSON(200, gin.H{"documents": docInfo})
}

func (sc *SearchController) PersonalSearchDocumentContentHandler(c *gin.Context) {
	searchText := c.Param("search_text")
	strDocIds, err := sc.dao.GetDocIDByContent(searchText)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	log.Println(strDocIds)

	c.JSON(200, gin.H{"data": strDocIds})
}
