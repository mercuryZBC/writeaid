package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"writeaid-backend/app/api/middleware"
	pb "writeaid-backend/app/userService"
	"writeaid-backend/config"
	"writeaid-backend/dao"
	"writeaid-backend/models"
	"writeaid-backend/util"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
}

var userDao *dao.UserDAO

func init() {
	userDao = dao.NewUserDao(util.GetDB(), util.GetElasticSearchClient())
}

func logFormat(funcName, result string, v ...interface{}) string {
	return fmt.Sprint(funcName, ":", "用户信息 ", v, result)
}

// Register方法
func (s *userServiceServer) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterResp, error) {
	user := models.User{
		Email:    req.Email,
		Nickname: req.Nickname,
		Password: req.Password,
	}

	tmpUser, err := userDao.GetUserByEmail(user.Email)
	if err != nil {
		log.Println(logFormat("用户注册", err.Error(), req))
		return &pb.RegisterResp{Status: pb.Status_REGISTER_ERROR}, err
	}
	if tmpUser != nil {
		log.Println(logFormat("用户注册", "用户已存在", req))
		return &pb.RegisterResp{Status: pb.Status_REGISTER_USER_EXIST}, nil
	}

	// 注册用户
	if err := userDao.CreateUser(user); err != nil {
		log.Println(logFormat("用户注册", err.Error(), req))
		return &pb.RegisterResp{Status: pb.Status_REGISTER_ERROR}, err
	}
	return &pb.RegisterResp{Status: pb.Status_REGISTER_OK}, nil
}

func (s *userServiceServer) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {

	// 验证用户名和密码
	user := models.User{Email: req.Email, Password: req.Password}
	pass, err := userDao.CheckPassword(user)
	if err != nil {
		log.Println(logFormat("用户登录", err.Error(), req))
		return &pb.LoginResp{Status: pb.Status_LOGIN_ERROR}, err
	}
	// 验证失败
	if !pass {
		log.Println(logFormat("用户登录", "用户信息错误", req))
		return &pb.LoginResp{Status: pb.Status_LOGIN_INFO_ERROR}, nil
	}

	//验证成功
	tmp_user, err := userDao.GetUserByEmail(user.Email)
	if err != nil {
		log.Println(logFormat("用户登录", "登录成功", req))
		return &pb.LoginResp{Status: pb.Status_LOGIN_ERROR}, err
	}
	// 生成JWT
	jwt, err := middleware.GenerateJWT(tmp_user.ID, tmp_user.Email)
	if err != nil {
		log.Println(logFormat("用户登录", "JWT生成失败"+err.Error(), req))
		return &pb.LoginResp{Status: pb.Status_LOGIN_ERROR}, err
	}
	log.Println(logFormat("用户登录", "成功", req))
	return &pb.LoginResp{
		Status:      pb.Status_LOGIN_OK,
		AccessToken: jwt,
		ExpiresIn:   int64(middleware.TokenExpireDuration),
	}, nil
}

func main() {
	logFile, err := os.OpenFile(config.GetUserServiceLogfilePath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetFlags(log.LstdFlags)
	log.SetOutput(logFile)
	listener, err := net.Listen("tcp", config.GetUserServicePort())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &userServiceServer{})
	fmt.Println("userService server is running on port 8001")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to userService serve: %v", err)
	}
}
