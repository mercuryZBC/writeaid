package dao

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"yuqueppbackend/service-base/db"
	"yuqueppbackend/service-base/models"
)

type UserDAO struct {
	DB *gorm.DB
}

// 校验密码是否匹配
func (dao *UserDAO) CheckPassword(user models.User) (bool, error) {
	tmpUser, err := dao.GetUserByEmail(user.Email)
	if err != nil {
		return false, err
	}
	if tmpUser.Password == user.Password {
		return true, nil
	} else {
		return false, nil
	}
}

// NewUserDAO 创建一个新的 UserDAO 实例
func NewUserDAO() *UserDAO {
	return &UserDAO{DB: db.GetDB()}
}

// CreateUser 创建一个新用户
func (dao *UserDAO) CreateUser(user models.User) error {
	user.RegisteredAt = time.Now()
	user.ExpiryAt = time.Now()
	user.LastLoginAt = time.Now()
	return dao.DB.Create(&user).Error
}

// GetUserByID 根据用户 ID 获取用户
func (dao *UserDAO) GetUserByID(id int64) (*models.User, error) {
	var user models.User
	if err := dao.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		return nil, err // 其他错误
	}
	return &user, nil
}

// GetUserByEmail 根据用户邮箱获取用户
func (dao *UserDAO) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := dao.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		return nil, err // 其他错误
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (dao *UserDAO) UpdateUser(user models.User) error {
	return dao.DB.Save(user).Error
}

func (dao *UserDAO) UpdateLastLoginTime(user models.User) error {
	user.LastLoginAt = time.Now()
	return dao.DB.Save(user).Error
}

// DeleteUserByID 根据 ID 删除用户
func (dao *UserDAO) DeleteUserByID(id int64) error {
	return dao.DB.Delete(&models.User{}, id).Error
}

// UpdateLastLogin 更新用户的最后登录时间
func (dao *UserDAO) UpdateLastLogin(userID int64) error {
	return dao.DB.Model(&models.User{}).Where("id = ?", userID).
		Update("last_login_at", time.Now()).Error
}
