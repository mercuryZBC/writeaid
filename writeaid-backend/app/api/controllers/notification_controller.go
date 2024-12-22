package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "writeaid-backend/app/notificationService"
	"writeaid-backend/config"
)

type NotificationController struct {
	notificationServiceClient pb.NotificationServiceClient
}

func NewNotificationController() *NotificationController {
	// 连接 gRPC 服务器
	conn, err := grpc.NewClient(config.GetNotificationServiceAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect knowledgeBase service: %v", err)
	}
	client := pb.NewNotificationServiceClient(conn)
	return &NotificationController{notificationServiceClient: client}
}

func (ns *NotificationController) GetNotificationsCountHandler(c *gin.Context) {
	userId, exist := c.Get("userid")
	if exist == false {
		c.JSON(200, gin.H{"nt_count": 0})
	}
	resp, err := ns.notificationServiceClient.PullNotificationCount(context.Background(),
		&pb.PullNotificationCountReq{UserId: userId.(int64)})
	if err != nil {
		c.JSON(200, gin.H{"nt_count": 0})
		return
	}
	c.JSON(200, gin.H{"nt_count": resp.Count})
}

func (ns *NotificationController) GetNotificationHandler(c *gin.Context) {
	userId, exist := c.Get("userid")
	if exist == false {
		c.JSON(200, gin.H{"error": "通知加载失败，请稍后再试"})
		return
	}
	resp, err := ns.notificationServiceClient.Pull(context.Background(),
		&pb.PullReq{UserId: userId.(int64)})
	if err != nil {
		c.JSON(200, gin.H{"error": "通知加载失败，请稍后再试"})
		return
	}

	var pullContents []map[string]interface{}
	for _, notification := range resp.PullContent {
		pullContents = append(pullContents, map[string]interface{}{
			"nt_type":            notification.Type,
			"nt_message_content": notification.MessageContent,
			"nt_time":            notification.Time,
			"nt_link":            notification.Link,
			"nt_link_describe":   notification.LinkDescribe,
			"nt_is_check":        notification.IsCheck,
			"nt_avatar_link":     notification.AvatarLink,
		})
	}
	c.JSON(200, gin.H{"nt_content": pullContents})
}
