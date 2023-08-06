package jwtoken

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/response"
)

// JWT 签名密钥，定义JWT中存放usrname和userID

var mySigningKey = []byte("mini_tiktok")

const (
	InvalidToken = "invalid token"
)

type JWTClaims struct {
	Username string `json:"user_name"`
	UserID   int64  `json:"user_id"`
	jwt.RegisteredClaims
}

// ParseToken 解析token
func ParseToken(tokenString string) (*JWTClaims, error) {
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New(InvalidToken)
}

// VerifyToken 校验token
func VerifyToken(token string) (int64, error) {
	if token == "" {
		return int64(0), nil
	}

	claims, err := ParseToken(token)
	if err != nil {
		return int64(0), err
	}
	return claims.UserID, nil

}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.PostForm("token")
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		userID, err := VerifyToken(tokenString)
		if err != nil || userID == int64(0) {
			log.Error("token error...")
			response.Fail(c, "auth error", nil)
			c.Abort()
		}

		c.Set("UserID", userID)
		c.Next()
	}
}

// JWTWithoutAuthMiddleware feed 针对登录和未登录采取不同的措施
func JWTWithoutAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.PostForm("token")
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		userID, err := VerifyToken(tokenString)
		if err != nil {
			log.Error("token error...")
			response.Fail(c, "auth error", nil)
			c.Abort()
		}
		// token 为空，传入0
		c.Set("UserID", userID)
		c.Next()
	}
}
