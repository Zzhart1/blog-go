package service

import (
	"errors"
	"strings"

	"blog/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserExists = errors.New("用户名已存在")
	ErrBadLogin   = errors.New("账号或密码错误")
)

// RegisterUser 注册用户：密码 bcrypt 哈希后存储，用户名重复返回 ErrUserExists
func RegisterUser(username, password string) (*model.User, error) {
	username = strings.TrimSpace(username)

	var count int64
	if err := model.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{Username: username, Password: string(hash)}
	if err := model.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// LoginUser 登录：按用户名查库 + bcrypt 校验密码，失败统一返回 ErrBadLogin
func LoginUser(username, password string) (*model.User, error) {
	var user model.User
	err := model.DB.Where("username = ?", strings.TrimSpace(username)).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBadLogin
		}
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, ErrBadLogin
	}
	return &user, nil
}
