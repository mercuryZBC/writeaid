package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"writeaid-backend/app/api/controllers"
	"writeaid-backend/app/api/middleware"
	"writeaid-backend/dao"
	"writeaid-backend/util"
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
	userDao := dao.NewUserDao(util.GetDB(), util.GetElasticSearchClient())
	userController := controllers.NewUserController(userDao)
	kbDao := dao.NewKBDAO(util.GetDB(), util.GetElasticSearchClient())
	kbController := controllers.NewKnowledgeBaseController(kbDao)
	docDao := dao.NewDocDao(util.GetDB(), util.GetElasticSearchClient())
	docController := controllers.NewDocumentController(docDao)
	dcDao := dao.NewCommentDAO(util.GetDB())
	dcController := controllers.NewCommentController(dcDao)
	scDao := dao.NewSearchDao(util.GetElasticSearchClient())
	scController := controllers.NewSearchController(scDao)
	notificationController := controllers.NewNotificationController()

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("register", userController.Register)
		authGroup.POST("login", userController.Login)
	}

	userGroup := r.Group("/api/user")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.GET("getUserInfo", userController.GetUserInfo)
		userGroup.POST("logout", userController.Logout)
	}

	utilGroup := r.Group("/api/util")
	{
		utilGroup.GET("getCaptcha", controllers.GetCaptcha)

	}

	knowledgeGroup := r.Group("/api/knowledge")
	knowledgeGroup.Use(middleware.AuthMiddleware()) // 使用认知中间件
	{
		knowledgeGroup.POST("/createKnowledgeBase", kbController.CreateKnowledgeBase)
		knowledgeGroup.GET("/getKnowledgeBaseList", kbController.GetKnowledgeBaseList)
		knowledgeGroup.GET("/:kb_id", kbController.GetKnowledgeBaseDetail)
		knowledgeGroup.POST("/updateKnowledgeBase", kbController.UpdateKnowledgeBase)
		knowledgeGroup.POST("/deleteKnowledgeBase", kbController.DeleteKnowledgeBase)
	}

	documentGroup := r.Group("/api/document")
	documentGroup.Use(middleware.AuthMiddleware())
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
	documentCommentGroup.Use(middleware.AuthMiddleware())
	{
		documentCommentGroup.POST("/createDocumentComment", dcController.CreateDocumentComment)
		documentCommentGroup.POST("/replyDocumentComment", dcController.ReplyDocumentComment)
		documentCommentGroup.GET("/getDocumentRootComment/:doc_id", dcController.GetDocumentRootComment)
		documentCommentGroup.GET("/getChildrenComment/:root_id", dcController.GetDocumentChildComment)
	}
	searchGroup := r.Group("/api/search")
	searchGroup.Use(middleware.AuthMiddleware())
	{
		searchGroup.GET("/personalKnowledgeSearch/:search_text", scController.PersonalSearchKnowledgeBaseHandler)
		searchGroup.GET("/personalDocumentSearch/:search_text", scController.PersonalSearchDocumentTitleHandler)
	}
	notificationGroup := r.Group("/api/notification")
	notificationGroup.Use(middleware.AuthMiddleware())
	{
		notificationGroup.GET("/total", notificationController.GetNotificationHandler)
		notificationGroup.GET("/count", notificationController.GetNotificationsCountHandler)
	}
	return r
}
