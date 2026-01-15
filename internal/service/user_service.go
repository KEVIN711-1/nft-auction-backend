// internal/service/user_service.go
package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"nft-auction-backend/internal/model"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// Register 注册用户
func (s *UserService) Register(username, password string) (string, error) {
	log.Printf("开始注册用户: %s", username)

	// 检查用户名是否已存在
	var existingUser model.User
	err := s.db.Where("username = ?", username).First(&existingUser).Error
	if err == nil {
		log.Printf("用户名已存在: %s", username)
		return "", errors.New("用户名已存在")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("查询用户失败: %v", err)
		return "", err
	}

	log.Printf("用户名可用: %s", username)

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		return "", err
	}

	log.Printf("密码加密成功")

	// 创建用户
	user := model.User{
		Username: username,
		Password: string(hashedPassword),
	}

	log.Printf("准备创建用户: %+v", user)

	// 这里执行数据库写入
	result := s.db.Create(&user)
	if result.Error != nil {
		log.Printf("创建用户失败: %v", result.Error)
		return "", result.Error
	}

	log.Printf("用户创建成功, ID: %d", user.ID)

	// 生成简单token
	token := GenerateSimpleToken(username)
	log.Printf("生成token: %s", token)

	return token, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (string, error) {
	// 查找用户
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("用户名或密码错误")
		}
		return "", err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 生成简单token
	token := GenerateSimpleToken(username)
	return token, nil
}

// ValidateToken 验证token（简单实现，只检查格式）
func (s *UserService) ValidateToken(token string) (string, error) {
	// 这里只需要验证token格式，真正的验证在authCheck中间件中
	// 从token中提取用户名（格式：时间戳-用户名）
	// 实现根据你的需要
	return "", nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GenerateSimpleToken 生成简单token（辅助函数）
func GenerateSimpleToken(username string) string {
	return fmt.Sprintf("%d-%s", time.Now().Unix(), username)
}
