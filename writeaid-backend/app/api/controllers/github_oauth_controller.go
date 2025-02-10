package controllers

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"net/http"
)

var githubOAuthConfig = &oauth2.Config{
	ClientID:     "your-client-id",
	ClientSecret: "your-client-secret",
	Scopes:       []string{"user:email"},
	Endpoint:     github.Endpoint,
	RedirectURL:  "http://localhost:8080/callback",
}

func githubLoginHandler(c *gin.Context) {
	url := githubOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// 回调处理器
func githubCallbackHandler(c *gin.Context) {
	// 验证state 防止CSRF攻击
	code := c.Query("code")
	// 使用授权码换取访问令牌
	token, err := githubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 获取用户信息
	client := githubOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer resp.Body.Close()

	// 解析用户信息
	var user struct {
		Login string `json:"login"`
		Email string `json:"email"`
		ID    int    `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 显示用户信息
	c.JSON(http.StatusOK, gin.H{"user": user})
}
