package model

import (
	"time"

	"gorm.io/gorm"
)

// Course 课程表
type Course struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Slug        string         `gorm:"uniqueIndex;size:200;not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	CoverImage  string         `gorm:"size:500" json:"cover_image"`
	Price       float64        `gorm:"not null" json:"price"`
	OrigPrice   float64        `gorm:"default:0" json:"orig_price"`
	IntroPath   string         `gorm:"size:500" json:"intro_path"` // Markdown 介绍文件路径
	SalesCount  int            `gorm:"default:0" json:"sales_count"`
	IsPublic    bool           `gorm:"default:true" json:"is_public"`
	Sort        int            `gorm:"default:0" json:"sort"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Files []CourseFile `gorm:"foreignKey:CourseID" json:"files,omitempty"`
}

func (Course) TableName() string {
	return "hpa_courses"
}

// CourseFile 课程文件表
type CourseFile struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CourseID  uint      `gorm:"index;not null" json:"course_id"`
	FileType  string    `gorm:"size:20;not null" json:"file_type"` // intro | resource
	FileName  string    `gorm:"size:200;not null" json:"file_name"`
	FilePath  string    `gorm:"size:500;not null" json:"file_path"`
	FileSize  int64     `gorm:"default:0" json:"file_size"`
	Sort      int       `gorm:"default:0" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
}

func (CourseFile) TableName() string {
	return "hpa_course_files"
}

// Order 订单表
type Order struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	OrderNo     string     `gorm:"uniqueIndex;size:50;not null" json:"order_no"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	CourseID    uint       `gorm:"index;not null" json:"course_id"`
	Amount      float64    `gorm:"not null" json:"amount"`
	Status      string     `gorm:"size:20;default:pending" json:"status"` // pending | paid | refunded | cancelled
	PayMethod   string     `gorm:"size:20" json:"pay_method"`             // wechat | alipay | invite_code
	PayTime     *time.Time `json:"pay_time"`
	InviteCode  string     `gorm:"size:50" json:"invite_code"` // 使用的邀请码
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

func (Order) TableName() string {
	return "hpa_orders"
}

// InviteCode 邀请码表
type InviteCode struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Code      string     `gorm:"uniqueIndex;size:50;not null" json:"code"`
	CourseID  uint       `gorm:"index;not null" json:"course_id"`
	MaxUses   int        `gorm:"default:1" json:"max_uses"`
	UsedCount int        `gorm:"default:0" json:"used_count"`
	ExpireAt  *time.Time `json:"expire_at"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// 关联
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

func (InviteCode) TableName() string {
	return "hpa_invite_codes"
}

// Download 下载记录表
type Download struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	CourseID  uint      `gorm:"index;not null" json:"course_id"`
	FileID    uint      `gorm:"default:0" json:"file_id"` // 指定要下载的文件ID，0表示下载所有
	Token     string    `gorm:"uniqueIndex;size:100;not null" json:"token"`
	ExpireAt  time.Time `json:"expire_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

func (Download) TableName() string {
	return "hpa_downloads"
}
