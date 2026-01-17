package service

import (
	"car4race/internal/model"
	"car4race/internal/repository"
)

type ContentService struct {
	repo *repository.ContentRepository
}

func NewContentService(repo *repository.ContentRepository) *ContentService {
	return &ContentService{repo: repo}
}

// ========== Category ==========

// GetAllCategories 获取所有分类
func (s *ContentService) GetAllCategories() ([]model.Category, error) {
	return s.repo.GetAllCategories()
}

// GetCategoryBySlug 根据 slug 获取分类
func (s *ContentService) GetCategoryBySlug(slug string) (*model.Category, error) {
	return s.repo.GetCategoryBySlug(slug)
}

// CreateCategory 创建分类
func (s *ContentService) CreateCategory(category *model.Category) error {
	return s.repo.CreateCategory(category)
}

// UpdateCategory 更新分类
func (s *ContentService) UpdateCategory(category *model.Category) error {
	return s.repo.UpdateCategory(category)
}

// DeleteCategory 删除分类
func (s *ContentService) DeleteCategory(id uint) error {
	return s.repo.DeleteCategory(id)
}

// ========== Note ==========

// GetNotes 获取笔记列表
func (s *ContentService) GetNotes(categoryID uint, page, pageSize int) ([]model.Note, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.GetNotes(categoryID, page, pageSize)
}

// GetNoteBySlug 根据 slug 获取笔记详情
func (s *ContentService) GetNoteBySlug(slug string, userID uint) (*model.Note, error) {
	note, err := s.repo.GetNoteBySlug(slug)
	if err != nil {
		return nil, err
	}

	// 增加浏览次数
	_ = s.repo.IncrementNoteViewCount(note.ID)

	// 记录浏览历史
	if userID > 0 {
		_ = s.repo.AddBrowseHistory(userID, note.ID)
	}

	return note, nil
}

// GetNoteByID 根据 ID 获取笔记
func (s *ContentService) GetNoteByID(id uint) (*model.Note, error) {
	return s.repo.GetNoteByID(id)
}

// CreateNote 创建笔记
func (s *ContentService) CreateNote(note *model.Note) error {
	return s.repo.CreateNote(note)
}

// UpdateNote 更新笔记
func (s *ContentService) UpdateNote(note *model.Note) error {
	return s.repo.UpdateNote(note)
}

// DeleteNote 删除笔记
func (s *ContentService) DeleteNote(id uint) error {
	return s.repo.DeleteNote(id)
}

// ========== BrowseHistory ==========

// GetUserBrowseHistory 获取用户浏览记录
func (s *ContentService) GetUserBrowseHistory(userID uint, page, pageSize int) ([]model.BrowseHistory, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.GetUserBrowseHistory(userID, page, pageSize)
}
