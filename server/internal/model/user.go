package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表 - 两个子应用共用
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Phone     string         `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Avatar    string         `gorm:"size:500" json:"avatar"`
	Role      string         `gorm:"size:20;default:user" json:"role"` // user | vip | admin
	Status    string         `gorm:"size:20;default:active" json:"status"` // active | banned
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 会员相关（私域视频网站）
	VIPExpireAt   *time.Time `json:"vip_expire_at"`
	YearlySpend   float64    `gorm:"default:0" json:"yearly_spend"` // 年消费金额
	CanDownload   bool       `gorm:"default:false" json:"can_download"` // 是否有下载权限
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// VerificationCode 验证码表
type VerificationCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Phone     string    `gorm:"index;size:20;not null" json:"phone"`
	Code      string    `gorm:"size:10;not null" json:"code"`
	Purpose   string    `gorm:"size:20;default:login" json:"purpose"` // login | register
	ExpireAt  time.Time `json:"expire_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (VerificationCode) TableName() string {
	return "verification_codes"
}
