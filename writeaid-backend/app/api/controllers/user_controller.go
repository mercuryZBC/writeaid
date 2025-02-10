package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strconv"
	"writeaid-backend/app/api/middleware"
	pb "writeaid-backend/app/userService"
	"writeaid-backend/config"
	"writeaid-backend/dao"
)

type UserController struct {
	dao               *dao.UserDAO
	userServiceClient pb.UserServiceClient
}

func NewUserController(dao *dao.UserDAO) *UserController {
	// 连接 gRPC 服务器
	conn, err := grpc.NewClient(config.GetUserServiceAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	client := pb.NewUserServiceClient(conn)
	return &UserController{dao: dao, userServiceClient: client}
}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	var requestData struct {
		Id    int64  `json:"userid"`
		Email string `json:"email"`
	}
	// 从上下文中获取数据
	if userID, exists := c.Get("userid"); exists {
		requestData.Id = userID.(int64) // 获取并赋值
	}
	if email, exists := c.Get("email"); exists {
		requestData.Email = email.(string) // 获取并赋值
	}
	user, err := uc.dao.GetUserByID(requestData.Id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误请稍后再试"})
	}
	c.JSON(http.StatusOK, gin.H{"id": strconv.FormatInt(user.ID, 10), "email": user.Email, "nickname": user.Nickname})
}

func (uc *UserController) Logout(c *gin.Context) {
	// 从上下文中获取数据
	if email, exists := c.Get("email"); exists {
		err := middleware.DeleteJWT(email.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误请稍后再试"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"error": "系统错误请稍后再试"})
}

// 用户注册
func (uc *UserController) Register(c *gin.Context) {
	var requestData struct {
		Email        string `json:"email" binding:"required"`
		Nickname     string `json:"nickname" binding:"required"`
		Password     string `json:"password" binding:"required"`
		CaptchaId    string `json:"captchaId" binding:"required"`
		CaptchaValue string `json:"captchaValue" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Println("请求参数错误" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}
	log.Println(requestData)
	// 确认验证码是否正确
	right, err := VerifyCaptcha(requestData.CaptchaId, requestData.CaptchaValue)
	if err != nil {
		log.Println("验证码模块错误：" + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误请稍后再试"})
		return
	}
	if !right {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	registerResp, err := uc.userServiceClient.Register(context.Background(),
		&pb.RegisterReq{
			Email:    requestData.Email,
			Nickname: requestData.Nickname,
			Password: requestData.Password,
		})
	if err != nil || registerResp == nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败请稍后重试"})
		return
	}

	if registerResp.Status == pb.Status_REGISTER_OK {
		c.JSON(http.StatusOK, gin.H{"message": "用户注册成功"})
		return
	} else if registerResp.Status == pb.Status_REGISTER_USER_EXIST {
		log.Println("用户已注册")
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户邮箱已注册，请直接登录"})
		return
	} else {
		log.Println("用户注册服务错误")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户注册失败"})
		return
	}
}

// 用户登录
func (uc *UserController) Login(c *gin.Context) {
	var loginData struct {
		Email        string `json:"email" binding:"required"`
		Password     string `json:"password" binding:"required"`
		CaptchaId    string `json:"captchaId" binding:"required"`
		CaptchaValue string `json:"captchaValue" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证验证码
	if !store.Verify(loginData.CaptchaId, loginData.CaptchaValue, true) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	// 验证用户名和密码
	loginResp, err := uc.userServiceClient.Login(context.Background(),
		&pb.LoginReq{
			Email:    loginData.Email,
			Password: loginData.Password,
		})
	if err != nil || loginResp == nil {
		if loginResp != nil {
			if loginResp.Status == pb.Status_LOGIN_ERROR {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败，请稍后再试"})
				return
			}
			if loginResp.Status == pb.Status_LOGIN_INFO_ERROR {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "用户名或密码错误"})
				return
			}
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if loginResp.Status == pb.Status_LOGIN_OK {
		c.JSON(http.StatusOK, gin.H{
			"access_token": loginResp.AccessToken,
			"expires_in":   loginResp.ExpiresIn,
		})
	} else if loginResp.Status == pb.Status_LOGIN_INFO_ERROR {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或密码错误"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败请稍后重试"})
	}
}
