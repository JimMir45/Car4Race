package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"car4race/internal/model"
	"car4race/internal/repository"
)

type CourseService struct {
	repo     *repository.CourseRepository
	userRepo *repository.UserRepository
}

func NewCourseService(repo *repository.CourseRepository, userRepo *repository.UserRepository) *CourseService {
	return &CourseService{repo: repo, userRepo: userRepo}
}

// ========== Course ==========

// GetCourses 获取课程列表
func (s *CourseService) GetCourses(page, pageSize int, sortBy string) ([]model.Course, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.GetCourses(page, pageSize, sortBy)
}

// GetCourseBySlug 根据 slug 获取课程
func (s *CourseService) GetCourseBySlug(slug string) (*model.Course, error) {
	return s.repo.GetCourseBySlug(slug)
}

// GetCourseByID 根据 ID 获取课程
func (s *CourseService) GetCourseByID(id uint) (*model.Course, error) {
	return s.repo.GetCourseByID(id)
}

// CreateCourse 创建课程
func (s *CourseService) CreateCourse(course *model.Course) error {
	return s.repo.CreateCourse(course)
}

// UpdateCourse 更新课程
func (s *CourseService) UpdateCourse(course *model.Course) error {
	return s.repo.UpdateCourse(course)
}

// DeleteCourse 删除课程
func (s *CourseService) DeleteCourse(id uint) error {
	return s.repo.DeleteCourse(id)
}

// ========== Order ==========

// CreateOrder 创建订单
func (s *CourseService) CreateOrder(userID, courseID uint) (*model.Order, error) {
	// 检查课程是否存在
	course, err := s.repo.GetCourseByID(courseID)
	if err != nil {
		return nil, errors.New("课程不存在")
	}

	// 检查是否已购买
	purchased, _ := s.repo.CheckUserPurchased(userID, courseID)
	if purchased {
		return nil, errors.New("您已购买过该课程")
	}

	// 生成订单号
	orderNo := generateOrderNo()

	order := &model.Order{
		OrderNo:  orderNo,
		UserID:   userID,
		CourseID: courseID,
		Amount:   course.Price,
		Status:   "pending",
	}

	if err := s.repo.CreateOrder(order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetUserOrders 获取用户订单列表
func (s *CourseService) GetUserOrders(userID uint, page, pageSize int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	return s.repo.GetUserOrders(userID, page, pageSize)
}

// CheckUserPurchased 检查用户是否已购买课程
func (s *CourseService) CheckUserPurchased(userID, courseID uint) (bool, error) {
	return s.repo.CheckUserPurchased(userID, courseID)
}

// RedeemInviteCode 使用邀请码兑换课程
func (s *CourseService) RedeemInviteCode(userID uint, code string) (*model.Order, error) {
	// 获取邀请码
	inviteCode, err := s.repo.GetInviteCode(code)
	if err != nil {
		return nil, errors.New("邀请码无效")
	}

	// 检查邀请码是否可用
	if !inviteCode.IsActive {
		return nil, errors.New("邀请码已失效")
	}
	if inviteCode.UsedCount >= inviteCode.MaxUses {
		return nil, errors.New("邀请码已用完")
	}
	if inviteCode.ExpireAt != nil && inviteCode.ExpireAt.Before(time.Now()) {
		return nil, errors.New("邀请码已过期")
	}

	// 检查是否已购买
	purchased, _ := s.repo.CheckUserPurchased(userID, inviteCode.CourseID)
	if purchased {
		return nil, errors.New("您已拥有该课程")
	}

	// 创建订单（已支付状态）
	orderNo := generateOrderNo()
	now := time.Now()
	order := &model.Order{
		OrderNo:    orderNo,
		UserID:     userID,
		CourseID:   inviteCode.CourseID,
		Amount:     0, // 邀请码免费
		Status:     "paid",
		PayMethod:  "invite_code",
		PayTime:    &now,
		InviteCode: code,
	}

	if err := s.repo.CreateOrder(order); err != nil {
		return nil, err
	}

	// 更新邀请码使用次数
	_ = s.repo.IncrementInviteCodeUsage(inviteCode.ID)

	// 更新课程销量
	_ = s.repo.IncrementSalesCount(inviteCode.CourseID)

	return order, nil
}

// ========== Download ==========

// CreateDownloadToken 创建下载令牌
func (s *CourseService) CreateDownloadToken(userID, courseID uint) (string, error) {
	// 检查是否已购买
	purchased, _ := s.repo.CheckUserPurchased(userID, courseID)
	if !purchased {
		// 检查是否是 VIP
		user, err := s.userRepo.FindByID(userID)
		if err != nil || !user.CanDownload {
			return "", errors.New("您没有下载权限")
		}
	}

	// 检查今日下载次数
	count, _ := s.repo.CountUserDownloadsToday(userID)
	if count >= 3 {
		return "", errors.New("今日下载次数已用完")
	}

	// 生成 token
	token := generateToken()

	download := &model.Download{
		UserID:   userID,
		CourseID: courseID,
		Token:    token,
		ExpireAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.CreateDownload(download); err != nil {
		return "", err
	}

	return token, nil
}

// ValidateDownloadToken 验证下载令牌
func (s *CourseService) ValidateDownloadToken(token string) (*model.Download, error) {
	download, err := s.repo.GetDownloadByToken(token)
	if err != nil {
		return nil, errors.New("下载链接无效")
	}

	if download.Used {
		return nil, errors.New("下载链接已使用")
	}

	if download.ExpireAt.Before(time.Now()) {
		return nil, errors.New("下载链接已过期")
	}

	// 标记为已使用
	_ = s.repo.MarkDownloadUsed(download.ID)

	return download, nil
}

// ========== Admin ==========

// CreateInviteCode 创建邀请码
func (s *CourseService) CreateInviteCode(courseID uint, maxUses int, expireAt *time.Time) (*model.InviteCode, error) {
	code := generateInviteCode()

	inviteCode := &model.InviteCode{
		Code:     code,
		CourseID: courseID,
		MaxUses:  maxUses,
		ExpireAt: expireAt,
		IsActive: true,
	}

	if err := s.repo.CreateInviteCode(inviteCode); err != nil {
		return nil, err
	}

	return inviteCode, nil
}

// GetAllCourses 获取所有课程（管理后台）
func (s *CourseService) GetAllCourses(page, pageSize int) ([]model.Course, int64, error) {
	return s.repo.GetAllCourses(page, pageSize)
}

// GetAllInviteCodes 获取所有邀请码（管理后台）
func (s *CourseService) GetAllInviteCodes(page, pageSize int) ([]model.InviteCode, int64, error) {
	return s.repo.GetAllInviteCodes(page, pageSize)
}

// ========== Helper ==========

func generateOrderNo() string {
	return fmt.Sprintf("ORD%d%s", time.Now().UnixNano()/1e6, randomString(6))
}

func generateToken() string {
	return randomString(32)
}

func generateInviteCode() string {
	return "INV" + randomString(8)
}

func randomString(n int) string {
	b := make([]byte, n/2+1)
	rand.Read(b)
	return hex.EncodeToString(b)[:n]
}
