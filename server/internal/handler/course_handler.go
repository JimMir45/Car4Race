package handler

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"car4race/internal/model"
	"car4race/internal/service"
	"car4race/pkg/errcode"
	"car4race/pkg/response"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service     *service.CourseService
	fileService *service.FileService
}

func NewCourseHandler(service *service.CourseService, fileService *service.FileService) *CourseHandler {
	return &CourseHandler{service: service, fileService: fileService}
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

	course, err := h.service.GetCourseBySlugWithFiles(slug)
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

	// 获取课程介绍 Markdown 内容
	introContent, _ := h.fileService.GetCourseIntroContent(course.ID)

	// 分离 intro 和 resource 文件
	var introFiles, resourceFiles []model.CourseFile
	for _, f := range course.Files {
		if f.FileType == "intro" {
			introFiles = append(introFiles, f)
		} else {
			resourceFiles = append(resourceFiles, f)
		}
	}

	response.Success(c, gin.H{
		"course":         course,
		"purchased":      purchased,
		"intro_content":  introContent,
		"intro_files":    introFiles,
		"resource_files": resourceFiles,
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
	FileID   uint `json:"file_id"` // 可选，指定要下载的文件ID
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

	token, err := h.service.CreateDownloadToken(userID, req.CourseID, req.FileID)
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

	// 如果有指定文件，返回文件下载信息
	if download.FileID > 0 {
		// 方式1：使用预签名 URL（推荐，减轻后端压力）
		presignedURL, err := h.fileService.GetPresignedURL(download.FileID, 1*time.Hour)
		if err != nil {
			response.Error(c, http.StatusNotFound, "文件不存在")
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, presignedURL)
		return

		// 方式2：通过后端代理下载（备用，适合内网场景）
		// obj, fileName, fileSize, err := h.fileService.GetFileObject(download.FileID)
		// if err != nil {
		// 	response.Error(c, http.StatusNotFound, "文件不存在")
		// 	return
		// }
		// defer obj.Close()
		// c.Header("Content-Disposition", "attachment; filename="+fileName)
		// c.Header("Content-Type", "application/octet-stream")
		// c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
		// io.Copy(c.Writer, obj)
		// return
	}

	// 否则返回课程的资源文件列表
	files, err := h.fileService.GetCourseFiles(download.CourseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取文件列表失败")
		return
	}

	// 过滤出资源文件
	var resourceFiles []model.CourseFile
	for _, f := range files {
		if f.FileType == "resource" {
			resourceFiles = append(resourceFiles, f)
		}
	}

	response.Success(c, gin.H{
		"course_id": download.CourseID,
		"title":     download.Course.Title,
		"files":     resourceFiles,
	})
}

// Unused import placeholder - will be used if proxy download is enabled
var _ = io.Copy
