package repository

import (
	"time"

	"car4race/internal/model"

	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// ========== Course ==========

// GetCourses 获取课程列表
func (r *CourseRepository) GetCourses(page, pageSize int, sortBy string) ([]model.Course, int64, error) {
	var courses []model.Course
	var total int64

	query := r.db.Model(&model.Course{}).Where("is_public = ?", true)
	query.Count(&total)

	// 排序方式
	orderBy := "created_at DESC"
	switch sortBy {
	case "price_asc":
		orderBy = "price ASC"
	case "price_desc":
		orderBy = "price DESC"
	case "sales":
		orderBy = "sales_count DESC"
	case "newest":
		orderBy = "created_at DESC"
	}

	err := query.Order(orderBy).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&courses).Error

	return courses, total, err
}

// GetCourseBySlug 根据 slug 获取课程
func (r *CourseRepository) GetCourseBySlug(slug string) (*model.Course, error) {
	var course model.Course
	err := r.db.Where("slug = ?", slug).First(&course).Error
	return &course, err
}

// GetCourseByID 根据 ID 获取课程
func (r *CourseRepository) GetCourseByID(id uint) (*model.Course, error) {
	var course model.Course
	err := r.db.First(&course, id).Error
	return &course, err
}

// CreateCourse 创建课程
func (r *CourseRepository) CreateCourse(course *model.Course) error {
	return r.db.Create(course).Error
}

// UpdateCourse 更新课程
func (r *CourseRepository) UpdateCourse(course *model.Course) error {
	return r.db.Save(course).Error
}

// DeleteCourse 删除课程
func (r *CourseRepository) DeleteCourse(id uint) error {
	return r.db.Delete(&model.Course{}, id).Error
}

// IncrementSalesCount 增加销量
func (r *CourseRepository) IncrementSalesCount(id uint) error {
	return r.db.Model(&model.Course{}).Where("id = ?", id).
		UpdateColumn("sales_count", gorm.Expr("sales_count + 1")).Error
}

// ========== Order ==========

// CreateOrder 创建订单
func (r *CourseRepository) CreateOrder(order *model.Order) error {
	return r.db.Create(order).Error
}

// GetOrderByNo 根据订单号获取订单
func (r *CourseRepository) GetOrderByNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Course").Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

// GetUserOrders 获取用户订单列表
func (r *CourseRepository) GetUserOrders(userID uint, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.db.Model(&model.Order{}).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.
		Preload("Course").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders).Error

	return orders, total, err
}

// UpdateOrderStatus 更新订单状态
func (r *CourseRepository) UpdateOrderStatus(orderNo, status string) error {
	updates := map[string]interface{}{"status": status}
	if status == "paid" {
		now := time.Now()
		updates["pay_time"] = &now
	}
	return r.db.Model(&model.Order{}).Where("order_no = ?", orderNo).Updates(updates).Error
}

// CheckUserPurchased 检查用户是否已购买课程
func (r *CourseRepository) CheckUserPurchased(userID, courseID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Order{}).
		Where("user_id = ? AND course_id = ? AND status = ?", userID, courseID, "paid").
		Count(&count).Error
	return count > 0, err
}

// ========== InviteCode ==========

// GetInviteCode 获取邀请码
func (r *CourseRepository) GetInviteCode(code string) (*model.InviteCode, error) {
	var inviteCode model.InviteCode
	err := r.db.Preload("Course").Where("code = ?", code).First(&inviteCode).Error
	return &inviteCode, err
}

// CreateInviteCode 创建邀请码
func (r *CourseRepository) CreateInviteCode(code *model.InviteCode) error {
	return r.db.Create(code).Error
}

// IncrementInviteCodeUsage 增加邀请码使用次数
func (r *CourseRepository) IncrementInviteCodeUsage(id uint) error {
	return r.db.Model(&model.InviteCode{}).Where("id = ?", id).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}

// ========== Download ==========

// CreateDownload 创建下载记录
func (r *CourseRepository) CreateDownload(download *model.Download) error {
	return r.db.Create(download).Error
}

// GetDownloadByToken 根据 token 获取下载记录
func (r *CourseRepository) GetDownloadByToken(token string) (*model.Download, error) {
	var download model.Download
	err := r.db.Preload("Course").Where("token = ?", token).First(&download).Error
	return &download, err
}

// MarkDownloadUsed 标记下载已使用
func (r *CourseRepository) MarkDownloadUsed(id uint) error {
	return r.db.Model(&model.Download{}).Where("id = ?", id).Update("used", true).Error
}

// CountUserDownloadsToday 统计用户今日创建的下载令牌数
func (r *CourseRepository) CountUserDownloadsToday(userID uint) (int64, error) {
	var count int64
	// 获取今日零点（本地时区）
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	err := r.db.Model(&model.Download{}).
		Where("user_id = ? AND created_at >= ?", userID, today).
		Count(&count).Error
	return count, err
}

// GetAllCourses 获取所有课程（管理后台用）
func (r *CourseRepository) GetAllCourses(page, pageSize int) ([]model.Course, int64, error) {
	var courses []model.Course
	var total int64

	query := r.db.Model(&model.Course{})
	query.Count(&total)

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&courses).Error

	return courses, total, err
}

// GetAllInviteCodes 获取所有邀请码（管理后台用）
func (r *CourseRepository) GetAllInviteCodes(page, pageSize int) ([]model.InviteCode, int64, error) {
	var codes []model.InviteCode
	var total int64

	query := r.db.Model(&model.InviteCode{})
	query.Count(&total)

	err := query.Preload("Course").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&codes).Error

	return codes, total, err
}
