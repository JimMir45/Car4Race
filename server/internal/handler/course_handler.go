package handler

import (
	"net/http"
	"strconv"

	"car4race/internal/service"
	"car4race/pkg/errcode"
	"car4race/pkg/response"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service *service.CourseService
}

func NewCourseHandler(service *service.CourseService) *CourseHandler {
	return &CourseHandler{service: service}
}

// GetCourses 获取课程列表
func (h *CourseHandler) GetCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort", "newest") // newest | price_asc | price_desc | sales

	courses, total, err := h.service.GetCourses(page, pageSize, sortBy)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取课程失败")
		return
	}

	response.Success(c, gin.H{
		"list":      courses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetCourse 获取课程详情
func (h *CourseHandler) GetCourse(c *gin.Context) {
	slug := c.Param("slug")

	course, err := h.service.GetCourseBySlug(slug)
	if err != nil {
		response.ErrorWithCode(c, http.StatusNotFound, errcode.CodeCourseNotFound, errcode.Message(errcode.CodeCourseNotFound))
		return
	}

	// 检查用户是否已购买（支持可选登录）
	var userID uint
	if val, exists := c.Get("user_id"); exists {
		if v, ok := val.(uint); ok {
			userID = v
		}
	}

	purchased := false
	if userID > 0 {
		purchased, _ = h.service.CheckUserPurchased(userID, course.ID)
	}

	response.Success(c, gin.H{
		"course":    course,
		"purchased": purchased,
	})
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	CourseID uint `json:"course_id" binding:"required"`
}

// CreateOrder 创建订单
func (h *CourseHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.ErrorWithCode(c, http.StatusUnauthorized, errcode.CodeUnauthorized, errcode.Message(errcode.CodeUnauthorized))
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errcode.CodeInvalidParam, "参数错误")
		return
	}

	order, err := h.service.CreateOrder(userID, req.CourseID)
	if err != nil {
		response.ErrorFromErr(c, err)
		return
	}

	response.Success(c, order)
}

// GetOrders 获取用户订单列表
func (h *CourseHandler) GetOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.ErrorWithCode(c, http.StatusUnauthorized, errcode.CodeUnauthorized, errcode.Message(errcode.CodeUnauthorized))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.service.GetUserOrders(userID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取订单失败")
		return
	}

	response.Success(c, gin.H{
		"list":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// RedeemCodeRequest 兑换邀请码请求
type RedeemCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

// RedeemCode 使用邀请码兑换课程
func (h *CourseHandler) RedeemCode(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.ErrorWithCode(c, http.StatusUnauthorized, errcode.CodeUnauthorized, errcode.Message(errcode.CodeUnauthorized))
		return
	}

	var req RedeemCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errcode.CodeInvalidParam, "参数错误")
		return
	}

	order, err := h.service.RedeemInviteCode(userID, req.Code)
	if err != nil {
		response.ErrorFromErr(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "兑换成功",
		"order":   order,
	})
}

// CreateDownloadRequest 创建下载请求
type CreateDownloadRequest struct {
	CourseID uint `json:"course_id" binding:"required"`
}

// CreateDownload 创建下载令牌
func (h *CourseHandler) CreateDownload(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.ErrorWithCode(c, http.StatusUnauthorized, errcode.CodeUnauthorized, errcode.Message(errcode.CodeUnauthorized))
		return
	}

	var req CreateDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errcode.CodeInvalidParam, "参数错误")
		return
	}

	token, err := h.service.CreateDownloadToken(userID, req.CourseID)
	if err != nil {
		response.ErrorFromErr(c, err)
		return
	}

	response.Success(c, gin.H{
		"token":        token,
		"expire_in":    86400, // 24小时
		"download_url": "/api/v1/hpa/download/" + token,
	})
}

// Download 下载文件
func (h *CourseHandler) Download(c *gin.Context) {
	token := c.Param("token")

	download, err := h.service.ValidateDownloadToken(token)
	if err != nil {
		response.ErrorFromErr(c, err)
		return
	}

	// 返回课程视频 URL（实际应用中应该返回文件流或临时 URL）
	response.Success(c, gin.H{
		"video_url": download.Course.VideoURL,
		"title":     download.Course.Title,
	})
}
