package repository

import (
	"car4race/internal/model"

	"gorm.io/gorm"
)

type ContentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

// ========== Category ==========

// GetAllCategories 获取所有分类（树形结构）
func (r *ContentRepository) GetAllCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("parent_id IS NULL").
		Preload("Children").
		Order("sort ASC").
		Find(&categories).Error
	return categories, err
}

// GetCategoryBySlug 根据 slug 获取分类
func (r *ContentRepository) GetCategoryBySlug(slug string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	return &category, err
}

// CreateCategory 创建分类
func (r *ContentRepository) CreateCategory(category *model.Category) error {
	return r.db.Create(category).Error
}

// UpdateCategory 更新分类
func (r *ContentRepository) UpdateCategory(category *model.Category) error {
	return r.db.Save(category).Error
}

// DeleteCategory 删除分类
func (r *ContentRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

// ========== Note ==========

// GetNotes 获取笔记列表
func (r *ContentRepository) GetNotes(categoryID uint, page, pageSize int) ([]model.Note, int64, error) {
	var notes []model.Note
	var total int64

	query := r.db.Model(&model.Note{}).Where("is_public = ?", true)
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	query.Count(&total)

	err := query.
		Preload("Category").
		Order("sort DESC, created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&notes).Error

	return notes, total, err
}

// GetNoteBySlug 根据 slug 获取笔记
func (r *ContentRepository) GetNoteBySlug(slug string) (*model.Note, error) {
	var note model.Note
	err := r.db.Preload("Category").Where("slug = ?", slug).First(&note).Error
	return &note, err
}

// GetNoteByID 根据 ID 获取笔记
func (r *ContentRepository) GetNoteByID(id uint) (*model.Note, error) {
	var note model.Note
	err := r.db.Preload("Category").First(&note, id).Error
	return &note, err
}

// CreateNote 创建笔记
func (r *ContentRepository) CreateNote(note *model.Note) error {
	return r.db.Create(note).Error
}

// UpdateNote 更新笔记
func (r *ContentRepository) UpdateNote(note *model.Note) error {
	return r.db.Save(note).Error
}

// DeleteNote 删除笔记
func (r *ContentRepository) DeleteNote(id uint) error {
	return r.db.Delete(&model.Note{}, id).Error
}

// IncrementNoteViewCount 增加笔记浏览次数
func (r *ContentRepository) IncrementNoteViewCount(id uint) error {
	return r.db.Model(&model.Note{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// ========== BrowseHistory ==========

// AddBrowseHistory 添加浏览记录
func (r *ContentRepository) AddBrowseHistory(userID, noteID uint) error {
	history := &model.BrowseHistory{
		UserID: userID,
		NoteID: noteID,
	}
	return r.db.Create(history).Error
}

// GetUserBrowseHistory 获取用户浏览记录
func (r *ContentRepository) GetUserBrowseHistory(userID uint, page, pageSize int) ([]model.BrowseHistory, int64, error) {
	var history []model.BrowseHistory
	var total int64

	query := r.db.Model(&model.BrowseHistory{}).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.
		Preload("Note").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&history).Error

	return history, total, err
}
