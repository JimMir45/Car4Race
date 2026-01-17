package model

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类表
type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Slug      string         `gorm:"uniqueIndex;size:50;not null" json:"slug"`
	ParentID  *uint          `gorm:"index" json:"parent_id"`
	Sort      int            `gorm:"default:0" json:"sort"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (Category) TableName() string {
	return "hpa_categories"
}

// Note 笔记表
type Note struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CategoryID uint           `gorm:"index;not null" json:"category_id"`
	Title      string         `gorm:"size:200;not null" json:"title"`
	Slug       string         `gorm:"uniqueIndex;size:200;not null" json:"slug"`
	Summary    string         `gorm:"size:500" json:"summary"`
	Content    string         `gorm:"type:text" json:"content"`
	CoverImage string         `gorm:"size:500" json:"cover_image"`
	ViewCount  int            `gorm:"default:0" json:"view_count"`
	IsPublic   bool           `gorm:"default:true" json:"is_public"`
	Sort       int            `gorm:"default:0" json:"sort"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Note) TableName() string {
	return "hpa_notes"
}

// BrowseHistory 浏览记录表
type BrowseHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	NoteID    uint      `gorm:"index;not null" json:"note_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Note Note `gorm:"foreignKey:NoteID" json:"note,omitempty"`
}

func (BrowseHistory) TableName() string {
	return "hpa_browse_history"
}
