package database

import (
	"go-jwt/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 替换成你的 MySQL 连接信息
	dsn := "root:123456@tcp(127.0.0.1:3306)/go_jwt_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 自动迁移，创建 users 表
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	DB = db

	// 初始化一个测试用户（可选）
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		db.Create(&models.User{
			Username: "admin",
			Password: "123456", // 会自动加密
		})
	}
}
