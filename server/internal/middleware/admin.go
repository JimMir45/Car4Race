package middleware

import (
	"net/http"

	"car4race/pkg/errcode"
	"car4race/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminAuth 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			response.ErrorWithCode(c, http.StatusForbidden, errcode.CodeAdminRequired, errcode.Message(errcode.CodeAdminRequired))
			c.Abort()
			return
		}
		c.Next()
	}
}
