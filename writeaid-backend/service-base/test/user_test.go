package test

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"yuqueppbackend/service-base/routes"
)

const prefix_url = "http://localhost:8080/api/user/"

func TestLogin(t *testing.T) {
	r := routes.SetupRouter()
	req, _ := http.NewRequest("GET", "/api/user/login", nil)
	// 使用 httptest.NewRecorder() 来记录响应
	w := httptest.NewRecorder()
	// 执行请求
	r.ServeHTTP(w, req)

	// 断言响应的状态码
	assert.Equal(t, http.StatusOK, w.Code)
	// 断言返回的数据
	expectedResponse := `{"status":"error"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestRegister(t *testing.T) {

}
