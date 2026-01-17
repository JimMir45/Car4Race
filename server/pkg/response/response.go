package response

import (
	"car4race/pkg/errcode"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpCode int, message string) {
	c.JSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
	})
}

// ErrorWithCode 带错误码的响应
func ErrorWithCode(c *gin.Context, httpCode int, code int, message string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorFromErr 根据 errcode.Error 自动响应
func ErrorFromErr(c *gin.Context, err error) {
	if e, ok := err.(*errcode.Error); ok {
		httpCode := getHTTPCode(e.Code)
		c.JSON(httpCode, Response{
			Code:    e.Code,
			Message: e.Message,
		})
		return
	}
	// 非 errcode.Error 类型，使用通用错误
	c.JSON(400, Response{
		Code:    400,
		Message: err.Error(),
	})
}

// getHTTPCode 根据业务错误码获取 HTTP 状态码
func getHTTPCode(code int) int {
	switch {
	case code >= 40001 && code < 40100:
		return 400 // 参数错误
	case code == errcode.CodeUnauthorized:
		return 401 // 未授权
	case code >= 40301 && code < 40400:
		return 403 // 禁止访问
	case code >= 40401 && code < 40500:
		return 404 // 资源不存在
	case code == errcode.CodeRateLimitExceed:
		return 429 // 请求过多
	default:
		return 400
	}
}
