package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"car4race/internal/model"
	"car4race/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// SendVerificationCode 发送验证码
func (s *UserService) SendVerificationCode(phone string) error {
	// 检查频率限制：1分钟内只能发送1次
	count, err := s.repo.CountRecentCodes(phone, time.Minute)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("请稍后再试")
	}

	// 生成6位随机验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 保存验证码
	vc := &model.VerificationCode{
		Phone:    phone,
		Code:     code,
		Purpose:  "login",
		ExpireAt: time.Now().Add(5 * time.Minute),
	}
	if err := s.repo.SaveVerificationCode(vc); err != nil {
		return err
	}

	// TODO: 调用短信服务发送验证码
	// 开发环境直接打印
	fmt.Printf("[DEV] Verification code for %s: %s\n", phone, code)

	return nil
}

// Login 登录/注册
func (s *UserService) Login(phone, code string) (string, *model.User, error) {
	// 验证验证码
	vc, err := s.repo.FindValidCode(phone, code, "login")
	if err != nil {
		return "", nil, errors.New("验证码无效或已过期")
	}

	// 标记验证码已使用
	_ = s.repo.MarkCodeUsed(vc.ID)

	// 查找或创建用户
	user, err := s.repo.FindByPhone(phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 新用户，自动注册
			user = &model.User{
				Phone:    phone,
				Username: generateUsername(),
				Nickname: "用户" + phone[len(phone)-4:],
				Role:     "user",
				Status:   "active",
			}
			if err := s.repo.Create(user); err != nil {
				return "", nil, err
			}
		} else {
			return "", nil, err
		}
	}

	// 检查用户状态
	if user.Status == "banned" {
		return "", nil, errors.New("账号已被禁用")
	}

	// 生成 JWT
	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(user *model.User) error {
	return s.repo.Update(user)
}

// generateToken 生成 JWT token
func (s *UserService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"phone":    user.Phone,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(), // 7天过期
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// generateUsername 生成随机用户名
func generateUsername() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return "user_" + string(b)
}
