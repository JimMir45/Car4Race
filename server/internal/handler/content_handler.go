package handler

import (
	"net/http"
	"strconv"

	"car4race/internal/service"
	"car4race/pkg/response"

	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	service *service.ContentService
}

func NewContentHandler(service *service.ContentService) *ContentHandler {
	return &ContentHandler{service: service}
}

// GetCategories 获取所有分类
func (h *ContentHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取分类失败")
		return
	}
	response.Success(c, categories)
}

// GetNotes 获取笔记列表
func (h *ContentHandler) GetNotes(c *gin.Context) {
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	notes, total, err := h.service.GetNotes(uint(categoryID), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取笔记失败")
		return
	}

	response.Success(c, gin.H{
		"list":      notes,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetNote 获取笔记详情
func (h *ContentHandler) GetNote(c *gin.Context) {
	slug := c.Param("slug")
	userID := c.GetUint("user_id") // 可能为 0（未登录）

	note, err := h.service.GetNoteBySlug(slug, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "笔记不存在")
		return
	}

	response.Success(c, note)
}

// GetBrowseHistory 获取浏览记录
func (h *ContentHandler) GetBrowseHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	history, total, err := h.service.GetUserBrowseHistory(userID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取浏览记录失败")
		return
	}

	response.Success(c, gin.H{
		"list":      history,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
