package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"yuqueppbackend/service-base/util"
)

func GetUserInfo(c *gin.Context) {
	var contextData struct {
		Id    int64  `json:"userid"`
		Email string `json:"email"`
	}
	// 从上下文中获取数据
	if userID, exists := c.Get("userid"); exists {
		contextData.Id = userID.(int64) // 获取并赋值
	}

	if email, exists := c.Get("email"); exists {
		contextData.Email = email.(string) // 获取并赋值
	}
	user, err := userDao.GetUserByID(contextData.Id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误请稍后再试"})
	}
	c.JSON(http.StatusOK, gin.H{"id": strconv.FormatInt(user.ID, 10), "email": user.Email, "nickname": user.Nickname})
}

func Logout(c *gin.Context) {
	// 从上下文中获取数据
	if email, exists := c.Get("email"); exists {
		err := util.DeleteJWT(email.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误请稍后再试"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"error": "系统错误请稍后再试"})
}
