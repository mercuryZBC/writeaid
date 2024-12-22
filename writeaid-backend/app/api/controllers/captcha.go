package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	id, b64s, err := GenerateCaptcha()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"captchaId": id, "captcha": b64s})
}

func GenerateCaptcha() (string, string, error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80) // 图形验证码配置
	captcha := base64Captcha.NewCaptcha(driver, store)

	id, b64s, _, err := captcha.Generate()
	if err != nil {
		return "", "", err
	}
	return id, b64s, nil
}

func VerifyCaptcha(captchaId string, value string) (bool, error) {
	if store.Verify(captchaId, value, true) {
		return true, nil
	} else {
		return false, nil
	}
}
