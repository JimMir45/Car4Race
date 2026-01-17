package middleware

import (
	"net/http"
	"strings"

	"car4race/pkg/errcode"
	"car4race/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth JWT 认证中间件
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorWithCode(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "未提供认证信息")
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorWithCode(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "认证格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			response.ErrorWithCode(c, http.StatusUnauthorized, errcode.CodeUnauthorized, errcode.Message(errcode.CodeUnauthorized))
			c.Abort()
			return
		}

		// 提取 claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.ErrorWithCode(c, http.StatusUnauthorized, errcode.CodeUnauthorized, "认证信息无效")
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(userID))
		}
		if phone, ok := claims["phone"].(string); ok {
			c.Set("phone", phone)
		}
		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
		}

		c.Next()
	}
}

// OptionalJWTAuth 可选 JWT 认证中间件
// 不要求必须登录，但如果提供了有效 token 则提取用户信息
func OptionalJWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]

		// 解析 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.Next()
			return
		}

		// 提取 claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Next()
			return
		}

		// 设置用户信息到上下文
		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", uint(userID))
		}
		if phone, ok := claims["phone"].(string); ok {
			c.Set("phone", phone)
		}
		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
		}

		c.Next()
	}
}
