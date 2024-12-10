package util

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func userHaveLoginCache(c *gin.Context) bool {
	session := sessions.Default(c)
	useremail := session.Get("email")
	nickname := session.Get("nickname")
	if useremail == nil || nickname == nil {
		return false
	}
	return true
}
