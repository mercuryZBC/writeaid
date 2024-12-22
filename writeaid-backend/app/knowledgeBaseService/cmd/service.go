package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"time"
	pb "writeaid-backend/app/knowledgeBaseService"
	notifyPb "writeaid-backend/app/notificationService"
	"writeaid-backend/config"
	"writeaid-backend/dao"
	"writeaid-backend/models"
	"writeaid-backend/util"
)

type knowledgeBaseServiceServer struct {
	pb.UnimplementedKnowledgeServiceServer
}

var kbDao *dao.KBDAO
var notifyClient notifyPb.NotificationServiceClient

func init() {
	kbDao = dao.NewKBDAO(util.GetDB(), util.GetElasticSearchClient())
	// 连接 gRPC 服务器
	conn, err := grpc.NewClient(config.GetNotificationServiceAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect notification service: %v", err)
	}
	notifyClient = notifyPb.NewNotificationServiceClient(conn)
}

func (ks *knowledgeBaseServiceServer) CreateKnowledgeBase(ctx context.Context, rep *pb.CreateKnowledgeBaseReq) (*pb.CreateKnowledgeBaseResp, error) {
	knowledgeBase := models.KnowledgeBase{
		Name:        rep.KnowledgeBase.KbName,
		Description: rep.KnowledgeBase.KbDescription,
		IsPublic:    rep.KnowledgeBase.KbIsPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		OwnerID:     rep.KnowledgeBase.KbOwnerId,
	}
	// 使用 KBDAO 来创建知识库
	if err := kbDao.CreateKB(&knowledgeBase); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("系统错误请稍后再试")
	}
	err := kbDao.InsertKBToEs(knowledgeBase)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("系统错误请稍后再试")
	}
	kbRes := pb.KnowledgeBase{
		KbName:        knowledgeBase.Name,
		KbId:          knowledgeBase.ID,
		KbOwnerId:     knowledgeBase.OwnerID,
		KbDescription: knowledgeBase.Description,
		KbIsPublic:    knowledgeBase.IsPublic,
		KbCreateAt:    knowledgeBase.CreatedAt.Unix(),
		KbUpdateAt:    knowledgeBase.UpdatedAt.Unix(),
	}

	_, err = notifyClient.PushKnowledgeBaseOperate(context.Background(), &notifyPb.PushKnowledgeBaseOperateReq{
		UserId:            knowledgeBase.OwnerID,
		KnowledgeBaseId:   knowledgeBase.ID,
		KnowledgeBaseName: knowledgeBase.Name,
		Time:              time.Now().Unix(),
		Operate:           notifyPb.KnowLedgeBaseOperate_KB_CREATE,
	})
	if err != nil {
		log.Println(err)
	}

	return &pb.CreateKnowledgeBaseResp{
		Status:        pb.Status_CREATE_OK,
		KnowledgeBase: &kbRes,
	}, nil
}
func (ks *knowledgeBaseServiceServer) GetKnowledgeBaseList(ctx context.Context, rep *pb.GetKnowledgeBaseListReq) (*pb.GetKnowledgeBaseListResp, error) {
	if kbList, err := kbDao.GetKBListByOwnerId(rep.UserId); err == nil {
		var kbListData []*pb.KnowledgeBase
		for _, kb := range kbList {
			kbResp := pb.KnowledgeBase{
				KbId:          kb.ID,
				KbName:        kb.Name,
				KbDescription: kb.Description,
				KbIsPublic:    kb.IsPublic,
				KbOwnerId:     kb.OwnerID,
				KbCreateAt:    kb.UpdatedAt.Unix(),
				KbUpdateAt:    kb.UpdatedAt.Unix(),
			}
			kbListData = append(kbListData, &kbResp)
		}
		return &pb.GetKnowledgeBaseListResp{KnowledgeBases: kbListData}, nil
	} else {
		log.Println(err)
		return nil, fmt.Errorf("系统错误请稍后再试")
	}
}
func (ks *knowledgeBaseServiceServer) GetKnowledgeBaseDetail(ctx context.Context, rep *pb.GetKnowledgeBaseDetailReq) (*pb.GetKnowledgeBaseDetailResp, error) {
	// 使用 KBDAO 查找知识库
	kb, err := kbDao.FindKB(rep.UserId, rep.KbId)
	if err != nil {
		return nil, fmt.Errorf("系统错误请稍后再试")
	}
	if kb == nil {
		return nil, fmt.Errorf("知识库不存在")
	}
	knowledgeBase := pb.KnowledgeBase{
		KbId:          kb.ID,
		KbName:        kb.Name,
		KbDescription: kb.Description,
		KbIsPublic:    kb.IsPublic,
		KbOwnerId:     kb.OwnerID,
		KbCreateAt:    kb.UpdatedAt.Unix(),
		KbUpdateAt:    kb.UpdatedAt.Unix(),
	}
	return &pb.GetKnowledgeBaseDetailResp{
		Status:        pb.Status_GET_DETAIL_OK,
		KnowledgeBase: &knowledgeBase,
	}, nil
}
func (ks *knowledgeBaseServiceServer) DeleteKnowledgeBase(ctx context.Context, rep *pb.DeleteKnowledgeBaseReq) (*pb.DeleteKnowledgeBaseResp, error) {
	knowledgeBase, err := kbDao.FindKB(rep.UserId, rep.KbId)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("系统错误请稍后再试")
	}
	// 使用 KBDAO 删除知识库
	if err := kbDao.DeleteKB(*knowledgeBase); err != nil {
		return nil, fmt.Errorf("系统错误请稍后再试")
	}
	if err := kbDao.DeleteKBFromES(knowledgeBase.ID); err != nil {
		return nil, fmt.Errorf("系统错误请稍后再试")
	}

	_, err = notifyClient.PushKnowledgeBaseOperate(context.Background(), &notifyPb.PushKnowledgeBaseOperateReq{
		UserId:            knowledgeBase.OwnerID,
		KnowledgeBaseId:   knowledgeBase.ID,
		KnowledgeBaseName: knowledgeBase.Name,
		Time:              time.Now().Unix(),
		Operate:           notifyPb.KnowLedgeBaseOperate_KB_DELETE,
	})
	if err != nil {
		log.Println(err)
	}

	return &pb.DeleteKnowledgeBaseResp{Status: pb.Status_DELETE_OK}, nil
}
func (ks *knowledgeBaseServiceServer) UpdateKnowledgeBase(ctx context.Context, rep *pb.UpdateKnowledgeBaseReq) (*pb.UpdateKnowledgeBaseResp, error) {
	kb := models.KnowledgeBase{
		ID:          rep.KnowledgeBase.KbId,
		Name:        rep.KnowledgeBase.KbName,
		Description: rep.KnowledgeBase.KbDescription,
		IsPublic:    rep.KnowledgeBase.KbIsPublic,
		CreatedAt:   time.Unix(rep.KnowledgeBase.KbCreateAt, 0),
		UpdatedAt:   time.Unix(rep.KnowledgeBase.KbUpdateAt, 0),
		OwnerID:     rep.KnowledgeBase.KbOwnerId,
	}
	resultKb, err := kbDao.UpdateKB(kb.OwnerID, kb.ID, kb)
	if err != nil {
		return nil, fmt.Errorf("系统错误请稍后再试")
	}
	respKb := pb.KnowledgeBase{
		KbId:          resultKb.ID,
		KbName:        resultKb.Name,
		KbDescription: resultKb.Description,
		KbIsPublic:    resultKb.IsPublic,
		KbCreateAt:    resultKb.UpdatedAt.Unix(),
		KbUpdateAt:    resultKb.UpdatedAt.Unix(),
		KbOwnerId:     resultKb.OwnerID,
	}

	_, err = notifyClient.PushKnowledgeBaseOperate(context.Background(), &notifyPb.PushKnowledgeBaseOperateReq{
		UserId:            kb.OwnerID,
		KnowledgeBaseId:   kb.ID,
		KnowledgeBaseName: kb.Name,
		Time:              time.Now().Unix(),
		Operate:           notifyPb.KnowLedgeBaseOperate_KB_CREATE,
	})
	if err != nil {
		log.Println(err)
	}

	return &pb.UpdateKnowledgeBaseResp{
		Status:        pb.Status_UPDATE_OK,
		KnowledgeBase: &respKb,
	}, nil
}

func main() {
	logFile, err := os.OpenFile(config.GetKnowledgeBaseServiceLogfilePath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetFlags(log.LstdFlags)
	log.SetOutput(logFile)
	listener, err := net.Listen("tcp", config.GetKnowledgeBaseServicePort())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterKnowledgeServiceServer(grpcServer, &knowledgeBaseServiceServer{})

	fmt.Println("knowledgeBaseService server is running on port 8003")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to knowledgeBaseService serve: %v", err)
	}
}
