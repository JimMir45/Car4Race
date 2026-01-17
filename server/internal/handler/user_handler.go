package handler

import (
	"net/http"
	"regexp"

	"car4race/internal/service"
	"car4race/pkg/errcode"
	"car4race/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// SendCode 发送验证码
func (h *UserHandler) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errcode.CodeInvalidParam, "参数错误")
		return
	}

	// 验证手机号格式
	if !isValidPhone(req.Phone) {
		response.ErrorWithCode(c, http.StatusBadRequest, errcode.CodeInvalidParam, "手机号格式不正确")
		return
	}

	if err := h.service.SendVerificationCode(req.Phone); err != nil {
		response.ErrorFromErr(c, err)
		return
	}

	response.Success(c, gin.H{"message": "验证码已发送"})
}

// Login 登录
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errcode.CodeInvalidParam, "参数错误")
		return
	}

	// 验证手机号格式
	if !isValidPhone(req.Phone) {
		response.ErrorWithCode(c, http.StatusBadRequest, errcode.CodeInvalidParam, "手机号格式不正确")
		return
	}

	// 验证验证码格式
	if len(req.Code) != 6 {
		response.ErrorWithCode(c, http.StatusBadRequest, errcode.CodeInvalidParam, "验证码格式不正确")
		return
	}

	token, user, err := h.service.Login(req.Phone, req.Code)
	if err != nil {
		response.ErrorFromErr(c, err)
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.ErrorWithCode(c, http.StatusUnauthorized, errcode.CodeUnauthorized, errcode.Message(errcode.CodeUnauthorized))
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		response.ErrorWithCode(c, http.StatusNotFound, errcode.CodeUserNotFound, errcode.Message(errcode.CodeUserNotFound))
		return
	}

	response.Success(c, user)
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.ErrorWithCode(c, http.StatusUnauthorized, errcode.CodeUnauthorized, errcode.Message(errcode.CodeUnauthorized))
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errcode.CodeInvalidParam, "参数错误")
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		response.ErrorWithCode(c, http.StatusNotFound, errcode.CodeUserNotFound, errcode.Message(errcode.CodeUserNotFound))
		return
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := h.service.UpdateUser(user); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(c, user)
}

// isValidPhone 验证手机号格式
func isValidPhone(phone string) bool {
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}
