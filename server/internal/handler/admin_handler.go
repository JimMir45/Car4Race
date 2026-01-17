package handler

import (
	"net/http"
	"strconv"
	"time"

	"car4race/internal/model"
	"car4race/internal/service"
	"car4race/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	contentService *service.ContentService
	courseService  *service.CourseService
	fileService    *service.FileService
}

func NewAdminHandler(contentService *service.ContentService, courseService *service.CourseService, fileService *service.FileService) *AdminHandler {
	return &AdminHandler{
		contentService: contentService,
		courseService:  courseService,
		fileService:    fileService,
	}
}

// ========== Category ==========

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	Slug     string `json:"slug" binding:"required"`
	ParentID *uint  `json:"parent_id"`
	Sort     int    `json:"sort"`
}

// CreateCategory 创建分类
func (h *AdminHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	category := &model.Category{
		Name:     req.Name,
		Slug:     req.Slug,
		ParentID: req.ParentID,
		Sort:     req.Sort,
	}

	if err := h.contentService.CreateCategory(category); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(c, category)
}

// UpdateCategory 更新分类
func (h *AdminHandler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	category := &model.Category{
		ID:       uint(id),
		Name:     req.Name,
		Slug:     req.Slug,
		ParentID: req.ParentID,
		Sort:     req.Sort,
	}

	if err := h.contentService.UpdateCategory(category); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(c, category)
}

// DeleteCategory 删除分类
func (h *AdminHandler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.contentService.DeleteCategory(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// ========== Note ==========

// CreateNoteRequest 创建笔记请求
type CreateNoteRequest struct {
	CategoryID uint   `json:"category_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Slug       string `json:"slug" binding:"required"`
	Summary    string `json:"summary"`
	Content    string `json:"content"`
	CoverImage string `json:"cover_image"`
	IsPublic   bool   `json:"is_public"`
	Sort       int    `json:"sort"`
}

// CreateNote 创建笔记
func (h *AdminHandler) CreateNote(c *gin.Context) {
	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	note := &model.Note{
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Slug:       req.Slug,
		Summary:    req.Summary,
		Content:    req.Content,
		CoverImage: req.CoverImage,
		IsPublic:   req.IsPublic,
		Sort:       req.Sort,
	}

	if err := h.contentService.CreateNote(note); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(c, note)
}

// UpdateNote 更新笔记
func (h *AdminHandler) UpdateNote(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	note, err := h.contentService.GetNoteByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "笔记不存在")
		return
	}

	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	note.CategoryID = req.CategoryID
	note.Title = req.Title
	note.Slug = req.Slug
	note.Summary = req.Summary
	note.Content = req.Content
	note.CoverImage = req.CoverImage
	note.IsPublic = req.IsPublic
	note.Sort = req.Sort

	if err := h.contentService.UpdateNote(note); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(c, note)
}

// DeleteNote 删除笔记
func (h *AdminHandler) DeleteNote(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.contentService.DeleteNote(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// ========== Course ==========

// CreateCourseRequest 创建课程请求
type CreateCourseRequest struct {
	Title       string  `json:"title" binding:"required"`
	Slug        string  `json:"slug" binding:"required"`
	Description string  `json:"description"`
	CoverImage  string  `json:"cover_image"`
	Price       float64 `json:"price" binding:"required"`
	OrigPrice   float64 `json:"orig_price"`
	IsPublic    bool    `json:"is_public"`
	Sort        int     `json:"sort"`
}

// GetCourses 获取课程列表（管理后台）
func (h *AdminHandler) GetCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	courses, total, err := h.courseService.GetAllCourses(page, pageSize)
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

// CreateCourse 创建课程
func (h *AdminHandler) CreateCourse(c *gin.Context) {
	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	course := &model.Course{
		Title:       req.Title,
		Slug:        req.Slug,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Price:       req.Price,
		OrigPrice:   req.OrigPrice,
		IsPublic:    req.IsPublic,
		Sort:        req.Sort,
	}

	if err := h.courseService.CreateCourse(course); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(c, course)
}

// UpdateCourse 更新课程
func (h *AdminHandler) UpdateCourse(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	course, err := h.courseService.GetCourseByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "课程不存在")
		return
	}

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	course.Title = req.Title
	course.Slug = req.Slug
	course.Description = req.Description
	course.CoverImage = req.CoverImage
	course.Price = req.Price
	course.OrigPrice = req.OrigPrice
	course.IsPublic = req.IsPublic
	course.Sort = req.Sort

	if err := h.courseService.UpdateCourse(course); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(c, course)
}

// DeleteCourse 删除课程
func (h *AdminHandler) DeleteCourse(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.courseService.DeleteCourse(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// ========== InviteCode ==========

// CreateInviteCodeRequest 创建邀请码请求
type CreateInviteCodeRequest struct {
	CourseID uint   `json:"course_id" binding:"required"`
	MaxUses  int    `json:"max_uses"`
	ExpireAt string `json:"expire_at"` // RFC3339 格式
}

// GetInviteCodes 获取邀请码列表
func (h *AdminHandler) GetInviteCodes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	codes, total, err := h.courseService.GetAllInviteCodes(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取邀请码失败")
		return
	}

	response.Success(c, gin.H{
		"list":      codes,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateInviteCode 创建邀请码
func (h *AdminHandler) CreateInviteCode(c *gin.Context) {
	var req CreateInviteCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	maxUses := req.MaxUses
	if maxUses <= 0 {
		maxUses = 1
	}

	var expireAt *time.Time
	if req.ExpireAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpireAt)
		if err == nil {
			expireAt = &t
		}
	}

	code, err := h.courseService.CreateInviteCode(req.CourseID, maxUses, expireAt)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(c, code)
}

// ========== CourseFile ==========

// UploadCourseFile 上传课程文件
func (h *AdminHandler) UploadCourseFile(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "课程ID无效")
		return
	}

	fileType := c.DefaultPostForm("file_type", "resource") // intro | resource
	if fileType != "intro" && fileType != "resource" {
		response.Error(c, http.StatusBadRequest, "文件类型无效")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择文件")
		return
	}

	courseFile, err := h.fileService.UploadCourseFile(uint(courseID), fileType, file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, courseFile)
}

// GetCourseFiles 获取课程文件列表
func (h *AdminHandler) GetCourseFiles(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "课程ID无效")
		return
	}

	files, err := h.fileService.GetCourseFiles(uint(courseID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取文件列表失败")
		return
	}

	response.Success(c, files)
}

// DeleteCourseFile 删除课程文件
func (h *AdminHandler) DeleteCourseFile(c *gin.Context) {
	fileID, err := strconv.ParseUint(c.Param("fileId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "文件ID无效")
		return
	}

	if err := h.fileService.DeleteCourseFile(uint(fileID)); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}
