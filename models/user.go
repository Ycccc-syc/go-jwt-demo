package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户表模型
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"` // 用户名唯一
	Password string `gorm:"not null"`        // 存储加密后的密码
}

// BeforeCreate GORM 钩子，在创建用户前加密密码
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
