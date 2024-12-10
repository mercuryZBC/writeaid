package util

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	jwtSecretKey = []byte("your_jwt_secret")
)

const TokenExpireDuration = time.Hour * 24 * 7

// 用于 JWT 的 Claims 结构体
type Claims struct {
	ID    int64  `json:"userid"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// 生成 JWT 的函数
func GenerateJWT(id int64, email string) (string, error) {
	claims := &Claims{
		ID:    id,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 设置过期时间
			IssuedAt:  time.Now().Unix(),                          // 签发时间
			Issuer:    "your_app_name",                            // 签发者
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 获取签名后的 token 字符串
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	// 保存 token 到 Redis（过期时间与 JWT 相同）
	err = GetRedisClient().Set(context.Background(), tokenString, email, TokenExpireDuration).Err()
	err = GetRedisClient().Set(context.Background(), email, tokenString, TokenExpireDuration).Err()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DeleteJWT(email string) error {
	jwtToken := GetRedisClient().Get(context.Background(), email).Val()
	err := GetRedisClient().Del(context.Background(), jwtToken).Err()
	if err != nil {
		return err
	}
	err = GetRedisClient().Del(context.Background(), email).Err()
	if err != nil {
		return err
	}
	return nil
}

// 验证 Token 的中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header invalid"})
			c.Abort()
			return
		}

		// 检查 token 是否存在于 Redis 中
		email, err := GetRedisClient().Get(context.Background(), tokenString).Result()
		if err == redis.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token is not valid or expired"})
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Redis error"})
			c.Abort()
			return
		}

		// 解析 JWT
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// 只验证签名部分，不验证 JWT 的其他部分
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("userid", claims.ID)
		c.Set("email", email)
		log.Printf("User ID: %d, Email: %s", claims.ID, claims.Email)
		// 如果 JWT 有效，可以继续处理请求
		c.Next()
	}
}
