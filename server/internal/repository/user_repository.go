package repository

import (
	"time"

	"car4race/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByPhone 根据手机号查找用户
func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// SaveVerificationCode 保存验证码
func (r *UserRepository) SaveVerificationCode(code *model.VerificationCode) error {
	return r.db.Create(code).Error
}

// FindValidCode 查找有效的验证码
func (r *UserRepository) FindValidCode(phone, code, purpose string) (*model.VerificationCode, error) {
	var vc model.VerificationCode
	err := r.db.Where(
		"phone = ? AND code = ? AND purpose = ? AND used = ? AND expire_at > ?",
		phone, code, purpose, false, time.Now(),
	).First(&vc).Error
	if err != nil {
		return nil, err
	}
	return &vc, nil
}

// MarkCodeUsed 标记验证码已使用
func (r *UserRepository) MarkCodeUsed(id uint) error {
	return r.db.Model(&model.VerificationCode{}).Where("id = ?", id).Update("used", true).Error
}

// CountRecentCodes 统计最近发送的验证码数量（用于频率限制）
func (r *UserRepository) CountRecentCodes(phone string, duration time.Duration) (int64, error) {
	var count int64
	err := r.db.Model(&model.VerificationCode{}).
		Where("phone = ? AND created_at > ?", phone, time.Now().Add(-duration)).
		Count(&count).Error
	return count, err
}
