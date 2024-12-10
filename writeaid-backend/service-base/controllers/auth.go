package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"yuqueppbackend/service-base/dao"
	"yuqueppbackend/service-base/models"
	"yuqueppbackend/service-base/util"
)

var userDao = dao.NewUserDAO()

// 用户注册
func Register(c *gin.Context) {
	var registerData struct {
		Email        string `json:"email" binding:"required"`
		Nickname     string `json:"nickname" binding:"required"`
		Password     string `json:"password" binding:"required"`
		CaptchaId    string `json:"captchaId" binding:"required"`
		CaptchaValue string `json:"captchaValue" binding:"required"`
	}
	if err := c.ShouldBindJSON(&registerData); err != nil {
		log.Println("数据绑定失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(registerData)
	// 确认验证码是否正确
	right, err := VerifyCaptcha(registerData.CaptchaId, registerData.CaptchaValue)
	if err != nil {
		log.Println("验证码模块错误")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误请稍后再试"})
		return
	}
	if !right {
		log.Println("验证码错误")
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}
	user := models.User{
		Email:    registerData.Email,
		Nickname: registerData.Nickname,
		Password: registerData.Password,
	}

	tmpUser, err := userDao.GetUserByEmail(user.Email)
	if err != nil {
		log.Println("系统错误")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误请稍后再试"})
		return
	}
	if tmpUser != nil {
		log.Println("用户已注册")
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户邮箱已注册，请直接登录"})
		return
	}
	log.Println(tmpUser)

	// 注册用户
	if err := userDao.CreateUser(user); err != nil {
		log.Println("系统错误")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户注册失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "用户注册成功"})
}

// 用户登录
func Login(c *gin.Context) {
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
	user := models.User{Email: loginData.Email, Password: loginData.Password}
	pass, err := userDao.CheckPassword(user)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误请稍后再试"})
	} else {
		if pass {
			tmp_user, err := userDao.GetUserByEmail(loginData.Email)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误请稍候再试"})
			}
			// 生成JWT
			if jwt, err := util.GenerateJWT(tmp_user.ID, tmp_user.Email); err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "系统错误请稍后再试"})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"access_token": jwt,
					"expires_in":   util.TokenExpireDuration,
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱或密码错误"})
		}
	}
}
