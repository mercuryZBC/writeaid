package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"time"
	pb "writeaid-backend/app/notificationService"
	"writeaid-backend/config"
	"writeaid-backend/dao"
	"writeaid-backend/models"
	"writeaid-backend/util"
)

type notificationServiceServer struct {
	pb.UnimplementedNotificationServiceServer
}

var notificationDao *dao.NotificationHistoryDAO
var kbDao *dao.KBDAO
var docDao *dao.DocDao
var userDao *dao.UserDAO

func init() {
	notificationDao = dao.NewNotificationHistoryDAO(util.GetDB(), util.GetRedisClient())
	kbDao = dao.NewKBDAO(util.GetDB(), util.GetElasticSearchClient())
	docDao = dao.NewDocDao(util.GetDB(), util.GetElasticSearchClient())
	userDao = dao.NewUserDao(util.GetDB(), util.GetElasticSearchClient())
}

func logFormat(funcName, result string, v ...interface{}) string {
	return fmt.Sprintf(funcName, ":", "用户信息 ", v, result)
}

func (ps *notificationServiceServer) PushKnowledgeBaseOperate(ctx context.Context, rep *pb.PushKnowledgeBaseOperateReq) (*pb.PushKnowledgeBaseOperateResp, error) {
	content := ""
	kb, err := kbDao.GetKnowledgeBaseById(rep.KnowledgeBaseId)
	if err != nil {
		log.Println(err)
		return &pb.PushKnowledgeBaseOperateResp{
			Status: pb.Status_PUSH_ERROR,
		}, fmt.Errorf("推送系统错误")
	}
	if kb == nil {
		content = "知识库不存在或已被删除"
	} else {
		content = kb.Name
		if rep.Operate == pb.KnowLedgeBaseOperate_KB_CREATE {
			content = "知识库" + content + "创建成功"
		} else if rep.Operate == pb.KnowLedgeBaseOperate_KB_VIEW {
			content = "知识库" + content + "已被查看"
		} else if rep.Operate == pb.KnowLedgeBaseOperate_KB_EDIT {
			content = "知识库" + content + "已被编辑"
		} else if rep.Operate == pb.KnowLedgeBaseOperate_KB_DELETE {
			content = "知识库" + content + "已被删除"
		}
	}

	notification := models.NotificationHistory{
		ReceiveUserID: rep.UserId,
		NotifyType:    models.KnowledgeBaseOper,
		Content:       content,
		Link:          "",
		LinkDescribe:  "",
		AvatarLink:    "",
		CreatedAt:     time.Unix(rep.Time, 0),
		IsCheck:       false,
	}
	err = notificationDao.CreateNotification(&notification)
	if err != nil {
		log.Println(err)
		return &pb.PushKnowledgeBaseOperateResp{
			Status: pb.Status_PUSH_ERROR,
		}, fmt.Errorf("推送系统错误")
	}
	log.Println("知识库相关通知push成功")
	return &pb.PushKnowledgeBaseOperateResp{Status: pb.Status_PUSH_OK}, nil
}

func (ps *notificationServiceServer) PushDocumentOperate(ctx context.Context, rep *pb.PushDocumentOperateReq) (*pb.PushDocumentOperateResp, error) {
	content := ""
	doc, err := docDao.GetDocumentByID(rep.DocumentId)
	if err != nil {
		log.Println(err)
		return &pb.PushDocumentOperateResp{
			Status: pb.Status_PUSH_ERROR,
		}, fmt.Errorf("推送系统错误")
	}
	if doc == nil {
		content = "文档不存在或已被删除"
	} else {
		content = doc.Title
		if rep.Operate == pb.DocumentOperate_DOC_CREATE {
			content = "文档" + content + "创建成功"
		} else if rep.Operate == pb.DocumentOperate_DOC_VIEW {
			content = "文档" + content + "已被查看"
		} else if rep.Operate == pb.DocumentOperate_DOC_EDIT {
			content = "文档" + content + "已被编辑"
		} else if rep.Operate == pb.DocumentOperate_DOC_DELETE {
			content = "文档" + content + "已被删除"
		}
	}

	notification := models.NotificationHistory{
		ReceiveUserID: rep.UserId,
		NotifyType:    models.DocumentOper,
		Content:       content,
		Link:          "",
		LinkDescribe:  "",
		AvatarLink:    "",
		CreatedAt:     time.Unix(rep.Time, 0),
		IsCheck:       false,
	}
	err = notificationDao.CreateNotification(&notification)
	if err != nil {
		log.Println(err)
		return &pb.PushDocumentOperateResp{
			Status: pb.Status_PUSH_ERROR,
		}, fmt.Errorf("推送系统错误")
	}
	log.Println("知识库相关通知push成功")
	return &pb.PushDocumentOperateResp{Status: pb.Status_PUSH_OK}, nil
}

func (ps *notificationServiceServer) PushComment(ctx context.Context, rep *pb.PushCommentReq) (*pb.PushCommentResp, error) {
	content := ""
	user, err := userDao.GetUserByID(rep.CommentUserId)
	if err != nil {
		log.Println(err)
		return &pb.PushCommentResp{
			Status: pb.Status_PUSH_ERROR,
		}, fmt.Errorf("推送系统错误")
	}
	if user == nil {
		content = "已注销用户回复了你的评论："
	} else {
		content = user.Nickname + "回复了你的评论："
	}

	notification := models.NotificationHistory{
		ReceiveUserID: rep.UserId,
		NotifyType:    models.Comment,
		Content:       content,
		Link:          strconv.FormatInt(rep.DocumentId, 10),
		LinkDescribe:  "点击查看",
		AvatarLink:    "",
		CreatedAt:     time.Unix(rep.CommentTime, 0),
		IsCheck:       false,
	}
	err = notificationDao.CreateNotification(&notification)
	if err != nil {
		log.Println(err)
		return &pb.PushCommentResp{
			Status: pb.Status_PUSH_ERROR,
		}, fmt.Errorf("推送系统错误")
	}
	log.Println("知识库相关通知push成功")
	return &pb.PushCommentResp{Status: pb.Status_PUSH_OK}, nil
}

func (ps *notificationServiceServer) Pull(ctx context.Context, rep *pb.PullReq) (*pb.PullResp, error) {
	notifications, err := notificationDao.GetNotificationsByUserID(rep.UserId, max(int(100)), 0)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("推送系统发送错误")
	}
	pullContents := []*pb.PullContent{}
	for _, notification := range notifications {
		pullContent := &pb.PullContent{
			Type:           pb.NotifyType(notification.NotifyType),
			MessageContent: notification.Content,
			Time:           notification.CreatedAt.Unix(),
			Link:           notification.Link,
			LinkDescribe:   notification.LinkDescribe,
			IsCheck:        notification.IsCheck,
			AvatarLink:     notification.AvatarLink,
		}
		pullContents = append(pullContents, pullContent)
	}
	return &pb.PullResp{PullContent: pullContents}, nil
}

func (ps *notificationServiceServer) PullNotificationCount(ctx context.Context, req *pb.PullNotificationCountReq) (*pb.PullNotificationCountResp, error) {
	if req == nil || req.UserId == 0 {
		log.Println("拉取通知服务调用失败")
		return nil, fmt.Errorf("拉取通知服务调用失败")
	}
	count, err := notificationDao.GetUnreadNotificationCount(req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.PullNotificationCountResp{Count: count}, nil
}

func main() {
	logFile, err := os.OpenFile(config.GetNotificationServiceLogfilePath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(logFile)
	listener, err := net.Listen("tcp", config.GetNotificationServicePort())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, &notificationServiceServer{})

	fmt.Println("notificationService server is running on port " + config.GetNotificationServicePort())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to notificationService serve: %v", err)
	}
}
