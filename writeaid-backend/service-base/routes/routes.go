package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"yuqueppbackend/service-base/controllers"
	"yuqueppbackend/service-base/dao"
	"yuqueppbackend/service-base/db"
	"yuqueppbackend/service-base/util"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 设置 CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // 允许的跨域来源（可以是 *，但不推荐用于生产环境）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的 HTTP 方法
		AllowHeaders:     []string{"Content-Type", "Authorization"},           // 允许的请求头
		AllowCredentials: true,                                                // 是否允许携带凭证（如 Cookies）
	}))

	// 初始化 DAO 和 Controller
	kbDao := dao.NewKBDAO(db.GetDB(), util.GetElasticSearchClient())
	kbController := controllers.NewKnowledgeBaseController(kbDao)
	docDao := dao.NewDocDao(db.GetDB(), util.GetElasticSearchClient())
	docController := controllers.NewDocumentController(docDao)
	dcDao := dao.NewCommentDAO(db.GetDB())
	dcController := controllers.NewCommentController(dcDao)
	scDao := dao.NewSearchDao(util.GetElasticSearchClient())
	scController := controllers.NewSearchController(scDao)

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("register", controllers.Register)
		authGroup.POST("login", controllers.Login)
	}

	userGroup := r.Group("/api/user")
	userGroup.Use(util.AuthMiddleware())
	{
		userGroup.GET("getUserInfo", controllers.GetUserInfo)
		userGroup.POST("logout", controllers.Logout)
	}

	utilGroup := r.Group("/api/util")
	{
		utilGroup.GET("getCaptcha", controllers.GetCaptcha)

	}

	knowledgeGroup := r.Group("/api/knowledge")
	knowledgeGroup.Use(util.AuthMiddleware()) // 使用认知中间件
	{
		knowledgeGroup.POST("/createKnowledgeBase", kbController.CreateKnowledgeBase)
		knowledgeGroup.GET("/getKnowledgeBaseList", kbController.GetKnowledgeBaseList)
		knowledgeGroup.GET("/:kb_id", kbController.GetKnowledgeBaseDetail)
		knowledgeGroup.POST("/updateKnowledgeBase", kbController.UpdateKnowledgeBase)
		knowledgeGroup.POST("/deleteKnowledgeBase", kbController.DeleteKnowledgeBase)
	}

	documentGroup := r.Group("/api/document")
	documentGroup.Use(util.AuthMiddleware())
	{
		// 文档相关路由
		documentGroup.POST("/createDocument", docController.CreateDocumentHandler)
		documentGroup.GET("/getDocument/:doc_id", docController.GetDocumentByIDHandler)
		documentGroup.GET("/getDocumentListByKbId/:kb_id", docController.GetDocumentsByKnowledgeBaseIDHandler)
		documentGroup.PUT("/updateDocument/:doc_id", docController.UpdateDocumentHandler)
		documentGroup.DELETE("/deleteDocument/:doc_id", docController.DeleteDocumentByIDHandler)
		documentGroup.POST("/documents/:doc_id/view", docController.IncrementViewCountHandler)
		documentGroup.GET("/recentViewDocument", docController.GetRecentViewDocumentsHandler)
		documentGroup.GET("/recentEditDocument", docController.GetRecentEditDocumentsHandler)
		documentGroup.GET("/recentCommentDocument", docController.GetRecentCommentDocumentsHandler)
		documentGroup.GET("/documentContentHash/:doc_id", docController.GetDocumenHashByIdHandler)
	}
	documentCommentGroup := r.Group("/api/comment")
	documentCommentGroup.Use(util.AuthMiddleware())
	{
		documentCommentGroup.POST("/createDocumentComment", dcController.CreateDocumentComment)
		documentCommentGroup.POST("/replyDocumentComment", dcController.ReplyDocumentComment)
		documentCommentGroup.GET("/getDocumentRootComment/:doc_id", dcController.GetDocumentRootComment)
		documentCommentGroup.GET("/getChildrenComment/:root_id", dcController.GetDocumentChildComment)
	}
	searchGroup := r.Group("/api/search")
	searchGroup.Use(util.AuthMiddleware())
	{
		searchGroup.GET("/personalKnowledgeSearch/:search_text", scController.PersonalSearchKnowledgeBaseHandler)
		searchGroup.GET("/personalDocumentSearch/:search_text", scController.PersonalSearchDocumentTitleHandler)
	}

	return r
}
